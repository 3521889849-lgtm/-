package coupon

import (
	"time"

	"gorm.io/gorm"
)

// 发放任务用户明细记录表
type CouponGrantTaskDetail struct {
	gorm.Model
	TaskID      int       `gorm:"type:int;comment:关联coupon_grant_task表的任务ID"`
	UserID      int64     `gorm:"type:int;comment:用户ID"`
	GrantResult int       `gorm:"type:int;comment:发放结果：1=成功 2=失败"`
	FailReason  string    `gorm:"type:varchar(100);comment:失败原因"`
	GrantTime   time.Time `gorm:"type:datetime;comment:发放时间"`
}
