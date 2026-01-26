package main

import (
	"context"
	customer "example_shop/service/customer/kitex_gen/customer"
)

// CustomerServiceImpl implements the last service interface defined in the IDL.
type CustomerServiceImpl struct{}

// GetCustomerService implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) GetCustomerService(ctx context.Context, req *customer.GetCustomerServiceReq) (resp *customer.GetCustomerServiceResp, err error) {
	// TODO: Your code here...
	return
}

// ListCustomerService implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ListCustomerService(ctx context.Context, req *customer.ListCustomerServiceReq) (resp *customer.ListCustomerServiceResp, err error) {
	// TODO: Your code here...
	return
}

// CreateShiftConfig implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) CreateShiftConfig(ctx context.Context, req *customer.CreateShiftConfigReq) (resp *customer.CreateShiftConfigResp, err error) {
	// TODO: Your code here...
	return
}

// ListShiftConfig implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ListShiftConfig(ctx context.Context, req *customer.ListShiftConfigReq) (resp *customer.ListShiftConfigResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateShiftConfig implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) UpdateShiftConfig(ctx context.Context, req *customer.UpdateShiftConfigReq) (resp *customer.UpdateShiftConfigResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteShiftConfig implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) DeleteShiftConfig(ctx context.Context, req *customer.DeleteShiftConfigReq) (resp *customer.DeleteShiftConfigResp, err error) {
	// TODO: Your code here...
	return
}

// AssignSchedule implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) AssignSchedule(ctx context.Context, req *customer.AssignScheduleReq) (resp *customer.AssignScheduleResp, err error) {
	// TODO: Your code here...
	return
}

// ListScheduleGrid implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ListScheduleGrid(ctx context.Context, req *customer.ListScheduleGridReq) (resp *customer.ListScheduleGridResp, err error) {
	// TODO: Your code here...
	return
}

// UpsertScheduleCell implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) UpsertScheduleCell(ctx context.Context, req *customer.UpsertScheduleCellReq) (resp *customer.UpsertScheduleCellResp, err error) {
	// TODO: Your code here...
	return
}

// AutoSchedule implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) AutoSchedule(ctx context.Context, req *customer.AutoScheduleReq) (resp *customer.AutoScheduleResp, err error) {
	// TODO: Your code here...
	return
}

// AssignCustomer implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) AssignCustomer(ctx context.Context, req *customer.AssignCustomerReq) (resp *customer.AssignCustomerResp, err error) {
	// TODO: Your code here...
	return
}

// CreateConversation implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) CreateConversation(ctx context.Context, req *customer.CreateConversationReq) (resp *customer.CreateConversationResp, err error) {
	// TODO: Your code here...
	return
}

// EndConversation implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) EndConversation(ctx context.Context, req *customer.EndConversationReq) (resp *customer.EndConversationResp, err error) {
	// TODO: Your code here...
	return
}

// TransferConversation implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) TransferConversation(ctx context.Context, req *customer.TransferConversationReq) (resp *customer.TransferConversationResp, err error) {
	// TODO: Your code here...
	return
}

// ApplyLeaveTransfer implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ApplyLeaveTransfer(ctx context.Context, req *customer.ApplyLeaveTransferReq) (resp *customer.ApplyLeaveTransferResp, err error) {
	// TODO: Your code here...
	return
}

// ApproveLeaveTransfer implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ApproveLeaveTransfer(ctx context.Context, req *customer.ApproveLeaveTransferReq) (resp *customer.ApproveLeaveTransferResp, err error) {
	// TODO: Your code here...
	return
}

// GetLeaveTransfer implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) GetLeaveTransfer(ctx context.Context, req *customer.GetLeaveTransferReq) (resp *customer.GetLeaveTransferResp, err error) {
	// TODO: Your code here...
	return
}

// ListLeaveTransfer implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ListLeaveTransfer(ctx context.Context, req *customer.ListLeaveTransferReq) (resp *customer.ListLeaveTransferResp, err error) {
	// TODO: Your code here...
	return
}

// GetLeaveAuditLog implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) GetLeaveAuditLog(ctx context.Context, req *customer.GetLeaveAuditLogReq) (resp *customer.GetLeaveAuditLogResp, err error) {
	// TODO: Your code here...
	return
}

