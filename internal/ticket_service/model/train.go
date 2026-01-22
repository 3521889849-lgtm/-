package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// TrainInfo 车次表
type TrainInfo struct {
	ID               string         `gorm:"column:train_id;type:VARCHAR(32);primaryKey;comment:车次编号（如G1234），唯一标识" json:"train_id"`
	TrainCode        string         `gorm:"column:train_code;type:VARCHAR(16);index:idx_train_code;comment:车次编码（如G1001），用于展示/搜索（多日期复用）" json:"train_code"`
	ServiceDate      sql.NullTime   `gorm:"column:service_date;type:DATE;index:idx_train_route_service,priority:3;comment:运行日期（按自然日），用于按天查询" json:"service_date,omitempty"`
	TrainType        string         `gorm:"column:train_type;type:VARCHAR(16);NOT NULL;comment:车次类型（高铁、动车、普速列车等）" json:"train_type"`
	DepartureStation string         `gorm:"column:departure_station;type:VARCHAR(64);NOT NULL;index:idx_train_route_date,priority:1;comment:出发站" json:"departure_station"`
	ArrivalStation   string         `gorm:"column:arrival_station;type:VARCHAR(64);NOT NULL;index:idx_train_route_date,priority:2;comment:到达站" json:"arrival_station"`
	DepartureTime    time.Time      `gorm:"column:departure_time;type:DATETIME;NOT NULL;index:idx_departure_time;index:idx_train_route_date,priority:3;index:idx_train_route_service,priority:4;comment:发车时间" json:"departure_time"`
	ArrivalTime      time.Time      `gorm:"column:arrival_time;type:DATETIME;NOT NULL;comment:到达时间" json:"arrival_time"`
	RuntimeMinutes   uint32         `gorm:"column:runtime_minutes;type:INT;NOT NULL;comment:运行时长（分钟）" json:"runtime_minutes"`
	SeatLayout       *JSON          `gorm:"column:seat_layout;type:JSON;NOT NULL;comment:座位布局（如硬座118座、硬卧66座，JSON存储各类型座位总数）" json:"seat_layout"`
	Status           string         `gorm:"column:status;type:VARCHAR(20);NOT NULL;default:'NORMAL';index:idx_train_status;comment:车次状态：NORMAL-正常运营，STOPPED-停运，TEMP-临时加开" json:"status"`
	CreatedAt        time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`
}

// TrainStationPass 车次途径站点表
type TrainStationPass struct {
	ID            uint64         `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:自增主键" json:"id"`
	TrainID       string         `gorm:"column:train_id;type:VARCHAR(32);NOT NULL;comment:车次编号，关联train_info.train_id" json:"train_id"`
	StationName   string         `gorm:"column:station_name;type:VARCHAR(64);NOT NULL;comment:途经站点名称" json:"station_name"`
	Sequence      uint32         `gorm:"column:sequence;type:INT;NOT NULL;comment:站点顺序（1-出发站，n-到达站）" json:"sequence"`
	ArrivalTime   sql.NullTime   `gorm:"column:arrival_time;type:DATETIME;comment:到站时间（出发站为NULL）" json:"arrival_time,omitempty"`
	DepartureTime sql.NullTime   `gorm:"column:departure_time;type:DATETIME;comment:离站时间（到达站为NULL）" json:"departure_time,omitempty"`
	StopMinutes   uint32         `gorm:"column:stop_minutes;type:INT;default:0;comment:停靠时长（分钟）" json:"stop_minutes"`
	CreatedAt     time.Time      `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`
}
