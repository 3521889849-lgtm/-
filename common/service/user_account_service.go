package service

import (
	"errors"
	"fmt"
	"time"

	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	AccountStatusActive   = "ACTIVE"
	AccountStatusInactive = "INACTIVE"
	AccountStatusLocked   = "LOCKED"
)

// UserAccountService 用户账号管理服务
// 负责处理酒店后台用户账号的创建、更新、查询、删除等核心业务逻辑，
// 包括账号与角色的关联校验、密码加密、用户名唯一性检查、账号状态管理等。
type UserAccountService struct{}

// CreateUserAccountReq 创建账号请求
type CreateUserAccountReq struct {
	Username     string  `json:"username" binding:"required"`
	Password     string  `json:"password" binding:"required"`
	RealName     string  `json:"real_name" binding:"required"`
	ContactPhone string  `json:"contact_phone" binding:"required"`
	RoleID       uint64  `json:"role_id" binding:"required"`
	BranchID     *uint64 `json:"branch_id,omitempty"`
	Status       string  `json:"status"`
}

// UpdateUserAccountReq 更新账号请求
type UpdateUserAccountReq struct {
	ID           uint64  `json:"id" binding:"required"`
	Username     *string `json:"username,omitempty"`
	Password     *string `json:"password,omitempty"`
	RealName     *string `json:"real_name,omitempty"`
	ContactPhone *string `json:"contact_phone,omitempty"`
	RoleID       *uint64 `json:"role_id,omitempty"`
	BranchID     *uint64 `json:"branch_id,omitempty"`
	Status       *string `json:"status,omitempty"`
}

// ListUserAccountsReq 账号列表查询请求
type ListUserAccountsReq struct {
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
	RoleID   *uint64 `json:"role_id,omitempty"`
	BranchID *uint64 `json:"branch_id,omitempty"`
	Status   *string `json:"status,omitempty"`
	Keyword  *string `json:"keyword,omitempty"`
}

// UserAccountInfo 账号信息
type UserAccountInfo struct {
	ID           uint64     `json:"id"`
	Username     string     `json:"username"`
	RealName     string     `json:"real_name"`
	ContactPhone string     `json:"contact_phone"`
	RoleID       uint64     `json:"role_id"`
	RoleName     string     `json:"role_name,omitempty"`
	BranchID     *uint64    `json:"branch_id,omitempty"`
	BranchName   string     `json:"branch_name,omitempty"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

// ListUserAccountsResp 账号列表响应
type ListUserAccountsResp struct {
	List     []UserAccountInfo `json:"list"`
	Total    uint64            `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}

// hashPassword 使用bcrypt加密密码
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码加密失败: %w", err)
	}
	return string(hash), nil
}

// buildUserAccountInfo 构建账号信息（避免代码重复）
func buildUserAccountInfo(account *hotel_admin.UserAccount) UserAccountInfo {
	info := UserAccountInfo{
		ID:           account.ID,
		Username:     account.Username,
		RealName:     account.RealName,
		ContactPhone: account.ContactPhone,
		RoleID:       account.RoleID,
		BranchID:     account.BranchID,
		Status:       account.Status,
		CreatedAt:    account.CreatedAt,
		LastLoginAt:  account.LastLoginAt,
	}

	if account.Role != nil {
		info.RoleName = account.Role.RoleName
	}
	if account.Branch != nil {
		info.BranchName = account.Branch.HotelName
	}

	return info
}

