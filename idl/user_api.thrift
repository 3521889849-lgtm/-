// user_api.thrift
//
// 说明：
// - 用户域 RPC 接口定义（注册/登录/实名/用户信息）。
// - 该文件只放“用户域”的请求/响应与服务，不与票务/订单混写。
namespace go userapi

include "common.thrift"

// VerifyRealNameReq 实名认证请求。
// 实名场景通常会校验：姓名 + 身份证号 + 手机号的一致性。
struct VerifyRealNameReq{
   1:string user_id
   2:string real_name
   3:string id_card
   4:string phone
}

// UserInfo 用户基本信息（对外展示用）。
struct UserInfo {
    1: string user_id
    2: string real_name
    3: string id_card
    4: string phone
    5: string user_level
    6: string real_name_verified
    7: string status
}

// GetUserInfoReq 获取用户信息请求。
struct GetUserInfoReq {
    1: string user_id
}

// GetUserInfoResp 获取用户信息响应。
struct GetUserInfoResp {
    1: common.BaseResp base_resp
    2: optional UserInfo user_info
}

// RegisterReq 用户注册请求。
struct RegisterReq {
    1: string user_name
    2: string password
    3: string phone
}

// RegisterResp 用户注册响应。
struct RegisterResp {
    1: common.BaseResp base_resp
    2: string user_id
}

// LoginReq 用户登录请求。
struct LoginReq {
    1: string phone
    2: string password
}

// LoginResp 用户登录响应。
struct LoginResp {
    1: common.BaseResp base_resp
    2: string user_id
    3: string token
    4: UserInfo user_info
}

// UserService 用户域 RPC 服务。
service UserService{
   // Register 用户注册。
   RegisterResp Register(1: RegisterReq req)
   // Login 用户登录，返回 token。
   LoginResp Login(1: LoginReq req)
   // VerifyRealName 实名认证。
   common.BaseResp VerifyRealName(1: VerifyRealNameReq req)
   // GetUserInfo 查询用户信息。
   GetUserInfoResp GetUserInfo(1: GetUserInfoReq req)
}
