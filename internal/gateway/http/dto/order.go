package dto

// TrainDetailHTTPReq 车次详情 HTTP 请求参数（通过 URL Query 传参）。
//
// 对应接口：GET /api/v1/train/detail
type TrainDetailHTTPReq struct {
	TrainID          string `query:"train_id,required"` // 车次ID
	DepartureStation string `query:"departure_station"` // 出发站（可选；为空则由后端根据车次推导）
	ArrivalStation   string `query:"arrival_station"`   // 到达站（可选；为空则由后端根据车次推导）
}

// SeatTypeRemainItem 某席别的余票与最低价信息。
type SeatTypeRemainItem struct {
	SeatType  string  `json:"seat_type"` // 席别
	Remaining int32   `json:"remaining"` // 余票数量
	MinPrice  float64 `json:"min_price"` // 最低票价（单位：元）
}

// TrainDetailHTTPResp 车次详情 HTTP 响应。
//
// 时间字段说明：
// - DepartureTimeUnix/ArrivalTimeUnix：Unix 时间戳（秒）
type TrainDetailHTTPResp struct {
	Code              int32                `json:"code"`                // 业务状态码
	Msg               string               `json:"msg"`                 // 状态说明/错误信息
	TrainID           string               `json:"train_id"`            // 车次ID
	TrainCode         string               `json:"train_code"`          // 车次编码（展示用，如 G1234）
	TrainType         string               `json:"train_type"`          // 车次类型（如 G/D/K/T 等）
	DepartureStation  string               `json:"departure_station"`   // 出发站
	ArrivalStation    string               `json:"arrival_station"`     // 到达站
	DepartureTimeUnix int64                `json:"departure_time_unix"` // 出发时间（Unix 秒）
	ArrivalTimeUnix   int64                `json:"arrival_time_unix"`   // 到达时间（Unix 秒）
	RuntimeMinutes    int32                `json:"runtime_minutes"`     // 运行时长（分钟）
	SeatTypes         []SeatTypeRemainItem `json:"seat_types"`          // 各席别余票与最低价
}

// CreateOrderPassengerHTTP 创建订单时的乘客信息。
type CreateOrderPassengerHTTP struct {
	PassengerID uint64 `json:"passenger_id"`          // 常用乘车人ID（可选；传了则后端从库里补全实名信息）
	UseSelf     bool   `json:"use_self"`              // 使用当前登录账号的实名信息（可选；true 则忽略 passenger_id/real_name/id_card）
	RealName    string `json:"real_name"`             // 乘客姓名（未选乘车人时必填）
	IDCard      string `json:"id_card"`               // 乘客身份证号（未选乘车人时必填）
	SeatType    string `json:"seat_type,required"`    // 乘客席别偏好/购票席别
}

// CreateOrderHTTPReq 创建订单 HTTP 请求（JSON Body）。
//
// 对应接口：POST /api/v1/order/create
type CreateOrderHTTPReq struct {
	UserID           string                     `json:"user_id"`                    // 用户ID（由Token推导；兼容旧客户端可传）
	TrainID          string                     `json:"train_id,required"`          // 车次ID
	DepartureStation string                     `json:"departure_station,required"` // 出发站
	ArrivalStation   string                     `json:"arrival_station,required"`   // 到达站
	Passengers       []CreateOrderPassengerHTTP `json:"passengers,required"`        // 乘客列表（至少 1 人）
}

// OrderSeatInfoHTTP 订单席位信息（出票/锁座结果）。
type OrderSeatInfoHTTP struct {
	SeatID      string  `json:"seat_id"`      // 座位ID（后端唯一标识）
	SeatType    string  `json:"seat_type"`    // 席别
	CarriageNum string  `json:"carriage_num"` // 车厢号（展示用）
	SeatNum     string  `json:"seat_num"`     // 座位号（展示用）
	SeatPrice   float64 `json:"seat_price"`   // 价格（单位：元）
}

