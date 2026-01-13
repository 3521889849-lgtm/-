package model

import (
	"time"

	"gorm.io/gorm"
)

// RoomTypeFacilityRel 房型-设施关联表
type RoomTypeFacilityRel struct {
	RoomTypeID int            `gorm:"primaryKey;comment:关联房型ID（room_type.id）" json:"room_type_id"`
	FacilityID int            `gorm:"primaryKey;index:idx_facility_id;comment:关联设施ID（room_facility.id）" json:"facility_id"`
	CreatedAt  time.Time      `gorm:"not null;comment:创建时间" json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index:idx_deleted_at;comment:软删除时间" json:"deleted_at"`

	// 关联关系
	RoomType *RoomType     `gorm:"foreignKey:RoomTypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"room_type,omitempty"`
	Facility *RoomFacility `gorm:"foreignKey:FacilityID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"facility,omitempty"`
}

// TableName 指定表名
func (RoomTypeFacilityRel) TableName() string {
	return "room_type_facility_rel"
}