// ApplyChainSwap implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ApplyChainSwap(ctx context.Context, req *customer.ApplyChainSwapReq) (resp *customer.ApplyChainSwapResp, err error) {
	// TODO: Your code here...
	return
}

// ApproveChainSwap implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ApproveChainSwap(ctx context.Context, req *customer.ApproveChainSwapReq) (resp *customer.ApproveChainSwapResp, err error) {
	// TODO: Your code here...
	return
}

// ListChainSwap implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ListChainSwap(ctx context.Context, req *customer.ListChainSwapReq) (resp *customer.ListChainSwapResp, err error) {
	// TODO: Your code here...
	return
}

// GetChainSwap implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) GetChainSwap(ctx context.Context, req *customer.GetChainSwapReq) (resp *customer.GetChainSwapResp, err error) {
	// TODO: Your code here...
	return
}

// Heartbeat implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) Heartbeat(ctx context.Context, req *customer.HeartbeatReq) (resp *customer.HeartbeatResp, err error) {
	// TODO: Your code here...
	return
}

// ListOnlineCustomers implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ListOnlineCustomers(ctx context.Context, req *customer.ListOnlineCustomersReq) (resp *customer.ListOnlineCustomersResp, err error) {
	// TODO: Your code here...
	return
}

// GetSwapCandidates implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) GetSwapCandidates(ctx context.Context, req *customer.GetSwapCandidatesReq) (resp *customer.GetSwapCandidatesResp, err error) {
	// TODO: Your code here...
	return
}

// CheckSwapConflict implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) CheckSwapConflict(ctx context.Context, req *customer.CheckSwapConflictReq) (resp *customer.CheckSwapConflictResp, err error) {
	// TODO: Your code here...
	return
}

// ListConversation implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ListConversation(ctx context.Context, req *customer.ListConversationReq) (resp *customer.ListConversationResp, err error) {
	// TODO: Your code here...
	return
}

// ListConversationHistory implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ListConversationHistory(ctx context.Context, req *customer.ListConversationHistoryReq) (resp *customer.ListConversationResp, err error) {
	// TODO: Your code here...
	return
}

// ListConversationMessage implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ListConversationMessage(ctx context.Context, req *customer.ListConversationMessageReq) (resp *customer.ListConversationMessageResp, err error) {
	// TODO: Your code here...
	return
}

// SendConversationMessage implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) SendConversationMessage(ctx context.Context, req *customer.SendConversationMessageReq) (resp *customer.SendConversationMessageResp, err error) {
	// TODO: Your code here...
	return
}

// ListQuickReply implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ListQuickReply(ctx context.Context, req *customer.ListQuickReplyReq) (resp *customer.ListQuickReplyResp, err error) {
	// TODO: Your code here...
	return
}

// CreateQuickReply implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) CreateQuickReply(ctx context.Context, req *customer.CreateQuickReplyReq) (resp *customer.CreateQuickReplyResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateQuickReply implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) UpdateQuickReply(ctx context.Context, req *customer.UpdateQuickReplyReq) (resp *customer.UpdateQuickReplyResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteQuickReply implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) DeleteQuickReply(ctx context.Context, req *customer.DeleteQuickReplyReq) (resp *customer.DeleteQuickReplyResp, err error) {
	// TODO: Your code here...
	return
}

// CreateConvCategory implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) CreateConvCategory(ctx context.Context, req *customer.CreateConvCategoryReq) (resp *customer.CreateConvCategoryResp, err error) {
	// TODO: Your code here...
	return
}

// ListConvCategory implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ListConvCategory(ctx context.Context, req *customer.ListConvCategoryReq) (resp *customer.ListConvCategoryResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateConversationClassify implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) UpdateConversationClassify(ctx context.Context, req *customer.UpdateConversationClassifyReq) (resp *customer.UpdateConversationClassifyResp, err error) {
	// TODO: Your code here...
	return
}

// CreateConvTag implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) CreateConvTag(ctx context.Context, req *customer.CreateConvTagReq) (resp *customer.CreateConvTagResp, err error) {
	// TODO: Your code here...
	return
}

// ListConvTag implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ListConvTag(ctx context.Context, req *customer.ListConvTagReq) (resp *customer.ListConvTagResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateConvTag implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) UpdateConvTag(ctx context.Context, req *customer.UpdateConvTagReq) (resp *customer.UpdateConvTagResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteConvTag implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) DeleteConvTag(ctx context.Context, req *customer.DeleteConvTagReq) (resp *customer.DeleteConvTagResp, err error) {
	// TODO: Your code here...
	return
}

