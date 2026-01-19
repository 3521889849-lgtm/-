package handler

import (
	"example_shop/common/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateFacility 创建设施字典
// @Summary 创建设施字典
// @Tags 设施管理
// @Accept json
// @Produce json
// @Param facility body service.CreateFacilityReq true "设施信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/facilities [post]
func (h *RoomHandler) CreateFacility(c *gin.Context) {
	var req service.CreateFacilityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	facility, err := h.FacilityService.CreateFacility(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": facility})
}

// UpdateFacility 更新设施字典
// @Summary 更新设施字典
// @Tags 设施管理
// @Accept json
// @Produce json
// @Param id path int true "设施ID"
// @Param facility body service.UpdateFacilityReq true "设施信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/facilities/{id} [put]
func (h *RoomHandler) UpdateFacility(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	var req service.UpdateFacilityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	if err := h.FacilityService.UpdateFacility(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功"})
}

// GetFacility 获取设施详情
// @Summary 获取设施详情
// @Tags 设施管理
// @Produce json
// @Param id path int true "设施ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/facilities/{id} [get]
func (h *RoomHandler) GetFacility(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	facility, err := h.FacilityService.GetFacility(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": facility})
}

// ListFacilities 获取设施列表
// @Summary 获取设施列表
// @Tags 设施管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param status query string false "状态筛选"
// @Param keyword query string false "关键词搜索"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/facilities [get]
func (h *RoomHandler) ListFacilities(c *gin.Context) {
	var req service.ListFacilityReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	facilities, total, err := h.FacilityService.ListFacilities(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":      facilities,
			"total":     total,
			"page":      req.Page,
			"page_size": req.PageSize,
		},
	})
}

// DeleteFacility 删除设施字典
// @Summary 删除设施字典
// @Tags 设施管理
// @Produce json
// @Param id path int true "设施ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/facilities/{id} [delete]
func (h *RoomHandler) DeleteFacility(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := h.FacilityService.DeleteFacility(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功"})
}
