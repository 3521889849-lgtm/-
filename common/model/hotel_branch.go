package model

import (
	"time"

	"gorm.io/gorm"
)

// HotelBranch 酒店/分店信息表
type HotelBranch struct {
	ID             int            `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	BranchName     string         `gorm:"type:varchar(100);not null;comment:分店名称（如：馨香之约客栈-洱海店）" json:"branch_name"`
	ContactPerson  string         `gorm:"type:varchar(50);comment:联系人" json:"contact_person"`
	ContactPhone   string         `gorm:"type:varchar(20);comment:联系电话" json:"contact_phone"`
	Address        string         `gorm:"type:varchar(255);comment:分店地址" json:"address"`
	Status         int8           `gorm:"type:tinyint(1);not null;default:1;comment:状态（0-停用，1-启用）" json:"status"`
	CreatedAt      time.Time      `gorm:"not null;comment:创建时间" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"not null;comment:更新时间" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index:idx_deleted_at;comment:软删除时间" json:"deleted_at"`

	// 关联关系
	RoomTypes          []RoomType          `gorm:"foreignKey:BranchID" json:"room_types,omitempty"`
	RoomInfos          []RoomInfo          `gorm:"foreignKey:BranchID" json:"room_infos,omitempty"`
	Orders             []OrderInfo         `gorm:"foreignKey:BranchID" json:"orders,omitempty"`
	FinanceFlows       []FinanceFlow       `gorm:"foreignKey:BranchID" json:"finance_flows,omitempty"`
	ShiftChangeRecords []ShiftChangeRecord `gorm:"foreignKey:BranchID" json:"shift_change_records,omitempty"`
	Users              []SysUser           `gorm:"foreignKey:BranchID" json:"users,omitempty"`
}

// TableName 指定表名
func (HotelBranch) TableName() string {
	return "hotel_branch"
}
