// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了财务流水的数据模型
//
// 功能说明：
//   - 记录所有财务收支明细（收入和支出）
//   - 支撑财务流水查询和统计分析
//   - 支持多维度统计（按日期、按分店、按收支类型等）
//   - 与订单、客人、房源关联，实现完整的财务追溯
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 财务流水模型 ====================

// FinancialFlow 收支流水表
//
// 业务用途：
//   - 财务管理：记录所有收入和支出明细
//   - 统计分析：支持按各种维度统计（日报、月报、年报）
//   - 审计追溯：每笔收支都有完整记录，可追溯到操作人
//   - 对账结算：与支付渠道对账、财务月结
//   - 经营分析：分析收入构成、支出结构、盈亏情况
//
// 设计说明：
//   - 每笔收支都生成一条流水记录
//   - 支持软删除，保证财务数据完整性
//   - 关联订单、客人、房源，方便多维度查询
//   - 记录操作人，实现责任追溯
type FinancialFlow struct {
	// ========== 基础字段 ==========
	
	// ID 流水ID，主键，自增
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:流水ID" json:"id"`
	
	// ========== 关联外键 ==========
	
	// OrderID 订单ID，外键关联 hotel_order_main 表，可选
	// NULL: 非订单相关的收支（如：杂项收入、日常支出）
	// 非NULL: 订单相关的收支（如：房费、押金）
	OrderID *uint64 `gorm:"column:order_id;type:BIGINT UNSIGNED;index:idx_order_id;comment:订单ID（外键，关联订单主表）" json:"order_id,omitempty"`
	
	// BranchID 分店ID，外键关联 hotel_branch 表，必填
	// 用途：实现分店级财务统计和数据隔离
	BranchID uint64 `gorm:"column:branch_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_branch_id;comment:分店ID（外键，关联分店表）" json:"branch_id"`
	
	// RoomID 房源ID，外键关联 room_info 表，可选
	// 用途：统计单个房间的收入情况
	RoomID *uint64 `gorm:"column:room_id;type:BIGINT UNSIGNED;index:idx_room_id;comment:房源ID（外键，关联房源表）" json:"room_id,omitempty"`
	
	// GuestID 客人ID，外键关联 guest_info 表，可选
	// 用途：统计客人的消费记录、欠款情况
	GuestID *uint64 `gorm:"column:guest_id;type:BIGINT UNSIGNED;index:idx_guest_id;comment:客人ID（外键，关联客人表）" json:"guest_id,omitempty"`
	
	// ========== 收支信息 ==========
	
	// FlowType 收支类型，标识这笔流水是收入还是支出
	// 可选值：
	//   - "INCOME": 收入（如：房费、加床费、损坏赔偿）
	//   - "EXPENSE": 支出（如：退款、赔偿、采购）
	// 用途：区分收支方向，计算盈亏
	FlowType string `gorm:"column:flow_type;type:VARCHAR(20);NOT NULL;index:idx_flow_type;comment:收支类型（收入/支出）" json:"flow_type"`
	
	// FlowItem 收支项目，详细的收支科目
	// 收入项目示例：
	//   - "ROOM_FEE": 房费（主要收入）
	//   - "DEPOSIT": 押金（临时收入，退房时退还）
	//   - "PENALTY": 违约金（取消订单收取）
	//   - "DAMAGE": 损坏赔偿
	//   - "EXTRA_BED": 加床费
	//   - "BREAKFAST": 早餐费
	//   - "LAUNDRY": 洗衣费
	//   - "MINIBAR": 迷你吧消费
	//   - "OTHER": 其他收入
	//
	// 支出项目示例：
	//   - "REFUND": 退款（取消订单或多收退还）
	//   - "DEPOSIT_REFUND": 押金退还
	//   - "COMPENSATION": 赔偿（如：服务问题赔偿客人）
	//   - "PROCUREMENT": 采购
	//   - "SALARY": 工资
	//   - "UTILITY": 水电费
	//   - "MAINTENANCE": 维修费
	//   - "OTHER": 其他支出
	//
	// 用途：财务科目统计、收入构成分析
	FlowItem string `gorm:"column:flow_item;type:VARCHAR(50);NOT NULL;index:idx_flow_item;comment:收支项目（房费/押金/违约金等）" json:"flow_item"`
	
	// PayType 支付方式，标识支付渠道
	// 可选值：
	//   - "CASH": 现金
	//   - "ALIPAY": 支付宝
	//   - "WECHAT": 微信支付
	//   - "UNIONPAY": 银联卡
	//   - "CREDIT_CARD": 信用卡
	//   - "DEBIT_CARD": 储蓄卡
	//   - "ACCOUNT": 挂账
	//   - "OTHER": 其他方式
	// 用途：对账、支付渠道分析
	PayType string `gorm:"column:pay_type;type:VARCHAR(20);NOT NULL;index:idx_pay_type;comment:支付方式" json:"pay_type"`
	
	// Amount 金额，收支金额
	// 单位：元（人民币）
	// 规则：始终为正数，通过FlowType区分收入/支出
	// 示例：298.00（表示298元）
	Amount float64 `gorm:"column:amount;type:DECIMAL(10,2);NOT NULL;comment:金额" json:"amount"`
	
	// ========== 时间和操作信息 ==========
	
	// OccurTime 发生时间，财务事件实际发生的时间
	// 用途：
	//   - 财务统计：按日期范围统计
	//   - 日夜审计：确定流水归属日期
	//   - 报表生成：日报、月报、年报
	OccurTime time.Time `gorm:"column:occur_time;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;index:idx_occur_time;comment:发生时间" json:"occur_time"`
	
	// OperatorID 操作人ID，关联 user_account 表
	// 用途：责任追溯、操作审计、业绩统计
	OperatorID uint64 `gorm:"column:operator_id;type:BIGINT UNSIGNED;NOT NULL;comment:操作人" json:"operator_id"`
	
	// Remark 备注，补充说明信息，可选
	// 用途：记录特殊情况、异常处理、补充信息
	// 示例："客人损坏电视遥控器赔偿"、"618活动优惠券抵扣"
	Remark *string `gorm:"column:remark;type:VARCHAR(500);comment:备注" json:"remark,omitempty"`
	
	// ========== 时间戳 ==========
	
	// CreatedAt 创建时间，流水记录创建的时间
	// 说明：通常与OccurTime相同，但对于补录的流水可能不同
	CreatedAt time.Time `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	
	// DeletedAt 软删除时间，非NULL表示已删除
	// 说明：财务流水原则上不应删除，仅在特殊情况下（如：录入错误）才软删除
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// ========== 关联关系 ==========
	
	// Order 关联的订单信息（可选）
	Order *OrderMain `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
	
	// Branch 所属分店信息
	Branch *HotelBranch `gorm:"foreignKey:BranchID;references:ID" json:"branch,omitempty"`
	
	// Room 关联的房源信息（可选）
	Room *RoomInfo `gorm:"foreignKey:RoomID;references:ID" json:"room,omitempty"`
	
	// Guest 关联的客人信息（可选）
	Guest *GuestInfo `gorm:"foreignKey:GuestID;references:ID" json:"guest,omitempty"`
}

// ==================== 表名配置 ====================

// TableName 指定数据库表名
//
// 返回：数据库表名 "financial_flow"
func (FinancialFlow) TableName() string {
	return "financial_flow"
}
