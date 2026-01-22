package ticketapi

import (
	"context"
	"example_shop/internal/ticket_service/app/ticketapp"
	"example_shop/kitex_gen/common"
	"example_shop/kitex_gen/ticketapi"
	legacy "example_shop/kitex_gen/user"
)

type TicketServiceImpl struct {
	ticketApp *ticketapp.Service
}

func NewTicketServiceImpl() *TicketServiceImpl {
	return &TicketServiceImpl{ticketApp: ticketapp.New()}
}

func (s *TicketServiceImpl) GetTrainDetail(ctx context.Context, req *ticketapi.GetTrainDetailReq) (*ticketapi.GetTrainDetailResp, error) {
	legacyReq := &legacy.GetTrainDetailReq{TrainId: req.TrainId}
	if req.DepartureStation != nil {
		legacyReq.DepartureStation = req.DepartureStation
	}
	if req.ArrivalStation != nil {
		legacyReq.ArrivalStation = req.ArrivalStation
	}

	resp, err := s.ticketApp.GetTrainDetail(ctx, legacyReq)
	if err != nil {
		return nil, err
	}

	seats := make([]*ticketapi.SeatTypeRemain, 0, len(resp.SeatTypes))
	for _, it := range resp.SeatTypes {
		if it == nil {
			continue
		}
		seats = append(seats, &ticketapi.SeatTypeRemain{SeatType: it.SeatType, Remaining: it.Remaining, MinPrice: it.MinPrice})
	}

	return &ticketapi.GetTrainDetailResp{
		BaseResp:          toCommon(resp.BaseResp),
		TrainId:           resp.TrainId,
		TrainCode:         resp.TrainCode,
		TrainType:         resp.TrainType,
		DepartureStation:  resp.DepartureStation,
		ArrivalStation:    resp.ArrivalStation,
		DepartureTimeUnix: resp.DepartureTimeUnix,
		ArrivalTimeUnix:   resp.ArrivalTimeUnix,
		RuntimeMinutes:    resp.RuntimeMinutes,
		SeatTypes:         seats,
	}, nil
}

func toCommon(b *legacy.BaseResp) *common.BaseResp {
	if b == nil {
		return nil
	}
	return &common.BaseResp{Code: b.Code, Msg: b.Msg}
}
