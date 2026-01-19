// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了房间设施字典的数据模型
//
// 功能说明：
//   - 维护可配置的房间设施选项
//   - 支持房源设施勾选功能
//   - 统一管理设施名称和说明
//   - 支持设施的增删改查
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 设施字典模型 ====================

// FacilityDict 房间设施字典表
//
// 业务用途：
//   - 设施标准化：统一定义所有可用的房间设施
//   - 设施配置：为房间勾选具备的设施
//   - 设施展示：在房间详情页展示设施列表
//   - 设施筛选：客人可以按设施筛选房间
//   - 设施统计：统计各设施的配备率
//
// 设计说明：
//   - 这是一个字典表，集中管理所有设施
//   - 通过中间表（room_facility_relation）关联到具体房间
//   - 新增设施时直接在此表添加即可
//   - 支持启用/停用，停用后不影响已配置的房间
type FacilityDict struct {
	// ========== 基础字段 ==========
	
	// ID 设施ID，主键，自增
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:设施ID" json:"id"`
	
	// ========== 设施信息 ==========
	
	// FacilityName 设施名称，设施的标准名称
	// 常见设施示例：
	//
	// 基础设施：
	//   - "无线WiFi"：免费WiFi
	//   - "有线网络"：网线上网
	//   - "空调"：冷暖空调
	//   - "暖气"：冬季供暖
	//   - "电视"：液晶电视
	//   - "电话"：客房电话
	//
	// 卫浴设施：
	//   - "24小时热水"：热水供应
	//   - "独立卫浴"：独立卫生间
	//   - "浴缸"：带浴缸
	//   - "淋浴"：淋浴设备
	//   - "吹风机"：电吹风
	//   - "洗漱用品"：牙刷牙膏等
	//
	// 电器设备：
	//   - "冰箱"：小冰箱
	//   - "微波炉"：加热设备
	//   - "电热水壶"：烧水壶
	//   - "保险箱"：贵重物品保管
	//
	// 舒适设施：
	//   - "沙发"：休息沙发
	//   - "书桌"：办公桌
	//   - "衣柜"：储物空间
	//   - "拖鞋"：一次性拖鞋
	//   - "窗户"：带窗户（采光通风）
	//
	// 特色设施：
	//   - "智能门锁"：电子门锁
	//   - "投影仪"：大屏投影
	//   - "智能音箱"：语音控制
	//   - "空气净化器"：空气净化
	//   - "胶囊咖啡机"：现磨咖啡
	//   - "按摩椅"：按摩设备
	//
	// 用途：设施展示、分类查询、筛选条件
	FacilityName string `gorm:"column:facility_name;type:VARCHAR(100);NOT NULL;index:idx_facility_name;comment:设施名称（无线wifi/空调/冰箱等）" json:"facility_name"`
	
	// Description 设施描述，设施的详细说明，可选
	// 用途：
	//   - 补充说明设施的规格型号
	//   - 说明设施的使用方法
	//   - 标注设施的特殊说明
	//
	// 示例：
	//   - "100M光纤免费WiFi，全覆盖无死角"
	//   - "智能变频空调，可独立控制温度"
	//   - "LG 43英寸液晶电视，支持投屏"
	//   - "King Koil床垫，舒适睡眠体验"
	Description *string `gorm:"column:description;type:VARCHAR(255);comment:设施描述" json:"description,omitempty"`
	
	// ========== 状态控制 ==========
	
	// Status 状态，控制设施是否可用
	// 可选值：
	//   - "ACTIVE": 启用，可以为房间配置此设施
	//   - "INACTIVE": 停用，不能配置新的（已配置的不受影响）
	//
	// 停用场景：
	//   - 设施淘汰：某些设施已不再提供（如：传真机）
	//   - 设施合并：多个设施合并为一个（如：有线电视并入智能电视）
	//   - 临时下线：设施暂时不对外宣传
	//
	// 用途：设施管理、系统维护
	Status string `gorm:"column:status;type:VARCHAR(20);NOT NULL;default:'ACTIVE';index:idx_status;comment:状态：ACTIVE-启用，INACTIVE-停用" json:"status"`
	
	// ========== 时间戳 ==========
	
	// CreatedAt 创建时间，设施首次添加的时间
	CreatedAt time.Time `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	
	// UpdatedAt 更新时间，设施最后修改的时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	
	// DeletedAt 软删除时间，非NULL表示已删除
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// ========== 关联关系 ==========
	
	// RoomFacilityRelations 房间-设施关联关系
	// 一对多关系：一个设施可以配置给多个房间
	// 用途：查询哪些房间配备了该设施
	RoomFacilityRelations []RoomFacilityRelation `gorm:"foreignKey:FacilityID;references:ID" json:"room_facility_relations,omitempty"`
}

// ==================== 表名配置 ====================

// TableName 指定数据库表名
//
// 返回：数据库表名 "facility_dict"
func (FacilityDict) TableName() string {
	return "facility_dict"
}
