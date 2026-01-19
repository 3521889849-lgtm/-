// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了订单主表的数据模型
//
// 功能说明：
//   - 存储订单的核心信息（客人、房间、时间、金额等）
//   - 支持订单的全生命周期管理（预订-入住-退房）
//   - 支持订单状态跟踪和查询
//   - 支持多种客人来源和支付方式
//   - 与客人、房源、财务等模块紧密关联
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 订单主表模型 ====================

// OrderMain 订单主表
//
// 业务用途：
//   - 记录客人的订房信息
//   - 作为整个订单流程的核心数据
//   - 关联房源、客人、支付、财务等多个业务模块
//   - 支持订单状态流转（预订→入住→退房）
//   - 支持订单财务管理（订金、欠款、违约金）
//
// 设计说明：
//   - 订单号全局唯一，用于对外展示和查询
//   - 支持多种客人来源（OTA渠道和自有渠道）
//   - 支持软删除，保证历史数据完整性
//   - 记录操作人，方便审计追溯
type OrderMain struct {
	// ========== 基础字段 ==========
	
	// ID 订单ID，主键，自增
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:订单ID" json:"id"`
	
	// OrderNo 订单号，全局唯一标识
	// 格式示例："ORD20260115123456"（订单+日期+序号）
	// 用途：对外展示、客人查询、财务对账
	OrderNo string `gorm:"column:order_no;type:VARCHAR(32);NOT NULL;uniqueIndex:uk_order_no;comment:订单号" json:"order_no"`
	
	// ========== 关联外键 ==========
	
	// BranchID 分店ID，外键关联 hotel_branch 表
	// 用途：订单所属分店，实现分店级数据隔离
	BranchID uint64 `gorm:"column:branch_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_branch_id;comment:分店ID（外键，关联分店表）" json:"branch_id"`
	
	// GuestID 客人ID，外键关联 guest_info 表
	// 用途：关联客人的实名制登记信息
	GuestID uint64 `gorm:"column:guest_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_guest_id;comment:客人ID（外键，关联客人表）" json:"guest_id"`
	
	// RoomID 房源ID，外键关联 room_info 表
	// 用途：记录客人预订的具体房间
	RoomID uint64 `gorm:"column:room_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_room_id;comment:房源ID（外键，关联房源表）" json:"room_id"`
	
	// RoomTypeID 房型ID，外键关联 room_type_dict 表
	// 用途：冗余字段，方便统计分析和报表查询
	// 原因：避免关联查询room_info表才能得到房型
	RoomTypeID uint64 `gorm:"column:room_type_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_room_type_id;comment:房型ID（外键，关联房型表）" json:"room_type_id"`
	
	// ========== 订单来源 ==========
	
	// GuestSource 客人来源，标识订单的渠道
	// 可选值：
	//   - "WALK_IN": 散客（前台直接登记）
	//   - "CTRIP": 携程
	//   - "TUYOU": 途游
	//   - "ELONG": 艺龙
	//   - "MEITUAN": 美团
	//   - "QUNAR": 去哪儿
	//   - "PHONE": 电话预订
	//   - "OFFICIAL": 官网预订
	// 用途：渠道分析、佣金结算、营销效果评估
	GuestSource string `gorm:"column:guest_source;type:VARCHAR(50);NOT NULL;index:idx_guest_source;comment:客人来源（散客/携程/途游/艺龙等）" json:"guest_source"`
	
	// ========== 时间信息 ==========
	
	// CheckInTime 入住时间，客人计划入住的时间
	// 用途：房态管理、房间分配、日夜审计
	CheckInTime time.Time `gorm:"column:check_in_time;type:DATETIME;NOT NULL;index:idx_check_in_time;comment:入住时间" json:"check_in_time"`
	
	// CheckOutTime 离店时间，客人计划退房的时间
	// 用途：房态管理、房间清洁安排、超时计费
	CheckOutTime time.Time `gorm:"column:check_out_time;type:DATETIME;NOT NULL;index:idx_check_out_time;comment:离店时间" json:"check_out_time"`
	
	// ReserveTime 预订时间，客人下单的时间
	// 用途：订单时效性判断、取消政策计算
	ReserveTime time.Time `gorm:"column:reserve_time;type:DATETIME;NOT NULL;comment:预定时间" json:"reserve_time"`
	
	// ========== 财务信息 ==========
	
	// OrderAmount 订单金额，客人应付的总金额
	// 单位：元（人民币）
	// 计算：房费 × 天数 + 其他费用
	// 用途：财务对账、收银结算
	OrderAmount float64 `gorm:"column:order_amount;type:DECIMAL(10,2);NOT NULL;comment:订单金额" json:"order_amount"`
	
	// DepositReceived 已收押金，客人已支付的押金金额
	// 单位：元（人民币）
	// 用途：退房时退还或抵扣
	// 说明：通常用于防止客人损坏物品或消费
	DepositReceived float64 `gorm:"column:deposit_received;type:DECIMAL(10,2);NOT NULL;default:0.00;comment:已收押金" json:"deposit_received"`
	
	// OutstandingAmount 欠补费用，客人尚未支付的金额
	// 单位：元（人民币）
	// 场景：
	//   - 预订时部分支付，入住时补齐
	//   - 入住期间产生额外消费
	//   - 退房时统一结算
	OutstandingAmount float64 `gorm:"column:outstanding_amount;type:DECIMAL(10,2);NOT NULL;default:0.00;comment:欠补费用" json:"outstanding_amount"`
	
	// ========== 订单状态 ==========
	
	// OrderStatus 订单状态，标识订单的当前阶段
	// 状态流转：
	//   - RESERVED: 已预订（客人下单但未入住）
	//   - CHECKED_IN: 已入住（客人已办理入住）
	//   - CHECKED_OUT: 已退房（客人已办理退房）
	//   - CANCELLED: 已取消（客人取消预订）
	//   - EXPIRED: 已失效（超时未入住，系统自动失效）
	// 用途：订单管理、房态更新、统计分析
	OrderStatus string `gorm:"column:order_status;type:VARCHAR(30);NOT NULL;default:'RESERVED';index:idx_order_status;comment:订单状态（已预定/已入住/已退房/已失效）" json:"order_status"`
	
	// ========== 支付方式 ==========
	
	// PayType 支付方式，客人使用的支付手段
	// 可选值：
	//   - "CASH": 现金
	//   - "ALIPAY": 支付宝
	//   - "WECHAT": 微信支付
	//   - "UNIONPAY": 银联卡
	//   - "CREDIT_CARD": 信用卡
	//   - "DEBIT_CARD": 储蓄卡
	//   - "ACCOUNT": 挂账（企业协议）
	// 用途：财务对账、支付渠道分析
	PayType string `gorm:"column:pay_type;type:VARCHAR(20);NOT NULL;index:idx_pay_type;comment:支付方式（现金/支付宝/微信/银联等）" json:"pay_type"`
	
	// ========== 退订政策 ==========
	
	// CancellationPolicyID 退订政策ID，外键关联 cancellation_policy 表
	// 可选字段，如果为NULL则使用房间的默认政策
	// 用途：判断客人取消订单时是否收取违约金
	CancellationPolicyID *uint64 `gorm:"column:cancellation_policy_id;type:BIGINT UNSIGNED;index:idx_cancellation_policy_id;comment:退订政策ID（外键，关联退订政策表）" json:"cancellation_policy_id,omitempty"`
	
	// PenaltyAmount 违约金金额，客人取消订单时应收的违约金
	// 单位：元（人民币）
	// 计算：根据退订政策和取消时间自动计算
	// 场景：客人在不允许免费取消的时间段内取消订单
	PenaltyAmount float64 `gorm:"column:penalty_amount;type:DECIMAL(10,2);NOT NULL;default:0.00;comment:违约金金额" json:"penalty_amount"`
	
	// ========== 时间戳 ==========
	
	// CreatedAt 创建时间，订单首次创建的时间
	// 用途：统计分析、订单排序
	CreatedAt time.Time `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;index:idx_create_time;comment:创建时间" json:"created_at"`
	
	// UpdatedAt 修改时间，订单最后修改的时间
	// 用途：变更追踪、同步判断
	UpdatedAt time.Time `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:修改时间" json:"updated_at"`
	
	// OperatorID 操作人ID，关联用户账号表
	// 用途：审计追溯，记录是谁创建/修改的订单
	OperatorID uint64 `gorm:"column:operator_id;type:BIGINT UNSIGNED;NOT NULL;comment:操作人" json:"operator_id"`
	
	// DeletedAt 软删除时间，非NULL表示已删除
	// 软删除：数据不会真正删除，只是标记为已删除状态
	// 好处：数据可恢复，财务数据完整性保证
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// ========== 关联关系 ==========
	// 以下字段不会存储在数据库中，仅用于GORM的关联查询
	
	// Branch 所属分店信息
	Branch *HotelBranch `gorm:"foreignKey:BranchID;references:ID" json:"branch,omitempty"`
	
	// Guest 客人信息
	Guest *GuestInfo `gorm:"foreignKey:GuestID;references:ID" json:"guest,omitempty"`
	
	// Room 房源信息
	Room *RoomInfo `gorm:"foreignKey:RoomID;references:ID" json:"room,omitempty"`
	
	// RoomType 房型信息
	RoomType *RoomTypeDict `gorm:"foreignKey:RoomTypeID;references:ID" json:"room_type,omitempty"`
	
	// CancellationPolicy 退订政策信息
	CancellationPolicy *CancellationPolicy `gorm:"foreignKey:CancellationPolicyID;references:ID" json:"cancellation_policy,omitempty"`
	
	// OrderExtension 订单扩展信息（如：特殊要求、备注等）
	// 一对一关系
	OrderExtension *OrderExtension `gorm:"foreignKey:OrderID;references:ID" json:"order_extension,omitempty"`
	
	// FinancialFlows 该订单相关的财务流水记录
	// 一对多关系：一个订单可以有多条财务记录（如：房费、押金、退款）
	FinancialFlows []FinancialFlow `gorm:"foreignKey:OrderID;references:ID" json:"financial_flows,omitempty"`
	
	// MemberPointsRecords 会员积分记录
	// 一对多关系：一个订单可以产生多条积分记录（如：消费积分、积分抵扣）
	MemberPointsRecords []MemberPointsRecord `gorm:"foreignKey:OrderID;references:ID" json:"member_points_records,omitempty"`
}

// ==================== 表名配置 ====================

// TableName 指定数据库表名
//
// GORM会自动调用此方法获取表名，用于生成SQL语句
//
// 返回：
//   - string: 数据库表名 "hotel_order_main"
func (OrderMain) TableName() string {
	return "hotel_order_main"
}
