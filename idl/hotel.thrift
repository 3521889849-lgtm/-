namespace go hotel

// 基础响应体
struct BaseResp {
    1: i32 code,
    2: string msg
}

// 房型字典
struct RoomType {
    1: i64 id,
    2: string room_type_name,
    3: string bed_spec,
    4: optional double area,
    5: bool has_breakfast,
    6: bool has_toiletries,
    7: double default_price,
    8: string status,
    9: string created_at,
    10: string updated_at
}

// 创建房型请求
struct CreateRoomTypeReq {
    1: string room_type_name,
    2: string bed_spec,
    3: optional double area,
    4: bool has_breakfast,
    5: bool has_toiletries,
    6: double default_price
}

// 更新房型请求
struct UpdateRoomTypeReq {
    1: i64 id,
    2: optional string room_type_name,
    3: optional string bed_spec,
    4: optional double area,
    5: optional bool has_breakfast,
    6: optional bool has_toiletries,
    7: optional double default_price,
    8: optional string status
}

// 房型列表请求
struct ListRoomTypeReq {
    1: i32 page = 1,
    2: i32 page_size = 10,
    3: optional string status,
    4: optional string keyword
}

// 房型列表响应
struct ListRoomTypeResp {
    1: list<RoomType> list,
    2: i64 total,
    3: i32 page,
    4: i32 page_size
}

// 房源信息
struct RoomInfo {
    1: i64 id,
    2: i64 branch_id,
    3: i64 room_type_id,
    4: string room_no,
    5: string room_name,
    6: double market_price,
    7: double calendar_price,
    8: i8 room_count,
    9: optional double area,
    10: string bed_spec,
    11: bool has_breakfast,
    12: bool has_toiletries,
    13: optional i64 cancellation_policy_id,
    14: string status,
    15: string created_at,
    16: string updated_at
}

// 创建房源请求
struct CreateRoomInfoReq {
    1: i64 branch_id,
    2: i64 room_type_id,
    3: string room_no,
    4: string room_name,
    5: double market_price,
    6: double calendar_price,
    7: i8 room_count,
    8: optional double area,
    9: string bed_spec,
    10: bool has_breakfast,
    11: bool has_toiletries,
    12: optional i64 cancellation_policy_id,
    13: i64 created_by
}

// 更新房源请求
struct UpdateRoomInfoReq {
    1: i64 id,
    2: optional string room_no,
    3: optional string room_name,
    4: optional double market_price,
    5: optional double calendar_price,
    6: optional i8 room_count,
    7: optional double area,
    8: optional string bed_spec,
    9: optional bool has_breakfast,
    10: optional bool has_toiletries,
    11: optional i64 cancellation_policy_id,
    12: optional string status
}

// 房源列表请求
struct ListRoomInfoReq {
    1: i32 page = 1,
    2: i32 page_size = 10,
    3: optional i64 branch_id,
    4: optional i64 room_type_id,
    5: optional string status,
    6: optional string keyword
}

// 房源列表响应
struct ListRoomInfoResp {
    1: list<RoomInfo> list,
    2: i64 total,
    3: i32 page,
    4: i32 page_size
}

// 更新房源状态请求
struct UpdateRoomStatusReq {
    1: i64 room_id,
    2: string status  // ACTIVE-启用, INACTIVE-停用, MAINTENANCE-维修
}

// 批量更新房源状态请求
struct BatchUpdateRoomStatusReq {
    1: list<i64> room_ids,
    2: string status
}

// 关联房绑定
struct RelatedRoomBinding {
    1: i64 id,
    2: i64 main_room_id,
    3: i64 related_room_id,
    4: optional string binding_desc,
    5: string created_at
}

// 创建关联房绑定请求
struct CreateRoomBindingReq {
    1: i64 main_room_id,
    2: i64 related_room_id,
    3: optional string binding_desc
}

// 批量创建关联房绑定请求
struct BatchCreateRoomBindingsReq {
    1: i64 main_room_id,
    2: list<i64> related_room_ids,
    3: optional string binding_desc
}

