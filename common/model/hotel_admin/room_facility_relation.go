// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了房间-设施关联的数据模型，实现多对多关系
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// RoomFacilityRelation 房源-设施关联表
//
// 业务用途：
//   - 实现房间和设施的多对多关联
//   - 记录单个房间已配置的设施列表
//   - 支持房间设施的灵活勾选
//   - 前端展示房间设施时查询此表
type RoomFacilityRelation struct {
	// ID 关联ID，主键
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:关联ID" json:"id"`
	
	// RoomID 房源ID，外键关联 room_info 表
	RoomID uint64 `gorm:"column:room_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_room_id;comment:房源ID（外键，关联房源表）" json:"room_id"`
	
	// FacilityID 设施ID，外键关联 facility_dict 表
	FacilityID uint64 `gorm:"column:facility_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_facility_id;comment:设施ID（外键，关联设施字典表）" json:"facility_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	Room     *RoomInfo     `gorm:"foreignKey:RoomID;references:ID" json:"room,omitempty"`         // 房间信息
	Facility *FacilityDict `gorm:"foreignKey:FacilityID;references:ID" json:"facility,omitempty"` // 设施信息
}

// TableName 指定数据库表名
func (RoomFacilityRelation) TableName() string {
	return "room_facility_relation"
}
