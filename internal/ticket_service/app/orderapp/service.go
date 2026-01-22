// Package orderapp 承载“订单域”的应用层用例（下单/支付/退改/查询等）。
//
// 注意：订单域会依赖票务域的“座席分配”边界能力（ticketapp.AllocateSeats），
// 从而保证依赖方向清晰：order -> ticket，而不是把票务细节散落在订单实现里。
package orderapp

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"example_shop/common/db"
	"example_shop/internal/ticket_service/app/shared"
	"example_shop/internal/ticket_service/app/ticketapp"
	"example_shop/internal/ticket_service/model"
	kitexuser "example_shop/kitex_gen/user"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service struct{}

// New 创建订单域应用服务。
func New() *Service {
	return &Service{}
}

// CreateOrder 下单用例：锁座 + 创建订单（待支付）+ 返回支付截止时间与座位信息。
func (s *Service) CreateOrder(ctx context.Context, req *kitexuser.CreateOrderReq) (*kitexuser.CreateOrderResp, error) {
	if strings.TrimSpace(req.UserId) == "" || strings.TrimSpace(req.TrainId) == "" {
		return &kitexuser.CreateOrderResp{BaseResp: baseResp(400, "user_id/train_id不能为空")}, nil
	}
	if strings.TrimSpace(req.DepartureStation) == "" || strings.TrimSpace(req.ArrivalStation) == "" {
		return &kitexuser.CreateOrderResp{BaseResp: baseResp(400, "departure_station/arrival_station不能为空")}, nil
	}
	if len(req.Passengers) == 0 {
		return &kitexuser.CreateOrderResp{BaseResp: baseResp(400, "passengers不能为空")}, nil
	}
	if len(req.Passengers) > 5 {
		return &kitexuser.CreateOrderResp{BaseResp: baseResp(400, "单次最多支持5名乘客")}, nil
	}

	var t model.TrainInfo
	if err := db.MysqlDB.Where("train_id = ?", req.TrainId).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &kitexuser.CreateOrderResp{BaseResp: baseResp(404, "车次不存在")}, nil
		}
		return &kitexuser.CreateOrderResp{BaseResp: baseResp(500, "查询车次失败: "+err.Error())}, nil
	}
	if time.Now().After(t.DepartureTime) {
		return &kitexuser.CreateOrderResp{BaseResp: baseResp(400, "车次已发车，无法下单")}, nil
	}

	fromStation := strings.TrimSpace(req.DepartureStation)
	toStation := strings.TrimSpace(req.ArrivalStation)
	var depStop model.TrainStationPass
	if err := db.MysqlDB.Where("train_id = ? AND station_name = ?", t.ID, fromStation).First(&depStop).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &kitexuser.CreateOrderResp{BaseResp: baseResp(404, "上车站不在该车次途经站列表")}, nil
		}
		return &kitexuser.CreateOrderResp{BaseResp: baseResp(500, "查询上车站失败: "+err.Error())}, nil
	}
	var arrStop model.TrainStationPass
	if err := db.MysqlDB.Where("train_id = ? AND station_name = ?", t.ID, toStation).First(&arrStop).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &kitexuser.CreateOrderResp{BaseResp: baseResp(404, "下车站不在该车次途经站列表")}, nil
		}
		return &kitexuser.CreateOrderResp{BaseResp: baseResp(500, "查询下车站失败: "+err.Error())}, nil
	}
	if depStop.Sequence >= arrStop.Sequence {
		return &kitexuser.CreateOrderResp{BaseResp: baseResp(400, "上车站必须早于下车站")}, nil
	}

	need := make(map[string]int)
	for _, p := range req.Passengers {
		if p == nil {
			return &kitexuser.CreateOrderResp{BaseResp: baseResp(400, "passengers包含空元素")}, nil
		}
		if strings.TrimSpace(p.RealName) == "" || strings.TrimSpace(p.IdCard) == "" {
			return &kitexuser.CreateOrderResp{BaseResp: baseResp(400, "乘客 real_name/id_card 不能为空")}, nil
		}
		if !seatTypeAllowed(p.SeatType) {
			return &kitexuser.CreateOrderResp{BaseResp: baseResp(400, "座位类型无效")}, nil
		}
		need[p.SeatType]++
	}

	orderID := uuid.New().String()
	payDeadline := time.Now().Add(15 * time.Minute)
	idemKey := buildOrderIdempotentKey(req)

	var allocated []*model.SeatInfo
	var total float64

	err := db.MysqlDB.Transaction(func(tx *gorm.DB) error {
		allocated = allocated[:0]
		total = 0
		got, sum, err := ticketapp.AllocateSeats(tx, t.ID, need, depStop.Sequence, arrStop.Sequence, payDeadline)
		if err != nil {
			return err
		}
		allocated = append(allocated, got...)
		total = sum

		order := model.OrderInfo{
			ID:               orderID,
			UserID:           req.UserId,
			TrainID:          t.ID,
			DepartureStation: fromStation,
			ArrivalStation:   toStation,
			FromSeq:          depStop.Sequence,
			ToSeq:            arrStop.Sequence,
			TotalAmount:      0,
			OrderStatus:      "PENDING_PAY",
			PayDeadline:      nullTime(payDeadline),
			RefundAmount:     0,
			RefundStatus:     "NO_REFUND",
			IdempotentKey:    idemKey,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		if err := tx.Create(&order).Error; err != nil {
			if isDuplicateKeyError(err) {
				return &idempotentHitError{Key: idemKey}
			}
			return err
		}
		writeOrderAuditLog(tx, orderID, "CREATE_ORDER", req.UserId, "", "PENDING_PAY", map[string]any{"train_id": t.ID, "from": fromStation, "to": toStation})

		occs := make([]model.SeatSegmentOccupancy, 0, len(allocated))
		for _, seat := range allocated {
			occs = append(occs, model.SeatSegmentOccupancy{
				TrainID:        t.ID,
				SeatID:         seat.ID,
				FromStation:    fromStation,
				ToStation:      toStation,
				FromSeq:        depStop.Sequence,
				ToSeq:          arrStop.Sequence,
				OrderID:        orderID,
				Status:         "LOCKED",
				LockExpireTime: nullTime(payDeadline),
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			})
		}
		if err := tx.CreateInBatches(&occs, 500).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.OrderInfo{}).Where("order_id = ?", orderID).Updates(map[string]any{"total_amount": shared.Money2(total)}).Error; err != nil {
			return err
		}

		rel := make([]model.OrderSeatRelation, 0, len(allocated))
		pass := make([]model.PassengerInfo, 0, len(req.Passengers))
		seatByType := map[string][]*model.SeatInfo{}
		for _, seat := range allocated {
			seatByType[seat.SeatType] = append(seatByType[seat.SeatType], seat)
		}
		for _, p := range req.Passengers {
			list := seatByType[p.SeatType]
			seat := list[0]
			seatByType[p.SeatType] = list[1:]
			pass = append(pass, model.PassengerInfo{
				OrderID:   orderID,
				UserID:    req.UserId,
				RealName:  p.RealName,
				IDCard:    p.IdCard,
				SeatID:    seat.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
			rel = append(rel, model.OrderSeatRelation{
				OrderID:    orderID,
				SeatID:     seat.ID,
				SeatType:   seat.SeatType,
				SeatPrice:  seat.SeatPrice,
				IsRefunded: "NO",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			})
		}
		if err := tx.Create(&pass).Error; err != nil {
			return err
		}
		if err := tx.Create(&rel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if hit, ok := err.(*idempotentHitError); ok {
			return loadIdempotentOrder(hit.Key)
		}
		return &kitexuser.CreateOrderResp{BaseResp: baseResp(400, err.Error())}, nil
	}

	seats := make([]*kitexuser.OrderSeatInfo, 0, len(allocated))
	for _, seat := range allocated {
		seats = append(seats, &kitexuser.OrderSeatInfo{SeatId: seat.ID, SeatType: seat.SeatType, CarriageNum: seat.CarriageNum, SeatNum: seat.SeatNum, SeatPrice: shared.Money2(seat.SeatPrice)})
	}
	return &kitexuser.CreateOrderResp{BaseResp: baseResp(200, "success"), OrderId: orderID, PayDeadlineUnix: payDeadline.Unix(), Seats: seats}, nil
}

// PayOrder 支付发起用例：将订单状态推进到 PAYING（简化版）。
func (s *Service) PayOrder(ctx context.Context, req *kitexuser.PayOrderReq) (*kitexuser.PayOrderResp, error) {
	if strings.TrimSpace(req.UserId) == "" || strings.TrimSpace(req.OrderId) == "" {
		return &kitexuser.PayOrderResp{BaseResp: baseResp(400, "user_id/order_id不能为空")}, nil
	}

	var payNo string
	var outStatus string
	err := db.MysqlDB.Transaction(func(tx *gorm.DB) error {
		var o model.OrderInfo
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_id = ? AND user_id = ?", req.OrderId, req.UserId).First(&o).Error; err != nil {
			return err
		}
		if o.OrderStatus == "ISSUED" {
			outStatus = "ISSUED"
			payNo = o.PayNo.String
			return nil
		}
		if o.OrderStatus != "PENDING_PAY" && o.OrderStatus != "PAYING" {
			return fmt.Errorf("订单状态不允许支付")
		}
		if o.PayDeadline.Valid && time.Now().After(o.PayDeadline.Time) {
			return fmt.Errorf("支付超时")
		}

		if o.PayNo.Valid {
			payNo = o.PayNo.String
		} else {
			payNo = "MOCKPAY-" + uuid.New().String()
		}
		updates := map[string]interface{}{
			"order_status": "PAYING",
			"pay_channel":  nullString(req.PayChannel),
			"pay_no":       nullString(payNo),
		}
		if err := tx.Model(&model.OrderInfo{}).Where("order_id = ?", o.ID).Updates(updates).Error; err != nil {
			return err
		}
		writeOrderAuditLog(tx, o.ID, "PAY_ORDER", req.UserId, o.OrderStatus, "PAYING", map[string]any{"pay_channel": req.PayChannel, "pay_no": payNo})
		outStatus = "PAYING"
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &kitexuser.PayOrderResp{BaseResp: baseResp(404, "订单不存在")}, nil
		}
		return &kitexuser.PayOrderResp{BaseResp: baseResp(400, err.Error())}, nil
	}
	resp := &kitexuser.PayOrderResp{BaseResp: baseResp(200, "success"), OrderStatus: outStatus, PayStatus: "PAYING"}
	if payNo != "" {
		resp.PayNo = &payNo
	}
	if outStatus == "ISSUED" {
		resp.PayStatus = "PAID"
	}
	return resp, nil
}

// ConfirmPay 支付回调确认用例：锁座占用变 SOLD，订单状态推进到 ISSUED。
func (s *Service) ConfirmPay(ctx context.Context, req *kitexuser.ConfirmPayReq) (*kitexuser.ConfirmPayResp, error) {
	if strings.TrimSpace(req.UserId) == "" || strings.TrimSpace(req.OrderId) == "" || strings.TrimSpace(req.PayNo) == "" {
		return &kitexuser.ConfirmPayResp{BaseResp: baseResp(400, "user_id/order_id/pay_no不能为空")}, nil
	}

	var outStatus string
	err := db.MysqlDB.Transaction(func(tx *gorm.DB) error {
		var o model.OrderInfo
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_id = ? AND user_id = ?", req.OrderId, req.UserId).First(&o).Error; err != nil {
			return err
		}
		if o.OrderStatus == "ISSUED" {
			outStatus = "ISSUED"
			return nil
		}
		if o.OrderStatus != "PAYING" {
			return fmt.Errorf("订单状态不允许确认支付")
		}
		if !o.PayNo.Valid || o.PayNo.String != req.PayNo {
			return fmt.Errorf("pay_no不匹配")
		}
		if o.PayDeadline.Valid && time.Now().After(o.PayDeadline.Time) {
			return fmt.Errorf("支付已超时")
		}
		third := strings.TrimSpace(req.GetThirdPartyStatus())
		if third != "" && strings.ToUpper(third) != "SUCCESS" {
			return fmt.Errorf("第三方支付未成功")
		}

		var occs []model.SeatSegmentOccupancy
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_id = ? AND status = ?", o.ID, "LOCKED").Find(&occs).Error; err != nil {
			return err
		}
		if len(occs) == 0 {
			return fmt.Errorf("未找到锁定座位")
		}
		now := time.Now()
		for _, oc := range occs {
			if oc.LockExpireTime.Valid && now.After(oc.LockExpireTime.Time) {
				return fmt.Errorf("锁座已过期")
			}
		}
		if err := tx.Model(&model.SeatSegmentOccupancy{}).
			Where("order_id = ? AND status = ?", o.ID, "LOCKED").
			Updates(map[string]interface{}{"status": "SOLD", "lock_expire_time": sql.NullTime{Valid: false}}).Error; err != nil {
			return err
		}

		payTime := time.Now()
		if err := tx.Model(&model.OrderInfo{}).Where("order_id = ?", o.ID).Updates(map[string]any{
			"order_status": "ISSUED",
			"pay_time":     nullTime(payTime),
		}).Error; err != nil {
			return err
		}
		writeOrderAuditLog(tx, o.ID, "PAY_CALLBACK", req.UserId, "PAYING", "ISSUED", map[string]any{"pay_no": req.PayNo, "third_party_status": third})
		outStatus = "ISSUED"
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &kitexuser.ConfirmPayResp{BaseResp: baseResp(404, "订单不存在")}, nil
		}
		return &kitexuser.ConfirmPayResp{BaseResp: baseResp(400, err.Error())}, nil
	}
	return &kitexuser.ConfirmPayResp{BaseResp: baseResp(200, "success"), OrderStatus: outStatus}, nil
}

