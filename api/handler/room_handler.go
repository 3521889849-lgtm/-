package handler

import (
	"example_shop/common/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	RoomTypeService           *service.RoomTypeService
	RoomInfoService           *service.RoomInfoService
	RoomStatusService         *service.RoomStatusService
	RoomBindingService        *service.RoomBindingService
	RoomImageService          *service.RoomImageService
	FacilityService           *service.FacilityService
	RoomFacilityService       *service.RoomFacilityService
	CancellationPolicyService *service.CancellationPolicyService
}

// NewRoomHandler 创建房源处理器
func NewRoomHandler() *RoomHandler {
	return &RoomHandler{
		RoomTypeService:           &service.RoomTypeService{},
		RoomInfoService:           &service.RoomInfoService{},
		RoomStatusService:         &service.RoomStatusService{},
		RoomBindingService:        &service.RoomBindingService{},
		RoomImageService:          &service.RoomImageService{},
		FacilityService:           &service.FacilityService{},
		RoomFacilityService:       &service.RoomFacilityService{},
		CancellationPolicyService: &service.CancellationPolicyService{},
	}
}

// CreateRoomType 创建房型字典
// @Summary 创建房型字典
// @Description 支持添加大床房、商务大床房、标准间等多种房型
// @Tags 房型管理
// @Accept json
// @Produce json
// @Param room_type body service.CreateRoomTypeReq true "房型信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-types [post]
func (h *RoomHandler) CreateRoomType(c *gin.Context) {
	var req service.CreateRoomTypeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	roomType, err := h.RoomTypeService.CreateRoomType(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": roomType})
}

// UpdateRoomType 更新房型字典
// @Summary 更新房型字典
// @Tags 房型管理
// @Accept json
// @Produce json
// @Param id path int true "房型ID"
// @Param room_type body service.UpdateRoomTypeReq true "房型信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-types/{id} [put]
func (h *RoomHandler) UpdateRoomType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	var req service.UpdateRoomTypeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	if err := h.RoomTypeService.UpdateRoomType(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功"})
}

// GetRoomType 获取房型详情
// @Summary 获取房型详情
// @Tags 房型管理
// @Produce json
// @Param id path int true "房型ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-types/{id} [get]
func (h *RoomHandler) GetRoomType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	roomType, err := h.RoomTypeService.GetRoomType(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": roomType})
}

// ListRoomTypes 获取房型列表
// @Summary 获取房型列表
// @Tags 房型管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param status query string false "状态筛选"
// @Param keyword query string false "关键词搜索"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-types [get]
func (h *RoomHandler) ListRoomTypes(c *gin.Context) {
	var req service.ListRoomTypeReq
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

	roomTypes, total, err := h.RoomTypeService.ListRoomTypes(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":      roomTypes,
			"total":     total,
			"page":      req.Page,
			"page_size": req.PageSize,
		},
	})
}

// DeleteRoomType 删除房型字典
// @Summary 删除房型字典
// @Tags 房型管理
// @Produce json
// @Param id path int true "房型ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-types/{id} [delete]
func (h *RoomHandler) DeleteRoomType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := h.RoomTypeService.DeleteRoomType(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功"})
}

// CreateRoomInfo 创建房源信息
// @Summary 创建房源信息
// @Description 支持设置房源名称、门市价、日历价、房间数量、房间号、面积、床型、是否含早/洗漱用品等基础属性
// @Tags 房源管理
// @Accept json
// @Produce json
// @Param room_info body service.CreateRoomInfoReq true "房源信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos [post]
func (h *RoomHandler) CreateRoomInfo(c *gin.Context) {
	var req service.CreateRoomInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	// 从上下文获取创建人ID（实际项目中应该从token中获取）
	// 如果前端没有传递 created_by，则从上下文获取或使用默认值
	if req.CreatedBy == 0 {
		if userID, exists := c.Get("user_id"); exists {
			if uid, ok := userID.(uint64); ok {
				req.CreatedBy = uid
			} else {
				req.CreatedBy = 1 // 默认值，实际项目中应该从token中获取
			}
		} else {
			req.CreatedBy = 1 // 默认值，实际项目中应该从token中获取
		}
	}

	roomInfo, err := h.RoomInfoService.CreateRoomInfo(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": roomInfo})
}

// UpdateRoomInfo 更新房源信息
// @Summary 更新房源信息
// @Tags 房源管理
// @Accept json
// @Produce json
// @Param id path int true "房源ID"
// @Param room_info body service.UpdateRoomInfoReq true "房源信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos/{id} [put]
func (h *RoomHandler) UpdateRoomInfo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	var req service.UpdateRoomInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	if err := h.RoomInfoService.UpdateRoomInfo(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功"})
}

// GetRoomInfo 获取房源详情
// @Summary 获取房源详情
// @Tags 房源管理
// @Produce json
// @Param id path int true "房源ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos/{id} [get]
func (h *RoomHandler) GetRoomInfo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	roomInfo, err := h.RoomInfoService.GetRoomInfo(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": roomInfo})
}

// ListRoomInfos 获取房源列表
// @Summary 获取房源列表
// @Tags 房源管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param branch_id query int false "分店ID"
// @Param room_type_id query int false "房型ID"
// @Param status query string false "状态筛选"
// @Param keyword query string false "关键词搜索"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos [get]
func (h *RoomHandler) ListRoomInfos(c *gin.Context) {
	var req service.ListRoomInfoReq
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

	roomInfos, total, err := h.RoomInfoService.ListRoomInfos(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":      roomInfos,
			"total":     total,
			"page":      req.Page,
			"page_size": req.PageSize,
		},
	})
}

// DeleteRoomInfo 删除房源信息
// @Summary 删除房源信息
// @Tags 房源管理
// @Produce json
// @Param id path int true "房源ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos/{id} [delete]
func (h *RoomHandler) DeleteRoomInfo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := h.RoomInfoService.DeleteRoomInfo(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功"})
}
