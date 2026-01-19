package service

import (
	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"
)

const (
	PermissionStatusActive   = "ACTIVE"
	PermissionStatusInactive = "INACTIVE"
)

// PermissionService 权限管理服务
// 负责处理后台权限的查询等核心业务逻辑，
// 包括权限树形结构构建、权限与角色的关联关系管理等。
type PermissionService struct{}

// ListPermissionsReq 权限列表查询请求
type ListPermissionsReq struct {
	PermissionType *string `json:"permission_type,omitempty"`
	ParentID       *uint64 `json:"parent_id,omitempty"`
	Status         *string `json:"status,omitempty"`
}

// PermissionInfo 权限信息
type PermissionInfo struct {
	ID             uint64           `json:"id"`
	PermissionName string           `json:"permission_name"`
	PermissionURL  string           `json:"permission_url"`
	PermissionType string           `json:"permission_type"`
	ParentID       *uint64          `json:"parent_id,omitempty"`
	Status         string           `json:"status"`
	Children       []PermissionInfo `json:"children,omitempty"`
}

// ListPermissionsResp 权限列表响应
type ListPermissionsResp struct {
	List []PermissionInfo `json:"list"`
}

// buildPermissionInfo 构建权限信息（避免代码重复）
func buildPermissionInfo(perm *hotel_admin.Permission) PermissionInfo {
	return PermissionInfo{
		ID:             perm.ID,
		PermissionName: perm.PermissionName,
		PermissionURL:  perm.PermissionURL,
		PermissionType: perm.PermissionType,
		ParentID:       perm.ParentID,
		Status:         perm.Status,
	}
}

// ListPermissions 获取权限列表（树形结构）
// 业务功能：查询权限列表并以树形结构返回，用于权限管理和权限分配场景
// 入参说明：
//   - req: 权限列表查询请求，支持按权限类型、父权限ID、状态筛选
//
// 返回值说明：
//   - *ListPermissionsResp: 权限列表（树形结构，包含子权限）
//   - error: 查询失败错误
//
// 业务规则：如果未指定父权限ID，则返回根节点权限（parent_id为NULL），否则返回指定父权限的子权限
func (s *PermissionService) ListPermissions(req ListPermissionsReq) (*ListPermissionsResp, error) {
	query := db.MysqlDB.Model(&hotel_admin.Permission{}).Where("deleted_at IS NULL") // 创建权限模型的查询构建器，添加软删除筛选条件（只查询未删除的权限）

	// 业务筛选：按权限类型筛选（如菜单权限、按钮权限等）
	if req.PermissionType != nil && *req.PermissionType != "" { // 如果请求中提供了权限类型（指针非空且值非空），则添加权限类型筛选条件
		query = query.Where("permission_type = ?", *req.PermissionType) // 添加权限类型筛选条件，只查询指定类型的权限（解引用指针获取值）
	}
	// 业务筛选：按权限状态筛选（如启用/停用）
	if req.Status != nil && *req.Status != "" { // 如果请求中提供了状态（指针非空且值非空），则添加状态筛选条件
		query = query.Where("status = ?", *req.Status) // 添加权限状态筛选条件，只查询指定状态的权限（解引用指针获取值）
	}

	// 业务逻辑：根据父权限ID筛选，如果未指定则只返回根节点权限，否则返回指定父权限的子权限
	if req.ParentID != nil { // 如果请求中提供了父权限ID（指针非空）
		query = query.Where("parent_id = ?", *req.ParentID) // 添加父权限ID筛选条件，只查询指定父权限的子权限（解引用指针获取值）
	} else {
		query = query.Where("parent_id IS NULL") // 添加父权限ID为空筛选条件，只查询根节点权限（没有父权限的权限）
	}

	var parentPermissions []hotel_admin.Permission                               // 声明父权限列表变量，用于存储查询到的父权限信息列表
	if err := query.Order("id ASC").Find(&parentPermissions).Error; err != nil { // 按ID正序排列并查询父权限列表，如果查询失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	if len(parentPermissions) == 0 { // 如果父权限列表为空（没有父权限）
		return &ListPermissionsResp{List: []PermissionInfo{}}, nil // 返回空列表响应对象和无错误
	}

	parentIDs := make([]uint64, len(parentPermissions)) // 创建父权限ID列表，长度为父权限数量
	for i, p := range parentPermissions {               // 遍历父权限列表，提取所有父权限ID
		parentIDs[i] = p.ID // 将父权限ID添加到列表中
	}

	var childPermissions []hotel_admin.Permission                                      // 声明子权限列表变量，用于存储查询到的子权限信息列表
	childQuery := db.MysqlDB.Where("parent_id IN ? AND deleted_at IS NULL", parentIDs) // 创建子权限查询构建器，查询所有父权限的子权限（使用IN子句）
	if req.Status != nil && *req.Status != "" {                                        // 如果请求中提供了状态（指针非空且值非空），则添加状态筛选条件
		childQuery = childQuery.Where("status = ?", *req.Status) // 添加权限状态筛选条件，只查询指定状态的子权限（解引用指针获取值）
	}
	if err := childQuery.Order("id ASC").Find(&childPermissions).Error; err != nil { // 按ID正序排列并查询子权限列表，如果查询失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	childrenMap := make(map[uint64][]PermissionInfo) // 创建子权限映射表，用于按父权限ID分组存储子权限信息
	for i := range childPermissions {                // 遍历子权限列表，构建子权限映射表
		child := buildPermissionInfo(&childPermissions[i]) // 调用构建函数，将子权限实体对象转换为权限信息对象
		if childPermissions[i].ParentID != nil {           // 如果子权限有父权限ID（指针非空）
			childrenMap[*childPermissions[i].ParentID] = append(childrenMap[*childPermissions[i].ParentID], child) // 将子权限信息添加到对应父权限ID的子权限列表中（使用映射表分组）
		}
	}

	permissionInfos := make([]PermissionInfo, len(parentPermissions)) // 创建权限信息列表，长度为父权限数量
	for i := range parentPermissions {                                // 遍历父权限列表，构建权限树形结构
		permissionInfos[i] = buildPermissionInfo(&parentPermissions[i]) // 调用构建函数，将父权限实体对象转换为权限信息对象
		if children, ok := childrenMap[parentPermissions[i].ID]; ok {   // 如果父权限ID在子权限映射表中存在（有子权限）
			permissionInfos[i].Children = children // 设置父权限的子权限列表（构建树形结构）
		}
	}

	return &ListPermissionsResp{List: permissionInfos}, nil // 返回权限列表响应对象（包含树形结构的权限信息）和无错误
}

// GetAllPermissions 获取所有权限（扁平列表）
func (s *PermissionService) GetAllPermissions() ([]PermissionInfo, error) {
	var permissions []hotel_admin.Permission                                                 // 声明权限列表变量，用于存储查询到的权限信息列表
	if err := db.MysqlDB.Where("deleted_at IS NULL AND status = ?", PermissionStatusActive). // 创建权限查询构建器，筛选未删除且状态为启用的权限
													Order("id ASC").Find(&permissions).Error; err != nil { // 按ID正序排列并查询权限列表，如果查询失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	permissionInfos := make([]PermissionInfo, len(permissions)) // 创建权限信息列表，长度为查询到的权限数量
	for i := range permissions {                                // 遍历查询到的权限列表
		permissionInfos[i] = buildPermissionInfo(&permissions[i]) // 调用构建函数，将每个权限实体对象转换为权限信息对象
	}

	return permissionInfos, nil // 返回权限信息列表（扁平结构）和无错误
}
