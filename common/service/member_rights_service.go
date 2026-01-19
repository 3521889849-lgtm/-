package service

import (
	"errors"
	"fmt"
	"time"

	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"

	"gorm.io/gorm"
)

const (
	RightsStatusActive   = "ACTIVE"
	RightsStatusInactive = "INACTIVE"
)

// MemberRightsService 会员权益管理服务
// 负责处理会员权益的创建、更新、查询、删除等核心业务逻辑，
// 包括会员权益与会员等级的关联、权益生效时间管理、权益折扣比例管理等。
type MemberRightsService struct{}

// CreateMemberRightsReq 创建会员权益请求
type CreateMemberRightsReq struct {
	MemberLevel   string     `json:"member_level" binding:"required"`
	RightsName    string     `json:"rights_name" binding:"required"`
	Description   *string    `json:"description,omitempty"`
	DiscountRatio *float64   `json:"discount_ratio,omitempty"`
	EffectiveTime time.Time  `json:"effective_time" binding:"required"`
	ExpireTime    *time.Time `json:"expire_time,omitempty"`
	Status        string     `json:"status"`
}

// UpdateMemberRightsReq 更新会员权益请求
type UpdateMemberRightsReq struct {
	ID            uint64     `json:"id" binding:"required"`
	MemberLevel   *string    `json:"member_level,omitempty"`
	RightsName    *string    `json:"rights_name,omitempty"`
	Description   *string    `json:"description,omitempty"`
	DiscountRatio *float64   `json:"discount_ratio,omitempty"`
	EffectiveTime *time.Time `json:"effective_time,omitempty"`
	ExpireTime    *time.Time `json:"expire_time,omitempty"`
	Status        *string    `json:"status,omitempty"`
}

// ListMemberRightsReq 会员权益列表查询请求
type ListMemberRightsReq struct {
	Page        int     `json:"page"`
	PageSize    int     `json:"page_size"`
	MemberLevel *string `json:"member_level,omitempty"`
	Status      *string `json:"status,omitempty"`
	Keyword     *string `json:"keyword,omitempty"`
}

