package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// TicketType 门票类型配置表-含退改规则/库存/价格，下单核心关联表，乐观锁防超卖
type TicketType struct {
	ID             uint64         `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:门票类型主键ID" json:"id"`
	SpotID         uint64         `gorm:"column:spot_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_spot_id;comment:所属景点ID" json:"spot_id"`
	TicketName     string         `gorm:"column:ticket_name;type:VARCHAR(100);NOT NULL;comment:门票名称（如成人票、儿童票、套票）" json:"ticket_name"`
	Price          float64        `gorm:"column:price;type:DECIMAL(10,2);NOT NULL;comment:门票售价" json:"price"`
	OriginalPrice  float64        `gorm:"column:original_price;type:DECIMAL(10,2);NOT NULL;comment:门票原价" json:"original_price"`
	Stock          uint32         `gorm:"column:stock;type:INT UNSIGNED;NOT NULL;default:0;comment:门票库存数量" json:"stock"`
	Version        uint32         `gorm:"column:version;type:INT UNSIGNED;NOT NULL;default:1;comment:乐观锁版本号，扣库存必用，防超卖核心字段" json:"version"`
	ValidStartTime sql.NullTime   `gorm:"column:valid_start_time;type:DATE;NOT NULL;comment:门票有效期开始时间" json:"valid_start_time"`
	ValidEndTime   sql.NullTime   `gorm:"column:valid_end_time;type:DATE;NOT NULL;comment:门票有效期结束时间" json:"valid_end_time"`
	TicketStatus   string         `gorm:"column:ticket_status;type:VARCHAR(20);NOT NULL;default:'ON_SALE';index:idx_ticket_status;comment:门票状态：ON_SALE-在售，OFF_SALE-下架，STOCK_OUT-售罄" json:"ticket_status"`
	RefundRule     string         `gorm:"column:refund_rule;type:TEXT;NOT NULL;comment:退改规则（如：游玩前24小时可退，逾期不退，改期限1次）" json:"refund_rule"`
	UseRule        string         `gorm:"column:use_rule;type:TEXT;NOT NULL;comment:使用规则（如：实名制入园、有效期内通用）" json:"use_rule"`
	ExtFields      *JSON          `gorm:"column:ext_fields;type:JSON;comment:扩展字段，如适用人群、免票政策等" json:"ext_fields,omitempty"`
	CreatedAt      time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// 关联关系
	Spot       *SpotInfo   `gorm:"foreignKey:SpotID;references:ID" json:"spot,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:TicketTypeID;references:ID" json:"order_items,omitempty"`
}

func (TicketType) TableName() string {
	return "ticket_type"
}
