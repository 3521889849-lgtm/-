package model

import (
	"time"

	"gorm.io/gorm"
)

// RoomInfo 房间信息表
type RoomInfo struct {
	ID           int            `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	BranchID     int            `gorm:"not null;index:idx_branch_room_no,priority:1;comment:关联分店ID（hotel_branch.id）" json:"branch_id"`
	RoomTypeID   int            `gorm:"not null;index:idx_room_type_id;comment:关联房型ID（room_type.id）" json:"room_type_id"`
	RoomNumber   string         `gorm:"type:varchar(20);not null;index:idx_branch_room_no,priority:2;comment:房间号（如：8101、8302）" json:"room_number"`
	ParentRoomID *int           `gorm:"index:idx_parent_room_id;comment:关联父房间ID（用于多房间绑定，关联自身id）" json:"parent_room_id"`
	Status       int8           `gorm:"type:tinyint(1);not null;default:1;comment:状态（0-停用，1-空净房，2-入住房，3-维修房，4-锁定房，5-空账房）" json:"status"`
	SortOrder    int            `gorm:"not null;default:0;comment:排序序号（支持拖动排序）" json:"sort_order"`
	CreatedAt    time.Time      `gorm:"not null;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"not null;comment:更新时间" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index:idx_deleted_at;comment:软删除时间" json:"deleted_at"`

	// 关联关系
	Branch     *HotelBranch `gorm:"foreignKey:BranchID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"branch,omitempty"`
	RoomType   *RoomType    `gorm:"foreignKey:RoomTypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"room_type,omitempty"`
	ParentRoom *RoomInfo    `gorm:"foreignKey:ParentRoomID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"parent_room,omitempty"`
	ChildRooms []RoomInfo   `gorm:"foreignKey:ParentRoomID" json:"child_rooms,omitempty"`
}

// TableName 指定表名
func (RoomInfo) TableName() string {
	return "room_info"
}
