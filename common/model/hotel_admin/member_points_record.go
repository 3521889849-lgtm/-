// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了会员积分记录的数据模型
//
// 功能说明：
//   - 记录会员积分的所有变动明细
//   - 支持积分获取和消费的全流程追溯
//   - 支持积分统计和分析
//   - 与订单关联，实现积分与消费的对应
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 会员积分记录模型 ====================

// MemberPointsRecord 会员积分记录表
//
// 业务用途：
//   - 积分明细：记录每一笔积分的来源和去向
//   - 积分追溯：查询积分的获取和使用历史
//   - 对账审计：确保积分余额与明细一致
//   - 统计分析：分析积分的使用情况、活跃度等
//   - 权益兑换：记录积分兑换礼品或服务的明细
//
// 设计说明：
//   - 每次积分变动（增加或减少）都生成一条记录
//   - ChangePoints 字段为正数表示增加，负数表示减少
//   - 与订单关联，记录积分来源
//   - 支持软删除，保证数据完整性
//   - 记录操作人，实现责任追溯
type MemberPointsRecord struct {
	// ========== 基础字段 ==========
	
	// ID 记录ID，主键，自增
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:记录ID" json:"id"`
	
	// ========== 会员关联 ==========
	
	// MemberID 会员ID，外键关联 member 表
	// 用途：标识这笔积分属于哪个会员
	MemberID uint64 `gorm:"column:member_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_member_id;comment:会员ID（外键，关联会员表）" json:"member_id"`
	
	// ========== 订单关联 ==========
	
	// OrderID 订单ID，外键关联 hotel_order_main 表，可选
	// NULL: 非订单相关的积分变动（如：活动赠送、手动调整、礼品兑换）
	// 非NULL: 订单相关的积分变动（如：消费获取积分、积分抵扣房费）
	// 用途：追溯积分来源，对账验证
	OrderID *uint64 `gorm:"column:order_id;type:BIGINT UNSIGNED;index:idx_order_id;comment:订单ID（外键，关联订单主表）" json:"order_id,omitempty"`
	
	// ========== 积分变动信息 ==========
	
	// ChangeType 积分变动类型，标识积分增加还是减少
	// 可选值：
	//   - "EARN": 获取积分（增加）
	//   - "CONSUME": 消费积分（减少）
	//   - "REFUND": 积分退还（增加，如：取消订单退还已扣积分）
	//   - "EXPIRE": 积分过期（减少）
	//   - "ADJUST": 人工调整（可增可减）
	// 用途：区分积分来源和去向，统计分析
	ChangeType string `gorm:"column:change_type;type:VARCHAR(20);NOT NULL;index:idx_change_type;comment:积分变动类型（获取/消费）" json:"change_type"`
	
	// ChangePoints 变动积分值
	// 规则：
	//   - 正数：积分增加（如：+100 表示获得100积分）
	//   - 负数：积分减少（如：-50 表示消费50积分）
	//   - 零：理论上不应出现，如出现可能是数据错误
	//
	// 计算：会员积分余额 = 历史所有ChangePoints的累加和
	// 说明：使用int64而不是uint64，是为了支持负数表示减少
	ChangePoints int64 `gorm:"column:change_points;type:BIGINT;NOT NULL;comment:变动积分值" json:"change_points"`
	
	// ChangeReason 变动原因，详细说明这笔积分变动的原因
	// 获取积分示例：
	//   - "订单ORD20260115001消费获得积分"
	//   - "新会员注册赠送积分"
	//   - "生日礼包赠送积分"
	//   - "推荐好友奖励积分"
	//   - "参与活动获得积分"
	//
	// 消费积分示例：
	//   - "订单ORD20260115002积分抵扣房费"
	//   - "兑换礼品消耗积分"
	//   - "兑换优惠券消耗积分"
	//
	// 其他示例：
	//   - "订单ORD20260115003取消退还积分"
	//   - "积分过期自动清零"
	//   - "人工调整：补偿客人积分"
	//
	// 用途：审计追溯、客户查询、问题排查
	ChangeReason string `gorm:"column:change_reason;type:VARCHAR(255);NOT NULL;comment:变动原因" json:"change_reason"`
	
	// ========== 时间和操作信息 ==========
	
	// ChangeTime 变动时间，积分变动实际发生的时间
	// 用途：
	//   - 时间排序：按时间顺序查询积分明细
	//   - 统计分析：按时间段统计积分获取和使用情况
	//   - 过期判断：计算积分的有效期
	ChangeTime time.Time `gorm:"column:change_time;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;index:idx_change_time;comment:变动时间" json:"change_time"`
	
	// OperatorID 操作人ID，关联 user_account 表
	// 场景：
	//   - 自动触发：订单完成自动增加积分，操作人为系统账号
	//   - 手动调整：管理员手动调整积分，操作人为管理员
	//   - 会员操作：会员自助兑换礼品，操作人为前台或系统
	// 用途：责任追溯、操作审计
	OperatorID uint64 `gorm:"column:operator_id;type:BIGINT UNSIGNED;NOT NULL;comment:操作人" json:"operator_id"`
	
	// ========== 时间戳 ==========
	
	// CreatedAt 创建时间，记录创建的时间
	// 说明：通常与ChangeTime相同，但对于补录的记录可能不同
	CreatedAt time.Time `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	
	// DeletedAt 软删除时间，非NULL表示已删除
	// 说明：积分记录原则上不应删除，仅在特殊情况下（如：重复记录）才软删除
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// ========== 关联关系 ==========
	
	// Member 关联的会员信息
	// 用途：获取会员的姓名、等级等信息
	Member *Member `gorm:"foreignKey:MemberID;references:ID" json:"member,omitempty"`
	
	// Order 关联的订单信息（可选）
	// 用途：追溯积分来源，查看订单详情
	Order *OrderMain `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
}

// ==================== 表名配置 ====================

// TableName 指定数据库表名
//
// 返回：数据库表名 "member_points_record"
func (MemberPointsRecord) TableName() string {
	return "member_points_record"
}
