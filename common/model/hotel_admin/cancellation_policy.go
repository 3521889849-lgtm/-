// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了退订政策的数据模型
//
// 功能说明：
//   - 存储自定义的退订规则
//   - 关联房源或房型生效
//   - 支持灵活配置违约金比例
//   - 保护酒店权益，规范退订流程
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 退订政策模型 ====================

// CancellationPolicy 退订政策表
//
// 业务用途：
//   - 规则管理：定义不同的退订规则
//   - 风险控制：通过违约金降低客人随意取消的风险
//   - 灵活配置：不同房型可以使用不同的退订政策
//   - 客人知情：预订时明确告知退订规则
//   - 纠纷处理：退订纠纷时有据可依
//
// 设计说明：
//   - 可以配置多个退订政策
//   - 可以关联到房型（该房型的所有房间使用此政策）
//   - 也可以直接关联到具体房间（覆盖房型的政策）
//   - 支持启用/停用，方便策略调整
type CancellationPolicy struct {
	// ========== 基础字段 ==========
	
	// ID 政策ID，主键，自增
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:政策ID" json:"id"`
	
	// ========== 政策信息 ==========
	
	// PolicyName 政策名称，政策的简短标识
	// 示例：
	//   - "灵活退订"：任何时候免费取消
	//   - "标准退订"：入住前24小时免费取消
	//   - "严格退订"：入住前7天免费取消
	//   - "不可退订"：预订后不可取消
	//   - "节假日政策"：节假日特殊政策
	// 用途：政策识别、前台展示、客人选择
	PolicyName string `gorm:"column:policy_name;type:VARCHAR(100);NOT NULL;comment:政策名称" json:"policy_name"`
	
	// RuleDescription 规则描述，政策的详细说明
	// 内容要求：
	//   - 清晰说明何时可以免费取消
	//   - 说明何时取消需要收取违约金
	//   - 说明违约金的计算方式
	//   - 特殊情况的处理方式
	//
	// 示例：
	//   - "入住前24小时免费取消，24小时内取消收取首晚房费作为违约金"
	//   - "入住前7天免费取消，7天内取消收取全额房费的50%作为违约金"
	//   - "入住前3天免费取消，3天内取消收取全额房费作为违约金"
	//   - "预订后不可取消，取消将收取全额房费"
	//   - "任何时候均可免费取消（特价房、活动房除外）"
	//
	// 用途：客人知情、前台执行、纠纷处理
	RuleDescription string `gorm:"column:rule_description;type:VARCHAR(500);NOT NULL;comment:规则描述（如：入住前X小时内不可取消）" json:"rule_description"`
	
	// PenaltyRatio 违约金比例，取消订单时收取的违约金比例
	// 单位：倍数（相对于房费）
	// 取值范围：
	//   - 0.0: 免费取消，不收取违约金
	//   - 0.5: 收取50%房费作为违约金
	//   - 1.0: 收取100%房费作为违约金（全额）
	//   - 2.0: 收取200%房费作为违约金（双倍赔偿，极少使用）
	//
	// 计算公式：
	//   违约金 = 订单金额 × PenaltyRatio
	//
	// 示例：
	//   - 订单金额298元，PenaltyRatio=0.5，违约金=149元
	//   - 订单金额298元，PenaltyRatio=1.0，违约金=298元
	//
	// 用途：自动计算违约金、财务结算
	PenaltyRatio float64 `gorm:"column:penalty_ratio;type:DECIMAL(5,2);NOT NULL;comment:违约金比例（X倍房费）" json:"penalty_ratio"`
	
	// ========== 适用范围 ==========
	
	// RoomTypeID 适用房型ID，外键关联 room_type_dict 表，可选
	// NULL: 通用政策，可以应用于任何房型
	// 非NULL: 专用政策，仅适用于指定房型
	//
	// 使用场景：
	//   - 高端房型：使用更严格的退订政策
	//   - 特价房型：使用不可退订政策
	//   - 普通房型：使用标准退订政策
	//
	// 用途：政策分类、自动应用
	RoomTypeID *uint64 `gorm:"column:room_type_id;type:BIGINT UNSIGNED;index:idx_room_type_id;comment:适用房型ID（外键，关联房型表）" json:"room_type_id,omitempty"`
	
	// ========== 状态控制 ==========
	
	// Status 状态，控制政策是否可用
	// 可选值：
	//   - "ACTIVE": 启用，可以应用此政策
	//   - "INACTIVE": 停用，不能应用新的（已应用的不受影响）
	//
	// 停用场景：
	//   - 政策调整：旧政策停用，启用新政策
	//   - 活动结束：活动期间的特殊政策
	//   - 季节变化：淡旺季不同政策
	//
	// 用途：政策管理、灵活调整
	Status string `gorm:"column:status;type:VARCHAR(20);NOT NULL;default:'ACTIVE';index:idx_status;comment:状态：ACTIVE-启用，INACTIVE-停用" json:"status"`
	
	// ========== 时间戳 ==========
	
	// CreatedAt 创建时间，政策首次创建的时间
	CreatedAt time.Time `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	
	// UpdatedAt 更新时间，政策最后修改的时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	
	// DeletedAt 软删除时间，非NULL表示已删除
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// ========== 关联关系 ==========
	
	// RoomType 适用的房型信息（可选）
	RoomType *RoomTypeDict `gorm:"foreignKey:RoomTypeID;references:ID" json:"room_type,omitempty"`
	
	// RoomInfos 使用此政策的房源
	// 一对多关系：一个政策可以应用于多个房间
	RoomInfos []RoomInfo `gorm:"foreignKey:CancellationPolicyID;references:ID" json:"room_infos,omitempty"`
	
	// Orders 使用此政策的订单
	// 一对多关系：订单创建时记录使用的退订政策
	Orders []OrderMain `gorm:"foreignKey:CancellationPolicyID;references:ID" json:"orders,omitempty"`
}

// ==================== 表名配置 ====================

// TableName 指定数据库表名
//
// 返回：数据库表名 "cancellation_policy"
func (CancellationPolicy) TableName() string {
	return "cancellation_policy"
}