// MemberRightsInfo 会员权益信息
type MemberRightsInfo struct {
	ID            uint64     `json:"id"`
	MemberLevel   string     `json:"member_level"`
	RightsName    string     `json:"rights_name"`
	Description   *string    `json:"description,omitempty"`
	DiscountRatio *float64   `json:"discount_ratio,omitempty"`
	EffectiveTime time.Time  `json:"effective_time"`
	ExpireTime    *time.Time `json:"expire_time,omitempty"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
}

// ListMemberRightsResp 会员权益列表响应
type ListMemberRightsResp struct {
	List     []MemberRightsInfo `json:"list"`
	Total    uint64             `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

// buildMemberRightsInfo 构建会员权益信息（避免代码重复）
func buildMemberRightsInfo(rights *hotel_admin.MemberRights) MemberRightsInfo {
	return MemberRightsInfo{
		ID:            rights.ID,
		MemberLevel:   rights.MemberLevel,
		RightsName:    rights.RightsName,
		Description:   rights.Description,
		DiscountRatio: rights.DiscountRatio,
		EffectiveTime: rights.EffectiveTime,
		ExpireTime:    rights.ExpireTime,
		Status:        rights.Status,
		CreatedAt:     rights.CreatedAt,
	}
}

// CreateMemberRights 创建会员权益
func (s *MemberRightsService) CreateMemberRights(req CreateMemberRightsReq) error {
	status := RightsStatusActive
	if req.Status != "" {
		status = req.Status
	}

	rights := hotel_admin.MemberRights{
		MemberLevel:   req.MemberLevel,
		RightsName:    req.RightsName,
		Description:   req.Description,
		DiscountRatio: req.DiscountRatio,
		EffectiveTime: req.EffectiveTime,
		ExpireTime:    req.ExpireTime,
		Status:        status,
	}

	return db.MysqlDB.Create(&rights).Error
}

// UpdateMemberRights 更新会员权益
func (s *MemberRightsService) UpdateMemberRights(req UpdateMemberRightsReq) error {
	updates := make(map[string]interface{})

	if req.MemberLevel != nil {
		updates["member_level"] = *req.MemberLevel
	}
	if req.RightsName != nil {
		updates["rights_name"] = *req.RightsName
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.DiscountRatio != nil {
		updates["discount_ratio"] = *req.DiscountRatio
	}
	if req.EffectiveTime != nil {
		updates["effective_time"] = *req.EffectiveTime
	}
	if req.ExpireTime != nil {
		updates["expire_time"] = *req.ExpireTime
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		return nil
	}

	result := db.MysqlDB.Model(&hotel_admin.MemberRights{}).Where("id = ?", req.ID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("会员权益ID %d 不存在", req.ID)
	}

	return nil
}

// GetMemberRights 获取会员权益详情
func (s *MemberRightsService) GetMemberRights(id uint64) (*MemberRightsInfo, error) {
	var rights hotel_admin.MemberRights
	if err := db.MysqlDB.First(&rights, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("会员权益ID %d 不存在", id)
		}
		return nil, err
	}

	info := buildMemberRightsInfo(&rights)
	return &info, nil
}

// ListMemberRights 获取会员权益列表
func (s *MemberRightsService) ListMemberRights(req ListMemberRightsReq) (*ListMemberRightsResp, error) {
	req.Page = max(req.Page, 1)                   // 如果页码小于1，则设置为1（使用max函数确保最小值）
	req.PageSize = min(max(req.PageSize, 1), 100) // 如果每页数量小于1则设置为1，如果大于100则设置为100（使用min和max函数确保范围）
	offset := (req.Page - 1) * req.PageSize       // 计算分页偏移量（跳过前面的记录数），用于SQL查询的OFFSET子句

	query := db.MysqlDB.Model(&hotel_admin.MemberRights{}).Where("deleted_at IS NULL") // 创建会员权益模型的查询构建器，添加软删除筛选条件（只查询未删除的会员权益）

	if req.MemberLevel != nil && *req.MemberLevel != "" { // 如果请求中提供了会员等级（指针非空且值非空），则添加会员等级筛选条件
		query = query.Where("member_level = ?", *req.MemberLevel) // 添加会员等级筛选条件，只查询指定会员等级的权益（解引用指针获取值）
	}
	if req.Status != nil && *req.Status != "" { // 如果请求中提供了状态（指针非空且值非空），则添加状态筛选条件
		query = query.Where("status = ?", *req.Status) // 添加权益状态筛选条件，只查询指定状态的权益（解引用指针获取值）
	}
	if req.Keyword != nil && *req.Keyword != "" { // 如果请求中提供了关键词（指针非空且值非空），则添加关键词搜索条件
		keyword := "%" + *req.Keyword + "%"                                               // 构建模糊搜索关键词（前后加%通配符）
		query = query.Where("rights_name LIKE ? OR description LIKE ?", keyword, keyword) // 添加关键词搜索条件，搜索权益名称或描述包含关键词的权益
	}

	var total int64                                   // 声明总数变量，用于存储符合条件的会员权益总数
	if err := query.Count(&total).Error; err != nil { // 统计符合条件的会员权益总数，如果统计失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	var rightsList []hotel_admin.MemberRights                                                                         // 声明会员权益列表变量，用于存储查询到的会员权益信息列表
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&rightsList).Error; err != nil { // 按创建时间倒序排列，添加分页限制（偏移量、每页数量）并查询会员权益列表，如果查询失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	rightsInfos := make([]MemberRightsInfo, len(rightsList)) // 创建会员权益信息列表，长度为查询到的会员权益数量
	for i := range rightsList {                              // 遍历查询到的会员权益列表
		rightsInfos[i] = buildMemberRightsInfo(&rightsList[i]) // 调用构建函数，将每个会员权益实体对象转换为会员权益信息对象
	}

	return &ListMemberRightsResp{ // 返回会员权益列表响应对象
		List:     rightsInfos,   // 设置会员权益列表（转换后的会员权益信息列表）
		Total:    uint64(total), // 设置总数（转换为uint64类型）
		Page:     req.Page,      // 设置当前页码（从请求中获取）
		PageSize: req.PageSize,  // 设置每页数量（从请求中获取）
	}, nil // 返回响应对象和无错误
}

// DeleteMemberRights 删除会员权益（软删除）
func (s *MemberRightsService) DeleteMemberRights(id uint64) error {
	return db.MysqlDB.Delete(&hotel_admin.MemberRights{}, id).Error // 执行软删除操作（设置deleted_at字段），根据会员权益ID删除会员权益记录，返回删除操作的结果（成功为nil，失败为error）
}

// GetRightsByMemberLevel 根据会员等级获取权益列表
// 业务功能：查询指定会员等级的有效权益列表，用于会员权益展示和权益计算
// 入参说明：
//   - memberLevel: 会员等级（如普通会员、黄金会员、钻石会员）
//
// 返回值说明：
//   - []MemberRightsInfo: 符合条件的会员权益列表（只返回当前时间生效且未过期的权益）
//   - error: 查询失败错误
//
// 业务规则：权益必须满足以下条件：1）属于指定会员等级；2）状态为启用；3）生效时间小于等于当前时间；4）失效时间为空或大于等于当前时间
func (s *MemberRightsService) GetRightsByMemberLevel(memberLevel string) ([]MemberRightsInfo, error) {
	var rightsList []hotel_admin.MemberRights // 声明会员权益列表变量，用于存储查询到的会员权益信息列表
	now := time.Now()                         // 获取当前时间（用于判断权益是否已生效且未过期）

	// 业务规则：查询有效权益的条件：属于指定会员等级、状态为启用、已生效且未过期、未被删除
	err := db.MysqlDB. // 创建数据库查询构建器
				Where("member_level = ? AND status = ? AND effective_time <= ? AND (expire_time IS NULL OR expire_time >= ?) AND deleted_at IS NULL", // 添加筛选条件：会员等级匹配、状态为启用、生效时间小于等于当前时间、失效时间为空或大于等于当前时间、未被删除
						memberLevel, RightsStatusActive, now, now). // 设置参数值：会员等级、启用状态常量、当前时间（生效时间判断）、当前时间（失效时间判断）
		Order("created_at ASC"). // 添加排序条件，按创建时间正序排列（最早创建的权益排在前面）
		Find(&rightsList).Error  // 执行查询并获取符合条件的会员权益列表，如果查询失败则返回错误

	if err != nil { // 如果查询失败，则返回错误
		return nil, err // 返回nil和错误信息
	}

	rightsInfos := make([]MemberRightsInfo, len(rightsList)) // 创建会员权益信息列表，长度为查询到的会员权益数量
	for i := range rightsList {                              // 遍历查询到的会员权益列表
		rightsInfos[i] = buildMemberRightsInfo(&rightsList[i]) // 调用构建函数，将每个会员权益实体对象转换为会员权益信息对象
	}

	return rightsInfos, nil // 返回会员权益信息列表和无错误
}
