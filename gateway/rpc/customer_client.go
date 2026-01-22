// Package rpc 提供 RPC 客户端封装
// 负责与后端Kitex RPC服务通信
// 主要功能：
// 1. 封装Kitex客户端连接
// 2. 提供业务方法调用
// 3. 集成链路追踪（创建 Client Span）
package rpc

import (
	"context"
	"fmt"

	"example_shop/pkg/logger"                                          // 日志工具
	"example_shop/pkg/trace"                                           // 链路追踪
	"example_shop/service/customer/kitex_gen/customer"                 // RPC服务接口定义
	"example_shop/service/customer/kitex_gen/customer/customerservice" // Kitex生成的客户端

	"github.com/cloudwego/kitex/client" // Kitex客户端配置
	"go.uber.org/zap"
)

// CustomerClient 客服服务RPC客户端封装
// 封装Kitex生成的客户端，提供统一的调用接口
type CustomerClient struct {
	client      customerservice.Client // Kitex生成的客户端实例
	serviceName string                 // 服务名称（用于追踪）
}

// NewCustomerClient 创建客服服务RPC客户端
// 参数:
//   - serviceName: 服务名称（用于服务发现）
//   - address: 服务地址（如 "127.0.0.1:9999"）
//
// 返回: 客户端实例和错误信息
func NewCustomerClient(serviceName, address string) (*CustomerClient, error) {
	if address == "" {
		return nil, fmt.Errorf("address is required")
	}
	if serviceName == "" {
		serviceName = "CustomerService" // 默认服务名
	}

	// 创建Kitex客户端实例
	c, err := customerservice.NewClient(
		serviceName,
		client.WithHostPorts(address), // 指定服务地址
	)
	if err != nil {
		return nil, fmt.Errorf("create kitex client failed (service=%q address=%q): %w", serviceName, address, err)
	}
	return &CustomerClient{
		client:      c,
		serviceName: serviceName,
	}, nil
}

// wrapContext 包装context，创建 Client Span 并添加RPC调用日志
// 在每次RPC调用前创建 span，记录方法名和TraceID
// 返回：新的 context 和 span（调用方需要在 RPC 调用后调用 span.End()）
func (c *CustomerClient) wrapContext(ctx context.Context, method string) (context.Context, *trace.Span) {
	// 创建 Client Span
	ctx, span := trace.StartClientSpan(ctx, c.serviceName, method)

	// 记录日志（兼容旧代码）
	traceID := trace.GetTraceID(ctx)
	if traceID != "" {
		logger.InfoWithTrace(ctx, "RPC Call",
			zap.String("method", method),
			zap.String("service", c.serviceName),
		)
	}

	return ctx, span
}

// ==================== 客服基础信息接口 ====================

