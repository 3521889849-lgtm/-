package orderapi

import (
	"context"
	"example_shop/internal/ticket_service/app/orderapp"
	"example_shop/kitex_gen/common"
	"example_shop/kitex_gen/orderapi"
	legacy "example_shop/kitex_gen/user"
)

type OrderServiceImpl struct {
	orderApp *orderapp.Service
}

func NewOrderServiceImpl() *OrderServiceImpl {
	return &OrderServiceImpl{orderApp: orderapp.New()}
}

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req *orderapi.CreateOrderReq) (*orderapi.CreateOrderResp, error) {
	resp, err := s.orderApp.CreateOrder(ctx, toLegacyCreateOrderReq(req))
	if err != nil {
		return nil, err
	}
	return fromLegacyCreateOrderResp(resp), nil
}

func (s *OrderServiceImpl) PayOrder(ctx context.Context, req *orderapi.PayOrderReq) (*orderapi.PayOrderResp, error) {
	resp, err := s.orderApp.PayOrder(ctx, toLegacyPayOrderReq(req))
	if err != nil {
		return nil, err
	}
	return fromLegacyPayOrderResp(resp), nil
}

func (s *OrderServiceImpl) ConfirmPay(ctx context.Context, req *orderapi.ConfirmPayReq) (*orderapi.ConfirmPayResp, error) {
	resp, err := s.orderApp.ConfirmPay(ctx, toLegacyConfirmPayReq(req))
	if err != nil {
		return nil, err
	}
	return fromLegacyConfirmPayResp(resp), nil
}

func (s *OrderServiceImpl) ChangeOrder(ctx context.Context, req *orderapi.ChangeOrderReq) (*orderapi.ChangeOrderResp, error) {
	resp, err := s.orderApp.ChangeOrder(ctx, toLegacyChangeOrderReq(req))
	if err != nil {
		return nil, err
	}
	return fromLegacyChangeOrderResp(resp), nil
}

func (s *OrderServiceImpl) CancelOrder(ctx context.Context, req *orderapi.CancelOrderReq) (*common.BaseResp, error) {
	resp, err := s.orderApp.CancelOrder(ctx, toLegacyCancelOrderReq(req))
	if err != nil {
		return nil, err
	}
	return toCommon(resp), nil
}

func (s *OrderServiceImpl) RefundOrder(ctx context.Context, req *orderapi.RefundOrderReq) (*orderapi.RefundOrderResp, error) {
	resp, err := s.orderApp.RefundOrder(ctx, toLegacyRefundOrderReq(req))
	if err != nil {
		return nil, err
	}
	return fromLegacyRefundOrderResp(resp), nil
}

func (s *OrderServiceImpl) GetOrder(ctx context.Context, req *orderapi.GetOrderReq) (*orderapi.GetOrderResp, error) {
	resp, err := s.orderApp.GetOrder(ctx, toLegacyGetOrderReq(req))
	if err != nil {
		return nil, err
	}
	return fromLegacyGetOrderResp(resp), nil
}

func (s *OrderServiceImpl) ListOrders(ctx context.Context, req *orderapi.ListOrdersReq) (*orderapi.ListOrdersResp, error) {
	resp, err := s.orderApp.ListOrders(ctx, toLegacyListOrdersReq(req))
	if err != nil {
		return nil, err
	}
	return fromLegacyListOrdersResp(resp), nil
}

func toLegacyCreateOrderReq(in *orderapi.CreateOrderReq) *legacy.CreateOrderReq {
	if in == nil {
		return nil
	}
	ps := make([]*legacy.CreateOrderPassenger, 0, len(in.Passengers))
	for _, p := range in.Passengers {
		if p == nil {
			ps = append(ps, nil)
			continue
		}
		ps = append(ps, &legacy.CreateOrderPassenger{RealName: p.RealName, IdCard: p.IdCard, SeatType: p.SeatType})
	}
	return &legacy.CreateOrderReq{
		UserId:           in.UserId,
		TrainId:          in.TrainId,
		DepartureStation: in.DepartureStation,
		ArrivalStation:   in.ArrivalStation,
		Passengers:       ps,
	}
}

func toLegacyPayOrderReq(in *orderapi.PayOrderReq) *legacy.PayOrderReq {
	if in == nil {
		return nil
	}
	return &legacy.PayOrderReq{UserId: in.UserId, OrderId: in.OrderId, PayChannel: in.PayChannel, PayNo: in.PayNo}
}

func toLegacyConfirmPayReq(in *orderapi.ConfirmPayReq) *legacy.ConfirmPayReq {
	if in == nil {
		return nil
	}
	return &legacy.ConfirmPayReq{UserId: in.UserId, OrderId: in.OrderId, PayNo: in.PayNo, ThirdPartyStatus: in.ThirdPartyStatus}
}

func toLegacyChangeOrderReq(in *orderapi.ChangeOrderReq) *legacy.ChangeOrderReq {
	if in == nil {
		return nil
	}
	return &legacy.ChangeOrderReq{
		UserId:               in.UserId,
		OrderId:              in.OrderId,
		NewTrainId_:          in.NewTrainId_,
		NewDepartureStation_: in.NewDepartureStation_,
		NewArrivalStation_:   in.NewArrivalStation_,
	}
}

func toLegacyCancelOrderReq(in *orderapi.CancelOrderReq) *legacy.CancelOrderReq {
	if in == nil {
		return nil
	}
	return &legacy.CancelOrderReq{UserId: in.UserId, OrderId: in.OrderId}
}

func toLegacyRefundOrderReq(in *orderapi.RefundOrderReq) *legacy.RefundOrderReq {
	if in == nil {
		return nil
	}
	return &legacy.RefundOrderReq{UserId: in.UserId, OrderId: in.OrderId, Reason: in.Reason}
}

