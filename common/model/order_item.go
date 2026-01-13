package model

import (
	"time"

	"gorm.io/gorm"
)

// OrderItem 订单详情表-主订单的明细，一个订单对应多个出行人/门票，一对多关系
type OrderItem struct {
	ID           uint64         `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:订单明细主键ID" json:"id"`
	OrderID      uint64         `gorm:"column:order_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_order_id;comment:关联主订单ID" json:"order_id"`
	TicketTypeID uint64         `gorm:"column:ticket_type_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_ticket_type_id;comment:关联门票类型ID" json:"ticket_type_id"`
	TravelerID   uint64         `gorm:"column:traveler_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_traveler_id;comment:关联出行人ID" json:"traveler_id"`
	TicketName   string         `gorm:"column:ticket_name;type:VARCHAR(100);NOT NULL;comment:门票名称（冗余存储，防止门票名称修改）" json:"ticket_name"`
	SinglePrice  float64        `gorm:"column:single_price;type:DECIMAL(10,2);NOT NULL;comment:单张门票价格" json:"single_price"`
	TicketNum    uint8          `gorm:"column:ticket_num;type:TINYINT UNSIGNED;NOT NULL;default:1;comment:购票数量" json:"ticket_num"`
	ExtFields    *JSON          `gorm:"column:ext_fields;type:JSON;comment:扩展字段，如门票有效期、入园须知等" json:"ext_fields,omitempty"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	Order      *OrderMain  `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
	TicketType *TicketType `gorm:"foreignKey:TicketTypeID;references:ID" json:"ticket_type,omitempty"`
	Traveler   *Traveler   `gorm:"foreignKey:TravelerID;references:ID" json:"traveler,omitempty"`
}

func (OrderItem) TableName() string {
	return "order_item"
}
