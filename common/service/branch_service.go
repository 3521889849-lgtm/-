package service

import (
	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"
	"fmt"
)

// BranchService 分店管理服务
// 负责处理酒店分店的创建、更新、查询、删除等核心业务逻辑，
// 包括分店编码自动生成、分店与房源的关联关系检查、分店状态管理等。
type BranchService struct{}

// ListBranchesReq 分店列表查询请求
type ListBranchesReq struct {
	Status *string `json:"status,omitempty"` // 状态筛选，可选
}

// BranchInfo 分店信息
type BranchInfo struct {
	ID           uint64 `json:"id"`
	HotelName    string `json:"hotel_name"`
	BranchCode   string `json:"branch_code"`
	Address      string `json:"address"`
	Contact      string `json:"contact"`
	ContactPhone string `json:"contact_phone"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// ListBranches 获取分店列表
// 业务功能：查询所有分店信息，支持按状态筛选，用于分店选择和统计场景
// 入参说明：
//   - req: 分店列表查询请求，支持按状态筛选（可选）
//
// 返回值说明：
//   - []BranchInfo: 符合条件的分店列表（按创建时间倒序）
//   - error: 查询失败错误
func (s *BranchService) ListBranches(req ListBranchesReq) ([]BranchInfo, error) {
	query := db.MysqlDB.Model(&hotel_admin.HotelBranch{}). // 创建分店模型的查询构建器
								Where("deleted_at IS NULL") // 添加软删除筛选条件（只查询未删除的分店）

	// 业务筛选：按分店状态筛选（如启用/停用），支持状态维度管理
	if req.Status != nil && *req.Status != "" { // 如果请求中提供了状态（指针非空且值非空），则添加状态筛选条件
		query = query.Where("status = ?", *req.Status) // 添加分店状态筛选条件，只查询指定状态的分店（解引用指针获取值）
	}

	// 业务排序：按创建时间倒序排列，最新创建的分店显示在最前面
	query = query.Order("created_at DESC") // 添加排序条件，按创建时间倒序排列（最新创建的分店排在前面）

	var branches []hotel_admin.HotelBranch              // 声明分店列表变量，用于存储查询到的分店信息列表
	if err := query.Find(&branches).Error; err != nil { // 执行查询并获取符合条件的分店列表，如果查询失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	// 转换为返回格式
	result := make([]BranchInfo, len(branches)) // 创建分店信息列表，长度为查询到的分店数量
	for i, branch := range branches {           // 遍历查询到的分店列表
		result[i] = BranchInfo{ // 创建分店信息对象
			ID:           branch.ID,                                      // 设置分店ID（从分店实体中获取）
			HotelName:    branch.HotelName,                               // 设置酒店名称（从分店实体中获取）
			BranchCode:   branch.BranchCode,                              // 设置分店编码（从分店实体中获取）
			Address:      branch.Address,                                 // 设置地址（从分店实体中获取）
			Contact:      branch.Contact,                                 // 设置联系人（从分店实体中获取）
			ContactPhone: branch.ContactPhone,                            // 设置联系电话（从分店实体中获取）
			Status:       branch.Status,                                  // 设置状态（从分店实体中获取）
			CreatedAt:    branch.CreatedAt.Format("2006-01-02 15:04:05"), // 设置创建时间（格式化时间字符串）
			UpdatedAt:    branch.UpdatedAt.Format("2006-01-02 15:04:05"), // 设置更新时间（格式化时间字符串）
		}
	}

	return result, nil // 返回分店信息列表和无错误
}

// GetBranch 获取分店详情
func (s *BranchService) GetBranch(branchID uint64) (*BranchInfo, error) {
	var branch hotel_admin.HotelBranch                                // 声明分店实体变量，用于存储查询到的分店信息
	if err := db.MysqlDB.First(&branch, branchID).Error; err != nil { // 通过分店ID查询分店信息，如果查询失败则说明分店不存在
		return nil, err // 返回nil和数据库查询错误
	}

	return &BranchInfo{ // 返回分店信息对象指针
		ID:           branch.ID,                                      // 设置分店ID（从分店实体中获取）
		HotelName:    branch.HotelName,                               // 设置酒店名称（从分店实体中获取）
		BranchCode:   branch.BranchCode,                              // 设置分店编码（从分店实体中获取）
		Address:      branch.Address,                                 // 设置地址（从分店实体中获取）
		Contact:      branch.Contact,                                 // 设置联系人（从分店实体中获取）
		ContactPhone: branch.ContactPhone,                            // 设置联系电话（从分店实体中获取）
		Status:       branch.Status,                                  // 设置状态（从分店实体中获取）
		CreatedAt:    branch.CreatedAt.Format("2006-01-02 15:04:05"), // 设置创建时间（格式化时间字符串）
		UpdatedAt:    branch.UpdatedAt.Format("2006-01-02 15:04:05"), // 设置更新时间（格式化时间字符串）
	}, nil // 返回分店信息指针和无错误
}

// CreateBranchReq 创建分店请求
type CreateBranchReq struct {
	HotelName    string `json:"hotel_name" binding:"required"`
	BranchCode   string `json:"branch_code"`
	Address      string `json:"address" binding:"required"`
	Contact      string `json:"contact" binding:"required"`
	ContactPhone string `json:"contact_phone" binding:"required"`
	Status       string `json:"status"`
}

// UpdateBranchReq 更新分店请求
type UpdateBranchReq struct {
	HotelName    string `json:"hotel_name"`
	BranchCode   string `json:"branch_code"`
	Address      string `json:"address"`
	Contact      string `json:"contact"`
	ContactPhone string `json:"contact_phone"`
	Status       string `json:"status"`
}

// CreateBranch 创建分店
// 业务功能：创建新的酒店分店，自动生成分店编码（如未提供），初始化分店状态为启用
// 入参说明：
//   - req: 创建分店请求，包含酒店名称、分店编码（可选，自动生成）、地址、联系人、联系电话、状态（可选，默认ACTIVE）
//
// 返回值说明：
//   - *hotel_admin.HotelBranch: 成功创建后的分店完整信息（包含自动生成的ID和编码）
//   - error: 数据库操作错误
func (s *BranchService) CreateBranch(req *CreateBranchReq) (*hotel_admin.HotelBranch, error) {
	branch := &hotel_admin.HotelBranch{ // 创建分店实体对象指针
		HotelName:    req.HotelName,    // 设置酒店名称（从请求中获取）
		BranchCode:   req.BranchCode,   // 设置分店编码（从请求中获取，可为空）
		Address:      req.Address,      // 设置地址（从请求中获取）
		Contact:      req.Contact,      // 设置联系人（从请求中获取）
		ContactPhone: req.ContactPhone, // 设置联系电话（从请求中获取）
		Status:       req.Status,       // 设置状态（从请求中获取，可为空）
		CreatedBy:    1,                // 设置创建人ID（TODO: 从登录用户获取）
	}

	// 业务规则：分店状态默认设置为启用（ACTIVE），如果请求中提供了状态则使用请求的状态
	if branch.Status == "" { // 如果分店状态为空（未提供），则设置为默认值
		branch.Status = "ACTIVE" // 设置分店状态为启用（默认值）
	}

	// 业务规则：如果未提供分店编码，自动生成唯一的分店编码（格式：BR001、BR002...）
	if branch.BranchCode == "" { // 如果分店编码为空（未提供），则自动生成
		// 业务逻辑：基于现有分店数量自动生成编码，确保编码唯一性
		var count int64                                            // 声明计数变量，用于存储现有分店数量
		db.MysqlDB.Model(&hotel_admin.HotelBranch{}).Count(&count) // 统计现有分店数量（包含所有分店，不排除已删除）
		branch.BranchCode = fmt.Sprintf("BR%03d", count+1)         // 根据现有数量生成新的分店编码（格式：BR001、BR002...），使用三位数字补零
	}

	if err := db.MysqlDB.Create(branch).Error; err != nil { // 将分店信息保存到数据库，如果保存失败则返回错误
		return nil, err // 返回nil和数据库操作错误
	}

	return branch, nil // 返回成功创建后的分店信息和无错误
}

// UpdateBranch 更新分店
func (s *BranchService) UpdateBranch(branchID uint64, req *UpdateBranchReq) error {
	var branch hotel_admin.HotelBranch
	if err := db.MysqlDB.First(&branch, branchID).Error; err != nil {
		return err
	}

	updates := make(map[string]interface{})
	if req.HotelName != "" {
		updates["hotel_name"] = req.HotelName
	}
	if req.BranchCode != "" {
		updates["branch_code"] = req.BranchCode
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Contact != "" {
		updates["contact"] = req.Contact
	}
	if req.ContactPhone != "" {
		updates["contact_phone"] = req.ContactPhone
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if len(updates) > 0 {
		if err := db.MysqlDB.Model(&branch).Updates(updates).Error; err != nil {
			return err
		}
	}

	return nil
}

// DeleteBranch 删除分店（软删除）
// 业务功能：逻辑删除分店记录，不物理删除数据，但需检查是否有房源关联到该分店
// 入参说明：
//   - branchID: 待删除的分店ID
//
// 返回值说明：
//   - error: 分店不存在、分店下有房源（无法删除）或数据库操作错误
func (s *BranchService) DeleteBranch(branchID uint64) error {
	// 业务规则：只能删除已存在的分店，验证分店是否存在
	var branch hotel_admin.HotelBranch
	if err := db.MysqlDB.First(&branch, branchID).Error; err != nil {
		return err
	}

	// 业务规则：如果分店下有关联的房源，不允许删除，避免破坏房源数据的完整性
	var roomCount int64
	db.MysqlDB.Model(&hotel_admin.RoomInfo{}).Where("branch_id = ? AND deleted_at IS NULL", branchID).Count(&roomCount)
	if roomCount > 0 {
		return fmt.Errorf("该分店下有 %d 个房源，无法删除", roomCount)
	}

	// 执行软删除：设置deleted_at字段，不物理删除记录
	return db.MysqlDB.Delete(&branch).Error
}
