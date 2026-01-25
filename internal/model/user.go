package model

import (
	"time"

	"gorm.io/gorm"
)

// UserInfo 用户表
type UserInfo struct {
	ID               string         `gorm:"column:user_id;type:VARCHAR(64);primaryKey;comment:用户唯一标识（如手机号、用户ID）" json:"user_id"`
	Password         string         `gorm:"column:password;type:VARCHAR(128);NOT NULL;comment:加密密码" json:"-"`
	RealName         string         `gorm:"column:real_name;type:VARCHAR(32);NOT NULL;comment:真实姓名（实名信息）" json:"real_name"`
	IDCard           *string        `gorm:"column:id_card;type:VARCHAR(32);DEFAULT NULL;uniqueIndex:idx_id_card;comment:身份证号（脱敏存储，如1426****414）" json:"id_card,omitempty"`
	Phone            string         `gorm:"column:phone;type:VARCHAR(16);NOT NULL;uniqueIndex;comment:手机号" json:"phone"`
	UserLevel        string         `gorm:"column:user_level;type:VARCHAR(20);NOT NULL;default:'ORDINARY';index:idx_user_level;comment:用户等级：ORDINARY-普通用户，VIP-VIP用户" json:"user_level"`
	RealNameVerified string         `gorm:"column:real_name_verified;type:VARCHAR(20);NOT NULL;default:'UNVERIFIED';comment:实名验证状态：UNVERIFIED-未验证，VERIFIED-已验证" json:"real_name_verified"`
	Status           string         `gorm:"column:status;type:VARCHAR(20);NOT NULL;default:'NORMAL';comment:用户状态：NORMAL-正常，DISABLED-禁用" json:"status"`
	CreatedAt        time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`
}

// PassengerInfo 乘客表
type PassengerInfo struct {
	ID        uint64         `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:自增主键" json:"id"`
	OrderID   string         `gorm:"column:order_id;type:VARCHAR(64);NOT NULL;comment:订单号，关联order_info.order_id" json:"order_id"`
	UserID    string         `gorm:"column:user_id;type:VARCHAR(64);NOT NULL;comment:关联用户ID，关联user_info.user_id" json:"user_id"`
	RealName  string         `gorm:"column:real_name;type:VARCHAR(32);NOT NULL;comment:乘客真实姓名" json:"real_name"`
	IDCard    string         `gorm:"column:id_card;type:VARCHAR(32);NOT NULL;comment:乘客身份证号（脱敏存储）" json:"id_card"`
	SeatID    string         `gorm:"column:seat_id;type:VARCHAR(64);NOT NULL;comment:乘客关联座位ID，关联seat_info.seat_id" json:"seat_id"`
	CreatedAt time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`
}
