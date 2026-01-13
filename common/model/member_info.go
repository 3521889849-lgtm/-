package model

import (
	"time"

	"gorm.io/gorm"
)

// MemberInfo 会员信息表
type MemberInfo struct {
	ID              int            `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	GuestID         int            `gorm:"not null;uniqueIndex:idx_guest_id;comment:关联客人ID（guest_info.id）" json:"guest_id"`
	LevelID         int            `gorm:"not null;index:idx_level_id;comment:关联等级ID（member_level.id）" json:"level_id"`
	MemberNo        string         `gorm:"type:varchar(50);not null;uniqueIndex:idx_member_no;comment:会员编号（唯一）" json:"member_no"`
	TotalPoints     int            `gorm:"not null;default:0;comment:总积分" json:"total_points"`
	UsedPoints      int            `gorm:"not null;default:0;comment:已使用积分" json:"used_points"`
	RegisterTime    time.Time      `gorm:"not null;comment:注册时间" json:"register_time"`
	LastCheckinTime *time.Time     `gorm:"comment:最后入住时间" json:"last_checkin_time"`
	Status          int8           `gorm:"type:tinyint(1);not null;default:1;comment:状态（0-冻结，1-正常）" json:"status"`
	CreatedAt       time.Time      `gorm:"not null;comment:创建时间" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"not null;comment:更新时间" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index:idx_deleted_at;comment:软删除时间" json:"deleted_at"`

	// 关联关系
	Guest *GuestInfo   `gorm:"foreignKey:GuestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"guest,omitempty"`
	Level *MemberLevel `gorm:"foreignKey:LevelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"level,omitempty"`
}

// TableName 指定表名
func (MemberInfo) TableName() string {
	return "member_info"
}
