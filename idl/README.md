# IDL 分层草案

- `idl/common.thrift`：跨域共享的基础类型（例如 BaseResp）。
- `idl/user_api.thrift`：用户域（注册/登录/实名/用户信息）。
- `idl/ticket_api.thrift`：票务域（车次详情等）。
- `idl/order_api.thrift`：订单域（下单/支付/改签/退票/查询）。

后续如果需要把 RPC 进程拆成 user-service/ticket-service/order-service，可基于这三份拆分草案生成各自的 Kitex 代码，并让网关分别持有不同服务的 client。
