# HTTP API 服务

## 项目结构

```
api/
├── handler/
│   └── room_handler.go    # 房源管理 API 处理器
├── router/
│   └── router.go          # 路由配置
└── main/
    └── main.go            # API 服务启动入口
```

## 启动服务

```bash
cd api/main
go run main.go
```

服务将在 `:8080` 端口启动。

## API 接口文档

### 基础信息

- **Base URL**: `http://localhost:8080`
- **API 版本**: `v1`
- **Content-Type**: `application/json`

### 房型管理接口

#### 1. 创建房型
- **接口**: `POST /api/v1/room-types`
- **请求体**:
```json
{
  "room_type_name": "大床房",
  "bed_spec": "1.8*2.0m",
  "area": 25.5,
  "has_breakfast": true,
  "has_toiletries": true,
  "default_price": 299.00
}
```

#### 2. 获取房型列表
- **接口**: `GET /api/v1/room-types`
- **查询参数**:
  - `page`: 页码（默认: 1）
  - `page_size`: 每页数量（默认: 10）
  - `status`: 状态筛选（可选）
  - `keyword`: 关键词搜索（可选）

#### 3. 获取房型详情
- **接口**: `GET /api/v1/room-types/:id`

#### 4. 更新房型
- **接口**: `PUT /api/v1/room-types/:id`
- **请求体**: 同创建房型（字段可选）

#### 5. 删除房型
- **接口**: `DELETE /api/v1/room-types/:id`

### 房源管理接口

#### 1. 创建房源
- **接口**: `POST /api/v1/room-infos`
- **请求体**:
```json
{
  "branch_id": 1,
  "room_type_id": 1,
  "room_no": "101",
  "room_name": "标准大床房-101",
  "market_price": 299.00,
  "calendar_price": 269.00,
  "room_count": 1,
  "area": 25.5,
  "bed_spec": "1.8*2.0m",
  "has_breakfast": true,
  "has_toiletries": true,
  "cancellation_policy_id": 1,
  "created_by": 1
}
```

#### 2. 获取房源列表
- **接口**: `GET /api/v1/room-infos`
- **查询参数**:
  - `page`: 页码（默认: 1）
  - `page_size`: 每页数量（默认: 10）
  - `branch_id`: 分店ID（可选）
  - `room_type_id`: 房型ID（可选）
  - `status`: 状态筛选（可选）
  - `keyword`: 关键词搜索（可选）

#### 3. 获取房源详情
- **接口**: `GET /api/v1/room-infos/:id`

#### 4. 更新房源
- **接口**: `PUT /api/v1/room-infos/:id`
- **请求体**: 同创建房源（字段可选）

#### 5. 删除房源
- **接口**: `DELETE /api/v1/room-infos/:id`

#### 6. 更新房源状态
- **接口**: `PUT /api/v1/room-infos/:id/status`
- **请求体**:
```json
{
  "status": "ACTIVE"  // ACTIVE-启用, INACTIVE-停用, MAINTENANCE-维修
}
```

#### 7. 批量更新房源状态
- **接口**: `PUT /api/v1/room-infos/batch-status`
- **请求体**:
```json
{
  "room_ids": [1, 2, 3],
  "status": "ACTIVE"
}
```

#### 8. 创建关联房绑定
- **接口**: `POST /api/v1/room-infos/bindings`
- **请求体**:
```json
{
  "main_room_id": 1,
  "related_room_id": 2,
  "binding_desc": "关联描述（可选）"
}
```

#### 9. 批量创建关联房绑定
- **接口**: `POST /api/v1/room-infos/batch-bindings`
- **请求体**:
```json
{
  "main_room_id": 1,
  "related_room_ids": [2, 3, 4],
  "binding_desc": "关联描述（可选）"
}
```

#### 10. 获取关联房列表
- **接口**: `GET /api/v1/room-infos/:id/bindings`

#### 11. 删除关联房绑定
- **接口**: `DELETE /api/v1/room-infos/bindings/:id`

#### 12. 批量上传房源图片
- **接口**: `POST /api/v1/room-infos/:id/images`
- **Content-Type**: `multipart/form-data`
- **参数**: `images` (文件数组，最多16张)
- **限制**:
  - 图片格式：jpg/png
  - 图片规格：自动调整为 400x300
  - 单张图片最大：5MB
  - 每个房源最多：16张

#### 13. 获取房源图片列表
- **接口**: `GET /api/v1/room-infos/:id/images`

#### 14. 删除房源图片
- **接口**: `DELETE /api/v1/room-infos/images/:id`

#### 15. 更新图片排序
- **接口**: `PUT /api/v1/room-infos/images/:id/sort`
- **请求体**:
```json
{
  "sort_order": 1
}
```

#### 16. 批量更新图片排序
- **接口**: `PUT /api/v1/room-infos/:id/images/sort`
- **请求体**:
```json
{
  "sort_orders": [
    {"image_id": 1, "sort_order": 0},
    {"image_id": 2, "sort_order": 1}
  ]
}
```

### 健康检查

- **接口**: `GET /health`
- **响应**: `{"status": "ok"}`

## 响应格式

### 成功响应
```json
{
  "code": 200,
  "msg": "操作成功",
  "data": {...}
}
```

### 错误响应
```json
{
  "code": 400,
  "msg": "错误信息"
}
```

## 测试示例

### 使用 curl 测试

```bash
# 创建房型
curl -X POST http://localhost:8080/api/v1/room-types \
  -H "Content-Type: application/json" \
  -d '{
    "room_type_name": "大床房",
    "bed_spec": "1.8*2.0m",
    "area": 25.5,
    "has_breakfast": true,
    "has_toiletries": true,
    "default_price": 299.00
  }'

# 获取房型列表
curl http://localhost:8080/api/v1/room-types?page=1&page_size=10

# 创建房源
curl -X POST http://localhost:8080/api/v1/room-infos \
  -H "Content-Type: application/json" \
  -d '{
    "branch_id": 1,
    "room_type_id": 1,
    "room_no": "101",
    "room_name": "标准大床房-101",
    "market_price": 299.00,
    "calendar_price": 269.00,
    "room_count": 1,
    "area": 25.5,
    "bed_spec": "1.8*2.0m",
    "has_breakfast": true,
    "has_toiletries": true,
    "created_by": 1
  }'

# 更新房源状态
curl -X PUT http://localhost:8080/api/v1/room-infos/1/status \
  -H "Content-Type: application/json" \
  -d '{"status": "ACTIVE"}'

# 创建关联房绑定
curl -X POST http://localhost:8080/api/v1/room-infos/bindings \
  -H "Content-Type: application/json" \
  -d '{
    "main_room_id": 1,
    "related_room_id": 2,
    "binding_desc": "关联描述"
  }'

# 上传房源图片（使用 multipart/form-data）
curl -X POST http://localhost:8080/api/v1/room-infos/1/images \
  -F "images=@/path/to/image1.jpg" \
  -F "images=@/path/to/image2.png"
```

## 注意事项

1. 所有接口都需要数据库连接正常
2. 创建房源时，需要先确保分店和房型存在
3. 删除操作会进行业务规则检查（如是否被使用）
4. 分页参数如果不提供，会使用默认值