// CancelOrder 取消订单用例：释放 LOCKED 的区间占用并将订单置为 CANCELLED。
func (s *Service) CancelOrder(ctx context.Context, req *kitexuser.CancelOrderReq) (*kitexuser.BaseResp, error) {
	if strings.TrimSpace(req.UserId) == "" || strings.TrimSpace(req.OrderId) == "" {
		return baseResp(400, "user_id/order_id不能为空"), nil
	}

	err := db.MysqlDB.Transaction(func(tx *gorm.DB) error {
		var o model.OrderInfo
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_id = ? AND user_id = ?", req.OrderId, req.UserId).First(&o).Error; err != nil {
			return err
		}
		if o.OrderStatus != "PENDING_PAY" && o.OrderStatus != "PAYING" {
			return fmt.Errorf("当前订单状态不允许取消")
		}
		if err := tx.Model(&model.SeatSegmentOccupancy{}).
			Where("order_id = ? AND status = ?", o.ID, "LOCKED").
			Updates(map[string]interface{}{"status": "CANCELLED", "lock_expire_time": sql.NullTime{Valid: false}}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.OrderInfo{}).Where("order_id = ?", o.ID).Updates(map[string]interface{}{"order_status": "CANCELLED"}).Error; err != nil {
			return err
		}
		writeOrderAuditLog(tx, o.ID, "CANCEL_ORDER", req.UserId, o.OrderStatus, "CANCELLED", nil)
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return baseResp(404, "订单不存在"), nil
		}
		return baseResp(400, err.Error()), nil
	}
	return baseResp(200, "success"), nil
}

