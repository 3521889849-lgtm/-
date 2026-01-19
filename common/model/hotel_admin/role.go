// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了角色的数据模型
//
// 功能说明：
//   - 支撑基于角色的权限管理（RBAC）
//   - 支持角色的增删改查
//   - 支持角色与权限的关联
//   - 支持角色与用户账号的关联
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 角色模型 ====================

// Role 角色表
//
// 业务用途：
//   - 权限管理：通过角色统一管理权限，避免为每个用户单独分配权限
//   - 分级管理：不同角色拥有不同的操作权限
//   - 批量授权：同一角色的所有用户自动拥有该角色的权限
//   - 灵活调整：修改角色权限可影响所有该角色的用户
//   - 职责分离：实现不同岗位的权限隔离
//
// 设计说明：
//   - 角色名称全局唯一
//   - 角色通过中间表（role_permission_relation）关联权限
//   - 用户通过 RoleID 字段关联角色
//   - 支持角色启用/停用，停用后该角色的用户暂时失去权限
//   - 支持软删除，保证数据完整性
//
// 典型角色示例：
//   - 超级管理员：拥有所有权限，不受分店限制
//   - 分店管理员：拥有本分店所有权限
//   - 前台接待：入住登记、退房结算、房态查看
//   - 财务人员：财务查询、财务统计、对账结算
//   - 客房部：房态管理、清洁状态、维修报修
//   - 会员专员：会员管理、积分管理、权益配置
//   - 营销人员：订单查询、客源分析、营销活动
//   - 系统审计员：查看所有日志、生成审计报告
type Role struct {
	// ========== 基础字段 ==========
	
	// ID 角色ID，主键，自增
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:角色ID" json:"id"`
	
	// ========== 角色信息 ==========
	
	// RoleName 角色名称，全局唯一标识
	// 命名规范：
	//   - 简洁明了：如"前台接待"、"财务经理"
	//   - 反映职责：名称应体现角色的主要职责
	//   - 便于理解：避免使用缩写或代号
	//
	// 示例：
	//   - "超级管理员": 系统最高权限
	//   - "分店经理": 分店管理权限
	//   - "前台接待": 入住退房权限
	//   - "财务人员": 财务查询权限
	//   - "客房部主管": 房态管理权限
	//   - "会员专员": 会员管理权限
	//
	// 用途：显示、查询、权限判断
	RoleName string `gorm:"column:role_name;type:VARCHAR(50);NOT NULL;uniqueIndex:uk_role_name;comment:角色名称（财务经理/管理员等）" json:"role_name"`
	
	// Description 角色描述，详细说明角色的职责和权限范围，可选
	// 用途：
	//   - 帮助管理员理解角色定位
	//   - 说明角色的主要职责
	//   - 列举角色的主要权限
	//
	// 示例：
	//   - "系统最高权限，可以管理所有功能和数据，不受分店限制"
	//   - "负责本分店的日常运营管理，包括人员管理、财务管理、客房管理等"
	//   - "负责办理客人入住和退房手续，处理日常接待工作"
	//   - "负责财务数据查询、财务报表生成、账务对账等工作"
	Description *string `gorm:"column:description;type:VARCHAR(255);comment:角色描述" json:"description,omitempty"`
	
	// ========== 状态控制 ==========
	
	// Status 角色状态，控制角色是否可用
	// 可选值：
	//   - "ACTIVE": 启用，角色正常使用（默认状态）
	//   - "INACTIVE": 停用，角色暂停使用
	//
	// 停用场景：
	//   - 角色废弃：某些角色不再需要时临时停用
	//   - 权限调整：大规模调整权限前临时停用
	//   - 安全管控：发现安全问题时紧急停用
	//
	// 停用影响：
	//   - 该角色的所有用户暂时失去权限
	//   - 无法为新用户分配该角色
	//   - 数据不会删除，可以随时启用恢复
	//   - 用户仍然可以登录，但操作会受限
	Status string `gorm:"column:status;type:VARCHAR(20);NOT NULL;default:'ACTIVE';index:idx_status;comment:状态：ACTIVE-启用，INACTIVE-停用" json:"status"`
	
	// ========== 时间戳 ==========
	
	// CreatedAt 创建时间，角色首次创建的时间
	CreatedAt time.Time `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	
	// UpdatedAt 更新时间，角色最后修改的时间
	// 自动更新：每次修改角色信息或权限配置时自动更新此字段
	// 用途：追踪角色变更历史
	UpdatedAt time.Time `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	
	// DeletedAt 软删除时间，非NULL表示已删除
	// 软删除：角色废弃时不真正删除数据，只是标记
	// 好处：数据可恢复，历史记录完整性保证
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// ========== 关联关系 ==========
	
	// UserAccounts 拥有该角色的所有用户账号
	// 一对多关系：一个角色可以分配给多个用户
	// 用途：查询该角色下的所有用户
	UserAccounts []UserAccount `gorm:"foreignKey:RoleID;references:ID" json:"user_accounts,omitempty"`
	
	// RolePermissionRelations 角色-权限关联关系
	// 一对多关系：一个角色可以拥有多个权限
	// 用途：查询该角色拥有的所有权限
	// 说明：通过中间表实现多对多关系
	RolePermissionRelations []RolePermissionRelation `gorm:"foreignKey:RoleID;references:ID" json:"role_permission_relations,omitempty"`
}

// ==================== 表名配置 ====================

// TableName 指定数据库表名
//
// 返回：数据库表名 "role"
func (Role) TableName() string {
	return "role"
}
