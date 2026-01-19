// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了房间图片的数据模型
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// RoomImage 房源图片表
//
// 业务用途：
//   - 存储房间的展示图片（最多16张）
//   - 支持图片排序，控制展示顺序
//   - 记录图片规格和格式
//   - 前端房间详情页展示图片轮播
type RoomImage struct {
	// ID 图片ID，主键
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:图片ID" json:"id"`
	
	// RoomID 房源ID，外键关联 room_info 表
	RoomID uint64 `gorm:"column:room_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_room_id;comment:房源ID（外键，关联房源表）" json:"room_id"`
	
	// ImageURL 图片URL，图片的访问地址
	// 可以是相对路径或完整URL
	ImageURL string `gorm:"column:image_url;type:VARCHAR(512);NOT NULL;comment:图片URL" json:"image_url"`
	
	// ImageSize 图片规格，默认："400x300"
	ImageSize string `gorm:"column:image_size;type:VARCHAR(20);NOT NULL;default:'400x300';comment:图片规格（400x300）" json:"image_size"`
	
	// ImageFormat 图片格式，如："jpg"、"png"
	ImageFormat string `gorm:"column:image_format;type:VARCHAR(10);NOT NULL;default:'jpg';comment:图片格式（jpg/png）" json:"image_format"`
	
	// SortOrder 排序序号，数字越小越靠前，控制展示顺序
	SortOrder uint8 `gorm:"column:sort_order;type:TINYINT UNSIGNED;NOT NULL;default:0;index:idx_sort_order;comment:排序序号" json:"sort_order"`
	
	// UploadTime 上传时间，图片上传的时间
	UploadTime time.Time `gorm:"column:upload_time;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:上传时间" json:"upload_time"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	Room *RoomInfo `gorm:"foreignKey:RoomID;references:ID" json:"room,omitempty"` // 房间信息
}

// TableName 指定数据库表名
func (RoomImage) TableName() string {
	return "room_image"
}
