package job

import (
	"database/sql"
	"example_shop/common/db"
	"example_shop/internal/ticket_service/model"
	"time"

	"gorm.io/gorm"
)

func StartOrderCleanup(stop <-chan struct{}) {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				cleanupOnce()
			case <-stop:
				return
			}
		}
	}()
}

func cleanupOnce() {
	now := time.Now()

	for {
		var ids []string
		err := db.MysqlDB.Model(&model.OrderInfo{}).
			Select("order_id").
			Where("order_status IN ('PENDING_PAY','PAYING') AND pay_deadline IS NOT NULL AND pay_deadline < ?", now).
			Limit(500).
			Scan(&ids).Error
		if err != nil || len(ids) == 0 {
			return
		}

		_ = db.MysqlDB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&model.SeatSegmentOccupancy{}).
				Where("order_id IN ? AND status = ?", ids, "LOCKED").
				Updates(map[string]any{
					"status":           "CANCELLED",
					"lock_expire_time": sql.NullTime{Valid: false},
				}).Error; err != nil {
				return err
			}

			if err := tx.Model(&model.OrderInfo{}).
				Where("order_id IN ? AND order_status IN ('PENDING_PAY','PAYING')", ids).
				Updates(map[string]any{"order_status": "CANCELLED"}).Error; err != nil {
				return err
			}
			return nil
		})
	}
}
