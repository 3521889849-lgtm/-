package coupon

import (
	"time"

	"gorm.io/gorm"
)

// CouponGrantTask 优惠券发放任务表
type CouponGrantTask struct {
	gorm.Model
	CouponID          int       `gorm:"type:int unsigned;comment:关联coupon_base表的优惠券ID"`
	GrantType         int       `gorm:"type:int;comment:发放类型：1=全员发放 2=定向发放 3=新人注册发放 4=邀请有礼发放"`
	GrantScope        int       `gorm:"type:int;comment:发放范围：1=全部用户 2=按注册时间 3=按用户等级 4=按用户标签"`
	UserLevelRange    string    `gorm:"type:varchar(50);comment:用户等级区间（如“1-3”，按等级发放时有效）"`
	RegisterTimeStart time.Time `gorm:"type:datetime;comment:注册时间起始（按注册时间发放时有效）"`
	RegisterTimeEnd   time.Time `gorm:"type:datetime;comment:注册时间结束（按注册时间发放时有效）"`
	UserTag           string    `gorm:"type:varchar(50);comment:用户标签（如“高频用户”，按标签发放时有效）"`
	IsNotify          int       `gorm:"type:int;comment:是否发送通知：0=否 1=是"`
	NotifyContent     string    `gorm:"type:varchar(200);comment:通知内容（短信/站内信）"`
	EstimateUserCount int       `gorm:"type:int;comment:预估发放人数"`
	SuccessCount      int       `gorm:"type:int;comment:发放成功数"`
	FailCount         int       `gorm:"type:int;comment:发放失败数"`
	TaskStatus        int       `gorm:"type:int;comment:任务状态：0=发放中 1=发放完成 2=发放失败"`
	OperateBy         string    `gorm:"type:varchar(50);comment:操作人（管理员ID）"`
}
