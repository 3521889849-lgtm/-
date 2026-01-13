package model

import "time"

// SysRole 系统角色表
type SysRole struct {
	ID        int       `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	RoleName  string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_role_name;comment:角色名称（如：财务经理、前台接待、超级管理员）" json:"role_name"`
	RoleDesc  string    `gorm:"type:varchar(255);comment:角色描述" json:"role_desc"`
	Status    int8      `gorm:"type:tinyint(1);not null;default:1;comment:状态（0-禁用，1-启用）" json:"status"`
	CreatedAt time.Time `gorm:"not null;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;comment:更新时间" json:"updated_at"`

	// 关联关系
	Users       []SysUser       `gorm:"many2many:sys_user_role_rel;" json:"users,omitempty"`
	Permissions []SysPermission `gorm:"many2many:sys_role_perm_rel;" json:"permissions,omitempty"`
}

// TableName 指定表名
func (SysRole) TableName() string {
	return "sys_role"
}
