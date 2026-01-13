package model

import "time"

// RoomFacility 房间设施表
type RoomFacility struct {
	ID           int       `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	FacilityName string    `gorm:"type:varchar(50);not null;comment:设施名称（如：无线wifi、空调、冰箱）" json:"facility_name"`
	FacilityDesc string    `gorm:"type:varchar(255);comment:设施描述" json:"facility_desc"`
	Status       int8      `gorm:"type:tinyint(1);not null;default:1;comment:状态（0-禁用，1-启用）" json:"status"`
	CreatedAt    time.Time `gorm:"not null;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;comment:更新时间" json:"updated_at"`

	// 关联关系
	RoomTypes []RoomType `gorm:"many2many:room_type_facility_rel;" json:"room_types,omitempty"`
}

// TableName 指定表名
func (RoomFacility) TableName() string {
	return "room_facility"
}
