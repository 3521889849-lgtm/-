package model

import (
	"time"

	"gorm.io/gorm"
)

// PayRecord 支付记录表-支付回调、退款的核心流水，对账必备，匹配支付与退款流程
type PayRecord struct {
	ID               uint64         `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:支付记录主键ID" json:"id"`
	OrderID          uint64         `gorm:"column:order_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_order_id;comment:关联订单ID" json:"order_id"`
	OrderNo          string         `gorm:"column:order_no;type:VARCHAR(32);NOT NULL;comment:订单编号（冗余）" json:"order_no"`
	PayType          string         `gorm:"column:pay_type;type:VARCHAR(20);NOT NULL;comment:支付方式：WECHAT-微信，ALIPAY-支付宝" json:"pay_type"`
	PayAmount        float64        `gorm:"column:pay_amount;type:DECIMAL(10,2);NOT NULL;comment:支付金额" json:"pay_amount"`
	PayStatus        string         `gorm:"column:pay_status;type:VARCHAR(20);NOT NULL;index:idx_pay_status;comment:支付状态：SUCCESS-成功，FAIL-失败，REFUND-退款，REFUNDING-退款中" json:"pay_status"`
	PlatformTradeNo  *string        `gorm:"column:platform_trade_no;type:VARCHAR(64);index:idx_platform_trade_no;comment:支付平台流水号（微信/支付宝返回）" json:"platform_trade_no,omitempty"`
	PlatformRefundNo *string        `gorm:"column:platform_refund_no;type:VARCHAR(64);comment:支付平台退款单号" json:"platform_refund_no,omitempty"`
	NotifyTime       *time.Time     `gorm:"column:notify_time;type:DATETIME;comment:支付平台回调时间" json:"notify_time,omitempty"`
	ExtFields        *JSON          `gorm:"column:ext_fields;type:JSON;comment:扩展字段，如支付签名、回调参数等" json:"ext_fields,omitempty"`
	CreatedAt        time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	Order *OrderMain `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
}

func (PayRecord) TableName() string {
	return "pay_record"
}
