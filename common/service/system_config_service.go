package service

import (
	"errors"
	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"
)

// SystemConfigService 系统配置管理服务
// 负责处理系统配置的创建、更新、查询等核心业务逻辑，
// 包括配置键唯一性检查、配置分类管理、配置状态管理等。
type SystemConfigService struct{}

// CreateSystemConfigReq 创建系统配置请求
type CreateSystemConfigReq struct {
	ConfigCategory string  `json:"config_category" binding:"required"`
	ConfigKey      string  `json:"config_key" binding:"required"`
	ConfigValue    string  `json:"config_value" binding:"required"`
	Description    *string `json:"description,omitempty"`
	Status         string  `json:"status"`
	UpdatedBy      uint64  `json:"updated_by" binding:"required"`
}

// UpdateSystemConfigReq 更新系统配置请求
type UpdateSystemConfigReq struct {
	ID             uint64  `json:"id" binding:"required"`
	ConfigCategory *string `json:"config_category,omitempty"`
	ConfigKey      *string `json:"config_key,omitempty"`
	ConfigValue    *string `json:"config_value,omitempty"`
	Description    *string `json:"description,omitempty"`
	Status         *string `json:"status,omitempty"`
	UpdatedBy      uint64  `json:"updated_by" binding:"required"`
}

// ListSystemConfigsReq 系统配置列表查询请求
type ListSystemConfigsReq struct {
	Page           int     `json:"page"`
	PageSize       int     `json:"page_size"`
	ConfigCategory *string `json:"config_category,omitempty"`
	Status         *string `json:"status,omitempty"`
	Keyword        *string `json:"keyword,omitempty"`
}

// SystemConfigInfo 系统配置信息
type SystemConfigInfo struct {
	ID             uint64  `json:"id"`
	ConfigCategory string  `json:"config_category"`
	ConfigKey      string  `json:"config_key"`
	ConfigValue    string  `json:"config_value"`
	Description    *string `json:"description,omitempty"`
	Status         string  `json:"status"`
	UpdatedAt      string  `json:"updated_at"`
	UpdatedBy      uint64  `json:"updated_by"`
}