// CreateOrderHTTPResp 创建订单 HTTP 响应。
//
// 支付截止时间说明：
// - PayDeadlineUnix：Unix 时间戳（秒），超过该时间未支付会被后端取消/释放锁座
type CreateOrderHTTPResp struct {
	Code            int32               `json:"code"`              // 业务状态码
	Msg             string              `json:"msg"`               // 状态说明/错误信息
	OrderID         string              `json:"order_id"`          // 订单ID
	PayDeadlineUnix int64               `json:"pay_deadline_unix"` // 支付截止时间（Unix 秒）
	Seats           []OrderSeatInfoHTTP `json:"seats"`             // 锁座/出票的席位信息
}

// PayOrderHTTPReq 支付订单 HTTP 请求（JSON Body）。
//
// 对应接口：POST /api/v1/order/pay
type PayOrderHTTPReq struct {
	UserID     string `json:"user_id"`              // 用户ID（由Token推导；兼容旧客户端可传）
	OrderID    string `json:"order_id,required"`    // 订单ID
	PayChannel string `json:"pay_channel,required"` // 支付渠道（如 ALIPAY/WECHAT/CARD；由后端定义）
	PayNo      string `json:"pay_no"`               // 支付流水号（可选；不传由后端生成/模拟）
}

// PayOrderHTTPResp 支付订单 HTTP 响应。
type PayOrderHTTPResp struct {
	Code        int32  `json:"code"`         // 业务状态码
	Msg         string `json:"msg"`          // 状态说明/错误信息
	OrderStatus string `json:"order_status"` // 订单状态（由后端定义）
	PayStatus   string `json:"pay_status"`   // 支付状态（由后端定义）
	PayNo       string `json:"pay_no"`       // 支付流水号
	PayURL      string `json:"pay_url,omitempty"`
}

// ConfirmPayHTTPReq 支付回调确认请求（JSON Body）。
//
// 对应接口：
// - POST /api/v1/pay/callback（模拟第三方异步回调）
// - POST /api/v1/pay/mock_notify（本地触发一条模拟回调）
type ConfirmPayHTTPReq struct {
	UserID           string `json:"user_id"`            // 用户ID（由Token推导；兼容旧客户端可传）
	OrderID          string `json:"order_id,required"`  // 订单ID
	PayNo            string `json:"pay_no,required"`    // 支付流水号
	ThirdPartyStatus string `json:"third_party_status"` // 第三方支付状态（如 SUCCESS/FAIL；由后端定义）
}

// ConfirmPayHTTPResp 支付确认 HTTP 响应。
type ConfirmPayHTTPResp struct {
	Code        int32  `json:"code"`         // 业务状态码
	Msg         string `json:"msg"`          // 状态说明/错误信息
	OrderStatus string `json:"order_status"` // 最新订单状态
}

// ChangeOrderHTTPReq 改签 HTTP 请求（JSON Body）。
//
// 对应接口：POST /api/v1/order/change
type ChangeOrderHTTPReq struct {
	UserID              string `json:"user_id"`                        // 用户ID（由Token推导；兼容旧客户端可传）
	OrderID             string `json:"order_id,required"`              // 原订单ID
	NewTrainID          string `json:"new_train_id,required"`          // 新车次ID
	NewDepartureStation string `json:"new_departure_station,required"` // 新出发站
	NewArrivalStation   string `json:"new_arrival_station,required"`   // 新到达站
}

// ChangeOrderHTTPResp 改签 HTTP 响应。
type ChangeOrderHTTPResp struct {
	Code             int32               `json:"code"`               // 业务状态码
	Msg              string              `json:"msg"`                // 状态说明/错误信息
	OldOrderStatus   string              `json:"old_order_status"`   // 原订单状态
	NewOrderID       string              `json:"new_order_id"`       // 改签后生成的新订单ID
	NewTotalAmount   float64             `json:"new_total_amount"`   // 新订单总金额（单位：元）
	RefundDiffAmount float64             `json:"refund_diff_amount"` // 差额退补金额（>0 退，<0 补）
	Seats            []OrderSeatInfoHTTP `json:"seats"`              // 新订单席位信息
}

// RefundOrderHTTPReq 退票 HTTP 请求（JSON Body）。
//
// 对应接口：POST /api/v1/order/refund
type RefundOrderHTTPReq struct {
	UserID  string `json:"user_id"`           // 用户ID（由Token推导；兼容旧客户端可传）
	OrderID string `json:"order_id,required"` // 订单ID
	Reason  string `json:"reason"`            // 退票原因（可选）
}

