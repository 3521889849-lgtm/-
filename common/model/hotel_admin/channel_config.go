// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了渠道配置的数据模型，支持与第三方OTA平台对接
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ChannelConfig 渠道配置表
//
// 业务用途：
//   - 配置第三方预订渠道（携程、美团、去哪儿等）
//   - 支持数据同步功能（房态、价格、订单）
//   - 统一管理多个渠道的对接参数
//   - 灵活控制渠道的启用/停用
type ChannelConfig struct {
	// ID 配置ID，主键
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:配置ID" json:"id"`
	
	// ChannelName 渠道名称，如："携程"、"美团"、"去哪儿"
	ChannelName string `gorm:"column:channel_name;type:VARCHAR(50);NOT NULL;uniqueIndex:uk_channel_name;comment:渠道名称（携程/途游/艺龙等）" json:"channel_name"`
	
	// ChannelCode 渠道编码，全局唯一标识
	// 示例："CTRIP"、"MEITUAN"、"QUNAR"
	ChannelCode string `gorm:"column:channel_code;type:VARCHAR(50);NOT NULL;uniqueIndex:uk_channel_code;comment:渠道编码" json:"channel_code"`
	
	// ApiURL 对接接口URL，第三方平台的API地址
	ApiURL string `gorm:"column:api_url;type:VARCHAR(255);NOT NULL;comment:对接接口URL" json:"api_url"`
	
	// SyncRule 同步规则
	// 可选值："REALTIME"（实时同步）、"SCHEDULED"（定时同步）
	SyncRule string `gorm:"column:sync_rule;type:VARCHAR(20);NOT NULL;default:'REALTIME';index:idx_sync_rule;comment:同步规则（实时/定时）" json:"sync_rule"`
	
	// Status 状态，控制渠道是否启用
	Status string `gorm:"column:status;type:VARCHAR(20);NOT NULL;default:'ACTIVE';index:idx_status;comment:状态（启用/停用）" json:"status"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:修改时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	SyncLogs []ChannelSyncLog `gorm:"foreignKey:ChannelID;references:ID" json:"sync_logs,omitempty"` // 同步日志记录
}

// TableName 指定数据库表名
func (ChannelConfig) TableName() string {
	return "channel_config"
}
