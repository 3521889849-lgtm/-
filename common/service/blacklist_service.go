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
	BlacklistStatusValid   = "VALID"
	BlacklistStatusInvalid = "INVALID"
)

// BlacklistService 黑名单管理服务
// 负责处理酒店黑名单的创建、更新、查询、删除等核心业务逻辑，
// 包括黑名单与客人的关联管理、黑名单状态管理、黑名单有效性检查等。
type BlacklistService struct{}

// CreateBlacklistReq 创建黑名单请求
type CreateBlacklistReq struct {
	GuestID    *uint64 `json:"guest_id,omitempty"`
	IDNumber   string  `json:"id_number" binding:"required"`
	Phone      string  `json:"phone" binding:"required"`
	Reason     string  `json:"reason" binding:"required"`
	OperatorID uint64  `json:"operator_id" binding:"required"`
	Status     string  `json:"status"`
}

// UpdateBlacklistReq 更新黑名单请求
type UpdateBlacklistReq struct {
	ID       uint64  `json:"id" binding:"required"`
	GuestID  *uint64 `json:"guest_id,omitempty"`
	IDNumber *string `json:"id_number,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Reason   *string `json:"reason,omitempty"`
	Status   *string `json:"status,omitempty"`
}

// ListBlacklistsReq 黑名单列表查询请求
type ListBlacklistsReq struct {
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
	Status   *string `json:"status,omitempty"`
	Keyword  *string `json:"keyword,omitempty"`
}

// BlacklistInfo 黑名单信息
type BlacklistInfo struct {
	ID         uint64    `json:"id"`
	GuestID    *uint64   `json:"guest_id,omitempty"`
	GuestName  string    `json:"guest_name,omitempty"`
	IDNumber   string    `json:"id_number"`
	Phone      string    `json:"phone"`
	Reason     string    `json:"reason"`
	BlackTime  time.Time `json:"black_time"`
	OperatorID uint64    `json:"operator_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// ListBlacklistsResp 黑名单列表响应
type ListBlacklistsResp struct {
	List     []BlacklistInfo `json:"list"`
	Total    uint64          `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

// buildBlacklistInfo 构建黑名单信息（避免代码重复）
func buildBlacklistInfo(blacklist *hotel_admin.Blacklist) BlacklistInfo {
	info := BlacklistInfo{
		ID:         blacklist.ID,
		GuestID:    blacklist.GuestID,
		IDNumber:   blacklist.IDNumber,
		Phone:      blacklist.Phone,
		Reason:     blacklist.Reason,
		BlackTime:  blacklist.BlackTime,
		OperatorID: blacklist.OperatorID,
		Status:     blacklist.Status,
		CreatedAt:  blacklist.CreatedAt,
	}

	if blacklist.Guest != nil {
		info.GuestName = blacklist.Guest.Name
	}

	return info
}

// CreateBlacklist 创建黑名单
func (s *BlacklistService) CreateBlacklist(req CreateBlacklistReq) error {
	var existingBlacklist hotel_admin.Blacklist                            // 声明黑名单实体变量，用于存储查询到的已存在黑名单信息
	err := db.MysqlDB.Where("(id_number = ? OR phone = ?) AND status = ?", // 查询是否已存在该证件号或手机号的有效黑名单记录（状态为有效）
		req.IDNumber, req.Phone, BlacklistStatusValid).First(&existingBlacklist).Error // 设置查询参数：证件号、手机号、有效状态常量
	if err == nil { // 如果查询成功（未报错），说明该证件号或手机号已在黑名单中
		return errors.New("该证件号或手机号已在黑名单中") // 返回错误信息，表示该证件号或手机号已在黑名单中
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) { // 如果错误不是记录不存在错误（是其他数据库错误）
		return fmt.Errorf("查询黑名单失败: %w", err) // 返回数据库查询错误信息
	}

	status := BlacklistStatusValid // 初始化黑名单状态为有效（默认值）
	if req.Status != "" {          // 如果请求中提供了状态（非空），则使用请求的状态
		status = req.Status // 更新黑名单状态为请求中提供的状态值
	}

	blacklist := hotel_admin.Blacklist{ // 创建黑名单实体对象
		GuestID:    req.GuestID,    // 设置客人ID（从请求中获取，可为空）
		IDNumber:   req.IDNumber,   // 设置证件号（从请求中获取）
		Phone:      req.Phone,      // 设置手机号（从请求中获取）
		Reason:     req.Reason,     // 设置拉黑原因（从请求中获取）
		BlackTime:  time.Now(),     // 设置拉黑时间为当前时间（自动生成）
		OperatorID: req.OperatorID, // 设置操作人ID（从请求中获取）
		Status:     status,         // 设置黑名单状态（使用计算后的状态值）
	}

	return db.MysqlDB.Create(&blacklist).Error // 将黑名单信息保存到数据库，返回保存操作的结果（成功为nil，失败为error）
}

// UpdateBlacklist 更新黑名单
func (s *BlacklistService) UpdateBlacklist(req UpdateBlacklistReq) error {
	updates := make(map[string]interface{}) // 创建更新字段映射表，用于存储需要更新的字段和值

	if req.GuestID != nil { // 如果请求中提供了客人ID（指针非空），则添加到更新映射表
		updates["guest_id"] = *req.GuestID // 添加客人ID到更新映射表（解引用指针获取值）
	}
	if req.IDNumber != nil { // 如果请求中提供了证件号（指针非空），则添加到更新映射表
		updates["id_number"] = *req.IDNumber // 添加证件号到更新映射表（解引用指针获取值）
	}
	if req.Phone != nil { // 如果请求中提供了手机号（指针非空），则添加到更新映射表
		updates["phone"] = *req.Phone // 添加手机号到更新映射表（解引用指针获取值）
	}
	if req.Reason != nil { // 如果请求中提供了拉黑原因（指针非空），则添加到更新映射表
		updates["reason"] = *req.Reason // 添加拉黑原因到更新映射表（解引用指针获取值）
	}
	if req.Status != nil { // 如果请求中提供了状态（指针非空），则添加到更新映射表
		updates["status"] = *req.Status // 添加状态到更新映射表（解引用指针获取值）
	}

	if len(updates) == 0 { // 如果更新映射表为空（没有需要更新的字段）
		return nil // 直接返回nil，表示更新成功（实际上没有更新任何字段）
	}

	result := db.MysqlDB.Model(&hotel_admin.Blacklist{}).Where("id = ?", req.ID).Updates(updates) // 根据黑名单ID更新黑名单信息，使用更新映射表中的字段和值
	if result.Error != nil {                                                                      // 如果更新操作失败，则返回错误
		return result.Error // 返回数据库操作错误
	}
	if result.RowsAffected == 0 { // 如果更新的记录数为0（没有记录被影响），说明黑名单不存在
		return fmt.Errorf("黑名单ID %d 不存在", req.ID) // 返回错误信息，表示黑名单不存在
	}

	return nil // 返回nil表示更新成功
}

// GetBlacklist 获取黑名单详情
func (s *BlacklistService) GetBlacklist(id uint64) (*BlacklistInfo, error) {
	var blacklist hotel_admin.Blacklist                                             // 声明黑名单实体变量，用于存储查询到的黑名单信息
	if err := db.MysqlDB.Preload("Guest").First(&blacklist, id).Error; err != nil { // 通过黑名单ID查询黑名单信息（预加载客人信息），如果查询失败则说明黑名单不存在
		return nil, errors.New("黑名单记录不存在") // 返回nil和错误信息，表示黑名单记录不存在
	}

	blacklistInfo := &BlacklistInfo{ // 创建黑名单信息对象指针
		ID:         blacklist.ID,         // 设置黑名单ID（从黑名单实体中获取）
		GuestID:    blacklist.GuestID,    // 设置客人ID（从黑名单实体中获取，可为空）
		IDNumber:   blacklist.IDNumber,   // 设置证件号（从黑名单实体中获取）
		Phone:      blacklist.Phone,      // 设置手机号（从黑名单实体中获取）
		Reason:     blacklist.Reason,     // 设置拉黑原因（从黑名单实体中获取）
		BlackTime:  blacklist.BlackTime,  // 设置拉黑时间（从黑名单实体中获取）
		OperatorID: blacklist.OperatorID, // 设置操作人ID（从黑名单实体中获取）
		Status:     blacklist.Status,     // 设置状态（从黑名单实体中获取）
		CreatedAt:  blacklist.CreatedAt,  // 设置创建时间（从黑名单实体中获取）
	}

	if blacklist.Guest != nil { // 如果黑名单关联了客人信息（预加载的数据）
		blacklistInfo.GuestName = blacklist.Guest.Name // 设置客人姓名（从关联的客人信息中获取）
	}

	return blacklistInfo, nil // 返回黑名单信息指针和无错误
}

// ListBlacklists 获取黑名单列表
func (s *BlacklistService) ListBlacklists(req ListBlacklistsReq) (*ListBlacklistsResp, error) {
	if req.Page <= 0 { // 如果页码小于等于0，则设置为1
		req.Page = 1 // 设置页码为1（使用默认值）
	}
	if req.PageSize <= 0 { // 如果每页数量小于等于0，则设置为10
		req.PageSize = 10 // 设置每页数量为10（使用默认值）
	}
	if req.PageSize > 100 { // 如果每页数量大于100，则设置为100
		req.PageSize = 100 // 设置每页数量为100（最大值限制）
	}

	offset := (req.Page - 1) * req.PageSize // 计算分页偏移量（跳过前面的记录数），用于SQL查询的OFFSET子句

	query := db.MysqlDB.Model(&hotel_admin.Blacklist{}). // 创建黑名单模型的查询构建器
								Preload("Guest").           // 预加载客人信息关联数据（JOIN查询客人信息）
								Where("deleted_at IS NULL") // 添加软删除筛选条件（只查询未删除的黑名单）

	if req.Status != nil && *req.Status != "" { // 如果请求中提供了状态（指针非空且值非空），则添加状态筛选条件
		query = query.Where("status = ?", *req.Status) // 添加黑名单状态筛选条件，只查询指定状态的黑名单（解引用指针获取值）
	}
	if req.Keyword != nil && *req.Keyword != "" { // 如果请求中提供了关键词（指针非空且值非空），则添加关键词搜索条件
		keyword := "%" + *req.Keyword + "%"                                                                 // 构建模糊搜索关键词（前后加%通配符）
		query = query.Where("id_number LIKE ? OR phone LIKE ? OR reason LIKE ?", keyword, keyword, keyword) // 添加关键词搜索条件，搜索证件号、手机号或拉黑原因包含关键词的黑名单
	}

	var total int64                                   // 声明总数变量，用于存储符合条件的黑名单总数
	if err := query.Count(&total).Error; err != nil { // 统计符合条件的黑名单总数，如果统计失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	var blacklists []hotel_admin.Blacklist                                                                            // 声明黑名单列表变量，用于存储查询到的黑名单信息列表
	if err := query.Order("black_time DESC").Offset(offset).Limit(req.PageSize).Find(&blacklists).Error; err != nil { // 按拉黑时间倒序排列，添加分页限制（偏移量、每页数量）并查询黑名单列表，如果查询失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	blacklistInfos := make([]BlacklistInfo, len(blacklists)) // 创建黑名单信息列表，长度为查询到的黑名单数量
	for i, blacklist := range blacklists {                   // 遍历查询到的黑名单列表
		blacklistInfos[i] = BlacklistInfo{ // 创建黑名单信息对象
			ID:         blacklist.ID,         // 设置黑名单ID（从黑名单实体中获取）
			GuestID:    blacklist.GuestID,    // 设置客人ID（从黑名单实体中获取，可为空）
			IDNumber:   blacklist.IDNumber,   // 设置证件号（从黑名单实体中获取）
			Phone:      blacklist.Phone,      // 设置手机号（从黑名单实体中获取）
			Reason:     blacklist.Reason,     // 设置拉黑原因（从黑名单实体中获取）
			BlackTime:  blacklist.BlackTime,  // 设置拉黑时间（从黑名单实体中获取）
			OperatorID: blacklist.OperatorID, // 设置操作人ID（从黑名单实体中获取）
			Status:     blacklist.Status,     // 设置状态（从黑名单实体中获取）
			CreatedAt:  blacklist.CreatedAt,  // 设置创建时间（从黑名单实体中获取）
		}

		if blacklist.Guest != nil { // 如果黑名单关联了客人信息（预加载的数据）
			blacklistInfos[i].GuestName = blacklist.Guest.Name // 设置客人姓名（从关联的客人信息中获取）
		}
	}

	return &ListBlacklistsResp{ // 返回黑名单列表响应对象
		List:     blacklistInfos, // 设置黑名单列表（转换后的黑名单信息列表）
		Total:    uint64(total),  // 设置总数（转换为uint64类型）
		Page:     req.Page,       // 设置当前页码（从请求中获取）
		PageSize: req.PageSize,   // 设置每页数量（从请求中获取）
	}, nil // 返回响应对象和无错误
}

// DeleteBlacklist 删除黑名单（软删除）
func (s *BlacklistService) DeleteBlacklist(id uint64) error {
	return db.MysqlDB.Delete(&hotel_admin.Blacklist{}, id).Error // 执行软删除操作（设置deleted_at字段），根据黑名单ID删除黑名单记录，返回删除操作的结果（成功为nil，失败为error）
}
