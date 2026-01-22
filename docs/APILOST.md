# APILost 测试指南（网关 HTTP API）

## 0. 前置说明
- 网关默认监听：`http://127.0.0.1:5200`（以 `conf/config.yaml -> Server.Gateway` 为准）
- 除登录/注册/查询类接口外，订单与用户信息接口需要携带 Token
- Token 建议放在请求 Header：`Authorization: Bearer {{token}}`

建议在 APILost 新建环境变量：
- `base_url`：`http://127.0.0.1:5200`
- `token`：登录后从响应复制

## 1. 用户注册
**请求**
- Method：POST
- URL：`{{base_url}}/api/v1/user/register`
- Body(JSON)：
```json
{
  "user_name": "u1",
  "password": "p1",
  "phone": "13800138000"
}
```

**预期响应**
- `code=200`，返回 `user_id`

## 2. 用户登录（获取 Token）
**请求**
- Method：POST
- URL：`{{base_url}}/api/v1/user/login`
- Body(JSON)：
```json
{
  "phone": "13800138000",
  "password": "p1"
}
```

**预期响应**
- `code=200`
- 复制 `token` 到环境变量 `{{token}}`

## 3. 获取用户信息（需要 Token）
**请求**
- Method：GET
- URL：`{{base_url}}/api/v1/user/info`
- Header：`Authorization: Bearer {{token}}`

## 4. 实名认证（需要 Token）
**请求**
- Method：POST
- URL：`{{base_url}}/api/v1/user/verify_realname`
- Header：`Authorization: Bearer {{token}}`
- Body(JSON)：
```json
{
  "real_name": "张三",
  "id_card": "110101199001011234",
  "phone": "13800138000"
}
```

## 5. 车次查询（不需要 Token）
**请求**
- Method：GET
- URL：`{{base_url}}/api/v1/train/search?departure_station=北京&arrival_station=上海&date=2026-01-19`

说明：具体 query 参数以你当前实现的 ticket handler 为准。

## 6. 下单锁座（需要 Token）
**请求**
- Method：POST
- URL：`{{base_url}}/api/v1/order/create`
- Header：`Authorization: Bearer {{token}}`
- Body(JSON)：
```json
{
  "train_id": "TRAIN-001",
  "departure_station": "北京",
  "arrival_station": "上海",
  "passengers": [
    { "real_name": "张三", "id_card": "110101199001011234", "seat_type": "SECOND_CLASS" }
  ]
}
```

**预期响应**
- `order_id`
- `pay_deadline_unix`

## 7. 发起支付（支付宝，返回 pay_url）
**请求**
- Method：POST
- URL：`{{base_url}}/api/v1/order/pay`
- Header：`Authorization: Bearer {{token}}`
- Body(JSON)：
```json
{
  "order_id": "{{order_id}}",
  "pay_channel": "ALIPAY"
}
```

**预期响应**
- `pay_no`：本项目用作支付宝 `out_trade_no`
- `pay_url`：复制到浏览器打开，进入支付宝收银台

## 8. 支付结果回调（两种测试方式）
### 8.1 本地模拟回调（无需真实支付宝）
**请求**
- Method：POST
- URL：`{{base_url}}/api/v1/pay/mock_notify`
- Header：`Authorization: Bearer {{token}}`
- Body(JSON)：
```json
{
  "order_id": "{{order_id}}",
  "pay_no": "{{pay_no}}",
  "third_party_status": "SUCCESS"
}
```

**预期**
- 返回 `success`
- 之后查询订单状态应为已出票/已支付（以服务端状态字段为准）

### 8.2 支付宝真实异步通知（生产/沙箱）
支付宝会以 `application/x-www-form-urlencoded` 方式 POST 到：
- `{{base_url}}/api/v1/pay/callback`

网关会做：
- 用 `AliPay.AliPublicKey` RSA2 验签
- 校验 `app_id` 是否与配置一致
- 校验 `total_amount` 是否与下单金额一致（从 Redis 的 out_trade_no 映射读取）
- 调用订单服务 `ConfirmPay` 推进订单

你在 APILost 手动模拟时（仅用于排查回调处理逻辑）：
- Method：POST
- URL：`{{base_url}}/api/v1/pay/callback`
- Header：`Content-Type: application/x-www-form-urlencoded`
- Body(Form)：至少要包含 `out_trade_no`、`trade_status`、`sign`（且 sign 必须与参数匹配，否则会被拒绝）

说明：手工构造一个正确签名不现实，建议用支付宝沙箱真正走一次支付，或用项目内 `/pay/mock_notify` 完成业务链路测试。

## 9. 查询订单
**请求**
- Method：GET
- URL：`{{base_url}}/api/v1/order/info?order_id={{order_id}}`
- Header：`Authorization: Bearer {{token}}`