// 关联房列表响应
struct ListRoomBindingsResp {
    1: list<RelatedRoomBinding> bindings
}

// 房源图片
struct RoomImage {
    1: i64 id,
    2: i64 room_id,
    3: string image_url,
    4: string image_size,
    5: string image_format,
    6: i8 sort_order,
    7: string upload_time
}

// 上传图片响应
struct UploadRoomImagesResp {
    1: list<RoomImage> images
}

// 图片列表响应
struct ListRoomImagesResp {
    1: list<RoomImage> images
}

// 更新图片排序请求
struct UpdateImageSortOrderReq {
    1: i64 image_id,
    2: i8 sort_order
}

// 批量更新图片排序请求
struct BatchUpdateImageSortOrderReq {
    1: i64 room_id,
    2: list<ImageSortOrder> sort_orders
}

// 图片排序项
struct ImageSortOrder {
    1: i64 image_id,
    2: i8 sort_order
}

// 设施字典
struct Facility {
    1: i64 id,
    2: string facility_name,
    3: optional string description,
    4: string status,
    5: string created_at,
    6: string updated_at
}

// 创建设施请求
struct CreateFacilityReq {
    1: string facility_name,
    2: optional string description
}

// 更新设施请求
struct UpdateFacilityReq {
    1: i64 id,
    2: optional string facility_name,
    3: optional string description,
    4: optional string status
}

// 设施列表请求
struct ListFacilityReq {
    1: i32 page = 1,
    2: i32 page_size = 10,
    3: optional string status,
    4: optional string keyword
}

// 设施列表响应
struct ListFacilitiesResp {
    1: list<Facility> list,
    2: i64 total,
    3: i32 page,
    4: i32 page_size
}

// 设置房源设施请求
struct SetRoomFacilitiesReq {
    1: i64 room_id,
    2: list<i64> facility_ids
}

// 添加房源设施请求
struct AddRoomFacilityReq {
    1: i64 room_id,
    2: i64 facility_id
}

// 移除房源设施请求
struct RemoveRoomFacilityReq {
    1: i64 room_id,
    2: i64 facility_id
}

// 退订政策
struct CancellationPolicy {
    1: i64 id,
    2: string policy_name,
    3: string rule_description,
    4: double penalty_ratio,
    5: optional i64 room_type_id,
    6: string status,
    7: string created_at,
    8: string updated_at
}

// 创建退订政策请求
struct CreateCancellationPolicyReq {
    1: string policy_name,
    2: string rule_description,
    3: double penalty_ratio,
    4: optional i64 room_type_id
}

// 更新退订政策请求
struct UpdateCancellationPolicyReq {
    1: i64 id,
    2: optional string policy_name,
    3: optional string rule_description,
    4: optional double penalty_ratio,
    5: optional i64 room_type_id,
    6: optional string status
}

// 退订政策列表请求
struct ListCancellationPolicyReq {
    1: i32 page = 1,
    2: i32 page_size = 10,
    3: optional i64 room_type_id,
    4: optional string status,
    5: optional string keyword
}

// 退订政策列表响应
struct ListCancellationPoliciesResp {
    1: list<CancellationPolicy> list,
    2: i64 total,
    3: i32 page,
    4: i32 page_size
}

// 日历化房态项
struct CalendarRoomStatusItem {
    1: i64 room_id,
    2: string room_no,
    3: string room_name,
    4: string date,  // YYYY-MM-DD
    5: string room_status,  // 空净房/入住房/维修房/锁定房/空账房/预定房
    6: i8 remaining_count,
    7: i8 checked_in_count,
    8: i8 check_out_pending_count,
    9: i8 reserved_pending_count
}

// 日历化房态查询请求
struct CalendarRoomStatusReq {
    1: optional i64 branch_id,
    2: string start_date,  // YYYY-MM-DD
    3: string end_date,     // YYYY-MM-DD
    4: optional string room_no,
    5: optional string status  // 房态筛选
}

// 日历化房态查询响应
struct CalendarRoomStatusResp {
    1: list<CalendarRoomStatusItem> items
}

