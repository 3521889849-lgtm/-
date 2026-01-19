package handler

import (
	"context"
	"example_shop/api/client"
	"example_shop/common/service"
	"example_shop/kitex_gen/hotel"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type BranchHandler struct {
	BranchService *service.BranchService
}

func NewBranchHandler() *BranchHandler {
	return &BranchHandler{
		BranchService: &service.BranchService{},
	}
}

// ListBranches 获取分店列表
func (h *BranchHandler) ListBranches(c *gin.Context) {
	var req service.ListBranchesReq

	if status := strings.TrimSpace(c.Query("status")); status != "" {
		statusUpper := strings.ToUpper(status)
		if statusUpper != "ALL" {
			req.Status = &statusUpper
		}
	}

	branches, err := h.BranchService.ListBranches(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": gin.H{
			"list":  branches,
			"total": len(branches),
		},
	})
}

// GetBranch 获取分店详情
func (h *BranchHandler) GetBranch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	branch, err := h.BranchService.GetBranch(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "分店不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": branch,
	})
}

// CreateBranch 创建分店
func (h *BranchHandler) CreateBranch(c *gin.Context) {
	var req service.CreateBranchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	branch, err := h.BranchService.CreateBranch(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": branch})
}

// UpdateBranch 更新分店
func (h *BranchHandler) UpdateBranch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	var req service.UpdateBranchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	if err := h.BranchService.UpdateBranch(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功"})
}

// DeleteBranch 删除分店
func (h *BranchHandler) DeleteBranch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := h.BranchService.DeleteBranch(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "该分店下有房源，无法删除"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功"})
}

// SyncRoomStatusToChannel 同步房态数据到渠道
// @Summary 同步房态数据到渠道
// @Description 支持房态数据同步至合作渠道（如途游）
// @Tags 渠道同步
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "同步请求" example({"branch_id":1,"channel_id":1,"start_date":"2024-01-15","end_date":"2024-01-21"})
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/sync-room-status [post]
func (h *RoomHandler) SyncRoomStatusToChannel(c *gin.Context) {
	var req struct {
		BranchID  uint64   `json:"branch_id" binding:"required"`
		ChannelID uint64   `json:"channel_id" binding:"required"`
		StartDate string   `json:"start_date" binding:"required"`
		EndDate   string   `json:"end_date" binding:"required"`
		RoomIDs   []uint64 `json:"room_ids,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	roomIDs := make([]int64, len(req.RoomIDs))
	for i, id := range req.RoomIDs {
		roomIDs[i] = int64(id)
	}

	rpcReq := &hotel.SyncRoomStatusToChannelReq{
		BranchId:  int64(req.BranchID),
		ChannelId: int64(req.ChannelID),
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}
	if len(roomIDs) > 0 {
		rpcReq.RoomIds = roomIDs
	}

	resp, err := client.HotelClient.SyncRoomStatusToChannel(context.Background(), rpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "RPC 调用失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "同步完成",
		"data": resp,
	})
}