// RefundOrder 退票用例：将 SOLD 的区间占用取消，并根据规则计算退款金额。
func (s *Service) RefundOrder(ctx context.Context, req *kitexuser.RefundOrderReq) (*kitexuser.RefundOrderResp, error) {
	if strings.TrimSpace(req.UserId) == "" || strings.TrimSpace(req.OrderId) == "" {
		return &kitexuser.RefundOrderResp{BaseResp: baseResp(400, "user_id/order_id不能为空")}, nil
	}

	var refundAmount float64
	var outStatus string
	err := db.MysqlDB.Transaction(func(tx *gorm.DB) error {
		var o model.OrderInfo
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_id = ? AND user_id = ?", req.OrderId, req.UserId).First(&o).Error; err != nil {
			return err
		}
		if o.OrderStatus != "ISSUED" {
			return fmt.Errorf("订单状态不允许退票")
		}
		var depStop model.TrainStationPass
		if err := tx.Where("train_id = ? AND station_name = ?", o.TrainID, o.DepartureStation).First(&depStop).Error; err == nil {
			if depStop.DepartureTime.Valid {
				if time.Now().After(depStop.DepartureTime.Time) {
					return fmt.Errorf("已发车，无法退票")
				}
				hoursBefore := depStop.DepartureTime.Time.Sub(time.Now()).Hours()
				rule := loadRefundFeeRule("")
				rate, ok := refundFeeRate(rule, hoursBefore)
				if !ok {
					return fmt.Errorf("距离发车过近，不允许退票")
				}
				refundAmount = shared.Money2(o.TotalAmount * (1 - rate))
			} else {
				rule := loadRefundFeeRule("")
				rate, _ := refundFeeRate(rule, 24)
				refundAmount = shared.Money2(o.TotalAmount * (1 - rate))
			}
		} else {
			refundAmount = shared.Money2(o.TotalAmount * 0.90)
		}
		if refundAmount < 0 {
			refundAmount = 0
		}

		if err := tx.Model(&model.SeatSegmentOccupancy{}).
			Where("order_id = ? AND status = ?", o.ID, "SOLD").
			Updates(map[string]any{"status": "CANCELLED"}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.OrderInfo{}).Where("order_id = ?", o.ID).Updates(map[string]any{
			"order_status":  "REFUNDED",
			"refund_amount": refundAmount,
			"refund_time":   nullTime(time.Now()),
			"refund_status": "REFUNDED",
		}).Error; err != nil {
			return err
		}
		writeOrderAuditLog(tx, o.ID, "REFUND_ORDER", req.UserId, "ISSUED", "REFUNDED", map[string]any{"reason": req.GetReason(), "refund_amount": refundAmount})
		outStatus = "REFUNDED"
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &kitexuser.RefundOrderResp{BaseResp: baseResp(404, "订单不存在")}, nil
		}
		return &kitexuser.RefundOrderResp{BaseResp: baseResp(400, err.Error())}, nil
	}
	return &kitexuser.RefundOrderResp{BaseResp: baseResp(200, "success"), RefundAmount: shared.Money2(refundAmount), RefundStatus: "REFUNDED", OrderStatus: outStatus}, nil
}

