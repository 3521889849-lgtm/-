package orderapp

import (
	"context"
	"errors"
	"example_shop/common/db"
	model2 "example_shop/internal/model"
	"example_shop/internal/ticket_service/app/shared"
	"example_shop/internal/ticket_service/app/ticketapp"
	kitexuser "example_shop/kitex_gen/user"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ChangeOrder 改签用例：创建一张新订单并将原订单置为 CHANGED（简化版）。
//
// 简化版改签策略（本项目假设）：
// - 仅支持“已出票”订单改签：ISSUED -> CHANGED，并生成一张新订单（新订单直接视为已出票 ISSUED）
// - 不支持补差价：如果新订单金额高于原订单金额，则直接拒绝
// - 支持退差价：如果新订单金额低于原订单金额，则把差额记录到原订单的 refund_amount（并标记 REFUNDED）
//
// 一致性要求：
// - 必须在同一事务内完成：
//   1) 为新车次分配座位并写入新订单（含新占用 SOLD、乘客与关联）
//   2) 释放原订单的已售占用（SOLD -> CANCELLED）
//   3) 推进原订单状态到 CHANGED，并记录差额退款（如有）
// - 对原订单行加 FOR UPDATE，避免并发下同时退票/改签/支付回调导致状态错乱
func (s *Service) ChangeOrder(ctx context.Context, req *kitexuser.ChangeOrderReq) (*kitexuser.ChangeOrderResp, error) {
	if strings.TrimSpace(req.UserId) == "" || strings.TrimSpace(req.OrderId) == "" || strings.TrimSpace(req.NewTrainId_) == "" {
		return &kitexuser.ChangeOrderResp{BaseResp: baseResp(400, "user_id/order_id/new_train_id不能为空")}, nil
	}
	fromStation := strings.TrimSpace(req.NewDepartureStation_)
	toStation := strings.TrimSpace(req.NewArrivalStation_)
	if fromStation == "" || toStation == "" {
		return &kitexuser.ChangeOrderResp{BaseResp: baseResp(400, "new_departure_station/new_arrival_station不能为空")}, nil
	}

	newOrderID := uuid.New().String()
	var newTotal float64
	var diffRefund float64
	var oldStatus string
	var outSeats []*kitexuser.OrderSeatInfo

	err := db.MysqlDB.Transaction(func(tx *gorm.DB) error {
		// 1) 锁定原订单并校验状态（只允许 ISSUED）
		var old model2.OrderInfo
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_id = ? AND user_id = ?", req.OrderId, req.UserId).First(&old).Error; err != nil {
			return err
		}
		oldStatus = old.OrderStatus
		if old.OrderStatus != "ISSUED" {
			return fmt.Errorf("仅支持已出票订单改签")
		}
		var depStop model2.TrainStationPass
		if err := tx.Where("train_id = ? AND station_name = ?", old.TrainID, old.DepartureStation).First(&depStop).Error; err == nil {
			if depStop.DepartureTime.Valid && time.Now().After(depStop.DepartureTime.Time) {
				return fmt.Errorf("已发车，无法改签")
			}
		}

		// 2) 校验新车次与新区间站序
		var newTrain model2.TrainInfo
		if err := tx.Where("train_id = ?", req.NewTrainId_).First(&newTrain).Error; err != nil {
			return err
		}
		var newDep model2.TrainStationPass
		if err := tx.Where("train_id = ? AND station_name = ?", newTrain.ID, fromStation).First(&newDep).Error; err != nil {
			return err
		}
		var newArr model2.TrainStationPass
		if err := tx.Where("train_id = ? AND station_name = ?", newTrain.ID, toStation).First(&newArr).Error; err != nil {
			return err
		}
		if newDep.Sequence >= newArr.Sequence {
			return fmt.Errorf("上车站必须早于下车站")
		}
		if newDep.DepartureTime.Valid && time.Now().After(newDep.DepartureTime.Time) {
			return fmt.Errorf("新车次已发车，无法改签")
		}

		// 3) 读取原订单座位席别分布（need 用来保证改签后的席别数量一致）
		var rel []model2.OrderSeatRelation
		if err := tx.Where("order_id = ?", old.ID).Find(&rel).Error; err != nil {
			return err
		}
		if len(rel) == 0 {
			return fmt.Errorf("原订单缺少座位信息")
		}
		need := map[string]int{}
		for _, r := range rel {
			need[r.SeatType]++
		}

		// 4) 在同一事务内为新车次分配座位
		//    这里 lockExpire 传 time.Now()：新订单直接视为“已出票”，不需要 LOCKED 等待支付。
		allocated, sum, err := ticketapp.AllocateSeats(tx, newTrain.ID, need, newDep.Sequence, newArr.Sequence, time.Now())
		if err != nil {
			return err
		}
		newTotal = shared.Money2(sum)
		if newTotal > shared.Money2(old.TotalAmount) {
			return fmt.Errorf("新订单金额更高，暂不支持补差价改签")
		}
		diffRefund = shared.Money2(old.TotalAmount - newTotal)

		// 5) 写新订单（直接 ISSUED），并写入新占用（直接 SOLD）
		newOrder := model2.OrderInfo{
			ID:               newOrderID,
			UserID:           old.UserID,
			TrainID:          newTrain.ID,
			DepartureStation: fromStation,
			ArrivalStation:   toStation,
			FromSeq:          newDep.Sequence,
			ToSeq:            newArr.Sequence,
			TotalAmount:      newTotal,
			OrderStatus:      "ISSUED",
			PayTime:          nullTime(time.Now()),
			PayChannel:       nullString("CHANGE"),
			PayNo:            nullString("CHANGE-" + uuid.New().String()),
			RefundAmount:     0,
			RefundStatus:     "NO_REFUND",
			IdempotentKey:    "CHANGE-" + uuid.New().String(),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		if err := tx.Create(&newOrder).Error; err != nil {
			return err
		}

		occs := make([]model2.SeatSegmentOccupancy, 0, len(allocated))
		for _, seat := range allocated {
			occs = append(occs, model2.SeatSegmentOccupancy{
				TrainID:     newTrain.ID,
				SeatID:      seat.ID,
				FromStation: fromStation,
				ToStation:   toStation,
				FromSeq:     newDep.Sequence,
				ToSeq:       newArr.Sequence,
				OrderID:     newOrderID,
				Status:      "SOLD",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			})
		}
		if err := tx.CreateInBatches(&occs, 500).Error; err != nil {
			return err
		}

		seatByType := map[string][]*model2.SeatInfo{}
		for _, seat := range allocated {
			seatByType[seat.SeatType] = append(seatByType[seat.SeatType], seat)
		}

		// 6) 复用原订单乘客信息，按“乘客原席别”在新分配的座位中做一一对应
		var passengers []model2.PassengerInfo
		if err := tx.Where("order_id = ?", old.ID).Find(&passengers).Error; err != nil {
			return err
		}
		oldSeatTypeBySeatID := make(map[string]string, len(rel))
		for _, r := range rel {
			oldSeatTypeBySeatID[r.SeatID] = r.SeatType
		}
		passOut := make([]model2.PassengerInfo, 0, len(passengers))
		relOut := make([]model2.OrderSeatRelation, 0, len(passengers))
		for _, p := range passengers {
			st := strings.TrimSpace(oldSeatTypeBySeatID[p.SeatID])
			if st == "" {
				return fmt.Errorf("改签失败：无法确定乘客席别")
			}
			list := seatByType[st]
			if len(list) == 0 {
				return fmt.Errorf("改签分配座位失败")
			}
			seat := list[0]
			seatByType[st] = list[1:]
			passOut = append(passOut, model2.PassengerInfo{OrderID: newOrderID, UserID: old.UserID, RealName: p.RealName, IDCard: p.IDCard, SeatID: seat.ID, CreatedAt: time.Now(), UpdatedAt: time.Now()})
			relOut = append(relOut, model2.OrderSeatRelation{OrderID: newOrderID, SeatID: seat.ID, SeatType: seat.SeatType, SeatPrice: seat.SeatPrice, IsRefunded: "NO", CreatedAt: time.Now(), UpdatedAt: time.Now()})
			outSeats = append(outSeats, &kitexuser.OrderSeatInfo{SeatId: seat.ID, SeatType: seat.SeatType, CarriageNum: seat.CarriageNum, SeatNum: seat.SeatNum, SeatPrice: shared.Money2(seat.SeatPrice)})
		}
		if err := tx.Create(&passOut).Error; err != nil {
			return err
		}
		if err := tx.Create(&relOut).Error; err != nil {
			return err
		}

		// 7) 释放原订单已售占用（让库存回到可售状态）
		if err := tx.Model(&model2.SeatSegmentOccupancy{}).Where("order_id = ? AND status = ?", old.ID, "SOLD").Updates(map[string]any{"status": "CANCELLED"}).Error; err != nil {
			return err
		}
		updates := map[string]any{"order_status": "CHANGED"}
		if diffRefund > 0 {
			updates["refund_amount"] = diffRefund
			updates["refund_time"] = nullTime(time.Now())
			updates["refund_status"] = "REFUNDED"
		}
		if err := tx.Model(&model2.OrderInfo{}).Where("order_id = ?", old.ID).Updates(updates).Error; err != nil {
			return err
		}
		writeOrderAuditLog(tx, old.ID, "CHANGE_ORDER", req.UserId, "ISSUED", "CHANGED", map[string]any{"new_order_id": newOrderID, "refund_diff": diffRefund})
		writeOrderAuditLog(tx, newOrderID, "CHANGE_ORDER_NEW", req.UserId, "", "ISSUED", map[string]any{"old_order_id": old.ID})
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &kitexuser.ChangeOrderResp{BaseResp: baseResp(404, "订单不存在")}, nil
		}
		return &kitexuser.ChangeOrderResp{BaseResp: baseResp(400, err.Error())}, nil
	}
	return &kitexuser.ChangeOrderResp{BaseResp: baseResp(200, "success"), OldOrderStatus: oldStatus, NewOrderId_: newOrderID, NewTotalAmount_: newTotal, RefundDiffAmount: diffRefund, Seats: outSeats}, nil
}
