package model

import (
	"time"

	"gorm.io/gorm"
)

// UserCoupon 用户优惠券表-用户领取后存储，预订门票时自动抵扣，匹配营销活动流程
type UserCoupon struct {
	ID             uint64         `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:用户优惠券主键ID" json:"id"`
	UserID         uint64         `gorm:"column:user_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_user_id;comment:用户ID" json:"user_id"`
	CouponID       uint64         `gorm:"column:coupon_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_coupon_id;comment:优惠券ID" json:"coupon_id"`
	CouponName     string         `gorm:"column:coupon_name;type:VARCHAR(100);NOT NULL;comment:优惠券名称（冗余）" json:"coupon_name"`
	Denomination   float64        `gorm:"column:denomination;type:DECIMAL(10,2);NOT NULL;comment:优惠券面额" json:"denomination"`
	MinUseAmount   float64        `gorm:"column:min_use_amount;type:DECIMAL(10,2);NOT NULL;comment:最低使用金额" json:"min_use_amount"`
	ValidStartTime time.Time      `gorm:"column:valid_start_time;type:DATETIME;NOT NULL;comment:有效期开始时间" json:"valid_start_time"`
	ValidEndTime   time.Time      `gorm:"column:valid_end_time;type:DATETIME;NOT NULL;comment:有效期结束时间" json:"valid_end_time"`
	UseStatus      string         `gorm:"column:use_status;type:VARCHAR(20);NOT NULL;default:'UNUSED';index:idx_use_status;comment:使用状态：UNUSED-未使用，USED-已使用，EXPIRED-已过期" json:"use_status"`
	OrderID        uint64         `gorm:"column:order_id;type:BIGINT UNSIGNED;default:0;comment:使用的订单ID，0=未使用" json:"order_id"`
	UseTime        *time.Time     `gorm:"column:use_time;type:DATETIME;comment:使用时间" json:"use_time,omitempty"`
	CreatedAt      time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	User   *SysUser `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Coupon *Coupon  `gorm:"foreignKey:CouponID;references:ID" json:"coupon,omitempty"`
}

func (UserCoupon) TableName() string {
	return "user_coupon"
}