// ChangeOrder 改签用例：创建一张新订单并将原订单置为 CHANGED（简化版）。
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
		var old model.OrderInfo
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_id = ? AND user_id = ?", req.OrderId, req.UserId).First(&old).Error; err != nil {
			return err
		}
		oldStatus = old.OrderStatus
		if old.OrderStatus != "ISSUED" {
			return fmt.Errorf("仅支持已出票订单改签")
		}
		var depStop model.TrainStationPass
		if err := tx.Where("train_id = ? AND station_name = ?", old.TrainID, old.DepartureStation).First(&depStop).Error; err == nil {
			if depStop.DepartureTime.Valid && time.Now().After(depStop.DepartureTime.Time) {
				return fmt.Errorf("已发车，无法改签")
			}
		}

		var newTrain model.TrainInfo
		if err := tx.Where("train_id = ?", req.NewTrainId_).First(&newTrain).Error; err != nil {
			return err
		}
		var newDep model.TrainStationPass
		if err := tx.Where("train_id = ? AND station_name = ?", newTrain.ID, fromStation).First(&newDep).Error; err != nil {
			return err
		}
		var newArr model.TrainStationPass
		if err := tx.Where("train_id = ? AND station_name = ?", newTrain.ID, toStation).First(&newArr).Error; err != nil {
			return err
		}
		if newDep.Sequence >= newArr.Sequence {
			return fmt.Errorf("上车站必须早于下车站")
		}
		if newDep.DepartureTime.Valid && time.Now().After(newDep.DepartureTime.Time) {
			return fmt.Errorf("新车次已发车，无法改签")
		}

		var rel []model.OrderSeatRelation
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

		allocated, sum, err := ticketapp.AllocateSeats(tx, newTrain.ID, need, newDep.Sequence, newArr.Sequence, time.Now())
		if err != nil {
			return err
		}
		newTotal = shared.Money2(sum)
		if newTotal > shared.Money2(old.TotalAmount) {
			return fmt.Errorf("新订单金额更高，暂不支持补差价改签")
		}
		diffRefund = shared.Money2(old.TotalAmount - newTotal)

		newOrder := model.OrderInfo{
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

		occs := make([]model.SeatSegmentOccupancy, 0, len(allocated))
		for _, seat := range allocated {
			occs = append(occs, model.SeatSegmentOccupancy{
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

		seatByType := map[string][]*model.SeatInfo{}
		for _, seat := range allocated {
			seatByType[seat.SeatType] = append(seatByType[seat.SeatType], seat)
		}

		var passengers []model.PassengerInfo
		if err := tx.Where("order_id = ?", old.ID).Find(&passengers).Error; err != nil {
			return err
		}
		oldSeatTypeBySeatID := make(map[string]string, len(rel))
		for _, r := range rel {
			oldSeatTypeBySeatID[r.SeatID] = r.SeatType
		}
		passOut := make([]model.PassengerInfo, 0, len(passengers))
		relOut := make([]model.OrderSeatRelation, 0, len(passengers))
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
			passOut = append(passOut, model.PassengerInfo{OrderID: newOrderID, UserID: old.UserID, RealName: p.RealName, IDCard: p.IDCard, SeatID: seat.ID, CreatedAt: time.Now(), UpdatedAt: time.Now()})
			relOut = append(relOut, model.OrderSeatRelation{OrderID: newOrderID, SeatID: seat.ID, SeatType: seat.SeatType, SeatPrice: seat.SeatPrice, IsRefunded: "NO", CreatedAt: time.Now(), UpdatedAt: time.Now()})
			outSeats = append(outSeats, &kitexuser.OrderSeatInfo{SeatId: seat.ID, SeatType: seat.SeatType, CarriageNum: seat.CarriageNum, SeatNum: seat.SeatNum, SeatPrice: shared.Money2(seat.SeatPrice)})
		}
		if err := tx.Create(&passOut).Error; err != nil {
			return err
		}
		if err := tx.Create(&relOut).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.SeatSegmentOccupancy{}).Where("order_id = ? AND status = ?", old.ID, "SOLD").Updates(map[string]any{"status": "CANCELLED"}).Error; err != nil {
			return err
		}
		updates := map[string]any{"order_status": "CHANGED"}
		if diffRefund > 0 {
			updates["refund_amount"] = diffRefund
			updates["refund_time"] = nullTime(time.Now())
			updates["refund_status"] = "REFUNDED"
		}
		if err := tx.Model(&model.OrderInfo{}).Where("order_id = ?", old.ID).Updates(updates).Error; err != nil {
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

// GetOrder 查询订单详情用例。
func (s *Service) GetOrder(ctx context.Context, req *kitexuser.GetOrderReq) (*kitexuser.GetOrderResp, error) {
	if strings.TrimSpace(req.UserId) == "" || strings.TrimSpace(req.OrderId) == "" {
		return &kitexuser.GetOrderResp{BaseResp: baseResp(400, "user_id/order_id不能为空")}, nil
	}

	var o model.OrderInfo
	if err := db.ReadDB().Where("order_id = ? AND user_id = ?", req.OrderId, req.UserId).First(&o).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &kitexuser.GetOrderResp{BaseResp: baseResp(404, "订单不存在")}, nil
		}
		return &kitexuser.GetOrderResp{BaseResp: baseResp(500, "查询失败: "+err.Error())}, nil
	}

	var rel []model.OrderSeatRelation
	_ = db.ReadDB().Where("order_id = ?", o.ID).Find(&rel).Error

	seatIDs := make([]string, 0, len(rel))
	for _, r := range rel {
		seatIDs = append(seatIDs, r.SeatID)
	}
	seatDetail := map[string]model.SeatInfo{}
	if len(seatIDs) > 0 {
		var ss []model.SeatInfo
		_ = db.ReadDB().Where("seat_id IN ?", seatIDs).Find(&ss).Error
		for _, seat := range ss {
			seatDetail[seat.ID] = seat
		}
	}

	seats := make([]*kitexuser.OrderSeatInfo, 0, len(rel))
	for _, r := range rel {
		d, ok := seatDetail[r.SeatID]
		if ok {
			seats = append(seats, &kitexuser.OrderSeatInfo{SeatId: r.SeatID, SeatType: r.SeatType, CarriageNum: d.CarriageNum, SeatNum: d.SeatNum, SeatPrice: shared.Money2(r.SeatPrice)})
			continue
		}
		seats = append(seats, &kitexuser.OrderSeatInfo{SeatId: r.SeatID, SeatType: r.SeatType, SeatPrice: shared.Money2(r.SeatPrice)})
	}

	payDeadline := int64(0)
	if o.PayDeadline.Valid {
		payDeadline = o.PayDeadline.Time.Unix()
	}
	var payTime *int64
	if o.PayTime.Valid {
		v := o.PayTime.Time.Unix()
		payTime = &v
	}
	var payChannel *string
	if o.PayChannel.Valid {
		v := o.PayChannel.String
		payChannel = &v
	}
	var payNo *string
	if o.PayNo.Valid {
		v := o.PayNo.String
		payNo = &v
	}

	info := &kitexuser.OrderInfo{
		OrderId:          o.ID,
		UserId:           o.UserID,
		TrainId:          o.TrainID,
		DepartureStation: o.DepartureStation,
		ArrivalStation:   o.ArrivalStation,
		FromSeq:          int32(o.FromSeq),
		ToSeq:            int32(o.ToSeq),
		TotalAmount:      shared.Money2(o.TotalAmount),
		OrderStatus:      o.OrderStatus,
		PayDeadlineUnix:  payDeadline,
		PayTimeUnix:      payTime,
		PayChannel:       payChannel,
		PayNo:            payNo,
		CreatedAtUnix:    o.CreatedAt.Unix(),
	}
	return &kitexuser.GetOrderResp{BaseResp: baseResp(200, "success"), Order: info, Seats: seats}, nil
}

// ListOrders 查询订单列表用例：支持按状态过滤与游标分页。
func (s *Service) ListOrders(ctx context.Context, req *kitexuser.ListOrdersReq) (*kitexuser.ListOrdersResp, error) {
	if strings.TrimSpace(req.UserId) == "" {
		return &kitexuser.ListOrdersResp{BaseResp: baseResp(400, "user_id不能为空")}, nil
	}
	limit := int(req.Limit)
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	status := ""
	if req.Status != nil {
		status = strings.TrimSpace(*req.Status)
	}

	q := db.ReadDB().Model(&model.OrderInfo{}).Where("user_id = ?", req.UserId)
	if status != "" {
		q = q.Where("order_status = ?", status)
	}
	if req.Cursor != nil {
		if t, id, ok := decodeOrderCursor(*req.Cursor); ok {
			q = q.Where("(created_at < ?) OR (created_at = ? AND order_id < ?)", t, t, id)
		}
	}
	q = q.Order("created_at DESC, order_id DESC").Limit(limit)

	var orders []model.OrderInfo
	if err := q.Find(&orders).Error; err != nil {
		return &kitexuser.ListOrdersResp{BaseResp: baseResp(500, "查询失败: "+err.Error()), Orders: []*kitexuser.OrderInfo{}}, nil
	}

	items := make([]*kitexuser.OrderInfo, 0, len(orders))
	for _, o := range orders {
		payDeadline := int64(0)
		if o.PayDeadline.Valid {
			payDeadline = o.PayDeadline.Time.Unix()
		}
		var payTime *int64
		if o.PayTime.Valid {
			v := o.PayTime.Time.Unix()
			payTime = &v
		}
		var payChannel *string
		if o.PayChannel.Valid {
			v := o.PayChannel.String
			payChannel = &v
		}
		var payNo *string
		if o.PayNo.Valid {
			v := o.PayNo.String
			payNo = &v
		}
		items = append(items, &kitexuser.OrderInfo{
			OrderId:          o.ID,
			UserId:           o.UserID,
			TrainId:          o.TrainID,
			DepartureStation: o.DepartureStation,
			ArrivalStation:   o.ArrivalStation,
			FromSeq:          int32(o.FromSeq),
			ToSeq:            int32(o.ToSeq),
			TotalAmount:      shared.Money2(o.TotalAmount),
			OrderStatus:      o.OrderStatus,
			PayDeadlineUnix:  payDeadline,
			PayTimeUnix:      payTime,
			PayChannel:       payChannel,
			PayNo:            payNo,
			CreatedAtUnix:    o.CreatedAt.Unix(),
		})
	}

	var next *string
	if len(orders) == limit {
		last := orders[len(orders)-1]
		c := encodeOrderCursor(last.CreatedAt, last.ID)
		next = &c
	}
	return &kitexuser.ListOrdersResp{BaseResp: baseResp(200, "success"), Orders: items, NextCursor: next}, nil
}

type idempotentHitError struct {
	Key string
}

func (e *idempotentHitError) Error() string {
	return "idempotent hit"
}

func baseResp(code int32, msg string) *kitexuser.BaseResp {
	return &kitexuser.BaseResp{Code: code, Msg: msg}
}

func nullTime(t time.Time) sql.NullTime {
	return sql.NullTime{Time: t, Valid: true}
}

func nullString(s string) sql.NullString {
	if strings.TrimSpace(s) == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func seatTypeAllowed(seatType string) bool {
	switch strings.TrimSpace(seatType) {
	case "硬座", "二等座", "一等座", "商务座", "硬卧", "软卧":
		return true
	default:
		return false
	}
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "Duplicate entry") || strings.Contains(msg, "duplicate")
}

type refundTier struct {
	HoursBefore int     `json:"hours_before"`
	FeeRate     float64 `json:"fee_rate"`
}

type refundFeeRule struct {
	Tiers      []refundTier `json:"tiers"`
	DefaultFee float64      `json:"default_fee"`
	NoRefundLT float64      `json:"no_refund_lt_hours"`
}

func loadRefundFeeRule(trainType string) refundFeeRule {
	var cfg model.TicketRuleConfig
	err := db.ReadDB().Where("rule_type = ? AND status = ?", "REFUND_FEE", "ENABLED").Order("version DESC").First(&cfg).Error
	if err == nil && cfg.RuleParam != nil {
		var rule refundFeeRule
		if jsonErr := json.Unmarshal([]byte(*cfg.RuleParam), &rule); jsonErr == nil && len(rule.Tiers) > 0 {
			sort.Slice(rule.Tiers, func(i, j int) bool { return rule.Tiers[i].HoursBefore > rule.Tiers[j].HoursBefore })
			return rule
		}
	}
	return refundFeeRule{
		Tiers: []refundTier{
			{HoursBefore: 48, FeeRate: 0.05},
			{HoursBefore: 24, FeeRate: 0.10},
			{HoursBefore: 2, FeeRate: 0.20},
			{HoursBefore: 0, FeeRate: 0.50},
		},
		DefaultFee: 0.20,
		NoRefundLT: 0,
	}
}

func refundFeeRate(rule refundFeeRule, hoursBefore float64) (float64, bool) {
	if rule.NoRefundLT > 0 && hoursBefore < rule.NoRefundLT {
		return 0, false
	}
	for _, t := range rule.Tiers {
		if hoursBefore >= float64(t.HoursBefore) {
			return t.FeeRate, true
		}
	}
	return rule.DefaultFee, true
}

func buildOrderIdempotentKey(req *kitexuser.CreateOrderReq) string {
	type p struct {
		IDCard   string `json:"id_card"`
		SeatType string `json:"seat_type"`
	}
	ps := make([]p, 0, len(req.Passengers))
	for _, it := range req.Passengers {
		if it == nil {
			continue
		}
		ps = append(ps, p{IDCard: strings.TrimSpace(it.IdCard), SeatType: strings.TrimSpace(it.SeatType)})
	}
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].SeatType == ps[j].SeatType {
			return ps[i].IDCard < ps[j].IDCard
		}
		return ps[i].SeatType < ps[j].SeatType
	})
	b, _ := json.Marshal(ps)
	src := req.UserId + "|" + req.TrainId + "|" + strings.TrimSpace(req.DepartureStation) + "|" + strings.TrimSpace(req.ArrivalStation) + "|" + string(b)
	h := md5.Sum([]byte(src))
	return fmt.Sprintf("%x", h)
}

