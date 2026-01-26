// Package handler 实现 Gateway 网关的 HTTP 请求处理器
// 该包负责：
// 1. 接收前端 HTTP 请求并进行参数校验
// 2. 将 HTTP 请求适配转换为 RPC 请求
// 3. 调用后端 Customer RPC 服务
// 4. 将 RPC 响应转换为 HTTP JSON 响应返回给前端
//
// 主要功能模块：
// - 客服信息管理（查询客服信息、客服列表）
// - 班次配置管理（创建、更新、删除、查询班次模板）
// - 排班管理（手动排班、自动排班、排班表格、Excel导出）
// - 请假调班（申请、审批、查询请假调班记录）
// - 会话管理（会话列表、历史会话、消息收发）
// - 会话分类与标签（分类管理、标签管理）
// - 统计看板（会话统计数据）
// - 用户认证（登录、注册、获取当前用户）
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"example_shop/gateway/config"
	"example_shop/gateway/middleware"
	"example_shop/gateway/rpc"
	"example_shop/service/customer/kitex_gen/customer"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xuri/excelize/v2"
)

// ============ 处理器结构体定义 ============

// CustomerHandler 网关客服服务处理器
// 作为 HTTP 请求与 RPC 服务之间的桥梁
// 负责处理前端 HTTP 请求并透传/适配到客服 RPC 服务
type CustomerHandler struct {
	client *rpc.CustomerClient // 客服 RPC 客户端，用于调用后端服务
}

// NewCustomerHandler 创建客服处理器实例
// 参数:
//   - client: 已初始化的 RPC 客户端
//
// 返回:
//   - *CustomerHandler: 处理器实例
func NewCustomerHandler(client *rpc.CustomerClient) *CustomerHandler {
	return &CustomerHandler{client: client}
}

// ============ 客服信息管理 ============

