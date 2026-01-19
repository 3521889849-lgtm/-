// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了交接班记录的数据模型
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ShiftChangeRecord 交接班记录表
//
// 业务用途：
//   - 记录工作人员交接班的财务数据
//   - 确保账目可追溯，责任明确
//   - 支持财务对账和审计
//   - 防止账目遗漏或错误
type ShiftChangeRecord struct {
	ID           uint64    `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:记录ID" json:"id"`
	BranchID     uint64    `gorm:"column:branch_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_branch_id;comment:分店ID（外键，关联分店表）" json:"branch_id"`
	HandoverID   uint64    `gorm:"column:handover_id;type:BIGINT UNSIGNED;NOT NULL;comment:交班人ID" json:"handover_id"`
	TakeoverID   uint64    `gorm:"column:takeover_id;type:BIGINT UNSIGNED;NOT NULL;comment:接班人ID" json:"takeover_id"`
	HandoverTime time.Time `gorm:"column:handover_time;type:DATETIME;NOT NULL;index:idx_handover_time;comment:交班时间" json:"handover_time"`
	TakeoverTime time.Time `gorm:"column:takeover_time;type:DATETIME;NOT NULL;index:idx_takeover_time;comment:接班时间" json:"takeover_time"`
	
	// 当班财务数据
	ShiftIncome     float64 `gorm:"column:shift_income;type:DECIMAL(10,2);NOT NULL;default:0.00;comment:当班收入总额" json:"shift_income"`
	ShiftExpense    float64 `gorm:"column:shift_expense;type:DECIMAL(10,2);NOT NULL;default:0.00;comment:当班支出总额" json:"shift_expense"`
	ShiftOrderCount uint32  `gorm:"column:shift_order_count;type:INT UNSIGNED;NOT NULL;default:0;comment:当班订单数" json:"shift_order_count"`
	
	// ReconciliationStatus 账目核对状态
	// 可选值："PENDING"（待核对）、"CONFIRMED"（已确认）、"DISPUTED"（有争议）
	ReconciliationStatus string  `gorm:"column:reconciliation_status;type:VARCHAR(20);NOT NULL;default:'PENDING';index:idx_reconciliation_status;comment:账目核对状态" json:"reconciliation_status"`
	Remark               *string `gorm:"column:remark;type:VARCHAR(500);comment:备注" json:"remark,omitempty"`
	
	// 时间戳
	CreatedAt time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	Branch *HotelBranch `gorm:"foreignKey:BranchID;references:ID" json:"branch,omitempty"` // 分店信息
}

// TableName 指定数据库表名
func (ShiftChangeRecord) TableName() string {
	return "shift_change_record"
}
