# 项目说明（example_shop / piaowu）

## 1. 项目定位
- 这是一个基于 CloudWeGo（Hertz + Kitex）的“车票/订单”示例项目
- 对外提供 HTTP API（Gateway），对内通过 RPC 调用后端服务（Ticket/User/Order）
- 支持：用户注册/登录（JWT）、车次/余票查询、下单锁座、支付推进、退票/改签

## 2. 目录结构（核心）
- cmd：各服务启动入口
  - cmd/gateway：HTTP 网关
  - cmd/ticket_service：票务/订单/用户后端（Kitex Server）
  - cmd/order_service：订单服务（当前实现复用了 ticket_service 的订单用例）
  - cmd/user_service：少量适配代码（当前订单 handler 也在这里做了一层转发）
- common：公共能力
  - common/config：配置加载（Viper）与全局配置对象
  - common/db：MySQL/Redis 初始化与全局连接
  - common/init：统一初始化入口（配置 -> MySQL -> Redis）
- internal/gateway：网关 HTTP 层（路由、handler、中间件）
- internal/ticket_service：后端业务用例（用户/票务/订单）
- pkg：可复用工具包（jwt、alipay、password、realname 等）
- conf/config.yaml：主配置文件
- idl / kitex_gen：RPC IDL 与生成代码

## 3. 配置加载与全局引用
### 3.1 配置入口
- 配置读取入口：[ViperInit](file:///d:/gowork/2304a/piaowu/common/config/viperInit.go#L33-L70)
- 统一初始化入口：[Init](file:///d:/gowork/2304a/piaowu/common/init/init.go#L29-L49)

### 3.2 关键点
- 配置对象：`common/config.Cfg`（全局可引用）
- 已增加 `ValidateAndNormalize()`：加载后填充默认值并做关键字段校验（避免“yaml 能读但运行时炸”）
- Redis 端口已改为 `int`，避免 yaml 数值端口导致反序列化失败

### 3.3 config.yaml 字段映射
对应结构体：[config.go](file:///d:/gowork/2304a/piaowu/common/config/config.go)
- Mysql：Host/Port/User/Password/Database/LogLevel/AutoMigrate/StartupFixes/ReadReplica
- Redis：Host/Port/Password/Database
- Server：Gateway/UserService/TicketService/OrderService
- JWT：Secret/Expire
- RealName：SecretID/SecretKey/DebugReturn
- AliPay：AppId/PrivateKey/AliPublicKey/NotifyURL/ReturnURL/IsProduction

## 4. 鉴权（JWT）
### 4.1 登录发 Token
- 登录 HTTP 接口：`POST /api/v1/user/login`
- 后端签发 Token：在 [userapp.Login](file:///d:/gowork/2304a/piaowu/internal/ticket_service/app/userapp/service.go#L73-L118) 内调用 `pkg/jwt.GenerateToken`

### 4.2 网关验 Token
- 中间件： [AuthRequired](file:///d:/gowork/2304a/piaowu/internal/gateway/http/middleware/auth.go)
- 客户端携带方式：
  - `Authorization: Bearer <token>`（推荐）
  - 或 `X-Token: <token>`
- 网关会把 user_id 写入上下文，同时写入请求头 `X-User-ID`（便于限流等中间件复用）

### 4.3 需要登录的接口
路由注册：[router.go](file:///d:/gowork/2304a/piaowu/internal/gateway/http/router/router.go)
- 用户：`/user/verify_realname`、`/user/info`
- 订单：`/order/*`、`/pay/mock_notify`
- 回调：`/pay/callback` 不需要登录（第三方回调场景）

## 5. 支付宝支付对接（WAP + 异步通知）
### 5.1 当前对接方式（符合本项目“订单推进”模型）
- 订单服务 `PayOrder`：把订单推进到 PAYING，并生成 pay_no（本项目将 pay_no 作为 out_trade_no）
- 网关在 `pay_channel=ALIPAY` 时生成支付宝收银台 URL，返回给客户端打开支付
- 支付宝异步通知回调到网关 `/api/v1/pay/callback`
- 网关验签通过后，调用订单服务 `ConfirmPay` 推进订单到 ISSUED（出票）

### 5.2 关键实现
- 生成支付 URL（手机网站支付 TradeWapPay）：[WapPayURL](file:///d:/gowork/2304a/piaowu/pkg/alipay/alipay.go#L59-L105)
- 异步通知验签（优先用 AliPublicKey 做 RSA2）：[Verify](file:///d:/gowork/2304a/piaowu/pkg/alipay/alipay.go#L107-L132)
- 网关支付接口返回 pay_url，并写 Redis 映射 out_trade_no -> 订单信息：[PayOrder](file:///d:/gowork/2304a/piaowu/internal/gateway/http/handler/order/handler.go#L80-L141)
- 网关回调：form 回调会验签、校验 app_id/金额一致性、再推进订单：[PayCallback](file:///d:/gowork/2304a/piaowu/internal/gateway/http/handler/order/handler.go#L143-L226)

### 5.3 必要配置
在 `conf/config.yaml` 的 `AliPay` 下配置：
- AppId：应用 APPID
- PrivateKey：应用私钥（用于请求签名/生成支付 URL）
- AliPublicKey：支付宝公钥（用于回调验签）
- NotifyURL：支付宝异步通知地址（公网可访问；本地调试可先留空）
- ReturnURL：同步跳转地址（可选）
- IsProduction：是否生产环境

## 6. 服务启动关系（常用）
- 网关：cmd/gateway
- 票务/用户/订单业务：cmd/ticket_service、cmd/order_service（当前订单服务做了适配转发）

## 7. 前端页面（仿携程 Demo）
- 目录：web/
- 由网关直接静态托管：访问 `http://127.0.0.1:5200/` 即可打开
- 页面能力：车次查询、登录/注册、下单锁座、生成支付宝 pay_url、模拟支付回调、订单列表

## 8. 架构图（流程图 / ER 图）

### 8.1 系统流程图

![系统流程图](diagrams/flowchart.svg)

Mermaid 源文件：`docs/diagrams/flowchart.mmd`

### 8.2 数据库 ER 图

![数据库ER图](diagrams/er.svg)

Mermaid 源文件：`docs/diagrams/er.mmd`
