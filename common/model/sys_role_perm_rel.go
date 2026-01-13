package model

import "time"

// SysRolePermRel 角色-权限关联表
type SysRolePermRel struct {
	RoleID    int       `gorm:"primaryKey;comment:关联角色ID（sys_role.id）" json:"role_id"`
	PermID    int       `gorm:"primaryKey;index:idx_perm_id;comment:关联权限ID（sys_permission.id）" json:"perm_id"`
	CreatedAt time.Time `gorm:"not null;comment:创建时间" json:"created_at"`

	// 关联关系
	Role       *SysRole       `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"role,omitempty"`
	Permission *SysPermission `gorm:"foreignKey:PermID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"permission,omitempty"`
}

// TableName 指定表名
func (SysRolePermRel) TableName() string {
	return "sys_role_perm_rel"
}
