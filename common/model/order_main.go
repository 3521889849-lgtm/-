package model

import (
	"time"

	"gorm.io/gorm"
)

// OrderMain 主订单表-核心业务表，订单状态流转全量记录，匹配工单所有状态
type OrderMain struct {
	ID           uint64         `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:订单主键ID" json:"id"`
	OrderNo      string         `gorm:"column:order_no;type:VARCHAR(32);NOT NULL;uniqueIndex:uk_order_no;comment:订单编号，唯一，生成规则：时间戳+用户ID+随机数" json:"order_no"`
	UserID       uint64         `gorm:"column:user_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_user_id;comment:下单用户ID" json:"user_id"`
	MerchantID   uint64         `gorm:"column:merchant_id;type:BIGINT UNSIGNED;NOT NULL;comment:所属商家ID" json:"merchant_id"`
	SpotID       uint64         `gorm:"column:spot_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_spot_id;comment:所属景点ID" json:"spot_id"`
	TotalAmount  float64        `gorm:"column:total_amount;type:DECIMAL(10,2);NOT NULL;comment:订单总金额" json:"total_amount"`
	PayAmount    float64        `gorm:"column:pay_amount;type:DECIMAL(10,2);NOT NULL;comment:实际支付金额（含优惠券抵扣）" json:"pay_amount"`
	CouponID     uint64         `gorm:"column:coupon_id;type:BIGINT UNSIGNED;default:0;comment:使用的优惠券ID，0=未使用" json:"coupon_id"`
	OrderStatus  string         `gorm:"column:order_status;type:VARCHAR(30);NOT NULL;default:'DRAFT';index:idx_order_status;comment:订单状态" json:"order_status"`
	PayType      *string        `gorm:"column:pay_type;type:VARCHAR(20);comment:支付方式：WECHAT-微信，ALIPAY-支付宝" json:"pay_type,omitempty"`
	PayTime      *time.Time     `gorm:"column:pay_time;type:DATETIME;comment:支付时间" json:"pay_time,omitempty"`
	VerifyCode   *string        `gorm:"column:verify_code;type:VARCHAR(64);uniqueIndex:uk_verify_code;comment:门票核销码，唯一，入园使用" json:"verify_code,omitempty"`
	VerifyTime   *time.Time     `gorm:"column:verify_time;type:DATETIME;comment:核销使用时间" json:"verify_time,omitempty"`
	CancelTime   *time.Time     `gorm:"column:cancel_time;type:DATETIME;comment:订单取消时间" json:"cancel_time,omitempty"`
	RefundAmount float64        `gorm:"column:refund_amount;type:DECIMAL(10,2);default:0.00;comment:退款金额" json:"refund_amount"`
	RefundTime   *time.Time     `gorm:"column:refund_time;type:DATETIME;comment:退款完成时间" json:"refund_time,omitempty"`
	ExtFields    *JSON          `gorm:"column:ext_fields;type:JSON;comment:扩展字段，如支付流水号、退款单号等" json:"ext_fields,omitempty"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;index:idx_create_time;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	User       *SysUser     `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Merchant   *SysMerchant `gorm:"foreignKey:MerchantID;references:ID" json:"merchant,omitempty"`
	Spot       *SpotInfo    `gorm:"foreignKey:SpotID;references:ID" json:"spot,omitempty"`
	OrderItems []OrderItem  `gorm:"foreignKey:OrderID;references:ID" json:"order_items,omitempty"`
	PayRecords []PayRecord  `gorm:"foreignKey:OrderID;references:ID" json:"pay_records,omitempty"`
}

func (OrderMain) TableName() string {
	return "order_main"
}
