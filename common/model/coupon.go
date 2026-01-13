package model

import (
	"time"

	"gorm.io/gorm"
)

// Coupon 优惠券配置表-管理端配置优惠券/限时活动核心表
type Coupon struct {
	ID             uint64         `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:优惠券主键ID" json:"id"`
	CouponName     string         `gorm:"column:coupon_name;type:VARCHAR(100);NOT NULL;comment:优惠券名称" json:"coupon_name"`
	CouponType     string         `gorm:"column:coupon_type;type:VARCHAR(20);NOT NULL;comment:优惠券类型：FIXED-满减券，DISCOUNT-折扣券" json:"coupon_type"`
	Denomination   float64        `gorm:"column:denomination;type:DECIMAL(10,2);NOT NULL;comment:优惠券面额/折扣率" json:"denomination"`
	MinUseAmount   float64        `gorm:"column:min_use_amount;type:DECIMAL(10,2);NOT NULL;comment:最低使用金额" json:"min_use_amount"`
	ValidStartTime time.Time      `gorm:"column:valid_start_time;type:DATETIME;NOT NULL;index:idx_valid_time;comment:有效期开始时间" json:"valid_start_time"`
	ValidEndTime   time.Time      `gorm:"column:valid_end_time;type:DATETIME;NOT NULL;index:idx_valid_time;comment:有效期结束时间" json:"valid_end_time"`
	Stock          uint32         `gorm:"column:stock;type:INT UNSIGNED;NOT NULL;default:0;comment:优惠券库存" json:"stock"`
	ApplySpotIDs   *string        `gorm:"column:apply_spot_ids;type:VARCHAR(512);comment:适用景点ID集合，逗号分隔，空=全景点通用" json:"apply_spot_ids,omitempty"`
	CouponStatus   string         `gorm:"column:coupon_status;type:VARCHAR(20);NOT NULL;default:'VALID';index:idx_coupon_status;comment:状态：VALID-有效，INVALID-失效" json:"coupon_status"`
	ExtFields      *JSON          `gorm:"column:ext_fields;type:JSON;comment:扩展字段，如使用规则、限制条件" json:"ext_fields,omitempty"`
	CreatedAt      time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	UserCoupons []UserCoupon `gorm:"foreignKey:CouponID;references:ID" json:"user_coupons,omitempty"`
}

func (Coupon) TableName() string {
	return "coupon"
}