// CreateUserAccount 创建账号
// 业务功能：创建新的后台用户账号，建立账号与角色的关联关系，对密码进行加密存储
// 入参说明：
//   - req: 创建账号请求，包含用户名、密码、姓名、联系电话、角色ID、分店ID（可选）、状态（可选，默认ACTIVE）
//
// 返回值说明：
//   - error: 用户名已存在、角色ID不存在、密码加密失败、业务校验失败或数据库操作错误
func (s *UserAccountService) CreateUserAccount(req CreateUserAccountReq) error {
	// 业务规则：用户名必须唯一，检查该用户名是否已被使用
	var existingAccount hotel_admin.UserAccount                                         // 声明账号实体变量，用于存储查询到的已存在账号信息
	err := db.MysqlDB.Where("username = ?", req.Username).First(&existingAccount).Error // 通过用户名查询是否已存在账号记录，如果查询失败则说明该用户名未被使用
	if err == nil {                                                                     // 如果查询成功（未报错），说明该用户名已被使用
		return fmt.Errorf("用户名 %s 已存在", req.Username) // 返回错误信息，表示用户名已存在
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) { // 如果查询失败但不是记录不存在错误，说明是其他数据库错误
		return fmt.Errorf("查询用户名失败: %w", err) // 返回数据库查询错误信息
	}

	// 业务规则：账号必须关联到有效的角色，验证角色是否存在
	var role hotel_admin.Role                                         // 声明角色实体变量，用于存储查询到的角色信息
	if err := db.MysqlDB.First(&role, req.RoleID).Error; err != nil { // 通过角色ID查询角色信息，如果查询失败则说明角色不存在
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果错误是记录不存在错误
			return fmt.Errorf("角色ID %d 不存在", req.RoleID) // 返回错误信息，表示角色不存在
		}
		return fmt.Errorf("查询角色失败: %w", err) // 返回数据库查询错误信息
	}

	// 业务规则：密码必须加密存储，使用bcrypt算法进行单向加密（不可逆）
	hashedPassword, err := hashPassword(req.Password) // 调用密码加密函数，对明文密码进行bcrypt加密
	if err != nil {                                   // 如果加密失败，则返回错误
		return err // 返回加密错误
	}

	// 业务规则：账号状态默认设置为启用（ACTIVE），如果请求中提供了状态则使用请求的状态
	status := AccountStatusActive // 初始化账号状态为启用（默认值）
	if req.Status != "" {         // 如果请求中提供了状态（非空），则使用请求的状态
		status = req.Status // 更新账号状态为请求中提供的状态值
	}

	account := hotel_admin.UserAccount{ // 创建账号实体对象
		Username:     req.Username,     // 设置用户名（从请求中获取）
		Password:     hashedPassword,   // 设置密码（存储加密后的密码，而不是明文）
		RealName:     req.RealName,     // 设置真实姓名（从请求中获取）
		ContactPhone: req.ContactPhone, // 设置联系电话（从请求中获取）
		RoleID:       req.RoleID,       // 设置角色ID（从请求中获取）
		BranchID:     req.BranchID,     // 设置分店ID（从请求中获取，可为空）
		Status:       status,           // 设置账号状态（使用计算后的状态值）
	}

	return db.MysqlDB.Create(&account).Error // 将账号信息保存到数据库，返回保存操作的结果（成功为nil，失败为error）
}

