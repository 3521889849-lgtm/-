package coupon

import (
	"time"

	"gorm.io/gorm"
)

// 优惠券基础配置表
type CouponBase struct {
	gorm.Model
	CouponName             string    `gorm:"type:varchar(100);comment:优惠卷名称"`
	CouponType             int       `gorm:"type:int;comment:优惠券类型：1=新人注册券 2=限时活动券 3=邀请有礼券 4=满减券 5=折扣券 6=无门槛券 7=全员福利券"`
	Denomination           float64   `gorm:"type:decimal(10,2);comment:面额：满减/无门槛券=减免额；折扣券=折扣比例（如8.0代表8折）"`
	FullReductionThreshold float64   `gorm:"type:decimal(10,2);comment:满减门槛（仅满减券有效，如“满100减30”的100）"`
	MaxDiscountAmount      float64   `gorm:"type:decimal(10,2);comment:最大优惠上限（仅折扣券有效，如8折最高减200元）"`
	ApplyRangeType         int       `gorm:"type:int;comment:使用范围类型：1=全部适用 2=部分指定"`
	ValidityType           int       `gorm:"type:int;comment:有效期类型：1=固定时间 2=领取后N天有效"`
	ValidityStart          time.Time `gorm:"type:datetime;comment:有效期开始时间（固定时间类型有效）"`
	ValidityEnd            time.Time `gorm:"type:datetime;comment:有效期结束时间（固定时间类型有效）"`
	ValidDays              int       `gorm:"type:int;comment:领取后有效天数（领取后N天类型有效）"`
	TotalStock             int       `gorm:"type:int;comment:总库存（0代表不限量）"`
	UserReceiveLimit       int       `gorm:"type:int;comment:单用户领取上限（每人限领张数）"`
	IsStack                int       `gorm:"type:int;comment:是否允许叠加：0=不允许 1=允许"`
	Status                 int       `gorm:"type:int;comment:优惠券状态：0=未生效 1=生效中 2=已停用 3=已过期"`
	CreateBy               string    `gorm:"type:varchar(50);comment:创建人（管理员ID）"`
}
