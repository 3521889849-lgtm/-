package handler

import (
	"example_shop/api/client"
	"example_shop/kitex_gen/hotel"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OperationLogHandler struct{}

func NewOperationLogHandler() *OperationLogHandler {
	return &OperationLogHandler{}
}

// CreateOperationLog 创建操作日志
func (h *OperationLogHandler) CreateOperationLog(c *gin.Context) {
	var req hotel.CreateOperationLogReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "请求参数错误: "+err.Error())
		return
	}

	resp, err := client.HotelClient.CreateOperationLog(c.Request.Context(), &req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListOperationLogs 查询操作日志列表
func (h *OperationLogHandler) ListOperationLogs(c *gin.Context) {
	req := &hotel.ListOperationLogsReq{
		Page:     1,
		PageSize: 10,
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.ParseInt(pageStr, 10, 32); err == nil {
			req.Page = int32(page)
		}
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32); err == nil {
			req.PageSize = int32(pageSize)
		}
	}
	if operatorIDStr := c.Query("operator_id"); operatorIDStr != "" {
		if operatorID, err := strconv.ParseInt(operatorIDStr, 10, 64); err == nil {
			id := operatorID
			req.OperatorId = &id
		}
	}
	if module := c.Query("module"); module != "" {
		req.Module = &module
	}
	if operationType := c.Query("operation_type"); operationType != "" {
		req.OperationType = &operationType
	}
	if startTime := c.Query("start_time"); startTime != "" {
		req.StartTime = &startTime
	}
	if endTime := c.Query("end_time"); endTime != "" {
		req.EndTime = &endTime
	}
	if isSuccessStr := c.Query("is_success"); isSuccessStr != "" {
		if isSuccess, err := strconv.ParseBool(isSuccessStr); err == nil {
			req.IsSuccess = &isSuccess
		}
	}

	resp, err := client.HotelClient.ListOperationLogs(c.Request.Context(), req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}
