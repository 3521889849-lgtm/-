package service

import (
	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"
	"time"
)

// OperationLogService 操作日志管理服务
// 负责处理系统操作日志的记录、查询等核心业务逻辑，
// 包括操作日志的记录（异步）、操作日志的查询和统计分析等。
type OperationLogService struct{}

// CreateOperationLogReq 创建操作日志请求
type CreateOperationLogReq struct {
	OperatorID    uint64  `json:"operator_id"`          // 操作人ID
	Module        string  `json:"module"`               // 操作模块
	OperationType string  `json:"operation_type"`       // 操作类型
	Content       string  `json:"content"`              // 操作内容
	OperationIP   string  `json:"operation_ip"`         // 操作IP
	RelatedID     *uint64 `json:"related_id,omitempty"` // 关联ID
	IsSuccess     bool    `json:"is_success"`           // 是否成功
}

// ListOperationLogsReq 查询操作日志请求
type ListOperationLogsReq struct {
	Page          int     `json:"page"`
	PageSize      int     `json:"page_size"`
	OperatorID    *uint64 `json:"operator_id,omitempty"`
	Module        *string `json:"module,omitempty"`
	OperationType *string `json:"operation_type,omitempty"`
	StartTime     *string `json:"start_time,omitempty"`
	EndTime       *string `json:"end_time,omitempty"`
	IsSuccess     *bool   `json:"is_success,omitempty"`
}

// OperationLogInfo 操作日志信息
type OperationLogInfo struct {
	ID            uint64    `json:"id"`
	OperatorID    uint64    `json:"operator_id"`
	OperatorName  string    `json:"operator_name,omitempty"`
	Module        string    `json:"module"`
	OperationType string    `json:"operation_type"`
	Content       string    `json:"content"`
	OperationTime time.Time `json:"operation_time"`
	OperationIP   string    `json:"operation_ip"`
	RelatedID     *uint64   `json:"related_id,omitempty"`
	IsSuccess     bool      `json:"is_success"`
	CreatedAt     time.Time `json:"created_at"`
}

