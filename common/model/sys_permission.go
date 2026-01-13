package model

import "time"

// SysPermission 系统权限表
type SysPermission struct {
	ID        int       `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	PermName  string    `gorm:"type:varchar(50);not null;comment:权限名称（如：房源管理-查询、财务流水-导出）" json:"perm_name"`
	PermCode  string    `gorm:"type:varchar(100);not null;uniqueIndex:idx_perm_code;comment:权限编码（如：ROOM:QUERY、FINANCE:EXPORT）" json:"perm_code"`
	PermType  int8      `gorm:"type:tinyint(1);not null;comment:权限类型（0-菜单，1-按钮）" json:"perm_type"`
	ParentID  *int      `gorm:"index:idx_parent_id;comment:父权限ID（关联自身id）" json:"parent_id"`
	CreatedAt time.Time `gorm:"not null;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;comment:更新时间" json:"updated_at"`

	// 关联关系
	Parent   *SysPermission   `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"parent,omitempty"`
	Children []SysPermission  `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Roles    []SysRole        `gorm:"many2many:sys_role_perm_rel;" json:"roles,omitempty"`
}

// TableName 指定表名
func (SysPermission) TableName() string {
	return "sys_permission"
}
