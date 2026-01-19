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
	PointsChangeTypeEarn    = "EARN"
	PointsChangeTypeConsume = "CONSUME"
)

// MemberPointsService 会员积分管理服务
// 负责处理会员积分的获取、消费、查询等核心业务逻辑，
// 包括积分变动记录管理、积分余额计算、积分余额不足检查、事务保证积分操作的原子性等。
type MemberPointsService struct{}

// CreatePointsRecordReq 创建积分记录请求
type CreatePointsRecordReq struct {
	MemberID     uint64  `json:"member_id" binding:"required"`
	OrderID      *uint64 `json:"order_id,omitempty"`
	ChangeType   string  `json:"change_type" binding:"required"`  // EARN-获取, CONSUME-消费
	PointsValue  int64   `json:"points_value" binding:"required"` // 正数表示获取，负数表示消费
	ChangeReason string  `json:"change_reason" binding:"required"`
	OperatorID   uint64  `json:"operator_id" binding:"required"`
}

// ListPointsRecordsReq 积分记录列表查询请求
type ListPointsRecordsReq struct {
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	MemberID   *uint64 `json:"member_id,omitempty"`
	OrderID    *uint64 `json:"order_id,omitempty"`
	ChangeType *string `json:"change_type,omitempty"`
	StartTime  *string `json:"start_time,omitempty"`
	EndTime    *string `json:"end_time,omitempty"`
}

// PointsRecordInfo 积分记录信息
type PointsRecordInfo struct {
	ID           uint64    `json:"id"`
	MemberID     uint64    `json:"member_id"`
	MemberName   string    `json:"member_name,omitempty"`
	OrderID      *uint64   `json:"order_id,omitempty"`
	ChangeType   string    `json:"change_type"`
	PointsValue  int64     `json:"points_value"`
	ChangeReason string    `json:"change_reason"`
	ChangeTime   time.Time `json:"change_time"`
	OperatorID   uint64    `json:"operator_id"`
}