// GetCustomerService 获取单个客服详细信息
// 请求方式: GET
// 请求参数:
//   - cs_id: 客服ID（必填）
//
// 响应: 客服详细信息，包括姓名、部门、状态等
func (h *CustomerHandler) GetCustomerService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	csID := strings.TrimSpace(r.URL.Query().Get("cs_id"))
	if csID == "" {
		respondJSON(w, http.StatusBadRequest, &customer.GetCustomerServiceResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "cs_id is required"},
		})
		return
	}

	req := &customer.GetCustomerServiceReq{
		CsId: csID,
	}

	resp, err := h.client.GetCustomerService(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// ListCustomerService 分页查询客服列表
// 请求方式: GET
// 请求参数:
//   - dept_id: 部门ID（可选，筛选指定部门）
//   - page: 页码（默认1）
//   - page_size: 每页数量（默认10）
//
// 响应: 客服列表及分页信息
func (h *CustomerHandler) ListCustomerService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	deptID := strings.TrimSpace(r.URL.Query().Get("dept_id"))
	page, pageSize := parsePaginationParams(r, 10)

	req := &customer.ListCustomerServiceReq{
		DeptId:   deptID,
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := h.client.ListCustomerService(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// ============ 班次配置管理 ============

// CreateShiftConfig 创建班次配置模板
// 请求方式: POST
// 请求体参数:
//   - shift_name: 班次名称（必填，如"早班"、"晚班"）
//   - start_time: 开始时间（必填，格式 HH:MM:SS）
//   - end_time: 结束时间（必填，格式 HH:MM:SS）
//   - min_staff: 最少在班人数（必填，>=0）
//   - is_holiday: 是否节假日班次（0-否，1-是）
//   - create_by: 创建人
//
// 响应: 创建结果及新班次ID
func (h *CustomerHandler) CreateShiftConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 请求体参数结构
	var body struct {
		ShiftName string `json:"shift_name"` // 班次名称
		StartTime string `json:"start_time"` // 开始时间
		EndTime   string `json:"end_time"`   // 结束时间
		MinStaff  int32  `json:"min_staff"`  // 最少在班人数
		IsHoliday int8   `json:"is_holiday"` // 是否节假日班次
		CreateBy  string `json:"create_by"`  // 创建人
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.CreateShiftConfigResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid json body"},
		})
		return
	}
	body.ShiftName = strings.TrimSpace(body.ShiftName)
	body.StartTime = strings.TrimSpace(body.StartTime)
	body.EndTime = strings.TrimSpace(body.EndTime)
	body.CreateBy = strings.TrimSpace(body.CreateBy)
	if body.ShiftName == "" || body.StartTime == "" || body.EndTime == "" {
		respondJSON(w, http.StatusBadRequest, &customer.CreateShiftConfigResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "shift_name, start_time and end_time are required"},
		})
		return
	}
	if body.MinStaff < 0 {
		respondJSON(w, http.StatusBadRequest, &customer.CreateShiftConfigResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "min_staff must be >= 0"},
		})
		return
	}
	if body.IsHoliday != 0 && body.IsHoliday != 1 {
		respondJSON(w, http.StatusBadRequest, &customer.CreateShiftConfigResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "is_holiday must be 0 or 1"},
		})
		return
	}

	req := &customer.CreateShiftConfigReq{
		Shift: &customer.ShiftConfig{
			ShiftName: body.ShiftName,
			StartTime: body.StartTime,
			EndTime:   body.EndTime,
			MinStaff:  body.MinStaff,
			IsHoliday: body.IsHoliday,
			CreateBy:  body.CreateBy,
		},
	}
	resp, err := h.client.CreateShiftConfig(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ListShiftConfig 查询班次配置列表
// 接收 GET 请求，根据条件查询班次配置信息
func (h *CustomerHandler) ListShiftConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	isHolidayStr := r.URL.Query().Get("is_holiday")
	shiftName := strings.TrimSpace(r.URL.Query().Get("shift_name"))
	isHoliday := int64(-1)
	if isHolidayStr != "" {
		if v, err := strconv.ParseInt(isHolidayStr, 10, 8); err == nil {
			if v == 0 || v == 1 {
				isHoliday = v
			}
		}
	}

	req := &customer.ListShiftConfigReq{
		IsHoliday: int8(isHoliday),
		ShiftName: shiftName,
	}
	resp, err := h.client.ListShiftConfig(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// UpdateShiftConfig 更新班次配置
// 接收 POST 请求，更新现有班次模板的信息
func (h *CustomerHandler) UpdateShiftConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var reqBody struct {
		ShiftID   int64  `json:"shift_id"`
		ShiftName string `json:"shift_name"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		MinStaff  int32  `json:"min_staff"`
		IsHoliday int8   `json:"is_holiday"`
		CreateBy  string `json:"create_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.UpdateShiftConfigResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid json body"},
		})
		return
	}

	reqBody.ShiftName = strings.TrimSpace(reqBody.ShiftName)
	reqBody.StartTime = strings.TrimSpace(reqBody.StartTime)
	reqBody.EndTime = strings.TrimSpace(reqBody.EndTime)
	reqBody.CreateBy = strings.TrimSpace(reqBody.CreateBy)
	if reqBody.ShiftID <= 0 || reqBody.ShiftName == "" || reqBody.StartTime == "" || reqBody.EndTime == "" {
		respondJSON(w, http.StatusBadRequest, &customer.UpdateShiftConfigResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "shift_id, shift_name, start_time and end_time are required"},
		})
		return
	}
	if reqBody.MinStaff < 0 {
		respondJSON(w, http.StatusBadRequest, &customer.UpdateShiftConfigResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "min_staff must be >= 0"},
		})
		return
	}
	if reqBody.IsHoliday != 0 && reqBody.IsHoliday != 1 {
		respondJSON(w, http.StatusBadRequest, &customer.UpdateShiftConfigResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "is_holiday must be 0 or 1"},
		})
		return
	}

	req := &customer.UpdateShiftConfigReq{
		Shift: &customer.ShiftConfig{
			ShiftId:   reqBody.ShiftID,
			ShiftName: reqBody.ShiftName,
			StartTime: reqBody.StartTime,
			EndTime:   reqBody.EndTime,
			MinStaff:  reqBody.MinStaff,
			IsHoliday: reqBody.IsHoliday,
			CreateBy:  reqBody.CreateBy,
		},
	}
	resp, err := h.client.UpdateShiftConfig(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// DeleteShiftConfig 删除班次配置
// 接收 POST 请求，根据 ID 删除特定班次模板
func (h *CustomerHandler) DeleteShiftConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		ShiftID int64 `json:"shift_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.DeleteShiftConfigResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid json body"},
		})
		return
	}
	if body.ShiftID <= 0 {
		respondJSON(w, http.StatusBadRequest, &customer.DeleteShiftConfigResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "shift_id is required"},
		})
		return
	}

	resp, err := h.client.DeleteShiftConfig(r.Context(), &customer.DeleteShiftConfigReq{ShiftId: body.ShiftID})
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ============ 排班管理 ============

// AssignSchedule 手动批量分配排班
// 请求方式: POST
// 请求体参数:
//   - schedule_date: 排班日期（必填，格式 YYYY-MM-DD）
//   - shift_id: 班次ID（必填）
//   - cs_ids: 客服ID列表（必填，批量分配）
//   - create_by: 操作人
//
// 响应: 排班分配结果
func (h *CustomerHandler) AssignSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 请求体参数结构
	var body struct {
		ScheduleDate string   `json:"schedule_date"` // 排班日期
		ShiftID      int64    `json:"shift_id"`      // 班次ID
		CsIDs        []string `json:"cs_ids"`        // 客服ID列表
		CreateBy     string   `json:"create_by"`     // 创建人
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.AssignScheduleResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid json body"},
		})
		return
	}
	body.ScheduleDate = strings.TrimSpace(body.ScheduleDate)
	body.CreateBy = strings.TrimSpace(body.CreateBy)
	if body.ScheduleDate == "" || body.ShiftID <= 0 || len(body.CsIDs) == 0 {
		respondJSON(w, http.StatusBadRequest, &customer.AssignScheduleResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "schedule_date, shift_id and cs_ids are required"},
		})
		return
	}
	if _, err := time.Parse("2006-01-02", body.ScheduleDate); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.AssignScheduleResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "schedule_date must be YYYY-MM-DD"},
		})
		return
	}

	req := &customer.AssignScheduleReq{
		ScheduleDate: body.ScheduleDate,
		ShiftId:      body.ShiftID,
		CsIds:        body.CsIDs,
		CreateBy:     body.CreateBy,
	}
	resp, err := h.client.AssignSchedule(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// AutoSchedule 自动排班
func (h *CustomerHandler) AutoSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		StartDate  string `json:"start_date"`
		EndDate    string `json:"end_date"`
		DeptID     string `json:"dept_id"`
		TeamID     string `json:"team_id"`
		OperatorID string `json:"operator_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.AutoScheduleResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid json body"},
		})
		return
	}
	body.StartDate = strings.TrimSpace(body.StartDate)
	body.EndDate = strings.TrimSpace(body.EndDate)
	body.OperatorID = strings.TrimSpace(body.OperatorID)
	if body.StartDate == "" || body.EndDate == "" {
		respondJSON(w, http.StatusBadRequest, &customer.AutoScheduleResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "start_date and end_date are required"},
		})
		return
	}
	if _, err := time.Parse("2006-01-02", body.StartDate); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.AutoScheduleResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "start_date must be YYYY-MM-DD"},
		})
		return
	}
	if _, err := time.Parse("2006-01-02", body.EndDate); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.AutoScheduleResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "end_date must be YYYY-MM-DD"},
		})
		return
	}
	if body.OperatorID == "" {
		body.OperatorID = "ADMIN"
	}

	req := &customer.AutoScheduleReq{
		StartDate:  body.StartDate,
		EndDate:    body.EndDate,
		DeptId:     body.DeptID,
		TeamId:     body.TeamID,
		OperatorId: body.OperatorID,
	}
	resp, err := h.client.AutoSchedule(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ============ 请假调班管理 ============

// ApplyLeaveTransfer 申请请假或调班
// 请求方式: POST
// 请求体参数:
//   - cs_id: 申请人客服ID（必填）
//   - apply_type: 申请类型（0-请假，1-调班）
//   - start_date: 开始日期（必填，格式 YYYY-MM-DD）
//   - end_date: 结束日期（必填，格式 YYYY-MM-DD）
//   - start_period: 开始时段（0-全天，1-上午，2-下午）
//   - end_period: 结束时段
//   - shift_id: 班次ID（调班时必填）
//   - target_cs_id: 调班目标客服ID（调班时必填）
//   - reason: 申请原因
//
// 响应: 申请提交结果及申请单ID
func (h *CustomerHandler) ApplyLeaveTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 从token获取操作人信息
	operatorInfo := getOperatorInfoFromContext(r.Context())
	if operatorInfo.ID == "" {
		respondJSON(w, http.StatusOK, &customer.ApplyLeaveTransferResp{
			BaseResp: &customer.BaseResp{Code: 401, Msg: "请先登录"},
		})
		return
	}

	// 请求体参数结构
	var body struct {
		ApplyType   int8   `json:"apply_type"`   // 申请类型：0-请假，1-调班
		StartDate   string `json:"start_date"`   // 开始日期
		EndDate     string `json:"end_date"`     // 结束日期
		StartPeriod int8   `json:"start_period"` // 开始时段
		EndPeriod   int8   `json:"end_period"`   // 结束时段
		ShiftID     int64  `json:"shift_id"`     // 班次ID
		TargetCsID  string `json:"target_cs_id"` // 调班目标客服ID
		Reason      string `json:"reason"`       // 申请原因
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.ApplyLeaveTransferResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid json body"},
		})
		return
	}
	body.StartDate = strings.TrimSpace(body.StartDate)
	body.EndDate = strings.TrimSpace(body.EndDate)
	body.TargetCsID = strings.TrimSpace(body.TargetCsID)
	body.Reason = strings.TrimSpace(body.Reason)

	if body.StartDate == "" {
		respondJSON(w, http.StatusBadRequest, &customer.ApplyLeaveTransferResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "start_date is required"},
		})
		return
	}
	if body.EndDate == "" {
		body.EndDate = body.StartDate
	}

	if _, err := time.Parse("2006-01-02", body.StartDate); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.ApplyLeaveTransferResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "start_date must be YYYY-MM-DD"},
		})
		return
	}
	if _, err := time.Parse("2006-01-02", body.EndDate); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.ApplyLeaveTransferResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "end_date must be YYYY-MM-DD"},
		})
		return
	}

	if body.ApplyType != 0 && body.ApplyType != 1 {
		respondJSON(w, http.StatusBadRequest, &customer.ApplyLeaveTransferResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "apply_type must be 0 or 1"},
		})
		return
	}
	if body.ApplyType == 1 && (body.TargetCsID == "" || body.ShiftID <= 0) {
		respondJSON(w, http.StatusBadRequest, &customer.ApplyLeaveTransferResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "target_cs_id and shift_id are required for transfer"},
		})
		return
	}

	// 使用token中的操作人信息作为申请人
	req := &customer.ApplyLeaveTransferReq{
		CsId:        operatorInfo.ID,
		ApplyType:   body.ApplyType,
		StartDate:   body.StartDate,
		EndDate:     body.EndDate,
		StartPeriod: body.StartPeriod,
		EndPeriod:   body.EndPeriod,
		ShiftId:     body.ShiftID,
		TargetCsId:  body.TargetCsID,
		Reason:      body.Reason,
		TargetDate:  body.StartDate, // 兼容旧字段
	}
	resp, err := h.client.ApplyLeaveTransfer(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ApproveLeaveTransfer 审批请假调班申请
// 请求方式: POST
// 请求体参数:
//   - apply_id: 申请单ID（必填）
//   - approval_status: 审批状态（1-通过，2-拒绝）
//   - approver_id: 审批人ID（必填）
//   - approver_name: 审批人姓名
//   - approval_remark: 审批备注
//
// 响应: 审批结果
func (h *CustomerHandler) ApproveLeaveTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 从token获取操作人信息
	operatorInfo := getOperatorInfoFromContext(r.Context())
	if operatorInfo.ID == "" {
		respondJSON(w, http.StatusOK, &customer.ApproveLeaveTransferResp{
			BaseResp: &customer.BaseResp{Code: 401, Msg: "请先登录"},
		})
		return
	}

	// 请求体参数结构
	var body struct {
		ApplyID        int64  `json:"apply_id"`        // 申请单ID
		ApprovalStatus int8   `json:"approval_status"` // 审批状态
		ApprovalRemark string `json:"approval_remark"` // 审批备注
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.ApproveLeaveTransferResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid json body"},
		})
		return
	}
	if body.ApplyID <= 0 {
		respondJSON(w, http.StatusBadRequest, &customer.ApproveLeaveTransferResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "apply_id is required"},
		})
		return
	}
	if body.ApprovalStatus != 1 && body.ApprovalStatus != 2 {
		respondJSON(w, http.StatusBadRequest, &customer.ApproveLeaveTransferResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "approval_status must be 1 or 2"},
		})
		return
	}

	// 使用token中的操作人信息
	req := &customer.ApproveLeaveTransferReq{
		ApplyId:        body.ApplyID,
		ApprovalStatus: body.ApprovalStatus,
		ApproverId:     operatorInfo.ID,
		ApproverName:   operatorInfo.Name,
		ApprovalRemark: body.ApprovalRemark,
	}
	resp, err := h.client.ApproveLeaveTransfer(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ApplyChainSwap 提交链式调班申请
func (h *CustomerHandler) ApplyChainSwap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 从token获取操作人信息
	operatorInfo := getOperatorInfoFromContext(r.Context())
	if operatorInfo.ID == "" {
		respondJSON(w, http.StatusOK, map[string]interface{}{"code": 401, "msg": "请先登录"})
		return
	}

	var body struct {
		DeptID string `json:"dept_id"`
		Reason string `json:"reason"`
		Items  []struct {
			CsID           string `json:"cs_id"`
			FromScheduleID int64  `json:"from_schedule_id"`
			ToScheduleID   int64  `json:"to_schedule_id"`
			Step           int32  `json:"step"`
		} `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{"code": 400, "msg": "invalid json"})
		return
	}

	// 使用token中的操作人信息作为申请人
	req := &customer.ApplyChainSwapReq{
		ApplicantId: operatorInfo.ID,
		DeptId:      body.DeptID,
		Reason:      body.Reason,
	}
	for _, it := range body.Items {
		req.Items = append(req.Items, &customer.ChainSwapItem{
			CsId:           it.CsID,
			FromScheduleId: it.FromScheduleID,
			ToScheduleId:   it.ToScheduleID,
			Step:           it.Step,
		})
	}

	resp, err := h.client.ApplyChainSwap(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 500, "msg": err.Error()})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ApproveChainSwap 审批链式调班申请
func (h *CustomerHandler) ApproveChainSwap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 从token获取操作人信息
	operatorInfo := getOperatorInfoFromContext(r.Context())
	if operatorInfo.ID == "" {
		respondJSON(w, http.StatusOK, map[string]interface{}{"code": 401, "msg": "请先登录"})
		return
	}

	var body struct {
		RequestID      int64  `json:"request_id"`
		ApprovalStatus int8   `json:"approval_status"`
		ApprovalRemark string `json:"approval_remark"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{"code": 400, "msg": "invalid json"})
		return
	}

	// 使用token中的操作人信息
	req := &customer.ApproveChainSwapReq{
		RequestId:      body.RequestID,
		ApprovalStatus: body.ApprovalStatus,
		ApproverId:     operatorInfo.ID,
		ApproverName:   operatorInfo.Name,
		ApprovalRemark: body.ApprovalRemark,
	}

	resp, err := h.client.ApproveChainSwap(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 500, "msg": err.Error()})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ListChainSwap 查询链式调班申请列表
