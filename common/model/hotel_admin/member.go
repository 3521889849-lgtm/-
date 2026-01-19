// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了会员信息的数据模型
//
// 功能说明：
//   - 存储会员的基础信息和等级
//   - 支持会员等级体系管理
//   - 支持积分余额管理
//   - 记录会员入住历史
//   - 支持会员状态控制（启用/冻结）
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 会员模型 ====================

// Member 会员表
//
// 业务用途：
//   - 会员等级管理：支持普通、黄金、钻石等多级会员体系
//   - 积分管理：记录会员的积分余额，支持积分获取和使用
//   - 会员权益：根据会员等级提供不同的权益（折扣、专属服务等）
//   - 客户关系：维护长期客户关系，提升客户忠诚度
//   - 营销分析：分析会员消费行为，支持精准营销
//
// 设计说明：
//   - 每个会员对应一个客人（一对一关系）
//   - 客人升级为会员时创建此记录
//   - 会员等级通过 MemberLevel 字段标识
//   - 积分余额实时更新，历史明细记录在积分流水表
//   - 支持会员冻结功能（如：违规行为、欠款等）
type Member struct {
	// ========== 基础字段 ==========
	
	// ID 会员ID，主键，自增
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:会员ID" json:"id"`
	
	// ========== 客人关联 ==========
	
	// GuestID 客人ID，外键关联 guest_info 表，全局唯一
	// 约束：一个客人只能有一个会员账号
	// 用途：关联客人的基本信息（姓名、手机号、身份证等）
	GuestID uint64 `gorm:"column:guest_id;type:BIGINT UNSIGNED;NOT NULL;uniqueIndex:uk_guest_id;comment:客人ID（外键，关联客人表）" json:"guest_id"`
	
	// ========== 会员等级 ==========
	
	// MemberLevel 会员等级，标识会员的级别
	// 可选值：
	//   - "NORMAL": 普通会员（默认等级，新注册会员）
	//   - "SILVER": 银卡会员（消费达到一定金额或次数）
	//   - "GOLD": 金卡会员（高级会员，享受更多权益）
	//   - "PLATINUM": 白金会员
	//   - "DIAMOND": 钻石会员（顶级会员，VIP待遇）
	//
	// 等级权益示例：
	//   - NORMAL: 基础积分、生日优惠
	//   - SILVER: 95折优惠、优先预订
	//   - GOLD: 9折优惠、免费升房、专属客服
	//   - PLATINUM: 85折优惠、免费早餐、延迟退房
	//   - DIAMOND: 8折优惠、套房升级、接送机服务
	//
	// 用途：
	//   - 权益判断：根据等级提供不同折扣和服务
	//   - 积分规则：不同等级的积分获取比例不同
	//   - 营销分析：高等级会员的价值分析
	MemberLevel string `gorm:"column:member_level;type:VARCHAR(30);NOT NULL;default:'NORMAL';index:idx_member_level;comment:会员等级（普通/黄金/钻石等）" json:"member_level"`
	
	// ========== 积分管理 ==========
	
	// PointsBalance 积分余额，当前可用积分数
	// 单位：积分（通常1元=1积分，可配置）
	// 用途：
	//   - 积分抵扣：积分可用于抵扣房费（如：100积分=10元）
	//   - 礼品兑换：积分可兑换礼品或服务
	//   - 会员升级：积分达到一定数量可升级会员等级
	//
	// 积分来源：
	//   - 消费获取：入住消费产生积分
	//   - 活动赠送：营销活动赠送积分
	//   - 推荐奖励：推荐新会员获得积分
	//
	// 积分消耗：
	//   - 积分抵扣：使用积分抵扣房费
	//   - 礼品兑换：使用积分兑换礼品
	//   - 积分过期：长期不活跃的积分可能过期
	//
	// 说明：积分明细记录在 member_points_record 表中
	PointsBalance uint64 `gorm:"column:points_balance;type:BIGINT UNSIGNED;NOT NULL;default:0;comment:积分余额" json:"points_balance"`
	
	// ========== 时间信息 ==========
	
	// RegisterTime 注册时间，客人成为会员的时间
	// 用途：
	//   - 会员年限计算
	//   - 会员生日优惠（注册周年庆）
	//   - 统计分析：会员增长趋势
	RegisterTime time.Time `gorm:"column:register_time;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:注册时间" json:"register_time"`
	
	// LastCheckInTime 最后入住时间，会员最近一次入住的时间，可选
	// NULL表示注册后从未入住
	// 非NULL表示有入住记录
	// 用途：
	//   - 活跃度判断：长期未入住的会员视为不活跃
	//   - 营销触达：对不活跃会员发送唤醒营销
	//   - 积分政策：不活跃会员的积分可能过期
	LastCheckInTime *time.Time `gorm:"column:last_check_in_time;type:DATETIME;index:idx_last_check_in_time;comment:最后入住时间" json:"last_check_in_time,omitempty"`
	
	// ========== 状态控制 ==========
	
	// Status 会员状态，控制会员是否可用
	// 可选值：
	//   - "ACTIVE": 启用，正常使用（默认状态）
	//   - "FROZEN": 冻结，暂停使用
	//
	// 冻结场景：
	//   - 违规行为：如恶意差评、欺诈行为
	//   - 欠款未还：长期欠款不还的会员
	//   - 会员申请：会员主动申请暂停
	//
	// 冻结影响：
	//   - 不能享受会员权益
	//   - 不能使用积分
	//   - 不能获取新积分
	//   - 可以正常入住（按散客处理）
	Status string `gorm:"column:status;type:VARCHAR(20);NOT NULL;default:'ACTIVE';index:idx_status;comment:会员状态（启用/冻结）" json:"status"`
	
	// ========== 时间戳 ==========
	
	// CreatedAt 创建时间，会员记录首次创建的时间
	CreatedAt time.Time `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	
	// UpdatedAt 更新时间，会员信息最后修改的时间
	// 自动更新：每次修改会员信息时自动更新此字段
	UpdatedAt time.Time `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	
	// DeletedAt 软删除时间，非NULL表示已删除
	// 软删除：会员注销时不真正删除数据，只是标记
	// 好处：数据可恢复，历史记录完整性保证
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// ========== 关联关系 ==========
	// 以下字段不会存储在数据库中，仅用于GORM的关联查询
	
	// Guest 关联的客人信息
	// 一对一关系：一个会员对应一个客人
	// 用途：获取会员的姓名、手机号、身份证等基本信息
	Guest *GuestInfo `gorm:"foreignKey:GuestID;references:ID" json:"guest,omitempty"`
	
	// PointsRecords 积分流水记录
	// 一对多关系：一个会员可以有多条积分流水
	// 用途：查询积分的获取和使用明细
	PointsRecords []MemberPointsRecord `gorm:"foreignKey:MemberID;references:ID" json:"points_records,omitempty"`
	
	// 注意：
	// MemberRights 是会员权益配置表，通过 member_level 字段逻辑关联
	// 不建立外键约束，因为 member_level 不是唯一键（多个会员可以是同一等级）
	// 查询会员权益时，需要通过 WHERE member_rights.level = member.member_level 进行关联
}

// ==================== 表名配置 ====================

// TableName 指定数据库表名
//
// GORM会自动调用此方法获取表名，用于生成SQL语句
//
// 返回：
//   - string: 数据库表名 "member"
func (Member) TableName() string {
	return "member"
}
