package ticket

import (
	"context"
	"encoding/json"
	"example_shop/internal/gateway/http/dto"
	ticketlogic "example_shop/internal/gateway/ticket"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/websocket"
)

var wsUpgrader = websocket.HertzUpgrader{
	CheckOrigin: func(c *app.RequestContext) bool { return true },
}

// remainPush 是余票 WebSocket 推送的消息体。
// 该结构体只用于网关向客户端推送，不参与 BindAndValidate。
type remainPush struct {
	TrainID    string    `json:"train_id"`
	SeatType   string    `json:"seat_type"`
	TravelDate string    `json:"travel_date"`
	Remaining  int64     `json:"remaining"`
	Time       time.Time `json:"time"`
}

// TicketRemainWS 余票 WebSocket 推送接口（长连接定时推送）。
//
// HTTP 接口：GET /api/v1/ticket/ws（升级为 WebSocket）
// 入参：dto.TicketRemainWSReq（Query 参数）
//
// 行为说明：
// - 建立连接后每 2 秒推送一次指定车次 + 席别 + 日期的余票数量
// - 客户端收到的消息体为 remainPush
func (h *Handler) TicketRemainWS(ctx context.Context, c *app.RequestContext) {
	var req dto.TicketRemainWSReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, dto.BaseHTTPResp{Code: 400, Msg: err.Error()})
		return
	}

	wsUpgrader.Upgrade(c, func(conn *websocket.Conn) {
		defer conn.Close()

		start, _, err := ticketlogic.ParseTravelDate(req.TravelDate)
		if err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"code":400,"msg":"travel_date格式错误"}`))
			return
		}

		seatType, msg := ticketlogic.ValidateSeatType(req.SeatType)
		if msg != "" || seatType == "" {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"code":400,"msg":"seat_type必填且有效"}`))
			return
		}

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			remain, err := ticketlogic.GetRemainingSeats(context.Background(), req.TrainID, seatType, start, 2*time.Second)
			if err != nil {
				_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"code":500,"msg":"余票查询失败"}`))
				return
			}
			msg := remainPush{
				TrainID:    req.TrainID,
				SeatType:   seatType,
				TravelDate: req.TravelDate,
				Remaining:  remain,
				Time:       time.Now(),
			}
			b, _ := json.Marshal(msg)
			if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
				return
			}
			<-ticker.C
		}
	})
}
