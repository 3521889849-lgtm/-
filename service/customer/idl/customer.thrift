namespace go customer

struct BaseResp {
    1: i32 code
    2: string msg
}

struct CustomerAgent {
    1: string cs_id
    2: string cs_name
    3: string dept_id
    4: string team_id
    5: string skill_tags
    6: i8 status
    7: i8 current_status
}

struct GetCustomerServiceReq {
    1: string cs_id
}

struct GetCustomerServiceResp {
    1: BaseResp base_resp
    2: CustomerAgent customer_service
}

struct ListCustomerServiceReq {
    1: string dept_id
    2: i32 page
    3: i32 page_size
}

struct ListCustomerServiceResp {
    1: BaseResp base_resp
    2: list<CustomerAgent> customer_services
    3: i64 total
}

struct ShiftConfig {
    1: i64 shift_id
    2: string shift_name
    3: string start_time
    4: string end_time
    5: i32 min_staff
    6: i8 is_holiday
    7: string create_by
}

struct CreateShiftConfigReq {
    1: ShiftConfig shift
}

struct CreateShiftConfigResp {
    1: BaseResp base_resp
    2: i64 shift_id
}

struct ListShiftConfigReq {
    1: i8 is_holiday
    2: string shift_name
}

struct ListShiftConfigResp {
    1: BaseResp base_resp
    2: list<ShiftConfig> shifts
    3: i64 total
}

struct UpdateShiftConfigReq {
    1: ShiftConfig shift
}

struct UpdateShiftConfigResp {
    1: BaseResp base_resp
}

struct DeleteShiftConfigReq {
    1: i64 shift_id
}

struct DeleteShiftConfigResp {
    1: BaseResp base_resp
}

struct AssignScheduleReq {
    1: string schedule_date
    2: i64 shift_id
    3: list<string> cs_ids
    4: string create_by
}

struct AssignScheduleResp {
    1: BaseResp base_resp
    2: list<string> conflict_cs_ids
}

// ScheduleCell 排班表格的单元格数据
struct ScheduleCell {
    1: string cs_id
    2: string schedule_date
    3: i64 shift_id
    4: i8 status
}

// ListScheduleGridReq 查询排班表格（按日期范围）
struct ListScheduleGridReq {
    1: string start_date
    2: string end_date
    3: string dept_id
    4: string team_id
}

struct ListScheduleGridResp {
    1: BaseResp base_resp
    2: list<string> dates
    3: list<CustomerAgent> customers
    4: list<ShiftConfig> shifts
    5: list<ScheduleCell> cells
}

// UpsertScheduleCellReq 更新/清空某个客服某天的排班
// shift_id=0 表示清空
struct UpsertScheduleCellReq {
    1: string cs_id
    2: string schedule_date
    3: i64 shift_id
    4: string operator_id
}

struct UpsertScheduleCellResp {
    1: BaseResp base_resp
}

struct AutoScheduleReq {
    1: string start_date
    2: string end_date
    3: string dept_id
    4: string team_id
    5: string operator_id
}

struct AutoScheduleResp {
    1: BaseResp base_resp
    2: i64 schedule_count
}

// ============ 客服自动分配 ============

// AssignCustomerReq 自动分配客服请求
struct AssignCustomerReq {
    1: string user_id           // 用户ID
    2: string user_nickname     // 用户昵称（可选）
    3: string source            // 来源渠道 (APP/Web/H5/WeChat)
}

// AssignCustomerResp 自动分配客服响应
struct AssignCustomerResp {
    1: BaseResp base_resp
    2: string cs_id             // 分配的客服ID
    3: string cs_name           // 客服姓名
    4: string conv_id           // 创建的会话ID
}

// ============ 会话管理 ============

// CreateConversationReq 创建会话请求
struct CreateConversationReq {
    1: string user_id           // 用户ID
    2: string user_nickname     // 用户昵称
    3: string source            // 来源渠道
    4: string cs_id             // 指定客服ID（可选，为空则自动分配）
    5: string first_msg         // 首条消息（可选）
}

// CreateConversationResp 创建会话响应
struct CreateConversationResp {
    1: BaseResp base_resp
    2: string conv_id           // 会话ID
    3: string cs_id             // 客服ID
    4: string cs_name           // 客服名称
    5: bool is_new              // 是否新创建（false表示返回已有会话）
}