func encodeOrderCursor(createdAt time.Time, orderID string) string {
	str := fmt.Sprintf("%d|%s", createdAt.UnixNano(), orderID)
	return base64.RawURLEncoding.EncodeToString([]byte(str))
}

func decodeOrderCursor(cursor string) (time.Time, string, bool) {
	if strings.TrimSpace(cursor) == "" {
		return time.Time{}, "", false
	}
	raw, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return time.Time{}, "", false
	}
	parts := strings.SplitN(string(raw), "|", 2)
	if len(parts) != 2 {
		return time.Time{}, "", false
	}
	ns, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return time.Time{}, "", false
	}
	return time.Unix(0, ns), parts[1], true
}

func writeOrderAuditLog(tx *gorm.DB, orderID, operateType, userID, before, after string, detail any) {
	var beforeNS, afterNS sql.NullString
	if strings.TrimSpace(before) != "" {
		beforeNS = sql.NullString{String: before, Valid: true}
	}
	if strings.TrimSpace(after) != "" {
		afterNS = sql.NullString{String: after, Valid: true}
	}
	var detailJSON *model.JSON
	if detail != nil {
		if b, err := json.Marshal(detail); err == nil {
			if j, err2 := model.ToJSON(json.RawMessage(b)); err2 == nil {
				detailJSON = &j
			}
		}
	}
	logRow := model.OrderAuditLog{
		OrderID:       orderID,
		OperateType:   operateType,
		OperateUser:   userID,
		BeforeStatus:  beforeNS,
		AfterStatus:   afterNS,
		OperateDetail: detailJSON,
		TraceID:       uuid.New().String(),
		CreatedAt:     time.Now(),
	}
	_ = tx.Create(&logRow).Error
}

