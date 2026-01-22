package dto

import "time"

// SearchTrainHTTPReq 车次查询 HTTP 请求参数（通过 URL Query 传参）。
//
// 对应接口：GET /api/v1/train/search
//
// 说明：
// - required：必填参数，由 Hertz 的 BindAndValidate 校验
// - 时间/日期格式：travel_date 建议使用 YYYY-MM-DD（如 2026-01-17）
type SearchTrainHTTPReq struct {
	DepartureStation string `query:"departure_station,required"` // 出发站（站点名/站点编码，取决于后端实现）
	ArrivalStation   string `query:"arrival_station,required"`   // 到达站（站点名/站点编码，取决于后端实现）
	TravelDate       string `query:"travel_date,required"`       // 出行日期（YYYY-MM-DD）
	TrainType        string `query:"train_type"`                 // 车次类型过滤（如 G/D/K/T 等；为空表示不过滤）
	SeatType         string `query:"seat_type"`                  // 席别过滤（如 二等座/一等座/硬卧 等；为空表示不过滤）
	DepartTimeStart  string `query:"depart_time_start"`          // 出发时间起（HH:MM，24小时制；为空表示不限）
	DepartTimeEnd    string `query:"depart_time_end"`            // 出发时间止（HH:MM，24小时制；为空表示不限）
	Sort             string `query:"sort"`                       // 排序字段（由后端定义，如 price/time）
	HasTicket        bool   `query:"has_ticket"`                 // 是否只看“有票”的车次（true=仅返回余票>0）
	Direction        string `query:"direction"`                  // 排序方向（asc/desc；由后端定义）
	Cursor           string `query:"cursor"`                     // 游标分页游标（上一页/下一页返回）
	Limit            int    `query:"limit"`                      // 每页条数（建议 1~50；不传由后端设默认）
}

// TrainSearchItem 车次查询结果中的单条车次信息。
type TrainSearchItem struct {
	TrainID            string    `json:"train_id"`             // 车次ID（后端唯一标识）
	TrainType          string    `json:"train_type"`           // 车次类型（如 G/D/K/T 等）
	DepartureStation   string    `json:"departure_station"`    // 出发站
	ArrivalStation     string    `json:"arrival_station"`      // 到达站
	DepartureTime      time.Time `json:"departure_time"`       // 出发时间（RFC3339 序列化）
	ArrivalTime        time.Time `json:"arrival_time"`         // 到达时间（RFC3339 序列化）
	RuntimeMinutes     uint32    `json:"runtime_minutes"`      // 运行时长（分钟）
	SeatType           string    `json:"seat_type"`            // 当前返回的席别（与请求 seat_type/后端聚合策略相关）
	SeatPrice          float64   `json:"seat_price"`           // 票价（单位：元）
	RemainingSeatCount int64     `json:"remaining_seat_count"` // 余票数量
}

// SearchTrainHTTPResp 车次查询 HTTP 响应。
//
// 分页说明：
// - PrevCursor/NextCursor：用于游标分页；不支持时可能为空
type SearchTrainHTTPResp struct {
	Code       int32             `json:"code"`        // 业务状态码（与 BaseHTTPResp 语义一致）
	Msg        string            `json:"msg"`         // 状态说明/错误信息
	Items      []TrainSearchItem `json:"items"`       // 车次列表
	PrevCursor string            `json:"prev_cursor"` // 上一页游标
	NextCursor string            `json:"next_cursor"` // 下一页游标
}
