package service

import (
	"errors"
	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"
)

// CancellationPolicyService 退订政策管理服务
// 负责处理酒店退订政策的创建、更新、查询、删除等核心业务逻辑，
// 包括退订政策与房型的关联校验、退订规则和违约金比例管理等。
type CancellationPolicyService struct{}

// CreateCancellationPolicy 创建退订政策
func (s *CancellationPolicyService) CreateCancellationPolicy(req *CreateCancellationPolicyReq) (*hotel_admin.CancellationPolicy, error) {
	// 如果指定了房型，检查房型是否存在
	if req.RoomTypeID != nil && *req.RoomTypeID > 0 { // 如果请求中提供了房型ID（指针非空且值大于0）
		var roomType hotel_admin.RoomTypeDict                                      // 声明房型实体变量，用于存储查询到的房型信息
		if err := db.MysqlDB.First(&roomType, *req.RoomTypeID).Error; err != nil { // 通过房型ID查询房型信息（解引用指针获取值），如果查询失败则说明房型不存在
			return nil, errors.New("房型不存在") // 返回nil和错误信息，表示房型不存在
		}
	}

	policy := &hotel_admin.CancellationPolicy{ // 创建退订政策实体对象指针
		PolicyName:      req.PolicyName,      // 设置政策名称（从请求中获取）
		RuleDescription: req.RuleDescription, // 设置规则描述（从请求中获取）
		PenaltyRatio:    req.PenaltyRatio,    // 设置违约金比例（从请求中获取）
		RoomTypeID:      req.RoomTypeID,      // 设置适用房型ID（从请求中获取，可为空）
		Status:          "ACTIVE",            // 设置政策状态为启用（默认值）
	}

	if err := db.MysqlDB.Create(policy).Error; err != nil { // 将退订政策信息保存到数据库，如果保存失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	return policy, nil // 返回退订政策实体指针和无错误
}

// UpdateCancellationPolicy 更新退订政策
func (s *CancellationPolicyService) UpdateCancellationPolicy(id uint64, req *UpdateCancellationPolicyReq) error {
	var policy hotel_admin.CancellationPolicy                   // 声明退订政策实体变量，用于存储查询到的政策信息
	if err := db.MysqlDB.First(&policy, id).Error; err != nil { // 通过政策ID查询政策信息，如果查询失败则说明政策不存在
		return errors.New("退订政策不存在") // 返回错误信息，表示退订政策不存在
	}

	if req.PolicyName != "" { // 如果请求中提供了政策名称（非空），则更新政策名称
		policy.PolicyName = req.PolicyName // 更新政策名称
	}
	if req.RuleDescription != "" { // 如果请求中提供了规则描述（非空），则更新规则描述
		policy.RuleDescription = req.RuleDescription // 更新规则描述
	}
	if req.PenaltyRatio > 0 { // 如果请求中提供了违约金比例（大于0），则更新违约金比例
		policy.PenaltyRatio = req.PenaltyRatio // 更新违约金比例
	}
	if req.RoomTypeID != nil { // 如果请求中提供了房型ID（指针非空），则需要验证房型是否存在
		// 如果指定了房型，检查房型是否存在
		if *req.RoomTypeID > 0 { // 如果房型ID大于0（解引用指针获取值）
			var roomType hotel_admin.RoomTypeDict                                      // 声明房型实体变量，用于存储查询到的房型信息
			if err := db.MysqlDB.First(&roomType, *req.RoomTypeID).Error; err != nil { // 通过房型ID查询房型信息（解引用指针获取值），如果查询失败则说明房型不存在
				return errors.New("房型不存在") // 返回错误信息，表示房型不存在
			}
		}
		policy.RoomTypeID = req.RoomTypeID // 更新适用房型ID（直接使用指针，因为RoomTypeID本身是指针类型）
	}
	if req.Status != "" { // 如果请求中提供了状态（非空），则更新状态
		policy.Status = req.Status // 更新状态
	}

	return db.MysqlDB.Save(&policy).Error // 保存退订政策信息到数据库，返回保存操作的结果（成功为nil，失败为error）
}

// GetCancellationPolicy 获取退订政策详情
func (s *CancellationPolicyService) GetCancellationPolicy(id uint64) (*hotel_admin.CancellationPolicy, error) {
	var policy hotel_admin.CancellationPolicy                                       // 声明退订政策实体变量，用于存储查询到的政策信息
	if err := db.MysqlDB.Preload("RoomType").First(&policy, id).Error; err != nil { // 通过政策ID查询政策信息（预加载房型信息），如果查询失败则说明政策不存在
		return nil, errors.New("退订政策不存在") // 返回nil和错误信息，表示退订政策不存在
	}
	return &policy, nil // 返回退订政策实体指针和无错误
}

// ListCancellationPolicies 获取退订政策列表
func (s *CancellationPolicyService) ListCancellationPolicies(req *ListCancellationPolicyReq) ([]hotel_admin.CancellationPolicy, int64, error) {
	var policies []hotel_admin.CancellationPolicy // 声明退订政策列表变量，用于存储查询到的政策信息列表
	var total int64                               // 声明总数变量，用于存储符合条件的政策总数

	query := db.MysqlDB.Model(&hotel_admin.CancellationPolicy{}) // 创建退订政策模型的查询构建器

	// 按房型筛选
	if req.RoomTypeID > 0 { // 如果请求中提供了房型ID（大于0），则添加房型筛选条件
		query = query.Where("room_type_id = ?", req.RoomTypeID) // 添加房型ID筛选条件，只查询指定房型的政策
	}

	// 按状态筛选
	if req.Status != "" { // 如果请求中提供了状态（非空），则添加状态筛选条件
		query = query.Where("status = ?", req.Status) // 添加政策状态筛选条件，只查询指定状态的政策
	}

	// 按名称搜索
	if req.Keyword != "" { // 如果请求中提供了关键词（非空），则添加关键词搜索条件
		query = query.Where("policy_name LIKE ? OR rule_description LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%") // 添加模糊搜索条件，搜索政策名称或规则描述包含关键词的政策
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil { // 统计符合条件的政策总数，如果统计失败则返回错误
		return nil, 0, err // 返回nil列表、0总数和错误信息
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize // 计算分页偏移量（跳过前面的记录数），用于SQL查询的OFFSET子句
	if err := query.Preload("RoomType").    // 预加载房型信息关联数据（JOIN查询房型信息）
						Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&policies).Error; err != nil { // 添加分页限制（偏移量、每页数量）、排序（按创建时间倒序）并查询政策列表，如果查询失败则返回错误
		return nil, 0, err // 返回nil列表、0总数和错误信息
	}

	return policies, total, nil // 返回政策列表、总数和无错误
}

// DeleteCancellationPolicy 删除退订政策（软删除）
func (s *CancellationPolicyService) DeleteCancellationPolicy(id uint64) error {
	var policy hotel_admin.CancellationPolicy                   // 声明退订政策实体变量，用于存储查询到的政策信息
	if err := db.MysqlDB.First(&policy, id).Error; err != nil { // 通过政策ID查询政策信息，如果查询失败则说明政策不存在
		return errors.New("退订政策不存在") // 返回错误信息，表示退订政策不存在
	}

	// 检查是否有房源或订单使用此政策
	var roomCount, orderCount int64                                                                                              // 声明计数变量，用于存储使用该政策的房源数量和订单数量
	db.MysqlDB.Model(&hotel_admin.RoomInfo{}).Where("cancellation_policy_id = ? AND deleted_at IS NULL", id).Count(&roomCount)   // 统计使用该政策的房源数量（排除已删除房源）
	db.MysqlDB.Model(&hotel_admin.OrderMain{}).Where("cancellation_policy_id = ? AND deleted_at IS NULL", id).Count(&orderCount) // 统计使用该政策的订单数量（排除已删除订单）
	if roomCount > 0 || orderCount > 0 {                                                                                         // 如果存在使用该政策的房源或订单（数量大于0），则不允许删除
		return errors.New("该退订政策正在被房源或订单使用，无法删除") // 返回错误信息，表示政策正在被使用
	}

	return db.MysqlDB.Delete(&policy).Error // 执行软删除操作（设置deleted_at字段），根据政策ID删除退订政策记录，返回删除操作的结果（成功为nil，失败为error）
}

// CreateCancellationPolicyReq 创建退订政策请求
type CreateCancellationPolicyReq struct {
	PolicyName      string  `json:"policy_name" binding:"required"`      // 政策名称
	RuleDescription string  `json:"rule_description" binding:"required"` // 规则描述
	PenaltyRatio    float64 `json:"penalty_ratio" binding:"required"`    // 违约金比例
	RoomTypeID      *uint64 `json:"room_type_id,omitempty"`              // 适用房型ID（可选）
}

// UpdateCancellationPolicyReq 更新退订政策请求
type UpdateCancellationPolicyReq struct {
	PolicyName      string  `json:"policy_name,omitempty"`      // 政策名称
	RuleDescription string  `json:"rule_description,omitempty"` // 规则描述
	PenaltyRatio    float64 `json:"penalty_ratio"`              // 违约金比例
	RoomTypeID      *uint64 `json:"room_type_id,omitempty"`     // 适用房型ID
	Status          string  `json:"status,omitempty"`           // 状态
}

// ListCancellationPolicyReq 退订政策列表请求
type ListCancellationPolicyReq struct {
	Page       int    `json:"page" form:"page"`                 // 页码
	PageSize   int    `json:"page_size" form:"page_size"`       // 每页数量
	RoomTypeID uint64 `json:"room_type_id" form:"room_type_id"` // 房型ID筛选
	Status     string `json:"status" form:"status"`             // 状态筛选
	Keyword    string `json:"keyword" form:"keyword"`           // 关键词搜索
}