// EndConversationReq 结束会话请求
struct EndConversationReq {
    1: string conv_id           // 会话ID
    2: string operator_id       // 操作人(客服ID或系统)
    3: string end_reason        // 结束原因（可选）
}

// EndConversationResp 结束会话响应
struct EndConversationResp {
    1: BaseResp base_resp
    2: i32 duration_seconds     // 会话时长(秒)
}

// TransferConversationReq 转接会话请求
struct TransferConversationReq {
    1: string conv_id           // 会话ID
    2: string from_cs_id        // 转出客服ID
    3: string to_cs_id          // 转入客服ID
    4: string transfer_reason   // 转接原因
    5: string context_remark    // 上下文备注(JSON结构)
}

// TransferConversationResp 转接会话响应
struct TransferConversationResp {
    1: BaseResp base_resp
    2: i64 transfer_id          // 转接记录ID
}

struct ApplyLeaveTransferReq {
    1: string cs_id
    2: i8 apply_type
    3: string target_date
    4: i64 shift_id
    5: string target_cs_id
    6: string reason
}

struct ApplyLeaveTransferResp {
    1: BaseResp base_resp
    2: i64 apply_id
}

struct ApproveLeaveTransferReq {
    1: i64 apply_id
    2: i8 approval_status
    3: string approver_id
}

struct ApproveLeaveTransferResp {
    1: BaseResp base_resp
}

struct LeaveTransferItem {
    1: i64 apply_id
    2: string cs_id
    3: string cs_name
    4: string dept_id
    5: string team_id
    6: i8 apply_type
    7: string target_date
    8: i64 shift_id
    9: string shift_name
    10: string target_cs_id
    11: string target_cs_name
    12: i8 approval_status
    13: string approver_id
    14: string approval_time
    15: string reason
    16: string create_time
}

struct GetLeaveTransferReq {
    1: i64 apply_id
}

struct GetLeaveTransferResp {
    1: BaseResp base_resp
    2: LeaveTransferItem item
}

struct ListLeaveTransferReq {
    1: i8 approval_status
    2: string keyword
    3: i32 page
    4: i32 page_size
}

struct ListLeaveTransferResp {
    1: BaseResp base_resp
    2: list<LeaveTransferItem> items
    3: i64 total
}

// ConversationItem 会话列表项（用于会话管理与记录查询）
struct ConversationItem {
    1: string conv_id
    2: string user_id
    3: string user_nickname
    4: string cs_id
    5: string source
    6: i8 status
    7: string last_msg
    8: string last_time
    9: i64 category_id
    10: string category_name
    11: string tags
    12: i8 is_core
}

struct ListConversationReq {
    1: string cs_id
    2: string keyword
    3: i8 status
    4: i32 page
    5: i32 page_size
}

struct ListConversationHistoryReq {
    1: string cs_id
    2: string keyword
    3: i8 status
    4: i32 page
    5: i32 page_size
}

struct ListConversationResp {
    1: BaseResp base_resp
    2: list<ConversationItem> conversations
    3: i64 total
}

struct ConvMessageItem {
    1: i64 msg_id
    2: string conv_id
    3: i8 sender_type
    4: string sender_id
    5: string msg_content
    6: i8 is_quick_reply
    7: i64 quick_reply_id
    8: string send_time
}

struct ListConversationMessageReq {
    1: string conv_id
    2: i32 page
    3: i32 page_size
    4: i8 order_asc
}

struct ListConversationMessageResp {
    1: BaseResp base_resp
    2: list<ConvMessageItem> messages
    3: i64 total
}

struct SendConversationMessageReq {
    1: string conv_id
    2: i8 sender_type
    3: string sender_id
    4: string msg_content
    5: i8 is_quick_reply
    6: i64 quick_reply_id
}

struct SendConversationMessageResp {
    1: BaseResp base_resp
    2: i64 msg_id
}

struct QuickReplyItem {
    1: i64 reply_id
    2: i8 reply_type
    3: string reply_content
    4: string create_by
    5: i8 is_public
    6: string update_time
}

