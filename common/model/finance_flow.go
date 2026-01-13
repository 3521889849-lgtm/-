package model

import (
	"time"

	"gorm.io/gorm"
)

// FinanceFlow 收支流水表
type FinanceFlow struct {
	ID              int64          `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	OrderID         int64          `gorm:"not null;index:idx_order_id;comment:关联订单ID（order_info.id）" json:"order_id"`
	BranchID        int            `gorm:"not null;index:idx_branch_id;comment:关联分店ID（hotel_branch.id）" json:"branch_id"`
	RoomID          int            `gorm:"not null;index:fk_flow_room;comment:关联房间ID（room_info.id）" json:"room_id"`
	GuestID         int            `gorm:"not null;index:idx_guest_id;comment:关联客人ID（guest_info.id）" json:"guest_id"`
	FlowType        int8           `gorm:"type:tinyint(1);not null;index:idx_flow_type;comment:流水类型（0-收入，1-支出）" json:"flow_type"`
	IncomeItem      string         `gorm:"type:varchar(50);not null;comment:收入项目（如：房费、押金）" json:"income_item"`
	ExpenditureItem string         `gorm:"type:varchar(50);comment:支出项目（如：退款、赔偿）" json:"expenditure_item"`
	Amount          float64        `gorm:"type:decimal(10,2);not null;comment:金额" json:"amount"`
	PaymentMethod   string         `gorm:"type:varchar(20);not null;comment:支付方式（如：现金、支付宝、微信、银联）" json:"payment_method"`
	Operator        string         `gorm:"type:varchar(50);not null;comment:操作人（用户名）" json:"operator"`
	FlowTime        time.Time      `gorm:"not null;index:idx_flow_time;comment:流水发生时间" json:"flow_time"`
	Remarks         string         `gorm:"type:varchar(500);comment:备注" json:"remarks"`
	CreatedAt       time.Time      `gorm:"not null;comment:创建时间" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"not null;comment:更新时间" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index:idx_deleted_at;comment:软删除时间" json:"deleted_at"`

	// 关联关系
	Order  *OrderInfo   `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"order,omitempty"`
	Branch *HotelBranch `gorm:"foreignKey:BranchID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"branch,omitempty"`
	Room   *RoomInfo    `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"room,omitempty"`
	Guest  *GuestInfo   `gorm:"foreignKey:GuestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"guest,omitempty"`
}

// TableName 指定表名
func (FinanceFlow) TableName() string {
	return "finance_flow"
}
