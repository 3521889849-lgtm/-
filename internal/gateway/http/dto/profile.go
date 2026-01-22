package dto

import "example_shop/kitex_gen/userapi"

// GetProfileHTTPResp 个人信息查询响应。
//
// 对应接口：GET /api/v1/user/profile
type GetProfileHTTPResp struct {
	Code     int32             `json:"code"`
	Msg      string            `json:"msg"`
	UserInfo *userapi.UserInfo `json:"user_info,omitempty"`
}

// UpdateProfileHTTPReq 个人信息更新请求（JSON Body）。
//
// 对应接口：POST /api/v1/user/profile
//
// 说明：
// - 本项目里 real_name 同时承担“显示名/实名姓名”的角色
// - 已实名(VERIFIED) 的用户不允许通过该接口直接修改 real_name，需要走实名流程
type UpdateProfileHTTPReq struct {
	RealName string `json:"real_name,required"`
}

// UpdateProfileHTTPResp 个人信息更新响应。
type UpdateProfileHTTPResp struct {
	Code     int32             `json:"code"`
	Msg      string            `json:"msg"`
	UserInfo *userapi.UserInfo `json:"user_info,omitempty"`
}

