package service

import (
	"errors"
	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"
)

// FacilityService 设施字典管理服务
// 负责处理酒店设施字典的创建、更新、查询、删除等核心业务逻辑，
// 包括设施名称唯一性检查、设施与房源的关联关系管理等。
type FacilityService struct{}

// CreateFacility 创建设施字典
func (s *FacilityService) CreateFacility(req *CreateFacilityReq) (*hotel_admin.FacilityDict, error) {
	// 检查设施名称是否已存在
	var existFacility hotel_admin.FacilityDict                                                                     // 声明设施实体变量，用于存储查询到的已存在设施信息
	result := db.MysqlDB.Where("facility_name = ? AND deleted_at IS NULL", req.FacilityName).First(&existFacility) // 查询是否已存在相同的设施名称（排除已删除记录）
	if result.Error == nil {                                                                                       // 如果查询成功（未报错），说明该设施名称已存在
		return nil, errors.New("设施名称已存在") // 返回nil和错误信息，表示设施名称已存在
	}

	facility := &hotel_admin.FacilityDict{ // 创建设施实体对象指针
		FacilityName: req.FacilityName, // 设置设施名称（从请求中获取）
		Description:  req.Description,  // 设置设施描述（从请求中获取，可为空）
		Status:       "ACTIVE",         // 设置设施状态为启用（默认值）
	}

	if err := db.MysqlDB.Create(facility).Error; err != nil { // 将设施信息保存到数据库，如果保存失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	return facility, nil // 返回设施实体指针和无错误
}

// UpdateFacility 更新设施字典
func (s *FacilityService) UpdateFacility(id uint64, req *UpdateFacilityReq) error {
	var facility hotel_admin.FacilityDict                         // 声明设施实体变量，用于存储查询到的设施信息
	if err := db.MysqlDB.First(&facility, id).Error; err != nil { // 通过设施ID查询设施信息，如果查询失败则说明设施不存在
		return errors.New("设施不存在") // 返回错误信息，表示设施不存在
	}

	// 如果修改了设施名称，检查是否重复
	if req.FacilityName != "" && req.FacilityName != facility.FacilityName { // 如果请求中提供了设施名称（非空）且与当前设施名称不同
		var existFacility hotel_admin.FacilityDict                                                                                     // 声明设施实体变量，用于存储查询到的已存在设施信息
		result := db.MysqlDB.Where("facility_name = ? AND id != ? AND deleted_at IS NULL", req.FacilityName, id).First(&existFacility) // 查询是否已存在相同设施名称的其他设施（排除当前设施和已删除记录）
		if result.Error == nil {                                                                                                       // 如果查询成功（未报错），说明该设施名称已被其他设施使用
			return errors.New("设施名称已存在") // 返回错误信息，表示设施名称已存在
		}
		facility.FacilityName = req.FacilityName // 更新设施名称
	}

	if req.Description != nil { // 如果请求中提供了描述（指针非空），则更新描述
		facility.Description = req.Description // 更新描述（直接使用指针，因为Description本身是指针类型）
	}
	if req.Status != "" { // 如果请求中提供了状态（非空），则更新状态
		facility.Status = req.Status // 更新状态
	}

	return db.MysqlDB.Save(&facility).Error // 保存设施信息到数据库，返回保存操作的结果（成功为nil，失败为error）
}

// GetFacility 获取设施详情
func (s *FacilityService) GetFacility(id uint64) (*hotel_admin.FacilityDict, error) {
	var facility hotel_admin.FacilityDict                         // 声明设施实体变量，用于存储查询到的设施信息
	if err := db.MysqlDB.First(&facility, id).Error; err != nil { // 通过设施ID查询设施信息，如果查询失败则说明设施不存在
		return nil, errors.New("设施不存在") // 返回nil和错误信息，表示设施不存在
	}
	return &facility, nil // 返回设施实体指针和无错误
}

// ListFacilities 获取设施列表
func (s *FacilityService) ListFacilities(req *ListFacilityReq) ([]hotel_admin.FacilityDict, int64, error) {
	var facilities []hotel_admin.FacilityDict // 声明设施列表变量，用于存储查询到的设施信息列表
	var total int64                           // 声明总数变量，用于存储符合条件的设施总数

	query := db.MysqlDB.Model(&hotel_admin.FacilityDict{}) // 创建设施模型的查询构建器

	// 按状态筛选
	if req.Status != "" { // 如果请求中提供了状态（非空），则添加状态筛选条件
		query = query.Where("status = ?", req.Status) // 添加设施状态筛选条件，只查询指定状态的设施
	}

	// 按名称搜索
	if req.Keyword != "" { // 如果请求中提供了关键词（非空），则添加关键词搜索条件
		query = query.Where("facility_name LIKE ?", "%"+req.Keyword+"%") // 添加模糊搜索条件，搜索设施名称包含关键词的设施
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil { // 统计符合条件的设施总数，如果统计失败则返回错误
		return nil, 0, err // 返回nil列表、0总数和错误信息
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize                                                                           // 计算分页偏移量（跳过前面的记录数），用于SQL查询的OFFSET子句
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&facilities).Error; err != nil { // 添加分页限制（偏移量、每页数量）、排序（按创建时间倒序）并查询设施列表，如果查询失败则返回错误
		return nil, 0, err // 返回nil列表、0总数和错误信息
	}

	return facilities, total, nil // 返回设施列表、总数和无错误
}

// DeleteFacility 删除设施字典（软删除）
func (s *FacilityService) DeleteFacility(id uint64) error {
	var facility hotel_admin.FacilityDict                         // 声明设施实体变量，用于存储查询到的设施信息
	if err := db.MysqlDB.First(&facility, id).Error; err != nil { // 通过设施ID查询设施信息，如果查询失败则说明设施不存在
		return errors.New("设施不存在") // 返回错误信息，表示设施不存在
	}

	// 检查是否有房源使用此设施
	var count int64                                                                                                         // 声明计数变量，用于存储使用该设施的房源数量
	db.MysqlDB.Model(&hotel_admin.RoomFacilityRelation{}).Where("facility_id = ? AND deleted_at IS NULL", id).Count(&count) // 统计使用该设施的房源数量（通过房源设施关联表，排除已删除关联）
	if count > 0 {                                                                                                          // 如果存在使用该设施的房源（数量大于0），则不允许删除
		return errors.New("该设施正在被房源使用，无法删除") // 返回错误信息，表示设施正在被使用
	}

	return db.MysqlDB.Delete(&facility).Error // 执行软删除操作（设置deleted_at字段），根据设施ID删除设施记录，返回删除操作的结果（成功为nil，失败为error）
}

// CreateFacilityReq 创建设施请求
type CreateFacilityReq struct {
	FacilityName string  `json:"facility_name" binding:"required"` // 设施名称
	Description  *string `json:"description,omitempty"`            // 设施描述
}

// UpdateFacilityReq 更新设施请求
type UpdateFacilityReq struct {
	FacilityName string  `json:"facility_name,omitempty"` // 设施名称
	Description  *string `json:"description,omitempty"`   // 设施描述
	Status       string  `json:"status,omitempty"`        // 状态
}

// ListFacilityReq 设施列表请求
type ListFacilityReq struct {
	Page     int    `json:"page" form:"page"`           // 页码
	PageSize int    `json:"page_size" form:"page_size"` // 每页数量
	Status   string `json:"status" form:"status"`       // 状态筛选
	Keyword  string `json:"keyword" form:"keyword"`     // 关键词搜索
}
