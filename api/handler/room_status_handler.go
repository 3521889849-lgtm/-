package handler

import (
	"context"
	"example_shop/api/client"
	"example_shop/common/service"
	"example_shop/kitex_gen/hotel"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// UpdateRoomStatus 更新房源状态
// @Summary 更新房源状态
// @Description 支持启用、停用、维修等状态切换
// @Tags 房源管理
// @Accept json
// @Produce json
// @Param id path int true "房源ID"
// @Param status body map[string]string true "状态信息" example({"status":"ACTIVE"})
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos/{id}/status [put]
func (h *RoomHandler) UpdateRoomStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	rpcReq := &hotel.UpdateRoomStatusReq{
		RoomId: id,
		Status: req.Status,
	}

	resp, err := client.HotelClient.UpdateRoomStatus(context.Background(), rpcReq)
	code, msg := handleRPCError(resp, err)
	if code != 200 {
		c.JSON(code, gin.H{"code": code, "msg": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": msg})
}

// BatchUpdateRoomStatus 批量更新房源状态
// @Summary 批量更新房源状态
// @Tags 房源管理
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "批量更新请求" example({"room_ids":[1,2,3],"status":"ACTIVE"})
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos/batch-status [put]
func (h *RoomHandler) BatchUpdateRoomStatus(c *gin.Context) {
	var req struct {
		RoomIDs []int64 `json:"room_ids" binding:"required"`
		Status  string  `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	rpcReq := &hotel.BatchUpdateRoomStatusReq{
		RoomIds: req.RoomIDs,
		Status:  req.Status,
	}

	resp, err := client.HotelClient.BatchUpdateRoomStatus(context.Background(), rpcReq)
	code, msg := handleRPCError(resp, err)
	if code != 200 {
		c.JSON(code, gin.H{"code": code, "msg": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": msg})
}

// GetCalendarRoomStatus 获取日历化房态
// @Summary 获取日历化房态
// @Description 按日期维度展示每日各房间的状态
// @Tags 房态管理
// @Accept json
// @Produce json
// @Param branch_id query int false "分店ID"
// @Param start_date query string true "开始日期 (YYYY-MM-DD)"
// @Param end_date query string true "结束日期 (YYYY-MM-DD)"
// @Param room_no query string false "房间号"
// @Param status query string false "房态筛选"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/calendar-room-status [get]
func (h *RoomHandler) GetCalendarRoomStatus(c *gin.Context) {
	// 解析日期参数
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "开始日期和结束日期不能为空"})
		return
	}

	// 解析日期
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "开始日期格式错误，应为 YYYY-MM-DD"})
		return
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "结束日期格式错误，应为 YYYY-MM-DD"})
		return
	}

	// 构建请求对象
	req := service.CalendarRoomStatusReq{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// 解析可选参数
	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		if branchID, err := strconv.ParseUint(branchIDStr, 10, 64); err == nil {
			req.BranchID = &branchID
		}
	}

	if roomNo := c.Query("room_no"); roomNo != "" {
		req.RoomNo = &roomNo
	}

	if status := c.Query("status"); status != "" {
		req.Status = &status
	}

	// 调用 Service 层
	items, err := h.RoomStatusService.GetCalendarRoomStatus(req)
	if err != nil {
		// 记录详细错误信息以便调试
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取日历化房态失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": items,
	})
}

// UpdateCalendarRoomStatus 更新日历化房态
// @Summary 更新日历化房态
// @Description 更新指定房间指定日期的房态
// @Tags 房态管理
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "更新请求" example({"room_id":1,"date":"2024-01-15","status":"入住房"})
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/calendar-room-status [put]
func (h *RoomHandler) UpdateCalendarRoomStatus(c *gin.Context) {
	var req struct {
		RoomID uint64 `json:"room_id" binding:"required"`
		Date   string `json:"date" binding:"required"`   // YYYY-MM-DD
		Status string `json:"status" binding:"required"` // 空净房/入住房/维修房/锁定房/空账房/预定房
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	rpcReq := &hotel.UpdateCalendarRoomStatusReq{
		RoomId: int64(req.RoomID),
		Date:   req.Date,
		Status: req.Status,
	}

	resp, err := client.HotelClient.UpdateCalendarRoomStatus(context.Background(), rpcReq)
	code, msg := handleRPCError(resp, err)
	if code != 200 {
		c.JSON(code, gin.H{"code": code, "msg": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": msg})
}

// BatchUpdateCalendarRoomStatus 批量更新日历化房态
// @Summary 批量更新日历化房态
// @Description 批量更新多个房间多个日期的房态
// @Tags 房态管理
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "批量更新请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/calendar-room-status/batch [put]
func (h *RoomHandler) BatchUpdateCalendarRoomStatus(c *gin.Context) {
	var req struct {
		Updates []struct {
			RoomID uint64 `json:"room_id" binding:"required"`
			Date   string `json:"date" binding:"required"`
			Status string `json:"status" binding:"required"`
		} `json:"updates" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	updates := make([]*hotel.UpdateCalendarRoomStatusReq, len(req.Updates))
	for i, update := range req.Updates {
		updates[i] = &hotel.UpdateCalendarRoomStatusReq{
			RoomId: int64(update.RoomID),
			Date:   update.Date,
			Status: update.Status,
		}
	}

	rpcReq := &hotel.BatchUpdateCalendarRoomStatusReq{
		Updates: updates,
	}

	resp, err := client.HotelClient.BatchUpdateCalendarRoomStatus(context.Background(), rpcReq)
	code, msg := handleRPCError(resp, err)
	if code != 200 {
		c.JSON(code, gin.H{"code": code, "msg": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": msg})
}

// GetRealTimeStatistics 获取实时数据统计
// @Summary 获取实时数据统计
// @Description 显示当日剩余房间数、已入住人数、预退房人数等核心数据
// @Tags 房态管理
// @Accept json
// @Produce json
// @Param branch_id query int false "分店ID"
// @Param date query string false "日期 (YYYY-MM-DD)，默认为今日"
// @Param room_no query string false "房间号"
// @Param room_type_id query int false "房型ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/real-time-statistics [get]
func (h *RoomHandler) GetRealTimeStatistics(c *gin.Context) {
	// 构建请求对象
	req := service.RealTimeStatisticsReq{}

	// 解析可选参数
	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		if branchID, err := strconv.ParseUint(branchIDStr, 10, 64); err == nil {
			req.BranchID = &branchID
		}
	}

	if date := c.Query("date"); date != "" {
		req.Date = &date
	}

	if roomNo := c.Query("room_no"); roomNo != "" {
		req.RoomNo = &roomNo
	}

	if roomTypeIDStr := c.Query("room_type_id"); roomTypeIDStr != "" {
		if roomTypeID, err := strconv.ParseUint(roomTypeIDStr, 10, 64); err == nil {
			req.RoomTypeID = &roomTypeID
		}
	}

	// 调用 Service 层
	resp, err := h.RoomStatusService.GetRealTimeStatistics(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取实时统计数据失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": resp,
	})
}

// 辅助函数：处理 RPC 错误
func handleRPCError(resp *hotel.BaseResp, err error) (int, string) {
	if err != nil {
		return 500, "RPC 调用失败: " + err.Error()
	}
	if resp.Code != 200 {
		return int(resp.Code), resp.Msg
	}
	return 200, resp.Msg
}
