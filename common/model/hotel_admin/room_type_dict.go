// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了房型字典的数据模型
//
// 功能说明：
//   - 统一维护房型的基础属性和标准
//   - 支撑房源配置功能
//   - 定义房型的默认价格和服务标准
//   - 作为房源管理的基础数据
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 房型字典模型 ====================

// RoomTypeDict 房型字典表
//
// 业务用途：
//   - 房型标准化：统一定义房型的属性和标准
//   - 房源配置：创建房源时选择对应的房型
//   - 价格管理：定义房型的默认价格
//   - 统计分析：按房型统计入住率、收入等
//
// 设计说明：
//   - 房型是房间的分类，如：标准间、豪华套房、商务大床房
//   - 同一房型的房间共享相同的基础属性
//   - 具体房间可以在房型基础上个性化调整
//   - 支持启用/停用，停用后不影响已有房源
type RoomTypeDict struct {
	// ========== 基础字段 ==========
	
	// ID 房型ID，主键，自增
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:房型ID" json:"id"`
	
	// ========== 房型信息 ==========
	
	// RoomTypeName 房型名称，房型的标准名称
	// 示例：
	//   - "标准间"：最基础的房型
	//   - "大床房"：带大床的房间
	//   - "商务大床房"：商务人士专用
	//   - "豪华套房"：高端房型
	//   - "家庭房"：适合家庭居住
	//   - "蜜月套房"：浪漫主题
	// 用途：房型展示、分类查询、统计分析
	RoomTypeName string `gorm:"column:room_type_name;type:VARCHAR(50);NOT NULL;index:idx_room_type_name;comment:房型名称（大床房/商务大床房/标准间等）" json:"room_type_name"`
	
	// ========== 房间规格 ==========
	
	// BedSpec 床型规格，标准床型配置
	// 示例：
	//   - "1张1.8m*2.0m大床"
	//   - "2张1.2m*2.0m单人床"
	//   - "1张2.0m*2.2m特大床"
	//   - "1张大床+1张沙发床"
	// 用途：房型说明、客人选房参考
	BedSpec string `gorm:"column:bed_spec;type:VARCHAR(50);NOT NULL;comment:床型规格（1.8*2.0m等）" json:"bed_spec"`
	
	// Area 标准面积，房型的标准面积，可选
	// 单位：平方米（㎡）
	// 示例：25.50（表示25.5平方米）
	// 说明：具体房间的面积可能略有差异
	Area *float64 `gorm:"column:area;type:DECIMAL(5,2);comment:面积（平方米）" json:"area,omitempty"`
	
	// ========== 服务标准 ==========
	
	// HasBreakfast 是否含早餐，房型的默认早餐配置
	// true: 该房型默认含早餐
	// false: 该房型默认不含早餐
	// 说明：具体房间可以覆盖此配置
	HasBreakfast bool `gorm:"column:has_breakfast;type:BOOLEAN;NOT NULL;default:false;comment:是否含早（0/1）" json:"has_breakfast"`
	
	// HasToiletries 是否提供洗漱用品，房型的默认配置
	// true: 该房型默认提供洗漱用品
	// false: 该房型默认不提供
	// 说明：具体房间可以覆盖此配置
	HasToiletries bool `gorm:"column:has_toiletries;type:BOOLEAN;NOT NULL;default:false;comment:是否提供洗漱用品（0/1）" json:"has_toiletries"`
	
	// ========== 价格信息 ==========
	
	// DefaultPrice 默认门市价，房型的标准价格
	// 单位：元（人民币）
	// 用途：
	//   - 创建房源时的默认价格
	//   - 价格参考基准
	//   - 统计分析时的标准价格
	// 说明：具体房间可以设置不同的价格
	DefaultPrice float64 `gorm:"column:default_price;type:DECIMAL(10,2);NOT NULL;comment:默认门市价" json:"default_price"`
	
	// ========== 状态控制 ==========
	
	// Status 状态，控制房型是否可用
	// 可选值：
	//   - "ACTIVE": 启用，可以创建该房型的房源
	//   - "INACTIVE": 停用，不能创建新房源（已有房源不受影响）
	// 用途：房型管理、系统维护
	Status string `gorm:"column:status;type:VARCHAR(20);NOT NULL;default:'ACTIVE';index:idx_status;comment:状态：ACTIVE-启用，INACTIVE-停用" json:"status"`
	
	// ========== 时间戳 ==========
	
	// CreatedAt 创建时间，房型首次创建的时间
	CreatedAt time.Time `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	
	// UpdatedAt 更新时间，房型最后修改的时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	
	// DeletedAt 软删除时间，非NULL表示已删除
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// ========== 关联关系 ==========
	
	// RoomInfos 该房型下的所有房源
	// 一对多关系：一个房型可以对应多个具体房间
	RoomInfos []RoomInfo `gorm:"foreignKey:RoomTypeID;references:ID" json:"room_infos,omitempty"`
	
	// CancellationPolicies 该房型适用的退订政策
	// 一对多关系：一个房型可以有多个退订政策（如：不同季节不同政策）
	CancellationPolicies []CancellationPolicy `gorm:"foreignKey:RoomTypeID;references:ID" json:"cancellation_policies,omitempty"`
	
	// Orders 该房型的订单记录
	// 一对多关系：用于统计分析
	Orders []OrderMain `gorm:"foreignKey:RoomTypeID;references:ID" json:"orders,omitempty"`
	
	// SalesStats 该房型的销售统计数据
	// 一对多关系：按房型统计销售情况
	SalesStats []SalesStatistics `gorm:"foreignKey:RoomTypeID;references:ID" json:"sales_stats,omitempty"`
}

// ==================== 表名配置 ====================

// TableName 指定数据库表名
//
// 返回：数据库表名 "room_type_dict"
func (RoomTypeDict) TableName() string {
	return "room_type_dict"
}