// ListSystemConfigsResp 系统配置列表响应
type ListSystemConfigsResp struct {
	List     []SystemConfigInfo `json:"list"`
	Total    uint64             `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

// CreateSystemConfig 创建系统配置
// 业务功能：创建新的系统配置项，用于系统参数的动态配置管理
// 入参说明：
//   - req: 创建系统配置请求，包含配置分类、配置键、配置值、描述、状态（可选，默认ACTIVE）、更新人ID
//
// 返回值说明：
//   - error: 配置键已存在（配置键必须唯一）、业务校验失败或数据库操作错误
func (s *SystemConfigService) CreateSystemConfig(req CreateSystemConfigReq) error {
	// 业务规则：配置键必须唯一，检查该配置键是否已被使用
	var existingConfig hotel_admin.SystemConfig                                                            // 声明系统配置实体变量，用于存储查询到的已存在配置信息
	if err := db.MysqlDB.Where("config_key = ?", req.ConfigKey).First(&existingConfig).Error; err == nil { // 通过配置键查询是否已存在配置记录，如果查询失败则说明该配置键未被使用
		return errors.New("配置键已存在") // 返回错误信息，表示配置键已存在
	}

	config := hotel_admin.SystemConfig{ // 创建系统配置实体对象
		ConfigCategory: req.ConfigCategory, // 设置配置分类（从请求中获取）
		ConfigKey:      req.ConfigKey,      // 设置配置键（从请求中获取）
		ConfigValue:    req.ConfigValue,    // 设置配置值（从请求中获取）
		Description:    req.Description,    // 设置描述（从请求中获取，可为空）
		Status:         "ACTIVE",           // 设置状态为启用（默认值）
		UpdatedBy:      req.UpdatedBy,      // 设置更新人ID（从请求中获取）
	}

	if req.Status != "" { // 如果请求中提供了状态（非空），则使用请求的状态
		config.Status = req.Status // 更新配置状态为请求中提供的状态值
	}

	return db.MysqlDB.Create(&config).Error // 将系统配置信息保存到数据库，返回保存操作的结果（成功为nil，失败为error）
}

// UpdateSystemConfig 更新系统配置
func (s *SystemConfigService) UpdateSystemConfig(req UpdateSystemConfigReq) error {
	var config hotel_admin.SystemConfig                             // 声明系统配置实体变量，用于存储查询到的配置信息
	if err := db.MysqlDB.First(&config, req.ID).Error; err != nil { // 通过配置ID查询配置信息，如果查询失败则说明配置不存在
		return errors.New("系统配置不存在") // 返回错误信息，表示系统配置不存在
	}

	if req.ConfigCategory != nil { // 如果请求中提供了配置分类（指针非空），则更新配置分类
		config.ConfigCategory = *req.ConfigCategory // 更新配置分类（解引用指针获取值）
	}
	if req.ConfigKey != nil { // 如果请求中提供了配置键（指针非空），则需要重新校验唯一性
		// 检查新配置键是否已被其他配置使用
		var existingConfig hotel_admin.SystemConfig                                                                                 // 声明系统配置实体变量，用于存储查询到的已存在配置信息
		if err := db.MysqlDB.Where("config_key = ? AND id != ?", *req.ConfigKey, req.ID).First(&existingConfig).Error; err == nil { // 通过配置键查询是否已存在其他配置使用该配置键（排除当前配置），如果查询失败则说明该配置键未被使用
			return errors.New("配置键已被使用") // 返回错误信息，表示配置键已被使用
		}
		config.ConfigKey = *req.ConfigKey // 更新配置键（解引用指针获取值）
	}
	if req.ConfigValue != nil { // 如果请求中提供了配置值（指针非空），则更新配置值
		config.ConfigValue = *req.ConfigValue // 更新配置值（解引用指针获取值）
	}
	if req.Description != nil { // 如果请求中提供了描述（指针非空），则更新描述
		config.Description = req.Description // 更新描述（直接使用指针，因为Description本身是指针类型）
	}
	if req.Status != nil { // 如果请求中提供了状态（指针非空），则更新状态
		config.Status = *req.Status // 更新状态（解引用指针获取值）
	}
	config.UpdatedBy = req.UpdatedBy // 更新更新人ID（从请求中获取）

	return db.MysqlDB.Save(&config).Error // 保存系统配置信息到数据库，返回保存操作的结果（成功为nil，失败为error）
}

// GetSystemConfig 获取系统配置详情
func (s *SystemConfigService) GetSystemConfig(id uint64) (*SystemConfigInfo, error) {
	var config hotel_admin.SystemConfig                         // 声明系统配置实体变量，用于存储查询到的配置信息
	if err := db.MysqlDB.First(&config, id).Error; err != nil { // 通过配置ID查询配置信息，如果查询失败则说明配置不存在
		return nil, errors.New("系统配置不存在") // 返回nil和错误信息，表示系统配置不存在
	}

	return &SystemConfigInfo{ // 返回系统配置信息对象指针
		ID:             config.ID,                                      // 设置配置ID（从系统配置实体中获取）
		ConfigCategory: config.ConfigCategory,                          // 设置配置分类（从系统配置实体中获取）
		ConfigKey:      config.ConfigKey,                               // 设置配置键（从系统配置实体中获取）
		ConfigValue:    config.ConfigValue,                             // 设置配置值（从系统配置实体中获取）
		Description:    config.Description,                             // 设置描述（从系统配置实体中获取，可为空）
		Status:         config.Status,                                  // 设置状态（从系统配置实体中获取）
		UpdatedAt:      config.UpdatedAt.Format("2006-01-02 15:04:05"), // 设置更新时间（格式化时间字符串）
		UpdatedBy:      config.UpdatedBy,                               // 设置更新人ID（从系统配置实体中获取）
	}, nil // 返回系统配置信息指针和无错误
}

// ListSystemConfigs 获取系统配置列表
func (s *SystemConfigService) ListSystemConfigs(req ListSystemConfigsReq) (*ListSystemConfigsResp, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	offset := (req.Page - 1) * req.PageSize

	query := db.MysqlDB.Model(&hotel_admin.SystemConfig{}).Where("deleted_at IS NULL")

	if req.ConfigCategory != nil && *req.ConfigCategory != "" {
		query = query.Where("config_category = ?", *req.ConfigCategory)
	}
	if req.Status != nil && *req.Status != "" {
		query = query.Where("status = ?", *req.Status)
	}
	if req.Keyword != nil && *req.Keyword != "" {
		keyword := "%" + *req.Keyword + "%"
		query = query.Where("config_key LIKE ? OR config_value LIKE ? OR description LIKE ?", keyword, keyword, keyword)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var configs []hotel_admin.SystemConfig
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&configs).Error; err != nil {
		return nil, err
	}

	configInfos := make([]SystemConfigInfo, len(configs))
	for i, config := range configs {
		configInfos[i] = SystemConfigInfo{
			ID:             config.ID,
			ConfigCategory: config.ConfigCategory,
			ConfigKey:      config.ConfigKey,
			ConfigValue:    config.ConfigValue,
			Description:    config.Description,
			Status:         config.Status,
			UpdatedAt:      config.UpdatedAt.Format("2006-01-02 15:04:05"),
			UpdatedBy:      config.UpdatedBy,
		}
	}

	return &ListSystemConfigsResp{
		List:     configInfos,
		Total:    uint64(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// DeleteSystemConfig 删除系统配置（软删除）
func (s *SystemConfigService) DeleteSystemConfig(id uint64) error {
	return db.MysqlDB.Delete(&hotel_admin.SystemConfig{}, id).Error // 执行软删除操作（设置deleted_at字段），根据配置ID删除系统配置记录，返回删除操作的结果（成功为nil，失败为error）
}

// GetConfigByCategory 按分类获取配置
func (s *SystemConfigService) GetConfigByCategory(category string) ([]SystemConfigInfo, error) {
	var configs []hotel_admin.SystemConfig
	if err := db.MysqlDB.Where("config_category = ? AND status = ? AND deleted_at IS NULL", category, "ACTIVE").
		Order("id ASC").Find(&configs).Error; err != nil {
		return nil, err
	}

	configInfos := make([]SystemConfigInfo, len(configs))
	for i, config := range configs {
		configInfos[i] = SystemConfigInfo{
			ID:             config.ID,
			ConfigCategory: config.ConfigCategory,
			ConfigKey:      config.ConfigKey,
			ConfigValue:    config.ConfigValue,
			Description:    config.Description,
			Status:         config.Status,
			UpdatedAt:      config.UpdatedAt.Format("2006-01-02 15:04:05"),
			UpdatedBy:      config.UpdatedBy,
		}
	}

	return configInfos, nil
}
