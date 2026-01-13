package model

import "time"

// SysUserRoleRel 用户-角色关联表
type SysUserRoleRel struct {
	UserID    int       `gorm:"primaryKey;comment:关联用户ID（sys_user.id）" json:"user_id"`
	RoleID    int       `gorm:"primaryKey;index:idx_role_id;comment:关联角色ID（sys_role.id）" json:"role_id"`
	CreatedAt time.Time `gorm:"not null;comment:创建时间" json:"created_at"`

	// 关联关系
	User *SysUser `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	Role *SysRole `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"role,omitempty"`
}

// TableName 指定表名
func (SysUserRoleRel) TableName() string {
	return "sys_user_role_rel"
}
