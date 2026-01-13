package model

import (
	"time"

	"gorm.io/gorm"
)

// OrderInfo 订单表
type OrderInfo struct {
	ID            int64          `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	OrderNo       string         `gorm:"type:varchar(50);not null;uniqueIndex:idx_order_no;comment:订单号（唯一，如：2018223233）" json:"order_no"`
	BranchID      int            `gorm:"not null;index:idx_branch_id;comment:关联分店ID（hotel_branch.id）" json:"branch_id"`
	ChannelID     int            `gorm:"not null;index:idx_channel_id;comment:关联渠道ID（booking_channel.id）" json:"channel_id"`
	GuestID       int            `gorm:"not null;index:idx_guest_id;comment:关联客人ID（guest_info.id）" json:"guest_id"`
	RoomID        int            `gorm:"not null;index:idx_room_id;comment:关联房间ID（room_info.id）" json:"room_id"`
	CheckinDate   time.Time      `gorm:"type:date;not null;index:idx_checkin_checkout,priority:1;comment:入住日期" json:"checkin_date"`
	CheckoutDate  time.Time      `gorm:"type:date;not null;index:idx_checkin_checkout,priority:2;comment:离店日期" json:"checkout_date"`
	OrderAmount   float64        `gorm:"type:decimal(10,2);not null;comment:订单总金额" json:"order_amount"`
	DepositAmount float64        `gorm:"type:decimal(10,2);not null;comment:已收押金" json:"deposit_amount"`
	BalanceAmount float64        `gorm:"type:decimal(10,2);not null;default:0.00;comment:欠补费用" json:"balance_amount"`
	OrderStatus   int8           `gorm:"type:tinyint(1);not null;index:idx_order_status;comment:订单状态（0-已预定，1-已入住，2-已退房，3-已失效）" json:"order_status"`
	BookTime      time.Time      `gorm:"not null;comment:预定时间" json:"book_time"`
	CheckinTime   *time.Time     `gorm:"comment:实际入住时间" json:"checkin_time"`
	CheckoutTime  *time.Time     `gorm:"comment:实际离店时间" json:"checkout_time"`
	Operator      string         `gorm:"type:varchar(50);comment:操作人（用户名）" json:"operator"`
	Remarks       string         `gorm:"type:varchar(500);comment:订单备注" json:"remarks"`
	CreatedAt     time.Time      `gorm:"not null;comment:创建时间" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null;comment:更新时间" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index:idx_deleted_at;comment:软删除时间" json:"deleted_at"`

	// 关联关系
	Branch        *HotelBranch    `gorm:"foreignKey:BranchID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"branch,omitempty"`
	Channel       *BookingChannel `gorm:"foreignKey:ChannelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"channel,omitempty"`
	Guest         *GuestInfo      `gorm:"foreignKey:GuestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"guest,omitempty"`
	Room          *RoomInfo       `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"room,omitempty"`
	FinanceFlows  []FinanceFlow   `gorm:"foreignKey:OrderID" json:"finance_flows,omitempty"`
}

// TableName 指定表名
func (OrderInfo) TableName() string {
	return "order_info"
}
