package service

import (
	"errors"
	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"
)

// RoomInfoService 房源信息管理服务
// 负责处理酒店房源的创建、更新、查询、删除等核心业务逻辑，
// 包括房源与分店、房型的关联校验，房间号唯一性检查，以及房源与订单的关联关系检查等。
type RoomInfoService struct{}

// CreateRoomInfo 创建房源信息
// 业务功能：在指定分店下创建新的房源记录，用于酒店客房信息管理

func (s *RoomInfoService) CreateRoomInfo(req *CreateRoomInfoReq) (*hotel_admin.RoomInfo, error) {
	// 业务规则：房源必须关联到有效的分店，验证分店是否存在且未被删除
	var branch hotel_admin.HotelBranch                                    // 声明分店实体变量，用于存储查询到的分店信息
	if err := db.MysqlDB.First(&branch, req.BranchID).Error; err != nil { // 通过分店ID查询分店信息，如果查询失败则说明分店不存在
		return nil, errors.New("分店不存在") // 返回错误信息，表示分店不存在
	}

	// 业务规则：房源必须关联到有效的房型，验证房型字典中是否存在该房型定义
	var roomType hotel_admin.RoomTypeDict                                     // 声明房型实体变量，用于存储查询到的房型信息
	if err := db.MysqlDB.First(&roomType, req.RoomTypeID).Error; err != nil { // 通过房型ID查询房型信息，如果查询失败则说明房型不存在
		return nil, errors.New("房型不存在") // 返回错误信息，表示房型不存在
	}

	// 业务规则：同一分店下房间号必须唯一，检查该房间号在当前分店中是否已被使用（排除已删除记录）
	var existRoom hotel_admin.RoomInfo                                                                                             // 声明房源实体变量，用于存储查询到的已存在房源信息
	result := db.MysqlDB.Where("branch_id = ? AND room_no = ? AND deleted_at IS NULL", req.BranchID, req.RoomNo).First(&existRoom) // 查询指定分店下是否已存在相同的房间号（排除已删除记录）
	if result.Error == nil {                                                                                                       // 如果查询成功（未报错），说明该房间号已存在
		return nil, errors.New("该房间号已存在") // 返回错误信息，表示房间号重复
	}

	// 构建房源实体对象，设置默认状态为启用，其他字段从请求中获取
	roomInfo := &hotel_admin.RoomInfo{ // 创建房源实体对象指针
		BranchID:             req.BranchID,
		RoomTypeID:           req.RoomTypeID,
		RoomNo:               req.RoomNo,
		RoomName:             req.RoomName,
		MarketPrice:          req.MarketPrice,
		CalendarPrice:        req.CalendarPrice,
		RoomCount:            req.RoomCount,
		Area:                 req.Area,
		BedSpec:              req.BedSpec,
		HasBreakfast:         req.HasBreakfast,
		HasToiletries:        req.HasToiletries,
		CancellationPolicyID: req.CancellationPolicyID,
		Status:               "ACTIVE",
		CreatedBy:            req.CreatedBy,
	}

	// 持久化房源信息到数据库
	if err := db.MysqlDB.Create(roomInfo).Error; err != nil { // 将房源信息保存到数据库，如果保存失败则返回错误
		return nil, err // 返回数据库操作错误
	}

	return roomInfo, nil // 返回成功创建后的房源信息，无错误
}

// UpdateRoomInfo 更新房源信息
// 业务功能：修改已存在房源的属性信息，支持部分字段更新

