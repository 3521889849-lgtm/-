package ticket

import (
	"context"
	"example_shop/internal/gateway/http/dto"
	"example_shop/kitex_gen/ticketapi/ticketservice"
	kitexticket "example_shop/kitex_gen/ticketapi"

	"github.com/cloudwego/hertz/pkg/app"
)

// Handler 承载票务域相关的 HTTP 接口（网关侧）。
type Handler struct {
	TicketClient ticketservice.Client
}

// TrainDetail 车次详情接口（余票与最低价）。
//
// HTTP 接口：GET /api/v1/train/detail
// 入参：dto.TrainDetailHTTPReq（Query 参数）
// 出参：dto.TrainDetailHTTPResp
//
// 说明：
// - departure_station / arrival_station 为可选参数；不传时由后端服务推导或按默认规则处理
func (h *Handler) TrainDetail(ctx context.Context, c *app.RequestContext) {
	var req dto.TrainDetailHTTPReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, dto.BaseHTTPResp{Code: 400, Msg: err.Error()})
		return
	}

	rpcReq := &kitexticket.GetTrainDetailReq{TrainId: req.TrainID}
	if req.DepartureStation != "" {
		rpcReq.DepartureStation = &req.DepartureStation
	}
	if req.ArrivalStation != "" {
		rpcReq.ArrivalStation = &req.ArrivalStation
	}
	rpcResp, err := h.TicketClient.GetTrainDetail(ctx, rpcReq)
	if err != nil {
		c.JSON(502, dto.BaseHTTPResp{Code: 502, Msg: err.Error()})
		return
	}

	items := make([]dto.SeatTypeRemainItem, 0, len(rpcResp.SeatTypes))
	for _, it := range rpcResp.SeatTypes {
		if it == nil {
			continue
		}
		items = append(items, dto.SeatTypeRemainItem{SeatType: it.SeatType, Remaining: it.Remaining, MinPrice: it.MinPrice})
	}

	c.JSON(200, dto.TrainDetailHTTPResp{
		Code:              rpcResp.BaseResp.Code,
		Msg:               rpcResp.BaseResp.Msg,
		TrainID:           rpcResp.TrainId,
		TrainCode:         rpcResp.TrainCode,
		TrainType:         rpcResp.TrainType,
		DepartureStation:  rpcResp.DepartureStation,
		ArrivalStation:    rpcResp.ArrivalStation,
		DepartureTimeUnix: rpcResp.DepartureTimeUnix,
		ArrivalTimeUnix:   rpcResp.ArrivalTimeUnix,
		RuntimeMinutes:    rpcResp.RuntimeMinutes,
		SeatTypes:         items,
	})
}
