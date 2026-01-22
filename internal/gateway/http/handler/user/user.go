/*
 * Gateway HTTP请求处理器
 *
 * 功能说明：
 * - 处理来自客户端的HTTP请求
 * - 将HTTP请求转换为RPC调用
 * - 处理请求参数验证和响应格式化
 *
 * 职责划分：
 * - Gateway Handler：负责HTTP层面的处理（参数绑定、响应格式化）
 * - RPC Service：负责业务逻辑处理（数据库操作、业务规则）
 */
package user

import (
	"context"
	"example_shop/internal/gateway/http/dto"
	"example_shop/internal/gateway/http/middleware"
	kitexuser "example_shop/kitex_gen/userapi"
	"example_shop/kitex_gen/userapi/userservice"

	"github.com/cloudwego/hertz/pkg/app"
)

// Handler 承载用户域相关的 HTTP 接口（网关侧）。
// 该层只做 HTTP <-> RPC 的协议转换与基础参数校验，业务逻辑在后端服务。
type Handler struct {
	UserClient userservice.Client
}

// Register 用户注册接口处理函数
//
// 处理流程：
// 1. 绑定并验证HTTP请求参数（DTO）
// 2. 转换为RPC请求参数
// 3. 调用User Service的Register方法
// 4. 将RPC响应转换为HTTP响应
//
// HTTP接口：POST /api/v1/user/register
func (h *Handler) Register(ctx context.Context, c *app.RequestContext) {
	// 1. 绑定HTTP请求参数到DTO结构体
	// BindAndValidate：自动解析JSON请求体并验证required标签
	var req dto.RegisterHTTPReq
	if err := c.BindAndValidate(&req); err != nil {
		// 参数验证失败，返回400错误
		c.JSON(400, dto.BaseHTTPResp{Code: 400, Msg: err.Error()})
		return
	}

	// 2. 调用RPC服务：将HTTP请求转换为RPC调用
	// ctx：传递上下文（可用于链路追踪、超时控制等）
	rpcResp, err := h.UserClient.Register(ctx, &kitexuser.RegisterReq{
		UserName: req.UserName,
		Password: req.Password,
		Phone:    req.Phone,
	})
	if err != nil {
		// RPC调用失败，返回502错误（网关错误）
		c.JSON(502, dto.BaseHTTPResp{Code: 502, Msg: err.Error()})
		return
	}

	// 3. 将RPC响应转换为HTTP响应
	// 返回200状态码和业务响应数据
	c.JSON(200, dto.RegisterHTTPResp{
		Code:   rpcResp.BaseResp.Code,
		Msg:    rpcResp.BaseResp.Msg,
		UserID: rpcResp.UserId,
	})
}

// Login 用户登录接口处理函数。
//
// HTTP接口：POST /api/v1/user/login
// 入参：dto.LoginHTTPReq（JSON Body）
// 出参：dto.LoginHTTPResp
func (h *Handler) Login(ctx context.Context, c *app.RequestContext) {
	var req dto.LoginHTTPReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, dto.BaseHTTPResp{Code: 400, Msg: err.Error()})
		return
	}

	rpcResp, err := h.UserClient.Login(ctx, &kitexuser.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(502, dto.BaseHTTPResp{Code: 502, Msg: err.Error()})
		return
	}

	// 返回登录信息和Token
	c.JSON(200, dto.LoginHTTPResp{
		Code:     rpcResp.BaseResp.Code,
		Msg:      rpcResp.BaseResp.Msg,
		UserID:   rpcResp.UserId,
		Token:    rpcResp.Token,    // 客户端后续请求需要携带此Token
		UserInfo: rpcResp.UserInfo, // 用户详细信息
	})
}

// VerifyRealName 实名认证接口处理函数
//
// HTTP接口：POST /api/v1/user/verify_realname
//
// 功能：用户提交身份证信息进行实名认证
func (h *Handler) VerifyRealName(ctx context.Context, c *app.RequestContext) {
	var req dto.VerifyRealNameHTTPReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, dto.BaseHTTPResp{Code: 400, Msg: err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(401, dto.BaseHTTPResp{Code: 401, Msg: "未登录"})
		return
	}
	if req.UserID != "" && req.UserID != userID {
		c.JSON(403, dto.BaseHTTPResp{Code: 403, Msg: "user_id不匹配"})
		return
	}

	rpcResp, err := h.UserClient.VerifyRealName(ctx, &kitexuser.VerifyRealNameReq{
		UserId:   userID,
		RealName: req.RealName,
		IdCard:   req.IDCard,
		Phone:    req.Phone,
	})
	if err != nil {
		c.JSON(502, dto.BaseHTTPResp{Code: 502, Msg: err.Error()})
		return
	}

	c.JSON(200, dto.BaseHTTPResp{Code: rpcResp.Code, Msg: rpcResp.Msg})
}

// GetUserInfo 获取用户信息接口处理函数
//
// HTTP接口：GET /api/v1/user/info?user_id=xxx
//
// 参数：user_id通过query参数传递
func (h *Handler) GetUserInfo(ctx context.Context, c *app.RequestContext) {
	// query参数绑定：从URL参数中获取user_id
	var req dto.GetUserInfoHTTPReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, dto.BaseHTTPResp{Code: 400, Msg: err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(401, dto.BaseHTTPResp{Code: 401, Msg: "未登录"})
		return
	}
	if req.UserID != "" && req.UserID != userID {
		c.JSON(403, dto.BaseHTTPResp{Code: 403, Msg: "user_id不匹配"})
		return
	}

	// 调用RPC服务查询用户信息
	rpcResp, err := h.UserClient.GetUserInfo(ctx, &kitexuser.GetUserInfoReq{
		UserId: userID,
	})
	if err != nil {
		c.JSON(502, dto.BaseHTTPResp{Code: 502, Msg: err.Error()})
		return
	}

	// 返回用户信息
	c.JSON(200, dto.GetUserInfoHTTPResp{
		Code:     rpcResp.BaseResp.Code,
		Msg:      rpcResp.BaseResp.Msg,
		UserInfo: rpcResp.UserInfo,
	})
}
