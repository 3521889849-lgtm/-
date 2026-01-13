package model

import "time"

// MemberLevel 会员等级表
type MemberLevel struct {
	ID           int       `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	LevelName    string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_level_name;comment:等级名称（如：普通会员、黄金会员、钻石会员）" json:"level_name"`
	LevelDesc    string    `gorm:"type:varchar(255);comment:等级描述" json:"level_desc"`
	PointsRule   string    `gorm:"type:text;comment:积分获取规则（如：消费1元积1分，入住1晚额外积10分）" json:"points_rule"`
	DiscountRate float64   `gorm:"type:decimal(5,2);not null;comment:房价折扣（如0.95=95折）" json:"discount_rate"`
	Status       int8      `gorm:"type:tinyint(1);not null;default:1;comment:状态（0-停用，1-启用）" json:"status"`
	CreatedAt    time.Time `gorm:"not null;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;comment:更新时间" json:"updated_at"`

	// 关联关系
	Members []MemberInfo `gorm:"foreignKey:LevelID" json:"members,omitempty"`
}

// TableName 指定表名
func (MemberLevel) TableName() string {
	return "member_level"
}
