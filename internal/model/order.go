package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// OrderInfo 车票订单表
type OrderInfo struct {
	ID            string         `gorm:"column:order_id;type:VARCHAR(64);primaryKey;comment:订单唯一标识（UUID生成）" json:"order_id"`
	UserID        string         `gorm:"column:user_id;type:VARCHAR(64);NOT NULL;comment:下单用户ID，关联user_info.user_id" json:"user_id"`
	TrainID       string         `gorm:"column:train_id;type:VARCHAR(32);NOT NULL;comment:车次编号，关联train_info.train_id" json:"train_id"`
	DepartureStation string      `gorm:"column:departure_station;type:VARCHAR(64);NOT NULL;comment:上车站（区间购票起点）" json:"departure_station"`
	ArrivalStation   string      `gorm:"column:arrival_station;type:VARCHAR(64);NOT NULL;comment:下车站（区间购票终点）" json:"arrival_station"`
	FromSeq          uint32      `gorm:"column:from_seq;type:INT;NOT NULL;comment:上车站序号（对应train_station_pass.sequence）" json:"from_seq"`
	ToSeq            uint32      `gorm:"column:to_seq;type:INT;NOT NULL;comment:下车站序号（对应train_station_pass.sequence）" json:"to_seq"`
	TotalAmount   float64        `gorm:"column:total_amount;type:DECIMAL(10,2);NOT NULL;comment:订单总金额" json:"total_amount"`
	OrderStatus   string         `gorm:"column:order_status;type:VARCHAR(20);NOT NULL;index:idx_order_status;comment:订单状态：PENDING_PAY-待支付，PAYING-支付中，ISSUED-已出票，CANCELLED-已取消，REFUNDED-已退票，CHANGED-已改签" json:"order_status"`
	PayDeadline   sql.NullTime   `gorm:"column:pay_deadline;type:DATETIME;index:idx_pay_deadline;comment:支付截止时间（待支付状态有效）" json:"pay_deadline,omitempty"`
	PayTime       sql.NullTime   `gorm:"column:pay_time;type:DATETIME;comment:支付时间（已支付状态有效）" json:"pay_time,omitempty"`
	PayChannel    sql.NullString `gorm:"column:pay_channel;type:VARCHAR(16);comment:支付渠道（微信、支付宝、银联）" json:"pay_channel,omitempty"`
	PayNo         sql.NullString `gorm:"column:pay_no;type:VARCHAR(64);comment:第三方支付流水号" json:"pay_no,omitempty"`
	RefundAmount  float64        `gorm:"column:refund_amount;type:DECIMAL(10,2);default:0.00;comment:退款金额（退票/改签时有效）" json:"refund_amount"`
	RefundTime    sql.NullTime   `gorm:"column:refund_time;type:DATETIME;comment:退款时间（退票/改签时有效）" json:"refund_time,omitempty"`
	RefundStatus  string         `gorm:"column:refund_status;type:VARCHAR(20);default:'NO_REFUND';comment:退款状态：NO_REFUND-未退款，REFUNDING-退款中，REFUNDED-已退款" json:"refund_status"`
	IdempotentKey string         `gorm:"column:idempotent_key;type:VARCHAR(128);NOT NULL;uniqueIndex;comment:幂等键（user_id+train_id+座位组合MD5）" json:"idempotent_key,omitempty"`
	CreatedAt     time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`
}

// OrderSeatRelation 订单座位关联表
type OrderSeatRelation struct {
	ID         uint64         `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:自增主键" json:"id"`
	OrderID    string         `gorm:"column:order_id;type:VARCHAR(64);NOT NULL;comment:订单号，关联order_info.order_id" json:"order_id"`
	SeatID     string         `gorm:"column:seat_id;type:VARCHAR(64);NOT NULL;comment:座位ID，关联seat_info.seat_id" json:"seat_id"`
	SeatType   string         `gorm:"column:seat_type;type:VARCHAR(16);NOT NULL;comment:座位类型（冗余字段，优化查询）" json:"seat_type"`
	SeatPrice  float64        `gorm:"column:seat_price;type:DECIMAL(10,2);NOT NULL;comment:座位单价（冗余字段，优化查询）" json:"seat_price"`
	IsRefunded string         `gorm:"column:is_refunded;type:VARCHAR(20);NOT NULL;default:'NO';comment: 是否已退款：NO-未退款，YES-已退款" json:"is_refunded"`
	CreatedAt  time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`
}

// OrderAuditLog 订单操作审计表
type OrderAuditLog struct {
	ID            uint64         `gorm:"column:log_id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:日志ID" json:"log_id"`
	OrderID       string         `gorm:"column:order_id;type:VARCHAR(64);NOT NULL;comment:订单号，关联order_info.order_id" json:"order_id"`
	OperateType   string         `gorm:"column:operate_type;type:VARCHAR(32);NOT NULL;comment:操作类型：CREATE_ORDER-创建订单，PAY_ORDER-支付订单，REFUND_ORDER-退票，CHANGE_ORDER-改签，CANCEL_ORDER-取消订单" json:"operate_type"`
	OperateUser   string         `gorm:"column:operate_user;type:VARCHAR(64);NOT NULL;comment:操作人（用户ID或管理员ID）" json:"operate_user"`
	BeforeStatus  sql.NullString `gorm:"column:before_status;type:VARCHAR(20);comment:操作前状态（订单状态/座位状态）" json:"before_status,omitempty"`
	AfterStatus   sql.NullString `gorm:"column:after_status;type:VARCHAR(20);comment:操作后状态（订单状态/座位状态）" json:"after_status,omitempty"`
	OperateDetail *JSON          `gorm:"column:operate_detail;type:JSON;comment:操作详情（如退票手续费、改签前后车次信息）" json:"operate_detail,omitempty"`
	OperateIP     sql.NullString `gorm:"column:operate_ip;type:VARCHAR(64);comment:操作IP地址" json:"operate_ip,omitempty"`
	TraceID       string         `gorm:"column:trace_id;type:VARCHAR(64);NOT NULL;index:idx_trace_id;comment:链路ID（全链路追踪）" json:"trace_id"`
	CreatedAt     time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:操作时间" json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`
}