// GetCustomerService 获取客服信息
func (c *CustomerClient) GetCustomerService(ctx context.Context, req *customer.GetCustomerServiceReq) (*customer.GetCustomerServiceResp, error) {
	ctx, span := c.wrapContext(ctx, "GetCustomerService")
	defer span.End()
	resp, err := c.client.GetCustomerService(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ListCustomerService 查询客服列表
func (c *CustomerClient) ListCustomerService(ctx context.Context, req *customer.ListCustomerServiceReq) (*customer.ListCustomerServiceResp, error) {
	ctx, span := c.wrapContext(ctx, "ListCustomerService")
	defer span.End()
	resp, err := c.client.ListCustomerService(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ==================== 班次配置接口 ====================

// CreateShiftConfig 创建班次配置
func (c *CustomerClient) CreateShiftConfig(ctx context.Context, req *customer.CreateShiftConfigReq) (*customer.CreateShiftConfigResp, error) {
	ctx, span := c.wrapContext(ctx, "CreateShiftConfig")
	defer span.End()
	resp, err := c.client.CreateShiftConfig(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ListShiftConfig 查询班次配置列表
func (c *CustomerClient) ListShiftConfig(ctx context.Context, req *customer.ListShiftConfigReq) (*customer.ListShiftConfigResp, error) {
	ctx, span := c.wrapContext(ctx, "ListShiftConfig")
	defer span.End()
	resp, err := c.client.ListShiftConfig(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// UpdateShiftConfig 更新班次配置
func (c *CustomerClient) UpdateShiftConfig(ctx context.Context, req *customer.UpdateShiftConfigReq) (*customer.UpdateShiftConfigResp, error) {
	ctx, span := c.wrapContext(ctx, "UpdateShiftConfig")
	defer span.End()
	resp, err := c.client.UpdateShiftConfig(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// DeleteShiftConfig 删除班次配置
func (c *CustomerClient) DeleteShiftConfig(ctx context.Context, req *customer.DeleteShiftConfigReq) (*customer.DeleteShiftConfigResp, error) {
	ctx, span := c.wrapContext(ctx, "DeleteShiftConfig")
	defer span.End()
	resp, err := c.client.DeleteShiftConfig(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ==================== 排班管理接口 ====================

// AssignSchedule 手动分配排班
func (c *CustomerClient) AssignSchedule(ctx context.Context, req *customer.AssignScheduleReq) (*customer.AssignScheduleResp, error) {
	ctx, span := c.wrapContext(ctx, "AssignSchedule")
	defer span.End()
	resp, err := c.client.AssignSchedule(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// AutoSchedule 自动排班
func (c *CustomerClient) AutoSchedule(ctx context.Context, req *customer.AutoScheduleReq) (*customer.AutoScheduleResp, error) {
	ctx, span := c.wrapContext(ctx, "AutoSchedule")
	defer span.End()
	resp, err := c.client.AutoSchedule(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ==================== 请假/调班接口 ====================

// ApplyLeaveTransfer 提交请假/调班申请
func (c *CustomerClient) ApplyLeaveTransfer(ctx context.Context, req *customer.ApplyLeaveTransferReq) (*customer.ApplyLeaveTransferResp, error) {
	ctx, span := c.wrapContext(ctx, "ApplyLeaveTransfer")
	defer span.End()
	resp, err := c.client.ApplyLeaveTransfer(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ApproveLeaveTransfer 审批请假/调班申请
func (c *CustomerClient) ApproveLeaveTransfer(ctx context.Context, req *customer.ApproveLeaveTransferReq) (*customer.ApproveLeaveTransferResp, error) {
	ctx, span := c.wrapContext(ctx, "ApproveLeaveTransfer")
	defer span.End()
	resp, err := c.client.ApproveLeaveTransfer(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// GetLeaveTransfer 获取请假/调班申请详情
func (c *CustomerClient) GetLeaveTransfer(ctx context.Context, req *customer.GetLeaveTransferReq) (*customer.GetLeaveTransferResp, error) {
	ctx, span := c.wrapContext(ctx, "GetLeaveTransfer")
	defer span.End()
	resp, err := c.client.GetLeaveTransfer(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ListLeaveTransfer 查询请假/调班申请列表
func (c *CustomerClient) ListLeaveTransfer(ctx context.Context, req *customer.ListLeaveTransferReq) (*customer.ListLeaveTransferResp, error) {
	ctx, span := c.wrapContext(ctx, "ListLeaveTransfer")
	defer span.End()
	resp, err := c.client.ListLeaveTransfer(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ==================== 排班表格接口 ====================

// ListScheduleGrid 查询排班表格数据
func (c *CustomerClient) ListScheduleGrid(ctx context.Context, req *customer.ListScheduleGridReq) (*customer.ListScheduleGridResp, error) {
	ctx, span := c.wrapContext(ctx, "ListScheduleGrid")
	defer span.End()
	resp, err := c.client.ListScheduleGrid(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// UpsertScheduleCell 更新排班单元格
func (c *CustomerClient) UpsertScheduleCell(ctx context.Context, req *customer.UpsertScheduleCellReq) (*customer.UpsertScheduleCellResp, error) {
	ctx, span := c.wrapContext(ctx, "UpsertScheduleCell")
	defer span.End()
	resp, err := c.client.UpsertScheduleCell(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ==================== 会话管理接口 ====================

// AssignCustomer 自动分配客服
func (c *CustomerClient) AssignCustomer(ctx context.Context, req *customer.AssignCustomerReq) (*customer.AssignCustomerResp, error) {
	ctx, span := c.wrapContext(ctx, "AssignCustomer")
	defer span.End()
	resp, err := c.client.AssignCustomer(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// CreateConversation 创建会话
// 支持用户发起新会话，可指定客服或自动分配
func (c *CustomerClient) CreateConversation(ctx context.Context, req *customer.CreateConversationReq) (*customer.CreateConversationResp, error) {
	ctx, span := c.wrapContext(ctx, "CreateConversation")
	defer span.End()
	// 设置业务属性（自动脱敏）
	if req.UserId != "" {
		span.SetBusinessAttrs(req.UserId, "")
	}
	resp, err := c.client.CreateConversation(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// EndConversation 结束会话
// 由客服或系统主动结束会话，支持乐观锁并发控制
func (c *CustomerClient) EndConversation(ctx context.Context, req *customer.EndConversationReq) (*customer.EndConversationResp, error) {
	ctx, span := c.wrapContext(ctx, "EndConversation")
	defer span.End()
	resp, err := c.client.EndConversation(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// TransferConversation 转接会话
// 将会话从当前客服转接给另一位客服，支持传递上下文备注
func (c *CustomerClient) TransferConversation(ctx context.Context, req *customer.TransferConversationReq) (*customer.TransferConversationResp, error) {
	ctx, span := c.wrapContext(ctx, "TransferConversation")
	defer span.End()
	resp, err := c.client.TransferConversation(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ListConversation 查询会话列表
func (c *CustomerClient) ListConversation(ctx context.Context, req *customer.ListConversationReq) (*customer.ListConversationResp, error) {
	ctx, span := c.wrapContext(ctx, "ListConversation")
	defer span.End()
	resp, err := c.client.ListConversation(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ListConversationHistory 查询历史会话记录
func (c *CustomerClient) ListConversationHistory(ctx context.Context, req *customer.ListConversationHistoryReq) (*customer.ListConversationResp, error) {
	ctx, span := c.wrapContext(ctx, "ListConversationHistory")
	defer span.End()
	resp, err := c.client.ListConversationHistory(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ListConversationMessage 查询会话消息列表
func (c *CustomerClient) ListConversationMessage(ctx context.Context, req *customer.ListConversationMessageReq) (*customer.ListConversationMessageResp, error) {
	ctx, span := c.wrapContext(ctx, "ListConversationMessage")
	defer span.End()
	resp, err := c.client.ListConversationMessage(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// SendConversationMessage 发送会话消息
func (c *CustomerClient) SendConversationMessage(ctx context.Context, req *customer.SendConversationMessageReq) (*customer.SendConversationMessageResp, error) {
	ctx, span := c.wrapContext(ctx, "SendConversationMessage")
	defer span.End()
	resp, err := c.client.SendConversationMessage(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ListQuickReply 查询快捷回复列表
func (c *CustomerClient) ListQuickReply(ctx context.Context, req *customer.ListQuickReplyReq) (*customer.ListQuickReplyResp, error) {
	ctx, span := c.wrapContext(ctx, "ListQuickReply")
	defer span.End()
	resp, err := c.client.ListQuickReply(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// CreateQuickReply 创建快捷回复
func (c *CustomerClient) CreateQuickReply(ctx context.Context, req *customer.CreateQuickReplyReq) (*customer.CreateQuickReplyResp, error) {
	ctx, span := c.wrapContext(ctx, "CreateQuickReply")
	defer span.End()
	resp, err := c.client.CreateQuickReply(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// UpdateQuickReply 更新快捷回复
func (c *CustomerClient) UpdateQuickReply(ctx context.Context, req *customer.UpdateQuickReplyReq) (*customer.UpdateQuickReplyResp, error) {
	ctx, span := c.wrapContext(ctx, "UpdateQuickReply")
	defer span.End()
	resp, err := c.client.UpdateQuickReply(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// DeleteQuickReply 删除快捷回复
func (c *CustomerClient) DeleteQuickReply(ctx context.Context, req *customer.DeleteQuickReplyReq) (*customer.DeleteQuickReplyResp, error) {
	ctx, span := c.wrapContext(ctx, "DeleteQuickReply")
	defer span.End()
	resp, err := c.client.DeleteQuickReply(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ==================== 会话分类接口 ====================

// CreateConvCategory 创建会话分类
func (c *CustomerClient) CreateConvCategory(ctx context.Context, req *customer.CreateConvCategoryReq) (*customer.CreateConvCategoryResp, error) {
	ctx, span := c.wrapContext(ctx, "CreateConvCategory")
	defer span.End()
	resp, err := c.client.CreateConvCategory(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ListConvCategory 查询会话分类列表
func (c *CustomerClient) ListConvCategory(ctx context.Context, req *customer.ListConvCategoryReq) (*customer.ListConvCategoryResp, error) {
	ctx, span := c.wrapContext(ctx, "ListConvCategory")
	defer span.End()
	resp, err := c.client.ListConvCategory(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// UpdateConversationClassify 更新会话分类/标签
func (c *CustomerClient) UpdateConversationClassify(ctx context.Context, req *customer.UpdateConversationClassifyReq) (*customer.UpdateConversationClassifyResp, error) {
	ctx, span := c.wrapContext(ctx, "UpdateConversationClassify")
	defer span.End()
	resp, err := c.client.UpdateConversationClassify(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ==================== 会话标签接口 ====================

// CreateConvTag 创建会话标签
func (c *CustomerClient) CreateConvTag(ctx context.Context, req *customer.CreateConvTagReq) (*customer.CreateConvTagResp, error) {
	ctx, span := c.wrapContext(ctx, "CreateConvTag")
	defer span.End()
	resp, err := c.client.CreateConvTag(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ListConvTag 查询会话标签列表
func (c *CustomerClient) ListConvTag(ctx context.Context, req *customer.ListConvTagReq) (*customer.ListConvTagResp, error) {
	ctx, span := c.wrapContext(ctx, "ListConvTag")
	defer span.End()
	resp, err := c.client.ListConvTag(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// UpdateConvTag 更新会话标签
func (c *CustomerClient) UpdateConvTag(ctx context.Context, req *customer.UpdateConvTagReq) (*customer.UpdateConvTagResp, error) {
	ctx, span := c.wrapContext(ctx, "UpdateConvTag")
	defer span.End()
	resp, err := c.client.UpdateConvTag(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// DeleteConvTag 删除会话标签
func (c *CustomerClient) DeleteConvTag(ctx context.Context, req *customer.DeleteConvTagReq) (*customer.DeleteConvTagResp, error) {
	ctx, span := c.wrapContext(ctx, "DeleteConvTag")
	defer span.End()
	resp, err := c.client.DeleteConvTag(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ==================== 统计接口 ====================

// GetConversationStats 获取会话统计数据
func (c *CustomerClient) GetConversationStats(ctx context.Context, req *customer.GetConversationStatsReq) (*customer.GetConversationStatsResp, error) {
	ctx, span := c.wrapContext(ctx, "GetConversationStats")
	defer span.End()
	resp, err := c.client.GetConversationStats(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ==================== 会话监控与导出接口 ====================

// GetConversationMonitor 获取会话监控数据
// 实时查看会话状态、客服在线状态、等待中会话数等
func (c *CustomerClient) GetConversationMonitor(ctx context.Context, req *customer.GetConversationMonitorReq) (*customer.GetConversationMonitorResp, error) {
	ctx, span := c.wrapContext(ctx, "GetConversationMonitor")
	defer span.End()
	resp, err := c.client.GetConversationMonitor(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ExportConversations 导出会话记录
// 支持按条件筛选导出会话记录为Excel/CSV格式
func (c *CustomerClient) ExportConversations(ctx context.Context, req *customer.ExportConversationsReq) (*customer.ExportConversationsResp, error) {
	ctx, span := c.wrapContext(ctx, "ExportConversations")
	defer span.End()
	resp, err := c.client.ExportConversations(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ==================== 消息分类管理接口 ====================

// MsgAutoClassify 消息自动分类
// 基于关键词匹配对会话消息进行自动分类
func (c *CustomerClient) MsgAutoClassify(ctx context.Context, req *customer.MsgAutoClassifyReq) (*customer.MsgAutoClassifyResp, error) {
	ctx, span := c.wrapContext(ctx, "MsgAutoClassify")
	defer span.End()
	resp, err := c.client.MsgAutoClassify(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// AdjustMsgClassify 人工调整消息分类
// 客服手动修正自动分类结果
func (c *CustomerClient) AdjustMsgClassify(ctx context.Context, req *customer.AdjustMsgClassifyReq) (*customer.AdjustMsgClassifyResp, error) {
	ctx, span := c.wrapContext(ctx, "AdjustMsgClassify")
	defer span.End()
	resp, err := c.client.AdjustMsgClassify(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// GetClassifyStats 获取分类统计数据
// 查询消息分类的统计信息，支持按日期范围和统计类型筛选
func (c *CustomerClient) GetClassifyStats(ctx context.Context, req *customer.GetClassifyStatsReq) (*customer.GetClassifyStatsResp, error) {
	ctx, span := c.wrapContext(ctx, "GetClassifyStats")
	defer span.End()
	resp, err := c.client.GetClassifyStats(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ==================== 消息分类维度CRUD ====================

// CreateMsgCategory 创建消息分类维度
func (c *CustomerClient) CreateMsgCategory(ctx context.Context, req *customer.CreateMsgCategoryReq) (*customer.CreateMsgCategoryResp, error) {
	ctx, span := c.wrapContext(ctx, "CreateMsgCategory")
	defer span.End()
	resp, err := c.client.CreateMsgCategory(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ListMsgCategory 查询消息分类维度列表
func (c *CustomerClient) ListMsgCategory(ctx context.Context, req *customer.ListMsgCategoryReq) (*customer.ListMsgCategoryResp, error) {
	ctx, span := c.wrapContext(ctx, "ListMsgCategory")
	defer span.End()
	resp, err := c.client.ListMsgCategory(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// UpdateMsgCategory 更新消息分类维度
func (c *CustomerClient) UpdateMsgCategory(ctx context.Context, req *customer.UpdateMsgCategoryReq) (*customer.UpdateMsgCategoryResp, error) {
	ctx, span := c.wrapContext(ctx, "UpdateMsgCategory")
	defer span.End()
	resp, err := c.client.UpdateMsgCategory(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// DeleteMsgCategory 删除消息分类维度
func (c *CustomerClient) DeleteMsgCategory(ctx context.Context, req *customer.DeleteMsgCategoryReq) (*customer.DeleteMsgCategoryResp, error) {
	ctx, span := c.wrapContext(ctx, "DeleteMsgCategory")
	defer span.End()
	resp, err := c.client.DeleteMsgCategory(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ==================== 用户认证接口 ====================

// Login 用户登录
func (c *CustomerClient) Login(ctx context.Context, req *customer.LoginReq) (*customer.LoginResp, error) {
	ctx, span := c.wrapContext(ctx, "Login")
	defer span.End()
	resp, err := c.client.Login(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// GetCurrentUser 获取当前登录用户信息
func (c *CustomerClient) GetCurrentUser(ctx context.Context, req *customer.GetCurrentUserReq) (*customer.GetCurrentUserResp, error) {
	ctx, span := c.wrapContext(ctx, "GetCurrentUser")
	defer span.End()
	resp, err := c.client.GetCurrentUser(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// Register 用户注册
func (c *CustomerClient) Register(ctx context.Context, req *customer.RegisterReq) (*customer.RegisterResp, error) {
	ctx, span := c.wrapContext(ctx, "Register")
	defer span.End()
	resp, err := c.client.Register(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ==================== 消息加密与脱敏接口 ====================

// EncryptMessage 加密消息内容
func (c *CustomerClient) EncryptMessage(ctx context.Context, req *customer.EncryptMessageReq) (*customer.EncryptMessageResp, error) {
	ctx, span := c.wrapContext(ctx, "EncryptMessage")
	defer span.End()
	resp, err := c.client.EncryptMessage(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// DecryptMessage 解密消息内容
func (c *CustomerClient) DecryptMessage(ctx context.Context, req *customer.DecryptMessageReq) (*customer.DecryptMessageResp, error) {
	ctx, span := c.wrapContext(ctx, "DecryptMessage")
	defer span.End()
	resp, err := c.client.DecryptMessage(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// DesensitizeMessage 消息脱敏处理
func (c *CustomerClient) DesensitizeMessage(ctx context.Context, req *customer.DesensitizeMessageReq) (*customer.DesensitizeMessageResp, error) {
	ctx, span := c.wrapContext(ctx, "DesensitizeMessage")
	defer span.End()
	resp, err := c.client.DesensitizeMessage(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// ==================== 数据归档管理接口 ====================

// ArchiveConversations 归档历史会话数据
func (c *CustomerClient) ArchiveConversations(ctx context.Context, req *customer.ArchiveConversationsReq) (*customer.ArchiveConversationsResp, error) {
	ctx, span := c.wrapContext(ctx, "ArchiveConversations")
	defer span.End()
	resp, err := c.client.ArchiveConversations(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// GetArchiveTask 查询归档任务状态
func (c *CustomerClient) GetArchiveTask(ctx context.Context, req *customer.GetArchiveTaskReq) (*customer.GetArchiveTaskResp, error) {
	ctx, span := c.wrapContext(ctx, "GetArchiveTask")
	defer span.End()
	resp, err := c.client.GetArchiveTask(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// QueryArchivedConversation 查询已归档会话
func (c *CustomerClient) QueryArchivedConversation(ctx context.Context, req *customer.QueryArchivedConversationReq) (*customer.QueryArchivedConversationResp, error) {
	ctx, span := c.wrapContext(ctx, "QueryArchivedConversation")
	defer span.End()
	resp, err := c.client.QueryArchivedConversation(ctx, req)
	if err != nil {
		span.SetError(err)
	}
	return resp, err
}

// Close 关闭RPC客户端连接
// 当前实现为空操作，Kitex客户端无需显式关闭
func (c *CustomerClient) Close() error {
	return nil
}