// 请求方式: GET
// 请求参数:
//   - status: 状态筛选（-1=全部, 0=待审批, 1=已通过, 2=已拒绝）
//   - keyword: 关键词搜索
//   - page: 页码（默认1）
//   - page_size: 每页数量（默认20）
func (h *CustomerHandler) ListChainSwap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	status := int8(-1)
	if v := strings.TrimSpace(r.URL.Query().Get("status")); v != "" {
		n, err := strconv.ParseInt(v, 10, 8)
		if err == nil {
			status = int8(n)
		}
	}
	keyword := strings.TrimSpace(r.URL.Query().Get("keyword"))

	page, pageSize := parsePaginationParams(r, 20)

	req := &customer.ListChainSwapReq{
		Status:   status,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	resp, err := h.client.ListChainSwap(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// GetChainSwap 获取链式调班申请详情
// 请求方式: GET
// 请求参数:
//   - swap_id: 链式调班申请ID（必填）
func (h *CustomerHandler) GetChainSwap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	swapIDStr := strings.TrimSpace(r.URL.Query().Get("swap_id"))
	swapID, err := strconv.ParseInt(swapIDStr, 10, 64)
	if err != nil || swapID <= 0 {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "swap_id is required",
		})
		return
	}

	req := &customer.GetChainSwapReq{SwapId: swapID}
	resp, err := h.client.GetChainSwap(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// GetLeaveTransfer 获取单个请假调班申请详情
// 请求方式: GET
// 请求参数:
//   - apply_id: 申请单ID（必填）
//
// 响应: 申请详情，包括申请人、类型、状态等
func (h *CustomerHandler) GetLeaveTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	applyIDStr := strings.TrimSpace(r.URL.Query().Get("apply_id"))
	applyID, err := strconv.ParseInt(applyIDStr, 10, 64)
	if err != nil || applyID <= 0 {
		respondJSON(w, http.StatusBadRequest, &customer.GetLeaveTransferResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "apply_id is required"},
		})
		return
	}

	req := &customer.GetLeaveTransferReq{ApplyId: applyID}
	resp, err := h.client.GetLeaveTransfer(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ListLeaveTransfer 分页查询请假调班申请列表
// 请求方式: GET
// 请求参数:
//   - approval_status: 审批状态筛选（-1-全部，0-待审批，1-已通过，2-已拒绝）
//   - keyword: 关键词搜索（搜索申请人姓名）
//   - page: 页码（默认1）
//   - page_size: 每页数量（默认20）
//
// 响应: 申请列表及分页信息
func (h *CustomerHandler) ListLeaveTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	approvalStatus := int8(-1)
	if v := strings.TrimSpace(r.URL.Query().Get("approval_status")); v != "" {
		n, err := strconv.ParseInt(v, 10, 8)
		if err == nil {
			approvalStatus = int8(n)
		}
	}
	keyword := strings.TrimSpace(r.URL.Query().Get("keyword"))

	page, pageSize := parsePaginationParams(r, 20)

	req := &customer.ListLeaveTransferReq{
		ApprovalStatus: approvalStatus,
		Keyword:        keyword,
		Page:           page,
		PageSize:       pageSize,
		OperatorId:     resolveCustomerCsID(r.Context(), ""), // 获取当前操作人ID
	}
	resp, err := h.client.ListLeaveTransfer(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// GetLeaveAuditLog 获取请假/调班申请的审计日志
// 请求方式: GET
// 请求参数:
//   - apply_id: 申请单ID（必填）
//
// 响应: 审计日志列表
func (h *CustomerHandler) GetLeaveAuditLog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	applyIDStr := strings.TrimSpace(r.URL.Query().Get("apply_id"))
	applyID, err := strconv.ParseInt(applyIDStr, 10, 64)
	if err != nil || applyID <= 0 {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "apply_id is required",
		})
		return
	}

	req := &customer.GetLeaveAuditLogReq{ApplyId: applyID}
	resp, err := h.client.GetLeaveAuditLog(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ListScheduleGrid 查询排班表格数据

func (h *CustomerHandler) ListScheduleGrid(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	startDate := strings.TrimSpace(r.URL.Query().Get("start_date"))
	endDate := strings.TrimSpace(r.URL.Query().Get("end_date"))
	deptID := strings.TrimSpace(r.URL.Query().Get("dept_id"))
	teamID := strings.TrimSpace(r.URL.Query().Get("team_id"))
	if startDate == "" || endDate == "" {
		respondJSON(w, http.StatusBadRequest, &customer.ListScheduleGridResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "start_date and end_date are required"},
		})
		return
	}

	req := &customer.ListScheduleGridReq{
		StartDate: startDate,
		EndDate:   endDate,
		DeptId:    deptID,
		TeamId:    teamID,
	}

	// 如果是客服角色，只能查看自己的排班
	operator := getOperatorInfoFromContext(r.Context())
	if operator.Role == "customer_service" {
		req.CsId = operator.ID
	}

	resp, err := h.client.ListScheduleGrid(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ExportScheduleExcel 导出排班表为Excel文件
// 请求方式: GET
// 请求参数:
//   - start_date: 开始日期（必填，格式 YYYY-MM-DD）
//   - end_date: 结束日期（必填，格式 YYYY-MM-DD）
//   - dept_id: 部门ID（可选）
//   - team_id: 小组ID（可选）
//
// 响应: Excel 文件下载流（Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet）
// 文件名格式: schedule_开始日期_结束日期.xlsx
func (h *CustomerHandler) ExportScheduleExcel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 解析查询参数
	startDate := strings.TrimSpace(r.URL.Query().Get("start_date"))
	endDate := strings.TrimSpace(r.URL.Query().Get("end_date"))
	deptID := strings.TrimSpace(r.URL.Query().Get("dept_id"))
	teamID := strings.TrimSpace(r.URL.Query().Get("team_id"))
	if startDate == "" || endDate == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "start_date and end_date are required",
		})
		return
	}

	grid, err := h.client.ListScheduleGrid(r.Context(), &customer.ListScheduleGridReq{
		StartDate: startDate,
		EndDate:   endDate,
		DeptId:    deptID,
		TeamId:    teamID,
	})
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	if grid == nil || grid.BaseResp == nil || grid.BaseResp.Code != 0 {
		code := 500
		msg := "export failed"
		if grid != nil && grid.BaseResp != nil {
			code = int(grid.BaseResp.Code)
			msg = grid.BaseResp.Msg
		}
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"code": code,
			"msg":  msg,
		})
		return
	}

	shiftNameByID := map[int64]string{}
	for _, s := range grid.Shifts {
		if s == nil {
			continue
		}
		shiftNameByID[s.ShiftId] = s.ShiftName
	}

	cellByKey := map[string]*customer.ScheduleCell{}
	for _, c := range grid.Cells {
		if c == nil {
			continue
		}
		cellByKey[c.CsId+"|"+c.ScheduleDate] = c
	}

	f := excelize.NewFile()
	sheet := "排班表"
	f.SetSheetName("Sheet1", sheet)

	header := make([]interface{}, 0, 2+len(grid.Dates))
	header = append(header, "客服ID", "客服姓名")
	for _, d := range grid.Dates {
		header = append(header, d)
	}
	if err := f.SetSheetRow(sheet, "A1", &header); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "export failed: " + err.Error(),
		})
		return
	}

	for i, cs := range grid.Customers {
		if cs == nil {
			continue
		}
		row := make([]interface{}, 0, 2+len(grid.Dates))
		row = append(row, cs.CsId, cs.CsName)
		for _, d := range grid.Dates {
			cell := cellByKey[cs.CsId+"|"+d]
			if cell == nil || cell.ShiftId <= 0 {
				row = append(row, "")
				continue
			}
			name := shiftNameByID[cell.ShiftId]
			if name == "" {
				name = fmt.Sprintf("班次%d", cell.ShiftId)
			}
			if cell.Status == 1 {
				name += "(请假)"
			} else if cell.Status == 2 {
				name += "(调班)"
			}
			row = append(row, name)
		}
		addr, _ := excelize.CoordinatesToCellName(1, i+2)
		if err := f.SetSheetRow(sheet, addr, &row); err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code": 500,
				"msg":  "export failed: " + err.Error(),
			})
			return
		}
	}

	f.SetColWidth(sheet, "A", "A", 14)
	f.SetColWidth(sheet, "B", "B", 16)
	if len(grid.Dates) > 0 {
		lastCol, _ := excelize.ColumnNumberToName(2 + len(grid.Dates))
		_ = f.SetColWidth(sheet, "C", lastCol, 14)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "export failed: " + err.Error(),
		})
		return
	}

	fileName := fmt.Sprintf("schedule_%s_%s.xlsx", startDate, endDate)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

