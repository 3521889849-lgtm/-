// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了订单扩展信息的数据模型
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// OrderExtension 订单明细扩展表
//
// 业务用途：
//   - 存储订单的补充信息（联系人、特殊需求等）
//   - 支持多渠道数据同步功能
//   - 记录订单与第三方平台的同步状态
//   - 一对一关联订单主表
type OrderExtension struct {
	ID             uint64  `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:扩展ID" json:"id"`
	OrderID        uint64  `gorm:"column:order_id;type:BIGINT UNSIGNED;NOT NULL;uniqueIndex:uk_order_id;comment:订单ID（外键，关联订单主表）" json:"order_id"`
	Contact        string  `gorm:"column:contact;type:VARCHAR(50);NOT NULL;comment:联系人" json:"contact"`
	ContactPhone   string  `gorm:"column:contact_phone;type:VARCHAR(20);NOT NULL;comment:联系电话" json:"contact_phone"`
	SpecialRequest *string `gorm:"column:special_request;type:VARCHAR(500);comment:特殊需求" json:"special_request,omitempty"`
	GuestCount     uint8   `gorm:"column:guest_count;type:TINYINT UNSIGNED;NOT NULL;default:1;comment:入住人数" json:"guest_count"`
	RoomCount      uint8   `gorm:"column:room_count;type:TINYINT UNSIGNED;NOT NULL;default:1;comment:房间数量" json:"room_count"`
	
	// 同步状态（用于第三方平台同步）
	SyncStatus string     `gorm:"column:sync_status;type:VARCHAR(20);NOT NULL;default:'PENDING';index:idx_sync_status;comment:数据同步状态（是否同步至第三方渠道）" json:"sync_status"`
	SyncTime   *time.Time `gorm:"column:sync_time;type:DATETIME;comment:同步时间" json:"sync_time,omitempty"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	Order *OrderMain `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"` // 订单主表
}

// TableName 指定数据库表名
func (OrderExtension) TableName() string {
	return "order_extension"
}
