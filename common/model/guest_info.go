package model

import (
	"time"

	"gorm.io/gorm"
)

// GuestInfo 客人信息表（实名制备案）
type GuestInfo struct {
	ID         int            `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	Name       string         `gorm:"type:varchar(50);not null;comment:客人姓名" json:"name"`
	IDCardType string         `gorm:"type:varchar(20);not null;comment:证件类型（如：身份证、护照）" json:"id_card_type"`
	IDCardNo   string         `gorm:"type:varchar(30);not null;comment:证件号（加密存储）" json:"id_card_no"`
	Nation     string         `gorm:"type:varchar(20);comment:民族" json:"nation"`
	Province   string         `gorm:"type:varchar(50);index:idx_province;comment:省份" json:"province"`
	Address    string         `gorm:"type:varchar(255);comment:详细地址" json:"address"`
	Phone      string         `gorm:"type:varchar(20);not null;index:idx_phone;comment:手机号" json:"phone"`
	CreatedAt  time.Time      `gorm:"not null;comment:创建时间" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"not null;comment:更新时间" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index:idx_deleted_at;comment:软删除时间" json:"deleted_at"`

	// 关联关系
	Orders       []OrderInfo   `gorm:"foreignKey:GuestID" json:"orders,omitempty"`
	FinanceFlows []FinanceFlow `gorm:"foreignKey:GuestID" json:"finance_flows,omitempty"`
	MemberInfo   *MemberInfo   `gorm:"foreignKey:GuestID" json:"member_info,omitempty"`
	Blacklist    *Blacklist    `gorm:"foreignKey:GuestID" json:"blacklist,omitempty"`
}

// TableName 指定表名
func (GuestInfo) TableName() string {
	return "guest_info"
}