// UpsertScheduleCell 更新/清空排班单元格（shift_id=0 清空）
func (h *CustomerHandler) UpsertScheduleCell(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		CsID         string `json:"cs_id"`
		ScheduleDate string `json:"schedule_date"`
		ShiftID      int64  `json:"shift_id"`
		OperatorID   string `json:"operator_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.UpsertScheduleCellResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid json body"},
		})
		return
	}
	body.CsID = strings.TrimSpace(body.CsID)
	body.ScheduleDate = strings.TrimSpace(body.ScheduleDate)
	body.OperatorID = strings.TrimSpace(body.OperatorID)
	if body.CsID == "" || body.ScheduleDate == "" {
		respondJSON(w, http.StatusBadRequest, &customer.UpsertScheduleCellResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "cs_id and schedule_date are required"},
		})
		return
	}
	if body.OperatorID == "" {
		body.OperatorID = "ADMIN"
	}

	req := &customer.UpsertScheduleCellReq{
		CsId:         body.CsID,
		ScheduleDate: body.ScheduleDate,
		ShiftId:      body.ShiftID,
		OperatorId:   body.OperatorID,
	}
	resp, err := h.client.UpsertScheduleCell(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ============ 会话管理 ============

// AssignCustomer 自动分配客服
// 为用户自动分配当前在线且负载最低的客服
// 请求方式: POST
// 请求体参数:
//   - user_id: 用户ID
//   - user_nickname: 用户昵称
//   - source: 来源渠道
//
// 响应: 分配的客服信息和新创建的会话ID
func (h *CustomerHandler) AssignCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		UserID       string `json:"user_id"`
		UserNickname string `json:"user_nickname"`
		Source       string `json:"source"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.AssignCustomerResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid json body"},
		})
		return
	}

	req := &customer.AssignCustomerReq{
		UserId:       body.UserID,
		UserNickname: body.UserNickname,
		Source:       body.Source,
	}

	resp, err := h.client.AssignCustomer(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, &customer.AssignCustomerResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: err.Error()},
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ListConversation 分页查询当前会话列表
// 查询客服的当前进行中会话（状态为进行中的会话）
// 请求方式: GET
// 请求参数:
//   - cs_id: 客服ID（可选，不填则查询全部）
//   - keyword: 关键词搜索（搜索用户昵称/ID）
//   - status: 会话状态筛选（-1-全部，0-进行中，1-已结束，2-已转接）
//   - page: 页码（默认1）
//   - page_size: 每页数量（默认20）
//
// 响应: 会话列表及分页信息、未读消息数
func (h *CustomerHandler) ListConversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	csID := strings.TrimSpace(r.URL.Query().Get("cs_id"))
	csID = resolveCustomerCsID(r.Context(), csID)
	keyword := strings.TrimSpace(r.URL.Query().Get("keyword"))
	statusStr := strings.TrimSpace(r.URL.Query().Get("status"))

	status := int64(-1)
	if statusStr != "" {
		if v, err := strconv.ParseInt(statusStr, 10, 8); err == nil {
			status = v
		}
	}

	page, pageSize := parsePaginationParams(r, 20)

	req := &customer.ListConversationReq{
		CsId:     csID,
		Keyword:  keyword,
		Status:   int8(status),
		Page:     page,
		PageSize: pageSize,
	}
	resp, err := h.client.ListConversation(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ListConversationHistory 查询会话历史记录列表
// 支持分页查询，可按客服ID、关键词、状态筛选
// 仅返回已结束或已转接的会话，且必须包含用户发送的消息
func (h *CustomerHandler) ListConversationHistory(w http.ResponseWriter, r *http.Request) {
	// 仅允许 GET 请求
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 解析查询参数
	csID := strings.TrimSpace(r.URL.Query().Get("cs_id"))
	csID = resolveCustomerCsID(r.Context(), csID)
	keyword := strings.TrimSpace(r.URL.Query().Get("keyword"))
	statusStr := strings.TrimSpace(r.URL.Query().Get("status"))

	// 状态参数处理：默认为-1（不筛选）
	status := int64(-1)
	if statusStr != "" {
		if v, err := strconv.ParseInt(statusStr, 10, 8); err == nil {
			status = v
		}
	}

	// 分页参数处理：默认第1页，每页20条
	page, pageSize := parsePaginationParams(r, 20)

	// 构建 RPC 请求对象
	req := &customer.ListConversationHistoryReq{
		CsId:     csID,
		Keyword:  keyword,
		Status:   int8(status),
		Page:     page,
		PageSize: pageSize,
	}
	// 调用 RPC 服务查询历史会话
	resp, err := h.client.ListConversationHistory(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ListConversationMessage 查询会话消息
// 支持分页查询，可指定排序方式（正序/倒序）
func (h *CustomerHandler) ListConversationMessage(w http.ResponseWriter, r *http.Request) {
	// 仅允许 GET 请求
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 校验必填参数 conv_id
	convID := strings.TrimSpace(r.URL.Query().Get("conv_id"))
	if convID == "" {
		respondJSON(w, http.StatusBadRequest, &customer.ListConversationMessageResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "conv_id is required"},
		})
		return
	}

	// 解析分页与排序参数
	orderAscStr := strings.TrimSpace(r.URL.Query().Get("order_asc"))

	page, pageSize := parsePaginationParams(r, 50)
	orderAsc := int64(0) // 默认倒序(0)
	if orderAscStr != "" {
		if v, err := strconv.ParseInt(orderAscStr, 10, 8); err == nil {
			orderAsc = v
		}
	}

	// 构建 RPC 请求对象
	req := &customer.ListConversationMessageReq{
		ConvId:   convID,
		Page:     page,
		PageSize: pageSize,
		OrderAsc: int8(orderAsc),
	}
	// 调用 RPC 服务查询消息
	resp, err := h.client.ListConversationMessage(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// SendConversationMessage 发送会话消息
// 客服或用户发送会话消息，支持普通消息和快捷回复
// 请求方式: POST
// 请求体参数:
//   - conv_id: 会话ID（必填）
//   - sender_type: 发送方类型（0-用户，1-客服）
//   - sender_id: 发送方ID
//   - msg_content: 消息内容（必填）
//   - is_quick_reply: 是否快捷回复（0-否，1-是）
//   - quick_reply_id: 快捷回复ID（快捷回复时必填）
//
// 响应: 发送结果及消息ID
func (h *CustomerHandler) SendConversationMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		ConvID       string `json:"conv_id"`
		SenderType   int8   `json:"sender_type"`
		SenderID     string `json:"sender_id"`
		MsgContent   string `json:"msg_content"`
		IsQuickReply int8   `json:"is_quick_reply"`
		QuickReplyID int64  `json:"quick_reply_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.SendConversationMessageResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid json body"},
		})
		return
	}
	body.ConvID = strings.TrimSpace(body.ConvID)
	body.SenderID = strings.TrimSpace(body.SenderID)
	body.MsgContent = strings.TrimSpace(body.MsgContent)
	if body.ConvID == "" || body.MsgContent == "" {
		respondJSON(w, http.StatusBadRequest, &customer.SendConversationMessageResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "conv_id and msg_content are required"},
		})
		return
	}
	if body.SenderID == "" {
		body.SenderID = "KF001"
	}

	req := &customer.SendConversationMessageReq{
		ConvId:       body.ConvID,
		SenderType:   body.SenderType,
		SenderId:     body.SenderID,
		MsgContent:   body.MsgContent,
		IsQuickReply: body.IsQuickReply,
		QuickReplyId: body.QuickReplyID,
	}
	resp, err := h.client.SendConversationMessage(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// CreateConversation 创建会话
// 用户发起新会话，可指定客服或自动分配
// 请求方式: POST
// 请求体参数:
//   - user_id: 用户ID（必填）
//   - user_nickname: 用户昵称
//   - source: 来源渠道（APP/Web/H5/WeChat）
//   - cs_id: 指定客服ID（可选，为空则自动分配）
//   - first_msg: 首条消息（可选）
//
// 响应: 会话eID、客服信息、是否新创建
func (h *CustomerHandler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		UserID       string `json:"user_id"`
		UserNickname string `json:"user_nickname"`
		Source       string `json:"source"`
		CsID         string `json:"cs_id"`
		FirstMsg     string `json:"first_msg"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400, "msg": "invalid json body",
		})
		return
	}
	body.UserID = strings.TrimSpace(body.UserID)
	body.UserNickname = strings.TrimSpace(body.UserNickname)
	body.Source = strings.TrimSpace(body.Source)
	body.CsID = strings.TrimSpace(body.CsID)
	body.FirstMsg = strings.TrimSpace(body.FirstMsg)

	if body.UserID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400, "msg": "user_id is required",
		})
		return
	}

	req := &customer.CreateConversationReq{
		UserId:       body.UserID,
		UserNickname: body.UserNickname,
		Source:       body.Source,
		CsId:         body.CsID,
		FirstMsg:     body.FirstMsg,
	}
	resp, err := h.client.CreateConversation(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500, "msg": "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// EndConversation 结束会话
// 客服或系统主动结束会话
// 请求方式: POST
// 请求体参数:
//   - conv_id: 会话ID（必填）
//   - operator_id: 操作人（客服ID或系统）
//   - end_reason: 结束原因（可选）
//
// 响应: 结束结果及会话时长
func (h *CustomerHandler) EndConversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		ConvID     string `json:"conv_id"`
		OperatorID string `json:"operator_id"`
		EndReason  string `json:"end_reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400, "msg": "invalid json body",
		})
		return
	}
	body.ConvID = strings.TrimSpace(body.ConvID)
	body.OperatorID = strings.TrimSpace(body.OperatorID)
	body.EndReason = strings.TrimSpace(body.EndReason)

	if body.ConvID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400, "msg": "conv_id is required",
		})
		return
	}

	req := &customer.EndConversationReq{
		ConvId:     body.ConvID,
		OperatorId: body.OperatorID,
		EndReason:  body.EndReason,
	}
	resp, err := h.client.EndConversation(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500, "msg": "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// TransferConversation 转接会话
// 将会话从当前客服转接给另一位客服
// 请求方式: POST
// 请求体参数:
//   - conv_id: 会话ID（必填）
//   - from_cs_id: 转出客服ID（必填）
//   - to_cs_id: 转入客服ID（必填）
//   - transfer_reason: 转接原因（可选）
//   - context_remark: 上下文备注（可选，JSON格式）
//
// 响应: 转接结果及转接记录ID
func (h *CustomerHandler) TransferConversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		ConvID         string `json:"conv_id"`
		FromCsID       string `json:"from_cs_id"`
		ToCsID         string `json:"to_cs_id"`
		TransferReason string `json:"transfer_reason"`
		ContextRemark  string `json:"context_remark"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400, "msg": "invalid json body",
		})
		return
	}
	body.ConvID = strings.TrimSpace(body.ConvID)
	body.FromCsID = strings.TrimSpace(body.FromCsID)
	body.ToCsID = strings.TrimSpace(body.ToCsID)
	body.TransferReason = strings.TrimSpace(body.TransferReason)

	if body.ConvID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400, "msg": "conv_id is required",
		})
		return
	}
	if body.FromCsID == "" || body.ToCsID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400, "msg": "from_cs_id and to_cs_id are required",
		})
		return
	}

	req := &customer.TransferConversationReq{
		ConvId:         body.ConvID,
		FromCsId:       body.FromCsID,
		ToCsId:         body.ToCsID,
		TransferReason: body.TransferReason,
		ContextRemark:  body.ContextRemark,
	}
	resp, err := h.client.TransferConversation(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500, "msg": "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ListQuickReply 分页查询快捷回复列表
// 查询系统预设和自定义的快捷回复语
// 请求方式: GET
// 请求参数:
//   - keyword: 关键词搜索（搜索快捷回复内容）
//   - reply_type: 回复类型筛选（-1-全部）
//   - is_public: 是否公开（-1-全部，0-私有，1-公开）
//   - page: 页码（默认1）
//   - page_size: 每页数量（默认50）
//
// 响应: 快捷回复列表及分页信息
func (h *CustomerHandler) ListQuickReply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	keyword := strings.TrimSpace(r.URL.Query().Get("keyword"))
	replyTypeStr := strings.TrimSpace(r.URL.Query().Get("reply_type"))
	isPublicStr := strings.TrimSpace(r.URL.Query().Get("is_public"))

	replyType := int64(-1)
	if replyTypeStr != "" {
		if v, err := strconv.ParseInt(replyTypeStr, 10, 8); err == nil {
			replyType = v
		}
	}
	isPublic := int64(-1)
	if isPublicStr != "" {
		if v, err := strconv.ParseInt(isPublicStr, 10, 8); err == nil {
			isPublic = v
		}
	}

	page, pageSize := parsePaginationParams(r, 50)

	req := &customer.ListQuickReplyReq{
		Keyword:   keyword,
		ReplyType: int8(replyType),
		IsPublic:  int8(isPublic),
		Page:      page,
		PageSize:  pageSize,
	}
	resp, err := h.client.ListQuickReply(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ============ 会话分类管理 ============

// CreateConvCategory 新增会话分类
// 创建一个新的会话分类（如"用户咨询"、"投诉建议"等）
// 请求方式: POST
// 请求体参数:
//   - category_name: 分类名称（必填）
//   - sort_no: 排序号（越小越前）
//   - create_by: 创建人
//
// 响应: 创建结果及新分类ID
func (h *CustomerHandler) CreateConvCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		CategoryName string `json:"category_name"`
		SortNo       int32  `json:"sort_no"`
		CreateBy     string `json:"create_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.CreateConvCategoryResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid json body"},
		})
		return
	}
	body.CategoryName = strings.TrimSpace(body.CategoryName)
	body.CreateBy = strings.TrimSpace(body.CreateBy)
	if body.CategoryName == "" {
		respondJSON(w, http.StatusBadRequest, &customer.CreateConvCategoryResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "category_name is required"},
		})
		return
	}
	if body.CreateBy == "" {
		body.CreateBy = "ADMIN"
	}

	req := &customer.CreateConvCategoryReq{
		CategoryName: body.CategoryName,
		SortNo:       body.SortNo,
		CreateBy:     body.CreateBy,
	}
	resp, err := h.client.CreateConvCategory(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ListConvCategory 查询所有会话分类
// 获取系统中所有已启用的会话分类，按排序号排序
// 请求方式: GET
// 请求参数: 无
//
// 响应: 分类列表
func (h *CustomerHandler) ListConvCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	req := &customer.ListConvCategoryReq{}
	resp, err := h.client.ListConvCategory(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// UpdateConversationClassify 更新会话分类、标签和核心标记
// 为指定会话设置/更新分类、标签和核心会话标记
// 请求方式: POST
// 请求体参数:
//   - conv_id: 会话ID（必填）
//   - category_id: 分类ID（0表示不更新）
//   - tags: 标签（逗号分隔的标签ID，空字符串清除标签）
//   - is_core: 是否核心会话（0-否，1-是）
//   - operator_id: 操作人ID
//
// 响应: 更新结果
func (h *CustomerHandler) UpdateConversationClassify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		ConvID     string `json:"conv_id"`
		CategoryID int64  `json:"category_id"`
		Tags       string `json:"tags"`
		IsCore     int8   `json:"is_core"`
		OperatorID string `json:"operator_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.UpdateConversationClassifyResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid json body"},
		})
		return
	}
	body.ConvID = strings.TrimSpace(body.ConvID)
	body.Tags = strings.TrimSpace(body.Tags)
	body.OperatorID = strings.TrimSpace(body.OperatorID)
	if body.ConvID == "" {
		respondJSON(w, http.StatusBadRequest, &customer.UpdateConversationClassifyResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "conv_id is required"},
		})
		return
	}
	if body.OperatorID == "" {
		body.OperatorID = "ADMIN"
	}

	req := &customer.UpdateConversationClassifyReq{
		ConvId:     body.ConvID,
		CategoryId: body.CategoryID,
		Tags:       body.Tags,
		IsCore:     body.IsCore,
		OperatorId: body.OperatorID,
	}
	resp, err := h.client.UpdateConversationClassify(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ============ HTTP 响应工具函数 ============

// respondJSON 返回JSON格式的HTTP响应
// 参数:
//   - w: HTTP响应写入器
//   - status: HTTP状态码
//   - data: 响应数据（将被序列化为JSON）
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError 返回错误响应
// 参数:
//   - w: HTTP响应写入器
//   - status: HTTP状态码
//   - message: 错误消息
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]interface{}{
		"code": status,
		"msg":  message,
	})
}

