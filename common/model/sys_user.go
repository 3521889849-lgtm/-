package model

import (
	"time"

	"gorm.io/gorm"
)

// SysUser 用户信息表-实名制信息存储，匹配安全合规要求
type SysUser struct {
	ID        uint64         `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:用户主键ID" json:"id"`
	UserName  string         `gorm:"column:user_name;type:VARCHAR(50);NOT NULL;comment:用户名" json:"user_name"`
	Phone     string         `gorm:"column:phone;type:VARCHAR(20);NOT NULL;uniqueIndex:uk_phone;comment:手机号（登录账号）" json:"phone"`
	IDCard    *string        `gorm:"column:id_card;type:VARCHAR(20);index:idx_id_card;comment:身份证号，实名制必填" json:"id_card,omitempty"`
	RealName  *string        `gorm:"column:real_name;type:VARCHAR(30);comment:真实姓名，实名制必填" json:"real_name,omitempty"`
	Avatar    *string        `gorm:"column:avatar;type:VARCHAR(512);comment:用户头像" json:"avatar,omitempty"`
	ExtFields *JSON          `gorm:"column:ext_fields;type:JSON;comment:扩展字段，如会员等级、积分等" json:"ext_fields,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	Travelers   []Traveler   `gorm:"foreignKey:UserID;references:ID" json:"travelers,omitempty"`
	Orders      []OrderMain  `gorm:"foreignKey:UserID;references:ID" json:"orders,omitempty"`
	UserCoupons []UserCoupon `gorm:"foreignKey:UserID;references:ID" json:"user_coupons,omitempty"`
}

func (SysUser) TableName() string {
	return "sys_user"
}
