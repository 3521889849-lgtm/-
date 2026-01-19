// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了系统配置的数据模型
//
// 功能说明：
//   - 存储系统的基础配置项
//   - 支撑短信模板、打印设置等功能
//   - 实现配置的集中管理
//   - 支持配置的动态更新（无需重启系统）
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 系统配置模型 ====================

// SystemConfig 系统配置表
//
// 业务用途：
//   - 集中管理：所有系统配置集中存储在一张表中
//   - 动态配置：修改配置后立即生效，无需重启系统
//   - 灵活扩展：新增配置项只需添加记录，无需修改代码
//   - 分类管理：通过ConfigCategory对配置进行分类
//   - 版本控制：记录修改人和修改时间，方便追溯
//
// 设计说明：
//   - ConfigKey全局唯一，作为配置的唯一标识
//   - ConfigValue使用TEXT类型，支持存储大量数据
//   - 支持JSON格式的复杂配置
//   - 通过Status控制配置是否生效
//   - 记录修改人，方便审计追溯
type SystemConfig struct {
	// ========== 基础字段 ==========
	
	// ID 配置ID，主键，自增
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:配置ID" json:"id"`
	
	// ========== 配置分类 ==========
	
	// ConfigCategory 配置分类，将配置按功能分类
	// 配置分类示例：
	//
	// 消费项设置：
	//   - "CONSUMPTION_ITEMS": 消费项目配置
	//   - 示例配置：早餐费、加床费、洗衣费、停车费等
	//
	// 短信模板：
	//   - "SMS_TEMPLATE": 短信模板配置
	//   - 示例配置：预订成功通知、入住提醒、退房通知、营销短信等
	//
	// 打印设置：
	//   - "PRINT_SETTINGS": 打印相关配置
	//   - 示例配置：小票格式、发票格式、标签格式、纸张尺寸等
	//
	// 会员规则：
	//   - "MEMBER_RULES": 会员管理规则
	//   - 示例配置：积分规则、升级规则、权益规则、有效期规则等
	//
	// 支付配置：
	//   - "PAYMENT_CONFIG": 支付相关配置
	//   - 示例配置：支付宝配置、微信配置、银联配置、手续费率等
	//
	// 通知配置：
	//   - "NOTIFICATION_CONFIG": 通知相关配置
	//   - 示例配置：邮件通知、短信通知、系统消息等
	//
	// 业务规则：
	//   - "BUSINESS_RULES": 业务流程规则
	//   - 示例配置：入住时间、退房时间、超时计费规则、押金标准等
	//
	// 系统参数：
	//   - "SYSTEM_PARAMS": 系统运行参数
	//   - 示例配置：分页大小、缓存时间、文件上传限制等
	//
	// 用途：配置分类、查询过滤、权限控制
	ConfigCategory string `gorm:"column:config_category;type:VARCHAR(50);NOT NULL;index:idx_config_category;comment:配置项（消费项设置/短信模板/打印设置等）" json:"config_category"`
	
	// ========== 配置键值 ==========
	
	// ConfigKey 配置键，配置项的唯一标识，全局唯一
	// 命名规范：
	//   - 使用下划线分隔单词
	//   - 全大写或小驼峰
	//   - 见名知意
	//
	// 配置键示例：
	//   - "BREAKFAST_FEE": 早餐费
	//   - "EXTRA_BED_FEE": 加床费
	//   - "SMS_BOOKING_SUCCESS": 预订成功短信模板
	//   - "SMS_CHECK_IN_REMIND": 入住提醒短信模板
	//   - "PRINT_RECEIPT_FORMAT": 小票打印格式
	//   - "MEMBER_POINTS_RATIO": 会员积分比例
	//   - "DEFAULT_CHECK_IN_TIME": 默认入住时间
	//   - "DEFAULT_CHECK_OUT_TIME": 默认退房时间
	//   - "ALIPAY_APP_ID": 支付宝APPID
	//   - "WECHAT_MCH_ID": 微信商户号
	//
	// 用途：代码中通过ConfigKey查询配置值
	ConfigKey string `gorm:"column:config_key;type:VARCHAR(100);NOT NULL;uniqueIndex:uk_config_key;comment:配置键" json:"config_key"`
	
	// ConfigValue 配置值，配置项的具体值
	// 数据类型：TEXT，支持大量数据
	// 数据格式：
	//   - 简单值：直接存储字符串、数字等
	//   - JSON格式：存储复杂的结构化数据
	//   - 模板内容：存储短信模板、打印模板等
	//
	// 配置值示例：
	//   - 简单值：
	//     ConfigKey="BREAKFAST_FEE", ConfigValue="15"（表示早餐费15元）
	//     ConfigKey="DEFAULT_CHECK_IN_TIME", ConfigValue="14:00"
	//
	//   - JSON格式：
	//     ConfigKey="CONSUMPTION_ITEMS", ConfigValue='[
	//       {"name":"早餐","price":15},
	//       {"name":"加床","price":50},
	//       {"name":"洗衣","price":20}
	//     ]'
	//
	//   - 模板内容：
	//     ConfigKey="SMS_BOOKING_SUCCESS", ConfigValue="尊敬的{name}，您已成功预订{hotelName}{roomType}，入住时间{checkInTime}，订单号{orderNo}"
	//     ConfigKey="PRINT_RECEIPT_FORMAT", ConfigValue="<html><body>...</body></html>"
	//
	// 使用说明：
	//   - 代码中读取ConfigValue后，根据需要进行类型转换
	//   - JSON格式需要反序列化后使用
	//   - 模板内容需要替换占位符后使用
	//
	// 用途：存储配置的实际值
	ConfigValue string `gorm:"column:config_value;type:TEXT;NOT NULL;comment:配置值" json:"config_value"`
	
	// ========== 配置说明 ==========
	
	// Description 配置描述，配置项的详细说明，可选
	// 内容要求：
	//   - 说明配置的用途
	//   - 说明配置值的格式和要求
	//   - 说明修改配置的注意事项
	//   - 提供配置示例
	//
	// 示例：
	//   - "早餐费标准，单位：元。前台收费时使用此价格"
	//   - "预订成功短信模板，支持变量：{name}客人姓名、{hotelName}酒店名称、{roomType}房型、{checkInTime}入住时间、{orderNo}订单号"
	//   - "小票打印格式，HTML格式，打印时自动填充数据"
	//   - "会员消费积分比例，1元=X积分。示例：1.0表示1元获得1积分，2.0表示1元获得2积分"
	//   - "支付宝APPID，修改后需要重新配置支付参数"
	//
	// 用途：配置说明、操作指引、问题排查
	Description *string `gorm:"column:description;type:VARCHAR(500);comment:配置描述" json:"description,omitempty"`
	
	// ========== 状态控制 ==========
	
	// Status 生效状态，控制配置是否启用
	// 可选值：
	//   - "ACTIVE": 启用，配置生效中
	//   - "INACTIVE": 停用，配置暂时不生效
	//
	// 停用场景：
	//   - 测试配置：新配置上线前先停用，测试通过后启用
	//   - 临时调整：临时停用某些配置
	//   - 功能开关：通过配置控制功能的开启关闭
	//
	// 用途：
	//   - 代码中只读取Status='ACTIVE'的配置
	//   - 方便配置的灵活切换
	Status string `gorm:"column:status;type:VARCHAR(20);NOT NULL;default:'ACTIVE';index:idx_status;comment:生效状态" json:"status"`
	
	// ========== 时间戳和操作人 ==========
	
	// UpdatedAt 修改时间，配置最后修改的时间
	// 自动更新：每次修改配置时自动更新此字段
	// 用途：追踪配置变更历史
	UpdatedAt time.Time `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:修改时间" json:"updated_at"`
	
	// UpdatedBy 修改人ID，关联 user_account 表
	// 用途：
	//   - 责任追溯：记录是谁修改的配置
	//   - 问题排查：配置异常时联系修改人
	//   - 操作审计：监控配置修改行为
	UpdatedBy uint64 `gorm:"column:updated_by;type:BIGINT UNSIGNED;NOT NULL;comment:修改人" json:"updated_by"`
	
	// CreatedAt 创建时间，配置首次创建的时间
	CreatedAt time.Time `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	
	// DeletedAt 软删除时间，非NULL表示已删除
	// 说明：配置原则上不删除，通过Status控制是否生效
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`
}

// ==================== 表名配置 ====================

// TableName 指定数据库表名
//
// 返回：数据库表名 "system_config"
func (SystemConfig) TableName() string {
	return "system_config"
}
