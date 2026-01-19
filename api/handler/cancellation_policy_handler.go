package handler

import (
	"example_shop/common/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateCancellationPolicy 创建退订政策
// @Summary 创建退订政策
// @Description 支持自定义退订规则（如"入住前X小时内不可取消，否则收取X倍房费"）
// @Tags 退订政策管理
// @Accept json
// @Produce json
// @Param policy body service.CreateCancellationPolicyReq true "退订政策信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/cancellation-policies [post]
func (h *RoomHandler) CreateCancellationPolicy(c *gin.Context) {
	var req service.CreateCancellationPolicyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	policy, err := h.CancellationPolicyService.CreateCancellationPolicy(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": policy})
}

// UpdateCancellationPolicy 更新退订政策
// @Summary 更新退订政策
// @Tags 退订政策管理
// @Accept json
// @Produce json
// @Param id path int true "政策ID"
// @Param policy body service.UpdateCancellationPolicyReq true "退订政策信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/cancellation-policies/{id} [put]
func (h *RoomHandler) UpdateCancellationPolicy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	var req service.UpdateCancellationPolicyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	if err := h.CancellationPolicyService.UpdateCancellationPolicy(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功"})
}

// GetCancellationPolicy 获取退订政策详情
// @Summary 获取退订政策详情
// @Tags 退订政策管理
// @Produce json
// @Param id path int true "政策ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/cancellation-policies/{id} [get]
func (h *RoomHandler) GetCancellationPolicy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	policy, err := h.CancellationPolicyService.GetCancellationPolicy(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": policy})
}

// ListCancellationPolicies 获取退订政策列表
// @Summary 获取退订政策列表
// @Tags 退订政策管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param room_type_id query int false "房型ID筛选"
// @Param status query string false "状态筛选"
// @Param keyword query string false "关键词搜索"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/cancellation-policies [get]
func (h *RoomHandler) ListCancellationPolicies(c *gin.Context) {
	var req service.ListCancellationPolicyReq
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

	policies, total, err := h.CancellationPolicyService.ListCancellationPolicies(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":      policies,
			"total":     total,
			"page":      req.Page,
			"page_size": req.PageSize,
		},
	})
}

// DeleteCancellationPolicy 删除退订政策
// @Summary 删除退订政策
// @Tags 退订政策管理
// @Produce json
// @Param id path int true "政策ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/cancellation-policies/{id} [delete]
func (h *RoomHandler) DeleteCancellationPolicy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := h.CancellationPolicyService.DeleteCancellationPolicy(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功"})
}
