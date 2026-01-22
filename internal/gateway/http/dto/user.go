/*
 * Gateway HTTP数据传输对象（DTO）
 *
 * 功能说明：
 * - 定义Gateway与客户端之间的数据传输格式
 * - 与RPC层的请求/响应结构体分离，便于版本控制和扩展
 *
 * 命名规范：
 * - HTTPReq: HTTP请求结构体
 * - HTTPResp: HTTP响应结构体
 * - required标签：表示该字段为必填项（Hertz自动验证）
 * - json标签：JSON序列化字段名
 * - query标签：URL查询参数绑定
 */
package dto

import "example_shop/kitex_gen/userapi"

// VerifyRealNameHTTPReq 实名认证HTTP请求结构体
type VerifyRealNameHTTPReq struct {
	UserID   string `json:"user_id"`            // 用户ID（由Token推导；兼容旧客户端可传）
	RealName string `json:"real_name,required"` // 真实姓名（必填）
	IDCard   string `json:"id_card,required"`   // 身份证号（必填）
	Phone    string `json:"phone,required"`     // 手机号（必填）
}

// BaseHTTPResp 基础HTTP响应结构体
// 所有HTTP响应都包含Code和Msg字段
type BaseHTTPResp struct {
	Code int32  `json:"code"` // 响应码：200成功，400参数错误，500服务器错误等
	Msg  string `json:"msg"`  // 响应消息：成功或错误描述
}

// RegisterHTTPReq 用户注册HTTP请求结构体
type RegisterHTTPReq struct {
	UserName string `json:"user_name,required"` // 用户名（必填）
	Password string `json:"password,required"`  // 密码（必填）
	Phone    string `json:"phone,required"`     // 手机号（必填，作为登录账号）
}

// RegisterHTTPResp 用户注册HTTP响应结构体
type RegisterHTTPResp struct {
	Code   int32  `json:"code"`    // 响应码
	Msg    string `json:"msg"`     // 响应消息
	UserID string `json:"user_id"` // 注册成功返回的用户ID
}

// LoginHTTPReq 用户登录HTTP请求结构体
type LoginHTTPReq struct {
	Phone    string `json:"phone,required"`    // 手机号（登录账号）
	Password string `json:"password,required"` // 密码
}

// LoginHTTPResp 用户登录HTTP响应结构体
type LoginHTTPResp struct {
	Code     int32             `json:"code"`      // 响应码
	Msg      string            `json:"msg"`       // 响应消息
	UserID   string            `json:"user_id"`   // 用户ID
	Token    string            `json:"token"`     // 登录令牌（后续请求需携带）
	UserInfo *userapi.UserInfo `json:"user_info"` // 用户详细信息
}

// GetUserInfoHTTPReq 获取用户信息HTTP请求结构体
// query标签：表示从URL查询参数中获取（如：/user/info?user_id=xxx）
type GetUserInfoHTTPReq struct {
	UserID string `json:"user_id"` // 用户ID（由Token推导；兼容旧客户端可传）
}

// GetUserInfoHTTPResp 获取用户信息HTTP响应结构体
type GetUserInfoHTTPResp struct {
	Code     int32             `json:"code"`      // 响应码
	Msg      string            `json:"msg"`       // 响应消息
	UserInfo *userapi.UserInfo `json:"user_info"` // 用户详细信息
}
