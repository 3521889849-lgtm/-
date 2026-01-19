package handler

import (
	"example_shop/api/client"
	"example_shop/kitex_gen/hotel"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SystemConfigHandler struct{}

func NewSystemConfigHandler() *SystemConfigHandler {
	return &SystemConfigHandler{}
}

// CreateSystemConfig 创建系统配置
func (h *SystemConfigHandler) CreateSystemConfig(c *gin.Context) {
	var req hotel.CreateSystemConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "请求参数错误: "+err.Error())
		return
	}

	resp, err := client.HotelClient.CreateSystemConfig(c.Request.Context(), &req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateSystemConfig 更新系统配置
func (h *SystemConfigHandler) UpdateSystemConfig(c *gin.Context) {
	var req hotel.UpdateSystemConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "请求参数错误: "+err.Error())
		return
	}

	resp, err := client.HotelClient.UpdateSystemConfig(c.Request.Context(), &req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetSystemConfig 获取系统配置详情
func (h *SystemConfigHandler) GetSystemConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		badRequest(c, "参数错误")
		return
	}

	resp, err := client.HotelClient.GetSystemConfig(c.Request.Context(), id)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListSystemConfigs 获取系统配置列表
func (h *SystemConfigHandler) ListSystemConfigs(c *gin.Context) {
	req := &hotel.ListSystemConfigsReq{
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
	if configCategory := c.Query("config_category"); configCategory != "" {
		req.ConfigCategory = &configCategory
	}
	if status := c.Query("status"); status != "" {
		req.Status = &status
	}
	if keyword := c.Query("keyword"); keyword != "" {
		req.Keyword = &keyword
	}

	resp, err := client.HotelClient.ListSystemConfigs(c.Request.Context(), req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteSystemConfig 删除系统配置
func (h *SystemConfigHandler) DeleteSystemConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		badRequest(c, "参数错误")
		return
	}

	resp, err := client.HotelClient.DeleteSystemConfig(c.Request.Context(), id)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetSystemConfigsByCategory 按分类获取系统配置
func (h *SystemConfigHandler) GetSystemConfigsByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		badRequest(c, "参数错误")
		return
	}

	resp, err := client.HotelClient.GetSystemConfigsByCategory(c.Request.Context(), category)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}
