// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了角色-权限关联的数据模型，实现多对多关系
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// RolePermissionRelation 角色-权限关联表
//
// 业务用途：
//   - 实现角色和权限的多对多关联
//   - 一个角色可以拥有多个权限
//   - 一个权限可以分配给多个角色
//   - 修改角色权限时，只需修改此表的记录
type RolePermissionRelation struct {
	// ID 关联ID，主键
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:关联ID" json:"id"`
	
	// RoleID 角色ID，外键关联 role 表
	RoleID uint64 `gorm:"column:role_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_role_id;comment:角色ID（外键，关联角色表）" json:"role_id"`
	
	// PermissionID 权限ID，外键关联 permission 表
	PermissionID uint64 `gorm:"column:permission_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_permission_id;comment:权限ID（外键，关联权限表）" json:"permission_id"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	Role       *Role       `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`             // 角色信息
	Permission *Permission `gorm:"foreignKey:PermissionID;references:ID" json:"permission,omitempty"` // 权限信息
}

// TableName 指定数据库表名
func (RolePermissionRelation) TableName() string {
	return "role_permission_relation"
}
