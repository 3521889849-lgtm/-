# 酒店住宿管理后台系统 - API接口文档

> 本文档基于 `api/router/router.go` 和 `api/handler/*.go` 生成  
> 生成时间：2025-01-16  
> API版本：v1  
> 基础路径：`/api/v1`

---

## 目录

1. [通用说明](#通用说明)
2. [响应格式](#响应格式)
3. [房型管理](#房型管理)
4. [房源管理](#房源管理)
5. [设施管理](#设施管理)
6. [退订政策管理](#退订政策管理)
7. [房态管理](#房态管理)
8. [分店管理](#分店管理)
9. [订单管理](#订单管理)
10. [客人管理](#客人管理)
11. [会员管理](#会员管理)
12. [财务管理](#财务管理)
13. [系统管理](#系统管理)
14. [权限管理](#权限管理)

---

## 通用说明

### 请求头

所有请求建议包含以下请求头：

```
Content-Type: application/json
X-Request-Id: <请求ID>（可选，系统会自动生成）
```

### 响应状态码

| HTTP状态码 | 说明 |
|-----------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

### CORS支持

API支持跨域请求，已配置CORS策略：
- 允许的方法：POST, GET, PUT, DELETE, OPTIONS
- 允许的请求头：Content-Type, Authorization等

---

## 响应格式

### 成功响应

```json
{
  "code": 200,
  "msg": "操作成功",
  "data": {
    // 响应数据
  }
}
```

### 错误响应

```json
{
  "code": 400,  // 或 404, 500 等
  "msg": "错误信息描述"
}
```

### 分页响应

```json
{
  "code": 200,
  "msg": "获取成功",
  "data": {
    "list": [...],        // 数据列表
    "total": 100,         // 总记录数
    "page": 1,            // 当前页码
    "page_size": 10       // 每页数量
  }
}
```

---

## 房型管理

### 1. 创建房型

**接口地址：** `POST /api/v1/room-types`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| room_type_name | string | 是 | 房型名称（全局唯一） |
| description | string | 否 | 房型描述 |

**请求示例：**

```json
{
  "room_type_name": "标准间",
  "description": "标准双床房，面积25㎡"
}
```

**响应示例：**

```json
{
  "code": 200,
  "msg": "创建成功",
  "data": {
    "id": 1,
    "room_type_name": "标准间",
    "description": "标准双床房，面积25㎡",
    "status": "ACTIVE",
    "created_at": "2025-01-16T10:00:00Z"
  }
}
```

---

### 2. 获取房型列表

**接口地址：** `GET /api/v1/room-types`

**查询参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认10 |
| status | string | 否 | 状态筛选（ACTIVE/INACTIVE） |
| keyword | string | 否 | 关键词搜索（房型名称） |

**响应示例：**

```json
{
  "code": 200,
  "msg": "获取成功",
  "data": {
    "list": [
      {
        "id": 1,
        "room_type_name": "标准间",
        "description": "标准双床房",
        "status": "ACTIVE"
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 10
  }
}
```

---

### 3. 获取房型详情

**接口地址：** `GET /api/v1/room-types/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 房型ID |

**响应示例：**

```json
{
  "code": 200,
  "msg": "获取成功",
  "data": {
    "id": 1,
    "room_type_name": "标准间",
    "description": "标准双床房，面积25㎡",
    "status": "ACTIVE"
  }
}
```

---

### 4. 更新房型

**接口地址：** `PUT /api/v1/room-types/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 房型ID |

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| room_type_name | string | 否 | 房型名称 |
| description | string | 否 | 房型描述 |
| status | string | 否 | 状态（ACTIVE/INACTIVE） |

---

### 5. 删除房型

**接口地址：** `DELETE /api/v1/room-types/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 房型ID |

**响应示例：**

```json
{
  "code": 200,
  "msg": "删除成功"
}
```

---

## 房源管理

### 1. 创建房源

**接口地址：** `POST /api/v1/room-infos`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| branch_id | uint64 | 是 | 分店ID |
| room_type_id | uint64 | 是 | 房型ID |
| room_no | string | 是 | 房间号 |
| room_name | string | 是 | 房间名称 |
| market_price | float64 | 是 | 门市价 |
| calendar_price | float64 | 是 | 日历价 |
| room_count | uint8 | 是 | 房间数量 |
| area | float64 | 否 | 面积（平方米） |
| bed_spec | string | 是 | 床型规格 |
| has_breakfast | bool | 否 | 是否含早餐 |
| has_toiletries | bool | 否 | 是否提供洗漱用品 |
| cancellation_policy_id | uint64 | 否 | 退订政策ID |
| created_by | uint64 | 是 | 创建人ID |

**请求示例：**

```json
{
  "branch_id": 1,
  "room_type_id": 1,
  "room_no": "101",
  "room_name": "舒适大床房",
  "market_price": 298.00,
  "calendar_price": 258.00,
  "room_count": 1,
  "area": 25.50,
  "bed_spec": "1张大床",
  "has_breakfast": true,
  "has_toiletries": true,
  "created_by": 1
}
```

---

### 2. 获取房源列表

**接口地址：** `GET /api/v1/room-infos`

**查询参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认10 |
| branch_id | uint64 | 否 | 分店ID筛选 |
| room_type_id | uint64 | 否 | 房型ID筛选 |
| status | string | 否 | 状态筛选（ACTIVE/INACTIVE/MAINTENANCE） |
| keyword | string | 否 | 关键词搜索（房间名称/房间号） |

---

### 3. 获取房源详情

**接口地址：** `GET /api/v1/room-infos/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 房源ID |

**响应说明：** 包含房源信息及关联的分店、房型、设施、图片等信息

---

### 4. 更新房源

**接口地址：** `PUT /api/v1/room-infos/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 房源ID |

**请求参数：** 所有字段均为可选，只更新传入的非空字段

---

### 5. 删除房源

**接口地址：** `DELETE /api/v1/room-infos/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 房源ID |

**业务规则：** 如果房源已被订单使用，不允许删除

---

### 6. 更新房源状态

**接口地址：** `PUT /api/v1/room-infos/:id/status`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 房源ID |

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| status | string | 是 | 状态（ACTIVE-启用，INACTIVE-停用，MAINTENANCE-维修） |

---

### 7. 批量更新房源状态

**接口地址：** `PUT /api/v1/room-infos/batch-status`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| room_ids | []uint64 | 是 | 房源ID列表 |
| status | string | 是 | 状态 |

**请求示例：**

```json
{
  "room_ids": [1, 2, 3],
  "status": "MAINTENANCE"
}
```

---

### 8. 设置房源设施

**接口地址：** `PUT /api/v1/room-infos/:id/facilities`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 房源ID |

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| facility_ids | []uint64 | 是 | 设施ID列表 |

**业务说明：** 先删除旧的关联关系，再创建新的关联关系

---

### 9. 获取房源设施列表

**接口地址：** `GET /api/v1/room-infos/:id/facilities`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 房源ID |

---

### 10. 添加房源设施

**接口地址：** `POST /api/v1/room-infos/:id/facilities`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 房源ID |

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| facility_id | uint64 | 是 | 设施ID |

---

### 11. 移除房源设施

**接口地址：** `DELETE /api/v1/room-infos/:id/facilities/:facility_id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 房源ID |
| facility_id | int | 是 | 设施ID |

---

### 12. 上传房源图片

**接口地址：** `POST /api/v1/room-infos/:id/images`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 房源ID |

**请求格式：** `multipart/form-data`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| files | []File | 是 | 图片文件列表（最多16张，支持jpg/png格式，单张最大5MB） |

**业务规则：**
- 图片会自动压缩为 400x300 尺寸
- 图片按上传顺序自动排序

---

### 13. 获取房源图片列表

**接口地址：** `GET /api/v1/room-infos/:id/images`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 房源ID |

**响应示例：**

```json
{
  "code": 200,
  "msg": "获取成功",
  "data": [
    {
      "id": 1,
      "room_id": 1,
      "image_url": "/uploads/room_images/1_1768448657_0.jpg",
      "image_size": "400x300",
      "image_format": "jpg",
      "sort_order": 0
    }
  ]
}
```

---

### 14. 删除房源图片

**接口地址：** `DELETE /api/v1/room-infos/images/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 图片ID |

---

### 15. 更新图片排序

**接口地址：** `PUT /api/v1/room-infos/images/:id/sort`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 图片ID |

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| sort_order | uint8 | 是 | 排序序号 |

---

### 16. 批量更新图片排序

**接口地址：** `PUT /api/v1/room-infos/:id/images/sort`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 房源ID |

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| image_sort_orders | array | 是 | 排序列表，每个元素包含 `image_id` 和 `sort_order` |

**请求示例：**

```json
{
  "image_sort_orders": [
    {"image_id": 1, "sort_order": 0},
    {"image_id": 2, "sort_order": 1}
  ]
}
```

---

### 17. 创建关联房绑定

**接口地址：** `POST /api/v1/room-infos/bindings`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| main_room_id | uint64 | 是 | 主房源ID |
| related_room_id | uint64 | 是 | 关联房源ID |
| binding_desc | string | 否 | 绑定描述 |

**业务规则：** 不能将房源关联到自己，不能重复绑定

---

### 18. 批量创建关联房绑定

**接口地址：** `POST /api/v1/room-infos/batch-bindings`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| main_room_id | uint64 | 是 | 主房源ID |
| related_room_ids | []uint64 | 是 | 关联房源ID列表 |
| binding_desc | string | 否 | 绑定描述 |

---

### 19. 获取关联房列表

**接口地址：** `GET /api/v1/room-infos/:id/bindings`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 主房源ID |

---

### 20. 删除关联房绑定

**接口地址：** `DELETE /api/v1/room-infos/bindings/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 绑定ID |

---

## 设施管理

### 1. 创建设施

**接口地址：** `POST /api/v1/facilities`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| facility_name | string | 是 | 设施名称（全局唯一） |
| description | string | 否 | 设施描述 |

**请求示例：**

```json
{
  "facility_name": "WiFi",
  "description": "免费无线网络"
}
```

---

### 2. 获取设施列表

**接口地址：** `GET /api/v1/facilities`

**查询参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认10 |
| status | string | 否 | 状态筛选（ACTIVE/INACTIVE） |
| keyword | string | 否 | 关键词搜索（设施名称） |

---

### 3. 获取设施详情

**接口地址：** `GET /api/v1/facilities/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 设施ID |

---

### 4. 更新设施

**接口地址：** `PUT /api/v1/facilities/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 设施ID |

**请求参数：** 所有字段均为可选

---

### 5. 删除设施

**接口地址：** `DELETE /api/v1/facilities/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 设施ID |

**业务规则：** 如果设施正在被房源使用，不允许删除

---

## 退订政策管理

### 1. 创建退订政策

**接口地址：** `POST /api/v1/cancellation-policies`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| policy_name | string | 是 | 政策名称 |
| rule_description | string | 是 | 规则描述 |
| penalty_ratio | float64 | 是 | 违约金比例 |
| room_type_id | uint64 | 否 | 适用房型ID |

**请求示例：**

```json
{
  "policy_name": "24小时免费取消",
  "rule_description": "提前24小时取消订单，不收取违约金",
  "penalty_ratio": 0.00,
  "room_type_id": 1
}
```

---

### 2. 获取退订政策列表

**接口地址：** `GET /api/v1/cancellation-policies`

**查询参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| room_type_id | uint64 | 否 | 房型ID筛选 |
| status | string | 否 | 状态筛选 |
| keyword | string | 否 | 关键词搜索 |

---

### 3. 获取退订政策详情

**接口地址：** `GET /api/v1/cancellation-policies/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 政策ID |

---

### 4. 更新退订政策

**接口地址：** `PUT /api/v1/cancellation-policies/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 政策ID |

---

### 5. 删除退订政策

**接口地址：** `DELETE /api/v1/cancellation-policies/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 政策ID |

**业务规则：** 如果政策正在被房源或订单使用，不允许删除

---

## 房态管理

### 1. 获取日历化房态

**接口地址：** `GET /api/v1/calendar-room-status`

**查询参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| branch_id | uint64 | 否 | 分店ID |
| start_date | string | 是 | 开始日期（YYYY-MM-DD） |
| end_date | string | 是 | 结束日期（YYYY-MM-DD） |
| room_no | string | 否 | 房间号筛选 |
| status | string | 否 | 房态筛选 |

**业务规则：**
- 日期范围不能超过90天
- 开始日期不能晚于结束日期

---

### 2. 更新日历化房态

**接口地址：** `PUT /api/v1/calendar-room-status`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| room_id | uint64 | 是 | 房源ID |
| date | string | 是 | 日期（YYYY-MM-DD） |
| status | string | 是 | 房态（空净房/入住房/维修房/锁定房/空账房/预定房） |

---

### 3. 批量更新日历化房态

**接口地址：** `PUT /api/v1/calendar-room-status/batch`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| updates | array | 是 | 更新列表，每个元素包含 `room_id`、`date`、`status` |

---

### 4. 获取实时数据统计

**接口地址：** `GET /api/v1/real-time-statistics`

**查询参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| branch_id | uint64 | 否 | 分店ID |
| date | string | 否 | 日期（YYYY-MM-DD），默认为今日 |
| room_no | string | 否 | 房间号筛选 |
| room_type_id | uint64 | 否 | 房型ID筛选 |

**响应示例：**

```json
{
  "code": 200,
  "msg": "获取成功",
  "data": {
    "date": "2025-01-16",
    "total_rooms": 50,
    "remaining_rooms": 30,
    "checked_in_count": 20,
    "check_out_pending_count": 5,
    "reserved_pending_count": 10,
    "occupied_rooms": 15,
    "maintenance_rooms": 2,
    "locked_rooms": 1,
    "empty_rooms": 30,
    "reserved_rooms": 2,
    "status_breakdown": [
      {"status": "空净房", "count": 30},
      {"status": "入住房", "count": 15}
    ],
    "room_details": [...]
  }
}
```

---

## 分店管理

### 1. 获取分店列表

**接口地址：** `GET /api/v1/branches`

**查询参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| status | string | 否 | 状态筛选（ACTIVE/INACTIVE/ALL） |

**响应示例：**

```json
{
  "code": 200,
  "msg": "查询成功",
  "data": {
    "list": [
      {
        "id": 1,
        "hotel_name": "锦江之星",
        "branch_code": "BJ001",
        "address": "北京市朝阳区建国路88号",
        "contact": "张经理",
        "contact_phone": "13812345678",
        "status": "ACTIVE"
      }
    ],
    "total": 1
  }
}
```

---

### 2. 获取分店详情

**接口地址：** `GET /api/v1/branches/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 分店ID |

---

### 3. 创建分店

**接口地址：** `POST /api/v1/branches`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| hotel_name | string | 是 | 酒店名称 |
| branch_code | string | 否 | 分店编码（不提供则自动生成） |
| address | string | 是 | 地址 |
| contact | string | 是 | 联系人 |
| contact_phone | string | 是 | 联系电话 |
| created_by | uint64 | 是 | 创建人ID |

---

### 4. 更新分店

**接口地址：** `PUT /api/v1/branches/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 分店ID |

**请求参数：** 所有字段均为可选

---

### 5. 删除分店

**接口地址：** `DELETE /api/v1/branches/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 分店ID |

**业务规则：** 如果分店下有房源，不允许删除

---

## 订单管理

### 1. 获取订单列表

**接口地址：** `GET /api/v1/orders`

**查询参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认10 |
| branch_id | uint64 | 否 | 分店ID |
| guest_source | string | 否 | 客人来源 |
| order_no | string | 否 | 订单号 |
| phone | string | 否 | 手机号 |
| keyword | string | 否 | 关键词（订单号/房间号/手机号/联系人） |
| order_status | string | 否 | 订单状态 |
| check_in_start | string | 否 | 入住开始时间（YYYY-MM-DD） |
| check_in_end | string | 否 | 入住结束时间（YYYY-MM-DD） |
| check_out_start | string | 否 | 离店开始时间（YYYY-MM-DD） |
| check_out_end | string | 否 | 离店结束时间（YYYY-MM-DD） |
| reserve_start | string | 否 | 预定开始时间（YYYY-MM-DD HH:mm:ss） |
| reserve_end | string | 否 | 预定结束时间（YYYY-MM-DD HH:mm:ss） |

**响应说明：** 返回订单列表，包含关联的客人、房间、分店等信息

---

### 2. 获取订单详情

**接口地址：** `GET /api/v1/orders/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 订单ID |

**响应说明：** 返回订单详情，包含关联的房间号、房型、客人信息、财务信息等

---

## 客人管理

### 1. 获取在住客人列表

**接口地址：** `GET /api/v1/in-house-guests`

**查询参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| branch_id | uint64 | 否 | 分店ID |
| keyword | string | 否 | 关键词（姓名/手机号/房间号） |

**响应说明：** 返回当前在店客人的列表，包含客人信息、房间信息、订单信息

---

## 会员管理

### 1. 创建会员

**接口地址：** `POST /api/v1/members`

**请求参数：** 通过RPC调用，具体参数参考RPC接口定义

---

### 2. 获取会员列表

**接口地址：** `GET /api/v1/members`

**查询参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认10 |
| member_level | string | 否 | 会员等级筛选 |
| status | string | 否 | 状态筛选 |
| keyword | string | 否 | 关键词搜索 |

---

### 3. 获取会员详情

**接口地址：** `GET /api/v1/members/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 会员ID |

---

### 4. 根据客人ID获取会员信息

**接口地址：** `GET /api/v1/members/guest/:guest_id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| guest_id | int | 是 | 客人ID |

---

### 5. 更新会员

**接口地址：** `PUT /api/v1/members/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 会员ID |

---

### 6. 删除会员

**接口地址：** `DELETE /api/v1/members/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 会员ID |

---

### 7. 获取会员积分余额

**接口地址：** `GET /api/v1/members/:id/points-balance`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 会员ID |

---

### 8. 创建积分记录

**接口地址：** `POST /api/v1/points-records`

**请求参数：** 通过RPC调用

---

### 9. 获取积分记录列表

**接口地址：** `GET /api/v1/points-records`

**查询参数：** 通过RPC调用

---

### 10. 创建会员权益

**接口地址：** `POST /api/v1/member-rights`

**请求参数：** 通过RPC调用

---

### 11. 获取会员权益列表

**接口地址：** `GET /api/v1/member-rights`

**查询参数：** 通过RPC调用

---

### 12. 根据会员等级获取权益

**接口地址：** `GET /api/v1/member-rights/level/:member_level`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| member_level | string | 是 | 会员等级 |

---

## 财务管理

### 1. 获取收支流水列表

**接口地址：** `GET /api/v1/financial-flows`

**查询参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| branch_id | uint64 | 否 | 分店ID |
| flow_type | string | 否 | 收支类型（INCOME/EXPENSE） |
| flow_item | string | 否 | 收支项目 |
| pay_type | string | 否 | 支付方式 |
| start_date | string | 否 | 开始日期 |
| end_date | string | 否 | 结束日期 |

**响应说明：** 返回财务流水列表，支持按日期范围、收支类型、支付方式等维度统计

---

## 系统管理

### 1. 创建用户账号

**接口地址：** `POST /api/v1/user-accounts`

**请求参数：** 通过RPC调用

---

### 2. 获取用户账号列表

**接口地址：** `GET /api/v1/user-accounts`

**查询参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| role_id | uint64 | 否 | 角色ID筛选 |
| branch_id | uint64 | 否 | 分店ID筛选 |
| status | string | 否 | 状态筛选 |
| keyword | string | 否 | 关键词搜索 |

---

### 3. 获取用户账号详情

**接口地址：** `GET /api/v1/user-accounts/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 账号ID |

---

### 4. 更新用户账号

**接口地址：** `PUT /api/v1/user-accounts/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 账号ID |

---

### 5. 删除用户账号

**接口地址：** `DELETE /api/v1/user-accounts/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 账号ID |

---

### 6. 创建角色

**接口地址：** `POST /api/v1/roles`

**请求参数：** 通过RPC调用

---

### 7. 获取角色列表

**接口地址：** `GET /api/v1/roles`

**查询参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| status | string | 否 | 状态筛选 |
| keyword | string | 否 | 关键词搜索 |

---

### 8. 获取角色详情

**接口地址：** `GET /api/v1/roles/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 角色ID |

---

### 9. 更新角色

**接口地址：** `PUT /api/v1/roles/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 角色ID |

---

### 10. 删除角色

**接口地址：** `DELETE /api/v1/roles/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 角色ID |

---

### 11. 获取权限列表

**接口地址：** `GET /api/v1/permissions`

**查询参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| permission_type | string | 否 | 权限类型筛选 |
| parent_id | uint64 | 否 | 父权限ID筛选 |
| status | string | 否 | 状态筛选 |

**响应说明：** 返回权限树形结构或扁平列表

---

### 12. 创建渠道配置

**接口地址：** `POST /api/v1/channel-configs`

**请求参数：** 通过Service层调用

---

### 13. 获取渠道配置列表

**接口地址：** `GET /api/v1/channel-configs`

**查询参数：** 通过Service层调用

---

### 14. 获取渠道配置详情

**接口地址：** `GET /api/v1/channel-configs/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 配置ID |

---

### 15. 更新渠道配置

**接口地址：** `PUT /api/v1/channel-configs/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 配置ID |

---

### 16. 删除渠道配置

**接口地址：** `DELETE /api/v1/channel-configs/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 配置ID |

---

### 17. 创建系统配置

**接口地址：** `POST /api/v1/system-configs`

**请求参数：** 通过Service层调用

---

### 18. 获取系统配置列表

**接口地址：** `GET /api/v1/system-configs`

**查询参数：** 通过Service层调用

---

### 19. 获取系统配置详情

**接口地址：** `GET /api/v1/system-configs/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 配置ID |

---

### 20. 更新系统配置

**接口地址：** `PUT /api/v1/system-configs/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 配置ID |

---

### 21. 删除系统配置

**接口地址：** `DELETE /api/v1/system-configs/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 配置ID |

---

### 22. 根据分类获取系统配置

**接口地址：** `GET /api/v1/system-configs/category/:category`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| category | string | 是 | 配置分类 |

---

### 23. 创建黑名单

**接口地址：** `POST /api/v1/blacklists`

**请求参数：** 通过Service层调用

---

### 24. 获取黑名单列表

**接口地址：** `GET /api/v1/blacklists`

**查询参数：** 通过Service层调用

---

### 25. 获取黑名单详情

**接口地址：** `GET /api/v1/blacklists/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 黑名单ID |

---

### 26. 更新黑名单

**接口地址：** `PUT /api/v1/blacklists/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 黑名单ID |

---

### 27. 删除黑名单

**接口地址：** `DELETE /api/v1/blacklists/:id`

**路径参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 黑名单ID |

---

### 28. 创建操作日志

**接口地址：** `POST /api/v1/operation-logs`

**请求参数：** 通过Service层调用

---

### 29. 获取操作日志列表

**接口地址：** `GET /api/v1/operation-logs`

**查询参数：** 通过Service层调用

---

## 渠道同步

### 1. 同步房态到渠道

**接口地址：** `POST /api/v1/sync-room-status`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| branch_id | uint64 | 是 | 分店ID |
| channel_id | uint64 | 是 | 渠道ID |
| start_date | string | 是 | 开始日期（YYYY-MM-DD） |
| end_date | string | 是 | 结束日期（YYYY-MM-DD） |
| room_ids | []uint64 | 否 | 房源ID列表（可选，不提供则同步该分店所有房源） |

**请求示例：**

```json
{
  "branch_id": 1,
  "channel_id": 1,
  "start_date": "2025-01-15",
  "end_date": "2025-01-21",
  "room_ids": [1, 2, 3]
}
```

**响应示例：**

```json
{
  "code": 200,
  "msg": "同步完成",
  "data": {
    "success_count": 50,
    "failed_count": 0
  }
}
```

---

## 健康检查

### 1. 健康检查

**接口地址：** `GET /api/v1/health` 或 `GET /health`

**响应示例：**

```json
{
  "status": "ok"
}
```

---

## 接口调用示例

### cURL示例

```bash
# 获取房型列表
curl -X GET "http://localhost:3000/api/v1/room-types?page=1&page_size=10" \
  -H "Content-Type: application/json"

# 创建房源
curl -X POST "http://localhost:3000/api/v1/room-infos" \
  -H "Content-Type: application/json" \
  -d '{
    "branch_id": 1,
    "room_type_id": 1,
    "room_no": "101",
    "room_name": "舒适大床房",
    "market_price": 298.00,
    "calendar_price": 258.00,
    "room_count": 1,
    "bed_spec": "1张大床",
    "created_by": 1
  }'

# 获取订单列表
curl -X GET "http://localhost:3000/api/v1/orders?page=1&page_size=10&branch_id=1" \
  -H "Content-Type: application/json"
```

### JavaScript示例

```javascript
// 获取房型列表
fetch('http://localhost:3000/api/v1/room-types?page=1&page_size=10', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json'
  }
})
  .then(response => response.json())
  .then(data => console.log(data));

// 创建房源
fetch('http://localhost:3000/api/v1/room-infos', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    branch_id: 1,
    room_type_id: 1,
    room_no: '101',
    room_name: '舒适大床房',
    market_price: 298.00,
    calendar_price: 258.00,
    room_count: 1,
    bed_spec: '1张大床',
    created_by: 1
  })
})
  .then(response => response.json())
  .then(data => console.log(data));
```

---

## 注意事项

1. **参数验证：** 所有必填参数必须提供，否则返回400错误
2. **数据格式：** 请求和响应均使用JSON格式
3. **分页参数：** 列表接口支持分页，默认page=1，page_size=10
4. **软删除：** 删除操作均为软删除，数据不会物理删除
5. **权限控制：** 部分接口可能需要权限验证（当前文档未包含认证相关接口）
6. **RPC接口：** 部分接口通过RPC调用实现，具体参数参考RPC接口定义
7. **文件上传：** 图片上传接口使用 `multipart/form-data` 格式

---

## 更新日志

| 日期 | 版本 | 说明 |
|------|------|------|
| 2025-01-16 | v1.0 | 初始版本，包含所有API接口说明 |

---

**文档维护：** 本文档基于代码自动生成，如接口有变更请及时更新文档。