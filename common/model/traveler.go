package model

import (
	"time"

	"gorm.io/gorm"
)

// Traveler 出行人信息表-用户下单时必填，匹配门票实名制要求
type Traveler struct {
	ID        uint64         `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:出行人主键ID" json:"id"`
	UserID    uint64         `gorm:"column:user_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_user_id;comment:所属用户ID" json:"user_id"`
	RealName  string         `gorm:"column:real_name;type:VARCHAR(30);NOT NULL;comment:出行人真实姓名" json:"real_name"`
	IDCard    string         `gorm:"column:id_card;type:VARCHAR(20);NOT NULL;index:idx_id_card;comment:出行人身份证号" json:"id_card"`
	Phone     string         `gorm:"column:phone;type:VARCHAR(20);NOT NULL;comment:出行人手机号" json:"phone"`
	IsDefault uint8          `gorm:"column:is_default;type:TINYINT UNSIGNED;NOT NULL;default:0;comment:是否默认出行人：0-否，1-是" json:"is_default"`
	CreatedAt time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	User       *SysUser    `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:TravelerID;references:ID" json:"order_items,omitempty"`
}

func (Traveler) TableName() string {
	return "traveler"
}
