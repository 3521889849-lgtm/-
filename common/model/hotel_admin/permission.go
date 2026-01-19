// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了权限的数据模型，支撑基于角色的权限管理（RBAC）
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// Permission 权限表 - 维护系统所有权限项，支撑角色权限分配
//
// 业务用途：
//   - 权限定义：定义系统中所有可操作的权限
//   - 权限分配：通过角色关联权限，实现批量授权
//   - 权限控制：前端根据权限显示菜单和按钮
//   - 权限树：支持父子级权限，构建权限树结构
type Permission struct {
	// ID 权限ID，主键
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:权限ID" json:"id"`
	
	// PermissionName 权限名称，如："收支流水查询"、"会员编辑"
	PermissionName string `gorm:"column:permission_name;type:VARCHAR(100);NOT NULL;index:idx_permission_name;comment:权限名称（收支流水查询/会员编辑等）" json:"permission_name"`
	
	// PermissionURL 权限URL，前端路由或API路径
	// 示例："/financial-flows"、"/api/v1/members"
	PermissionURL string `gorm:"column:permission_url;type:VARCHAR(255);NOT NULL;comment:权限URL" json:"permission_url"`
	
	// PermissionType 权限类型
	// 可选值："MENU"（菜单）、"BUTTON"（按钮/操作）
	PermissionType string `gorm:"column:permission_type;type:VARCHAR(20);NOT NULL;index:idx_permission_type;comment:权限类型（菜单/按钮）" json:"permission_type"`
	
	// ParentID 父权限ID，构建权限树，顶级权限为NULL
	ParentID *uint64 `gorm:"column:parent_id;type:BIGINT UNSIGNED;index:idx_parent_id;comment:父权限ID" json:"parent_id,omitempty"`
	
	// Status 状态，控制权限是否可用
	Status string `gorm:"column:status;type:VARCHAR(20);NOT NULL;default:'ACTIVE';index:idx_status;comment:状态：ACTIVE-启用，INACTIVE-停用" json:"status"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	Parent                  *Permission              `gorm:"foreignKey:ParentID;references:ID" json:"parent,omitempty"`                         // 父权限
	Children                []Permission             `gorm:"foreignKey:ParentID;references:ID" json:"children,omitempty"`                       // 子权限列表
	RolePermissionRelations []RolePermissionRelation `gorm:"foreignKey:PermissionID;references:ID" json:"role_permission_relations,omitempty"` // 角色权限关联
}

// TableName 指定数据库表名
func (Permission) TableName() string {
	return "permission"
}
