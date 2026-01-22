package ticketapp

import (
	"database/sql"
	"errors"
	"example_shop/internal/ticket_service/model"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AllocateSeats 负责在事务内按席别分配可用座席，并返回分配到的 seat 列表与总价。
//
// 该函数代表“票务域对外的座席分配边界能力”，订单域通过它来完成选座/锁座的核心步骤，
// 从而避免订单域直接实现复杂的座席可用性判定 SQL。
func AllocateSeats(tx *gorm.DB, trainID string, need map[string]int, fromSeq, toSeq uint32, lockExpire time.Time) ([]*model.SeatInfo, float64, error) {
	if tx == nil {
		return nil, 0, errors.New("tx is nil")
	}
	if fromSeq == 0 || toSeq == 0 || fromSeq >= toSeq {
		return nil, 0, errors.New("invalid segment")
	}

	allocated := make([]*model.SeatInfo, 0, 8)
	total := 0.0
	now := time.Now()

	_ = tx.Model(&model.SeatSegmentOccupancy{}).
		Where("status = ? AND lock_expire_time IS NOT NULL AND lock_expire_time < ?", "LOCKED", now).
		Updates(map[string]interface{}{"status": "CANCELLED", "lock_expire_time": sql.NullTime{Valid: false}}).Error

	for seatType, cnt := range need {
		if cnt <= 0 {
			continue
		}
		var seats []model.SeatInfo
		if err := tx.Table("seat_infos si").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("si.train_id = ? AND si.seat_type = ? AND si.status = ? AND si.deleted_at IS NULL", trainID, seatType, "AVAILABLE").
			Where(
				"NOT EXISTS (SELECT 1 FROM seat_segment_occupancies o WHERE o.deleted_at IS NULL AND o.train_id = si.train_id AND o.seat_id = si.seat_id AND o.status IN ('SOLD','LOCKED') AND (o.status <> 'LOCKED' OR o.lock_expire_time > ?) AND o.from_seq < ? AND o.to_seq > ?)",
				now, toSeq, fromSeq,
			).
			Order("si.seat_id ASC").
			Limit(cnt).
			Scan(&seats).Error; err != nil {
			return nil, 0, err
		}
		if len(seats) < cnt {
			return nil, 0, fmt.Errorf("座位余量不足: %s", seatType)
		}
		for i := range seats {
			allocated = append(allocated, &seats[i])
			total += seats[i].SeatPrice
		}
	}

	return allocated, total, nil
}