// 更新日历化房态请求
struct UpdateCalendarRoomStatusReq {
    1: i64 room_id,
    2: string date,  // YYYY-MM-DD
    3: string status  // 空净房/入住房/维修房/锁定房/空账房/预定房
}

// 批量更新日历化房态请求
struct BatchUpdateCalendarRoomStatusReq {
    1: list<UpdateCalendarRoomStatusReq> updates
}

// 实时数据统计请求
struct RealTimeStatisticsReq {
    1: optional i64 branch_id,
    2: optional string date,  // YYYY-MM-DD，默认为今日
    3: optional string room_no,
    4: optional i64 room_type_id
}

// 房态分组统计
struct StatusBreakdown {
    1: string status,
    2: i64 count
}

// 房间明细统计
struct RoomDetailStat {
    1: i64 room_id,
    2: string room_no,
    3: string room_name,
    4: string room_status,
    5: i8 remaining_count,
    6: i8 checked_in_count,
    7: i8 check_out_pending_count,
    8: i8 reserved_pending_count
}

// 实时数据统计响应
struct RealTimeStatisticsResp {
    1: string date,
    2: i64 total_rooms,
    3: i64 remaining_rooms,
    4: i64 checked_in_count,
    5: i64 check_out_pending_count,
    6: i64 reserved_pending_count,
    7: i64 occupied_rooms,
    8: i64 maintenance_rooms,
    9: i64 locked_rooms,
    10: i64 empty_rooms,
    11: i64 reserved_rooms,
    12: list<StatusBreakdown> status_breakdown,
    13: optional list<RoomDetailStat> room_details
}

// 分店信息
struct Branch {
    1: i64 id,
    2: string hotel_name,
    3: string branch_code,
    4: string address,
    5: string contact,
    6: string contact_phone,
    7: string status,
    8: string created_at,
    9: string updated_at
}

// 分店列表请求
struct ListBranchesReq {
    1: optional string status
}

// 分店列表响应
struct ListBranchesResp {
    1: list<Branch> branches
}

// 同步房态到渠道请求
struct SyncRoomStatusToChannelReq {
    1: i64 branch_id,
    2: i64 channel_id,
    3: string start_date,  // YYYY-MM-DD
    4: string end_date,    // YYYY-MM-DD
    5: optional list<i64> room_ids
}

// 同步房态到渠道响应
struct SyncRoomStatusToChannelResp {
    1: i32 success_count,
    2: i32 fail_count,
    3: list<i64> sync_logs
}

// 订单信息
struct Order {
    1: i64 id,
    2: string order_no,
    3: i64 branch_id,
    4: optional string branch_name,
    5: i64 guest_id,
    6: optional string guest_name,
    7: i64 room_id,
    8: optional string room_no,
    9: optional string room_name,
    10: i64 room_type_id,
    11: optional string room_type_name,
    12: string guest_source,
    13: string check_in_time,
    14: string check_out_time,
    15: string reserve_time,
    16: double order_amount,
    17: double deposit_received,
    18: double outstanding_amount,
    19: string order_status,
    20: string pay_type,
    21: double penalty_amount,
    22: string created_at,
    23: string updated_at,
    24: optional string contact,
    25: optional string contact_phone,
    26: optional string special_request,
    27: optional i8 guest_count,
    28: optional i8 room_count,
    29: optional list<string> room_nos
}

// 订单列表查询请求
struct ListOrdersReq {
    1: i32 page = 1,
    2: i32 page_size = 10,
    3: optional i64 branch_id,
    4: optional string guest_source,
    5: optional string order_no,
    6: optional string phone,
    7: optional string keyword,
    8: optional string order_status,
    9: optional string check_in_start,
    10: optional string check_in_end,
    11: optional string check_out_start,
    12: optional string check_out_end,
    13: optional string reserve_start,
    14: optional string reserve_end
}

// 订单列表响应
struct ListOrdersResp {
    1: list<Order> list,
    2: i64 total,
    3: i32 page,
    4: i32 page_size
}

