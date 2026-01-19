// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了黑名单的数据模型
//
// 功能说明：
//   - 维护黑名单用户信息
//   - 支撑黑名单管理功能
//   - 防止问题客人再次入住
//   - 保护酒店权益和其他客人安全
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 黑名单模型 ====================

// Blacklist 黑名单表
//
// 业务用途：
//   - 风险控制：记录有不良记录的客人
//   - 入住拦截：前台接待时自动检查黑名单
//   - 安全保障：防止问题客人再次入住
//   - 证据保存：记录拉黑原因和操作人
//   - 申诉处理：黑名单可以解除（状态改为无效）
//
// 设计说明：
//   - 通过身份证号和手机号双重识别
//   - 敏感信息加密存储
//   - 支持黑名单的启用/解除（通过Status字段）
//   - 记录拉黑原因和操作人，方便追溯
//   - 可以关联到客人记录，也可以独立存在
type Blacklist struct {
	// ========== 基础字段 ==========
	
	// ID 黑名单ID，主键，自增
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:黑名单ID" json:"id"`
	
	// ========== 客人关联 ==========
	
	// GuestID 客人ID，外键关联 guest_info 表，可选
	// NULL: 客人信息尚未录入系统（可能是预防性拉黑）
	// 非NULL: 关联到已有客人记录
	//
	// 说明：
	//   - 客人首次入住时被拉黑：GuestID为NULL，只有身份证和手机号
	//   - 客人入住后被拉黑：GuestID非NULL，关联到客人记录
	//
	// 用途：关联客人的历史入住记录
	GuestID *uint64 `gorm:"column:guest_id;type:BIGINT UNSIGNED;index:idx_guest_id;comment:客人ID（外键，关联客人表）" json:"guest_id,omitempty"`
	
	// ========== 识别信息（敏感信息，加密存储）==========
	
	// IDNumber 证件号码，必填，加密存储
	// ⚠️ 安全要求：
	//   - 存储前必须使用 AES 加密
	//   - 查询时需解密后使用
	//   - 展示时必须脱敏
	//   - 访问需记录操作日志
	//
	// 用途：
	//   - 唯一身份标识
	//   - 入住时自动检查黑名单
	//   - 防止客人使用不同手机号重新入住
	//
	// 检查逻辑：
	//   客人办理入住时，系统自动查询黑名单表
	//   如果身份证号匹配且状态为VALID，则拒绝入住
	IDNumber string `gorm:"column:id_number;type:VARCHAR(50);NOT NULL;index:idx_id_number;comment:证件号（加密存储）" json:"id_number"`
	
	// Phone 手机号码，必填，加密存储
	// ⚠️ 安全要求：同IDNumber
	//
	// 用途：
	//   - 辅助识别客人
	//   - 防止客人更换身份证后重新入住
	//   - 电话订房时检查黑名单
	//
	// 检查逻辑：
	//   电话订房时，系统自动查询黑名单表
	//   如果手机号匹配且状态为VALID，则拒绝预订
	Phone string `gorm:"column:phone;type:VARCHAR(20);NOT NULL;index:idx_phone;comment:手机号（加密存储）" json:"phone"`
	
	// ========== 拉黑信息 ==========
	
	// Reason 拉黑原因，详细说明拉黑的具体原因
	// 内容要求：
	//   - 描述具体的违规行为
	//   - 列举相关证据
	//   - 记录处理过程
	//   - 说明拉黑依据
	//
	// 常见拉黑原因示例：
	//   - "恶意损坏房间设施，造成经济损失约500元，拒绝赔偿"
	//   - "多次无故取消预订，影响酒店正常经营"
	//   - "入住期间扰乱酒店秩序，影响其他客人休息"
	//   - "使用虚假证件入住，涉嫌违法"
	//   - "欠款不还，累计欠款金额1000元"
	//   - "盗窃酒店财物，已报警处理"
	//   - "恶意差评敲诈，要求不合理赔偿"
	//   - "在房间内从事违法活动，已移交公安机关"
	//
	// 用途：
	//   - 拉黑依据
	//   - 申诉处理参考
	//   - 纠纷证据
	//   - 其他分店参考
	Reason string `gorm:"column:reason;type:VARCHAR(500);NOT NULL;comment:拉黑原因" json:"reason"`
	
	// BlackTime 拉黑时间，客人被加入黑名单的时间
	// 用途：
	//   - 时间排序
	//   - 统计分析
	//   - 申诉处理（如：一年后自动解除）
	BlackTime time.Time `gorm:"column:black_time;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:拉黑时间" json:"black_time"`
	
	// OperatorID 拉黑操作人ID，关联 user_account 表
	// 用途：
	//   - 责任追溯：记录是谁拉黑的
	//   - 申诉处理：联系操作人了解情况
	//   - 操作审计：监控拉黑操作是否合理
	OperatorID uint64 `gorm:"column:operator_id;type:BIGINT UNSIGNED;NOT NULL;comment:拉黑操作人" json:"operator_id"`
	
	// ========== 状态控制 ==========
	
	// Status 状态，控制黑名单是否生效
	// 可选值：
	//   - "VALID": 有效，黑名单生效中（默认状态）
	//   - "INVALID": 无效，黑名单已解除
	//
	// 解除场景：
	//   - 客人申诉成功：经调查确认拉黑不当
	//   - 赔偿完成：客人补偿了损失
	//   - 时间到期：设定期限（如：一年）后自动解除
	//   - 特殊情况：酒店主动解除
	//
	// 用途：
	//   - 入住检查：只检查状态为VALID的记录
	//   - 申诉处理：申诉通过后改为INVALID
	//   - 数据保留：解除后不删除记录，只改状态
	Status string `gorm:"column:status;type:VARCHAR(20);NOT NULL;default:'VALID';index:idx_status;comment:状态（有效/无效）" json:"status"`
	
	// ========== 时间戳 ==========
	
	// CreatedAt 创建时间，黑名单记录创建的时间
	// 说明：通常与BlackTime相同
	CreatedAt time.Time `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	
	// UpdatedAt 更新时间，黑名单记录最后修改的时间
	// 用途：记录状态变更时间（如：解除黑名单的时间）
	UpdatedAt time.Time `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	
	// DeletedAt 软删除时间，非NULL表示已删除
	// 说明：黑名单记录原则上不删除，通过Status控制是否生效
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// ========== 关联关系 ==========
	
	// Guest 关联的客人信息（可选）
	// 用途：查看客人的历史入住记录和详细信息
	Guest *GuestInfo `gorm:"foreignKey:GuestID;references:ID" json:"guest,omitempty"`
}

// ==================== 表名配置 ====================

// TableName 指定数据库表名
//
// 返回：数据库表名 "blacklist"
func (Blacklist) TableName() string {
	return "blacklist"
}
