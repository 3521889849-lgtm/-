// Package ticket 是网关内的“票务业务模块”（Application/UseCase）。
//
// 目标：
// - 将票务查询类逻辑（站点联想、车次查询、缓存/熔断策略）从 HTTP Handler 下沉出来
// - 让 handler 只做“参数绑定 -> 调用用例 -> 写响应”，避免业务代码散落在网关各处
package ticket

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"example_shop/common/cache"
	"example_shop/common/db"
	"example_shop/internal/gateway/http/dto"
	"example_shop/internal/ticket_service/model"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Service struct{}

// New 创建票务业务模块实例。
func New() *Service {
	return &Service{}
}

// Result 是网关业务模块返回给 HTTP Handler 的统一结果。
// Status 为 HTTP 状态码，Body 为要输出的 JSON 结构体。
type Result struct {
	Status int
	Body   any
}

var stationListCache = cache.NewTTLCache[[]string]()

func normalizeStationKeyword(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "")
	return strings.ToLower(s)
}

func stationAlias(keyword string) string {
	switch normalizeStationKeyword(keyword) {
	case "沪", "sh", "shanghai", "shangh":
		return "上海"
	case "京", "bj", "beijing":
		return "北京"
	case "广", "gz", "guangzhou":
		return "广州"
	case "深", "sz", "shenzhen":
		return "深圳"
	default:
		return ""
	}
}

func loadStations(ctx context.Context) ([]string, error) {
	if v, ok := stationListCache.Get("all"); ok && len(v) > 0 {
		return v, nil
	}

	type row struct {
		Name string
	}
	stations := make(map[string]struct{})

	var all []row
	if err := db.ReadDB().Model(&model.TrainStationPass{}).Select("DISTINCT station_name AS name").Scan(&all).Error; err != nil {
		return nil, err
	}
	for _, r := range all {
		name := strings.TrimSpace(r.Name)
		if name != "" {
			stations[name] = struct{}{}
		}
	}

	items := make([]string, 0, len(stations))
	for s := range stations {
		items = append(items, s)
	}
	sort.Strings(items)

	stationListCache.Set("all", items, 10*time.Minute)
	return items, nil
}

// StationSuggest 是“站点联想”用例：返回最多 limit 条候选站点名称。
func (s *Service) StationSuggest(ctx context.Context, req dto.StationSuggestHTTPReq) Result {
	limit := req.Limit
	if limit <= 0 || limit > 20 {
		limit = 10
	}

	kw := normalizeStationKeyword(req.Keyword)
	if len(kw) > 64 {
		return Result{Status: 400, Body: dto.BaseHTTPResp{Code: 400, Msg: "keyword过长"}}
	}

	stations, err := loadStations(ctx)
	if err != nil {
		return Result{Status: 500, Body: dto.BaseHTTPResp{Code: 500, Msg: "加载站点失败: " + err.Error()}}
	}

	alias := stationAlias(req.Keyword)
	res := make([]string, 0, limit)
	seen := make(map[string]struct{}, limit)

	if alias != "" {
		seen[alias] = struct{}{}
		res = append(res, alias)
	}

	for _, st := range stations {
		if len(res) >= limit {
			break
		}
		if _, ok := seen[st]; ok {
			continue
		}
		if strings.Contains(normalizeStationKeyword(st), kw) || strings.Contains(st, req.Keyword) {
			seen[st] = struct{}{}
			res = append(res, st)
		}
	}

	return Result{Status: 200, Body: dto.StationSuggestHTTPResp{Code: 200, Msg: "success", Items: res}}
}

type trainQueryCursor struct {
	DepartureTime time.Time `json:"departure_time"`
	TrainID       string    `json:"train_id"`
}

var safeTextRe = regexp.MustCompile(`^[\p{Han}A-Za-z0-9_\-]+$`)

type breakerState struct {
	mu        sync.Mutex
	openUntil time.Time
	errCount  int
}

var searchBreaker breakerState

func encodeCursor(c trainQueryCursor) string {
	b, _ := json.Marshal(c)
	return base64.RawURLEncoding.EncodeToString(b)
}

func decodeCursor(s string) (trainQueryCursor, bool) {
	if strings.TrimSpace(s) == "" {
		return trainQueryCursor{}, false
	}
	raw, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return trainQueryCursor{}, false
	}
	var c trainQueryCursor
	if err := json.Unmarshal(raw, &c); err != nil {
		return trainQueryCursor{}, false
	}
	if c.TrainID == "" || c.DepartureTime.IsZero() {
		return trainQueryCursor{}, false
	}
	return c, true
}

// parseTravelDate 解析 YYYY-MM-DD 为 [dayStart, dayEnd)。
func parseTravelDate(dateStr string) (time.Time, time.Time, error) {
	d, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	start := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
	end := start.Add(24 * time.Hour)
	return start, end, nil
}

// ParseTravelDate 解析 YYYY-MM-DD 为 [dayStart, dayEnd)。
func ParseTravelDate(dateStr string) (time.Time, time.Time, error) {
	return parseTravelDate(dateStr)
}