// 在住客人信息
struct InHouseGuest {
    1: i64 id,
    2: i64 guest_id,
    3: string name,
    4: string id_type,
    5: string id_number,
    6: string phone,
    7: optional string province,
    8: optional string address,
    9: optional string ethnicity,
    10: string check_in_time,
    11: string check_out_time,
    12: i64 order_id,
    13: string order_no,
    14: string guest_source,
    15: i64 room_id,
    16: string room_no,
    17: i64 room_type_id,
    18: string room_type_name,
    19: double order_amount,
    20: double deposit_received,
    21: double outstanding_amount
}

// 在住客人列表查询请求
struct ListInHouseGuestsReq {
    1: i32 page = 1,
    2: i32 page_size = 200,
    3: optional i64 branch_id,
    4: optional string province,
    5: optional string city,
    6: optional string district,
    7: optional string name,
    8: optional string phone,
    9: optional string id_number,
    10: optional string room_no
}

// 在住客人列表响应
struct ListInHouseGuestsResp {
    1: list<InHouseGuest> list,
    2: i64 total,
    3: i32 page,
    4: i32 page_size
}

// 财务流水信息
struct FinancialFlow {
    1: i64 id,
    2: optional i64 order_id,
    3: i64 branch_id,
    4: optional i64 room_id,
    5: optional i64 guest_id,
    6: string flow_type,
    7: string flow_item,
    8: string pay_type,
    9: double amount,
    10: string occur_time,
    11: i64 operator_id,
    12: optional string remark,
    13: string created_at,
    14: optional string room_no,
    15: optional string guest_name,
    16: optional string contact_phone,
    17: optional string operator_name,
    18: optional string order_no
}

// 财务汇总（按支付方式）
struct FinancialSummary {
    1: double total,
    2: double cash,
    3: double alipay,
    4: double wechat,
    5: double unionpay,
    6: double card_swipe,
    7: double tuyou_collection,
    8: double ctrip_collection,
    9: double qunar_collection
}

// 财务汇总响应
struct FinancialSummaryResp {
    1: FinancialSummary income,
    2: FinancialSummary expense,
    3: FinancialSummary balance
}

// 收支流水列表查询请求
struct ListFinancialFlowsReq {
    1: i32 page = 1,
    2: i32 page_size = 200,
    3: optional i64 branch_id,
    4: optional string flow_type,
    5: optional string flow_item,
    6: optional string pay_type,
    7: optional i64 operator_id,
    8: optional string occur_start,
    9: optional string occur_end
}

// 收支流水列表响应
struct ListFinancialFlowsResp {
    1: list<FinancialFlow> list,
    2: i64 total,
    3: i32 page,
    4: i32 page_size,
    5: FinancialSummaryResp summary
}

// 账号信息
struct UserAccount {
    1: i64 id,
    2: string username,
    3: string real_name,
    4: string contact_phone,
    5: i64 role_id,
    6: optional string role_name,
    7: optional i64 branch_id,
    8: optional string branch_name,
    9: string status,
    10: string created_at,
    11: optional string last_login_at
}

// 创建账号请求
struct CreateUserAccountReq {
    1: string username,
    2: string password,
    3: string real_name,
    4: string contact_phone,
    5: i64 role_id,
    6: optional i64 branch_id,
    7: optional string status
}

// 更新账号请求
struct UpdateUserAccountReq {
    1: i64 id,
    2: optional string username,
    3: optional string password,
    4: optional string real_name,
    5: optional string contact_phone,
    6: optional i64 role_id,
    7: optional i64 branch_id,
    8: optional string status
}

// 账号列表查询请求
struct ListUserAccountsReq {
    1: i32 page = 1,
    2: i32 page_size = 10,
    3: optional i64 role_id,
    4: optional i64 branch_id,
    5: optional string status,
    6: optional string keyword
}

// 账号列表响应
struct ListUserAccountsResp {
    1: list<UserAccount> list,
    2: i64 total,
    3: i32 page,
    4: i32 page_size
}

// 角色信息
struct Role {
    1: i64 id,
    2: string role_name,
    3: optional string description,
    4: string status,
    5: string created_at,
    6: string updated_at,
    7: optional list<i64> permission_ids
}

