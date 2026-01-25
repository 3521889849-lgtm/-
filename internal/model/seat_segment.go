package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type SeatSegmentOccupancy struct {
	ID             uint64         `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:自增主键" json:"id"`
	TrainID        string         `gorm:"column:train_id;type:VARCHAR(32);NOT NULL;index:idx_occ_train_seat,priority:1;index:idx_occ_train_seg,priority:1;comment:车次ID" json:"train_id"`
	SeatID         string         `gorm:"column:seat_id;type:VARCHAR(64);NOT NULL;index:idx_occ_train_seat,priority:2;comment:座位ID" json:"seat_id"`
	FromStation    string         `gorm:"column:from_station;type:VARCHAR(64);NOT NULL;comment:上车站" json:"from_station"`
	ToStation      string         `gorm:"column:to_station;type:VARCHAR(64);NOT NULL;comment:下车站" json:"to_station"`
	FromSeq        uint32         `gorm:"column:from_seq;type:INT;NOT NULL;index:idx_occ_train_seg,priority:2;comment:上车站序号" json:"from_seq"`
	ToSeq          uint32         `gorm:"column:to_seq;type:INT;NOT NULL;index:idx_occ_train_seg,priority:3;comment:下车站序号" json:"to_seq"`
	OrderID        string         `gorm:"column:order_id;type:VARCHAR(64);NOT NULL;index:idx_occ_order;comment:订单ID" json:"order_id"`
	Status         string         `gorm:"column:status;type:VARCHAR(20);NOT NULL;index:idx_occ_status_expire,priority:1;comment:占用状态：LOCKED-锁定，SOLD-已售，CANCELLED-已取消" json:"status"`
	LockExpireTime sql.NullTime   `gorm:"column:lock_expire_time;type:DATETIME;index:idx_occ_status_expire,priority:2;comment:锁定到期时间（仅LOCKED有效）" json:"lock_expire_time,omitempty"`
	CreatedAt      time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`
}