func (s *RoomInfoService) UpdateRoomInfo(id uint64, req *UpdateRoomInfoReq) error {
	// 业务规则：只能更新已存在的房源，验证房源是否存在
	var roomInfo hotel_admin.RoomInfo                             // 声明房源实体变量，用于存储查询到的房源信息
	if err := db.MysqlDB.First(&roomInfo, id).Error; err != nil { // 通过房源ID查询房源信息，如果查询失败则说明房源不存在
		return errors.New("房源不存在") // 返回错误信息，表示房源不存在
	}

	// 业务规则：如果修改房间号，需重新校验唯一性，确保新房间号在当前分店下未被其他房源使用
	if req.RoomNo != "" && req.RoomNo != roomInfo.RoomNo { // 如果请求中提供了房间号且与当前房间号不同，则需要重新校验
		var existRoom hotel_admin.RoomInfo                                                                                                                  // 声明房源实体变量，用于存储查询到的已存在房源信息
		result := db.MysqlDB.Where("branch_id = ? AND room_no = ? AND id != ? AND deleted_at IS NULL", roomInfo.BranchID, req.RoomNo, id).First(&existRoom) // 查询同一分店下是否已存在相同房间号的其他房源（排除当前房源和已删除记录）
		if result.Error == nil {                                                                                                                            // 如果查询成功（未报错），说明该房间号已被其他房源使用
			return errors.New("该房间号已存在") // 返回错误信息，表示房间号重复
		}
		roomInfo.RoomNo = req.RoomNo // 房间号校验通过，更新房源房间号
	}

	// 业务逻辑：采用部分更新策略，只更新请求中提供的非空字段，保持其他字段不变
	if req.RoomName != "" { // 如果请求中提供了房间名称（非空），则更新
		roomInfo.RoomName = req.RoomName // 更新房间名称
	}
	if req.MarketPrice > 0 { // 如果请求中提供了门市价（大于0），则更新
		roomInfo.MarketPrice = req.MarketPrice // 更新门市价
	}
	if req.CalendarPrice > 0 { // 如果请求中提供了日历价（大于0），则更新
		roomInfo.CalendarPrice = req.CalendarPrice // 更新日历价
	}
	if req.RoomCount > 0 { // 如果请求中提供了房间数量（大于0），则更新
		roomInfo.RoomCount = req.RoomCount // 更新房间数量
	}
	if req.Area != nil { // 如果请求中提供了面积（指针非空），则更新
		roomInfo.Area = req.Area // 更新面积
	}
	if req.BedSpec != "" { // 如果请求中提供了床型规格（非空），则更新
		roomInfo.BedSpec = req.BedSpec // 更新床型规格
	}
	if req.HasBreakfast != nil { // 如果请求中提供了是否含早餐（指针非空），则更新
		roomInfo.HasBreakfast = *req.HasBreakfast // 更新是否含早餐（解引用指针获取值）
	}
	if req.HasToiletries != nil { // 如果请求中提供了是否提供洗漱用品（指针非空），则更新
		roomInfo.HasToiletries = *req.HasToiletries // 更新是否提供洗漱用品（解引用指针获取值）
	}
	if req.CancellationPolicyID != nil { // 如果请求中提供了退订政策ID（指针非空），则更新
		roomInfo.CancellationPolicyID = req.CancellationPolicyID // 更新退订政策ID
	}
	if req.Status != "" { // 如果请求中提供了状态（非空），则更新
		roomInfo.Status = req.Status // 更新房源状态
	}

	// 保存更新后的房源信息
	return db.MysqlDB.Save(&roomInfo).Error // 将更新后的房源信息保存到数据库，返回保存操作的结果（成功为nil，失败为error）
}

// GetRoomInfo 获取房源详情
// 业务功能：根据房源ID查询房源的完整信息，包含关联的分店、房型、退订政策、设施、图片等数据

