package service

import (
	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"
	"time"
)

// FinancialService 财务管理服务
// 负责处理酒店收支流水的查询、统计等核心业务逻辑，
// 包括收支流水的多条件筛选、按支付方式的汇总统计、收入支出结余计算等。
type FinancialService struct{}

// ListFinancialFlowsReq 收支流水列表查询请求
type ListFinancialFlowsReq struct {
	Page       int     `json:"page"`                  // 页码
	PageSize   int     `json:"page_size"`             // 每页数量
	BranchID   *uint64 `json:"branch_id,omitempty"`   // 分店ID，可选
	FlowType   *string `json:"flow_type,omitempty"`   // 收支类型（收入/支出），可选
	FlowItem   *string `json:"flow_item,omitempty"`   // 收支项目，可选
	PayType    *string `json:"pay_type,omitempty"`    // 支付方式，可选
	OperatorID *uint64 `json:"operator_id,omitempty"` // 操作人ID，可选
	OccurStart *string `json:"occur_start,omitempty"` // 发生开始时间 YYYY-MM-DD
	OccurEnd   *string `json:"occur_end,omitempty"`   // 发生结束时间 YYYY-MM-DD
}

// FinancialFlowInfo 收支流水信息
type FinancialFlowInfo struct {
	ID         uint64    `json:"id"`
	OrderID    *uint64   `json:"order_id,omitempty"`
	BranchID   uint64    `json:"branch_id"`
	RoomID     *uint64   `json:"room_id,omitempty"`
	GuestID    *uint64   `json:"guest_id,omitempty"`
	FlowType   string    `json:"flow_type"`
	FlowItem   string    `json:"flow_item"`
	PayType    string    `json:"pay_type"`
	Amount     float64   `json:"amount"`
	OccurTime  time.Time `json:"occur_time"`
	OperatorID uint64    `json:"operator_id"`
	Remark     *string   `json:"remark,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	// 关联信息
	RoomNo       string `json:"room_no,omitempty"`
	GuestName    string `json:"guest_name,omitempty"`
	ContactPhone string `json:"contact_phone,omitempty"`
	OperatorName string `json:"operator_name,omitempty"`
	OrderNo      string `json:"order_no,omitempty"`
}

// FinancialSummary 财务汇总（按支付方式）
type FinancialSummary struct {
	Total           float64 `json:"total"`            // 合计
	Cash            float64 `json:"cash"`             // 现金
	Alipay          float64 `json:"alipay"`           // 支付宝
	WeChat          float64 `json:"wechat"`           // 微信
	UnionPay        float64 `json:"unionpay"`         // 银联
	CardSwipe       float64 `json:"card_swipe"`       // 刷卡
	TuyouCollection float64 `json:"tuyou_collection"` // 途游代收
	CtripCollection float64 `json:"ctrip_collection"` // 携程代收
	QunarCollection float64 `json:"qunar_collection"` // 去哪儿代收
}

// FinancialSummaryResp 财务汇总响应
type FinancialSummaryResp struct {
	Income  FinancialSummary `json:"income"`  // 收入汇总
	Expense FinancialSummary `json:"expense"` // 支出汇总
	Balance FinancialSummary `json:"balance"` // 结余汇总
}

// ListFinancialFlowsResp 收支流水列表响应
type ListFinancialFlowsResp struct {
	List     []FinancialFlowInfo  `json:"list"`
	Total    uint64               `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
	Summary  FinancialSummaryResp `json:"summary"` // 汇总统计
}

