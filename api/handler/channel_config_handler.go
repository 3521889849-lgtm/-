package handler

import (
	"example_shop/api/client"
	"example_shop/kitex_gen/hotel"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChannelConfigHandler struct{}

func NewChannelConfigHandler() *ChannelConfigHandler {
	return &ChannelConfigHandler{}
}

// CreateChannelConfig 创建渠道配置
func (h *ChannelConfigHandler) CreateChannelConfig(c *gin.Context) {
	var req hotel.CreateChannelConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "请求参数错误: "+err.Error())
		return
	}

	resp, err := client.HotelClient.CreateChannelConfig(c.Request.Context(), &req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateChannelConfig 更新渠道配置
func (h *ChannelConfigHandler) UpdateChannelConfig(c *gin.Context) {
	var req hotel.UpdateChannelConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "请求参数错误: "+err.Error())
		return
	}

	resp, err := client.HotelClient.UpdateChannelConfig(c.Request.Context(), &req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetChannelConfig 获取渠道配置详情
func (h *ChannelConfigHandler) GetChannelConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		badRequest(c, "参数错误")
		return
	}

	resp, err := client.HotelClient.GetChannelConfig(c.Request.Context(), id)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListChannelConfigs 获取渠道配置列表
func (h *ChannelConfigHandler) ListChannelConfigs(c *gin.Context) {
	req := &hotel.ListChannelConfigsReq{
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
	if status := c.Query("status"); status != "" {
		req.Status = &status
	}
	if keyword := c.Query("keyword"); keyword != "" {
		req.Keyword = &keyword
	}

	resp, err := client.HotelClient.ListChannelConfigs(c.Request.Context(), req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteChannelConfig 删除渠道配置
func (h *ChannelConfigHandler) DeleteChannelConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		badRequest(c, "参数错误")
		return
	}

	resp, err := client.HotelClient.DeleteChannelConfig(c.Request.Context(), id)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}
