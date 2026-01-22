// order_api.thrift
//
// 说明：
// - 订单域 RPC 接口定义（下单/支付/改签/退票/查询）。
// - 该文件只放“订单域”的请求/响应与服务，不与用户/票务混写。
namespace go orderapi

include "common.thrift"

// OrderSeatInfo 订单中的席位信息（锁座/出票结果）。
struct OrderSeatInfo {
    1: string seat_id
    2: string seat_type
    3: string carriage_num
    4: string seat_num
    5: double seat_price
}

// CreateOrderPassenger 创建订单时的乘客信息（实名制）。
struct CreateOrderPassenger {
    1: string real_name
    2: string id_card
    3: string seat_type
}

// CreateOrderReq 下单请求：锁座 + 创建订单（待支付）。
struct CreateOrderReq {
    1: string user_id
    2: string train_id
    3: string departure_station
    4: string arrival_station
    5: list<CreateOrderPassenger> passengers
}

// CreateOrderResp 下单响应。
// pay_deadline_unix：支付截止时间（Unix 秒），超时未付会被取消/释放锁座。
struct CreateOrderResp {
    1: common.BaseResp base_resp
    2: string order_id
    3: i64 pay_deadline_unix
    4: list<OrderSeatInfo> seats
}

// PayOrderReq 支付请求。
struct PayOrderReq {
    1: string user_id
    2: string order_id
    3: string pay_channel
    4: optional string pay_no
}

// PayOrderResp 支付响应。
struct PayOrderResp {
    1: common.BaseResp base_resp
    2: string order_status
    3: string pay_status
    4: optional string pay_no
}

// ConfirmPayReq 支付回调确认请求（模拟第三方异步通知）。
struct ConfirmPayReq {
    1: string user_id
    2: string order_id
    3: string pay_no
    4: optional string third_party_status
}

// ConfirmPayResp 支付回调确认响应。
struct ConfirmPayResp {
    1: common.BaseResp base_resp
    2: string order_status
}

// ChangeOrderReq 改签请求：更换车次/区间（可能生成新订单）。
struct ChangeOrderReq {
    1: string user_id
    2: string order_id
    3: string new_train_id
    4: string new_departure_station
    5: string new_arrival_station
}

// ChangeOrderResp 改签响应。
struct ChangeOrderResp {
    1: common.BaseResp base_resp
    2: string old_order_status
    3: string new_order_id
    4: double new_total_amount
    5: double refund_diff_amount
    6: list<OrderSeatInfo> seats
}

// CancelOrderReq 取消订单请求（释放锁座/占用）。
struct CancelOrderReq {
    1: string user_id
    2: string order_id
}

// RefundOrderReq 退票请求（触发退票与退款流程）。
struct RefundOrderReq {
    1: string user_id
    2: string order_id
    3: optional string reason
}

// RefundOrderResp 退票响应。
struct RefundOrderResp {
    1: common.BaseResp base_resp
    2: double refund_amount
    3: string refund_status
    4: string order_status
}

// ListOrdersReq 查询订单列表请求（游标分页）。
struct ListOrdersReq {
    1: string user_id
    2: optional string status
    3: optional string cursor
    4: i32 limit
}

// OrderInfo 订单基本信息（列表/详情通用）。
struct OrderInfo {
    1: string order_id
    2: string user_id
    3: string train_id
    4: string departure_station
    5: string arrival_station
    6: i32 from_seq
    7: i32 to_seq
    8: double total_amount
    9: string order_status
    10: i64 pay_deadline_unix
    11: optional i64 pay_time_unix
    12: optional string pay_channel
    13: optional string pay_no
    14: i64 created_at_unix
}

// ListOrdersResp 查询订单列表响应。
struct ListOrdersResp {
    1: common.BaseResp base_resp
    2: list<OrderInfo> orders
    3: optional string next_cursor
}

// GetOrderReq 查询订单详情请求。
struct GetOrderReq {
    1: string user_id
    2: string order_id
}

// GetOrderResp 查询订单详情响应。
struct GetOrderResp {
    1: common.BaseResp base_resp
    2: optional OrderInfo order
    3: list<OrderSeatInfo> seats
}

// OrderService 订单域 RPC 服务。
service OrderService{
   // CreateOrder 创建订单并锁座。
   CreateOrderResp CreateOrder(1: CreateOrderReq req)
   // PayOrder 支付订单（简化版：可能直接出票）。
   PayOrderResp PayOrder(1: PayOrderReq req)
   // ConfirmPay 支付回调确认（模拟第三方异步通知）。
   ConfirmPayResp ConfirmPay(1: ConfirmPayReq req)
   // ChangeOrder 改签（可能生成新订单）。
   ChangeOrderResp ChangeOrder(1: ChangeOrderReq req)
   // CancelOrder 取消订单（释放锁座/占用）。
   common.BaseResp CancelOrder(1: CancelOrderReq req)
   // RefundOrder 退票（触发退票与退款流程）。
   RefundOrderResp RefundOrder(1: RefundOrderReq req)
   // GetOrder 查询订单详情。
   GetOrderResp GetOrder(1: GetOrderReq req)
   // ListOrders 查询订单列表（游标分页）。
   ListOrdersResp ListOrders(1: ListOrdersReq req)
}
