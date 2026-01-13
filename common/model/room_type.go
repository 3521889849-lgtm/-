package model

import (
	"time"

	"gorm.io/gorm"
)

// RoomType 房型配置表
type RoomType struct {
	ID                int            `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	BranchID          int            `gorm:"not null;index:idx_branch_id;comment:关联分店ID（hotel_branch.id）" json:"branch_id"`
	TypeName          string         `gorm:"type:varchar(50);not null;comment:房型名称（如：大床房、商务标准间）" json:"type_name"`
	MarketPrice       float64        `gorm:"type:decimal(10,2);not null;comment:门市价" json:"market_price"`
	CalendarPrice     float64        `gorm:"type:decimal(10,2);not null;comment:日历价（默认售价）" json:"calendar_price"`
	RoomCount         int            `gorm:"not null;comment:该房型总房间数" json:"room_count"`
	Area              float64        `gorm:"type:decimal(6,2);comment:房间面积（单位：㎡）" json:"area"`
	BedSpec           string         `gorm:"type:varchar(50);comment:床型规格（如：1.8*2.0m、两张1.2m床）" json:"bed_spec"`
	IncludeBreakfast  int8           `gorm:"type:tinyint(1);not null;default:0;comment:是否含早（0-否，1-是）" json:"include_breakfast"`
	IncludeToiletries int8           `gorm:"type:tinyint(1);not null;default:1;comment:是否含洗漱用品（0-否，1-是）" json:"include_toiletries"`
	CancelPolicy      string         `gorm:"type:text;comment:退订规则（如：入住前24小时内不可取消，否则收取1倍房费）" json:"cancel_policy"`
	Status            int8           `gorm:"type:tinyint(1);not null;default:1;comment:状态（0-停用，1-启用，2-维修）" json:"status"`
	CreatedAt         time.Time      `gorm:"not null;comment:创建时间" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"not null;comment:更新时间" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index:idx_deleted_at;comment:软删除时间" json:"deleted_at"`

	// 关联关系
	Branch     *HotelBranch   `gorm:"foreignKey:BranchID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"branch,omitempty"`
	RoomInfos  []RoomInfo     `gorm:"foreignKey:RoomTypeID" json:"room_infos,omitempty"`
	RoomImages []RoomImage    `gorm:"foreignKey:RoomTypeID" json:"room_images,omitempty"`
	Facilities []RoomFacility `gorm:"many2many:room_type_facility_rel;" json:"facilities,omitempty"`
}

// TableName 指定表名
func (RoomType) TableName() string {
	return "room_type"
}
