package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// SeatInfo 座位表
type SeatInfo struct {
	ID             string         `gorm:"column:seat_id;type:VARCHAR(64);primaryKey;comment:座位唯一标识（格式：train_id+车厢号+座位号，如G1234_03_01A）" json:"seat_id"`
	TrainID        string         `gorm:"column:train_id;type:VARCHAR(32);NOT NULL;index:idx_train_seat_status,priority:1;comment:车次编号，关联train_info.train_id" json:"train_id"`
	CarriageNum    string         `gorm:"column:carriage_num;type:VARCHAR(8);NOT NULL;comment:车厢号（如03、10）" json:"carriage_num"`
	SeatNum        string         `gorm:"column:seat_num;type:VARCHAR(8);NOT NULL;comment:座位号（如01A、12下）" json:"seat_num"`
	SeatType       string         `gorm:"column:seat_type;type:VARCHAR(16);NOT NULL;index:idx_train_seat_status,priority:2;comment:座位类型（硬座、软座、硬卧、软卧、商务座等）" json:"seat_type"`
	SeatPrice      float64        `gorm:"column:seat_price;type:DECIMAL(10,2);NOT NULL;comment:座位单价" json:"seat_price"`
	Status         string         `gorm:"column:status;type:VARCHAR(20);NOT NULL;default:'AVAILABLE';index:idx_seat_status;index:idx_train_seat_status,priority:3;comment:座位状态：AVAILABLE-可售，LOCKED-已锁定，SOLD-已售出，REFUNDED-已退票" json:"status"`
	LockExpireTime sql.NullTime   `gorm:"column:lock_expire_time;type:DATETIME;index:idx_lock_expire_time;comment:锁定过期时间（状态为LOCKED时有效）" json:"lock_expire_time,omitempty"`
	LockOrderID    sql.NullString `gorm:"column:lock_order_id;type:VARCHAR(64);comment:锁定关联的订单号（状态为LOCKED时有效）" json:"lock_order_id,omitempty"`
	CreatedAt      time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`
}

// TicketInventoryLog 余票变更日志表
type TicketInventoryLog struct {
	ID             uint64         `gorm:"column:log_id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:日志ID" json:"log_id"`
	TrainID        string         `gorm:"column:train_id;type:VARCHAR(32);NOT NULL;comment:车次编号，关联train_info.train_id" json:"train_id"`
	SeatType       string         `gorm:"column:seat_type;type:VARCHAR(16);NOT NULL;comment:座位类型" json:"seat_type"`
	ChangeType     string         `gorm:"column:change_type;type:VARCHAR(20);NOT NULL;comment:变更类型：LOCK-锁定，RELEASE-释放，SALE-售卖，REFUND-退票" json:"change_type"`
	BeforeCount    uint32         `gorm:"column:before_count;type:INT;NOT NULL;comment:变更前余票数量" json:"before_count"`
	AfterCount     uint32         `gorm:"column:after_count;type:INT;NOT NULL;comment:变更后余票数量" json:"after_count"`
	RelatedOrderID sql.NullString `gorm:"column:related_order_id;type:VARCHAR(64);comment:关联订单号（变更触发源）" json:"related_order_id,omitempty"`
	TraceID        string         `gorm:"column:trace_id;type:VARCHAR(64);NOT NULL;comment:链路ID" json:"trace_id"`
	CreatedAt      time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:变更时间" json:"created_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`
}

// TicketRuleConfig 购票规则配置表
type TicketRuleConfig struct {
	ID              string         `gorm:"column:rule_id;type:VARCHAR(64);primaryKey;comment:规则ID（如rule_limit_ticket、rule_refund_fee）" json:"rule_id"`
	RuleType        string         `gorm:"column:rule_type;type:VARCHAR(32);NOT NULL;comment:规则类型：LIMIT-限购规则，PAY_TIMEOUT-支付超时规则，REFUND_FEE-退票手续费规则，SEAT_ALLOC-座位分配策略" json:"rule_type"`
	RuleParam       *JSON          `gorm:"column:rule_param;type:JSON;NOT NULL;comment:规则参数（JSON格式，如{\"single_train_limit\":3,\"daily_limit\":5}）" json:"rule_param"`
	ApplyCities     string         `gorm:"column:apply_cities;type:VARCHAR(512);default:'all';comment:适用城市（逗号分隔，all表示全部）" json:"apply_cities"`
	ApplyUserLevels string         `gorm:"column:apply_user_levels;type:VARCHAR(64);default:'ORDINARY,VIP';comment:适用用户等级（逗号分隔，ORDINARY-普通用户，VIP-VIP用户）" json:"apply_user_levels"`
	ApplyTrainTypes string         `gorm:"column:apply_train_types;type:VARCHAR(128);default:'all';comment:适用车次类型（逗号分隔，all表示全部）" json:"apply_train_types"`
	GrayRatio       uint32         `gorm:"column:gray_ratio;type:INT;default:100;comment:灰度比例（0-100，100表示全量生效）" json:"gray_ratio"`
	Status          string         `gorm:"column:status;type:VARCHAR(20);NOT NULL;default:'ENABLED';comment:规则状态：ENABLED-启用，DISABLED-禁用" json:"status"`
	Version         uint32         `gorm:"column:version;type:INT;NOT NULL;default:1;comment:规则版本（用于回滚）" json:"version"`
	Operator        string         `gorm:"column:operator;type:VARCHAR(64);NOT NULL;comment:操作人（后台管理员ID）" json:"operator"`
	CreatedAt       time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`
}