// 创建角色请求
struct CreateRoleReq {
    1: string role_name,
    2: optional string description,
    3: optional string status,
    4: optional list<i64> permission_ids
}

// 更新角色请求
struct UpdateRoleReq {
    1: i64 id,
    2: optional string role_name,
    3: optional string description,
    4: optional string status,
    5: optional list<i64> permission_ids
}

// 角色列表查询请求
struct ListRolesReq {
    1: i32 page = 1,
    2: i32 page_size = 10,
    3: optional string status,
    4: optional string keyword
}

// 角色列表响应
struct ListRolesResp {
    1: list<Role> list,
    2: i64 total,
    3: i32 page,
    4: i32 page_size
}

// 权限信息
struct Permission {
    1: i64 id,
    2: string permission_name,
    3: string permission_url,
    4: string permission_type,
    5: optional i64 parent_id,
    6: string status,
    7: optional list<Permission> children
}

// 权限列表查询请求
struct ListPermissionsReq {
    1: optional string permission_type,
    2: optional i64 parent_id,
    3: optional string status
}

// 权限列表响应
struct ListPermissionsResp {
    1: list<Permission> list
}

// 渠道配置信息
struct ChannelConfig {
    1: i64 id,
    2: string channel_name,
    3: string channel_code,
    4: string api_url,
    5: string sync_rule,
    6: string status,
    7: string created_at,
    8: string updated_at
}

// 创建渠道配置请求
struct CreateChannelConfigReq {
    1: string channel_name,
    2: string channel_code,
    3: string api_url,
    4: optional string sync_rule,
    5: optional string status
}

// 更新渠道配置请求
struct UpdateChannelConfigReq {
    1: i64 id,
    2: optional string channel_name,
    3: optional string channel_code,
    4: optional string api_url,
    5: optional string sync_rule,
    6: optional string status
}

// 渠道配置列表查询请求
struct ListChannelConfigsReq {
    1: i32 page = 1,
    2: i32 page_size = 10,
    3: optional string status,
    4: optional string keyword
}

// 渠道配置列表响应
struct ListChannelConfigsResp {
    1: list<ChannelConfig> list,
    2: i64 total,
    3: i32 page,
    4: i32 page_size
}

// 系统配置信息
struct SystemConfig {
    1: i64 id,
    2: string config_category,
    3: string config_key,
    4: string config_value,
    5: optional string description,
    6: string status,
    7: string updated_at,
    8: i64 updated_by
}

// 创建系统配置请求
struct CreateSystemConfigReq {
    1: string config_category,
    2: string config_key,
    3: string config_value,
    4: optional string description,
    5: optional string status,
    6: i64 updated_by
}

// 更新系统配置请求
struct UpdateSystemConfigReq {
    1: i64 id,
    2: optional string config_category,
    3: optional string config_key,
    4: optional string config_value,
    5: optional string description,
    6: optional string status,
    7: i64 updated_by
}

// 系统配置列表查询请求
struct ListSystemConfigsReq {
    1: i32 page = 1,
    2: i32 page_size = 10,
    3: optional string config_category,
    4: optional string status,
    5: optional string keyword
}

// 系统配置列表响应
struct ListSystemConfigsResp {
    1: list<SystemConfig> list,
    2: i64 total,
    3: i32 page,
    4: i32 page_size
}

// 黑名单信息
struct Blacklist {
    1: i64 id,
    2: optional i64 guest_id,
    3: optional string guest_name,
    4: string id_number,
    5: string phone,
    6: string reason,
    7: string black_time,
    8: i64 operator_id,
    9: string status,
    10: string created_at
}

// 创建黑名单请求
struct CreateBlacklistReq {
    1: optional i64 guest_id,
    2: string id_number,
    3: string phone,
    4: string reason,
    5: i64 operator_id,
    6: optional string status
}

// 更新黑名单请求
struct UpdateBlacklistReq {
    1: i64 id,
    2: optional i64 guest_id,
    3: optional string id_number,
    4: optional string phone,
    5: optional string reason,
    6: optional string status
}

