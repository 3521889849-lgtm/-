package ticket

import (
	"context"
	"example_shop/common/db"
	model2 "example_shop/internal/model"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type breakerState struct {
	mu        sync.Mutex
	openUntil time.Time
	errCount  int
}

var searchBreaker breakerState

func isBreakerOpen() bool {
	searchBreaker.mu.Lock()
	defer searchBreaker.mu.Unlock()
	return time.Now().Before(searchBreaker.openUntil)
}

func markBreakerErr() {
	searchBreaker.mu.Lock()
	defer searchBreaker.mu.Unlock()
	// 轻量熔断策略（用于保护 MySQL）：
	// - 连续出现多次查询错误后，短时间“开路”，让请求快速失败/降级
	// - 这不是完整的熔断器实现（缺少半开探测/按错误类型分类等），但足以避免错误风暴
	searchBreaker.errCount++
	if searchBreaker.errCount >= 5 {
		searchBreaker.openUntil = time.Now().Add(3 * time.Second)
		searchBreaker.errCount = 0
	}
}

func markBreakerOK() {
	searchBreaker.mu.Lock()
	defer searchBreaker.mu.Unlock()
	searchBreaker.errCount = 0
}

func hotKey(dateStr, dep, arr, trainType string) string {
	tt := strings.TrimSpace(trainType)
	if tt == "" {
		tt = "all"
	}
	return fmt.Sprintf("ticket:hot:%s:%s:%s:%s", dateStr, dep, arr, tt)
}

func calcTTLByHotness(n int64) time.Duration {
	// 热度越高（同一搜索条件被访问越频繁），TTL 越短：
	// - 热点 TTL 短：降低“余票/价格变化”的缓存滞后时间
	// - 冷门 TTL 长：减少 DB 压力
	if n >= 100 {
		return 1 * time.Minute
	}
	if n >= 20 {
		return 5 * time.Minute
	}
	return 30 * time.Minute
}

func redisRemainKey(trainID, seatType string, departureDate time.Time) string {
	return fmt.Sprintf("ticket:remain:%s:%s:%s", trainID, seatType, departureDate.Format("2006-01-02"))
}

func redisRemainKeySegment(trainID, seatType string, departureDate time.Time, fromSeq, toSeq uint32) string {
	return fmt.Sprintf("ticket:remainseg:%s:%s:%s:%d:%d", trainID, seatType, departureDate.Format("2006-01-02"), fromSeq, toSeq)
}

func redisMinPriceKey(trainID, seatType string) string {
	return fmt.Sprintf("ticket:price:min:%s:%s", trainID, seatType)
}

func getRemainingSeats(ctx context.Context, trainID, seatType string, departureDate time.Time, ttl time.Duration) (int64, error) {
	if strings.TrimSpace(seatType) == "" {
		return 0, nil
	}

	key := redisRemainKey(trainID, seatType, departureDate)
	if db.Rdb != nil {
		val, err := db.Rdb.Get(ctx, key).Result()
		if err == nil {
			n, parseErr := strconv.ParseInt(val, 10, 64)
			if parseErr == nil {
				return n, nil
			}
		}
	}

	var t model2.TrainInfo
	if err := db.ReadDB().Select("train_id, departure_station, arrival_station").Where("train_id = ?", trainID).First(&t).Error; err != nil {
		return 0, err
	}
	var depStop model2.TrainStationPass
	var arrStop model2.TrainStationPass
	seqOK := true
	if err := db.ReadDB().Select("sequence").Where("train_id = ? AND station_name = ? AND deleted_at IS NULL", t.ID, t.DepartureStation).First(&depStop).Error; err != nil {
		seqOK = false
	}
	if err := db.ReadDB().Select("sequence").Where("train_id = ? AND station_name = ? AND deleted_at IS NULL", t.ID, t.ArrivalStation).First(&arrStop).Error; err != nil {
		seqOK = false
	}
	if depStop.Sequence == 0 || arrStop.Sequence == 0 || depStop.Sequence >= arrStop.Sequence {
		seqOK = false
	}

	var count int64
	now := time.Now()
	if seqOK {
		if err := db.ReadDB().Table("seat_infos si").
			Where("si.train_id = ? AND si.seat_type = ? AND si.status = ? AND si.deleted_at IS NULL", trainID, seatType, "AVAILABLE").
			Where(
				"NOT EXISTS (SELECT 1 FROM seat_segment_occupancies o WHERE o.deleted_at IS NULL AND o.train_id = si.train_id AND o.seat_id = si.seat_id AND o.status IN ('SOLD','LOCKED') AND (o.status <> 'LOCKED' OR o.lock_expire_time > ?) AND o.from_seq < ? AND o.to_seq > ?)",
				now, arrStop.Sequence, depStop.Sequence,
			).
			Count(&count).Error; err != nil {
			return 0, err
		}
	} else {
		q := db.ReadDB().Model(&model2.SeatInfo{}).
			Where(
				"train_id = ? AND seat_type = ? AND ((status = ?) OR (status = ? AND lock_expire_time IS NOT NULL AND lock_expire_time < ?))",
				trainID, seatType, "AVAILABLE", "LOCKED", now,
			)
		if err := q.Count(&count).Error; err != nil {
			return 0, err
		}
	}

	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	if db.Rdb != nil {
		_ = db.Rdb.Set(ctx, key, strconv.FormatInt(count, 10), ttl).Err()
	}
	return count, nil
}

func getRemainingSeatsSegment(ctx context.Context, trainID, seatType string, departureDate time.Time, fromSeq, toSeq uint32, ttl time.Duration) (int64, error) {
	if strings.TrimSpace(seatType) == "" {
		return 0, nil
	}
	if fromSeq == 0 || toSeq == 0 || fromSeq >= toSeq {
		return 0, nil
	}

	key := redisRemainKeySegment(trainID, seatType, departureDate, fromSeq, toSeq)
	if db.Rdb != nil {
		val, err := db.Rdb.Get(ctx, key).Result()
		if err == nil {
			n, parseErr := strconv.ParseInt(val, 10, 64)
			if parseErr == nil {
				return n, nil
			}
		}
	}

	var count int64
	now := time.Now()

	// 用“区间重叠”规则过滤掉被 SOLD/未过期 LOCKED 占用的座位。
	// 说明：
	// - 余票计算使用 COUNT，不分配具体 seat_id，适合查询场景
	// - 与下单分配座位的 SQL 逻辑保持一致，避免“展示有票但下单无票”的偏差过大
	if err := db.ReadDB().Table("seat_infos si").
		Where("si.train_id = ? AND si.seat_type = ? AND si.status = ? AND si.deleted_at IS NULL", trainID, seatType, "AVAILABLE").
		Where(
			"NOT EXISTS (SELECT 1 FROM seat_segment_occupancies o WHERE o.deleted_at IS NULL AND o.train_id = si.train_id AND o.seat_id = si.seat_id AND o.status IN ('SOLD','LOCKED') AND (o.status <> 'LOCKED' OR o.lock_expire_time > ?) AND o.from_seq < ? AND o.to_seq > ?)",
			now, toSeq, fromSeq,
		).
		Count(&count).Error; err != nil {
		return 0, err
	}

	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	if db.Rdb != nil {
		_ = db.Rdb.Set(ctx, key, strconv.FormatInt(count, 10), ttl).Err()
	}
	return count, nil
}

// GetRemainingSeats 获取指定车次+席别在某天的余票数（带缓存）。
func GetRemainingSeats(ctx context.Context, trainID, seatType string, departureDate time.Time, ttl time.Duration) (int64, error) {
	return getRemainingSeats(ctx, trainID, seatType, departureDate, ttl)
}

func getMinPrice(ctx context.Context, trainID, seatType string, ttl time.Duration) (float64, error) {
	if strings.TrimSpace(seatType) == "" {
		return 0, nil
	}
	key := redisMinPriceKey(trainID, seatType)
	if db.Rdb != nil {
		val, err := db.Rdb.Get(ctx, key).Result()
		if err == nil {
			f, parseErr := strconv.ParseFloat(val, 64)
			if parseErr == nil {
				return f, nil
			}
		}
	}

	var price float64
	if err := db.ReadDB().Model(&model2.SeatInfo{}).
		Select("MIN(seat_price)").
		Where("train_id = ? AND seat_type = ?", trainID, seatType).
		Scan(&price).Error; err != nil {
		return 0, err
	}

	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	if db.Rdb != nil {
		_ = db.Rdb.Set(ctx, key, strconv.FormatFloat(price, 'f', -1, 64), ttl).Err()
	}
	return price, nil
}
