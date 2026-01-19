// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了会员权益的数据模型
//
// 功能说明：
//   - 维护会员的专属权益配置
//   - 支持会员权益管理功能
//   - 根据会员等级提供不同权益
//   - 支持权益的灵活配置和调整
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 会员权益模型 ====================

// MemberRights 会员权益表
//
// 业务用途：
//   - 权益配置：为不同等级的会员配置专属权益
//   - 权益展示：在会员中心展示可享受的权益
//   - 权益判断：订房时自动应用会员权益
//   - 营销工具：通过权益吸引客人成为会员
//   - 等级激励：高等级会员享受更多权益，激励消费
//
// 设计说明：
//   - 通过MemberLevel字段关联到会员
//   - 同一等级可以有多个权益
//   - 支持权益的生效和失效时间
//   - 支持权益的启用/停用
//   - 折扣比例字段支持灵活配置优惠力度
type MemberRights struct {
	// ========== 基础字段 ==========
	
	// ID 权益ID，主键，自增
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:权益ID" json:"id"`
	
	// ========== 权益归属 ==========
	
	// MemberLevel 会员等级，标识该权益属于哪个等级的会员
	// 可选值：
	//   - "NORMAL": 普通会员
	//   - "SILVER": 银卡会员
	//   - "GOLD": 金卡会员
	//   - "PLATINUM": 白金会员
	//   - "DIAMOND": 钻石会员
	//
	// 说明：
	//   - 同一等级可以配置多个权益
	//   - 高等级会员可以享受低等级的所有权益
	//   - 查询会员权益时，需要查询该等级的所有权益记录
	//
	// 用途：权益分类、权益查询、权益展示
	MemberLevel string `gorm:"column:member_level;type:VARCHAR(30);NOT NULL;index:idx_member_level;comment:会员等级" json:"member_level"`
	
	// ========== 权益信息 ==========
	
	// RightsName 权益名称，权益的标准名称
	// 权益类型示例：
	//
	// 预订权益：
	//   - "专属预订通道"：会员专属预订渠道
	//   - "优先预订权"：紧俏房型优先预订
	//   - "预订优惠"：预订时享受折扣
	//   - "免预订费"：免收预订手续费
	//
	// 房价折扣：
	//   - "95折优惠"：房费享95折
	//   - "9折优惠"：房费享9折
	//   - "85折优惠"：房费享85折
	//   - "8折优惠"：房费享8折
	//
	// 积分权益：
	//   - "积分兑换"：积分可兑换房费或礼品
	//   - "双倍积分"：消费获得双倍积分
	//   - "积分抵扣"：积分可抵扣部分房费
	//   - "积分赠送"：每月赠送固定积分
	//
	// 升级权益：
	//   - "免费升房"：免费升级更高房型
	//   - "优先升房"：房型充足时优先升级
	//   - "套房升级"：升级至套房
	//
	// 服务权益：
	//   - "免费早餐"：享受免费早餐
	//   - "延迟退房"：可延迟至下午2点退房
	//   - "专属客服"：一对一专属客服
	//   - "贵宾通道"：快速办理入住退房
	//   - "免费停车"：享受免费停车位
	//   - "接送机服务"：免费接送机服务
	//
	// 特殊权益：
	//   - "生日特权"：生日当月享受特殊优惠
	//   - "免费取消"：任何时候免费取消预订
	//   - "积分永久有效"：积分不会过期
	//   - "会籍延期"：会籍有效期自动延长
	//
	// 用途：权益展示、权益判断、营销宣传
	RightsName string `gorm:"column:rights_name;type:VARCHAR(100);NOT NULL;comment:权益名称（专属预订/积分兑换/房价折扣等）" json:"rights_name"`
	
	// Description 权益描述，权益的详细说明，可选
	// 内容要求：
	//   - 详细说明权益的具体内容
	//   - 说明权益的使用方法
	//   - 标注权益的使用限制
	//   - 特殊情况的处理说明
	//
	// 示例：
	//   - "享受房费95折优惠，特价房、活动房除外"
	//   - "每消费1元获得1积分，积分可用于兑换礼品或抵扣房费"
	//   - "入住时如有同等级或更低级房型可用，免费升级至更高房型"
	//   - "每天7-10点提供免费自助早餐，需提前预约"
	//   - "标准退房时间12点，会员可延迟至下午2点退房，不额外收费"
	//   - "提供24小时一对一专属客服，随时解答疑问和处理问题"
	//
	// 用途：客人知情、权益解释、纠纷处理
	Description *string `gorm:"column:description;type:VARCHAR(500);comment:权益描述" json:"description,omitempty"`
	
	// ========== 折扣配置 ==========
	
	// DiscountRatio 折扣比例，享受的折扣优惠，可选
	// 单位：折扣系数
	// 取值范围：
	//   - 1.0: 不打折（原价）
	//   - 0.95: 95折（5%优惠）
	//   - 0.9: 9折（10%优惠）
	//   - 0.85: 85折（15%优惠）
	//   - 0.8: 8折（20%优惠）
	//   - 0.0: 免费（特殊权益，如：生日免费住）
	//
	// 计算公式：
	//   实付金额 = 原价 × DiscountRatio
	//
	// 示例：
	//   - 原价298元，DiscountRatio=0.95，实付283.1元
	//   - 原价298元，DiscountRatio=0.9，实付268.2元
	//
	// NULL表示该权益不涉及折扣（如：免费早餐、延迟退房）
	//
	// 用途：自动计算折扣、价格展示
	DiscountRatio *float64 `gorm:"column:discount_ratio;type:DECIMAL(5,2);comment:折扣比例" json:"discount_ratio,omitempty"`
	
	// ========== 时间控制 ==========
	
	// EffectiveTime 生效时间，权益开始生效的时间
	// 用途：
	//   - 新权益上线：设定未来时间，到期自动生效
	//   - 活动权益：活动开始时间
	//   - 季节权益：淡旺季不同权益
	//   - 权益升级：新老权益的切换时间
	//
	// 检查逻辑：
	//   当前时间 >= EffectiveTime 且 当前时间 < ExpireTime（或ExpireTime为NULL）
	//   且 Status = 'ACTIVE' 时，权益生效
	EffectiveTime time.Time `gorm:"column:effective_time;type:DATETIME;NOT NULL;index:idx_effective_time;comment:生效时间" json:"effective_time"`
	
	// ExpireTime 失效时间，权益结束的时间，可选
	// NULL: 永久有效（除非手动停用）
	// 非NULL: 到期自动失效
	//
	// 用途：
	//   - 活动权益：活动结束时间
	//   - 临时权益：短期权益的有效期
	//   - 季节权益：淡旺季切换时间
	//
	// 检查逻辑：
	//   如果ExpireTime不为NULL且当前时间 >= ExpireTime，权益失效
	ExpireTime *time.Time `gorm:"column:expire_time;type:DATETIME;index:idx_expire_time;comment:失效时间" json:"expire_time,omitempty"`
	
	// ========== 状态控制 ==========
	
	// Status 状态，控制权益是否可用
	// 可选值：
	//   - "ACTIVE": 启用，权益生效中（需同时满足时间条件）
	//   - "INACTIVE": 停用，权益暂停使用
	//
	// 停用场景：
	//   - 权益调整：临时停用旧权益，测试新权益
	//   - 紧急情况：发现权益配置错误，紧急停用
	//   - 成本控制：控制成本时临时停用部分权益
	//
	// 用途：权益管理、灵活调整
	Status string `gorm:"column:status;type:VARCHAR(20);NOT NULL;default:'ACTIVE';index:idx_status;comment:状态（启用/停用）" json:"status"`
	
	// ========== 时间戳 ==========
	
	// CreatedAt 创建时间，权益首次创建的时间
	CreatedAt time.Time `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	
	// UpdatedAt 更新时间，权益最后修改的时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	
	// DeletedAt 软删除时间，非NULL表示已删除
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`
}

// ==================== 表名配置 ====================

// TableName 指定数据库表名
//
// 返回：数据库表名 "member_rights"
func (MemberRights) TableName() string {
	return "member_rights"
}
