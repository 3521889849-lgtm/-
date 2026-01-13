package coupon

import "gorm.io/gorm"

// CouponApplyRange 优惠券适用范围关联表
type CouponApplyRange struct {
	gorm.Model
	CouponID   int    `gorm:"type:bigint unsigned;comment:关联coupon_base表的优惠券ID"`
	ApplyType  int    `gorm:"type:tinyint;comment:适用对象类型：1=票务类型 2=商家 3=线路 4=房型"`
	ApplyValue string `gorm:"type:varchar(100);comment:适用对象值（如“景点门票”“商家A”“线路B”）"`
}