// parseTimeOfDay 解析 HH:mm 为一天内的偏移量。
func parseTimeOfDay(s string) (time.Duration, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	t, err := time.Parse("15:04", s)
	if err != nil {
		return 0, false
	}
	return time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute, true
}

func applyDepartTimeWindow(dayStart, dayEnd time.Time, startStr, endStr string) (time.Time, time.Time, string) {
	start := dayStart
	end := dayEnd
	if d, ok := parseTimeOfDay(startStr); ok {
		start = dayStart.Add(d)
	}
	if d, ok := parseTimeOfDay(endStr); ok {
		end = dayStart.Add(d)
		if end.Before(start) {
			return time.Time{}, time.Time{}, "出发时段不合法"
		}
	}
	return start, end, ""
}

func validateTextField(val, fieldName string) (string, string) {
	v := strings.TrimSpace(val)
	if v == "" {
		return "", fieldName + "不能为空"
	}
	if len(v) > 64 {
		return "", fieldName + "过长"
	}
	if !safeTextRe.MatchString(strings.ReplaceAll(v, " ", "")) {
		return "", fieldName + "包含非法字符"
	}
	return v, ""
}

func validateSeatType(seatType string) (string, string) {
	v := strings.TrimSpace(seatType)
	if v == "" {
		return "", ""
	}
	allow := map[string]struct{}{
		"硬座":  {},
		"二等座": {},
		"一等座": {},
		"商务座": {},
		"硬卧":  {},
		"软卧":  {},
	}
	if _, ok := allow[v]; !ok {
		return "", "座位类型无效"
	}
	return v, ""
}

// ValidateSeatType 校验座位类型并返回规范化后的值（白名单）。
func ValidateSeatType(seatType string) (string, string) {
	return validateSeatType(seatType)
}

func hotKey(dateStr, dep, arr, trainType string) string {
	tt := strings.TrimSpace(trainType)
	if tt == "" {
		tt = "all"
	}
	return fmt.Sprintf("ticket:hot:%s:%s:%s:%s", dateStr, dep, arr, tt)
}

func calcTTLByHotness(n int64) time.Duration {
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

	var t model.TrainInfo
	if err := db.ReadDB().Select("train_id, departure_station, arrival_station").Where("train_id = ?", trainID).First(&t).Error; err != nil {
		return 0, err
	}
	var depStop model.TrainStationPass
	var arrStop model.TrainStationPass
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
		q := db.ReadDB().Model(&model.SeatInfo{}).
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
	if err := db.ReadDB().Model(&model.SeatInfo{}).
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

func isBreakerOpen() bool {
	searchBreaker.mu.Lock()
	defer searchBreaker.mu.Unlock()
	return time.Now().Before(searchBreaker.openUntil)
}

func markBreakerErr() {
	searchBreaker.mu.Lock()
	defer searchBreaker.mu.Unlock()
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

// SearchTrain 是“车次查询”用例：按站点、日期等过滤，支持游标分页与可选席别余票/最低价。
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
		TrainID        string
		TrainType      string
		DepartureTime  time.Time
		ArrivalTime    time.Time
		RuntimeMinutes int
	}

	q := db.ReadDB().Model(&model.TrainInfo{}).
		Select("train_id, train_type, departure_time, arrival_time, TIMESTAMPDIFF(MINUTE, departure_time, arrival_time) AS runtime_minutes").
		Where("departure_station = ? AND arrival_station = ? AND departure_time >= ? AND departure_time < ?", dep, arr, windowStart, windowEnd)
	if strings.TrimSpace(req.TrainType) != "" {
		q = q.Where("train_type = ?", strings.TrimSpace(req.TrainType))
	}
	if !cursor.DepartureTime.IsZero() && cursor.TrainID != "" {
		q = q.Where("(departure_time > ?) OR (departure_time = ? AND train_id > ?)", cursor.DepartureTime, cursor.DepartureTime, cursor.TrainID)
	}
	q = q.Order("departure_time ASC, train_id ASC").Limit(limit + 1)

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
		items = append(items, dto.TrainSearchItem{
			TrainID:            r.TrainID,
			TrainType:          r.TrainType,
			DepartureStation:   dep,
			ArrivalStation:     arr,
			DepartureTime:      r.DepartureTime,
			ArrivalTime:        r.ArrivalTime,
			RuntimeMinutes:     uint32(r.RuntimeMinutes),
			SeatType:           seatType,
			SeatPrice:          0,
			RemainingSeatCount: 0,
		})

		if seatType == "" {
			continue
		}
		wg.Add(1)
		go func(idx int, trainID string) {
			defer wg.Done()
			remain, err1 := getRemainingSeats(ctx, trainID, seatType, dayStart, ttl)
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
		}(i, r.TrainID)
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
		last := rows[len(rows)-1]
		nextCursor = encodeCursor(trainQueryCursor{DepartureTime: last.DepartureTime, TrainID: last.TrainID})
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
