// common.thrift
//
// 说明：
// - 公共类型定义（各域 Thrift 可 include 复用）。
// - 仅放“跨域共享且稳定”的结构体，避免把各域请求/响应堆在同一个文件里。
namespace go common

// BaseResp 统一响应基础结构。
// code：业务状态码（常见：200 成功；400 参数错误；401 未授权；404 不存在；500 服务端错误；502 下游错误）
// msg：状态说明/错误信息（用于排查与展示）
struct BaseResp{
   1: i32 code
   2: string msg
}

