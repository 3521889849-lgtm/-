// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了关联房间绑定的数据模型
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// RelatedRoomBinding 关联房绑定表
//
// 业务用途：
//   - 满足多房间绑定需求（如：套房由多个房间组成）
//   - 主房间预订时，自动预订关联房间
//   - 灵活配置房间组合关系
//   - 示例：豪华套房 = 主卧（101） + 次卧（102）
type RelatedRoomBinding struct {
	ID            uint64  `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:绑定ID" json:"id"`
	MainRoomID    uint64  `gorm:"column:main_room_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_main_room_id;comment:主房源ID（外键，关联房源表）" json:"main_room_id"`
	RelatedRoomID uint64  `gorm:"column:related_room_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_related_room_id;comment:关联房源ID（外键，关联房源表）" json:"related_room_id"`
	BindingDesc   *string `gorm:"column:binding_desc;type:VARCHAR(255);comment:绑定描述" json:"binding_desc,omitempty"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	MainRoom    *RoomInfo `gorm:"foreignKey:MainRoomID;references:ID" json:"main_room,omitempty"`    // 主房间
	RelatedRoom *RoomInfo `gorm:"foreignKey:RelatedRoomID;references:ID" json:"related_room,omitempty"` // 关联房间
}

// TableName 指定数据库表名
func (RelatedRoomBinding) TableName() string {
	return "related_room_binding"
}