// ============ 会话标签管理 ============

// CreateConvTag 创建会话标签
func (h *CustomerHandler) CreateConvTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		TagName  string `json:"tag_name"`
		TagColor string `json:"tag_color"`
		SortNo   int32  `json:"sort_no"`
		CreateBy string `json:"create_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.CreateConvTagResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid json body"},
		})
		return
	}
	body.TagName = strings.TrimSpace(body.TagName)
	body.TagColor = strings.TrimSpace(body.TagColor)
	body.CreateBy = strings.TrimSpace(body.CreateBy)
	if body.TagName == "" {
		respondJSON(w, http.StatusBadRequest, &customer.CreateConvTagResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "tag_name is required"},
		})
		return
	}

	req := &customer.CreateConvTagReq{
		TagName:  body.TagName,
		TagColor: body.TagColor,
		SortNo:   body.SortNo,
		CreateBy: body.CreateBy,
	}
	resp, err := h.client.CreateConvTag(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ListConvTag 查询会话标签列表
func (h *CustomerHandler) ListConvTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	resp, err := h.client.ListConvTag(r.Context(), &customer.ListConvTagReq{})
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// UpdateConvTag 更新会话标签
func (h *CustomerHandler) UpdateConvTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		TagID    int64  `json:"tag_id"`
		TagName  string `json:"tag_name"`
		TagColor string `json:"tag_color"`
		SortNo   int32  `json:"sort_no"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.UpdateConvTagResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid json body"},
		})
		return
	}
	if body.TagID <= 0 {
		respondJSON(w, http.StatusBadRequest, &customer.UpdateConvTagResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "tag_id is required"},
		})
		return
	}

	req := &customer.UpdateConvTagReq{
		TagId:    body.TagID,
		TagName:  strings.TrimSpace(body.TagName),
		TagColor: strings.TrimSpace(body.TagColor),
		SortNo:   body.SortNo,
	}
	resp, err := h.client.UpdateConvTag(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// DeleteConvTag 删除会话标签
func (h *CustomerHandler) DeleteConvTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		TagID int64 `json:"tag_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.DeleteConvTagResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid json body"},
		})
		return
	}
	if body.TagID <= 0 {
		respondJSON(w, http.StatusBadRequest, &customer.DeleteConvTagResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "tag_id is required"},
		})
		return
	}

	resp, err := h.client.DeleteConvTag(r.Context(), &customer.DeleteConvTagReq{TagId: body.TagID})
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ============ 会话统计看板 ============

// GetConversationStats 获取会话统计数据
func (h *CustomerHandler) GetConversationStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	startDate := strings.TrimSpace(r.URL.Query().Get("start_date"))
	endDate := strings.TrimSpace(r.URL.Query().Get("end_date"))
	statType := strings.TrimSpace(r.URL.Query().Get("stat_type"))

	req := &customer.GetConversationStatsReq{
		StartDate: startDate,
		EndDate:   endDate,
		StatType:  statType,
	}
	resp, err := h.client.GetConversationStats(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ============ 用户认证相关接口 ============

// Login 用户登录
// 调用RPC验证用户名密码，成功后生成JWT Token返回
func (h *CustomerHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "invalid json body",
		})
		return
	}

	body.UserName = strings.TrimSpace(body.UserName)
	body.Password = strings.TrimSpace(body.Password)
	if body.UserName == "" || body.Password == "" {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"code": 400,
			"msg":  "用户名和密码不能为空",
		})
		return
	}

	// 调用RPC验证登录
	req := &customer.LoginReq{
		UserName: body.UserName,
		Password: body.Password,
	}
	resp, err := h.client.Login(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}

	// 登录验证失败，直接返回RPC响应
	if resp.BaseResp == nil || resp.BaseResp.Code != 0 {
		respondJSON(w, http.StatusOK, resp)
		return
	}

	// 登录成功，生成JWT Token
	token, err := generateJWTToken(resp.UserInfo.Id, resp.UserInfo.UserName, resp.UserInfo.RoleCode)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "生成Token失败: " + err.Error(),
		})
		return
	}

	// 返回登录成功响应（包含Token）
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "登录成功",
		"data": map[string]interface{}{
			"token":     token,
			"user_info": resp.UserInfo,
		},
	})
}

// GetCurrentUser 获取当前登录用户信息
func (h *CustomerHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 从上下文获取用户ID（由认证中间件写入）
	userID := getUserIDFromContext(r.Context())
	if userID <= 0 {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"code": 401,
			"msg":  "请先登录",
		})
		return
	}

	// 调用RPC获取用户信息
	req := &customer.GetCurrentUserReq{
		UserId: userID,
	}
	resp, err := h.client.GetCurrentUser(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// ============ JWT Token 工具函数 ============

// generateJWTToken 生成JWT Token的包装函数
// 参数:
//   - userID: 用户ID
//   - userName: 用户名
//   - roleCode: 角色编码
//
// 返回:
//   - string: 生成的Token字符串
//   - error: 错误信息
func generateJWTToken(userID int64, userName, roleCode string) (string, error) {
	return generateToken(userID, userName, roleCode)
}

// generateToken 生成JWT Token的实际实现
// 使用HS256算法签名，包含用户ID、用户名、角色编码、过期时间等声明
// 参数:
//   - userID: 用户ID
//   - userName: 用户名
//   - roleCode: 角色编码
//
// 返回:
//   - string: 生成的Token字符串
//   - error: 错误信息
func generateToken(userID int64, userName, roleCode string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_name": userName,
		"role_code": roleCode,
		"exp":       now.Add(time.Duration(config.GetJWTExpireHours()) * time.Hour).Unix(),
		"iat":       now.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetJWTSecret()))
}

// ============ 上下文工具函数 ============

// getUserIDFromContext 从上下文中获取用户ID
// 由认证中间件将用户ID存入context，此函数用于提取
// 参数:
//   - ctx: 请求上下文
//
// 返回:
//   - int64: 用户ID，未找到返回0
func getUserIDFromContext(ctx context.Context) int64 {
	if userID, ok := ctx.Value("user_id").(int64); ok {
		return userID
	}
	return 0
}

// resolveCustomerCsID 解析客服ID
// 根据当前登录用户角色解析客服ID：
// - 管理员：返回传入的csID
// - 客服：返回当前登录用户名或根据用户ID生成
// 参数:
//   - ctx: 请求上下文（包含用户信息）
//   - csID: 传入的客服ID
//
// 返回:
//   - string: 解析后的客服ID
func resolveCustomerCsID(ctx context.Context, csID string) string {
	roleCode := middleware.GetRoleCodeFromContext(ctx)
	if roleCode != middleware.RoleCustomerService {
		return strings.TrimSpace(csID)
	}
	userName := strings.TrimSpace(middleware.GetUserNameFromContext(ctx))
	if userName != "" {
		upper := strings.ToUpper(userName)
		if strings.HasPrefix(upper, "CS") || strings.HasPrefix(upper, "KF") {
			return userName
		}
	}
	userID := middleware.GetUserIDFromContext(ctx)
	if userID > 0 {
		return fmt.Sprintf("CS%d", userID)
	}
	return strings.TrimSpace(csID)
}

// OperatorInfo 操作人信息结构
type OperatorInfo struct {
	ID   string // 操作人ID（客服ID格式）
	Name string // 操作人姓名
	Role string // 操作人角色（admin/manager/customer_service）
}

// getOperatorInfoFromContext 从token上下文获取操作人信息
// 自动从JWT token中提取当前登录用户的ID、姓名和角色
// 用于审批、创建等操作的身份记录
func getOperatorInfoFromContext(ctx context.Context) OperatorInfo {
	info := OperatorInfo{}

	// 获取姓名
	info.Name = strings.TrimSpace(middleware.GetUserNameFromContext(ctx))

	// 获取角色和用户ID
	roleCode := middleware.GetRoleCodeFromContext(ctx)
	userID := middleware.GetUserIDFromContext(ctx)

	switch roleCode {
	case middleware.RoleAdmin:
		info.Role = "admin"
		if info.Name != "" {
			info.ID = info.Name
		} else {
			info.ID = fmt.Sprintf("ADMIN%d", userID)
		}
	case middleware.RoleCustomerService:
		info.Role = "customer_service"
		// 优先使用用户名（如果符合客服ID格式）
		if info.Name != "" {
			upper := strings.ToUpper(info.Name)
			if strings.HasPrefix(upper, "CS") || strings.HasPrefix(upper, "KF") {
				info.ID = info.Name
			}
		}
		// 如果用户名不符合格式，使用用户ID生成
		if info.ID == "" && userID > 0 {
			info.ID = fmt.Sprintf("CS%d", userID)
		}
	default:
		info.Role = roleCode
		if userID > 0 {
			info.ID = fmt.Sprintf("USER%d", userID)
		}
	}

	return info
}

// parsePaginationParams 从请求中解析分页参数
// 返回 page, pageSize (int32)
func parsePaginationParams(r *http.Request, defaultPageSize int) (int32, int32) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	return int32(page), int32(pageSize)
}

// Register 用户注册
// 仅允许注册客服账号，调用RPC完成注册
func (h *CustomerHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
		RealName string `json:"real_name"`
		Phone    string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "invalid json body",
		})
		return
	}

	body.UserName = strings.TrimSpace(body.UserName)
	body.Password = strings.TrimSpace(body.Password)
	body.RealName = strings.TrimSpace(body.RealName)
	body.Phone = strings.TrimSpace(body.Phone)

	// 调用RPC注册
	req := &customer.RegisterReq{
		UserName: body.UserName,
		Password: body.Password,
		RealName: body.RealName,
		Phone:    body.Phone,
	}
	resp, err := h.client.Register(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// ============ 快捷回复管理 ============

// CreateQuickReply 创建快捷回复
// 管理员或客服创建预设的快捷回复话术
// 请求方式: POST
// 请求体参数:
//   - reply_type: 回复类型（0-通用, 1-售前, 2-售后, 3-投诉）
//   - reply_content: 回复内容（必填）
//   - create_by: 创建人
//   - is_public: 是否公开（0-私有, 1-公开）
//
// 响应: 创建结果及新回复ID
func (h *CustomerHandler) CreateQuickReply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		ReplyType    int8   `json:"reply_type"`
		ReplyContent string `json:"reply_content"`
		CreateBy     string `json:"create_by"`
		IsPublic     int8   `json:"is_public"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "invalid json body",
		})
		return
	}

	// 如果未指定创建人，使用当前登录用户
	if strings.TrimSpace(body.CreateBy) == "" {
		userName := middleware.GetUserNameFromContext(r.Context())
		if userName != "" {
			body.CreateBy = userName
		} else {
			body.CreateBy = "ADMIN"
		}
	}

	req := &customer.CreateQuickReplyReq{
		ReplyType:    body.ReplyType,
		ReplyContent: strings.TrimSpace(body.ReplyContent),
		CreateBy:     strings.TrimSpace(body.CreateBy),
		IsPublic:     body.IsPublic,
	}

	resp, err := h.client.CreateQuickReply(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// UpdateQuickReply 更新快捷回复
// 修改已有的快捷回复内容或属性
// 请求方式: POST
// 请求体参数:
//   - reply_id: 回复ID（必填）
//   - reply_type: 回复类型
//   - reply_content: 回复内容（必填）
//   - is_public: 是否公开
//
// 响应: 更新结果
func (h *CustomerHandler) UpdateQuickReply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		ReplyID      int64  `json:"reply_id"`
		ReplyType    int8   `json:"reply_type"`
		ReplyContent string `json:"reply_content"`
		IsPublic     int8   `json:"is_public"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "invalid json body",
		})
		return
	}

	req := &customer.UpdateQuickReplyReq{
		ReplyId:      body.ReplyID,
		ReplyType:    body.ReplyType,
		ReplyContent: strings.TrimSpace(body.ReplyContent),
		IsPublic:     body.IsPublic,
	}

	resp, err := h.client.UpdateQuickReply(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// DeleteQuickReply 删除快捷回复
// 删除指定的快捷回复记录
// 请求方式: POST
// 请求体参数:
//   - reply_id: 回复ID（必填）
//
// 响应: 删除结果
func (h *CustomerHandler) DeleteQuickReply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		ReplyID int64 `json:"reply_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "invalid json body",
		})
		return
	}

	req := &customer.DeleteQuickReplyReq{
		ReplyId: body.ReplyID,
	}

	resp, err := h.client.DeleteQuickReply(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ============ 会话监控与导出 ============

// GetConversationMonitor 获取会话监控数据
// 实时查看会话状态、客服在线状态、等待中会话数等
// 请求方式: GET
// 请求参数:
//   - dept_id: 部门ID（可选，筛选指定部门）
//   - status_filter: 状态筛选 -1-全部 0-等待 1-进行中
//
// 响应: 客服状态列表、会话列表、统计数据
func (h *CustomerHandler) GetConversationMonitor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	deptID := strings.TrimSpace(r.URL.Query().Get("dept_id"))
	statusFilterStr := r.URL.Query().Get("status_filter")
	statusFilter := int8(-1) // 默认全部
	if statusFilterStr != "" {
		if v, err := strconv.Atoi(statusFilterStr); err == nil {
			statusFilter = int8(v)
		}
	}

	req := &customer.GetConversationMonitorReq{
		DeptId:       deptID,
		StatusFilter: statusFilter,
	}

	resp, err := h.client.GetConversationMonitor(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ExportConversations 导出会话记录
// 支持按条件筛选导出会话记录为Excel/CSV格式
// 请求方式: GET
// 请求参数:
//   - cs_id: 客服ID筛选（可选）
//   - user_id: 用户ID筛选（可选）
//   - start_date: 开始日期 YYYY-MM-DD
//   - end_date: 结束日期 YYYY-MM-DD
//   - status: 状态筛选 -1-全部
//   - keyword: 关键词搜索
//   - export_format: 导出格式 excel/csv（默认excel）
//
// 响应: 文件流下载
func (h *CustomerHandler) ExportConversations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	csID := strings.TrimSpace(r.URL.Query().Get("cs_id"))
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))
	startDate := strings.TrimSpace(r.URL.Query().Get("start_date"))
	endDate := strings.TrimSpace(r.URL.Query().Get("end_date"))
	keyword := strings.TrimSpace(r.URL.Query().Get("keyword"))
	exportFormat := strings.TrimSpace(r.URL.Query().Get("export_format"))
	statusStr := r.URL.Query().Get("status")

	status := int8(-1)
	if statusStr != "" {
		if v, err := strconv.Atoi(statusStr); err == nil {
			status = int8(v)
		}
	}

	if exportFormat == "" {
		exportFormat = "excel"
	}

	req := &customer.ExportConversationsReq{
		CsId:         csID,
		UserId:       userID,
		StartDate:    startDate,
		EndDate:      endDate,
		Status:       status,
		Keyword:      keyword,
		ExportFormat: exportFormat,
	}

	resp, err := h.client.ExportConversations(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}

	if resp.BaseResp != nil && resp.BaseResp.Code != 0 {
		respondJSON(w, http.StatusOK, resp)
		return
	}

	// 返回文件流
	contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	if exportFormat == "csv" {
		contentType = "text/csv"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", resp.FileName))
	w.Write(resp.FileData)
}