func loadIdempotentOrder(idemKey string) (*kitexuser.CreateOrderResp, error) {
	var o model.OrderInfo
	if err := db.ReadDB().Where("idempotent_key = ?", idemKey).First(&o).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &kitexuser.CreateOrderResp{BaseResp: baseResp(409, "重复请求但未找到原订单")}, nil
		}
		return &kitexuser.CreateOrderResp{BaseResp: baseResp(500, "查询原订单失败: "+err.Error())}, nil
	}
	var rel []model.OrderSeatRelation
	_ = db.ReadDB().Where("order_id = ?", o.ID).Find(&rel).Error
	ids := make([]string, 0, len(rel))
	for _, r := range rel {
		ids = append(ids, r.SeatID)
	}
	seatDetail := map[string]model.SeatInfo{}
	if len(ids) > 0 {
		var ss []model.SeatInfo
		_ = db.ReadDB().Where("seat_id IN ?", ids).Find(&ss).Error
		for _, seat := range ss {
			seatDetail[seat.ID] = seat
		}
	}
	seats := make([]*kitexuser.OrderSeatInfo, 0, len(rel))
	for _, r := range rel {
		d, ok := seatDetail[r.SeatID]
		if ok {
			seats = append(seats, &kitexuser.OrderSeatInfo{SeatId: r.SeatID, SeatType: r.SeatType, CarriageNum: d.CarriageNum, SeatNum: d.SeatNum, SeatPrice: shared.Money2(r.SeatPrice)})
			continue
		}
		seats = append(seats, &kitexuser.OrderSeatInfo{SeatId: r.SeatID, SeatType: r.SeatType, SeatPrice: shared.Money2(r.SeatPrice)})
	}
	deadline := int64(0)
	if o.PayDeadline.Valid {
		deadline = o.PayDeadline.Time.Unix()
	}
	return &kitexuser.CreateOrderResp{BaseResp: baseResp(200, "success"), OrderId: o.ID, PayDeadlineUnix: deadline, Seats: seats}, nil
}
