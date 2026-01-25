package ticket

import (
	"context"
	"database/sql"
	"encoding/json"
	"example_shop/common/db"
	"example_shop/internal/gateway/http/dto"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// SearchTrain 车次查询用例：按站点、日期等过滤，支持游标分页与可选席别余票/最低价。
func (s *Service) SearchTrain(ctx context.Context, req dto.SearchTrainHTTPReq) Result {
	if isBreakerOpen() {
		return Result{Status: 503, Body: dto.BaseHTTPResp{Code: 503, Msg: "系统繁忙，请稍后再试"}}
	}

	dep, msg := validateTextField(req.DepartureStation, "departure_station")
	if msg != "" {
		return Result{Status: 400, Body: dto.BaseHTTPResp{Code: 400, Msg: msg}}
	}
	arr, msg := validateTextField(req.ArrivalStation, "arrival_station")
	if msg != "" {
		return Result{Status: 400, Body: dto.BaseHTTPResp{Code: 400, Msg: msg}}
	}
	seatType, msg := validateSeatType(req.SeatType)
	if msg != "" {
		return Result{Status: 400, Body: dto.BaseHTTPResp{Code: 400, Msg: msg}}
	}

	dayStart, dayEnd, err := parseTravelDate(req.TravelDate)
	if err != nil {
		return Result{Status: 400, Body: dto.BaseHTTPResp{Code: 400, Msg: "travel_date格式错误"}}
	}
	windowStart, windowEnd, msg := applyDepartTimeWindow(dayStart, dayEnd, req.DepartTimeStart, req.DepartTimeEnd)
	if msg != "" {
		return Result{Status: 400, Body: dto.BaseHTTPResp{Code: 400, Msg: msg}}
	}

	limit := req.Limit
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	var cursor trainQueryCursor
	if v, ok := decodeCursor(req.Cursor); ok {
		cursor = v
	}

	hot := hotKey(req.TravelDate, dep, arr, req.TrainType)
	var hotness int64 = 1
	if db.Rdb != nil {
		hotness, _ = db.Rdb.Incr(ctx, hot).Result()
		_ = db.Rdb.Expire(ctx, hot, 10*time.Minute).Err()
	}
	ttl := calcTTLByHotness(hotness)

	key := fmt.Sprintf(
		"ticket:search:%s:%s:%s:%s:%s:%s:%s:%t:%s:%s:%d:%s",
		req.TravelDate,
		dep,
		arr,
		seatType,
		strings.TrimSpace(req.TrainType),
		req.DepartTimeStart,
		req.DepartTimeEnd,
		req.HasTicket,
		strings.TrimSpace(req.Sort),
		strings.TrimSpace(req.Direction),
		limit,
		req.Cursor,
	)
	if db.Rdb != nil {
		if raw, err := db.Rdb.Get(ctx, key).Result(); err == nil && raw != "" {
			var cached dto.SearchTrainHTTPResp
			if jsonErr := json.Unmarshal([]byte(raw), &cached); jsonErr == nil {
				return Result{Status: 200, Body: cached}
			}
		}
	}

	type baseRow struct {
		TrainID       string
		TrainType     string
		FromSeq       uint32
		ToSeq         uint32
		DepartureTime sql.NullTime
		ArrivalTime   sql.NullTime
	}

	q := db.ReadDB().Table("train_station_passes dep").
		Select("dep.train_id AS train_id, t.train_type AS train_type, dep.sequence AS from_seq, arr.sequence AS to_seq, dep.departure_time AS departure_time, arr.arrival_time AS arrival_time").
		Joins("JOIN train_station_passes arr ON arr.train_id = dep.train_id AND arr.station_name = ? AND arr.deleted_at IS NULL", arr).
		Joins("JOIN train_infos t ON t.train_id = dep.train_id AND t.deleted_at IS NULL").
		Where("dep.station_name = ? AND dep.deleted_at IS NULL", dep).
		Where("dep.sequence < arr.sequence").
		Where("dep.departure_time >= ? AND dep.departure_time < ?", windowStart, windowEnd)
	if strings.TrimSpace(req.TrainType) != "" {
		q = q.Where("train_type = ?", strings.TrimSpace(req.TrainType))
	}
	if !cursor.DepartureTime.IsZero() && cursor.TrainID != "" {
		q = q.Where("(dep.departure_time > ?) OR (dep.departure_time = ? AND dep.train_id > ?)", cursor.DepartureTime, cursor.DepartureTime, cursor.TrainID)
	}
	q = q.Order("dep.departure_time ASC, dep.train_id ASC").Limit(limit + 1)

	var rows []baseRow
	if err := q.Scan(&rows).Error; err != nil {
		markBreakerErr()
		return Result{Status: 500, Body: dto.BaseHTTPResp{Code: 500, Msg: "查询失败: " + err.Error()}}
	}
	markBreakerOK()

	hasMore := false
	if len(rows) > limit {
		hasMore = true
		rows = rows[:limit]
	}

	items := make([]dto.TrainSearchItem, 0, len(rows))

	var wg sync.WaitGroup
	type remainResult struct {
		i        int
		remain   int64
		minPrice float64
		err      error
	}
	results := make(chan remainResult, len(rows))

	for i := range rows {
		r := rows[i]
		if !r.DepartureTime.Valid || !r.ArrivalTime.Valid {
			continue
		}
		rm := int(r.ArrivalTime.Time.Sub(r.DepartureTime.Time).Minutes())
		if rm < 0 {
			rm = 0
		}
		itemIdx := len(items)
		items = append(items, dto.TrainSearchItem{
			TrainID:            r.TrainID,
			TrainType:          r.TrainType,
			DepartureStation:   dep,
			ArrivalStation:     arr,
			DepartureTime:      r.DepartureTime.Time,
			ArrivalTime:        r.ArrivalTime.Time,
			RuntimeMinutes:     uint32(rm),
			SeatType:           seatType,
			SeatPrice:          0,
			RemainingSeatCount: 0,
		})

		if seatType == "" {
			continue
		}
		wg.Add(1)
		go func(idx int, trainID string, fromSeq, toSeq uint32) {
			defer wg.Done()
			remain, err1 := getRemainingSeatsSegment(ctx, trainID, seatType, dayStart, fromSeq, toSeq, ttl)
			if err1 != nil {
				results <- remainResult{i: idx, err: err1}
				return
			}
			minP, err2 := getMinPrice(ctx, trainID, seatType, ttl)
			if err2 != nil {
				results <- remainResult{i: idx, err: err2}
				return
			}
			results <- remainResult{i: idx, remain: remain, minPrice: minP}
		}(itemIdx, r.TrainID, r.FromSeq, r.ToSeq)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for rr := range results {
		if rr.err != nil {
			continue
		}
		if rr.i >= 0 && rr.i < len(items) {
			items[rr.i].SeatType = seatType
			items[rr.i].SeatPrice = rr.minPrice
			items[rr.i].RemainingSeatCount = rr.remain
		}
	}

	if req.HasTicket && seatType != "" {
		filtered := make([]dto.TrainSearchItem, 0, len(items))
		for _, it := range items {
			if it.RemainingSeatCount > 0 {
				filtered = append(filtered, it)
			}
		}
		items = filtered
	}

	switch strings.ToLower(strings.TrimSpace(req.Sort)) {
	case "price":
		if seatType != "" {
			sort.Slice(items, func(i, j int) bool { return items[i].SeatPrice < items[j].SeatPrice })
		}
	case "remain", "remaining":
		if seatType != "" {
			sort.Slice(items, func(i, j int) bool { return items[i].RemainingSeatCount < items[j].RemainingSeatCount })
		}
	case "", "time", "departure_time", "depart_time":
	default:
	}
	if strings.EqualFold(strings.TrimSpace(req.Direction), "desc") {
		for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
			items[i], items[j] = items[j], items[i]
		}
	}

	nextCursor := ""
	if hasMore && len(rows) > 0 {
		for i := len(rows) - 1; i >= 0; i-- {
			if rows[i].DepartureTime.Valid {
				nextCursor = encodeCursor(trainQueryCursor{DepartureTime: rows[i].DepartureTime.Time, TrainID: rows[i].TrainID})
				break
			}
		}
	}

	resp := dto.SearchTrainHTTPResp{
		Code:       200,
		Msg:        "success",
		Items:      items,
		PrevCursor: "",
		NextCursor: nextCursor,
	}

	if db.Rdb != nil {
		if b, err := json.Marshal(resp); err == nil {
			_ = db.Rdb.Set(ctx, key, string(b), ttl).Err()
		}
	}
	return Result{Status: 200, Body: resp}
}
