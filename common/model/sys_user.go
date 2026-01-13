package model

import (
	"time"

	"gorm.io/gorm"
)

// SysUser 系统用户账号表
type SysUser struct {
	ID            int            `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	Username      string         `gorm:"type:varchar(50);not null;uniqueIndex:idx_username;comment:用户名（唯一）" json:"username"`
	Password      string         `gorm:"type:varchar(100);not null;comment:密码（加密存储）" json:"password"`
	RealName      string         `gorm:"type:varchar(50);comment:真实姓名" json:"real_name"`
	Phone         string         `gorm:"type:varchar(20);comment:手机号" json:"phone"`
	BranchID      *int           `gorm:"index:idx_branch_id;comment:关联分店ID（hotel_branch.id，超级管理员为空）" json:"branch_id"`
	Status        int8           `gorm:"type:tinyint(1);not null;default:1;comment:状态（0-停用，1-启用）" json:"status"`
	LastLoginTime *time.Time     `gorm:"comment:最后登录时间" json:"last_login_time"`
	CreatedAt     time.Time      `gorm:"not null;comment:创建时间" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null;comment:更新时间" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index:idx_deleted_at;comment:软删除时间" json:"deleted_at"`

	// 关联关系
	Branch         *HotelBranch     `gorm:"foreignKey:BranchID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"branch,omitempty"`
	Roles          []SysRole        `gorm:"many2many:sys_user_role_rel;" json:"roles,omitempty"`
	OperationLogs  []SysOperationLog `gorm:"foreignKey:UserID" json:"operation_logs,omitempty"`
}

// TableName 指定表名
func (SysUser) TableName() string {
	return "sys_user"
}