// GetConversationStats implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) GetConversationStats(ctx context.Context, req *customer.GetConversationStatsReq) (resp *customer.GetConversationStatsResp, err error) {
	// TODO: Your code here...
	return
}

// GetConversationMonitor implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) GetConversationMonitor(ctx context.Context, req *customer.GetConversationMonitorReq) (resp *customer.GetConversationMonitorResp, err error) {
	// TODO: Your code here...
	return
}

// ExportConversations implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ExportConversations(ctx context.Context, req *customer.ExportConversationsReq) (resp *customer.ExportConversationsResp, err error) {
	// TODO: Your code here...
	return
}

// MsgAutoClassify implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) MsgAutoClassify(ctx context.Context, req *customer.MsgAutoClassifyReq) (resp *customer.MsgAutoClassifyResp, err error) {
	// TODO: Your code here...
	return
}

// AdjustMsgClassify implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) AdjustMsgClassify(ctx context.Context, req *customer.AdjustMsgClassifyReq) (resp *customer.AdjustMsgClassifyResp, err error) {
	// TODO: Your code here...
	return
}

// GetClassifyStats implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) GetClassifyStats(ctx context.Context, req *customer.GetClassifyStatsReq) (resp *customer.GetClassifyStatsResp, err error) {
	// TODO: Your code here...
	return
}

// CreateMsgCategory implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) CreateMsgCategory(ctx context.Context, req *customer.CreateMsgCategoryReq) (resp *customer.CreateMsgCategoryResp, err error) {
	// TODO: Your code here...
	return
}

// ListMsgCategory implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ListMsgCategory(ctx context.Context, req *customer.ListMsgCategoryReq) (resp *customer.ListMsgCategoryResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateMsgCategory implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) UpdateMsgCategory(ctx context.Context, req *customer.UpdateMsgCategoryReq) (resp *customer.UpdateMsgCategoryResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteMsgCategory implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) DeleteMsgCategory(ctx context.Context, req *customer.DeleteMsgCategoryReq) (resp *customer.DeleteMsgCategoryResp, err error) {
	// TODO: Your code here...
	return
}

// Login implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) Login(ctx context.Context, req *customer.LoginReq) (resp *customer.LoginResp, err error) {
	// TODO: Your code here...
	return
}

// GetCurrentUser implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) GetCurrentUser(ctx context.Context, req *customer.GetCurrentUserReq) (resp *customer.GetCurrentUserResp, err error) {
	// TODO: Your code here...
	return
}

// Register implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) Register(ctx context.Context, req *customer.RegisterReq) (resp *customer.RegisterResp, err error) {
	// TODO: Your code here...
	return
}

// Logout implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) Logout(ctx context.Context, req *customer.LogoutReq) (resp *customer.LogoutResp, err error) {
	// TODO: Your code here...
	return
}

// EncryptMessage implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) EncryptMessage(ctx context.Context, req *customer.EncryptMessageReq) (resp *customer.EncryptMessageResp, err error) {
	// TODO: Your code here...
	return
}

// DecryptMessage implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) DecryptMessage(ctx context.Context, req *customer.DecryptMessageReq) (resp *customer.DecryptMessageResp, err error) {
	// TODO: Your code here...
	return
}

// DesensitizeMessage implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) DesensitizeMessage(ctx context.Context, req *customer.DesensitizeMessageReq) (resp *customer.DesensitizeMessageResp, err error) {
	// TODO: Your code here...
	return
}

// ArchiveConversations implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) ArchiveConversations(ctx context.Context, req *customer.ArchiveConversationsReq) (resp *customer.ArchiveConversationsResp, err error) {
	// TODO: Your code here...
	return
}

// GetArchiveTask implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) GetArchiveTask(ctx context.Context, req *customer.GetArchiveTaskReq) (resp *customer.GetArchiveTaskResp, err error) {
	// TODO: Your code here...
	return
}

// QueryArchivedConversation implements the CustomerServiceImpl interface.
func (s *CustomerServiceImpl) QueryArchivedConversation(ctx context.Context, req *customer.QueryArchivedConversationReq) (resp *customer.QueryArchivedConversationResp, err error) {
	// TODO: Your code here...
	return
}
