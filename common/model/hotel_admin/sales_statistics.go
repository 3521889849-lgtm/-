// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了销售统计的数据模型
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// SalesStatistics 销售统计表
//
// 业务用途：
//   - 按日期统计房间销售数据
//   - 按房型、渠道等维度统计
//   - 支持销售报表展示
//   - 分析收入构成（房费 + 非房费）
type SalesStatistics struct {
	ID         uint64    `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:统计ID" json:"id"`
	StatDate   time.Time `gorm:"column:stat_date;type:DATE;NOT NULL;index:idx_stat_date;comment:统计日期" json:"stat_date"`
	BranchID   uint64    `gorm:"column:branch_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_branch_id;comment:分店ID（外键，关联分店表）" json:"branch_id"`
	RoomTypeID *uint64   `gorm:"column:room_type_id;type:BIGINT UNSIGNED;index:idx_room_type_id;comment:房型ID（外键，关联房型表）" json:"room_type_id,omitempty"`
	
	// SalesType 销售类型，如："房间销售"、"渠道销售"
	SalesType string `gorm:"column:sales_type;type:VARCHAR(30);NOT NULL;index:idx_sales_type;comment:销售类型（房间销售/渠道销售）" json:"sales_type"`
	
	// 统计数据
	SalesCount    uint32  `gorm:"column:sales_count;type:INT UNSIGNED;NOT NULL;default:0;comment:销售数量" json:"sales_count"`
	SalesAmount   float64 `gorm:"column:sales_amount;type:DECIMAL(10,2);NOT NULL;default:0.00;comment:销售金额" json:"sales_amount"`
	NonRoomIncome float64 `gorm:"column:non_room_income;type:DECIMAL(10,2);NOT NULL;default:0.00;comment:非房费收入金额" json:"non_room_income"`
	
	// 时间戳
	StatTime  time.Time      `gorm:"column:stat_time;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:统计时间" json:"stat_time"`
	CreatedAt time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	Branch   *HotelBranch  `gorm:"foreignKey:BranchID;references:ID" json:"branch,omitempty"`       // 分店信息
	RoomType *RoomTypeDict `gorm:"foreignKey:RoomTypeID;references:ID" json:"room_type,omitempty"` // 房型信息
}

// TableName 指定数据库表名
func (SalesStatistics) TableName() string {
	return "sales_statistics"
}
