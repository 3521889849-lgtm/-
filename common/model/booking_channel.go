package model

import "time"

// BookingChannel 预订渠道表
type BookingChannel struct {
	ID          int       `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	ChannelName string    `gorm:"type:varchar(50);not null;comment:渠道名称（如：散客、携程、途游、艺龙）" json:"channel_name"`
	ChannelCode string    `gorm:"type:varchar(20);not null;uniqueIndex:idx_channel_code;comment:渠道编码（如：DIRECT、CTRIP、TUYOU）" json:"channel_code"`
	SyncStatus  int8      `gorm:"type:tinyint(1);not null;default:1;comment:数据同步状态（0-不同步，1-同步）" json:"sync_status"`
	CreatedAt   time.Time `gorm:"not null;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;comment:更新时间" json:"updated_at"`

	// 关联关系
	Orders []OrderInfo `gorm:"foreignKey:ChannelID" json:"orders,omitempty"`
}

// TableName 指定表名
func (BookingChannel) TableName() string {
	return "booking_channel"
}
