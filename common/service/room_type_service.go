package service

import (
	"errors"
	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"
	"fmt"
)

// RoomTypeService 房型字典管理服务
// 负责处理酒店房型字典的创建、更新、查询、删除等核心业务逻辑，
// 包括房型名称唯一性检查、房型与房源的关联关系管理等。
type RoomTypeService struct{}

// CreateRoomType 创建房型字典
func (s *RoomTypeService) CreateRoomType(req *CreateRoomTypeReq) (*hotel_admin.RoomTypeDict, error) {
	// 检查房型名称是否已存在
	var count int64                                                                                                                               // 声明计数变量，用于存储查询到的已存在房型数量
	err := db.MysqlDB.Model(&hotel_admin.RoomTypeDict{}).Where("room_type_name = ? AND deleted_at IS NULL", req.RoomTypeName).Count(&count).Error // 统计是否已存在相同的房型名称（排除已删除记录），如果查询失败则说明数据库操作异常
	if err != nil {                                                                                                                               // 如果查询失败，说明数据库操作异常
		return nil, fmt.Errorf("查询房型失败: %w", err) // 返回nil和错误信息
	}
	if count > 0 { // 如果记录数量大于0，说明该房型名称已存在
		return nil, errors.New("房型名称已存在") // 返回nil和错误信息，表示房型名称已存在
	}

	roomType := &hotel_admin.RoomTypeDict{ // 创建房型实体对象指针
		RoomTypeName:  req.RoomTypeName,  // 设置房型名称（从请求中获取）
		BedSpec:       req.BedSpec,       // 设置床型规格（从请求中获取）
		Area:          req.Area,          // 设置面积（从请求中获取，可为空）
		HasBreakfast:  req.HasBreakfast,  // 设置是否含早（从请求中获取）
		HasToiletries: req.HasToiletries, // 设置是否提供洗漱用品（从请求中获取）
		DefaultPrice:  req.DefaultPrice,  // 设置默认价格（从请求中获取）
		Status:        "ACTIVE",          // 设置房型状态为启用（默认值）
	}

	if err := db.MysqlDB.Create(roomType).Error; err != nil { // 将房型信息保存到数据库，如果保存失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	return roomType, nil // 返回房型实体指针和无错误
}

// UpdateRoomType 更新房型字典
func (s *RoomTypeService) UpdateRoomType(id uint64, req *UpdateRoomTypeReq) error {
	var roomType hotel_admin.RoomTypeDict                         // 声明房型实体变量，用于存储查询到的房型信息
	if err := db.MysqlDB.First(&roomType, id).Error; err != nil { // 通过房型ID查询房型信息，如果查询失败则说明房型不存在
		return errors.New("房型不存在") // 返回错误信息，表示房型不存在
	}

	// 如果修改了房型名称，检查是否重复
	if req.RoomTypeName != nil && *req.RoomTypeName != "" && *req.RoomTypeName != roomType.RoomTypeName { // 如果请求中提供了房型名称（指针非空且值非空）且与当前房型名称不同
		var existRoomType hotel_admin.RoomTypeDict                                                                                       // 声明房型实体变量，用于存储查询到的已存在房型信息
		result := db.MysqlDB.Where("room_type_name = ? AND id != ? AND deleted_at IS NULL", *req.RoomTypeName, id).First(&existRoomType) // 查询是否已存在相同房型名称的其他房型（排除当前房型和已删除记录）
		if result.Error == nil {                                                                                                         // 如果查询成功（未报错），说明该房型名称已被其他房型使用
			return errors.New("房型名称已存在") // 返回错误信息，表示房型名称已存在
		}
		roomType.RoomTypeName = *req.RoomTypeName // 更新房型名称（解引用指针获取值）
	}

	// 更新字段
	if req.BedSpec != "" { // 如果请求中提供了床型规格（非空），则更新床型规格
		roomType.BedSpec = req.BedSpec // 更新床型规格
	}
	if req.Area != nil { // 如果请求中提供了面积（指针非空），则更新面积
		roomType.Area = req.Area // 更新面积（直接使用指针，因为Area本身是指针类型）
	}
	if req.HasBreakfast != nil { // 如果请求中提供了是否含早（指针非空），则更新是否含早
		roomType.HasBreakfast = *req.HasBreakfast // 更新是否含早（解引用指针获取值）
	}
	if req.HasToiletries != nil { // 如果请求中提供了是否提供洗漱用品（指针非空），则更新是否提供洗漱用品
		roomType.HasToiletries = *req.HasToiletries // 更新是否提供洗漱用品（解引用指针获取值）
	}
	if req.DefaultPrice > 0 { // 如果请求中提供了默认价格（大于0），则更新默认价格
		roomType.DefaultPrice = req.DefaultPrice // 更新默认价格
	}
	if req.Status != "" { // 如果请求中提供了状态（非空），则更新状态
		roomType.Status = req.Status // 更新状态
	}

	return db.MysqlDB.Save(&roomType).Error // 保存房型信息到数据库，返回保存操作的结果（成功为nil，失败为error）
}

// GetRoomType 获取房型详情
func (s *RoomTypeService) GetRoomType(id uint64) (*hotel_admin.RoomTypeDict, error) {
	var roomType hotel_admin.RoomTypeDict                         // 声明房型实体变量，用于存储查询到的房型信息
	if err := db.MysqlDB.First(&roomType, id).Error; err != nil { // 通过房型ID查询房型信息，如果查询失败则说明房型不存在
		return nil, errors.New("房型不存在") // 返回nil和错误信息，表示房型不存在
	}
	return &roomType, nil // 返回房型实体指针和无错误
}

// ListRoomTypes 获取房型列表
func (s *RoomTypeService) ListRoomTypes(req *ListRoomTypeReq) ([]hotel_admin.RoomTypeDict, int64, error) {
	var roomTypes []hotel_admin.RoomTypeDict // 声明房型列表变量，用于存储查询到的房型信息列表
	var total int64                          // 声明总数变量，用于存储符合条件的房型总数

	query := db.MysqlDB.Model(&hotel_admin.RoomTypeDict{}) // 创建房型模型的查询构建器

	// 按状态筛选
	if req.Status != "" { // 如果请求中提供了状态（非空），则添加状态筛选条件
		query = query.Where("status = ?", req.Status) // 添加房型状态筛选条件，只查询指定状态的房型
	}

	// 按名称搜索
	if req.Keyword != "" { // 如果请求中提供了关键词（非空），则添加关键词搜索条件
		query = query.Where("room_type_name LIKE ?", "%"+req.Keyword+"%") // 添加模糊搜索条件，搜索房型名称包含关键词的房型
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil { // 统计符合条件的房型总数，如果统计失败则返回错误
		return nil, 0, err // 返回nil列表、0总数和错误信息
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize                                                                          // 计算分页偏移量（跳过前面的记录数），用于SQL查询的OFFSET子句
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&roomTypes).Error; err != nil { // 添加分页限制（偏移量、每页数量）、排序（按创建时间倒序）并查询房型列表，如果查询失败则返回错误
		return nil, 0, err // 返回nil列表、0总数和错误信息
	}

	return roomTypes, total, nil // 返回房型列表、总数和无错误
}

// DeleteRoomType 删除房型字典（软删除）
func (s *RoomTypeService) DeleteRoomType(id uint64) error {
	var roomType hotel_admin.RoomTypeDict                         // 声明房型实体变量，用于存储查询到的房型信息
	if err := db.MysqlDB.First(&roomType, id).Error; err != nil { // 通过房型ID查询房型信息，如果查询失败则说明房型不存在
		return errors.New("房型不存在") // 返回错误信息，表示房型不存在
	}

	// 检查是否有房源使用此房型
	var count int64                                                                                              // 声明计数变量，用于存储使用该房型的房源数量
	db.MysqlDB.Model(&hotel_admin.RoomInfo{}).Where("room_type_id = ? AND deleted_at IS NULL", id).Count(&count) // 统计使用该房型的房源数量（排除已删除房源）
	if count > 0 {                                                                                               // 如果存在使用该房型的房源（数量大于0），则不允许删除
		return errors.New("该房型正在被使用，无法删除") // 返回错误信息，表示房型正在被使用
	}

	return db.MysqlDB.Delete(&roomType).Error // 执行软删除操作（设置deleted_at字段），根据房型ID删除房型记录，返回删除操作的结果（成功为nil，失败为error）
}

// CreateRoomTypeReq 创建房型请求
type CreateRoomTypeReq struct {
	RoomTypeName  string   `json:"room_type_name" binding:"required"` // 房型名称
	BedSpec       string   `json:"bed_spec" binding:"required"`       // 床型规格
	Area          *float64 `json:"area,omitempty"`                    // 面积
	HasBreakfast  bool     `json:"has_breakfast"`                     // 是否含早
	HasToiletries bool     `json:"has_toiletries"`                    // 是否提供洗漱用品
	DefaultPrice  float64  `json:"default_price" binding:"required"`  // 默认门市价
}

// UpdateRoomTypeReq 更新房型请求
type UpdateRoomTypeReq struct {
	RoomTypeName  *string  `json:"room_type_name,omitempty"` // 房型名称
	BedSpec       string   `json:"bed_spec,omitempty"`       // 床型规格
	Area          *float64 `json:"area,omitempty"`           // 面积
	HasBreakfast  *bool    `json:"has_breakfast,omitempty"`  // 是否含早
	HasToiletries *bool    `json:"has_toiletries,omitempty"` // 是否提供洗漱用品
	DefaultPrice  float64  `json:"default_price"`            // 默认门市价
	Status        string   `json:"status,omitempty"`         // 状态
}

// ListRoomTypeReq 房型列表请求
type ListRoomTypeReq struct {
	Page     int    `json:"page" form:"page"`           // 页码
	PageSize int    `json:"page_size" form:"page_size"` // 每页数量
	Status   string `json:"status" form:"status"`       // 状态筛选
	Keyword  string `json:"keyword" form:"keyword"`     // 关键词搜索
}
