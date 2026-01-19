package handler

import (
	"example_shop/api/client"
	"example_shop/kitex_gen/hotel"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MemberHandler struct{}

func NewMemberHandler() *MemberHandler {
	return &MemberHandler{}
}

// CreateMember 创建会员
func (h *MemberHandler) CreateMember(c *gin.Context) {
	var req hotel.CreateMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "请求参数错误: "+err.Error())
		return
	}

	resp, err := client.HotelClient.CreateMember(c.Request.Context(), &req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateMember 更新会员
func (h *MemberHandler) UpdateMember(c *gin.Context) {
	var req hotel.UpdateMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "请求参数错误: "+err.Error())
		return
	}

	resp, err := client.HotelClient.UpdateMember(c.Request.Context(), &req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetMember 获取会员详情
func (h *MemberHandler) GetMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		badRequest(c, "参数错误")
		return
	}

	resp, err := client.HotelClient.GetMember(c.Request.Context(), id)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetMemberByGuestID 根据客人ID获取会员信息
func (h *MemberHandler) GetMemberByGuestID(c *gin.Context) {
	guestIDStr := c.Param("guest_id")
	guestID, err := strconv.ParseInt(guestIDStr, 10, 64)
	if err != nil {
		badRequest(c, "参数错误")
		return
	}

	resp, err := client.HotelClient.GetMemberByGuestID(c.Request.Context(), guestID)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListMembers 获取会员列表
func (h *MemberHandler) ListMembers(c *gin.Context) {
	req := &hotel.ListMembersReq{
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

	resp, err := client.HotelClient.ListMembers(c.Request.Context(), req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteMember 删除会员
func (h *MemberHandler) DeleteMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		badRequest(c, "参数错误")
		return
	}

	resp, err := client.HotelClient.DeleteMember(c.Request.Context(), id)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}
