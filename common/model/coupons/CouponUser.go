package coupon

import (
	"time"

	"gorm.io/gorm"
)

// 用户已领取的优惠券表
type UserCoupon struct {
	gorm.Model
	UserID        int       `gorm:"type:int;comment:用户ID"`
	CouponID      int       `gorm:"type:int unsigned;comment:关联coupon_base表的优惠券ID"`
	GrantSource   int       `gorm:"type:int;comment:发放来源：1=系统全员发放 2=定向发放 3=新人注册自动发放 4=活动手动领取 5=邀请有礼发放"`
	Status        int       `gorm:"type:int;comment:用户券状态：0=未使用 1=已使用 2=已过期"`
	ReceiveTime   time.Time `gorm:"type:datetime;comment:领取时间"`
	EffectiveTime time.Time `gorm:"type:datetime;comment:生效时间（领取后N天计算得出）"`
	ExpireTime    time.Time `gorm:"type:datetime;comment:失效时间（领取后N天计算得出）"`
	IsReturned    int       `gorm:"type:int;comment:是否退款退回：0=未退回 1=已退回"`
}
