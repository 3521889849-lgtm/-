package model

import "time"

// RoomImage 房源图片表（≤16张，规格400x300，格式jpg/png）
type RoomImage struct {
	ID         int       `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	RoomTypeID int       `gorm:"not null;index:idx_room_type_id;comment:关联房型ID（room_type.id）" json:"room_type_id"`
	ImageURL   string    `gorm:"type:varchar(255);not null;comment:图片URL" json:"image_url"`
	ImageSort  int       `gorm:"not null;default:0;comment:图片排序" json:"image_sort"`
	CreatedAt  time.Time `gorm:"not null;comment:创建时间" json:"created_at"`
	UpdatedAt  time.Time `gorm:"not null;comment:更新时间" json:"updated_at"`

	// 关联关系
	RoomType *RoomType `gorm:"foreignKey:RoomTypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"room_type,omitempty"`
}

// TableName 指定表名
func (RoomImage) TableName() string {
	return "room_image"
}
