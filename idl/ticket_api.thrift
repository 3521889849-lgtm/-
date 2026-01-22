// ticket_api.thrift
//
// 说明：
// - 票务域 RPC 接口定义（车次详情/余票/最低价等）。
// - 该文件只放“票务域”的请求/响应与服务，不与用户/订单混写。
namespace go ticketapi

include "common.thrift"

// SeatTypeRemain 某席别的余票与最低价信息。
struct SeatTypeRemain {
    1: string seat_type
    2: i32 remaining
    3: double min_price
}

// GetTrainDetailReq 获取车次详情请求。
struct GetTrainDetailReq {
    1: string train_id
    2: optional string departure_station
    3: optional string arrival_station
}

// GetTrainDetailResp 获取车次详情响应（包含各席别余票与最低价）。
// 时间字段说明：departure_time_unix/arrival_time_unix 为 Unix 时间戳（秒）。
struct GetTrainDetailResp {
    1: common.BaseResp base_resp
    2: string train_id
    3: string train_code
    4: string train_type
    5: string departure_station
    6: string arrival_station
    7: i64 departure_time_unix
    8: i64 arrival_time_unix
    9: i32 runtime_minutes
    10: list<SeatTypeRemain> seat_types
}

// TicketService 票务域 RPC 服务。
service TicketService{
   // GetTrainDetail 获取车次详情（余票与最低价）。
   GetTrainDetailResp GetTrainDetail(1: GetTrainDetailReq req)
}
