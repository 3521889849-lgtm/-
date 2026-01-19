package handler

import (
	"example_shop/api/client"
	"example_shop/kitex_gen/hotel"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MemberPointsHandler struct{}

func NewMemberPointsHandler() *MemberPointsHandler {
	return &MemberPointsHandler{}
}

// CreatePointsRecord 创建积分记录
func (h *MemberPointsHandler) CreatePointsRecord(c *gin.Context) {
	var req hotel.CreatePointsRecordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "请求参数错误: "+err.Error())
		return
	}

	resp, err := client.HotelClient.CreatePointsRecord(c.Request.Context(), &req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListPointsRecords 获取积分记录列表
func (h *MemberPointsHandler) ListPointsRecords(c *gin.Context) {
	req := &hotel.ListPointsRecordsReq{
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
	if memberIDStr := c.Query("member_id"); memberIDStr != "" {
		if memberID, err := strconv.ParseInt(memberIDStr, 10, 64); err == nil {
			req.MemberId = &memberID
		}
	}
	if orderIDStr := c.Query("order_id"); orderIDStr != "" {
		if orderID, err := strconv.ParseInt(orderIDStr, 10, 64); err == nil {
			req.OrderId = &orderID
		}
	}
	if changeType := c.Query("change_type"); changeType != "" {
		req.ChangeType = &changeType
	}
	if startTime := c.Query("start_time"); startTime != "" {
		req.StartTime = &startTime
	}
	if endTime := c.Query("end_time"); endTime != "" {
		req.EndTime = &endTime
	}

	resp, err := client.HotelClient.ListPointsRecords(c.Request.Context(), req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetMemberPointsBalance 获取会员积分余额
func (h *MemberPointsHandler) GetMemberPointsBalance(c *gin.Context) {
	memberIDStr := c.Param("id")
	memberID, err := strconv.ParseInt(memberIDStr, 10, 64)
	if err != nil {
		badRequest(c, "参数错误")
		return
	}

	balance, err := client.HotelClient.GetMemberPointsBalance(c.Request.Context(), memberID)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