struct ListQuickReplyReq {
    1: string keyword
    2: i8 reply_type
    3: i8 is_public
    4: i32 page
    5: i32 page_size
}

struct ListQuickReplyResp {
    1: BaseResp base_resp
    2: list<QuickReplyItem> replies
    3: i64 total
}

// CreateQuickReplyReq 创建快捷回复请求
struct CreateQuickReplyReq {
    1: i8 reply_type           // 回复类型
    2: string reply_content    // 回复内容
    3: string create_by        // 创建人
    4: i8 is_public            // 是否公开 0-私有 1-公开
}

// CreateQuickReplyResp 创建快捷回复响应
struct CreateQuickReplyResp {
    1: BaseResp base_resp
    2: i64 reply_id            // 新创建的回复ID
}

// UpdateQuickReplyReq 更新快捷回复请求
struct UpdateQuickReplyReq {
    1: i64 reply_id            // 回复ID
    2: i8 reply_type           // 回复类型
    3: string reply_content    // 回复内容
    4: i8 is_public            // 是否公开 0-私有 1-公开
}

// UpdateQuickReplyResp 更新快捷回复响应
struct UpdateQuickReplyResp {
    1: BaseResp base_resp
}

// DeleteQuickReplyReq 删除快捷回复请求
struct DeleteQuickReplyReq {
    1: i64 reply_id            // 回复ID
}

// DeleteQuickReplyResp 删除快捷回复响应
struct DeleteQuickReplyResp {
    1: BaseResp base_resp
}

// ConvCategory 会话分类（用于记录存储/分类）
struct ConvCategory {
    1: i64 category_id
    2: string category_name
    3: i32 sort_no
    4: string create_by
}

struct CreateConvCategoryReq {
    1: string category_name
    2: i32 sort_no
    3: string create_by
}

struct CreateConvCategoryResp {
    1: BaseResp base_resp
    2: i64 category_id
}

struct ListConvCategoryReq {
}

struct ListConvCategoryResp {
    1: BaseResp base_resp
    2: list<ConvCategory> categories
    3: i64 total
}

struct UpdateConversationClassifyReq {
    1: string conv_id
    2: i64 category_id
    3: string tags
    4: i8 is_core
    5: string operator_id
}

struct UpdateConversationClassifyResp {
    1: BaseResp base_resp
}

// ============ 会话标签管理 ============

// ConvTag 会话标签
struct ConvTag {
    1: i64 tag_id
    2: string tag_name
    3: string tag_color
    4: i32 sort_no
}

struct CreateConvTagReq {
    1: string tag_name
    2: string tag_color
    3: i32 sort_no
    4: string create_by
}

struct CreateConvTagResp {
    1: BaseResp base_resp
    2: i64 tag_id
}

struct ListConvTagReq {
}

struct ListConvTagResp {
    1: BaseResp base_resp
    2: list<ConvTag> tags
    3: i64 total
}

struct UpdateConvTagReq {
    1: i64 tag_id
    2: string tag_name
    3: string tag_color
    4: i32 sort_no
}

struct UpdateConvTagResp {
    1: BaseResp base_resp
}

struct DeleteConvTagReq {
    1: i64 tag_id
}

struct DeleteConvTagResp {
    1: BaseResp base_resp
}

// ============ 会话统计看板 ============

struct GetConversationStatsReq {
    1: string start_date   // 开始日期 YYYY-MM-DD
    2: string end_date     // 结束日期 YYYY-MM-DD
    3: string stat_type    // 统计类型: day/week/month
}

// 标签统计
struct TagStat {
    1: string tag_name
    2: i64 count
    3: double ratio  // 占比百分比
}

// 分类统计
struct CategoryStat {
    1: string category_name
    2: i64 count
    3: double ratio
}

// 处理时长统计
struct DurationStat {
    1: string date
    2: double avg_duration_minutes  // 平均处理时长（分钟）
    3: i64 conv_count
}

struct GetConversationStatsResp {
    1: BaseResp base_resp
    2: list<TagStat> top_tags           // Top问题（按标签）
    3: list<CategoryStat> top_categories // Top问题（按分类）
    4: list<DurationStat> duration_trend // 处理时长趋势
    5: i64 total_conversations           // 总会话数
    6: i64 core_conversations            // 核心会话数
    7: double core_ratio                 // 核心会话占比
}

