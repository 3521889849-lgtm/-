package handler

import (
	"context"
	"example_shop/api/client"
	"example_shop/kitex_gen/hotel"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListInHouseGuests 获取在住客人列表
// @Summary 获取在住客人列表
// @Description 支持按省份、姓名、手机号、身份证号、房间号等条件筛选在住客人信息
// @Tags 客人管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(200)
// @Param branch_id query int false "分店ID"
// @Param province query string false "省份"
// @Param city query string false "城市"
// @Param district query string false "区县"
// @Param name query string false "姓名"
// @Param phone query string false "手机号"
// @Param id_number query string false "身份证号"
// @Param room_no query string false "房间号"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/in-house-guests [get]
func (h *RoomHandler) ListInHouseGuests(c *gin.Context) {
	rpcReq := &hotel.ListInHouseGuestsReq{
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

	if province := c.Query("province"); province != "" {
		rpcReq.Province = &province
	}

	if city := c.Query("city"); city != "" {
		rpcReq.City = &city
	}

	if district := c.Query("district"); district != "" {
		rpcReq.District = &district
	}

	if name := c.Query("name"); name != "" {
		rpcReq.Name = &name
	}

	if phone := c.Query("phone"); phone != "" {
		rpcReq.Phone = &phone
	}

	if idNumber := c.Query("id_number"); idNumber != "" {
		rpcReq.IdNumber = &idNumber
	}

	if roomNo := c.Query("room_no"); roomNo != "" {
		rpcReq.RoomNo = &roomNo
	}

	resp, err := client.HotelClient.ListInHouseGuests(context.Background(), rpcReq)
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
