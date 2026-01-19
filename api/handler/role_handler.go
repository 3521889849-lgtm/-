package handler

import (
	"example_shop/api/client"
	"example_shop/kitex_gen/hotel"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct{}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{}
}

// CreateRole 创建角色
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req hotel.CreateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "请求参数错误: "+err.Error())
		return
	}

	resp, err := client.HotelClient.CreateRole(c.Request.Context(), &req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateRole 更新角色
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	var req hotel.UpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "请求参数错误: "+err.Error())
		return
	}

	resp, err := client.HotelClient.UpdateRole(c.Request.Context(), &req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetRole 获取角色详情
func (h *RoleHandler) GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		badRequest(c, "参数错误")
		return
	}

	resp, err := client.HotelClient.GetRole(c.Request.Context(), id)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListRoles 获取角色列表
func (h *RoleHandler) ListRoles(c *gin.Context) {
	req := &hotel.ListRolesReq{
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

	resp, err := client.HotelClient.ListRoles(c.Request.Context(), req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteRole 删除角色
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		badRequest(c, "参数错误")
		return
	}

	resp, err := client.HotelClient.DeleteRole(c.Request.Context(), id)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListPermissions 获取权限列表
func (h *RoleHandler) ListPermissions(c *gin.Context) {
	req := &hotel.ListPermissionsReq{}

	if permissionType := c.Query("permission_type"); permissionType != "" {
		req.PermissionType = &permissionType
	}
	if parentIDStr := c.Query("parent_id"); parentIDStr != "" {
		if parentID, err := strconv.ParseInt(parentIDStr, 10, 64); err == nil {
			req.ParentId = &parentID
		}
	}
	if status := c.Query("status"); status != "" {
		req.Status = &status
	}

	resp, err := client.HotelClient.ListPermissions(c.Request.Context(), req)
	if err != nil {
		internalError(c, "RPC 调用失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}
