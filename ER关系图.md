# 核心模块ER关系图

本文档展示了房源管理、房态管理、订单管理、客人管理四个核心模块的数据库ER关系图。

---

## ER关系图

```mermaid
erDiagram
    hotel_branch ||--o{ room_info : "一个分店有多个房源"
    room_type_dict ||--o{ room_info : "一个房型有多个房源"
    room_info ||--o{ room_status_detail : "一个房源有多个房态明细"
    room_info ||--o{ room_image : "一个房源有多个图片"
    room_info ||--o{ order_main : "一个房源有多个订单"
    hotel_branch ||--o{ order_main : "一个分店有多个订单"
    guest_info ||--o{ order_main : "一个客人有多个订单"
    room_type_dict ||--o{ order_main : "一个房型有多个订单"
    order_main ||--|| order_extension : "一个订单有一个扩展信息"

    hotel_branch {
        int id PK "分店ID"
        string hotel_name "酒店名称"
        string branch_code UK "分店编码"
        string address "地址"
        string contact "联系人"
        string contact_phone "联系电话"
        string status "状态"
        datetime created_at "创建时间"
        datetime updated_at "更新时间"
        datetime deleted_at "软删除时间"
    }

    room_type_dict {
        int id PK "房型ID"
        string room_type_name "房型名称"
        string bed_spec "床型规格"
        decimal area "面积"
        boolean has_breakfast "是否含早"
        boolean has_toiletries "是否提供洗漱用品"
        decimal default_price "默认价格"
        string status "状态"
        datetime created_at "创建时间"
        datetime updated_at "更新时间"
        datetime deleted_at "软删除时间"
    }

    room_info {
        int id PK "房源ID"
        int branch_id FK "分店ID"
        int room_type_id FK "房型ID"
        string room_no "房间号"
        string room_name "房间名称"
        decimal market_price "门市价"
        decimal calendar_price "日历价"
        int room_count "房间数量"
        decimal area "面积"
        string bed_spec "床型规格"
        boolean has_breakfast "是否含早"
        boolean has_toiletries "是否提供洗漱用品"
        int cancellation_policy_id FK "退订政策ID"
        string status "状态"
        datetime created_at "创建时间"
        datetime updated_at "更新时间"
        datetime deleted_at "软删除时间"
    }

    room_status_detail {
        int id PK "记录ID"
        int room_id FK "房源ID"
        date date "日期"
        string room_status "房态"
        int remaining_count "剩余数量"
        int checked_in_count "已入住人数"
        int check_out_pending_count "预退房人数"
        int reserved_pending_count "预定待入住人数"
        datetime created_at "创建时间"
        datetime updated_at "更新时间"
        datetime deleted_at "软删除时间"
    }

    room_image {
        int id PK "图片ID"
        int room_id FK "房源ID"
        string image_url "图片URL"
        string image_size "图片规格"
        string image_format "图片格式"
        int sort_order "排序序号"
        datetime upload_time "上传时间"
        datetime created_at "创建时间"
        datetime deleted_at "软删除时间"
    }

    guest_info {
        int id PK "客人ID"
        string name "姓名"
        string id_type "证件类型"
        string id_number "证件号(加密)"
        string phone "手机号(加密)"
        string gender "性别"
        string ethnicity "民族"
        string province "省份"
        string address "地址"
        datetime check_in_time "入住时间"
        datetime check_out_time "离店时间"
        int room_id FK "房间ID"
        int order_id FK "订单ID"
        int register_by "登记人ID"
        datetime register_time "登记时间"
        boolean is_member "是否会员"
        int member_id FK "会员ID"
        datetime created_at "创建时间"
        datetime updated_at "更新时间"
        datetime deleted_at "软删除时间"
    }

    order_main {
        int id PK "订单ID"
        string order_no UK "订单号"
        int branch_id FK "分店ID"
        int guest_id FK "客人ID"
        int room_id FK "房源ID"
        int room_type_id FK "房型ID"
        string guest_source "客人来源"
        datetime check_in_time "入住时间"
        datetime check_out_time "离店时间"
        datetime reserve_time "预定时间"
        decimal order_amount "订单金额"
        decimal deposit_received "已收押金"
        decimal outstanding_amount "未付金额"
        string order_status "订单状态"
        string pay_type "支付方式"
        decimal penalty_amount "违约金"
        datetime created_at "创建时间"
        datetime updated_at "更新时间"
        datetime deleted_at "软删除时间"
    }

    order_extension {
        int id PK "扩展ID"
        int order_id FK UK "订单ID"
        string contact "联系人"
        string contact_phone "联系电话"
        string special_request "特殊需求"
        int guest_count "入住人数"
        int room_count "房间数量"
        string sync_status "同步状态"
        datetime sync_time "同步时间"
        datetime created_at "创建时间"
        datetime updated_at "更新时间"
        datetime deleted_at "软删除时间"
    }
```

