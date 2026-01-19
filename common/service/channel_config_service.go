package service

import (
	"errors"
	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"
)

// ChannelConfigService 渠道配置管理服务
// 负责处理外部渠道（如OTA平台）配置的创建、更新、查询、删除等核心业务逻辑，
// 包括渠道API URL配置、同步规则管理、渠道状态管理等。
type ChannelConfigService struct{}

// CreateChannelConfigReq 创建渠道配置请求
type CreateChannelConfigReq struct {
	ChannelName string `json:"channel_name" binding:"required"`
	ChannelCode string `json:"channel_code" binding:"required"`
	ApiURL      string `json:"api_url" binding:"required"`
	SyncRule    string `json:"sync_rule"`
	Status      string `json:"status"`
}

// UpdateChannelConfigReq 更新渠道配置请求
type UpdateChannelConfigReq struct {
	ID          uint64  `json:"id" binding:"required"`
	ChannelName *string `json:"channel_name,omitempty"`
	ChannelCode *string `json:"channel_code,omitempty"`
	ApiURL      *string `json:"api_url,omitempty"`
	SyncRule    *string `json:"sync_rule,omitempty"`
	Status      *string `json:"status,omitempty"`
}

// ListChannelConfigsReq 渠道配置列表查询请求
type ListChannelConfigsReq struct {
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
	Status   *string `json:"status,omitempty"`
	Keyword  *string `json:"keyword,omitempty"`
}

// ChannelConfigInfo 渠道配置信息
type ChannelConfigInfo struct {
	ID          uint64 `json:"id"`
	ChannelName string `json:"channel_name"`
	ChannelCode string `json:"channel_code"`
	ApiURL      string `json:"api_url"`
	SyncRule    string `json:"sync_rule"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ListChannelConfigsResp 渠道配置列表响应
type ListChannelConfigsResp struct {
	List     []ChannelConfigInfo `json:"list"`
	Total    uint64              `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
}

// CreateChannelConfig 创建渠道配置
func (s *ChannelConfigService) CreateChannelConfig(req CreateChannelConfigReq) error {
	// 检查渠道名称是否已存在
	var existingConfig hotel_admin.ChannelConfig                                                               // 声明渠道配置实体变量，用于存储查询到的已存在配置信息
	if err := db.MysqlDB.Where("channel_name = ?", req.ChannelName).First(&existingConfig).Error; err == nil { // 通过渠道名称查询是否已存在配置记录，如果查询失败则说明该渠道名称未被使用
		return errors.New("渠道名称已存在") // 返回错误信息，表示渠道名称已存在
	}

	// 检查渠道编码是否已存在
	if err := db.MysqlDB.Where("channel_code = ?", req.ChannelCode).First(&existingConfig).Error; err == nil { // 通过渠道编码查询是否已存在配置记录，如果查询失败则说明该渠道编码未被使用
		return errors.New("渠道编码已存在") // 返回错误信息，表示渠道编码已存在
	}

	config := hotel_admin.ChannelConfig{ // 创建渠道配置实体对象
		ChannelName: req.ChannelName, // 设置渠道名称（从请求中获取）
		ChannelCode: req.ChannelCode, // 设置渠道编码（从请求中获取）
		ApiURL:      req.ApiURL,      // 设置API URL（从请求中获取）
		SyncRule:    "REALTIME",      // 设置同步规则为实时同步（默认值）
		Status:      "ACTIVE",        // 设置状态为启用（默认值）
	}

	if req.SyncRule != "" { // 如果请求中提供了同步规则（非空），则使用请求的同步规则
		config.SyncRule = req.SyncRule // 更新同步规则为请求中提供的同步规则值
	}
	if req.Status != "" { // 如果请求中提供了状态（非空），则使用请求的状态
		config.Status = req.Status // 更新配置状态为请求中提供的状态值
	}

	return db.MysqlDB.Create(&config).Error // 将渠道配置信息保存到数据库，返回保存操作的结果（成功为nil，失败为error）
}