// ============ 会话监控 ============

// 客服状态信息
struct CsStatusInfo {
    1: string cs_id
    2: string cs_name
    3: i8 online_status          // 0-离线 1-在线 2-忙碌
    4: i32 current_conv_count    // 当前会话数
    5: i32 today_conv_count      // 今日处理总数
    6: string last_active_time   // 最后活跃时间
}

// 会话监控项
struct MonitorConvItem {
    1: string conv_id
    2: string user_id
    3: string user_nickname
    4: string cs_id
    5: string cs_name
    6: i8 status                 // 0-等待 1-进行中 2-已结束
    7: i32 wait_seconds          // 等待时长(秒)
    8: i32 duration_seconds      // 会话时长(秒)
    9: string last_msg           // 最后一条消息
    10: string start_time        // 开始时间
}

struct GetConversationMonitorReq {
    1: string dept_id            // 部门筛选（可选）
    2: i8 status_filter          // 状态筛选 -1-全部 0-等待 1-进行中
}

struct GetConversationMonitorResp {
    1: BaseResp base_resp
    2: list<CsStatusInfo> cs_list       // 客服状态列表
    3: list<MonitorConvItem> conv_list  // 会话列表
    4: i32 waiting_count                // 等待中会话数
    5: i32 ongoing_count                // 进行中会话数
    6: i32 online_cs_count              // 在线客服数
}

// ============ 会话记录导出 ============

struct ExportConversationsReq {
    1: string cs_id              // 客服ID筛选（可选）
    2: string user_id            // 用户ID筛选（可选）
    3: string start_date         // 开始日期
    4: string end_date           // 结束日期
    5: i8 status                 // 状态筛选 -1-全部
    6: string keyword            // 关键词搜索
    7: string export_format      // 导出格式: excel/csv
}

struct ExportConversationsResp {
    1: BaseResp base_resp
    2: binary file_data          // 文件二进制数据
    3: string file_name          // 文件名
    4: i64 total_count           // 导出记录数
}

// ============ 消息分类管理 ============

// 消息分类维度
struct MsgCategory {
    1: i64 category_id
    2: string category_name      // 分类名称: 咨询类/投诉类/建议类/其他类
    3: string keywords           // 关键词列表(JSON)
    4: i32 sort_no
}

// 自动分类请求
struct MsgAutoClassifyReq {
    1: string conv_id            // 会话ID
    2: list<string> msg_contents // 消息内容列表
}

// 自动分类响应
struct MsgAutoClassifyResp {
    1: BaseResp base_resp
    2: i64 category_id           // 分类 ID
    3: string category_name      // 分类名称
    4: double confidence         // 置信度 0-1
    5: bool need_manual_confirm  // 是否需要人工确认
    6: list<string> matched_keywords // 匹配的关键词
}

// 人工调整分类请求
struct AdjustMsgClassifyReq {
    1: string conv_id            // 会话ID
    2: i64 original_category_id  // 原分类ID
    3: i64 new_category_id       // 新分类ID
    4: string operator_id        // 操作人ID
    5: string adjust_reason      // 调整原因
}

// 人工调整分类响应
struct AdjustMsgClassifyResp {
    1: BaseResp base_resp
    2: i64 adjust_log_id         // 调整记录ID
}

// 分类统计查询请求
struct GetClassifyStatsReq {
    1: string start_date         // 开始日期
    2: string end_date           // 结束日期
    3: string stat_type          // 统计类型: day/week/month
}

// 分类统计项
struct ClassifyStatItem {
    1: string date               // 日期
    2: i64 category_id
    3: string category_name
    4: i64 count                 // 数量
    5: double ratio              // 占比
}

// 分类统计查询响应
struct GetClassifyStatsResp {
    1: BaseResp base_resp
    2: list<ClassifyStatItem> daily_stats    // 每日统计
    3: list<CategoryStat> category_summary   // 分类汇总
    4: i64 total_classified                  // 已分类总数
    5: i64 manual_adjusted                   // 人工调整数
    6: double auto_accuracy                  // 自动分类准确率
}