func (s *RoomInfoService) GetRoomInfo(id uint64) (*hotel_admin.RoomInfo, error) {
	// 使用预加载（Preload）机制一次性加载所有关联数据，避免N+1查询问题
	// 关联数据包括：所属分店、房型字典、退订政策、设施关联关系（含设施详情）、房源图片列表
	var roomInfo hotel_admin.RoomInfo                                                         // 声明房源实体变量，用于存储查询到的房源信息
	if err := db.MysqlDB.Preload("Branch").Preload("RoomType").Preload("CancellationPolicy"). // 预加载所属分店、房型字典、退订政策关联数据
		Preload("RoomFacilityRelations.Facility").Preload("RoomImages"). // 预加载设施关联关系（含设施详情）、房源图片列表关联数据
		First(&roomInfo, id).Error; err != nil { // 通过房源ID查询房源信息（包含所有预加载的关联数据），如果查询失败则说明房源不存在
		return nil, errors.New("房源不存在") // 返回nil和错误信息，表示房源不存在
	}
	return &roomInfo, nil // 返回房源信息指针和无错误
}

// ListRoomInfos 获取房源列表
// 业务功能：支持多条件筛选和分页查询房源列表，用于房源管理和搜索场景

func (s *RoomInfoService) ListRoomInfos(req *ListRoomInfoReq) ([]hotel_admin.RoomInfo, int64, error) {
	var roomInfos []hotel_admin.RoomInfo // 声明房源列表变量，用于存储查询到的房源信息列表
	var total int64                      // 声明总数变量，用于存储符合条件的房源总数

	query := db.MysqlDB.Model(&hotel_admin.RoomInfo{}) // 创建房源模型的查询构建器

	// 业务筛选：按分店ID筛选，支持查看特定分店下的所有房源
	if req.BranchID > 0 { // 如果请求中提供了分店ID（大于0），则添加分店筛选条件
		query = query.Where("branch_id = ?", req.BranchID) // 添加分店ID筛选条件，只查询指定分店下的房源
	}

	// 业务筛选：按房型ID筛选，支持查看特定房型的所有房源
	if req.RoomTypeID > 0 { // 如果请求中提供了房型ID（大于0），则添加房型筛选条件
		query = query.Where("room_type_id = ?", req.RoomTypeID) // 添加房型ID筛选条件，只查询指定房型的房源
	}

	// 业务筛选：按状态筛选（如启用/停用/维护中等），支持状态维度的房源管理
	if req.Status != "" { // 如果请求中提供了状态（非空），则添加状态筛选条件
		query = query.Where("status = ?", req.Status) // 添加状态筛选条件，只查询指定状态的房源
	}

	// 业务搜索：支持模糊搜索房间名称或房间号，提升房源查找效率
	if req.Keyword != "" { // 如果请求中提供了关键词（非空），则添加关键词搜索条件
		query = query.Where("room_name LIKE ? OR room_no LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%") // 添加模糊搜索条件，搜索房间名称或房间号包含关键词的房源
	}

	// 先统计符合条件的房源总数，用于前端分页组件显示总条数
	if err := query.Count(&total).Error; err != nil { // 统计符合条件的房源总数，如果统计失败则返回错误
		return nil, 0, err // 返回nil列表、0总数和错误信息
	}

	// 分页查询：按创建时间倒序排列，加载关联的分店和房型信息，避免后续单独查询
	offset := (req.Page - 1) * req.PageSize                // 计算分页偏移量（跳过前面的记录数）
	if err := query.Preload("Branch").Preload("RoomType"). // 预加载所属分店和房型字典关联数据
		Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&roomInfos).Error; err != nil { // 添加分页限制（偏移量、每页数量）、排序（按创建时间倒序）并查询房源列表，如果查询失败则返回错误
		return nil, 0, err // 返回nil列表、0总数和错误信息
	}

	return roomInfos, total, nil // 返回房源列表、总数和无错误
}

// DeleteRoomInfo 删除房源信息（软删除）
// 业务功能：逻辑删除房源记录，不物理删除数据，保留历史订单关联关系

