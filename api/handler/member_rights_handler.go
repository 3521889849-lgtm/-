package handler

import (
	"example_shop/api/client"
	"example_shop/kitex_gen/hotel"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MemberRightsHandler struct{}

func NewMemberRightsHandler() *MemberRightsHandler {
	return &MemberRightsHandler{}
}

// CreateMemberRights 创建会员权益
func (h *MemberRightsHandler) CreateMemberRights(c *gin.Context) {
	var req hotel.CreateMemberRightsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "请求参数错误: "+err.Error())
		return
	}

	resp, err := client.HotelClient.CreateMemberRights(c.Request.Context(), &req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateMemberRights 更新会员权益
func (h *MemberRightsHandler) UpdateMemberRights(c *gin.Context) {
	var req hotel.UpdateMemberRightsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "请求参数错误: "+err.Error())
		return
	}

	resp, err := client.HotelClient.UpdateMemberRights(c.Request.Context(), &req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetMemberRights 获取会员权益详情
func (h *MemberRightsHandler) GetMemberRights(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		badRequest(c, "参数错误")
		return
	}

	resp, err := client.HotelClient.GetMemberRights(c.Request.Context(), id)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListMemberRights 获取会员权益列表
func (h *MemberRightsHandler) ListMemberRights(c *gin.Context) {
	req := &hotel.ListMemberRightsReq{
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
	if memberLevel := c.Query("member_level"); memberLevel != "" {
		req.MemberLevel = &memberLevel
	}
	if status := c.Query("status"); status != "" {
		req.Status = &status
	}
	if keyword := c.Query("keyword"); keyword != "" {
		req.Keyword = &keyword
	}

	resp, err := client.HotelClient.ListMemberRights(c.Request.Context(), req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetRightsByMemberLevel 根据会员等级获取权益列表
func (h *MemberRightsHandler) GetRightsByMemberLevel(c *gin.Context) {
	memberLevel := c.Param("member_level")
	if memberLevel == "" {
		badRequest(c, "参数错误")
		return
	}

	resp, err := client.HotelClient.GetRightsByMemberLevel(c.Request.Context(), memberLevel)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteMemberRights 删除会员权益
func (h *MemberRightsHandler) DeleteMemberRights(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		badRequest(c, "参数错误")
		return
	}

	resp, err := client.HotelClient.DeleteMemberRights(c.Request.Context(), id)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}
