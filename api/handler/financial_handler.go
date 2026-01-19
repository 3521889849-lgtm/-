package handler

import (
	"context"
	"example_shop/api/client"
	"example_shop/kitex_gen/hotel"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListFinancialFlows 获取收支流水列表
// @Summary 获取收支流水列表
// @Description 支持按分店、收支类型、收支项目、支付方式、操作人、发生时间等条件筛选，包含财务汇总统计
// @Tags 财务管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(200)
// @Param branch_id query int false "分店ID"
// @Param flow_type query string false "收支类型（收入/支出）"
// @Param flow_item query string false "收支项目"
// @Param pay_type query string false "支付方式"
// @Param operator_id query int false "操作人ID"
// @Param occur_start query string false "发生开始时间 YYYY-MM-DD"
// @Param occur_end query string false "发生结束时间 YYYY-MM-DD"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/financial-flows [get]
func (h *RoomHandler) ListFinancialFlows(c *gin.Context) {
	rpcReq := &hotel.ListFinancialFlowsReq{
		Page:     1,
		PageSize: 200,
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

	if flowType := c.Query("flow_type"); flowType != "" {
		rpcReq.FlowType = &flowType
	}

	if flowItem := c.Query("flow_item"); flowItem != "" {
		rpcReq.FlowItem = &flowItem
	}

	if payType := c.Query("pay_type"); payType != "" {
		rpcReq.PayType = &payType
	}

	if operatorIDStr := c.Query("operator_id"); operatorIDStr != "" {
		if operatorID, err := strconv.ParseInt(operatorIDStr, 10, 64); err == nil {
			rpcReq.OperatorId = &operatorID
		}
	}

	if occurStart := c.Query("occur_start"); occurStart != "" {
		rpcReq.OccurStart = &occurStart
	}

	if occurEnd := c.Query("occur_end"); occurEnd != "" {
		rpcReq.OccurEnd = &occurEnd
	}

	resp, err := client.HotelClient.ListFinancialFlows(context.Background(), rpcReq)
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