// ============ 消息分类管理 ============

// MsgAutoClassify 消息自动分类
// 基于关键词匹配对会话消息进行自动分类
// 请求方式: POST
// 请求体参数:
//   - conv_id: 会话ID
//   - msg_contents: 消息内容列表
//
// 响应: 分类ID、分类名称、置信度、匹配的关键词
func (h *CustomerHandler) MsgAutoClassify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		ConvID      string   `json:"conv_id"`
		MsgContents []string `json:"msg_contents"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "invalid json body",
		})
		return
	}

	req := &customer.MsgAutoClassifyReq{
		ConvId:      strings.TrimSpace(body.ConvID),
		MsgContents: body.MsgContents,
	}

	resp, err := h.client.MsgAutoClassify(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// AdjustMsgClassify 人工调整消息分类
// 客服手动修正自动分类结果
// 请求方式: POST
// 请求体参数:
//   - conv_id: 会话ID
//   - original_category_id: 原分类ID
//   - new_category_id: 新分类ID
//   - operator_id: 操作人id
//   - adjust_reason: 调整原因
//
// 响应: 调整记录ID
func (h *CustomerHandler) AdjustMsgClassify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		ConvID             string `json:"conv_id"`
		OriginalCategoryID int64  `json:"original_category_id"`
		NewCategoryID      int64  `json:"new_category_id"`
		OperatorID         string `json:"operator_id"`
		AdjustReason       string `json:"adjust_reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "invalid json body",
		})
		return
	}

	// 如果未指定操作人，使用当前登录用户
	if strings.TrimSpace(body.OperatorID) == "" {
		userName := middleware.GetUserNameFromContext(r.Context())
		if userName != "" {
			body.OperatorID = userName
		}
	}

	req := &customer.AdjustMsgClassifyReq{
		ConvId:             strings.TrimSpace(body.ConvID),
		OriginalCategoryId: body.OriginalCategoryID,
		NewCategoryId_:     body.NewCategoryID,
		OperatorId:         strings.TrimSpace(body.OperatorID),
		AdjustReason:       strings.TrimSpace(body.AdjustReason),
	}

	resp, err := h.client.AdjustMsgClassify(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// GetClassifyStats 获取分类统计数据
// 查询消息分类的统计信息
// 请求方式: GET
// 请求参数:
//   - start_date: 开始日期 YYYY-MM-DD
//   - end_date: 结束日期 YYYY-MM-DD
//   - stat_type: 统计类型 day/week/month
//
// 响应: 每日统计、分类汇总、准确率等
func (h *CustomerHandler) GetClassifyStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	startDate := strings.TrimSpace(r.URL.Query().Get("start_date"))
	endDate := strings.TrimSpace(r.URL.Query().Get("end_date"))
	statType := strings.TrimSpace(r.URL.Query().Get("stat_type"))

	req := &customer.GetClassifyStatsReq{
		StartDate: startDate,
		EndDate:   endDate,
		StatType:  statType,
	}

	resp, err := h.client.GetClassifyStats(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ============ 消息分类维度CRUD ============

// CreateMsgCategory 创建消息分类维度
// 创建新的消息分类类型（如咨询类、投诉类、建议类等）
// 请求方式: POST
// 请求体参数:
//   - category_name: 分类名称（必填）
//   - keywords: 关键词列表(JSON)
//   - sort_no: 排序号
//   - create_by: 创建人
//
// 响应: 创建结果及新分类ID
func (h *CustomerHandler) CreateMsgCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		CategoryName string `json:"category_name"`
		Keywords     string `json:"keywords"`
		SortNo       int32  `json:"sort_no"`
		CreateBy     string `json:"create_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "invalid json body",
		})
		return
	}

	// 如果未指定创建人，使用当前登录用户
	if strings.TrimSpace(body.CreateBy) == "" {
		userName := middleware.GetUserNameFromContext(r.Context())
		if userName != "" {
			body.CreateBy = userName
		} else {
			body.CreateBy = "ADMIN"
		}
	}

	req := &customer.CreateMsgCategoryReq{
		CategoryName: strings.TrimSpace(body.CategoryName),
		Keywords:     body.Keywords,
		SortNo:       body.SortNo,
		CreateBy:     strings.TrimSpace(body.CreateBy),
	}

	resp, err := h.client.CreateMsgCategory(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ListMsgCategory 查询消息分类维度列表
// 获取所有消息分类类型
// 请求方式: GET
// 响应: 分类列表
func (h *CustomerHandler) ListMsgCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	req := &customer.ListMsgCategoryReq{}

	resp, err := h.client.ListMsgCategory(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// UpdateMsgCategory 更新消息分类维度
// 修改消息分类的名称、关键词等属性
// 请求方式: POST
// 请求体参数:
//   - category_id: 分类ID（必填）
//   - category_name: 分类名称
//   - keywords: 关键词列表(JSON)
//   - sort_no: 排序号
//
// 响应: 更新结果
func (h *CustomerHandler) UpdateMsgCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		CategoryID   int64  `json:"category_id"`
		CategoryName string `json:"category_name"`
		Keywords     string `json:"keywords"`
		SortNo       int32  `json:"sort_no"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "invalid json body",
		})
		return
	}

	req := &customer.UpdateMsgCategoryReq{
		CategoryId:   body.CategoryID,
		CategoryName: strings.TrimSpace(body.CategoryName),
		Keywords:     body.Keywords,
		SortNo:       body.SortNo,
	}

	resp, err := h.client.UpdateMsgCategory(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// DeleteMsgCategory 删除消息分类维度
// 删除指定的消息分类类型
// 请求方式: POST
// 请求体参数:
//   - category_id: 分类ID（必填）
//
// 响应: 删除结果
func (h *CustomerHandler) DeleteMsgCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		CategoryID int64 `json:"category_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "invalid json body",
		})
		return
	}

	req := &customer.DeleteMsgCategoryReq{
		CategoryId: body.CategoryID,
	}

	resp, err := h.client.DeleteMsgCategory(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ============ 消息加密与脱敏 ============

// EncryptMessage 加密消息内容
// 使用AES-256-GCM算法加密敏感消息
// 请求方式: POST
// 请求体参数:
//   - msg_content: 待加密的消息内容（必填）
//
// 响应: 加密后的内容
func (h *CustomerHandler) EncryptMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		MsgContent string `json:"msg_content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "invalid json body",
		})
		return
	}

	if body.MsgContent == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "msg_content is required",
		})
		return
	}

	req := &customer.EncryptMessageReq{
		MsgContent: body.MsgContent,
	}

	resp, err := h.client.EncryptMessage(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// DecryptMessage 解密消息内容
// 解密已加密的消息内容
// 请求方式: POST
// 请求体参数:
//   - encrypted_content: 加密后的消息内容（必填）
//
// 响应: 解密后的原始内容
func (h *CustomerHandler) DecryptMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		EncryptedContent string `json:"encrypted_content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "invalid json body",
		})
		return
	}

	if body.EncryptedContent == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "encrypted_content is required",
		})
		return
	}

	req := &customer.DecryptMessageReq{
		EncryptedContent: body.EncryptedContent,
	}

	resp, err := h.client.DecryptMessage(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// DesensitizeMessage 消息脱敏处理
// 对消息中的敏感信息（手机号、身份证、银行卡、邮箱）进行脱敏
// 请求方式: POST
// 请求体参数:
//   - msg_content: 待脱敏的消息内容（必填）
//
// 响应: 脱敏后的内容和检测到的敏感信息类型
func (h *CustomerHandler) DesensitizeMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		MsgContent string `json:"msg_content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "invalid json body",
		})
		return
	}

	if body.MsgContent == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "msg_content is required",
		})
		return
	}

	req := &customer.DesensitizeMessageReq{
		MsgContent: body.MsgContent,
	}

	resp, err := h.client.DesensitizeMessage(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ============ 数据归档管理 ============

// ArchiveConversations 归档历史会话
// 将指定日期之前的会话数据归档到归档表
// 请求方式: POST
// 请求体参数:
//   - end_date: 截止日期，归档此日期之前的数据（格式：2006-01-02）
//   - retention_days: 归档数据保留天数（默认365天）
//   - operator_id: 操作人ID
//
// 响应: 归档任务ID和预计归档数量
func (h *CustomerHandler) ArchiveConversations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		EndDate       string `json:"end_date"`
		RetentionDays int32  `json:"retention_days"`
		OperatorID    string `json:"operator_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "invalid json body",
		})
		return
	}

	if body.EndDate == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "end_date is required",
		})
		return
	}

	// 获取操作人信息
	operatorInfo := getOperatorInfoFromContext(r.Context())

	req := &customer.ArchiveConversationsReq{
		EndDate:       body.EndDate,
		RetentionDays: body.RetentionDays,
		OperatorId:    operatorInfo.ID,
	}

	resp, err := h.client.ArchiveConversations(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// GetArchiveTask 获取归档任务状态
// 查询归档任务的执行进度和状态
// 请求方式: GET
// 查询参数:
//   - task_id: 归档任务ID（必填）
//
// 响应: 任务状态、进度、已归档数量等信息
func (h *CustomerHandler) GetArchiveTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	taskIDStr := r.URL.Query().Get("task_id")
	if taskIDStr == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "task_id is required",
		})
		return
	}

	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "invalid task_id",
		})
		return
	}

	req := &customer.GetArchiveTaskReq{
		TaskId: taskID,
	}

	resp, err := h.client.GetArchiveTask(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// QueryArchivedConversation 查询归档会话
// 从归档表中查询历史会话数据
// 请求方式: GET
// 查询参数:
//   - user_id: 用户ID（可选，按用户查询）
//   - cs_id: 客服ID（可选，按客服查询）
//   - start_date: 开始日期（格式：2006-01-02）
//   - end_date: 结束日期（格式：2006-01-02）
//   - page: 页码（默认1）
//   - page_size: 每页条数（默认20）
//
// 响应: 归档会话列表和总数
func (h *CustomerHandler) QueryArchivedConversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(query.Get("page_size"))
	if pageSize <= 0 {
		pageSize = 20
	}

	req := &customer.QueryArchivedConversationReq{
		UserId:    query.Get("user_id"),
		CsId:      query.Get("cs_id"),
		StartDate: query.Get("start_date"),
		EndDate:   query.Get("end_date"),
		Page:      int32(page),
		PageSize:  int32(pageSize),
	}

	resp, err := h.client.QueryArchivedConversation(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ============ 心跳与在线状态管理 ============

// Heartbeat 客服心跳上报
// 客服端定期调用此接口保持在线状态
// 请求方式: POST
// 请求体参数:
//   - cs_id: 客服ID（必填）
//
// 响应: 在线状态确认
func (h *CustomerHandler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		CsID string `json:"cs_id"` // 客服ID
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.HeartbeatResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid json body"},
		})
		return
	}
	body.CsID = strings.TrimSpace(body.CsID)
	if body.CsID == "" {
		respondJSON(w, http.StatusBadRequest, &customer.HeartbeatResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "cs_id is required"},
		})
		return
	}

	req := &customer.HeartbeatReq{
		CsId: body.CsID,
	}
	resp, err := h.client.Heartbeat(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ListOnlineCustomers 获取当前在线客服列表
// 请求方式: GET
// 查询参数:
//   - dept_id: 部门ID（可选，按部门筛选）
//
// 响应: 在线客服列表
func (h *CustomerHandler) ListOnlineCustomers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	deptID := strings.TrimSpace(r.URL.Query().Get("dept_id"))

	req := &customer.ListOnlineCustomersReq{
		DeptId: deptID,
	}
	resp, err := h.client.ListOnlineCustomers(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ============ 调班辅助接口 ============

// GetSwapCandidates 获取调班候选人
// 返回指定日期有排班的其他客服及班次信息
// 请求方式: GET
// 查询参数:
//   - cs_id: 发起人客服ID（必填）
//   - target_date: 调班日期（必填）
//
// 响应: 可调班候选人列表
func (h *CustomerHandler) GetSwapCandidates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	csID := strings.TrimSpace(r.URL.Query().Get("cs_id"))
	targetDate := strings.TrimSpace(r.URL.Query().Get("target_date"))

	if csID == "" || targetDate == "" {
		respondJSON(w, http.StatusBadRequest, &customer.GetSwapCandidatesResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "cs_id and target_date are required"},
		})
		return
	}

	req := &customer.GetSwapCandidatesReq{
		CsId:       csID,
		TargetDate: targetDate,
	}
	resp, err := h.client.GetSwapCandidates(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// CheckSwapConflict 检测调班冲突
// 检测发起人与目标客服之间是否存在调班冲突
// 请求方式: POST
// 请求体参数:
//   - initiator_cs_id: 发起人客服ID（必填）
//   - target_cs_id: 目标客服ID（必填）
//   - target_date: 调班日期（必填）
//
// 响应: 冲突检测结果
func (h *CustomerHandler) CheckSwapConflict(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		InitiatorCsID string `json:"initiator_cs_id"` // 发起人客服ID
		TargetCsID    string `json:"target_cs_id"`    // 目标客服ID
		TargetDate    string `json:"target_date"`     // 调班日期
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, &customer.CheckSwapConflictResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid json body"},
		})
		return
	}
	body.InitiatorCsID = strings.TrimSpace(body.InitiatorCsID)
	body.TargetCsID = strings.TrimSpace(body.TargetCsID)
	body.TargetDate = strings.TrimSpace(body.TargetDate)

	if body.InitiatorCsID == "" || body.TargetCsID == "" || body.TargetDate == "" {
		respondJSON(w, http.StatusBadRequest, &customer.CheckSwapConflictResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "initiator_cs_id, target_cs_id and target_date are required"},
		})
		return
	}

	req := &customer.CheckSwapConflictReq{
		InitiatorCsId: body.InitiatorCsID,
		TargetCsId:    body.TargetCsID,
		TargetDate:    body.TargetDate,
	}
	resp, err := h.client.CheckSwapConflict(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// ============ 退出登录 ============

// Logout 客服退出登录
// 将客服置为离线状态
// 请求方式: POST
// 请求体参数:
//   - cs_id: 客服ID（必填）
//
// 响应: 退出结果
func (h *CustomerHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body struct {
		CsID string `json:"cs_id"` // 客服ID
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "invalid json body",
		})
		return
	}
	body.CsID = strings.TrimSpace(body.CsID)
	if body.CsID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "cs_id is required",
		})
		return
	}

	req := &customer.LogoutReq{
		CsId: body.CsID,
	}
	resp, err := h.client.Logout(r.Context(), req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "Internal server error: " + err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, resp)
}
