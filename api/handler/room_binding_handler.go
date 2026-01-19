package handler

import (
	"context"
	"example_shop/api/client"
	"example_shop/kitex_gen/hotel"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateRoomBinding 创建关联房绑定
// @Summary 创建关联房绑定
// @Tags 房源管理
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "关联房信息" example({"main_room_id":1,"related_room_id":2,"binding_desc":"关联描述"})
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos/bindings [post]
func (h *RoomHandler) CreateRoomBinding(c *gin.Context) {
	var req struct {
		MainRoomID    int64   `json:"main_room_id" binding:"required"`
		RelatedRoomID int64   `json:"related_room_id" binding:"required"`
		BindingDesc   *string `json:"binding_desc,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	rpcReq := &hotel.CreateRoomBindingReq{
		MainRoomId:    req.MainRoomID,
		RelatedRoomId: req.RelatedRoomID,
		BindingDesc:   req.BindingDesc,
	}

	resp, err := client.HotelClient.CreateRoomBinding(context.Background(), rpcReq)
	code, msg := handleRPCError(resp, err)
	if code != 200 {
		c.JSON(code, gin.H{"code": code, "msg": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": msg})
}

// BatchCreateRoomBindings 批量创建关联房绑定
// @Summary 批量创建关联房绑定
// @Tags 房源管理
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "批量关联信息" example({"main_room_id":1,"related_room_ids":[2,3,4],"binding_desc":"关联描述"})
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos/batch-bindings [post]
func (h *RoomHandler) BatchCreateRoomBindings(c *gin.Context) {
	var req struct {
		MainRoomID     int64    `json:"main_room_id" binding:"required"`
		RelatedRoomIDs []int64  `json:"related_room_ids" binding:"required"`
		BindingDesc    *string  `json:"binding_desc,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	rpcReq := &hotel.BatchCreateRoomBindingsReq{
		MainRoomId:     req.MainRoomID,
		RelatedRoomIds: req.RelatedRoomIDs,
		BindingDesc:    req.BindingDesc,
	}

	resp, err := client.HotelClient.BatchCreateRoomBindings(context.Background(), rpcReq)
	code, msg := handleRPCError(resp, err)
	if code != 200 {
		c.JSON(code, gin.H{"code": code, "msg": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": msg})
}

// GetRoomBindings 获取房源的关联房列表
// @Summary 获取房源的关联房列表
// @Tags 房源管理
// @Produce json
// @Param id path int true "房源ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos/{id}/bindings [get]
func (h *RoomHandler) GetRoomBindings(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	resp, err := client.HotelClient.GetRoomBindings(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "RPC 调用失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": resp.Bindings})
}

// DeleteRoomBinding 删除关联房绑定
// @Summary 删除关联房绑定
// @Tags 房源管理
// @Produce json
// @Param id path int true "绑定ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos/bindings/{id} [delete]
func (h *RoomHandler) DeleteRoomBinding(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	resp, err := client.HotelClient.DeleteRoomBinding(context.Background(), id)
	code, msg := handleRPCError(resp, err)
	if code != 200 {
		c.JSON(code, gin.H{"code": code, "msg": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": msg})
}
