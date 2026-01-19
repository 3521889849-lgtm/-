// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了系统用户账号的数据模型
//
// 功能说明：
//   - 存储系统用户的账号信息
//   - 支持账号管理和权限控制
//   - 支持角色绑定和分店绑定
//   - 密码采用bcrypt加密存储
//   - 记录登录历史
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 用户账号模型 ====================

// UserAccount 用户账号表
//
// 业务用途：
//   - 系统登录：前台、管理员登录系统的凭证
//   - 权限管理：通过角色控制用户的操作权限
//   - 分店管理：用户可以绑定到特定分店，实现分店级数据隔离
//   - 操作审计：记录用户的操作历史，追溯责任
//   - 账号管理：支持账号的增删改查、启用/停用
//
// 设计说明：
//   - 用户名全局唯一，作为登录凭证
//   - 密码使用bcrypt加密存储，不可逆，安全性高
//   - 每个用户绑定一个角色，角色定义了权限集合
//   - 管理员可以不绑定分店（BranchID为NULL）
//   - 普通员工必须绑定分店，只能操作本分店数据
//
// 安全说明：
//   - Password 字段使用 bcrypt 加密存储
//   - json标签设置为 "-"，API响应中不返回密码字段
//   - 登录时使用 bcrypt.CompareHashAndPassword 验证密码
//   - 修改密码时使用 bcrypt.GenerateFromPassword 生成新密码
type UserAccount struct {
	// ========== 基础字段 ==========
	
	// ID 账号ID，主键，自增
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:账号ID" json:"id"`
	
	// ========== 登录凭证 ==========
	
	// Username 用户名，全局唯一标识，用于登录
	// 规则：
	//   - 长度：4-50字符
	//   - 字符：字母、数字、下划线
	//   - 唯一性：全系统唯一，不能重复
	// 示例："admin"、"zhangsan"、"beijing_frontdesk"
	// 用途：登录、展示、查询
	Username string `gorm:"column:username;type:VARCHAR(50);NOT NULL;uniqueIndex:uk_username;comment:用户名" json:"username"`
	
	// Password 密码，使用bcrypt加密存储
	// ⚠️ 安全要求：
	//   - 存储前必须使用 bcrypt.GenerateFromPassword 加密
	//   - 验证时使用 bcrypt.CompareHashAndPassword 对比
	//   - 不可逆加密，即使数据库泄露也无法还原原始密码
	//   - API响应中不返回密码字段（json:"-"）
	//
	// 密码规则建议：
	//   - 长度：8-20字符
	//   - 复杂度：包含大小写字母、数字、特殊字符
	//   - 定期更换：建议90天更换一次
	//
	// 示例（加密后）：
	//   "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
	Password string `gorm:"column:password;type:VARCHAR(255);NOT NULL;comment:密码（加密存储）" json:"-"`
	
	// ========== 用户信息 ==========
	
	// RealName 真实姓名，用户的真实姓名
	// 用途：
	//   - 系统显示：前台显示操作员姓名
	//   - 操作日志：记录操作人姓名
	//   - 沟通交流：同事间相互识别
	// 示例："张三"、"李经理"
	RealName string `gorm:"column:real_name;type:VARCHAR(50);NOT NULL;comment:姓名" json:"real_name"`
	
	// ContactPhone 联系电话，用户的手机号或座机号
	// 用途：
	//   - 紧急联系：通知重要信息
	//   - 身份验证：找回密码、修改敏感信息
	//   - 短信通知：系统通知、任务提醒
	// 格式：手机号（11位）或座机号（带区号）
	// 示例："13812345678"、"010-12345678"
	ContactPhone string `gorm:"column:contact_phone;type:VARCHAR(20);NOT NULL;index:idx_contact_phone;comment:联系电话" json:"contact_phone"`
	
	// ========== 角色权限 ==========
	
	// RoleID 角色ID，外键关联 role 表
	// 用途：
	//   - 权限控制：角色定义了用户可以执行的操作
	//   - 批量授权：同一角色的用户拥有相同的权限
	//   - 灵活管理：修改角色权限可影响所有该角色的用户
	//
	// 角色示例：
	//   - 超级管理员：所有权限，不受分店限制
	//   - 分店管理员：本分店所有权限
	//   - 前台接待：入住、退房、收款
	//   - 财务人员：查看财务报表、对账
	//   - 客房部：房态管理、清洁状态更新
	RoleID uint64 `gorm:"column:role_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_role_id;comment:角色ID（外键，关联角色表）" json:"role_id"`
	
	// ========== 分店绑定 ==========
	
	// BranchID 分店ID，外键关联 hotel_branch 表，可选
	// NULL: 超级管理员或总部员工，可以跨分店操作
	// 非NULL: 绑定到特定分店，只能操作本分店数据
	//
	// 数据隔离规则：
	//   - 绑定分店的用户：
	//     - 只能查看和操作本分店的数据
	//     - 不能查看其他分店的数据
	//     - 登录后自动切换到本分店
	//   - 未绑定分店的用户：
	//     - 可以查看所有分店的数据
	//     - 可以在不同分店间切换
	//     - 通常是系统管理员或总部人员
	//
	// 用途：实现多分店管理、数据隔离、权限控制
	BranchID *uint64 `gorm:"column:branch_id;type:BIGINT UNSIGNED;index:idx_branch_id;comment:分店ID（外键，关联分店表，管理员可为空）" json:"branch_id,omitempty"`
	
	// ========== 状态控制 ==========
	
	// Status 账号状态，控制账号是否可用
	// 可选值：
	//   - "ACTIVE": 启用，可以正常登录使用（默认状态）
	//   - "INACTIVE": 停用，禁止登录
	//
	// 停用场景：
	//   - 员工离职：离职后立即停用账号
	//   - 长期休假：长期请假期间临时停用
	//   - 安全风险：账号异常时紧急停用
	//   - 违规操作：违反规定时暂停账号
	//
	// 停用影响：
	//   - 无法登录系统
	//   - 已登录的会话会被踢出
	//   - 数据不会删除，可以随时启用恢复
	Status string `gorm:"column:status;type:VARCHAR(20);NOT NULL;default:'ACTIVE';index:idx_status;comment:账号状态（启用/停用）" json:"status"`
	
	// ========== 时间戳 ==========
	
	// CreatedAt 创建时间，账号首次创建的时间
	// 用途：统计分析、账号排序
	CreatedAt time.Time `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	
	// LastLoginAt 最后登录时间，用户最近一次登录的时间，可选
	// NULL表示从未登录
	// 非NULL表示有登录记录
	// 用途：
	//   - 安全监控：长期未登录的账号可能被盗用
	//   - 活跃度分析：判断账号的使用频率
	//   - 自动清理：长期不用的账号可以自动停用
	LastLoginAt *time.Time `gorm:"column:last_login_at;type:DATETIME;index:idx_last_login_at;comment:最后登录时间" json:"last_login_at,omitempty"`
	
	// UpdatedAt 更新时间，账号信息最后修改的时间
	// 自动更新：每次修改账号信息时自动更新此字段
	UpdatedAt time.Time `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	
	// DeletedAt 软删除时间，非NULL表示已删除
	// 软删除：账号注销时不真正删除数据，只是标记
	// 好处：数据可恢复，操作日志完整性保证
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// ========== 关联关系 ==========
	// 以下字段不会存储在数据库中，仅用于GORM的关联查询
	
	// Role 关联的角色信息
	// 多对一关系：多个用户可以属于同一个角色
	// 用途：获取用户的权限列表
	Role *Role `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
	
	// Branch 关联的分店信息
	// 多对一关系（可选）：多个用户可以属于同一个分店
	// 用途：显示用户所属分店、数据隔离控制
	Branch *HotelBranch `gorm:"foreignKey:BranchID;references:ID" json:"branch,omitempty"`
	
	// OperationLogs 该用户的操作日志记录
	// 一对多关系：一个用户可以有多条操作日志
	// 用途：审计追溯、责任认定、行为分析
	OperationLogs []OperationLog `gorm:"foreignKey:OperatorID;references:ID" json:"operation_logs,omitempty"`
}

// ==================== 表名配置 ====================

// TableName 指定数据库表名
//
// GORM会自动调用此方法获取表名，用于生成SQL语句
//
// 返回：
//   - string: 数据库表名 "user_account"
func (UserAccount) TableName() string {
	return "user_account"
}
