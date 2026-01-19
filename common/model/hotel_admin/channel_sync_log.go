// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了渠道同步日志的数据模型
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ChannelSyncLog 渠道数据同步日志表
//
// 业务用途：
//   - 记录与第三方渠道的数据同步情况
//   - 追溯同步历史和失败原因
//   - 支持同步失败后的重试
//   - 监控同步成功率
type ChannelSyncLog struct {
	ID         uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:日志ID" json:"id"`
	ChannelID  uint64 `gorm:"column:channel_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_channel_id;comment:渠道ID（外键，关联渠道配置表）" json:"channel_id"`
	
	// SyncType 同步类型，如："订单同步"、"房态同步"、"价格同步"
	SyncType   string `gorm:"column:sync_type;type:VARCHAR(30);NOT NULL;index:idx_sync_type;comment:同步类型（订单同步/房态同步）" json:"sync_type"`
	SyncDataID uint64 `gorm:"column:sync_data_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_sync_data_id;comment:同步数据ID（如订单ID/房源ID）" json:"sync_data_id"`
	
	// SyncStatus 同步状态，"SUCCESS"（成功）或"FAILED"（失败）
	SyncStatus string    `gorm:"column:sync_status;type:VARCHAR(20);NOT NULL;index:idx_sync_status;comment:同步状态（成功/失败）" json:"sync_status"`
	SyncTime   time.Time `gorm:"column:sync_time;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;index:idx_sync_time;comment:同步时间" json:"sync_time"`
	FailReason *string   `gorm:"column:fail_reason;type:VARCHAR(500);comment:失败原因" json:"fail_reason,omitempty"`
	RetryCount uint8     `gorm:"column:retry_count;type:TINYINT UNSIGNED;NOT NULL;default:0;comment:重试次数" json:"retry_count"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	Channel *ChannelConfig `gorm:"foreignKey:ChannelID;references:ID" json:"channel,omitempty"` // 渠道配置
}

// TableName 指定数据库表名
func (ChannelSyncLog) TableName() string {
	return "channel_sync_log"
}