// 消息分类维度CRUD
struct CreateMsgCategoryReq {
    1: string category_name
    2: string keywords           // 关键词列表(JSON)
    3: i32 sort_no
    4: string create_by
}

struct CreateMsgCategoryResp {
    1: BaseResp base_resp
    2: i64 category_id
}

struct ListMsgCategoryReq {
}

struct ListMsgCategoryResp {
    1: BaseResp base_resp
    2: list<MsgCategory> categories
    3: i64 total
}

struct UpdateMsgCategoryReq {
    1: i64 category_id
    2: string category_name
    3: string keywords
    4: i32 sort_no
}

struct UpdateMsgCategoryResp {
    1: BaseResp base_resp
}

struct DeleteMsgCategoryReq {
    1: i64 category_id
}

struct DeleteMsgCategoryResp {
    1: BaseResp base_resp
}

// ============ 用户认证 ============

// UserInfo 用户信息
struct UserInfo {
    1: i64 id
    2: string user_name      // 登录账号
    3: string real_name      // 真实姓名
    4: string phone          // 手机号
    5: string role_code      // 角色编码
    6: string role_name      // 角色名称
    7: i8 status             // 状态 1-正常 0-禁用
}

// LoginReq 登录请求
struct LoginReq {
    1: string user_name      // 登录账号
    2: string password       // 密码
}

// LoginResp 登录响应
struct LoginResp {
    1: BaseResp base_resp
    2: UserInfo user_info    // 用户信息
}

// GetCurrentUserReq 获取当前用户请求
struct GetCurrentUserReq {
    1: i64 user_id           // 用户ID
}

// GetCurrentUserResp 获取当前用户响应
struct GetCurrentUserResp {
    1: BaseResp base_resp
    2: UserInfo user_info
}

// RegisterReq 注册请求（仅允许注册客服账号）
struct RegisterReq {
    1: string user_name      // 登录账号
    2: string password       // 密码
    3: string real_name      // 真实姓名
    4: string phone          // 手机号（可选）
}

// RegisterResp 注册响应
struct RegisterResp {
    1: BaseResp base_resp
    2: i64 user_id           // 新用户ID
}

// ============ 消息加密与脱敏 ============

struct EncryptMessageReq {
    1: string msg_content        // 原始消息内容
}

struct EncryptMessageResp {
    1: BaseResp base_resp
    2: string encrypted_content  // 加密后内容
}

struct DecryptMessageReq {
    1: string encrypted_content  // 加密内容
}

struct DecryptMessageResp {
    1: BaseResp base_resp
    2: string msg_content        // 解密后内容
}

struct DesensitizeMessageReq {
    1: string msg_content        // 原始消息内容
}

struct DesensitizeMessageResp {
    1: BaseResp base_resp
    2: string desensitized_content // 脱敏后内容
    3: list<string> detected_types  // 检测到的敏感信息类型
}

// ============ 数据归档管理 ============

struct ArchiveConversationsReq {
    1: string end_date           // 归档截止日期（该日期之前的数据）
    2: i32 retention_days        // 归档数据保留天数
    3: string operator_id        // 操作人
}

struct ArchiveConversationsResp {
    1: BaseResp base_resp
    2: i64 task_id               // 归档任务ID
    3: i64 archived_count        // 已归档数量
}

struct GetArchiveTaskReq {
    1: i64 task_id               // 任务ID
}

struct GetArchiveTaskResp {
    1: BaseResp base_resp
    2: i64 task_id
    3: string task_type
    4: string start_date
    5: string end_date
    6: i64 archived_count
    7: i64 deleted_count
    8: i8 status                 // 0-进行中 1-完成 2-失败
    9: string error_msg
}

struct QueryArchivedConversationReq {
    1: string user_id            // 用户ID筛选
    2: string cs_id              // 客服ID筛选
    3: string start_date         // 开始日期
    4: string end_date           // 结束日期
    5: i32 page
    6: i32 page_size
}

struct ArchivedConvItem {
    1: string conv_id
    2: string user_id
    3: string cs_id
    4: i32 msg_count
    5: string original_date
    6: string archive_time
}

struct QueryArchivedConversationResp {
    1: BaseResp base_resp
    2: list<ArchivedConvItem> items
    3: i64 total
}

