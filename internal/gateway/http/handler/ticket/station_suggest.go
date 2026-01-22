package ticket

import (
	"context"
	"example_shop/internal/gateway/http/dto"
	ticketlogic "example_shop/internal/gateway/ticket"

	"github.com/cloudwego/hertz/pkg/app"
)

// StationSuggest 站点联想接口（用于输入出发站/到达站时的候选提示）。
//
// HTTP 接口：GET /api/v1/station/suggest
// 入参：dto.StationSuggestHTTPReq（Query 参数）
// 出参：dto.StationSuggestHTTPResp
func (h *Handler) StationSuggest(ctx context.Context, c *app.RequestContext) {
	var req dto.StationSuggestHTTPReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, dto.BaseHTTPResp{Code: 400, Msg: err.Error()})
		return
	}
	res := ticketlogic.New().StationSuggest(ctx, req)
	c.JSON(res.Status, res.Body)
}