// ListPointsRecordsResp 积分记录列表响应
type ListPointsRecordsResp struct {
	List     []PointsRecordInfo `json:"list"`
	Total    uint64             `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

// CreatePointsRecord 创建积分记录
// 业务功能：记录会员积分变动（获取或消费），同步更新会员积分余额，用于积分管理和积分兑换
// 入参说明：
//   - req: 创建积分记录请求，包含会员ID、订单ID（可选）、变动类型（EARN/CONSUME）、积分值、变动原因、操作人ID
//
// 返回值说明：
//   - error: 会员不存在、积分变动值为0、积分类型不正确、积分余额不足（消费时）或数据库操作错误
//
// 业务规则：使用数据库事务确保积分记录创建和积分余额更新的原子性
func (s *MemberPointsService) CreatePointsRecord(req CreatePointsRecordReq) error {
	// 业务规则：积分变动值不能为0，确保每次积分操作都有实际意义
	if req.PointsValue == 0 { // 如果积分变动值为0（无效值），则返回错误
		return errors.New("积分变动值不能为0") // 返回错误信息，表示积分变动值不能为0
	}

	// 业务逻辑：使用数据库事务确保积分记录创建和积分余额更新的原子性（要么全部成功，要么全部回滚）
	return db.MysqlDB.Transaction(func(tx *gorm.DB) error { // 开启数据库事务，传入事务处理函数，返回事务执行结果
		// 业务规则：会员必须存在，验证会员是否存在
		var member hotel_admin.Member                                 // 声明会员实体变量，用于存储查询到的会员信息
		if err := tx.First(&member, req.MemberID).Error; err != nil { // 通过会员ID查询会员信息（使用事务连接），如果查询失败则说明会员不存在
			if errors.Is(err, gorm.ErrRecordNotFound) { // 如果错误是记录不存在错误
				return fmt.Errorf("会员ID %d 不存在", req.MemberID) // 返回错误信息，表示会员不存在
			}
			return err // 返回其他数据库查询错误
		}

		// 业务逻辑：根据积分变动类型（获取或消费）校验积分值的正负性，并计算积分变动量
		var pointsChange int64  // 声明积分变动量变量，用于存储计算后的积分变动量
		switch req.ChangeType { // 根据积分变动类型（获取或消费）进行不同的处理
		case PointsChangeTypeEarn: // 如果是获取积分类型
			// 业务规则：获取积分时，积分值必须为正数
			if req.PointsValue < 0 { // 如果积分值小于0（负数），则返回错误
				return errors.New("获取积分值必须为正数") // 返回错误信息，表示获取积分值必须为正数
			}
			pointsChange = req.PointsValue // 获取积分时，积分变动量等于积分值（正数）
		case PointsChangeTypeConsume: // 如果是消费积分类型
			// 业务规则：消费积分时，积分值必须为负数
			if req.PointsValue > 0 { // 如果积分值大于0（正数），则返回错误
				return errors.New("消费积分值必须为负数") // 返回错误信息，表示消费积分值必须为负数
			}
			pointsChange = req.PointsValue // 消费积分时，积分变动量等于积分值（负数）
			// 业务规则：消费积分时，必须检查积分余额是否充足
			if member.PointsBalance < uint64(-pointsChange) { // 如果会员积分余额小于需要消费的积分数量（pointsChange是负数，需要取绝对值）
				return fmt.Errorf("积分余额不足，当前余额: %d，需要: %d", member.PointsBalance, -pointsChange) // 返回错误信息，表示积分余额不足，显示当前余额和需要的积分
			}
		default: // 如果是其他无效的积分变动类型
			return fmt.Errorf("无效的积分变动类型: %s", req.ChangeType) // 返回错误信息，表示无效的积分变动类型
		}

		// 业务计算：计算新的积分余额，确保余额不为负数（双重检查）
		newBalance := int64(member.PointsBalance) + pointsChange // 计算新的积分余额（当前余额 + 积分变动量）
		if newBalance < 0 {                                      // 如果新余额小于0（负数），则返回错误
			return errors.New("积分余额不能为负数") // 返回错误信息，表示积分余额不能为负数（双重检查）
		}

		record := hotel_admin.MemberPointsRecord{ // 创建积分记录实体对象
			MemberID:     req.MemberID,     // 设置会员ID（从请求中获取）
			OrderID:      req.OrderID,      // 设置订单ID（从请求中获取，可为空）
			ChangeType:   req.ChangeType,   // 设置变动类型（从请求中获取：EARN-获取或CONSUME-消费）
			ChangePoints: pointsChange,     // 设置变动积分（计算后的积分变动量）
			ChangeReason: req.ChangeReason, // 设置变动原因（从请求中获取）
			ChangeTime:   time.Now(),       // 设置变动时间（当前时间）
			OperatorID:   req.OperatorID,   // 设置操作人ID（从请求中获取）
		}

		if err := tx.Create(&record).Error; err != nil { // 将积分记录保存到数据库（使用事务连接），如果保存失败则返回错误
			return fmt.Errorf("创建积分记录失败: %w", err) // 返回错误信息，表示创建积分记录失败（包含原始错误）
		}

		// 业务逻辑：同步更新会员的积分余额，确保积分记录和积分余额的一致性
		if err := tx.Model(&member).Update("points_balance", newBalance).Error; err != nil { // 更新会员的积分余额（使用事务连接），如果更新失败则返回错误
			return fmt.Errorf("更新积分余额失败: %w", err) // 返回错误信息，表示更新积分余额失败（包含原始错误）
		}

		return nil // 返回nil表示事务执行成功（积分记录创建和积分余额更新都成功）
	})
}

// buildPointsRecordInfo 构建积分记录信息（避免代码重复）
func buildPointsRecordInfo(record *hotel_admin.MemberPointsRecord) PointsRecordInfo {
	info := PointsRecordInfo{ // 创建积分记录信息对象
		ID:           record.ID,           // 设置记录ID（从积分记录实体中获取）
		MemberID:     record.MemberID,     // 设置会员ID（从积分记录实体中获取）
		OrderID:      record.OrderID,      // 设置订单ID（从积分记录实体中获取，可为空）
		ChangeType:   record.ChangeType,   // 设置变动类型（从积分记录实体中获取：EARN-获取或CONSUME-消费）
		PointsValue:  record.ChangePoints, // 设置积分值（从积分记录实体中获取，变动积分）
		ChangeReason: record.ChangeReason, // 设置变动原因（从积分记录实体中获取）
		ChangeTime:   record.ChangeTime,   // 设置变动时间（从积分记录实体中获取）
		OperatorID:   record.OperatorID,   // 设置操作人ID（从积分记录实体中获取）
	}

	if record.Member != nil && record.Member.Guest != nil { // 如果积分记录关联了会员信息且会员关联了客人信息（预加载的数据）
		info.MemberName = record.Member.Guest.Name // 设置会员姓名（从关联的客人信息中获取）
	}

	return info
}

// ListPointsRecords 获取积分记录列表
func (s *MemberPointsService) ListPointsRecords(req ListPointsRecordsReq) (*ListPointsRecordsResp, error) {
	req.Page = max(req.Page, 1)                   // 如果页码小于1，则设置为1（使用max函数确保最小值）
	req.PageSize = min(max(req.PageSize, 1), 100) // 如果每页数量小于1则设置为1，如果大于100则设置为100（使用min和max函数确保范围）
	offset := (req.Page - 1) * req.PageSize       // 计算分页偏移量（跳过前面的记录数），用于SQL查询的OFFSET子句

	query := db.MysqlDB.Model(&hotel_admin.MemberPointsRecord{}).Where("deleted_at IS NULL") // 创建积分记录模型的查询构建器，添加软删除筛选条件（只查询未删除的积分记录）

	if req.MemberID != nil { // 如果请求中提供了会员ID（指针非空），则添加会员筛选条件
		query = query.Where("member_id = ?", *req.MemberID) // 添加会员ID筛选条件，只查询指定会员的积分记录（解引用指针获取值）
	}
	if req.OrderID != nil { // 如果请求中提供了订单ID（指针非空），则添加订单筛选条件
		query = query.Where("order_id = ?", *req.OrderID) // 添加订单ID筛选条件，只查询指定订单的积分记录（解引用指针获取值）
	}
	if req.ChangeType != nil && *req.ChangeType != "" { // 如果请求中提供了变动类型（指针非空且值非空），则添加变动类型筛选条件
		query = query.Where("change_type = ?", *req.ChangeType) // 添加变动类型筛选条件，只查询指定变动类型的积分记录（解引用指针获取值）
	}
	if req.StartTime != nil && *req.StartTime != "" { // 如果请求中提供了开始时间（指针非空且值非空），则添加开始时间筛选条件
		query = query.Where("change_time >= ?", *req.StartTime) // 添加开始时间筛选条件，只查询变动时间大于等于开始时间的积分记录（解引用指针获取值）
	}
	if req.EndTime != nil && *req.EndTime != "" { // 如果请求中提供了结束时间（指针非空且值非空），则添加结束时间筛选条件
		query = query.Where("change_time <= ?", *req.EndTime) // 添加结束时间筛选条件，只查询变动时间小于等于结束时间的积分记录（解引用指针获取值）
	}

	var total int64                                   // 声明总数变量，用于存储符合条件的积分记录总数
	if err := query.Count(&total).Error; err != nil { // 统计符合条件的积分记录总数，如果统计失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	var records []hotel_admin.MemberPointsRecord                       // 声明积分记录列表变量，用于存储查询到的积分记录信息列表
	if err := query.Preload("Member.Guest").Order("change_time DESC"). // 预加载会员信息和客人信息关联数据（嵌套预加载），按变动时间倒序排列
										Offset(offset).Limit(req.PageSize).Find(&records).Error; err != nil { // 添加分页限制（偏移量、每页数量）并查询积分记录列表，如果查询失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	recordInfos := make([]PointsRecordInfo, len(records)) // 创建积分记录信息列表，长度为查询到的积分记录数量
	for i := range records {                              // 遍历查询到的积分记录列表
		recordInfos[i] = buildPointsRecordInfo(&records[i]) // 调用构建函数，将每个积分记录实体对象转换为积分记录信息对象（包含会员姓名）
	}

	return &ListPointsRecordsResp{ // 返回积分记录列表响应对象
		List:     recordInfos,   // 设置积分记录列表（转换后的积分记录信息列表）
		Total:    uint64(total), // 设置总数（转换为uint64类型）
		Page:     req.Page,      // 设置当前页码（从请求中获取）
		PageSize: req.PageSize,  // 设置每页数量（从请求中获取）
	}, nil // 返回响应对象和无错误
}

// GetMemberPointsBalance 获取会员积分余额
func (s *MemberPointsService) GetMemberPointsBalance(memberID uint64) (uint64, error) {
	var member hotel_admin.Member                                     // 声明会员实体变量，用于存储查询到的会员信息
	if err := db.MysqlDB.First(&member, memberID).Error; err != nil { // 通过会员ID查询会员信息，如果查询失败则说明会员不存在
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果错误是记录不存在错误
			return 0, fmt.Errorf("会员ID %d 不存在", memberID) // 返回0和错误信息，表示会员不存在
		}
		return 0, err // 返回0和其他数据库查询错误
	}
	return member.PointsBalance, nil // 返回会员积分余额和无错误
}
