package model

import (
	"time"

	"gorm.io/gorm"
)

// SysAdmin 管理员信息表-管理端所有操作的执行人
type SysAdmin struct {
	ID        uint64         `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:管理员主键ID" json:"id"`
	AdminName string         `gorm:"column:admin_name;type:VARCHAR(50);NOT NULL;comment:管理员姓名" json:"admin_name"`
	Account   string         `gorm:"column:account;type:VARCHAR(50);NOT NULL;uniqueIndex:uk_account;comment:登录账号" json:"account"`
	Password  string         `gorm:"column:password;type:VARCHAR(64);NOT NULL;comment:登录密码（加密存储）" json:"-"`
	Phone     string         `gorm:"column:phone;type:VARCHAR(20);NOT NULL;comment:联系电话" json:"phone"`
	Role      string         `gorm:"column:role;type:VARCHAR(20);NOT NULL;comment:角色：SUPER-超级管理员，OPERATOR-运营管理员" json:"role"`
	Status    uint8          `gorm:"column:status;type:TINYINT UNSIGNED;NOT NULL;default:1;comment:状态：0-禁用，1-启用" json:"status"`
	CreatedAt time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	Merchants []SysMerchant `gorm:"foreignKey:AdminID;references:ID" json:"merchants,omitempty"`
	OperLogs  []SysOperLog  `gorm:"foreignKey:OperAdminID;references:ID" json:"oper_logs,omitempty"`
}

func (SysAdmin) TableName() string {
	return "sys_admin"
}
