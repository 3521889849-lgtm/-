package coupon

import "gorm.io/gorm"

// 优惠券核销/撤销记录表
type CouponVerifyRecord struct {
	gorm.Model
	UserCouponID    int     `gorm:"type:int unsigned;comment:关联user_coupon表的用户券ID"`
	OrderID         int     `gorm:"type:int;comment:关联订单ID"`
	PayAmount       float64 `gorm:"type:decimal(10,2);comment:订单原价金额"`
	DiscountAmount  float64 `gorm:"type:decimal(10,2);comment:优惠券抵扣金额"`
	ActualPayAmount float64 `gorm:"type:decimal(10,2);comment:订单实付金额"`
	VerifyStatus    int     `gorm:"type:int;comment:核销状态：1=已核销 2=已撤销"`
	OperateType     int     `gorm:"type:int;comment:操作类型：1=支付核销 2=退款撤销"`
	RefundOrderNo   string  `gorm:"type:varchar(50);comment:退款单号（仅撤销时有效）"`
}
