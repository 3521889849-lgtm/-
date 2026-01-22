package user

import (
	"context"
	"example_shop/common/db"
	"example_shop/common/encrypt"
	"example_shop/internal/gateway/http/dto"
	"example_shop/internal/gateway/http/middleware"
	"example_shop/internal/ticket_service/model"

	"github.com/cloudwego/hertz/pkg/app"
)

func (h *Handler) ListPassengers(ctx context.Context, c *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(401, dto.BaseHTTPResp{Code: 401, Msg: "未登录"})
		return
	}

	var rows []model.PassengerInfo
	if err := db.MysqlDB.Where("user_id = ?", userID).Order("id desc").Limit(200).Find(&rows).Error; err != nil {
		c.JSON(500, dto.BaseHTTPResp{Code: 500, Msg: "查询乘车人失败: " + err.Error()})
		return
	}

	seen := make(map[string]struct{}, 32)
	out := make([]dto.PassengerBriefHTTP, 0, 20)
	for _, r := range rows {
		key := r.RealName + "|" + r.IDCard
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, dto.PassengerBriefHTTP{
			PassengerID: r.ID,
			RealName:    r.RealName,
			IDCard:      encrypt.MaskIDCard(r.IDCard),
		})
		if len(out) >= 20 {
			break
		}
	}

	c.JSON(200, dto.ListPassengersHTTPResp{Code: 200, Msg: "success", Passengers: out})
}