func toLegacyGetOrderReq(in *orderapi.GetOrderReq) *legacy.GetOrderReq {
	if in == nil {
		return nil
	}
	return &legacy.GetOrderReq{UserId: in.UserId, OrderId: in.OrderId}
}

func toLegacyListOrdersReq(in *orderapi.ListOrdersReq) *legacy.ListOrdersReq {
	if in == nil {
		return nil
	}
	return &legacy.ListOrdersReq{UserId: in.UserId, Status: in.Status, Cursor: in.Cursor, Limit: in.Limit}
}

func fromLegacyCreateOrderResp(in *legacy.CreateOrderResp) *orderapi.CreateOrderResp {
	if in == nil {
		return nil
	}
	seats := make([]*orderapi.OrderSeatInfo, 0, len(in.Seats))
	for _, s := range in.Seats {
		if s == nil {
			continue
		}
		seats = append(seats, &orderapi.OrderSeatInfo{SeatId: s.SeatId, SeatType: s.SeatType, CarriageNum: s.CarriageNum, SeatNum: s.SeatNum, SeatPrice: s.SeatPrice})
	}
	return &orderapi.CreateOrderResp{BaseResp: toCommon(in.BaseResp), OrderId: in.OrderId, PayDeadlineUnix: in.PayDeadlineUnix, Seats: seats}
}

func fromLegacyPayOrderResp(in *legacy.PayOrderResp) *orderapi.PayOrderResp {
	if in == nil {
		return nil
	}
	return &orderapi.PayOrderResp{BaseResp: toCommon(in.BaseResp), OrderStatus: in.OrderStatus, PayStatus: in.PayStatus, PayNo: in.PayNo}
}

func fromLegacyConfirmPayResp(in *legacy.ConfirmPayResp) *orderapi.ConfirmPayResp {
	if in == nil {
		return nil
	}
	return &orderapi.ConfirmPayResp{BaseResp: toCommon(in.BaseResp), OrderStatus: in.OrderStatus}
}

func fromLegacyChangeOrderResp(in *legacy.ChangeOrderResp) *orderapi.ChangeOrderResp {
	if in == nil {
		return nil
	}
	seats := make([]*orderapi.OrderSeatInfo, 0, len(in.Seats))
	for _, s := range in.Seats {
		if s == nil {
			continue
		}
		seats = append(seats, &orderapi.OrderSeatInfo{SeatId: s.SeatId, SeatType: s.SeatType, CarriageNum: s.CarriageNum, SeatNum: s.SeatNum, SeatPrice: s.SeatPrice})
	}
	return &orderapi.ChangeOrderResp{
		BaseResp:         toCommon(in.BaseResp),
		OldOrderStatus:   in.OldOrderStatus,
		NewOrderId_:      in.NewOrderId_,
		NewTotalAmount_:  in.NewTotalAmount_,
		RefundDiffAmount: in.RefundDiffAmount,
		Seats:            seats,
	}
}

func fromLegacyRefundOrderResp(in *legacy.RefundOrderResp) *orderapi.RefundOrderResp {
	if in == nil {
		return nil
	}
	return &orderapi.RefundOrderResp{BaseResp: toCommon(in.BaseResp), RefundAmount: in.RefundAmount, RefundStatus: in.RefundStatus, OrderStatus: in.OrderStatus}
}

func fromLegacyGetOrderResp(in *legacy.GetOrderResp) *orderapi.GetOrderResp {
	if in == nil {
		return nil
	}
	seats := make([]*orderapi.OrderSeatInfo, 0, len(in.Seats))
	for _, s := range in.Seats {
		if s == nil {
			continue
		}
		seats = append(seats, &orderapi.OrderSeatInfo{SeatId: s.SeatId, SeatType: s.SeatType, CarriageNum: s.CarriageNum, SeatNum: s.SeatNum, SeatPrice: s.SeatPrice})
	}
	return &orderapi.GetOrderResp{BaseResp: toCommon(in.BaseResp), Order: fromLegacyOrderInfo(in.Order), Seats: seats}
}

func fromLegacyListOrdersResp(in *legacy.ListOrdersResp) *orderapi.ListOrdersResp {
	if in == nil {
		return nil
	}
	items := make([]*orderapi.OrderInfo, 0, len(in.Orders))
	for _, o := range in.Orders {
		items = append(items, fromLegacyOrderInfo(o))
	}
	return &orderapi.ListOrdersResp{BaseResp: toCommon(in.BaseResp), Orders: items, NextCursor: in.NextCursor}
}

func fromLegacyOrderInfo(in *legacy.OrderInfo) *orderapi.OrderInfo {
	if in == nil {
		return nil
	}
	return &orderapi.OrderInfo{
		OrderId:          in.OrderId,
		UserId:           in.UserId,
		TrainId:          in.TrainId,
		DepartureStation: in.DepartureStation,
		ArrivalStation:   in.ArrivalStation,
		FromSeq:          in.FromSeq,
		ToSeq:            in.ToSeq,
		TotalAmount:      in.TotalAmount,
		OrderStatus:      in.OrderStatus,
		PayDeadlineUnix:  in.PayDeadlineUnix,
		PayTimeUnix:      in.PayTimeUnix,
		PayChannel:       in.PayChannel,
		PayNo:            in.PayNo,
		CreatedAtUnix:    in.CreatedAtUnix,
	}
}

func toCommon(b *legacy.BaseResp) *common.BaseResp {
	if b == nil {
		return nil
	}
	return &common.BaseResp{Code: b.Code, Msg: b.Msg}
}