// UpdateUserAccount 更新账号
// 业务功能：修改已存在账号的属性信息，支持部分字段更新（用户名、密码、姓名、联系电话、角色、分店、状态）
// 入参说明：
//   - req: 账号更新请求，所有字段均为可选，只更新传入的非空字段
//
// 返回值说明：
//   - error: 账号不存在、用户名已被使用（修改用户名时）、角色ID不存在（修改角色时）、密码加密失败或数据库操作错误
func (s *UserAccountService) UpdateUserAccount(req UpdateUserAccountReq) error {
	updates := make(map[string]interface{}) // 创建更新字段映射表，用于存储需要更新的字段和值

	// 业务逻辑：如果修改用户名，需重新校验唯一性，确保新用户名未被其他账号使用
	if req.Username != nil { // 如果请求中提供了用户名（指针非空），则需要重新校验唯一性
		var existingAccount hotel_admin.UserAccount                                                              // 声明账号实体变量，用于存储查询到的已存在账号信息
		err := db.MysqlDB.Where("username = ? AND id != ?", *req.Username, req.ID).First(&existingAccount).Error // 通过用户名查询是否已存在其他账号使用该用户名（排除当前账号），如果查询失败则说明该用户名未被使用
		if err == nil {                                                                                          // 如果查询成功（未报错），说明该用户名已被其他账号使用
			return fmt.Errorf("用户名 %s 已被使用", *req.Username) // 返回错误信息，表示用户名已被使用
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) { // 如果错误不是记录不存在错误（是其他数据库错误）
			return fmt.Errorf("查询用户名失败: %w", err) // 返回数据库查询错误信息
		}
		updates["username"] = *req.Username // 添加用户名到更新映射表（解引用指针获取值）
	}

	// 业务逻辑：如果修改密码，需重新加密后存储，确保密码安全性
	if req.Password != nil { // 如果请求中提供了密码（指针非空），则需要重新加密
		hashedPassword, err := hashPassword(*req.Password) // 调用密码加密函数，对明文密码进行bcrypt加密（解引用指针获取值）
		if err != nil {                                    // 如果加密失败，则返回错误
			return err // 返回加密错误
		}
		updates["password"] = hashedPassword // 添加加密后的密码到更新映射表
	}

	if req.RealName != nil { // 如果请求中提供了真实姓名（指针非空），则添加到更新映射表
		updates["real_name"] = *req.RealName // 添加真实姓名到更新映射表（解引用指针获取值）
	}

	if req.ContactPhone != nil { // 如果请求中提供了联系电话（指针非空），则添加到更新映射表
		updates["contact_phone"] = *req.ContactPhone // 添加联系电话到更新映射表（解引用指针获取值）
	}

	// 业务逻辑：如果修改角色，需验证新角色是否存在
	if req.RoleID != nil { // 如果请求中提供了角色ID（指针非空），则需要验证角色是否存在
		var role hotel_admin.Role                                          // 声明角色实体变量，用于存储查询到的角色信息
		if err := db.MysqlDB.First(&role, *req.RoleID).Error; err != nil { // 通过角色ID查询角色信息（解引用指针获取值），如果查询失败则说明角色不存在
			if errors.Is(err, gorm.ErrRecordNotFound) { // 如果错误是记录不存在错误
				return fmt.Errorf("角色ID %d 不存在", *req.RoleID) // 返回错误信息，表示角色不存在
			}
			return fmt.Errorf("查询角色失败: %w", err) // 返回数据库查询错误信息
		}
		updates["role_id"] = *req.RoleID // 添加角色ID到更新映射表（解引用指针获取值）
	}

	// 业务逻辑：采用部分更新策略，只更新请求中提供的非空字段
	if req.BranchID != nil { // 如果请求中提供了分店ID（指针非空），则添加到更新映射表
		updates["branch_id"] = *req.BranchID // 添加分店ID到更新映射表（解引用指针获取值）
	}

	if req.Status != nil { // 如果请求中提供了状态（指针非空），则添加到更新映射表
		updates["status"] = *req.Status // 添加状态到更新映射表（解引用指针获取值）
	}

	// 业务规则：如果没有需要更新的字段，直接返回，避免无效的数据库操作
	if len(updates) == 0 { // 如果更新映射表为空（没有需要更新的字段）
		return nil // 直接返回nil，表示更新成功（实际上没有更新任何字段）
	}

	// 执行更新，并检查是否有记录被影响（用于判断账号是否存在）
	result := db.MysqlDB.Model(&hotel_admin.UserAccount{}).Where("id = ?", req.ID).Updates(updates) // 根据账号ID更新账号信息，使用更新映射表中的字段和值
	if result.Error != nil {                                                                        // 如果更新操作失败，则返回错误
		return result.Error // 返回数据库操作错误
	}
	if result.RowsAffected == 0 { // 如果更新的记录数为0（没有记录被影响），说明账号不存在
		return fmt.Errorf("账号ID %d 不存在", req.ID) // 返回错误信息，表示账号不存在
	}

	return nil // 返回nil表示更新成功
}

// GetUserAccount 获取账号详情
// 业务功能：根据账号ID查询账号的完整信息，包含关联的角色和分店信息
// 入参说明：
//   - id: 账号ID
//
// 返回值说明：
//   - *UserAccountInfo: 账号完整信息（包含角色名称、分店名称等关联数据）
//   - error: 账号不存在或查询失败
func (s *UserAccountService) GetUserAccount(id uint64) (*UserAccountInfo, error) {
	var account hotel_admin.UserAccount                                                            // 声明账号实体变量，用于存储查询到的账号信息
	if err := db.MysqlDB.Preload("Role").Preload("Branch").First(&account, id).Error; err != nil { // 通过账号ID查询账号信息（预加载角色和分店信息），如果查询失败则说明账号不存在
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果错误是记录不存在错误
			return nil, fmt.Errorf("账号ID %d 不存在", id) // 返回nil和错误信息，表示账号不存在
		}
		return nil, err // 返回nil和其他数据库查询错误
	}

	info := buildUserAccountInfo(&account) // 调用构建函数，将账号实体对象转换为账号信息对象（包含角色和分店信息）
	return &info, nil                      // 返回账号信息指针和无错误
}