func (s *RoomInfoService) DeleteRoomInfo(id uint64) error {
	// 业务规则：只能删除已存在的房源，验证房源是否存在
	var roomInfo hotel_admin.RoomInfo                             // 声明房源实体变量，用于存储查询到的房源信息
	if err := db.MysqlDB.First(&roomInfo, id).Error; err != nil { // 通过房源ID查询房源信息，如果查询失败则说明房源不存在
		return errors.New("房源不存在") // 返回错误信息，表示房源不存在
	}

	// 业务规则：如果房源已被订单使用，不允许删除，避免破坏订单数据的完整性
	var count int64                                                                                          // 声明计数变量，用于存储关联订单的数量
	db.MysqlDB.Model(&hotel_admin.OrderMain{}).Where("room_id = ? AND deleted_at IS NULL", id).Count(&count) // 统计使用该房源的订单数量（排除已删除订单）
	if count > 0 {                                                                                           // 如果存在关联订单（数量大于0），则不允许删除
		return errors.New("该房源正在被订单使用，无法删除") // 返回错误信息，表示房源正在被使用
	}

	// 执行软删除：设置deleted_at字段，不物理删除记录
	return db.MysqlDB.Delete(&roomInfo).Error // 执行软删除操作（设置deleted_at字段），返回删除操作的结果（成功为nil，失败为error）
}

// CreateRoomInfoReq 创建房源请求
type CreateRoomInfoReq struct {
	BranchID             uint64   `json:"branch_id" binding:"required"`      // 分店ID
	RoomTypeID           uint64   `json:"room_type_id" binding:"required"`   // 房型ID
	RoomNo               string   `json:"room_no" binding:"required"`        // 房间号
	RoomName             string   `json:"room_name" binding:"required"`      // 房间名称
	MarketPrice          float64  `json:"market_price" binding:"required"`   // 门市价
	CalendarPrice        float64  `json:"calendar_price" binding:"required"` // 日历价
	RoomCount            uint8    `json:"room_count" binding:"required"`     // 房间数量
	Area                 *float64 `json:"area,omitempty"`                    // 面积
	BedSpec              string   `json:"bed_spec" binding:"required"`       // 床型规格
	HasBreakfast         bool     `json:"has_breakfast"`                     // 是否含早
	HasToiletries        bool     `json:"has_toiletries"`                    // 是否提供洗漱用品
	CancellationPolicyID *uint64  `json:"cancellation_policy_id,omitempty"`  // 退订政策ID
	RoomNos              []string `json:"room_nos,omitempty"`                // 房间号列表（用于关联房间表）
	CreatedBy            uint64   `json:"created_by" binding:"required"`     // 创建人
}

// UpdateRoomInfoReq 更新房源请求
type UpdateRoomInfoReq struct {
	RoomNo               string   `json:"room_no,omitempty"`                // 房间号
	RoomName             string   `json:"room_name,omitempty"`              // 房间名称
	MarketPrice          float64  `json:"market_price"`                     // 门市价
	CalendarPrice        float64  `json:"calendar_price"`                   // 日历价
	RoomCount            uint8    `json:"room_count"`                       // 房间数量
	Area                 *float64 `json:"area,omitempty"`                   // 面积
	BedSpec              string   `json:"bed_spec,omitempty"`               // 床型规格
	HasBreakfast         *bool    `json:"has_breakfast,omitempty"`          // 是否含早
	HasToiletries        *bool    `json:"has_toiletries,omitempty"`         // 是否提供洗漱用品
	CancellationPolicyID *uint64  `json:"cancellation_policy_id,omitempty"` // 退订政策ID
	Status               string   `json:"status,omitempty"`                 // 状态
}

// ListRoomInfoReq 房源列表请求
type ListRoomInfoReq struct {
	Page       int    `json:"page" form:"page"`                 // 页码
	PageSize   int    `json:"page_size" form:"page_size"`       // 每页数量
	BranchID   uint64 `json:"branch_id" form:"branch_id"`       // 分店ID
	RoomTypeID uint64 `json:"room_type_id" form:"room_type_id"` // 房型ID
	Status     string `json:"status" form:"status"`             // 状态筛选
	Keyword    string `json:"keyword" form:"keyword"`           // 关键词搜索
}
