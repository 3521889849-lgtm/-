package service

import (
	"errors"
	"fmt"
	"time"

	"example_shop/common/db"
	hotel_admin "example_shop/common/model/hotel_admin"

	"gorm.io/gorm"
)

const (
	MemberStatusActive   = "ACTIVE"
	MemberStatusInactive = "INACTIVE"
)

// MemberService 会员管理服务
// 负责处理酒店会员的创建、更新、查询、删除等核心业务逻辑，
// 包括会员与客人信息的关联校验、会员等级管理、积分余额管理、会员状态管理等。
type MemberService struct{}

// CreateMemberReq 创建会员请求
type CreateMemberReq struct {
	GuestID       uint64 `json:"guest_id" binding:"required"`
	MemberLevel   string `json:"member_level" binding:"required"`
	PointsBalance uint64 `json:"points_balance"`
	Status        string `json:"status"`
}

// UpdateMemberReq 更新会员请求
type UpdateMemberReq struct {
	ID            uint64  `json:"id" binding:"required"`
	MemberLevel   *string `json:"member_level,omitempty"`
	PointsBalance *uint64 `json:"points_balance,omitempty"`
	Status        *string `json:"status,omitempty"`
}

// ListMembersReq 会员列表查询请求
type ListMembersReq struct {
	Page        int     `json:"page"`
	PageSize    int     `json:"page_size"`
	MemberLevel *string `json:"member_level,omitempty"`
	Status      *string `json:"status,omitempty"`
	Keyword     *string `json:"keyword,omitempty"`
}

// MemberInfo 会员信息
type MemberInfo struct {
	ID              uint64     `json:"id"`
	GuestID         uint64     `json:"guest_id"`
	GuestName       string     `json:"guest_name,omitempty"`
	GuestPhone      string     `json:"guest_phone,omitempty"`
	MemberLevel     string     `json:"member_level"`
	PointsBalance   uint64     `json:"points_balance"`
	RegisterTime    time.Time  `json:"register_time"`
	LastCheckInTime *time.Time `json:"last_check_in_time,omitempty"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
}

// ListMembersResp 会员列表响应
type ListMembersResp struct {
	List     []MemberInfo `json:"list"`
	Total    uint64       `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
}

// buildMemberInfo 构建会员信息（避免代码重复）
func buildMemberInfo(member *hotel_admin.Member) MemberInfo {
	info := MemberInfo{
		ID:            member.ID,
		GuestID:       member.GuestID,
		MemberLevel:   member.MemberLevel,
		PointsBalance: member.PointsBalance,
		RegisterTime:  member.RegisterTime,
		Status:        member.Status,
		CreatedAt:     member.CreatedAt,
	}

	if member.Guest != nil {
		info.GuestName = member.Guest.Name
		info.GuestPhone = member.Guest.Phone
	}
	if member.LastCheckInTime != nil {
		info.LastCheckInTime = member.LastCheckInTime
	}

	return info
}

// CreateMember 创建会员
// 业务功能：将客人注册为会员，建立客人与会员的关联关系，初始化会员等级和积分余额
// 入参说明：
//   - req: 创建会员请求，包含客人ID、会员等级、积分余额（可选，默认0）、状态（可选，默认ACTIVE）
//
// 返回值说明：
//   - error: 客人ID不存在、客人已是会员、业务校验失败或数据库操作错误
func (s *MemberService) CreateMember(req CreateMemberReq) error {
	// 业务规则：一个客人只能注册一次会员，检查该客人是否已存在会员记录
	var existingMember hotel_admin.Member                                             // 声明会员实体变量，用于存储查询到的已存在会员信息
	err := db.MysqlDB.Where("guest_id = ?", req.GuestID).First(&existingMember).Error // 通过客人ID查询是否已存在会员记录，如果查询失败则说明该客人未注册会员
	if err == nil {                                                                   // 如果查询成功（未报错），说明该客人已是会员
		return fmt.Errorf("客人ID %d 已是会员", req.GuestID) // 返回错误信息，表示客人已是会员
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) { // 如果查询失败但不是记录不存在错误，说明是其他数据库错误
		return fmt.Errorf("查询会员失败: %w", err) // 返回数据库查询错误信息
	}

	// 业务规则：会员必须关联到有效的客人信息，验证客人是否存在
	var guest hotel_admin.GuestInfo                                     // 声明客人实体变量，用于存储查询到的客人信息
	if err := db.MysqlDB.First(&guest, req.GuestID).Error; err != nil { // 通过客人ID查询客人信息，如果查询失败则说明客人不存在
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果错误是记录不存在错误
			return fmt.Errorf("客人ID %d 不存在", req.GuestID) // 返回错误信息，表示客人不存在
		}
		return fmt.Errorf("查询客人失败: %w", err) // 返回数据库查询错误信息
	}

	// 业务规则：会员状态默认设置为启用（ACTIVE），如果请求中提供了状态则使用请求的状态
	status := MemberStatusActive // 初始化会员状态为启用（默认值）
	if req.Status != "" {        // 如果请求中提供了状态（非空），则使用请求的状态
		status = req.Status // 更新会员状态为请求中提供的状态值
	}

	member := hotel_admin.Member{ // 创建会员实体对象
		GuestID:       req.GuestID,       // 设置客人ID（从请求中获取）
		MemberLevel:   req.MemberLevel,   // 设置会员等级（从请求中获取）
		PointsBalance: req.PointsBalance, // 设置积分余额（从请求中获取，默认0）
		Status:        status,            // 设置会员状态（使用计算后的状态值）
	}

	return db.MysqlDB.Create(&member).Error // 将会员信息保存到数据库，返回保存操作的结果（成功为nil，失败为error）
}