// ListFinancialFlows 获取收支流水列表
// 业务功能：支持多条件筛选和分页查询收支流水列表，并按支付方式汇总统计收入、支出、结余，用于财务管理和分析场景
// 入参说明：
//   - req: 收支流水列表查询请求，支持按分店、收支类型、收支项目、支付方式、操作人、发生时间范围等多维度筛选，支持分页
//
// 返回值说明：
//   - *ListFinancialFlowsResp: 符合条件的收支流水列表（包含关联的订单、房源、客人、操作人信息）、分页信息及汇总统计（按支付方式分类）
//   - error: 查询失败错误
func (s *FinancialService) ListFinancialFlows(req ListFinancialFlowsReq) (*ListFinancialFlowsResp, error) {
	// 业务规则：分页参数默认值设置，页码最小为1，每页数量默认200条（财务流水通常数据量大），最大不超过500条
	if req.Page <= 0 { // 如果页码小于等于0（无效值），则设置为默认值
		req.Page = 1 // 设置页码为1（第一页）
	}
	if req.PageSize <= 0 { // 如果每页数量小于等于0（无效值），则设置为默认值
		req.PageSize = 200 // 设置每页数量为200条（默认值，财务流水通常数据量大）
	}
	if req.PageSize > 500 { // 如果每页数量超过500条（防止查询过大数据集），则限制为最大值
		req.PageSize = 500 // 设置每页数量为500条（最大值）
	}

	offset := (req.Page - 1) * req.PageSize // 计算分页偏移量（跳过前面的记录数），用于SQL查询的OFFSET子句

	// 构建查询：使用预加载（Preload）机制一次性加载所有关联数据，避免N+1查询问题
	query := db.MysqlDB.Model(&hotel_admin.FinancialFlow{}). // 创建收支流水模型的查询构建器
									Preload("Order").           // 预加载订单信息关联数据（JOIN查询订单信息）
									Preload("Room").            // 预加载房源信息关联数据（JOIN查询房源信息）
									Preload("Guest").           // 预加载客人信息关联数据（JOIN查询客人信息）
									Where("deleted_at IS NULL") // 添加软删除筛选条件，只查询未删除的收支流水

	// 分店筛选
	if req.BranchID != nil {
		query = query.Where("branch_id = ?", *req.BranchID)
	}

	// 收支类型筛选
	if req.FlowType != nil && *req.FlowType != "" {
		query = query.Where("flow_type = ?", *req.FlowType)
	}

	// 收支项目筛选
	if req.FlowItem != nil && *req.FlowItem != "" {
		query = query.Where("flow_item = ?", *req.FlowItem)
	}

	// 支付方式筛选
	if req.PayType != nil && *req.PayType != "" {
		query = query.Where("pay_type = ?", *req.PayType)
	}

	// 操作人筛选
	if req.OperatorID != nil {
		query = query.Where("operator_id = ?", *req.OperatorID)
	}

	// 发生时间筛选
	if req.OccurStart != nil && *req.OccurStart != "" {
		query = query.Where("DATE(occur_time) >= ?", *req.OccurStart)
	}
	if req.OccurEnd != nil && *req.OccurEnd != "" {
		query = query.Where("DATE(occur_time) <= ?", *req.OccurEnd)
	}

	// 获取总数
	var total int64                                   // 声明总数变量，用于存储符合条件的收支流水总数
	if err := query.Count(&total).Error; err != nil { // 统计符合条件的收支流水总数，如果统计失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	// 排序：按发生时间倒序
	query = query.Order("occur_time DESC") // 添加排序条件，按收支发生时间倒序排列（最新发生的流水排在前面）

	// 分页
	query = query.Offset(offset).Limit(req.PageSize) // 添加分页限制（偏移量、每页数量），用于SQL查询的OFFSET和LIMIT子句

	// 查询流水列表
	var flows []hotel_admin.FinancialFlow            // 声明收支流水列表变量，用于存储查询到的收支流水信息列表
	if err := query.Find(&flows).Error; err != nil { // 执行查询并获取符合条件的收支流水列表（包含所有预加载的关联数据），如果查询失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	// 转换为返回格式
	flowInfos := make([]FinancialFlowInfo, len(flows)) // 创建收支流水信息列表，长度为查询到的收支流水数量
	for i, flow := range flows {                       // 遍历查询到的收支流水列表
		flowInfo := FinancialFlowInfo{ // 创建收支流水信息对象
			ID:         flow.ID,         // 设置流水ID（从收支流水实体中获取）
			BranchID:   flow.BranchID,   // 设置分店ID（从收支流水实体中获取）
			FlowType:   flow.FlowType,   // 设置收支类型（从收支流水实体中获取：收入/支出）
			FlowItem:   flow.FlowItem,   // 设置收支项目（从收支流水实体中获取）
			PayType:    flow.PayType,    // 设置支付方式（从收支流水实体中获取）
			Amount:     flow.Amount,     // 设置金额（从收支流水实体中获取）
			OccurTime:  flow.OccurTime,  // 设置发生时间（从收支流水实体中获取）
			OperatorID: flow.OperatorID, // 设置操作人ID（从收支流水实体中获取）
			CreatedAt:  flow.CreatedAt,  // 设置创建时间（从收支流水实体中获取）
		}

		if flow.OrderID != nil { // 如果收支流水关联了订单ID（指针非空）
			flowInfo.OrderID = flow.OrderID // 设置订单ID（解引用指针获取值）
			if flow.Order != nil {          // 如果收支流水关联了订单信息（预加载的数据）
				flowInfo.OrderNo = flow.Order.OrderNo // 设置订单号（从关联的订单信息中获取）
			}
		}
		if flow.RoomID != nil { // 如果收支流水关联了房源ID（指针非空）
			flowInfo.RoomID = flow.RoomID // 设置房源ID（解引用指针获取值）
			if flow.Room != nil {         // 如果收支流水关联了房源信息（预加载的数据）
				flowInfo.RoomNo = flow.Room.RoomNo // 设置房间号（从关联的房源信息中获取）
			}
		}
		if flow.GuestID != nil { // 如果收支流水关联了客人ID（指针非空）
			flowInfo.GuestID = flow.GuestID // 设置客人ID（解引用指针获取值）
			if flow.Guest != nil {          // 如果收支流水关联了客人信息（预加载的数据）
				flowInfo.GuestName = flow.Guest.Name     // 设置客人姓名（从关联的客人信息中获取）
				flowInfo.ContactPhone = flow.Guest.Phone // 设置联系电话（从关联的客人信息中获取）
			}
		}
		if flow.Remark != nil { // 如果收支流水有备注（指针非空）
			flowInfo.Remark = flow.Remark // 设置备注（解引用指针获取值）
		}

		flowInfos[i] = flowInfo // 将转换后的收支流水信息对象添加到列表中
	}

	// 计算汇总统计
	summaryQuery := db.MysqlDB.Model(&hotel_admin.FinancialFlow{}).
		Where("deleted_at IS NULL")

	// 应用相同的筛选条件
	if req.BranchID != nil {
		summaryQuery = summaryQuery.Where("branch_id = ?", *req.BranchID)
	}
	if req.FlowType != nil && *req.FlowType != "" {
		summaryQuery = summaryQuery.Where("flow_type = ?", *req.FlowType)
	}
	if req.FlowItem != nil && *req.FlowItem != "" {
		summaryQuery = summaryQuery.Where("flow_item = ?", *req.FlowItem)
	}
	if req.PayType != nil && *req.PayType != "" {
		summaryQuery = summaryQuery.Where("pay_type = ?", *req.PayType)
	}
	if req.OperatorID != nil {
		summaryQuery = summaryQuery.Where("operator_id = ?", *req.OperatorID)
	}
	if req.OccurStart != nil && *req.OccurStart != "" {
		summaryQuery = summaryQuery.Where("DATE(occur_time) >= ?", *req.OccurStart)
	}
	if req.OccurEnd != nil && *req.OccurEnd != "" {
		summaryQuery = summaryQuery.Where("DATE(occur_time) <= ?", *req.OccurEnd)
	}

	// 收入汇总 - 使用 GORM 查询
	var incomeSummary FinancialSummary
	incomeQuery := summaryQuery.Where("flow_type = ?", "收入")

	// 使用分组查询获取各支付方式的汇总
	var incomeResults []struct {
		PayType string  `gorm:"column:pay_type"`
		Total   float64 `gorm:"column:total"`
	}
	incomeQuery.Select("pay_type, SUM(amount) as total").
		Group("pay_type").
		Scan(&incomeResults)

	// 汇总各支付方式
	for _, result := range incomeResults {
		switch result.PayType {
		case "现金":
			incomeSummary.Cash = result.Total
		case "支付宝":
			incomeSummary.Alipay = result.Total
		case "微信":
			incomeSummary.WeChat = result.Total
		case "银联":
			incomeSummary.UnionPay = result.Total
		case "刷卡":
			incomeSummary.CardSwipe = result.Total
		case "途游代收":
			incomeSummary.TuyouCollection = result.Total
		case "携程代收":
			incomeSummary.CtripCollection = result.Total
		case "去哪儿代收":
			incomeSummary.QunarCollection = result.Total
		}
		incomeSummary.Total += result.Total
	}

	// 支出汇总
	var expenseSummary FinancialSummary
	expenseQuery := summaryQuery.Where("flow_type = ?", "支出")

	var expenseResults []struct {
		PayType string  `gorm:"column:pay_type"`
		Total   float64 `gorm:"column:total"`
	}
	expenseQuery.Select("pay_type, SUM(amount) as total").
		Group("pay_type").
		Scan(&expenseResults)

	// 汇总各支付方式
	for _, result := range expenseResults {
		switch result.PayType {
		case "现金":
			expenseSummary.Cash = result.Total
		case "支付宝":
			expenseSummary.Alipay = result.Total
		case "微信":
			expenseSummary.WeChat = result.Total
		case "银联":
			expenseSummary.UnionPay = result.Total
		case "刷卡":
			expenseSummary.CardSwipe = result.Total
		case "途游代收":
			expenseSummary.TuyouCollection = result.Total
		case "携程代收":
			expenseSummary.CtripCollection = result.Total
		case "去哪儿代收":
			expenseSummary.QunarCollection = result.Total
		}
		expenseSummary.Total += result.Total
	}

	// 业务计算：结余汇总 = 收入汇总 - 支出汇总，按支付方式分别计算各支付方式的结余金额
	balanceSummary := FinancialSummary{ // 创建结余汇总对象
		Total:           incomeSummary.Total - expenseSummary.Total,                     // 设置合计结余（收入合计 - 支出合计）
		Cash:            incomeSummary.Cash - expenseSummary.Cash,                       // 设置现金结余（现金收入 - 现金支出）
		Alipay:          incomeSummary.Alipay - expenseSummary.Alipay,                   // 设置支付宝结余（支付宝收入 - 支付宝支出）
		WeChat:          incomeSummary.WeChat - expenseSummary.WeChat,                   // 设置微信结余（微信收入 - 微信支出）
		UnionPay:        incomeSummary.UnionPay - expenseSummary.UnionPay,               // 设置银联结余（银联收入 - 银联支出）
		CardSwipe:       incomeSummary.CardSwipe - expenseSummary.CardSwipe,             // 设置刷卡结余（刷卡收入 - 刷卡支出）
		TuyouCollection: incomeSummary.TuyouCollection - expenseSummary.TuyouCollection, // 设置途游代收结余（途游代收收入 - 途游代收支出）
		CtripCollection: incomeSummary.CtripCollection - expenseSummary.CtripCollection, // 设置携程代收结余（携程代收收入 - 携程代收支出）
		QunarCollection: incomeSummary.QunarCollection - expenseSummary.QunarCollection, // 设置去哪儿代收结余（去哪儿代收收入 - 去哪儿代收支出）
	}

	return &ListFinancialFlowsResp{ // 返回收支流水列表响应对象
		List:     flowInfos,     // 设置收支流水列表（转换后的收支流水信息列表）
		Total:    uint64(total), // 设置总数（转换为uint64类型）
		Page:     req.Page,      // 设置当前页码（从请求中获取）
		PageSize: req.PageSize,  // 设置每页数量（从请求中获取）
		Summary: FinancialSummaryResp{ // 设置汇总统计对象
			Income:  incomeSummary,  // 设置收入汇总（按支付方式分类的收入统计）
			Expense: expenseSummary, // 设置支出汇总（按支付方式分类的支出统计）
			Balance: balanceSummary, // 设置结余汇总（按支付方式分类的结余统计）
		},
	}, nil // 返回响应对象和无错误
}
