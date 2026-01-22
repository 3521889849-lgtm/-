package dto

// TicketRemainWSReq 余票 WebSocket 订阅请求参数（通过 URL Query 传参）。
//
// 对应接口：GET /api/v1/ticket/ws
//
// 说明：
// - 该接口会升级为 WebSocket 连接
// - 建议 travel_date 使用 YYYY-MM-DD（如 2026-01-17）
type TicketRemainWSReq struct {
	TrainID    string `query:"train_id,required"`    // 车次ID（后端唯一标识）
	SeatType   string `query:"seat_type,required"`   // 席别（用于查询指定席别余票）
	TravelDate string `query:"travel_date,required"` // 出行日期（YYYY-MM-DD）
}
