package model

import "time"

// ShiftChangeRecord 交接班记录表
type ShiftChangeRecord struct {
	ID               int       `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	BranchID         int       `gorm:"not null;index:idx_branch_id;comment:关联分店ID（hotel_branch.id）" json:"branch_id"`
	OnDutyPerson     string    `gorm:"type:varchar(50);not null;comment:当班人" json:"on_duty_person"`
	OffDutyPerson    string    `gorm:"type:varchar(50);not null;comment:接班人" json:"off_duty_person"`
	OnDutyTime       time.Time `gorm:"not null;index:idx_on_off_time,priority:1;comment:上班时间" json:"on_duty_time"`
	OffDutyTime      time.Time `gorm:"not null;index:idx_on_off_time,priority:2;comment:下班时间" json:"off_duty_time"`
	TotalIncome      float64   `gorm:"type:decimal(12,2);not null;comment:当班总收入" json:"total_income"`
	TotalExpenditure float64   `gorm:"type:decimal(12,2);not null;comment:当班总支出" json:"total_expenditure"`
	BalanceAmount    float64   `gorm:"type:decimal(12,2);not null;comment:当班结余" json:"balance_amount"`
	HandoverRemarks  string    `gorm:"type:text;comment:交接备注" json:"handover_remarks"`
	CreatedAt        time.Time `gorm:"not null;comment:创建时间" json:"created_at"`
	UpdatedAt        time.Time `gorm:"not null;comment:更新时间" json:"updated_at"`

	// 关联关系
	Branch *HotelBranch `gorm:"foreignKey:BranchID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"branch,omitempty"`
}

// TableName 指定表名
func (ShiftChangeRecord) TableName() string {
	return "shift_change_record"
}
