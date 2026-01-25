package orderapp

import (
	"context"
	"encoding/base64"
	"errors"
	"example_shop/common/db"
	model2 "example_shop/internal/model"
	"example_shop/internal/ticket_service/app/shared"
	kitexuser "example_shop/kitex_gen/user"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// GetOrder 查询订单详情用例。
func (s *Service) GetOrder(ctx context.Context, req *kitexuser.GetOrderReq) (*kitexuser.GetOrderResp, error) {
	if strings.TrimSpace(req.UserId) == "" || strings.TrimSpace(req.OrderId) == "" {
		return &kitexuser.GetOrderResp{BaseResp: baseResp(400, "user_id/order_id不能为空")}, nil
	}

	var o model2.OrderInfo
	if err := db.ReadDB().Where("order_id = ? AND user_id = ?", req.OrderId, req.UserId).First(&o).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &kitexuser.GetOrderResp{BaseResp: baseResp(404, "订单不存在")}, nil
		}
		return &kitexuser.GetOrderResp{BaseResp: baseResp(500, "查询失败: "+err.Error())}, nil
	}

	var rel []model2.OrderSeatRelation
	_ = db.ReadDB().Where("order_id = ?", o.ID).Find(&rel).Error

	seatIDs := make([]string, 0, len(rel))
	for _, r := range rel {
		seatIDs = append(seatIDs, r.SeatID)
	}
	seatDetail := map[string]model2.SeatInfo{}
	if len(seatIDs) > 0 {
		var ss []model2.SeatInfo
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

	q := db.ReadDB().Model(&model2.OrderInfo{}).Where("user_id = ?", req.UserId)
	if status != "" {
		q = q.Where("order_status = ?", status)
	}
	if req.Cursor != nil {
		if t, id, ok := decodeOrderCursor(*req.Cursor); ok {
			q = q.Where("(created_at < ?) OR (created_at = ? AND order_id < ?)", t, t, id)
		}
	}
	q = q.Order("created_at DESC, order_id DESC").Limit(limit)

	var orders []model2.OrderInfo
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

