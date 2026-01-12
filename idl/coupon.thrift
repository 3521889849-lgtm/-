// 必须的Thrift文档头部，解决「not document」核心问题
namespace go coupon

// 极简响应体（仅占位）
struct BaseResp {
    1: i32 code,
    2: string msg
}

// 极简入参（仅占位）
struct EmptyReq {
    1: i64 id
}

// Kitex必须的Service定义（仅占位，无业务逻辑）
service CouponService {
    BaseResp Test(1: EmptyReq req)
}