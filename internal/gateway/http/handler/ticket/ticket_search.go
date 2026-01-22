package ticket

import (
	"context"
	"example_shop/internal/gateway/http/dto"
	ticketlogic "example_shop/internal/gateway/ticket"

	"github.com/cloudwego/hertz/pkg/app"
)

// SearchTrain 车次查询接口（票务查询）。
//
// HTTP 接口：GET /api/v1/train/search
// 入参：dto.SearchTrainHTTPReq（Query 参数）
// 出参：dto.SearchTrainHTTPResp
func (h *Handler) SearchTrain(ctx context.Context, c *app.RequestContext) {
	var req dto.SearchTrainHTTPReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, dto.BaseHTTPResp{Code: 400, Msg: err.Error()})
		return
	}
	res := ticketlogic.New().SearchTrain(ctx, req)
	c.JSON(res.Status, res.Body)
}
