# 后端接口返回结构对照表（HTTP /api/v1）

## 背景
当前项目同时存在两类接口实现方式：
- **HTTP 直连数据库（Gin + Service + GORM）**：典型是 `branches`，返回 `{code,msg,data:{...}}`。
- **HTTP 代理 RPC（Gin -> Kitex RPC）**：典型是 `members/points-records/...`，成功时直接返回 Thrift 结构体，失败时此前返回 `{error:...}`（已统一修复为 `{code,msg}`）。

这会导致前端在不同页面需要用不同方式解包响应。

## 返回结构类型
### A 类：标准包裹（code/msg/data）
```json
{ "code": 200, "msg": "查询成功", "data": { "list": [], "total": 0 } }
```

### B 类：RPC 列表结构（不带 code/msg）
```json
{ "list": [], "total": 0, "page": 1, "page_size": 10 }
```

### C 类：RPC 对象结构（不带 code/msg）
```json
{ "id": 1, "member_level": "NORMAL", "...": "..." }
```

### D 类：RPC BaseResp（带 code/msg）
```json
{ "code": 200, "msg": "创建成功" }
```

### 统一错误结构（已落地）
无论接口属于 A/B/C/D 类，**只要是参数错误或 RPC 调用失败**，HTTP 现在统一返回：
```json
{ "code": 400/500, "msg": "..." }
```

## 重点接口对照（前端已使用）
| 资源 | 路由 | 成功返回 | 失败返回（当前） | 前端消费方式（现状） |
|---|---|---|---|---|
| 分店 | GET /branches | A 类：`{code,msg,data:{list,total}}` | `{code,msg}` | `room.js` 返回 axios response，页面用 `response.data` |
| 分店 | POST/PUT/DELETE /branches | A 类：`{code,msg,data?}` | `{code,msg}` | `HotelManagement.jsx` 用 `response.data.code` |
| 会员 | GET /members | B 类：`{list,total,page,page_size}` | `{code,msg}` | `member.js` 返回 `response.data`，页面用 `response.list` |
| 会员 | GET /members/:id | C 类：Member 结构体 | `{code,msg}` | 页面直接读取字段（如 `guest_id`） |
| 会员 | POST/PUT/DELETE /members | D 类：`{code,msg}` | `{code,msg}` | 页面用 `response.code` |
| 积分记录 | GET /points-records | B 类：`{list,total,page,page_size}` | `{code,msg}` | `MemberPointsManagement.jsx` 用 `response.list` |
| 积分记录 | POST /points-records | D 类：`{code,msg}` | `{code,msg}` | 页面用 `response.code` |
| 积分余额 | GET /members/:id/points-balance | `{balance:<int64>}` | `{code,msg}` | `member.js` 返回 `response.data` |
| 会员权益 | GET /member-rights | B 类 | `{code,msg}` | `member.js` 返回 `response.data` |
| 黑名单 | /blacklists | B/C/D 类 | `{code,msg}` | `room.js` 风格页面通常用 `response.data` |
| 系统配置 | /system-configs | B/C/D 类 | `{code,msg}` | `room.js` 风格页面通常用 `response.data` |
| 角色/权限 | /roles /permissions | B/C/D 类 | `{code,msg}` | `room.js` 风格页面通常用 `response.data` |

## 前端适配建议（推荐做法）
为避免每个页面都要猜返回结构，建议在前端 `web/src/api` 增加统一的解包层：
- 对 A 类：返回 `data`（也就是把 `{code,msg,data}` 解成 `data`）。  
- 对 D 类：原样返回（保留 `{code,msg}`）。  
- 对 B/C 类：原样返回。  
- 对 `{code,msg}` 错误：统一抛错或返回并由页面提示 `msg`。

这样可以逐步把所有页面的判断逻辑统一起来，减少“接口改一点就全站炸”的风险。

## 分店接口补充说明
- `GET /api/v1/branches?status=ACTIVE`：仅返回启用分店。
- `GET /api/v1/branches?status=ALL` 或不传 `status`：返回全部状态分店（ACTIVE/INACTIVE）。