service CustomerService {
    // ============ 客服管理 ============
    // GetCustomerService 获取单个客服详细信息
    GetCustomerServiceResp GetCustomerService(1: GetCustomerServiceReq req)
    // ListCustomerService 分页查询客服列表
    ListCustomerServiceResp ListCustomerService(1: ListCustomerServiceReq req)

    // ============ 班次配置 ============
    // CreateShiftConfig 创建班次模板（如早班、晚班）
    CreateShiftConfigResp CreateShiftConfig(1: CreateShiftConfigReq req)
    // ListShiftConfig 查询班次模板列表
    ListShiftConfigResp ListShiftConfig(1: ListShiftConfigReq req)
    // UpdateShiftConfig 更新现有班次模板
    UpdateShiftConfigResp UpdateShiftConfig(1: UpdateShiftConfigReq req)
    // DeleteShiftConfig 删除指定的班次模板
    DeleteShiftConfigResp DeleteShiftConfig(1: DeleteShiftConfigReq req)

    // ============ 排班管理 ============
    // AssignSchedule 手动为客服分配排班日期和班次
    AssignScheduleResp AssignSchedule(1: AssignScheduleReq req)
    // ListScheduleGrid 获取排班表视图数据（按日期范围）
    ListScheduleGridResp ListScheduleGrid(1: ListScheduleGridReq req)
    // UpsertScheduleCell 更新或清空单个排班单元格
    UpsertScheduleCellResp UpsertScheduleCell(1: UpsertScheduleCellReq req)
    // AutoSchedule 执行自动排班算法生成排班计划
    AutoScheduleResp AutoSchedule(1: AutoScheduleReq req)
    
    // ============ 会话分配 ============
    // AssignCustomer 自动为用户分配当前在线的客服
    AssignCustomerResp AssignCustomer(1: AssignCustomerReq req)
    
    // ============ 会话生命周期 ============
    // CreateConversation 创建或恢复用户与客服的会话
    CreateConversationResp CreateConversation(1: CreateConversationReq req)
    // EndConversation 主动结束进行中的会话
    EndConversationResp EndConversation(1: EndConversationReq req)
    // TransferConversation 将会话从当前客服转交给另一位客服
    TransferConversationResp TransferConversation(1: TransferConversationReq req)
    
    // ============ 请假调班 ============
    // ApplyLeaveTransfer 提交请假或调班申请
    ApplyLeaveTransferResp ApplyLeaveTransfer(1: ApplyLeaveTransferReq req)
    // ApproveLeaveTransfer 审批（通过/拒绝）请假调班申请
    ApproveLeaveTransferResp ApproveLeaveTransfer(1: ApproveLeaveTransferReq req)
    // GetLeaveTransfer 获取单个申请单详情
    GetLeaveTransferResp GetLeaveTransfer(1: GetLeaveTransferReq req)
    // ListLeaveTransfer 查询申请记录列表
    ListLeaveTransferResp ListLeaveTransfer(1: ListLeaveTransferReq req)

    // ============ 会话查询与消息 ============
    // ListConversation 查询客服当前进行中的会话列表
    ListConversationResp ListConversation(1: ListConversationReq req)
    // ListConversationHistory 查询已结束或已转接的历史会话记录
    ListConversationResp ListConversationHistory(1: ListConversationHistoryReq req)
    // ListConversationMessage 获取指定会话的所有消息历史
    ListConversationMessageResp ListConversationMessage(1: ListConversationMessageReq req)
    // SendConversationMessage 发送聊天消息记录
    SendConversationMessageResp SendConversationMessage(1: SendConversationMessageReq req)

    // ============ 快捷回复 ============
    // ListQuickReply 查询常用语快捷回复列表
    ListQuickReplyResp ListQuickReply(1: ListQuickReplyReq req)
    // CreateQuickReply 新增快捷回复内容
    CreateQuickReplyResp CreateQuickReply(1: CreateQuickReplyReq req)
    // UpdateQuickReply 修改快捷回复内容
    UpdateQuickReplyResp UpdateQuickReply(1: UpdateQuickReplyReq req)
    // DeleteQuickReply 删除快捷回复内容
    DeleteQuickReplyResp DeleteQuickReply(1: DeleteQuickReplyReq req)

    // ============ 会话分类管理 ============
    // CreateConvCategory 创建会话业务分类
    CreateConvCategoryResp CreateConvCategory(1: CreateConvCategoryReq req)
    // ListConvCategory 获取所有会话业务分类列表
    ListConvCategoryResp ListConvCategory(1: ListConvCategoryReq req)
    // UpdateConversationClassify 为会话手动设置分类和标签
    UpdateConversationClassifyResp UpdateConversationClassify(1: UpdateConversationClassifyReq req)
    
    // ============ 标签管理 ============
    // CreateConvTag 创建会话辅助标签
    CreateConvTagResp CreateConvTag(1: CreateConvTagReq req)
    // ListConvTag 获取会话标签列表
    ListConvTagResp ListConvTag(1: ListConvTagReq req)
    // UpdateConvTag 更新会话标签信息
    UpdateConvTagResp UpdateConvTag(1: UpdateConvTagReq req)
    // DeleteConvTag 删除指定的会话标签
    DeleteConvTagResp DeleteConvTag(1: DeleteConvTagReq req)
    
    // ============ 统计与监控 ============
    // GetConversationStats 获取会话数据统计分析报告
    GetConversationStatsResp GetConversationStats(1: GetConversationStatsReq req)
    // GetConversationMonitor 获取实时会话看板和客服在线状态监控
    GetConversationMonitorResp GetConversationMonitor(1: GetConversationMonitorReq req)
    // ExportConversations 导出会话历史记录为文件
    ExportConversationsResp ExportConversations(1: ExportConversationsReq req)
    
    // ============ 消息分类与 NLP ============
    // MsgAutoClassify 基于 NLP 自动对消息内容进行智能分类
    MsgAutoClassifyResp MsgAutoClassify(1: MsgAutoClassifyReq req)
    // AdjustMsgClassify 人工修正系统自动分类的结果
    AdjustMsgClassifyResp AdjustMsgClassify(1: AdjustMsgClassifyReq req)
    // GetClassifyStats 获取分类准确率和分布统计
    GetClassifyStatsResp GetClassifyStats(1: GetClassifyStatsReq req)
    
    // ============ 消息分类维度 ============
    // CreateMsgCategory 创建消息分类的维度定义
    CreateMsgCategoryResp CreateMsgCategory(1: CreateMsgCategoryReq req)
    // ListMsgCategory 获取所有消息分类维度
    ListMsgCategoryResp ListMsgCategory(1: ListMsgCategoryReq req)
    // UpdateMsgCategory 更新消息分类维度定义
    UpdateMsgCategoryResp UpdateMsgCategory(1: UpdateMsgCategoryReq req)
    // DeleteMsgCategory 删除消息分类维度定义
    DeleteMsgCategoryResp DeleteMsgCategory(1: DeleteMsgCategoryReq req)
    
    // ============ 身份认证 ============
    // Login 用户登录验证
    LoginResp Login(1: LoginReq req)
    // GetCurrentUser 获取当前登录用户的基本信息
    GetCurrentUserResp GetCurrentUser(1: GetCurrentUserReq req)
    // Register 注册新用户（仅限客服角色）
    RegisterResp Register(1: RegisterReq req)
    
    // ============ 安全与加密 ============
    // EncryptMessage 加密敏感消息内容
    EncryptMessageResp EncryptMessage(1: EncryptMessageReq req)
    // DecryptMessage 解密已加密的消息内容
    DecryptMessageResp DecryptMessage(1: DecryptMessageReq req)
    // DesensitizeMessage 对消息中的敏感信息（如手机号）进行脱敏处理
    DesensitizeMessageResp DesensitizeMessage(1: DesensitizeMessageReq req)
    
    // ============ 数据归档 ============
    // ArchiveConversations 将历史会话数据转存到归档库
    ArchiveConversationsResp ArchiveConversations(1: ArchiveConversationsReq req)
    // GetArchiveTask 查询数据归档任务的进度
    GetArchiveTaskResp GetArchiveTask(1: GetArchiveTaskReq req)
    // QueryArchivedConversation 在归档库中搜索历史会话记录
    QueryArchivedConversationResp QueryArchivedConversation(1: QueryArchivedConversationReq req)
}
