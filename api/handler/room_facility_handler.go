package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SetRoomFacilities 设置房源的设施
// @Summary 设置房源的设施
// @Description 批量设置房源的设施（会先删除旧的，再创建新的）
// @Tags 房源管理
// @Accept json
// @Produce json
// @Param id path int true "房源ID"
// @Param request body map[string]interface{} true "设施ID列表" example({"facility_ids":[1,2,3]})
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos/{id}/facilities [put]
func (h *RoomHandler) SetRoomFacilities(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	var req struct {
		FacilityIDs []uint64 `json:"facility_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	if err := h.RoomFacilityService.SetRoomFacilities(id, req.FacilityIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "设置成功"})
}

// GetRoomFacilities 获取房源的设施列表
// @Summary 获取房源的设施列表
// @Tags 房源管理
// @Produce json
// @Param id path int true "房源ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos/{id}/facilities [get]
func (h *RoomHandler) GetRoomFacilities(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	facilities, err := h.RoomFacilityService.GetRoomFacilities(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": facilities})
}

// AddRoomFacility 为房源添加单个设施
// @Summary 为房源添加单个设施
// @Tags 房源管理
// @Accept json
// @Produce json
// @Param id path int true "房源ID"
// @Param request body map[string]interface{} true "设施ID" example({"facility_id":1})
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos/{id}/facilities [post]
func (h *RoomHandler) AddRoomFacility(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	var req struct {
		FacilityID uint64 `json:"facility_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	if err := h.RoomFacilityService.AddRoomFacility(id, req.FacilityID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "添加成功"})
}

// RemoveRoomFacility 移除房源的单个设施
// @Summary 移除房源的单个设施
// @Tags 房源管理
// @Produce json
// @Param id path int true "房源ID"
// @Param facility_id path int true "设施ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos/{id}/facilities/{facility_id} [delete]
func (h *RoomHandler) RemoveRoomFacility(c *gin.Context) {
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	facilityID, err := strconv.ParseUint(c.Param("facility_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := h.RoomFacilityService.RemoveRoomFacility(roomID, facilityID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "移除成功"})
}
