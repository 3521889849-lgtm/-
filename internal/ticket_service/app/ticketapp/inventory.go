package ticketapp

import (
	"database/sql"
	"errors"
	model2 "example_shop/internal/model"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AllocateSeats 负责在事务内按席别分配可用座席，并返回分配到的 seat 列表与总价。
//
// 该函数代表“票务域对外的座席分配边界能力”，订单域通过它来完成选座/锁座的核心步骤，
// 从而避免订单域直接实现复杂的座席可用性判定 SQL。
func AllocateSeats(tx *gorm.DB, trainID string, need map[string]int, fromSeq, toSeq uint32, lockExpire time.Time) ([]*model2.SeatInfo, float64, error) {
	if tx == nil {
		return nil, 0, errors.New("tx is nil")
	}
	if fromSeq == 0 || toSeq == 0 || fromSeq >= toSeq {
		return nil, 0, errors.New("invalid segment")
	}

	allocated := make([]*model2.SeatInfo, 0, 8)
	total := 0.0
	now := time.Now()

	// 重要：锁座过期的“回收”
	// - 这里是“尽力而为”的清理：在分配座位前，把已过期的 LOCKED 占用改成 CANCELLED，减少后续冲突概率
	// - 为什么不是强一致的“定时任务唯一回收”：
	//   - 项目内既有后台清理 job，也允许在热点路径做轻量回收，避免大量遗留 LOCKED 影响购票体验
	// - 注意：这里不返回错误，避免“清理失败导致无法下单”
	_ = tx.Model(&model2.SeatSegmentOccupancy{}).
		Where("status = ? AND lock_expire_time IS NOT NULL AND lock_expire_time < ?", "LOCKED", now).
		Updates(map[string]interface{}{"status": "CANCELLED", "lock_expire_time": sql.NullTime{Valid: false}}).Error

	for seatType, cnt := range need {
		if cnt <= 0 {
			continue
		}
		var seats []model2.SeatInfo

		// 座位可用性的核心判定（区间冲突过滤）：
		// - 座位“可用”不仅取决于 seat_infos.status，还取决于 seat_segment_occupancies 的区间占用
		// - 冲突条件：已有占用区间与新购票区间重叠，即：
		//     existing.from_seq < new.to_seq AND existing.to_seq > new.from_seq
		// - 对 LOCKED 还需要判断是否过期：未过期的 LOCKED 等价 SOLD（都应该阻塞分配）
		//
		// 并发安全策略（简化版）：
		// - 对候选 seat_infos 记录加 FOR UPDATE（GORM clause.Locking）
		// - 让并发事务在“同一批 seat 行”上互斥，降低同一座位被同时分配的概率
		if err := tx.Table("seat_infos si").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("si.train_id = ? AND si.seat_type = ? AND si.status = ? AND si.deleted_at IS NULL", trainID, seatType, "AVAILABLE").
			Where(
				"NOT EXISTS (SELECT 1 FROM seat_segment_occupancies o WHERE o.deleted_at IS NULL AND o.train_id = si.train_id AND o.seat_id = si.seat_id AND o.status IN ('SOLD','LOCKED') AND (o.status <> 'LOCKED' OR o.lock_expire_time > ?) AND o.from_seq < ? AND o.to_seq > ?)",
				now, toSeq, fromSeq,
			).
			Order("si.seat_id ASC").
			Limit(cnt).
			Scan(&seats).Error; err != nil {
			return nil, 0, err
		}
		if len(seats) < cnt {
			return nil, 0, fmt.Errorf("座位余量不足: %s", seatType)
		}
		for i := range seats {
			allocated = append(allocated, &seats[i])
			total += seats[i].SeatPrice
		}
	}

	return allocated, total, nil
}
