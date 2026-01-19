// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了房态明细的数据模型，按日期维度存储房态数据
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// RoomStatusDetail 房态明细表
//
// 业务用途：
//   - 按日期维度记录房间状态
//   - 支撑日历化房态展示（房态图）
//   - 实时统计各状态房间数量
//   - 支持房态历史查询
//
// 房态类型：
//   - 空净房：已清洁，可入住
//   - 入住房：客人已入住
//   - 维修房：维修中，不可用
//   - 锁定房：临时锁定，不可预订
//   - 空账房：客人已退房但未结账
//   - 预定房：已预订但未入住
type RoomStatusDetail struct {
	ID     uint64    `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:记录ID" json:"id"`
	RoomID uint64    `gorm:"column:room_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_room_id;comment:房源ID（外键，关联房源表）" json:"room_id"`
	Date   time.Time `gorm:"column:date;type:DATE;NOT NULL;index:idx_date;comment:日期（YYYY-MM-DD）" json:"date"`
	
	// RoomStatus 房态类型
	RoomStatus string `gorm:"column:room_status;type:VARCHAR(20);NOT NULL;index:idx_room_status;comment:房态（空净房/入住房/维修房/锁定房/空账房/预定房）" json:"room_status"`
	
	// 统计字段
	RemainingCount       uint8 `gorm:"column:remaining_count;type:TINYINT UNSIGNED;NOT NULL;default:0;comment:当日剩余数量" json:"remaining_count"`
	CheckedInCount       uint8 `gorm:"column:checked_in_count;type:TINYINT UNSIGNED;NOT NULL;default:0;comment:已入住人数" json:"checked_in_count"`
	CheckOutPendingCount uint8 `gorm:"column:check_out_pending_count;type:TINYINT UNSIGNED;NOT NULL;default:0;comment:预退房人数" json:"check_out_pending_count"`
	ReservedPendingCount uint8 `gorm:"column:reserved_pending_count;type:TINYINT UNSIGNED;NOT NULL;default:0;comment:预定待入住人数" json:"reserved_pending_count"`
	
	// 时间戳
	UpdatedAt time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;index:idx_updated_at;comment:更新时间" json:"updated_at"`
	CreatedAt time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	Room *RoomInfo `gorm:"foreignKey:RoomID;references:ID" json:"room,omitempty"`
}

// TableName 指定数据库表名
func (RoomStatusDetail) TableName() string {
	return "room_status_detail"
}
