package user

import (
	"context"
	"example_shop/common/db"
	"example_shop/common/encrypt"
	"example_shop/internal/gateway/http/dto"
	"example_shop/internal/gateway/http/middleware"
	"example_shop/internal/ticket_service/model"
	"example_shop/kitex_gen/userapi"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

// GetProfile 获取当前登录用户的个人信息。
//
// HTTP 接口：GET /api/v1/user/profile
func (h *Handler) GetProfile(ctx context.Context, c *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(401, dto.BaseHTTPResp{Code: 401, Msg: "未登录"})
		return
	}

	var u model.UserInfo
	if err := db.MysqlDB.Where("user_id = ?", userID).First(&u).Error; err != nil {
		c.JSON(500, dto.BaseHTTPResp{Code: 500, Msg: "查询失败: " + err.Error()})
		return
	}

	idCard := ""
	if u.IDCard != nil {
		idCard = *u.IDCard
	}

	c.JSON(200, dto.GetProfileHTTPResp{
		Code: 200,
		Msg:  "success",
		UserInfo: &userapi.UserInfo{
			UserId:           u.ID,
			RealName:         u.RealName,
			IdCard:           encrypt.MaskIDCard(idCard),
			Phone:            encrypt.MaskPhone(u.Phone),
			UserLevel:        u.UserLevel,
			RealNameVerified: u.RealNameVerified,
			Status:           u.Status,
		},
	})
}

// UpdateProfile 更新当前登录用户的个人信息（简化版：仅更新 real_name）。
//
// HTTP 接口：POST /api/v1/user/profile
func (h *Handler) UpdateProfile(ctx context.Context, c *app.RequestContext) {
	var req dto.UpdateProfileHTTPReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, dto.BaseHTTPResp{Code: 400, Msg: err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(401, dto.BaseHTTPResp{Code: 401, Msg: "未登录"})
		return
	}

	realName := strings.TrimSpace(req.RealName)
	if len(realName) < 2 || len(realName) > 32 {
		c.JSON(400, dto.BaseHTTPResp{Code: 400, Msg: "real_name长度不合法"})
		return
	}

	var u model.UserInfo
	if err := db.MysqlDB.Where("user_id = ?", userID).First(&u).Error; err != nil {
		c.JSON(500, dto.BaseHTTPResp{Code: 500, Msg: "查询失败: " + err.Error()})
		return
	}
	if strings.EqualFold(strings.TrimSpace(u.RealNameVerified), "VERIFIED") {
		c.JSON(400, dto.BaseHTTPResp{Code: 400, Msg: "已实名用户不允许直接修改姓名"})
		return
	}

	if err := db.MysqlDB.Model(&model.UserInfo{}).Where("user_id = ?", userID).
		Updates(map[string]any{"real_name": realName, "updated_at": time.Now()}).Error; err != nil {
		c.JSON(500, dto.BaseHTTPResp{Code: 500, Msg: "更新失败: " + err.Error()})
		return
	}

	u.RealName = realName
	idCard := ""
	if u.IDCard != nil {
		idCard = *u.IDCard
	}

	c.JSON(200, dto.UpdateProfileHTTPResp{
		Code: 200,
		Msg:  "success",
		UserInfo: &userapi.UserInfo{
			UserId:           u.ID,
			RealName:         u.RealName,
			IdCard:           encrypt.MaskIDCard(idCard),
			Phone:            encrypt.MaskPhone(u.Phone),
			UserLevel:        u.UserLevel,
			RealNameVerified: u.RealNameVerified,
			Status:           u.Status,
		},
	})
}