// RefundOrderHTTPResp 退票 HTTP 响应。
type RefundOrderHTTPResp struct {
	Code         int32   `json:"code"`          // 业务状态码
	Msg          string  `json:"msg"`           // 状态说明/错误信息
	RefundAmount float64 `json:"refund_amount"` // 退款金额（单位：元）
	RefundStatus string  `json:"refund_status"` // 退款状态（由后端定义）
	OrderStatus  string  `json:"order_status"`  // 最新订单状态
}

// CancelOrderHTTPReq 取消订单 HTTP 请求（JSON Body）。
//
// 对应接口：POST /api/v1/order/cancel
type CancelOrderHTTPReq struct {
	UserID  string `json:"user_id"`           // 用户ID（由Token推导；兼容旧客户端可传）
	OrderID string `json:"order_id,required"` // 订单ID
}

// GetOrderHTTPReq 查询订单详情 HTTP 请求参数（通过 URL Query 传参）。
//
// 对应接口：GET /api/v1/order/info
type GetOrderHTTPReq struct {
	UserID  string `query:"user_id"`            // 用户ID（由Token推导；兼容旧客户端可传）
	OrderID string `query:"order_id,required"`  // 订单ID
}

// ListOrdersHTTPReq 查询订单列表 HTTP 请求参数（通过 URL Query 传参）。
//
// 对应接口：GET /api/v1/order/list
type ListOrdersHTTPReq struct {
	UserID string `query:"user_id"`          // 用户ID（由Token推导；兼容旧客户端可传）
	Status string `query:"status"`           // 状态过滤（为空表示不过滤）
	Cursor string `query:"cursor"`           // 游标分页游标
	Limit  int    `query:"limit"`            // 每页条数（建议 1~50；不传由后端设默认）
}

// ListOrdersHTTPResp 查询订单列表 HTTP 响应。
type ListOrdersHTTPResp struct {
	Code       int32           `json:"code"`        // 业务状态码
	Msg        string          `json:"msg"`         // 状态说明/错误信息
	Orders     []OrderInfoHTTP `json:"orders"`      // 订单列表
	NextCursor *string         `json:"next_cursor"` // 下一页游标（没有下一页时可能为 nil）
}

// OrderInfoHTTP 订单基本信息（列表/详情通用）。
type OrderInfoHTTP struct {
	OrderID          string  `json:"order_id"`          // 订单ID
	UserID           string  `json:"user_id"`           // 用户ID
	TrainID          string  `json:"train_id"`          // 车次ID
	DepartureStation string  `json:"departure_station"` // 出发站
	ArrivalStation   string  `json:"arrival_station"`   // 到达站
	FromSeq          int32   `json:"from_seq"`          // 区间起始站序（用于占座/计价）
	ToSeq            int32   `json:"to_seq"`            // 区间终止站序（用于占座/计价）
	TotalAmount      float64 `json:"total_amount"`      // 订单总金额（单位：元）
	OrderStatus      string  `json:"order_status"`      // 订单状态
	PayDeadlineUnix  int64   `json:"pay_deadline_unix"` // 支付截止时间（Unix 秒）
	PayTimeUnix      *int64  `json:"pay_time_unix"`     // 支付时间（Unix 秒；未支付为 nil）
	PayChannel       *string `json:"pay_channel"`       // 支付渠道（未支付为 nil）
	PayNo            *string `json:"pay_no"`            // 支付流水号（未支付为 nil）
	CreatedAtUnix    int64   `json:"created_at_unix"`   // 创建时间（Unix 秒）
}

// GetOrderHTTPResp 查询订单详情 HTTP 响应。
type GetOrderHTTPResp struct {
	Code  int32               `json:"code"`  // 业务状态码
	Msg   string              `json:"msg"`   // 状态说明/错误信息
	Order *OrderInfoHTTP      `json:"order"` // 订单信息（可能为 nil）
	Seats []OrderSeatInfoHTTP `json:"seats"` // 席位信息（详情接口返回）
}
