package model

import "time"

// SysConfig 系统配置表
type SysConfig struct {
	ID          int       `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	ConfigKey   string    `gorm:"type:varchar(100);not null;uniqueIndex:idx_config_key;comment:配置键（如：SMS_TEMPLATE_CHECKIN、PRINT_SETTING_ROOM）" json:"config_key"`
	ConfigValue string    `gorm:"type:text;comment:配置值" json:"config_value"`
	ConfigDesc  string    `gorm:"type:varchar(255);comment:配置描述" json:"config_desc"`
	CreatedAt   time.Time `gorm:"not null;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;comment:更新时间" json:"updated_at"`
}

// TableName 指定表名
func (SysConfig) TableName() string {
	return "sys_config"
}
