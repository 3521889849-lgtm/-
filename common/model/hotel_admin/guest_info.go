// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了客人信息的数据模型
//
// 功能说明：
//   - 存储客人的实名制登记信息
//   - 支撑公安系统合规备案要求
//   - 支持在住客人查询和历史入住记录
//   - 敏感信息（身份证号、手机号）加密存储
//   - 支持会员关联和积分管理
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 客人信息模型 ====================

// GuestInfo 客人信息表
//
// 业务用途：
//   - 实名制登记：满足公安部门的住宿登记要求
//   - 客人管理：记录所有入住过的客人信息
//   - 快速入住：老客人再次入住时可快速调取信息
//   - 会员关联：普通客人升级为会员时建立关联
//   - 历史追溯：查询客人的入住历史记录
//
// 设计说明：
//   - 身份证号和手机号必须加密存储，保护客人隐私
//   - 支持多种证件类型（身份证、护照、港澳通行证等）
//   - 与订单、会员、财务模块关联
//   - 支持软删除，保证数据完整性
//
// 安全说明：
//   - IDNumber 和 Phone 字段使用 AES 加密存储
//   - 前端展示时需要脱敏处理（如：138****1234）
//   - 访问敏感信息需要记录操作日志
type GuestInfo struct {
	// ========== 基础字段 ==========
	
	// ID 客人ID，主键，自增
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:客人ID" json:"id"`
	
	// ========== 基本信息 ==========
	
	// Name 客人姓名，必填
	// 用途：登记、显示、查询
	// 示例："张三"、"John Smith"
	Name string `gorm:"column:name;type:VARCHAR(50);NOT NULL;index:idx_name;comment:姓名" json:"name"`
	
	// Gender 性别，可选
	// 可选值："男"、"女"、"MALE"、"FEMALE"
	// 用途：统计分析、客人偏好记录
	Gender *string `gorm:"column:gender;type:VARCHAR(10);comment:性别" json:"gender,omitempty"`
	
	// Ethnicity 民族，可选
	// 用途：统计分析、特殊服务需求（如：清真餐）
	// 示例："汉族"、"回族"、"维吾尔族"
	Ethnicity *string `gorm:"column:ethnicity;type:VARCHAR(30);comment:民族" json:"ethnicity,omitempty"`
	
	// Province 省份，客人身份证所属省份，可选
	// 用途：客源分析、营销区域划分
	// 示例："北京市"、"广东省"
	Province *string `gorm:"column:province;type:VARCHAR(30);comment:省份" json:"province,omitempty"`
	
	// Address 详细地址，客人的居住地址，可选
	// 用途：公安登记、紧急联系
	// 示例："北京市朝阳区建国路88号"
	Address *string `gorm:"column:address;type:VARCHAR(255);comment:地址" json:"address,omitempty"`
	
	// ========== 证件信息（重要：涉及实名制合规）==========
	
	// IDType 证件类型，必填
	// 可选值：
	//   - "ID_CARD": 身份证（最常用）
	//   - "PASSPORT": 护照（境外客人）
	//   - "HK_MACAU_PASS": 港澳通行证
	//   - "TAIWAN_PASS": 台湾通行证
	//   - "MILITARY_ID": 军官证
	//   - "OTHER": 其他证件
	// 用途：公安系统备案、身份验证
	IDType string `gorm:"column:id_type;type:VARCHAR(20);NOT NULL;comment:证件类型（身份证等）" json:"id_type"`
	
	// IDNumber 证件号码，必填，加密存储
	// ⚠️ 安全要求：
	//   - 存储前必须使用 AES 加密
	//   - 查询时需解密后使用
	//   - 展示时必须脱敏（如：110108********1234）
	//   - 访问需记录操作日志
	// 用途：唯一身份标识、公安系统备案、防止重复登记
	// 索引：支持快速查询（基于加密后的值）
	IDNumber string `gorm:"column:id_number;type:VARCHAR(50);NOT NULL;index:idx_id_number;comment:证件号（加密存储）" json:"id_number"`
	
	// ========== 联系方式（敏感信息）==========
	
	// Phone 手机号，必填，加密存储
	// ⚠️ 安全要求：
	//   - 存储前必须使用 AES 加密
	//   - 查询时需解密后使用
	//   - 展示时必须脱敏（如：138****1234）
	//   - 访问需记录操作日志
	// 用途：联系客人、短信通知、会员识别
	// 索引：支持快速查询（基于加密后的值）
	Phone string `gorm:"column:phone;type:VARCHAR(20);NOT NULL;index:idx_phone;comment:手机号（加密存储）" json:"phone"`
	
	// ========== 入住信息 ==========
	
	// CheckInTime 入住时间，可选
	// NULL表示预订但未入住
	// 非NULL表示已办理入住
	// 用途：判断客人状态、统计在店客人
	CheckInTime *time.Time `gorm:"column:check_in_time;type:DATETIME;index:idx_check_in_time;comment:入住时间" json:"check_in_time,omitempty"`
	
	// CheckOutTime 离店时间，可选
	// NULL表示客人仍在店（未退房）
	// 非NULL表示已办理退房
	// 用途：判断客人状态、房态管理
	CheckOutTime *time.Time `gorm:"column:check_out_time;type:DATETIME;index:idx_check_out_time;comment:离店时间" json:"check_out_time,omitempty"`
	
	// RoomID 当前所住房间ID，外键关联 room_info 表，可选
	// NULL表示客人未入住
	// 非NULL表示客人正在此房间
	// 用途：房态图显示、在住客人查询
	RoomID *uint64 `gorm:"column:room_id;type:BIGINT UNSIGNED;index:idx_room_id;comment:房源ID（外键，关联房源表）" json:"room_id,omitempty"`
	
	// OrderID 关联的订单ID，外键关联 hotel_order_main 表，可选
	// 用途：查询客人的订单信息
	OrderID *uint64 `gorm:"column:order_id;type:BIGINT UNSIGNED;index:idx_order_id;comment:订单ID（外键，关联订单主表）" json:"order_id,omitempty"`
	
	// ========== 登记信息 ==========
	
	// RegisterTime 登记时间，客人信息录入系统的时间
	// 用途：公安系统备案、记录追溯
	RegisterTime time.Time `gorm:"column:register_time;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:登记时间" json:"register_time"`
	
	// RegisterBy 登记人ID，关联用户账号表
	// 用途：审计追溯，记录是谁办理的登记
	RegisterBy uint64 `gorm:"column:register_by;type:BIGINT UNSIGNED;NOT NULL;comment:登记人" json:"register_by"`
	
	// ========== 会员关联 ==========
	
	// IsMember 是否会员标识
	// true: 该客人是会员
	// false: 该客人是普通散客
	// 用途：快速筛选会员、权益判断
	IsMember bool `gorm:"column:is_member;type:BOOLEAN;NOT NULL;default:false;comment:是否会员（0/1）" json:"is_member"`
	
	// MemberID 会员ID，外键关联 member 表，可选
	// NULL表示普通散客
	// 非NULL表示会员，可享受会员权益
	// 用途：会员识别、积分管理、权益兑换
	MemberID *uint64 `gorm:"column:member_id;type:BIGINT UNSIGNED;index:idx_member_id;comment:会员ID（外键，关联会员表，非会员可为空）" json:"member_id,omitempty"`
	
	// ========== 时间戳 ==========
	
	// CreatedAt 创建时间，客人信息首次录入的时间
	CreatedAt time.Time `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	
	// UpdatedAt 更新时间，客人信息最后修改的时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	
	// DeletedAt 软删除时间，非NULL表示已删除
	// 软删除：数据不会真正删除，只是标记为已删除状态
	// 好处：符合数据保留政策，历史数据完整性保证
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// ========== 关联关系 ==========
	// 以下字段不会存储在数据库中，仅用于GORM的关联查询
	
	// Room 当前所住房间信息
	// 多对一关系（可选）
	Room *RoomInfo `gorm:"foreignKey:RoomID;references:ID" json:"room,omitempty"`
	
	// Order 当前关联的订单信息
	// 多对一关系（可选）
	Order *OrderMain `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
	
	// Member 会员信息
	// 多对一关系（可选，仅会员客人有值）
	Member *Member `gorm:"foreignKey:MemberID;references:ID" json:"member,omitempty"`
	
	// Orders 该客人的所有历史订单
	// 一对多关系：一个客人可以有多个订单
	// 用途：查询客人的入住历史
	Orders []OrderMain `gorm:"foreignKey:GuestID;references:ID" json:"orders,omitempty"`
	
	// FinancialFlows 该客人相关的财务流水
	// 一对多关系：一个客人可以产生多条财务记录
	// 用途：查询客人的消费记录、欠款情况
	FinancialFlows []FinancialFlow `gorm:"foreignKey:GuestID;references:ID" json:"financial_flows,omitempty"`
}

// ==================== 表名配置 ====================

// TableName 指定数据库表名
//
// GORM会自动调用此方法获取表名，用于生成SQL语句
//
// 返回：
//   - string: 数据库表名 "guest_info"
func (GuestInfo) TableName() string {
	return "guest_info"
}