// UpdateMember 更新会员
// 业务功能：修改已存在会员的属性信息，支持部分字段更新（会员等级、积分余额、状态）
// 入参说明：
//   - req: 会员更新请求，所有字段均为可选，只更新传入的非空字段
//
// 返回值说明：
//   - error: 会员不存在或数据库操作错误
func (s *MemberService) UpdateMember(req UpdateMemberReq) error {
	updates := make(map[string]interface{}) // 创建更新字段映射表，用于存储需要更新的字段和值

	// 业务逻辑：采用部分更新策略，只更新请求中提供的非空字段
	if req.MemberLevel != nil { // 如果请求中提供了会员等级（指针非空），则添加到更新映射表
		updates["member_level"] = *req.MemberLevel // 添加会员等级到更新映射表（解引用指针获取值）
	}
	if req.PointsBalance != nil { // 如果请求中提供了积分余额（指针非空），则添加到更新映射表
		updates["points_balance"] = *req.PointsBalance // 添加积分余额到更新映射表（解引用指针获取值）
	}
	if req.Status != nil { // 如果请求中提供了状态（指针非空），则添加到更新映射表
		updates["status"] = *req.Status // 添加状态到更新映射表（解引用指针获取值）
	}

	// 业务规则：如果没有需要更新的字段，直接返回，避免无效的数据库操作
	if len(updates) == 0 { // 如果更新映射表为空（没有需要更新的字段）
		return nil // 直接返回nil，表示更新成功（实际上没有更新任何字段）
	}

	// 执行更新，并检查是否有记录被影响（用于判断会员是否存在）
	result := db.MysqlDB.Model(&hotel_admin.Member{}).Where("id = ?", req.ID).Updates(updates) // 根据会员ID更新会员信息，使用更新映射表中的字段和值
	if result.Error != nil {                                                                   // 如果更新操作失败，则返回错误
		return result.Error // 返回数据库操作错误
	}
	if result.RowsAffected == 0 { // 如果更新的记录数为0（没有记录被影响），说明会员不存在
		return fmt.Errorf("会员ID %d 不存在", req.ID) // 返回错误信息，表示会员不存在
	}

	return nil // 返回nil表示更新成功
}

// GetMember 获取会员详情
// 业务功能：根据会员ID查询会员的完整信息，包含关联的客人信息
// 入参说明：
//   - id: 会员ID
//
// 返回值说明：
//   - *MemberInfo: 会员完整信息（包含客人姓名、手机号等关联数据）
//   - error: 会员不存在或查询失败
func (s *MemberService) GetMember(id uint64) (*MemberInfo, error) {
	var member hotel_admin.Member                                                // 声明会员实体变量，用于存储查询到的会员信息
	if err := db.MysqlDB.Preload("Guest").First(&member, id).Error; err != nil { // 通过会员ID查询会员信息（预加载客人信息），如果查询失败则说明会员不存在
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果错误是记录不存在错误
			return nil, fmt.Errorf("会员ID %d 不存在", id) // 返回nil和错误信息，表示会员不存在
		}
		return nil, err // 返回nil和其他数据库查询错误
	}

	info := buildMemberInfo(&member) // 调用构建函数，将会员实体对象转换为会员信息对象（包含客人信息）
	return &info, nil                // 返回会员信息指针和无错误
}

