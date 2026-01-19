package handler

import (
	"context"
	"example_shop/api/client"
	"example_shop/kitex_gen/hotel"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListOrders 获取订单列表
// @Summary 获取订单列表
// @Description 支持按客人来源、订单号、手机号、入住/离店时间等多条件查询
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param branch_id query int false "分店ID"
// @Param guest_source query string false "客人来源"
// @Param order_no query string false "订单号"
// @Param phone query string false "手机号"
// @Param keyword query string false "关键词（订单号/房间号/手机号/联系人）"
// @Param order_status query string false "订单状态"
// @Param check_in_start query string false "入住开始时间 YYYY-MM-DD"
// @Param check_in_end query string false "入住结束时间 YYYY-MM-DD"
// @Param check_out_start query string false "离店开始时间 YYYY-MM-DD"
// @Param check_out_end query string false "离店结束时间 YYYY-MM-DD"
// @Param reserve_start query string false "预定开始时间 YYYY-MM-DD HH:mm:ss"
// @Param reserve_end query string false "预定结束时间 YYYY-MM-DD HH:mm:ss"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/orders [get]
func (h *RoomHandler) ListOrders(c *gin.Context) {
	rpcReq := &hotel.ListOrdersReq{
		Page:     1,
		PageSize: 10,
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			rpcReq.Page = int32(page)
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			rpcReq.PageSize = int32(pageSize)
		}
	}

	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		if branchID, err := strconv.ParseInt(branchIDStr, 10, 64); err == nil {
			rpcReq.BranchId = &branchID
		}
	}

	if guestSource := c.Query("guest_source"); guestSource != "" {
		rpcReq.GuestSource = &guestSource
	}

	if orderNo := c.Query("order_no"); orderNo != "" {
		rpcReq.OrderNo = &orderNo
	}

	if phone := c.Query("phone"); phone != "" {
		rpcReq.Phone = &phone
	}

	if keyword := c.Query("keyword"); keyword != "" {
		rpcReq.Keyword = &keyword
	}

	if orderStatus := c.Query("order_status"); orderStatus != "" {
		rpcReq.OrderStatus = &orderStatus
	}

	if checkInStart := c.Query("check_in_start"); checkInStart != "" {
		rpcReq.CheckInStart = &checkInStart
	}

	if checkInEnd := c.Query("check_in_end"); checkInEnd != "" {
		rpcReq.CheckInEnd = &checkInEnd
	}

	if checkOutStart := c.Query("check_out_start"); checkOutStart != "" {
		rpcReq.CheckOutStart = &checkOutStart
	}

	if checkOutEnd := c.Query("check_out_end"); checkOutEnd != "" {
		rpcReq.CheckOutEnd = &checkOutEnd
	}

	if reserveStart := c.Query("reserve_start"); reserveStart != "" {
		rpcReq.ReserveStart = &reserveStart
	}

	if reserveEnd := c.Query("reserve_end"); reserveEnd != "" {
		rpcReq.ReserveEnd = &reserveEnd
	}

	resp, err := client.HotelClient.ListOrders(context.Background(), rpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "RPC 调用失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": resp,
	})
}

// GetOrder 获取订单详情
// @Summary 获取订单详情
// @Description 查看订单详情，关联房间号与房型信息，记录订单关键节点数据
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/orders/{id} [get]
func (h *RoomHandler) GetOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	resp, err := client.HotelClient.GetOrder(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "RPC 调用失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": resp,
	})
}