// 黑名单列表查询请求
struct ListBlacklistsReq {
    1: i32 page = 1,
    2: i32 page_size = 10,
    3: optional string status,
    4: optional string keyword
}

// 黑名单列表响应
struct ListBlacklistsResp {
    1: list<Blacklist> list,
    2: i64 total,
    3: i32 page,
    4: i32 page_size
}

// 会员信息
struct Member {
    1: i64 id,
    2: i64 guest_id,
    3: optional string guest_name,
    4: optional string guest_phone,
    5: string member_level,
    6: i64 points_balance,
    7: string register_time,
    8: optional string last_check_in_time,
    9: string status,
    10: string created_at
}

// 创建会员请求
struct CreateMemberReq {
    1: i64 guest_id,
    2: string member_level,
    3: optional i64 points_balance,
    4: optional string status
}

// 更新会员请求
struct UpdateMemberReq {
    1: i64 id,
    2: optional string member_level,
    3: optional i64 points_balance,
    4: optional string status
}

// 会员列表查询请求
struct ListMembersReq {
    1: i32 page = 1,
    2: i32 page_size = 10,
    3: optional string member_level,
    4: optional string status,
    5: optional string keyword
}

// 会员列表响应
struct ListMembersResp {
    1: list<Member> list,
    2: i64 total,
    3: i32 page,
    4: i32 page_size
}

// 会员权益信息
struct MemberRights {
    1: i64 id,
    2: string member_level,
    3: string rights_name,
    4: optional string description,
    5: optional double discount_ratio,
    6: string effective_time,
    7: optional string expire_time,
    8: string status,
    9: string created_at
}

// 创建会员权益请求
struct CreateMemberRightsReq {
    1: string member_level,
    2: string rights_name,
    3: optional string description,
    4: optional double discount_ratio,
    5: string effective_time,
    6: optional string expire_time,
    7: optional string status
}

// 更新会员权益请求
struct UpdateMemberRightsReq {
    1: i64 id,
    2: optional string member_level,
    3: optional string rights_name,
    4: optional string description,
    5: optional double discount_ratio,
    6: optional string effective_time,
    7: optional string expire_time,
    8: optional string status
}

// 会员权益列表查询请求
struct ListMemberRightsReq {
    1: i32 page = 1,
    2: i32 page_size = 10,
    3: optional string member_level,
    4: optional string status,
    5: optional string keyword
}

// 会员权益列表响应
struct ListMemberRightsResp {
    1: list<MemberRights> list,
    2: i64 total,
    3: i32 page,
    4: i32 page_size
}

// 积分记录信息
struct PointsRecord {
    1: i64 id,
    2: i64 member_id,
    3: optional string member_name,
    4: optional i64 order_id,
    5: string change_type,
    6: i64 points_value,
    7: string change_reason,
    8: string change_time,
    9: i64 operator_id
}

// 创建积分记录请求
struct CreatePointsRecordReq {
    1: i64 member_id,
    2: optional i64 order_id,
    3: string change_type,
    4: i64 points_value,
    5: string change_reason,
    6: i64 operator_id
}

// 积分记录列表查询请求
struct ListPointsRecordsReq {
    1: i32 page = 1,
    2: i32 page_size = 10,
    3: optional i64 member_id,
    4: optional i64 order_id,
    5: optional string change_type,
    6: optional string start_time,
    7: optional string end_time
}

// 积分记录列表响应
struct ListPointsRecordsResp {
    1: list<PointsRecord> list,
    2: i64 total,
    3: i32 page,
    4: i32 page_size
}

// 操作日志信息
struct OperationLog {
    1: i64 id,
    2: i64 operator_id,
    3: optional string operator_name,
    4: string module,
    5: string operation_type,
    6: string content,
    7: string operation_time,
    8: string operation_ip,
    9: optional i64 related_id,
    10: bool is_success,
    11: string created_at
}

// 创建操作日志请求
struct CreateOperationLogReq {
    1: i64 operator_id,
    2: string module,
    3: string operation_type,
    4: string content,
    5: string operation_ip,
    6: optional i64 related_id,
    7: bool is_success
}