// ListOperationLogsResp 操作日志列表响应
type ListOperationLogsResp struct {
	List     []OperationLogInfo `json:"list"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

// CreateOperationLog 创建操作日志
func (s *OperationLogService) CreateOperationLog(req CreateOperationLogReq) error {
	log := hotel_admin.OperationLog{ // 创建操作日志实体对象
		OperatorID:    req.OperatorID,    // 设置操作人ID（从请求中获取）
		Module:        req.Module,        // 设置操作模块（从请求中获取）
		OperationType: req.OperationType, // 设置操作类型（从请求中获取）
		Content:       req.Content,       // 设置操作内容（从请求中获取）
		OperationIP:   req.OperationIP,   // 设置操作IP（从请求中获取）
		RelatedID:     req.RelatedID,     // 设置关联ID（从请求中获取，可为空）
		IsSuccess:     req.IsSuccess,     // 设置是否成功（从请求中获取）
		OperationTime: time.Now(),        // 设置操作时间为当前时间（自动生成）
	}

	return db.MysqlDB.Create(&log).Error // 将操作日志保存到数据库，返回保存操作的结果（成功为nil，失败为error）
}

// ListOperationLogs 查询操作日志列表
func (s *OperationLogService) ListOperationLogs(req ListOperationLogsReq) (*ListOperationLogsResp, error) {
	var logs []hotel_admin.OperationLog // 声明操作日志列表变量，用于存储查询到的操作日志信息列表
	var total int64                     // 声明总数变量，用于存储符合条件的操作日志总数

	query := db.MysqlDB.Model(&hotel_admin.OperationLog{}) // 创建操作日志模型的查询构建器

	// 条件过滤
	if req.OperatorID != nil { // 如果请求中提供了操作人ID（指针非空），则添加操作人筛选条件
		query = query.Where("operator_id = ?", *req.OperatorID) // 添加操作人ID筛选条件，只查询指定操作人的日志（解引用指针获取值）
	}
	if req.Module != nil && *req.Module != "" { // 如果请求中提供了模块（指针非空且值非空），则添加模块筛选条件
		query = query.Where("module = ?", *req.Module) // 添加模块筛选条件，只查询指定模块的日志（解引用指针获取值）
	}
	if req.OperationType != nil && *req.OperationType != "" { // 如果请求中提供了操作类型（指针非空且值非空），则添加操作类型筛选条件
		query = query.Where("operation_type = ?", *req.OperationType) // 添加操作类型筛选条件，只查询指定操作类型的日志（解引用指针获取值）
	}
	if req.StartTime != nil && *req.StartTime != "" { // 如果请求中提供了开始时间（指针非空且值非空），则添加开始时间筛选条件
		query = query.Where("operation_time >= ?", *req.StartTime) // 添加开始时间筛选条件，只查询操作时间大于等于开始时间的日志（解引用指针获取值）
	}
	if req.EndTime != nil && *req.EndTime != "" { // 如果请求中提供了结束时间（指针非空且值非空），则添加结束时间筛选条件
		query = query.Where("operation_time <= ?", *req.EndTime) // 添加结束时间筛选条件，只查询操作时间小于等于结束时间的日志（解引用指针获取值）
	}
	if req.IsSuccess != nil { // 如果请求中提供了是否成功（指针非空），则添加成功状态筛选条件
		query = query.Where("is_success = ?", *req.IsSuccess) // 添加成功状态筛选条件，只查询指定成功状态的日志（解引用指针获取值）
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil { // 统计符合条件的操作日志总数，如果统计失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize // 计算分页偏移量（跳过前面的记录数），用于SQL查询的OFFSET子句
	if err := query.                        // 继续构建查询
						Preload("Operator").            // 预加载操作人信息关联数据（JOIN查询操作人信息）
						Order("operation_time DESC").   // 添加排序条件，按操作时间倒序排列（最新操作的日志排在前面）
						Offset(offset).                 // 添加分页偏移量
						Limit(req.PageSize).            // 添加每页数量限制
						Find(&logs).Error; err != nil { // 执行查询并获取符合条件的操作日志列表（包含所有预加载的关联数据），如果查询失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	// 转换为响应格式
	list := make([]OperationLogInfo, len(logs)) // 创建操作日志信息列表，长度为查询到的操作日志数量
	for i, log := range logs {                  // 遍历查询到的操作日志列表
		operatorName := ""       // 声明操作人姓名变量，初始化为空字符串
		if log.Operator != nil { // 如果操作日志关联了操作人信息（预加载的数据）
			operatorName = log.Operator.Username // 设置操作人姓名（从关联的操作人信息中获取用户名）
		}
		list[i] = OperationLogInfo{ // 创建操作日志信息对象
			ID:            log.ID,            // 设置日志ID（从操作日志实体中获取）
			OperatorID:    log.OperatorID,    // 设置操作人ID（从操作日志实体中获取）
			OperatorName:  operatorName,      // 设置操作人姓名（从关联的操作人信息中获取）
			Module:        log.Module,        // 设置操作模块（从操作日志实体中获取）
			OperationType: log.OperationType, // 设置操作类型（从操作日志实体中获取）
			Content:       log.Content,       // 设置操作内容（从操作日志实体中获取）
			OperationTime: log.OperationTime, // 设置操作时间（从操作日志实体中获取）
			OperationIP:   log.OperationIP,   // 设置操作IP（从操作日志实体中获取）
			RelatedID:     log.RelatedID,     // 设置关联ID（从操作日志实体中获取，可为空）
			IsSuccess:     log.IsSuccess,     // 设置是否成功（从操作日志实体中获取）
			CreatedAt:     log.CreatedAt,     // 设置创建时间（从操作日志实体中获取）
		}
	}

	return &ListOperationLogsResp{ // 返回操作日志列表响应对象
		List:     list,         // 设置操作日志列表（转换后的操作日志信息列表）
		Total:    total,        // 设置总数（int64类型）
		Page:     req.Page,     // 设置当前页码（从请求中获取）
		PageSize: req.PageSize, // 设置每页数量（从请求中获取）
	}, nil // 返回响应对象和无错误
}
