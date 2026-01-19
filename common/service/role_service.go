package service

import (
	"errors"
	"fmt"

	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"

	"gorm.io/gorm"
)

const (
	RoleStatusActive   = "ACTIVE"
	RoleStatusInactive = "INACTIVE"
)

// RoleService 角色管理服务
// 负责处理酒店后台角色的创建、更新、查询、删除等核心业务逻辑，
// 包括角色与权限的关联管理、角色名称唯一性检查、角色使用情况检查等。
type RoleService struct{}

// CreateRoleReq 创建角色请求
type CreateRoleReq struct {
	RoleName      string   `json:"role_name" binding:"required"`
	Description   *string  `json:"description,omitempty"`
	Status        string   `json:"status"`
	PermissionIDs []uint64 `json:"permission_ids,omitempty"` // 权限ID列表
}

// UpdateRoleReq 更新角色请求
type UpdateRoleReq struct {
	ID            uint64   `json:"id" binding:"required"`
	RoleName      *string  `json:"role_name,omitempty"`
	Description   *string  `json:"description,omitempty"`
	Status        *string  `json:"status,omitempty"`
	PermissionIDs []uint64 `json:"permission_ids,omitempty"` // 权限ID列表
}

// ListRolesReq 角色列表查询请求
type ListRolesReq struct {
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
	Status   *string `json:"status,omitempty"`
	Keyword  *string `json:"keyword,omitempty"`
}

// RoleInfo 角色信息
type RoleInfo struct {
	ID          uint64           `json:"id"`
	RoleName    string           `json:"role_name"`
	Description *string          `json:"description,omitempty"`
	Status      string           `json:"status"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`
	Permissions []PermissionInfo `json:"permissions,omitempty"` // 权限列表
}

