// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了酒店分店信息的数据模型
//
// 功能说明：
//   - 支持多分店管理，每个酒店可以有多个分店
//   - 提供分店的基本信息管理（名称、地址、联系方式等）
//   - 支持分店的启用/停用状态控制
//   - 支持软删除，数据安全可恢复
//   - 关联房源、订单、财务等业务数据
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 酒店分店模型 ====================

// HotelBranch 酒店/分店信息表
//
// 业务用途：
//   - 管理酒店的多个分店信息
//   - 支持分店切换功能
//   - 数据隔离：每个分店的房源、订单、财务数据独立管理
//   - 权限控制：用户可以绑定到特定分店，实现分店级别的权限管理
//
// 设计说明：
//   - 使用分店编码（branch_code）作为唯一标识
//   - 支持软删除，删除后数据仍保留在数据库中
//   - 记录创建人和创建时间，方便追溯
type HotelBranch struct {
	// ========== 基础字段 ==========
	
	// ID 分店ID，主键，自增
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:分店ID" json:"id"`
	
	// HotelName 酒店名称，用于显示和搜索
	// 示例："锦江之星"、"如家快捷"
	HotelName string `gorm:"column:hotel_name;type:VARCHAR(100);NOT NULL;index:idx_hotel_name;comment:酒店名称" json:"hotel_name"`
	
	// BranchCode 分店编码，全局唯一标识
	// 示例："BJ001"（北京分店）、"SH002"（上海分店）
	// 用于：系统集成、数据同步、报表统计等场景
	BranchCode string `gorm:"column:branch_code;type:VARCHAR(50);NOT NULL;uniqueIndex:uk_branch_code;comment:分店编码" json:"branch_code"`
	
	// Address 分店地址，完整的详细地址
	// 示例："北京市朝阳区建国路88号"
	Address string `gorm:"column:address;type:VARCHAR(255);NOT NULL;comment:地址" json:"address"`
	
	// ========== 联系信息 ==========
	
	// Contact 联系人姓名，通常是分店经理或负责人
	Contact string `gorm:"column:contact;type:VARCHAR(50);NOT NULL;comment:联系人" json:"contact"`
	
	// ContactPhone 联系电话，用于客户咨询和内部沟通
	// 格式：手机号或座机号
	ContactPhone string `gorm:"column:contact_phone;type:VARCHAR(20);NOT NULL;comment:联系电话" json:"contact_phone"`
	
	// ========== 状态控制 ==========
	
	// Status 分店状态，控制分店是否可用
	// 可选值：
	//   - ACTIVE: 启用，正常营业
	//   - INACTIVE: 停用，暂停营业（如装修、节假日关闭等）
	// 用途：停用的分店不允许新增订单，但可以查看历史数据
	Status string `gorm:"column:status;type:VARCHAR(20);NOT NULL;default:'ACTIVE';index:idx_status;comment:状态：ACTIVE-启用，INACTIVE-停用" json:"status"`
	
	// ========== 时间戳 ==========
	
	// CreatedAt 创建时间，记录分店信息首次创建的时间
	CreatedAt time.Time `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	
	// CreatedBy 创建人ID，关联用户账号表
	// 用于审计追溯，记录是谁创建的这个分店
	CreatedBy uint64 `gorm:"column:created_by;type:BIGINT UNSIGNED;NOT NULL;comment:创建人" json:"created_by"`
	
	// UpdatedAt 更新时间，每次修改分店信息时自动更新
	UpdatedAt time.Time `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	
	// DeletedAt 软删除时间，非NULL表示已删除
	// 软删除：数据不会真正删除，只是标记为已删除状态
	// 好处：数据可恢复，历史数据完整性保证
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// ========== 关联关系 ==========
	// 以下字段不会存储在数据库中，仅用于GORM的关联查询
	
	// RoomInfos 该分店下的所有房源信息
	// 一对多关系：一个分店可以有多个房源
	RoomInfos []RoomInfo `gorm:"foreignKey:BranchID;references:ID" json:"room_infos,omitempty"`
	
	// Orders 该分店下的所有订单
	// 一对多关系：一个分店可以有多个订单
	Orders []OrderMain `gorm:"foreignKey:BranchID;references:ID" json:"orders,omitempty"`
	
	// FinancialFlows 该分店的财务流水记录
	// 一对多关系：一个分店可以有多条财务流水
	FinancialFlows []FinancialFlow `gorm:"foreignKey:BranchID;references:ID" json:"financial_flows,omitempty"`
	
	// SalesStats 该分店的销售统计数据
	// 一对多关系：一个分店可以有多条统计记录（按日、按月等）
	SalesStats []SalesStatistics `gorm:"foreignKey:BranchID;references:ID" json:"sales_stats,omitempty"`
	
	// ShiftRecords 该分店的交接班记录
	// 一对多关系：一个分店可以有多条交接班记录
	ShiftRecords []ShiftChangeRecord `gorm:"foreignKey:BranchID;references:ID" json:"shift_records,omitempty"`
	
	// UserAccounts 该分店的用户账号
	// 一对多关系：一个分店可以有多个员工账号
	UserAccounts []UserAccount `gorm:"foreignKey:BranchID;references:ID" json:"user_accounts,omitempty"`
}

// ==================== 表名配置 ====================

// TableName 指定数据库表名
//
// GORM会自动调用此方法获取表名，用于生成SQL语句
//
// 返回：
//   - string: 数据库表名 "hotel_branch"
func (HotelBranch) TableName() string {
	return "hotel_branch"
}
