package model

import (
	"time"

	"gorm.io/gorm"
)

// SpotInfo 景点信息表-门票查询的基础数据
type SpotInfo struct {
	ID           uint64         `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:景点主键ID" json:"id"`
	MerchantID   uint64         `gorm:"column:merchant_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_merchant_id;comment:所属商家ID" json:"merchant_id"`
	SpotName     string         `gorm:"column:spot_name;type:VARCHAR(100);NOT NULL;index:idx_spot_name;comment:景点名称" json:"spot_name"`
	SpotDesc     *string        `gorm:"column:spot_desc;type:TEXT;comment:景点介绍" json:"spot_desc,omitempty"`
	Province     string         `gorm:"column:province;type:VARCHAR(30);NOT NULL;index:idx_province_city;comment:省份" json:"province"`
	City         string         `gorm:"column:city;type:VARCHAR(30);NOT NULL;index:idx_province_city;comment:城市" json:"city"`
	Address      string         `gorm:"column:address;type:VARCHAR(255);NOT NULL;comment:景点详细地址" json:"address"`
	CoverImg     string         `gorm:"column:cover_img;type:VARCHAR(512);NOT NULL;comment:景点封面图" json:"cover_img"`
	OpenTime     string         `gorm:"column:open_time;type:VARCHAR(100);NOT NULL;comment:开放时间" json:"open_time"`
	ContactPhone *string        `gorm:"column:contact_phone;type:VARCHAR(20);comment:景点联系电话" json:"contact_phone,omitempty"`
	ExtFields    *JSON          `gorm:"column:ext_fields;type:JSON;comment:扩展字段，如评分、特色标签等" json:"ext_fields,omitempty"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	Merchant    *SysMerchant `gorm:"foreignKey:MerchantID;references:ID" json:"merchant,omitempty"`
	TicketTypes []TicketType `gorm:"foreignKey:SpotID;references:ID" json:"ticket_types,omitempty"`
}

func (SpotInfo) TableName() string {
	return "spot_info"
}