// ListRolesResp 角色列表响应
type ListRolesResp struct {
	List     []RoleInfo `json:"list"`
	Total    uint64     `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// buildRoleInfo 构建角色信息（避免代码重复）
func buildRoleInfo(role *hotel_admin.Role, includePermissions bool) RoleInfo {
	info := RoleInfo{
		ID:          role.ID,
		RoleName:    role.RoleName,
		Description: role.Description,
		Status:      role.Status,
		CreatedAt:   role.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if includePermissions && len(role.RolePermissionRelations) > 0 {
		permissions := make([]PermissionInfo, 0, len(role.RolePermissionRelations))
		for _, relation := range role.RolePermissionRelations {
			if relation.Permission != nil {
				permissions = append(permissions, PermissionInfo{
					ID:             relation.Permission.ID,
					PermissionName: relation.Permission.PermissionName,
					PermissionURL:  relation.Permission.PermissionURL,
					PermissionType: relation.Permission.PermissionType,
					ParentID:       relation.Permission.ParentID,
					Status:         relation.Permission.Status,
				})
			}
		}
		info.Permissions = permissions
	}

	return info
}

// CreateRole 创建角色
// 业务功能：创建新的后台角色，建立角色与权限的关联关系，用于权限管理和权限分配
// 入参说明：
//   - req: 创建角色请求，包含角色名称、描述、状态（可选，默认ACTIVE）、权限ID列表
//
// 返回值说明：
//   - error: 角色名称已存在、权限分配失败或数据库操作错误
func (s *RoleService) CreateRole(req CreateRoleReq) error {
	// 业务规则：角色名称必须唯一，检查该角色名称是否已被使用
	var count int64                                                                                       // 声明计数器变量，用于存储查询到的记录数量
	err := db.MysqlDB.Model(&hotel_admin.Role{}).Where("role_name = ?", req.RoleName).Count(&count).Error // 通过角色名称统计是否存在角色记录，如果查询失败则说明数据库操作异常
	if err != nil {                                                                                       // 如果查询失败，说明数据库操作异常
		return fmt.Errorf("查询角色失败: %w", err) // 返回数据库查询错误信息
	}
	if count > 0 { // 如果记录数量大于0，说明该角色名称已被使用
		return fmt.Errorf("角色名称 %s 已存在", req.RoleName) // 返回错误信息，表示角色名称已存在
	}

	// 业务规则：角色状态默认设置为启用（ACTIVE），如果请求中提供了状态则使用请求的状态
	status := RoleStatusActive // 初始化角色状态为启用（默认值）
	if req.Status != "" {      // 如果请求中提供了状态（非空），则使用请求的状态
		status = req.Status // 更新角色状态为请求中提供的状态值
	}

	// 业务逻辑：使用数据库事务确保角色创建和权限分配的原子性（要么全部成功，要么全部回滚）
	return db.MysqlDB.Transaction(func(tx *gorm.DB) error { // 开启数据库事务，传入事务处理函数，返回事务执行结果
		role := hotel_admin.Role{ // 创建角色实体对象
			RoleName:    req.RoleName,    // 设置角色名称（从请求中获取）
			Description: req.Description, // 设置角色描述（从请求中获取，可为空）
			Status:      status,          // 设置角色状态（使用计算后的状态值）
		}

		if err := tx.Create(&role).Error; err != nil { // 将角色信息保存到数据库（使用事务连接），如果保存失败则返回错误
			return fmt.Errorf("创建角色失败: %w", err) // 返回错误信息，表示创建角色失败（包含原始错误）
		}

		// 业务逻辑：如果提供了权限ID列表，批量创建角色权限关联记录，建立角色与权限的多对多关系
		if len(req.PermissionIDs) > 0 { // 如果请求中提供了权限ID列表（长度大于0），则批量创建权限关联
			relations := make([]hotel_admin.RolePermissionRelation, len(req.PermissionIDs)) // 创建角色权限关联记录切片，长度为权限ID列表长度
			for i, permissionID := range req.PermissionIDs {                                // 遍历权限ID列表，为每个权限创建关联记录
				relations[i] = hotel_admin.RolePermissionRelation{ // 创建角色权限关联实体对象
					RoleID:       role.ID,      // 设置角色ID（使用刚创建的角色ID）
					PermissionID: permissionID, // 设置权限ID（使用当前遍历到的权限ID）
				}
			}
			if err := tx.Create(&relations).Error; err != nil { // 批量创建角色权限关联记录（使用事务连接），如果保存失败则返回错误
				return fmt.Errorf("分配权限失败: %w", err) // 返回错误信息，表示分配权限失败（包含原始错误）
			}
		}

		return nil // 返回nil表示事务执行成功（角色创建和权限分配都成功）
	})
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(req UpdateRoleReq) error {
	return db.MysqlDB.Transaction(func(tx *gorm.DB) error {
		updates := make(map[string]interface{})

		if req.RoleName != nil {
			var existingRole hotel_admin.Role
			err := tx.Where("role_name = ? AND id != ?", *req.RoleName, req.ID).First(&existingRole).Error
			if err == nil {
				return fmt.Errorf("角色名称 %s 已被使用", *req.RoleName)
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("查询角色失败: %w", err)
			}
			updates["role_name"] = *req.RoleName
		}

		if req.Description != nil {
			updates["description"] = *req.Description
		}

		if req.Status != nil {
			updates["status"] = *req.Status
		}

		if len(updates) > 0 {
			result := tx.Model(&hotel_admin.Role{}).Where("id = ?", req.ID).Updates(updates)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return fmt.Errorf("角色ID %d 不存在", req.ID)
			}
		}

		// 业务逻辑：如果提供了权限ID列表，先删除该角色的所有旧权限关联，再批量创建新的权限关联（实现权限的全量更新）
		if req.PermissionIDs != nil {
			if err := tx.Where("role_id = ?", req.ID).Delete(&hotel_admin.RolePermissionRelation{}).Error; err != nil {
				return fmt.Errorf("删除旧权限失败: %w", err)
			}

			if len(req.PermissionIDs) > 0 {
				relations := make([]hotel_admin.RolePermissionRelation, len(req.PermissionIDs))
				for i, permissionID := range req.PermissionIDs {
					relations[i] = hotel_admin.RolePermissionRelation{
						RoleID:       req.ID,
						PermissionID: permissionID,
					}
				}
				if err := tx.Create(&relations).Error; err != nil {
					return fmt.Errorf("分配新权限失败: %w", err)
				}
			}
		}

		return nil
	})
}

// GetRole 获取角色详情
func (s *RoleService) GetRole(id uint64) (*RoleInfo, error) {
	var role hotel_admin.Role                                                                               // 声明角色实体变量，用于存储查询到的角色信息
	if err := db.MysqlDB.Preload("RolePermissionRelations.Permission").First(&role, id).Error; err != nil { // 通过角色ID查询角色信息（预加载角色权限关联关系和权限详情），如果查询失败则说明角色不存在
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果错误是记录不存在错误
			return nil, fmt.Errorf("角色ID %d 不存在", id) // 返回nil和错误信息，表示角色不存在
		}
		return nil, err // 返回nil和其他数据库查询错误
	}

	info := buildRoleInfo(&role, true) // 调用构建函数，将角色实体对象转换为角色信息对象（包含权限树，第二个参数true表示需要构建权限树）
	return &info, nil                  // 返回角色信息指针和无错误
}

// ListRoles 获取角色列表
func (s *RoleService) ListRoles(req ListRolesReq) (*ListRolesResp, error) {
	req.Page = max(req.Page, 1)                   // 如果页码小于1，则设置为1（使用max函数确保最小值）
	req.PageSize = min(max(req.PageSize, 1), 100) // 如果每页数量小于1则设置为1，如果大于100则设置为100（使用min和max函数确保范围）
	offset := (req.Page - 1) * req.PageSize       // 计算分页偏移量（跳过前面的记录数），用于SQL查询的OFFSET子句

	query := db.MysqlDB.Model(&hotel_admin.Role{}).Where("deleted_at IS NULL") // 创建角色模型的查询构建器，添加软删除筛选条件（只查询未删除的角色）

	if req.Status != nil && *req.Status != "" { // 如果请求中提供了状态（指针非空且值非空），则添加状态筛选条件
		query = query.Where("status = ?", *req.Status) // 添加角色状态筛选条件，只查询指定状态的角色（解引用指针获取值）
	}
	if req.Keyword != nil && *req.Keyword != "" { // 如果请求中提供了关键词（指针非空且值非空），则添加关键词搜索条件
		keyword := "%" + *req.Keyword + "%"                                             // 构建模糊搜索关键词（前后加%通配符）
		query = query.Where("role_name LIKE ? OR description LIKE ?", keyword, keyword) // 添加关键词搜索条件，搜索角色名称或描述包含关键词的角色
	}

	var total int64                                   // 声明总数变量，用于存储符合条件的角色总数
	if err := query.Count(&total).Error; err != nil { // 统计符合条件的角色总数，如果统计失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	var roles []hotel_admin.Role                                                                                 // 声明角色列表变量，用于存储查询到的角色信息列表
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&roles).Error; err != nil { // 按创建时间倒序排列，添加分页限制（偏移量、每页数量）并查询角色列表，如果查询失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	roleInfos := make([]RoleInfo, len(roles)) // 创建角色信息列表，长度为查询到的角色数量
	for i := range roles {                    // 遍历查询到的角色列表
		roleInfos[i] = buildRoleInfo(&roles[i], false) // 调用构建函数，将每个角色实体对象转换为角色信息对象（不构建权限树，第二个参数false表示不需要构建权限树，提高列表查询性能）
	}

	return &ListRolesResp{ // 返回角色列表响应对象
		List:     roleInfos,     // 设置角色列表（转换后的角色信息列表）
		Total:    uint64(total), // 设置总数（转换为uint64类型）
		Page:     req.Page,      // 设置当前页码（从请求中获取）
		PageSize: req.PageSize,  // 设置每页数量（从请求中获取）
	}, nil // 返回响应对象和无错误
}

// DeleteRole 删除角色（软删除）
// 业务功能：逻辑删除角色记录，不物理删除数据，但需检查是否有账号关联到该角色
// 入参说明：
//   - id: 待删除的角色ID
//
// 返回值说明：
//   - error: 角色不存在、角色正在被账号使用（无法删除）或数据库操作错误
func (s *RoleService) DeleteRole(id uint64) error {
	// 业务规则：如果角色被账号使用，不允许删除，避免破坏账号数据的完整性
	var count int64                                                                                                   // 声明计数变量，用于存储使用该角色的账号数量
	if err := db.MysqlDB.Model(&hotel_admin.UserAccount{}).Where("role_id = ?", id).Count(&count).Error; err != nil { // 统计使用该角色的账号数量，如果统计失败则返回错误
		return fmt.Errorf("查询角色使用情况失败: %w", err) // 返回错误信息，表示查询角色使用情况失败
	}

	if count > 0 { // 如果存在使用该角色的账号（数量大于0），则不允许删除
		return fmt.Errorf("该角色正在被 %d 个账号使用，无法删除", count) // 返回错误信息，表示角色正在被使用（包含使用该角色的账号数量）
	}

	result := db.MysqlDB.Delete(&hotel_admin.Role{}, id) // 执行软删除操作（设置deleted_at字段），根据角色ID删除角色记录
	if result.Error != nil {                             // 如果删除操作失败，则返回错误
		return result.Error // 返回数据库操作错误
	}
	if result.RowsAffected == 0 { // 如果删除的记录数为0（没有记录被影响），说明角色不存在
		return fmt.Errorf("角色ID %d 不存在", id) // 返回错误信息，表示角色不存在
	}

	return nil // 返回nil表示删除成功
}
