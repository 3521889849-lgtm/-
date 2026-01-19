// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了操作日志的数据模型
//
// 功能说明：
//   - 记录系统关键操作日志
//   - 支撑审计追溯和责任认定
//   - 支持操作行为分析
//   - 支持安全监控和异常检测
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 操作日志模型 ====================

// OperationLog 操作日志表
//
// 业务用途：
//   - 审计追溯：记录谁在什么时间做了什么操作
//   - 责任认定：出现问题时追溯操作人
//   - 安全监控：监控异常操作、频繁操作、非法访问
//   - 数据恢复：通过日志恢复误删除或误修改的数据
//   - 行为分析：分析用户操作习惯、系统使用情况
//   - 合规要求：满足监管部门对操作日志的要求
//
// 设计说明：
//   - 记录关键业务操作（房源修改、订单处理、财务操作等）
//   - 不记录频繁的查询操作（避免日志表过大）
//   - 记录操作是否成功，失败时记录原因
//   - 记录操作IP，支持安全分析
//   - 支持软删除，保证日志完整性
//
// 记录原则：
//   - 增删改操作必须记录
//   - 敏感信息查询需要记录
//   - 普通查询操作不记录
//   - 记录内容要包含足够的上下文信息
type OperationLog struct {
	// ========== 基础字段 ==========
	
	// ID 日志ID，主键，自增
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:日志ID" json:"id"`
	
	// ========== 操作人信息 ==========
	
	// OperatorID 操作人ID，外键关联 user_account 表
	// 用途：
	//   - 记录是谁执行的操作
	//   - 责任追溯
	//   - 用户行为分析
	OperatorID uint64 `gorm:"column:operator_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_operator_id;comment:操作人ID（外键，关联账号表）" json:"operator_id"`
	
	// ========== 操作分类 ==========
	
	// Module 操作模块，标识操作发生在哪个业务模块
	// 可选值：
	//   - "ROOM": 房源管理
	//   - "ORDER": 订单处理
	//   - "GUEST": 客人管理
	//   - "MEMBER": 会员管理
	//   - "FINANCIAL": 财务管理
	//   - "USER": 账号管理
	//   - "ROLE": 角色权限管理
	//   - "SYSTEM": 系统配置
	//   - "CHANNEL": 渠道配置
	//   - "REPORT": 报表统计
	// 用途：日志分类、模块使用统计
	Module string `gorm:"column:module;type:VARCHAR(50);NOT NULL;index:idx_module;comment:操作模块（房源管理/订单处理等）" json:"module"`
	
	// OperationType 操作类型，标识具体的操作行为
	// 可选值：
	//   - "QUERY": 查询（一般不记录，除非是敏感查询）
	//   - "CREATE": 添加/创建
	//   - "UPDATE": 修改/更新
	//   - "DELETE": 删除
	//   - "LOGIN": 登录
	//   - "LOGOUT": 登出
	//   - "EXPORT": 导出数据
	//   - "IMPORT": 导入数据
	//   - "APPROVE": 审批
	//   - "REJECT": 拒绝
	// 用途：操作行为分类、安全监控
	OperationType string `gorm:"column:operation_type;type:VARCHAR(20);NOT NULL;index:idx_operation_type;comment:操作类型（查询/添加/修改/删除等）" json:"operation_type"`
	
	// ========== 操作内容 ==========
	
	// Content 操作内容，详细描述这次操作的内容
	// 格式建议：JSON格式或结构化文本
	//
	// 示例：
	//   - 添加房源：{"action":"CREATE","room_no":"101","room_name":"舒适大床房","price":298.00}
	//   - 修改订单：{"action":"UPDATE","order_no":"ORD001","field":"status","old":"RESERVED","new":"CHECKED_IN"}
	//   - 删除客人：{"action":"DELETE","guest_id":123,"guest_name":"张三"}
	//   - 查询敏感信息：{"action":"QUERY","guest_id":123,"fields":["id_number","phone"]}
	//   - 登录系统：{"action":"LOGIN","username":"admin","result":"success"}
	//
	// 注意：
	//   - 不要记录敏感信息的明文（如：密码）
	//   - 对于修改操作，建议记录修改前后的值
	//   - 内容要包含足够的上下文，方便后续追溯
	//
	// 用途：详细记录、数据恢复、问题排查
	Content string `gorm:"column:content;type:TEXT;NOT NULL;comment:操作内容" json:"content"`
	
	// ========== 时间和IP信息 ==========
	
	// OperationTime 操作时间，操作实际发生的时间
	// 用途：
	//   - 时间排序：按时间顺序查询日志
	//   - 统计分析：按时间段统计操作频率
	//   - 异常检测：检测非工作时间的操作
	OperationTime time.Time `gorm:"column:operation_time;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;index:idx_operation_time;comment:操作时间" json:"operation_time"`
	
	// OperationIP 操作IP，操作者的IP地址
	// 格式：IPv4或IPv6地址
	// 示例："192.168.1.100"、"::1"
	// 用途：
	//   - 安全监控：检测异常IP访问
	//   - 地域分析：分析操作来源
	//   - 异常检测：同一账号从不同IP登录
	OperationIP string `gorm:"column:operation_ip;type:VARCHAR(50);NOT NULL;comment:操作IP" json:"operation_ip"`
	
	// ========== 关联信息 ==========
	
	// RelatedID 关联ID，操作对象的ID，可选
	// 用途：快速定位操作对象
	// 示例：
	//   - 修改房源时，记录房源ID
	//   - 处理订单时，记录订单ID
	//   - 删除客人时，记录客人ID
	// 说明：具体含义根据Module确定
	RelatedID *uint64 `gorm:"column:related_id;type:BIGINT UNSIGNED;index:idx_related_id;comment:关联ID（如房源ID/订单ID）" json:"related_id,omitempty"`
	
	// ========== 操作结果 ==========
	
	// IsSuccess 操作是否成功
	// true: 操作成功完成
	// false: 操作失败（如：权限不足、数据验证失败、系统错误）
	// 用途：
	//   - 统计成功率
	//   - 监控失败操作
	//   - 问题排查（失败的操作需要重点关注）
	// 说明：失败时，Content字段应包含失败原因
	IsSuccess bool `gorm:"column:is_success;type:BOOLEAN;NOT NULL;default:true;index:idx_is_success;comment:是否成功（0/1）" json:"is_success"`
	
	// ========== 时间戳 ==========
	
	// CreatedAt 创建时间，日志记录创建的时间
	// 说明：通常与OperationTime相同
	CreatedAt time.Time `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	
	// DeletedAt 软删除时间，非NULL表示已删除
	// 说明：操作日志原则上不应删除，仅在特殊情况下（如：测试数据清理）才软删除
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// ========== 关联关系 ==========
	
	// Operator 操作人信息
	// 用途：获取操作人的姓名、角色等信息
	Operator *UserAccount `gorm:"foreignKey:OperatorID;references:ID" json:"operator,omitempty"`
}

// ==================== 表名配置 ====================

// TableName 指定数据库表名
//
// 返回：数据库表名 "operation_log"
func (OperationLog) TableName() string {
	return "operation_log"
}