// UpdateChannelConfig 更新渠道配置
func (s *ChannelConfigService) UpdateChannelConfig(req UpdateChannelConfigReq) error {
	var config hotel_admin.ChannelConfig                            // 声明渠道配置实体变量，用于存储查询到的配置信息
	if err := db.MysqlDB.First(&config, req.ID).Error; err != nil { // 通过配置ID查询配置信息，如果查询失败则说明配置不存在
		return errors.New("渠道配置不存在") // 返回错误信息，表示渠道配置不存在
	}

	if req.ChannelName != nil { // 如果请求中提供了渠道名称（指针非空），则需要重新校验唯一性
		// 检查新渠道名称是否已被其他配置使用
		var existingConfig hotel_admin.ChannelConfig                                                                                    // 声明渠道配置实体变量，用于存储查询到的已存在配置信息
		if err := db.MysqlDB.Where("channel_name = ? AND id != ?", *req.ChannelName, req.ID).First(&existingConfig).Error; err == nil { // 通过渠道名称查询是否已存在其他配置使用该渠道名称（排除当前配置），如果查询失败则说明该渠道名称未被使用
			return errors.New("渠道名称已被使用") // 返回错误信息，表示渠道名称已被使用
		}
		config.ChannelName = *req.ChannelName // 更新渠道名称（解引用指针获取值）
	}
	if req.ChannelCode != nil { // 如果请求中提供了渠道编码（指针非空），则需要重新校验唯一性
		// 检查新渠道编码是否已被其他配置使用
		var existingConfig hotel_admin.ChannelConfig                                                                                    // 声明渠道配置实体变量，用于存储查询到的已存在配置信息
		if err := db.MysqlDB.Where("channel_code = ? AND id != ?", *req.ChannelCode, req.ID).First(&existingConfig).Error; err == nil { // 通过渠道编码查询是否已存在其他配置使用该渠道编码（排除当前配置），如果查询失败则说明该渠道编码未被使用
			return errors.New("渠道编码已被使用") // 返回错误信息，表示渠道编码已被使用
		}
		config.ChannelCode = *req.ChannelCode // 更新渠道编码（解引用指针获取值）
	}
	if req.ApiURL != nil { // 如果请求中提供了API URL（指针非空），则更新API URL
		config.ApiURL = *req.ApiURL // 更新API URL（解引用指针获取值）
	}
	if req.SyncRule != nil { // 如果请求中提供了同步规则（指针非空），则更新同步规则
		config.SyncRule = *req.SyncRule // 更新同步规则（解引用指针获取值）
	}
	if req.Status != nil { // 如果请求中提供了状态（指针非空），则更新状态
		config.Status = *req.Status // 更新状态（解引用指针获取值）
	}

	return db.MysqlDB.Save(&config).Error // 保存渠道配置信息到数据库，返回保存操作的结果（成功为nil，失败为error）
}

// GetChannelConfig 获取渠道配置详情
func (s *ChannelConfigService) GetChannelConfig(id uint64) (*ChannelConfigInfo, error) {
	var config hotel_admin.ChannelConfig
	if err := db.MysqlDB.First(&config, id).Error; err != nil {
		return nil, errors.New("渠道配置不存在")
	}

	return &ChannelConfigInfo{
		ID:          config.ID,
		ChannelName: config.ChannelName,
		ChannelCode: config.ChannelCode,
		ApiURL:      config.ApiURL,
		SyncRule:    config.SyncRule,
		Status:      config.Status,
		CreatedAt:   config.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   config.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// ListChannelConfigs 获取渠道配置列表
func (s *ChannelConfigService) ListChannelConfigs(req ListChannelConfigsReq) (*ListChannelConfigsResp, error) {
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

	query := db.MysqlDB.Model(&hotel_admin.ChannelConfig{}).Where("deleted_at IS NULL")

	if req.Status != nil && *req.Status != "" {
		query = query.Where("status = ?", *req.Status)
	}
	if req.Keyword != nil && *req.Keyword != "" {
		keyword := "%" + *req.Keyword + "%"
		query = query.Where("channel_name LIKE ? OR channel_code LIKE ?", keyword, keyword)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var configs []hotel_admin.ChannelConfig
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&configs).Error; err != nil {
		return nil, err
	}

	configInfos := make([]ChannelConfigInfo, len(configs))
	for i, config := range configs {
		configInfos[i] = ChannelConfigInfo{
			ID:          config.ID,
			ChannelName: config.ChannelName,
			ChannelCode: config.ChannelCode,
			ApiURL:      config.ApiURL,
			SyncRule:    config.SyncRule,
			Status:      config.Status,
			CreatedAt:   config.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   config.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &ListChannelConfigsResp{
		List:     configInfos,
		Total:    uint64(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// DeleteChannelConfig 删除渠道配置（软删除）
func (s *ChannelConfigService) DeleteChannelConfig(id uint64) error {
	return db.MysqlDB.Delete(&hotel_admin.ChannelConfig{}, id).Error // 执行软删除操作（设置deleted_at字段），根据配置ID删除渠道配置记录，返回删除操作的结果（成功为nil，失败为error）
}