// 操作日志列表查询请求
struct ListOperationLogsReq {
    1: i32 page = 1,
    2: i32 page_size = 10,
    3: optional i64 operator_id,
    4: optional string module,
    5: optional string operation_type,
    6: optional string start_time,
    7: optional string end_time,
    8: optional bool is_success
}

// 操作日志列表响应
struct ListOperationLogsResp {
    1: list<OperationLog> list,
    2: i64 total,
    3: i32 page,
    4: i32 page_size
}

// 房源服务接口
service HotelService {
    // 房型管理
    BaseResp CreateRoomType(1: CreateRoomTypeReq req),
    BaseResp UpdateRoomType(1: UpdateRoomTypeReq req),
    RoomType GetRoomType(1: i64 id),
    ListRoomTypeResp ListRoomTypes(1: ListRoomTypeReq req),
    BaseResp DeleteRoomType(1: i64 id),
    
    // 房源管理
    BaseResp CreateRoomInfo(1: CreateRoomInfoReq req),
    BaseResp UpdateRoomInfo(1: UpdateRoomInfoReq req),
    RoomInfo GetRoomInfo(1: i64 id),
    ListRoomInfoResp ListRoomInfos(1: ListRoomInfoReq req),
    BaseResp DeleteRoomInfo(1: i64 id),
    
    // 房源状态管理
    BaseResp UpdateRoomStatus(1: UpdateRoomStatusReq req),
    BaseResp BatchUpdateRoomStatus(1: BatchUpdateRoomStatusReq req),
    
    // 关联房管理
    BaseResp CreateRoomBinding(1: CreateRoomBindingReq req),
    BaseResp BatchCreateRoomBindings(1: BatchCreateRoomBindingsReq req),
    ListRoomBindingsResp GetRoomBindings(1: i64 room_id),
    BaseResp DeleteRoomBinding(1: i64 binding_id),
    
    // 房源图片管理（注意：图片上传需要通过 HTTP API，RPC 不支持文件上传）
    ListRoomImagesResp GetRoomImages(1: i64 room_id),
    BaseResp DeleteRoomImage(1: i64 image_id),
    BaseResp UpdateImageSortOrder(1: UpdateImageSortOrderReq req),
    BaseResp BatchUpdateImageSortOrder(1: BatchUpdateImageSortOrderReq req),
    
    // 设施管理
    BaseResp CreateFacility(1: CreateFacilityReq req),
    BaseResp UpdateFacility(1: UpdateFacilityReq req),
    Facility GetFacility(1: i64 id),
    ListFacilitiesResp ListFacilities(1: ListFacilityReq req),
    BaseResp DeleteFacility(1: i64 id),
    
    // 房源设施关联
    BaseResp SetRoomFacilities(1: SetRoomFacilitiesReq req),
    ListFacilitiesResp GetRoomFacilities(1: i64 room_id),
    BaseResp AddRoomFacility(1: AddRoomFacilityReq req),
    BaseResp RemoveRoomFacility(1: RemoveRoomFacilityReq req),
    
    // 退订政策管理
    BaseResp CreateCancellationPolicy(1: CreateCancellationPolicyReq req),
    BaseResp UpdateCancellationPolicy(1: UpdateCancellationPolicyReq req),
    CancellationPolicy GetCancellationPolicy(1: i64 id),
    ListCancellationPoliciesResp ListCancellationPolicies(1: ListCancellationPolicyReq req),
    BaseResp DeleteCancellationPolicy(1: i64 id),
    
    // 日历化房态管理
    CalendarRoomStatusResp GetCalendarRoomStatus(1: CalendarRoomStatusReq req),
    BaseResp UpdateCalendarRoomStatus(1: UpdateCalendarRoomStatusReq req),
    BaseResp BatchUpdateCalendarRoomStatus(1: BatchUpdateCalendarRoomStatusReq req),
    
    // 实时数据统计
    RealTimeStatisticsResp GetRealTimeStatistics(1: RealTimeStatisticsReq req),
    
    // 分店管理
    ListBranchesResp ListBranches(1: ListBranchesReq req),
    Branch GetBranch(1: i64 branch_id),
    
    // 渠道同步
    SyncRoomStatusToChannelResp SyncRoomStatusToChannel(1: SyncRoomStatusToChannelReq req),
    
    // 订单管理
    ListOrdersResp ListOrders(1: ListOrdersReq req),
    Order GetOrder(1: i64 order_id),
    
    // 在住客人管理
    ListInHouseGuestsResp ListInHouseGuests(1: ListInHouseGuestsReq req),
    
    // 财务管理
    ListFinancialFlowsResp ListFinancialFlows(1: ListFinancialFlowsReq req),
    
    // 账号管理
    BaseResp CreateUserAccount(1: CreateUserAccountReq req),
    BaseResp UpdateUserAccount(1: UpdateUserAccountReq req),
    UserAccount GetUserAccount(1: i64 id),
    ListUserAccountsResp ListUserAccounts(1: ListUserAccountsReq req),
    BaseResp DeleteUserAccount(1: i64 id),
    
    // 角色管理
    BaseResp CreateRole(1: CreateRoleReq req),
    BaseResp UpdateRole(1: UpdateRoleReq req),
    Role GetRole(1: i64 id),
    ListRolesResp ListRoles(1: ListRolesReq req),
    BaseResp DeleteRole(1: i64 id),
    
    // 权限管理
    ListPermissionsResp ListPermissions(1: ListPermissionsReq req),
    
    // 渠道配置管理
    BaseResp CreateChannelConfig(1: CreateChannelConfigReq req),
    BaseResp UpdateChannelConfig(1: UpdateChannelConfigReq req),
    ChannelConfig GetChannelConfig(1: i64 id),
    ListChannelConfigsResp ListChannelConfigs(1: ListChannelConfigsReq req),
    BaseResp DeleteChannelConfig(1: i64 id),
    
    // 系统配置管理
    BaseResp CreateSystemConfig(1: CreateSystemConfigReq req),
    BaseResp UpdateSystemConfig(1: UpdateSystemConfigReq req),
    SystemConfig GetSystemConfig(1: i64 id),
    ListSystemConfigsResp ListSystemConfigs(1: ListSystemConfigsReq req),
    BaseResp DeleteSystemConfig(1: i64 id),
    ListSystemConfigsResp GetSystemConfigsByCategory(1: string category),
    
    // 黑名单管理
    BaseResp CreateBlacklist(1: CreateBlacklistReq req),
    BaseResp UpdateBlacklist(1: UpdateBlacklistReq req),
    Blacklist GetBlacklist(1: i64 id),
    ListBlacklistsResp ListBlacklists(1: ListBlacklistsReq req),
    BaseResp DeleteBlacklist(1: i64 id),
    
    // 会员管理
    BaseResp CreateMember(1: CreateMemberReq req),
    BaseResp UpdateMember(1: UpdateMemberReq req),
    Member GetMember(1: i64 id),
    Member GetMemberByGuestID(1: i64 guest_id),
    ListMembersResp ListMembers(1: ListMembersReq req),
    BaseResp DeleteMember(1: i64 id),
    
    // 会员权益管理
    BaseResp CreateMemberRights(1: CreateMemberRightsReq req),
    BaseResp UpdateMemberRights(1: UpdateMemberRightsReq req),
    MemberRights GetMemberRights(1: i64 id),
    ListMemberRightsResp ListMemberRights(1: ListMemberRightsReq req),
    ListMemberRightsResp GetRightsByMemberLevel(1: string member_level),
    BaseResp DeleteMemberRights(1: i64 id),
    
    // 会员积分管理
    BaseResp CreatePointsRecord(1: CreatePointsRecordReq req),
    ListPointsRecordsResp ListPointsRecords(1: ListPointsRecordsReq req),
    i64 GetMemberPointsBalance(1: i64 member_id),
    
    // 操作日志管理
    BaseResp CreateOperationLog(1: CreateOperationLogReq req),
    ListOperationLogsResp ListOperationLogs(1: ListOperationLogsReq req),
}