---

## 关系说明

### 1. 分店与房源 (hotel_branch → room_info)
- **关系类型**：一对多 (1:N)
- **说明**：一个分店可以有多个房源，每个房源属于一个分店
- **外键**：`room_info.branch_id` → `hotel_branch.id`

### 2. 房型与房源 (room_type_dict → room_info)
- **关系类型**：一对多 (1:N)
- **说明**：一个房型可以有多个房源，每个房源对应一个房型
- **外键**：`room_info.room_type_id` → `room_type_dict.id`

### 3. 房源与房态明细 (room_info → room_status_detail)
- **关系类型**：一对多 (1:N)
- **说明**：一个房源可以有多个日期的房态明细记录
- **外键**：`room_status_detail.room_id` → `room_info.id`

### 4. 房源与图片 (room_info → room_image)
- **关系类型**：一对多 (1:N)
- **说明**：一个房源可以有多个图片
- **外键**：`room_image.room_id` → `room_info.id`

### 5. 房源与订单 (room_info → order_main)
- **关系类型**：一对多 (1:N)
- **说明**：一个房源可以有多个订单（不同时间段的预订）
- **外键**：`order_main.room_id` → `room_info.id`

### 6. 分店与订单 (hotel_branch → order_main)
- **关系类型**：一对多 (1:N)
- **说明**：一个分店可以有多个订单
- **外键**：`order_main.branch_id` → `hotel_branch.id`

### 7. 客人与订单 (guest_info → order_main)
- **关系类型**：一对多 (1:N)
- **说明**：一个客人可以有多个订单（多次入住记录）
- **外键**：`order_main.guest_id` → `guest_info.id`

### 8. 房型与订单 (room_type_dict → order_main)
- **关系类型**：一对多 (1:N)
- **说明**：一个房型可以有多个订单（冗余字段，方便统计）
- **外键**：`order_main.room_type_id` → `room_type_dict.id`

### 9. 订单与订单扩展 (order_main → order_extension)
- **关系类型**：一对一 (1:1)
- **说明**：一个订单有一个扩展信息（联系人、特殊需求等）
- **外键**：`order_extension.order_id` → `order_main.id`（唯一索引）

---

## 核心业务实体

### 房源管理模块
- **hotel_branch**：分店信息，房源的组织单位
- **room_type_dict**：房型字典，房源的分类标准
- **room_info**：房源信息，核心业务实体
- **room_image**：房源图片，房源的展示信息

### 房态管理模块
- **room_status_detail**：房态明细，按日期记录房态信息

### 订单管理模块
- **order_main**：订单主表，核心业务实体
- **order_extension**：订单扩展表，订单的补充信息

### 客人管理模块
- **guest_info**：客人信息，实名制登记信息

---

## 数据流向

### 创建订单流程
1. 选择分店 (`hotel_branch`)
2. 选择房型 (`room_type_dict`)
3. 选择房源 (`room_info`)
4. 创建客人信息 (`guest_info`)
5. 创建订单 (`order_main`)
6. 创建订单扩展 (`order_extension`)

### 查询在住客人流程
1. 查询订单 (`order_main`) - 状态为已入住
2. 关联客人 (`guest_info`) - 获取客人信息
3. 关联房源 (`room_info`) - 获取房间信息
4. 关联房型 (`room_type_dict`) - 获取房型信息

### 房态查询流程
1. 查询房态明细 (`room_status_detail`) - 按日期筛选
2. 关联房源 (`room_info`) - 获取房间信息
3. 关联分店 (`hotel_branch`) - 获取分店信息

---

## 索引说明

### 主要索引
- **分店ID索引**：`room_info.branch_id`、`order_main.branch_id`
- **房型ID索引**：`room_info.room_type_id`、`order_main.room_type_id`
- **客人ID索引**：`order_main.guest_id`
- **房源ID索引**：`room_status_detail.room_id`、`room_image.room_id`、`order_main.room_id`
- **订单号唯一索引**：`order_main.order_no`
- **日期索引**：`room_status_detail.date`
- **房态索引**：`room_status_detail.room_status`

### 唯一约束
- **分店编码**：`hotel_branch.branch_code` (UK)
- **订单号**：`order_main.order_no` (UK)
- **订单扩展**：`order_extension.order_id` (UK)

---

*注：本ER图仅包含房源管理、房态管理、订单管理、客人管理四个核心模块相关的数据表。*
