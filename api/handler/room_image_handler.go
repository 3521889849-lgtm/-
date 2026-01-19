package handler

import (
	"context"
	"example_shop/api/client"
	"example_shop/kitex_gen/hotel"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UploadRoomImages 批量上传房源图片
// @Summary 批量上传房源图片
// @Description 支持批量上传房源图片（≤16张，规格400x300，格式jpg/png）
// @Tags 房源管理
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "房源ID"
// @Param images formData file true "图片文件" multiple
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos/{id}/images [post]
func (h *RoomHandler) UploadRoomImages(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "获取表单失败: " + err.Error()})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请选择要上传的图片"})
		return
	}

	uploadedImages, err := h.RoomImageService.UploadRoomImages(id, files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "上传成功", "data": uploadedImages})
}

// GetRoomImages 获取房源图片列表
// @Summary 获取房源图片列表
// @Tags 房源管理
// @Produce json
// @Param id path int true "房源ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos/{id}/images [get]
func (h *RoomHandler) GetRoomImages(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	resp, err := client.HotelClient.GetRoomImages(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "RPC 调用失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": resp.Images})
}

// DeleteRoomImage 删除房源图片
// @Summary 删除房源图片
// @Tags 房源管理
// @Produce json
// @Param id path int true "图片ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos/images/{id} [delete]
func (h *RoomHandler) DeleteRoomImage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	resp, err := client.HotelClient.DeleteRoomImage(context.Background(), id)
	code, msg := handleRPCError(resp, err)
	if code != 200 {
		c.JSON(code, gin.H{"code": code, "msg": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": msg})
}

// UpdateImageSortOrder 更新图片排序
// @Summary 更新图片排序
// @Tags 房源管理
// @Accept json
// @Produce json
// @Param id path int true "图片ID"
// @Param request body map[string]interface{} true "排序信息" example({"sort_order":1})
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos/images/{id}/sort [put]
func (h *RoomHandler) UpdateImageSortOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	var req struct {
		SortOrder int8 `json:"sort_order" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	rpcReq := &hotel.UpdateImageSortOrderReq{
		ImageId:   id,
		SortOrder: req.SortOrder,
	}

	resp, err := client.HotelClient.UpdateImageSortOrder(context.Background(), rpcReq)
	code, msg := handleRPCError(resp, err)
	if code != 200 {
		c.JSON(code, gin.H{"code": code, "msg": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": msg})
}

// BatchUpdateImageSortOrder 批量更新图片排序
// @Summary 批量更新图片排序
// @Tags 房源管理
// @Accept json
// @Produce json
// @Param id path int true "房源ID"
// @Param request body map[string]interface{} true "排序列表" example({"sort_orders":[{"image_id":1,"sort_order":0},{"image_id":2,"sort_order":1}]})
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/room-infos/{id}/images/sort [put]
func (h *RoomHandler) BatchUpdateImageSortOrder(c *gin.Context) {
	roomID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	var req struct {
		SortOrders []struct {
			ImageID  int64 `json:"image_id"`
			SortOrder int8  `json:"sort_order"`
		} `json:"sort_orders" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误: " + err.Error()})
		return
	}

	sortOrders := make([]*hotel.ImageSortOrder, len(req.SortOrders))
	for i, so := range req.SortOrders {
		sortOrders[i] = &hotel.ImageSortOrder{
			ImageId:   so.ImageID,
			SortOrder: so.SortOrder,
		}
	}

	rpcReq := &hotel.BatchUpdateImageSortOrderReq{
		RoomId:     roomID,
		SortOrders: sortOrders,
	}

	resp, err := client.HotelClient.BatchUpdateImageSortOrder(context.Background(), rpcReq)
	code, msg := handleRPCError(resp, err)
	if code != 200 {
		c.JSON(code, gin.H{"code": code, "msg": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": msg})
}