// ListMembers 获取会员列表
// 业务功能：支持多条件筛选和分页查询会员列表，用于会员管理和统计分析场景
// 入参说明：
//   - req: 会员列表查询请求，支持按会员等级、状态筛选，支持关键词搜索（会员姓名/手机号），支持分页
//
// 返回值说明：
//   - *ListMembersResp: 符合条件的会员列表（包含客人信息）及分页信息
//   - error: 查询失败错误
func (s *MemberService) ListMembers(req ListMembersReq) (*ListMembersResp, error) {
	// 业务规则：分页参数默认值设置，页码最小为1，每页数量最小1条，最大不超过100条
	req.Page = max(req.Page, 1)                   // 如果页码小于1，则设置为1（使用max函数确保最小值）
	req.PageSize = min(max(req.PageSize, 1), 100) // 如果每页数量小于1则设置为1，如果大于100则设置为100（使用min和max函数确保范围）
	offset := (req.Page - 1) * req.PageSize       // 计算分页偏移量（跳过前面的记录数），用于SQL查询的OFFSET子句

	query := db.MysqlDB.Model(&hotel_admin.Member{}).Where("member.deleted_at IS NULL") // 创建会员模型的查询构建器，添加软删除筛选条件（只查询未删除的会员）

	// 业务筛选：按会员等级筛选（如普通/黄金/钻石会员），支持等级维度统计
	if req.MemberLevel != nil && *req.MemberLevel != "" { // 如果请求中提供了会员等级（指针非空且值非空），则添加等级筛选条件
		query = query.Where("member.member_level = ?", *req.MemberLevel) // 添加会员等级筛选条件，只查询指定等级的会员（解引用指针获取值）
	}
	// 业务筛选：按会员状态筛选（如启用/停用），支持状态维度管理
	if req.Status != nil && *req.Status != "" { // 如果请求中提供了状态（指针非空且值非空），则添加状态筛选条件
		query = query.Where("member.status = ?", *req.Status) // 添加会员状态筛选条件，只查询指定状态的会员（解引用指针获取值）
	}
	// 业务搜索：关键词多字段模糊搜索，通过JOIN客人表搜索会员姓名或手机号
	if req.Keyword != nil && *req.Keyword != "" { // 如果请求中提供了关键词（指针非空且值非空），则添加关键词搜索条件
		keyword := "%" + *req.Keyword + "%"                                        // 构建模糊搜索关键词（前后加%通配符）
		query = query.Joins("JOIN guest_info ON member.guest_id = guest_info.id"). // 通过JOIN客人表关联查询（内连接）
												Where("guest_info.name LIKE ? OR guest_info.phone LIKE ?", keyword, keyword) // 添加关键词搜索条件，搜索客人姓名或手机号包含关键词的会员
	}

	var total int64                                   // 声明总数变量，用于存储符合条件的会员总数
	if err := query.Count(&total).Error; err != nil { // 统计符合条件的会员总数，如果统计失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	var members []hotel_admin.Member                                  // 声明会员列表变量，用于存储查询到的会员信息列表
	if err := query.Preload("Guest").Order("member.created_at DESC"). // 预加载客人信息关联数据，按创建时间倒序排列
										Offset(offset).Limit(req.PageSize).Find(&members).Error; err != nil { // 添加分页限制（偏移量、每页数量）并查询会员列表，如果查询失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	memberInfos := make([]MemberInfo, len(members)) // 创建会员信息列表，长度为查询到的会员数量
	for i := range members {                        // 遍历查询到的会员列表
		memberInfos[i] = buildMemberInfo(&members[i]) // 调用构建函数，将每个会员实体对象转换为会员信息对象（包含客人信息）
	}

	return &ListMembersResp{ // 返回会员列表响应对象
		List:     memberInfos,   // 设置会员列表（转换后的会员信息列表）
		Total:    uint64(total), // 设置总数（转换为uint64类型）
		Page:     req.Page,      // 设置当前页码（从请求中获取）
		PageSize: req.PageSize,  // 设置每页数量（从请求中获取）
	}, nil // 返回响应对象和无错误
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// DeleteMember 删除会员（软删除）
// 业务功能：逻辑删除会员记录，不物理删除数据，保留历史积分和权益记录
// 入参说明：
//   - id: 待删除的会员ID
//
// 返回值说明：
//   - error: 会员不存在或数据库操作错误
func (s *MemberService) DeleteMember(id uint64) error {
	return db.MysqlDB.Delete(&hotel_admin.Member{}, id).Error // 执行软删除操作（设置deleted_at字段），根据会员ID删除会员记录，返回删除操作的结果（成功为nil，失败为error）
}

// GetMemberByGuestID 根据客人ID获取会员信息
// 业务功能：通过客人ID查询该客人的会员信息，用于客人注册会员时的重复检查或会员信息展示
// 入参说明：
//   - guestID: 客人ID
//
// 返回值说明：
//   - *MemberInfo: 会员完整信息（包含客人信息）
//   - error: 客人无会员记录或查询失败
func (s *MemberService) GetMemberByGuestID(guestID uint64) (*MemberInfo, error) {
	var member hotel_admin.Member                                                                           // 声明会员实体变量，用于存储查询到的会员信息
	if err := db.MysqlDB.Preload("Guest").Where("guest_id = ?", guestID).First(&member).Error; err != nil { // 通过客人ID查询会员信息（预加载客人信息），如果查询失败则说明该客人无会员记录
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果错误是记录不存在错误
			return nil, fmt.Errorf("客人ID %d 无会员记录", guestID) // 返回nil和错误信息，表示客人无会员记录
		}
		return nil, err // 返回nil和其他数据库查询错误
	}

	info := buildMemberInfo(&member) // 调用构建函数，将会员实体对象转换为会员信息对象（包含客人信息）
	return &info, nil                // 返回会员信息指针和无错误
}
