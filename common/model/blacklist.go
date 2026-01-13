package model

import "time"

// Blacklist 黑名单表
type Blacklist struct {
	ID              int        `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	GuestID         int        `gorm:"not null;uniqueIndex:idx_guest_id;comment:关联客人ID（guest_info.id）" json:"guest_id"`
	BlacklistReason string     `gorm:"type:text;comment:拉黑原因" json:"blacklist_reason"`
	BlacklistTime   time.Time  `gorm:"not null;index:idx_blacklist_time;comment:拉黑时间" json:"blacklist_time"`
	ExpireTime      *time.Time `gorm:"index:idx_expire_time;comment:过期时间（NULL表示永久拉黑）" json:"expire_time"`
	Operator        string     `gorm:"type:varchar(50);not null;comment:操作人（用户名）" json:"operator"`
	CreatedAt       time.Time  `gorm:"not null;comment:创建时间" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"not null;comment:更新时间" json:"updated_at"`

	// 关联关系
	Guest *GuestInfo `gorm:"foreignKey:GuestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"guest,omitempty"`
}

// TableName 指定表名
func (Blacklist) TableName() string {
	return "blacklist"
}