// ListUserAccounts 获取账号列表
// 业务功能：支持多条件筛选和分页查询账号列表，用于账号管理和权限分配场景
// 入参说明：
//   - req: 账号列表查询请求，支持按角色ID、分店ID、状态筛选，支持关键词搜索（用户名/姓名/联系电话），支持分页
//
// 返回值说明：
//   - *ListUserAccountsResp: 符合条件的账号列表（包含角色和分店信息）及分页信息
//   - error: 查询失败错误
func (s *UserAccountService) ListUserAccounts(req ListUserAccountsReq) (*ListUserAccountsResp, error) {
	// 业务规则：分页参数默认值设置，页码最小为1，每页数量最小1条，最大不超过100条
	req.Page = max(req.Page, 1)                   // 如果页码小于1，则设置为1（使用max函数确保最小值）
	req.PageSize = min(max(req.PageSize, 1), 100) // 如果每页数量小于1则设置为1，如果大于100则设置为100（使用min和max函数确保范围）
	offset := (req.Page - 1) * req.PageSize       // 计算分页偏移量（跳过前面的记录数），用于SQL查询的OFFSET子句

	query := db.MysqlDB.Model(&hotel_admin.UserAccount{}).Where("deleted_at IS NULL") // 创建账号模型的查询构建器，添加软删除筛选条件（只查询未删除的账号）

	// 业务筛选：按角色ID筛选，支持查看特定角色下的所有账号
	if req.RoleID != nil { // 如果请求中提供了角色ID（指针非空），则添加角色筛选条件
		query = query.Where("role_id = ?", *req.RoleID) // 添加角色ID筛选条件，只查询指定角色的账号（解引用指针获取值）
	}
	// 业务筛选：按分店ID筛选，支持查看特定分店下的所有账号
	if req.BranchID != nil { // 如果请求中提供了分店ID（指针非空），则添加分店筛选条件
		query = query.Where("branch_id = ?", *req.BranchID) // 添加分店ID筛选条件，只查询指定分店的账号（解引用指针获取值）
	}
	// 业务筛选：按账号状态筛选（如启用/停用/锁定），支持状态维度管理
	if req.Status != nil && *req.Status != "" { // 如果请求中提供了状态（指针非空且值非空），则添加状态筛选条件
		query = query.Where("status = ?", *req.Status) // 添加账号状态筛选条件，只查询指定状态的账号（解引用指针获取值）
	}
	// 业务搜索：关键词多字段模糊搜索，支持用户名、姓名、联系电话三个维度同时搜索
	if req.Keyword != nil && *req.Keyword != "" { // 如果请求中提供了关键词（指针非空且值非空），则添加关键词搜索条件
		keyword := "%" + *req.Keyword + "%"                                                                           // 构建模糊搜索关键词（前后加%通配符）
		query = query.Where("username LIKE ? OR real_name LIKE ? OR contact_phone LIKE ?", keyword, keyword, keyword) // 添加关键词搜索条件，搜索用户名、姓名或联系电话包含关键词的账号
	}

	var total int64                                   // 声明总数变量，用于存储符合条件的账号总数
	if err := query.Count(&total).Error; err != nil { // 统计符合条件的账号总数，如果统计失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	var accounts []hotel_admin.UserAccount                                      // 声明账号列表变量，用于存储查询到的账号信息列表
	if err := query.Preload("Role").Preload("Branch").Order("created_at DESC"). // 预加载角色和分店信息关联数据，按创建时间倒序排列
											Offset(offset).Limit(req.PageSize).Find(&accounts).Error; err != nil { // 添加分页限制（偏移量、每页数量）并查询账号列表，如果查询失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	accountInfos := make([]UserAccountInfo, len(accounts)) // 创建账号信息列表，长度为查询到的账号数量
	for i := range accounts {                              // 遍历查询到的账号列表
		accountInfos[i] = buildUserAccountInfo(&accounts[i]) // 调用构建函数，将每个账号实体对象转换为账号信息对象（包含角色和分店信息）
	}

	return &ListUserAccountsResp{ // 返回账号列表响应对象
		List:     accountInfos,  // 设置账号列表（转换后的账号信息列表）
		Total:    uint64(total), // 设置总数（转换为uint64类型）
		Page:     req.Page,      // 设置当前页码（从请求中获取）
		PageSize: req.PageSize,  // 设置每页数量（从请求中获取）
	}, nil // 返回响应对象和无错误
}

// DeleteUserAccount 删除账号（软删除）
// 业务功能：逻辑删除账号记录，不物理删除数据，保留历史操作记录和权限关联关系
// 入参说明：
//   - id: 待删除的账号ID
//
// 返回值说明：
//   - error: 账号不存在或数据库操作错误
func (s *UserAccountService) DeleteUserAccount(id uint64) error {
	return db.MysqlDB.Delete(&hotel_admin.UserAccount{}, id).Error // 执行软删除操作（设置deleted_at字段），根据账号ID删除账号记录，返回删除操作的结果（成功为nil，失败为error）
}
