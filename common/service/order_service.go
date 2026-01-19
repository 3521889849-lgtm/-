package service

import (
	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"
	"time"
)

// OrderService 订单管理服务
// 负责处理酒店订单的查询、管理等核心业务逻辑，
// 包括订单列表多条件筛选（分店、来源、状态、时间范围等）、订单详情查询、
// 订单关联数据（客人、房源、房型、扩展信息等）的加载等。
type OrderService struct{}

// ListOrdersReq 订单列表查询请求
type ListOrdersReq struct {
	Page          int     `json:"page"`                      // 页码
	PageSize      int     `json:"page_size"`                 // 每页数量
	BranchID      *uint64 `json:"branch_id,omitempty"`       // 分店ID，可选
	GuestSource   *string `json:"guest_source,omitempty"`    // 客人来源，可选
	OrderNo       *string `json:"order_no,omitempty"`        // 订单号，可选
	Phone         *string `json:"phone,omitempty"`           // 手机号，可选
	Keyword       *string `json:"keyword,omitempty"`         // 关键词（订单号/房间号/手机号/联系人），可选
	OrderStatus   *string `json:"order_status,omitempty"`    // 订单状态，可选
	CheckInStart  *string `json:"check_in_start,omitempty"`  // 入住开始时间 YYYY-MM-DD
	CheckInEnd    *string `json:"check_in_end,omitempty"`    // 入住结束时间 YYYY-MM-DD
	CheckOutStart *string `json:"check_out_start,omitempty"` // 离店开始时间 YYYY-MM-DD
	CheckOutEnd   *string `json:"check_out_end,omitempty"`   // 离店结束时间 YYYY-MM-DD
	ReserveStart  *string `json:"reserve_start,omitempty"`   // 预定开始时间 YYYY-MM-DD HH:mm:ss
	ReserveEnd    *string `json:"reserve_end,omitempty"`     // 预定结束时间 YYYY-MM-DD HH:mm:ss
}

// OrderInfo 订单信息
type OrderInfo struct {
	ID                uint64    `json:"id"`
	OrderNo           string    `json:"order_no"`
	BranchID          uint64    `json:"branch_id"`
	BranchName        string    `json:"branch_name,omitempty"`
	GuestID           uint64    `json:"guest_id"`
	GuestName         string    `json:"guest_name,omitempty"`
	RoomID            uint64    `json:"room_id"`
	RoomNo            string    `json:"room_no,omitempty"`
	RoomName          string    `json:"room_name,omitempty"`
	RoomTypeID        uint64    `json:"room_type_id"`
	RoomTypeName      string    `json:"room_type_name,omitempty"`
	GuestSource       string    `json:"guest_source"`
	CheckInTime       time.Time `json:"check_in_time"`
	CheckOutTime      time.Time `json:"check_out_time"`
	ReserveTime       time.Time `json:"reserve_time"`
	OrderAmount       float64   `json:"order_amount"`
	DepositReceived   float64   `json:"deposit_received"`
	OutstandingAmount float64   `json:"outstanding_amount"`
	OrderStatus       string    `json:"order_status"`
	PayType           string    `json:"pay_type"`
	PenaltyAmount     float64   `json:"penalty_amount"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	// 扩展信息
	Contact        string  `json:"contact,omitempty"`
	ContactPhone   string  `json:"contact_phone,omitempty"`
	SpecialRequest *string `json:"special_request,omitempty"`
	GuestCount     uint8   `json:"guest_count,omitempty"`
	RoomCount      uint8   `json:"room_count,omitempty"`
	// 关联房间号列表（一个订单可能包含多个房间）
	RoomNos []string `json:"room_nos,omitempty"`
}

// ListOrdersResp 订单列表响应
type ListOrdersResp struct {
	List     []OrderInfo `json:"list"`
	Total    uint64      `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// ListOrders 获取订单列表
// 业务功能：支持多条件筛选和分页查询订单列表，用于订单管理和统计分析场景
// 入参说明：
//   - req: 订单列表查询请求，支持按分店、客人来源、订单状态、订单号、手机号、关键词（订单号/房间号/手机号/联系人）、
//     入住时间范围、离店时间范围、预定时间范围等多维度筛选，支持分页
//
// 返回值说明：
//   - *ListOrdersResp: 符合条件的订单列表（包含关联的分店、客人、房源、房型、扩展信息等）及分页信息
//   - error: 查询失败错误
func (s *OrderService) ListOrders(req ListOrdersReq) (*ListOrdersResp, error) {
	// 业务规则：分页参数默认值设置，页码最小为1，每页数量默认10条，最大不超过100条（防止查询过大数据集）
	if req.Page <= 0 { // 如果页码小于等于0（无效值），则设置为默认值
		req.Page = 1 // 设置页码为1（第一页）
	}
	if req.PageSize <= 0 { // 如果每页数量小于等于0（无效值），则设置为默认值
		req.PageSize = 10 // 设置每页数量为10条（默认值）
	}
	if req.PageSize > 100 { // 如果每页数量超过100条（防止查询过大数据集），则限制为最大值
		req.PageSize = 100 // 设置每页数量为100条（最大值）
	}

	offset := (req.Page - 1) * req.PageSize // 计算分页偏移量（跳过前面的记录数），用于SQL查询的OFFSET子句

	// 构建查询：使用预加载（Preload）机制一次性加载所有关联数据，避免N+1查询问题
	// 关联数据包括：所属分店、客人信息、房源信息、房型信息、订单扩展信息
	query := db.MysqlDB.Model(&hotel_admin.OrderMain{}). // 创建订单模型的查询构建器
								Preload("Branch").          // 预加载所属分店关联数据（JOIN查询分店信息）
								Preload("Guest").           // 预加载客人信息关联数据（JOIN查询客人信息）
								Preload("Room").            // 预加载房源信息关联数据（JOIN查询房源信息）
								Preload("RoomType").        // 预加载房型信息关联数据（JOIN查询房型信息）
								Preload("OrderExtension").  // 预加载订单扩展信息关联数据（JOIN查询订单扩展信息）
								Where("deleted_at IS NULL") // 添加软删除筛选条件，只查询未删除的订单

	// 业务筛选：按分店ID筛选，支持查看特定分店下的所有订单
	if req.BranchID != nil { // 如果请求中提供了分店ID（指针非空），则添加分店筛选条件
		query = query.Where("branch_id = ?", *req.BranchID) // 添加分店ID筛选条件，只查询指定分店的订单（解引用指针获取值）
	}

	// 业务筛选：按客人来源筛选（如线上/线下/OTA等），支持来源维度统计
	if req.GuestSource != nil && *req.GuestSource != "" { // 如果请求中提供了客人来源（指针非空且值非空），则添加来源筛选条件
		query = query.Where("guest_source = ?", *req.GuestSource) // 添加客人来源筛选条件，只查询指定来源的订单（解引用指针获取值）
	}

	// 业务筛选：按订单状态筛选（如已确认/已入住/已离店/已取消等），支持状态维度管理
	if req.OrderStatus != nil && *req.OrderStatus != "" {
		query = query.Where("order_status = ?", *req.OrderStatus)
	}

	// 业务筛选：按订单号模糊搜索，支持部分订单号查询
	if req.OrderNo != nil && *req.OrderNo != "" {
		query = query.Where("order_no LIKE ?", "%"+*req.OrderNo+"%")
	}

	// 业务筛选：按手机号搜索，需要关联订单扩展表和客人信息表
	if req.Phone != nil && *req.Phone != "" {
		query = query.Joins("LEFT JOIN order_extension ON hotel_order_main.id = order_extension.order_id").
			Joins("LEFT JOIN guest_info ON hotel_order_main.guest_id = guest_info.id").
			Where("order_extension.contact_phone LIKE ? OR guest_info.phone LIKE ?", "%"+*req.Phone+"%", "%"+*req.Phone+"%")
	}

	// 业务搜索：关键词多字段模糊搜索，支持订单号、房间号、手机号、联系人姓名等多个维度同时搜索
	if req.Keyword != nil && *req.Keyword != "" {
		keyword := "%" + *req.Keyword + "%"
		query = query.Joins("LEFT JOIN room_info ON hotel_order_main.room_id = room_info.id").
			Joins("LEFT JOIN order_extension ON hotel_order_main.id = order_extension.order_id").
			Joins("LEFT JOIN guest_info ON hotel_order_main.guest_id = guest_info.id").
			Where("hotel_order_main.order_no LIKE ? OR room_info.room_no LIKE ? OR order_extension.contact_phone LIKE ? OR guest_info.phone LIKE ? OR order_extension.contact LIKE ? OR guest_info.name LIKE ?",
				keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// 业务筛选：按入住时间范围筛选，使用DATE函数只比较日期部分，忽略时分秒
	if req.CheckInStart != nil && *req.CheckInStart != "" {
		query = query.Where("DATE(check_in_time) >= ?", *req.CheckInStart)
	}
	if req.CheckInEnd != nil && *req.CheckInEnd != "" {
		query = query.Where("DATE(check_in_time) <= ?", *req.CheckInEnd)
	}

	// 业务筛选：按离店时间范围筛选，使用DATE函数只比较日期部分，忽略时分秒
	if req.CheckOutStart != nil && *req.CheckOutStart != "" {
		query = query.Where("DATE(check_out_time) >= ?", *req.CheckOutStart)
	}
	if req.CheckOutEnd != nil && *req.CheckOutEnd != "" {
		query = query.Where("DATE(check_out_time) <= ?", *req.CheckOutEnd)
	}

	// 业务筛选：按预定时间范围筛选（精确到时分秒），用于查询特定时间段内的预定记录
	if req.ReserveStart != nil && *req.ReserveStart != "" {
		query = query.Where("reserve_time >= ?", *req.ReserveStart)
	}
	if req.ReserveEnd != nil && *req.ReserveEnd != "" {
		query = query.Where("reserve_time <= ?", *req.ReserveEnd)
	}

	// 获取总数：由于JOIN操作可能导致数据重复，需要单独构建计数查询，确保总数准确
	var total int64
	countQuery := db.MysqlDB.Model(&hotel_admin.OrderMain{}).Where("deleted_at IS NULL")

	// 业务规则：计数查询需应用与主查询相同的筛选条件（但不包括JOIN，避免重复计数）
	if req.BranchID != nil {
		countQuery = countQuery.Where("branch_id = ?", *req.BranchID)
	}
	if req.GuestSource != nil && *req.GuestSource != "" {
		countQuery = countQuery.Where("guest_source = ?", *req.GuestSource)
	}
	if req.OrderStatus != nil && *req.OrderStatus != "" {
		countQuery = countQuery.Where("order_status = ?", *req.OrderStatus)
	}
	if req.OrderNo != nil && *req.OrderNo != "" {
		countQuery = countQuery.Where("order_no LIKE ?", "%"+*req.OrderNo+"%")
	}
	if req.Phone != nil && *req.Phone != "" {
		countQuery = countQuery.Joins("LEFT JOIN order_extension ON hotel_order_main.id = order_extension.order_id").
			Joins("LEFT JOIN guest_info ON hotel_order_main.guest_id = guest_info.id").
			Where("order_extension.contact_phone LIKE ? OR guest_info.phone LIKE ?", "%"+*req.Phone+"%", "%"+*req.Phone+"%")
	}
	if req.Keyword != nil && *req.Keyword != "" {
		keyword := "%" + *req.Keyword + "%"
		countQuery = countQuery.Joins("LEFT JOIN room_info ON hotel_order_main.room_id = room_info.id").
			Joins("LEFT JOIN order_extension ON hotel_order_main.id = order_extension.order_id").
			Joins("LEFT JOIN guest_info ON hotel_order_main.guest_id = guest_info.id").
			Where("hotel_order_main.order_no LIKE ? OR room_info.room_no LIKE ? OR order_extension.contact_phone LIKE ? OR guest_info.phone LIKE ? OR order_extension.contact LIKE ? OR guest_info.name LIKE ?",
				keyword, keyword, keyword, keyword, keyword, keyword)
	}
	if req.CheckInStart != nil && *req.CheckInStart != "" {
		countQuery = countQuery.Where("DATE(check_in_time) >= ?", *req.CheckInStart)
	}
	if req.CheckInEnd != nil && *req.CheckInEnd != "" {
		countQuery = countQuery.Where("DATE(check_in_time) <= ?", *req.CheckInEnd)
	}
	if req.CheckOutStart != nil && *req.CheckOutStart != "" {
		countQuery = countQuery.Where("DATE(check_out_time) >= ?", *req.CheckOutStart)
	}
	if req.CheckOutEnd != nil && *req.CheckOutEnd != "" {
		countQuery = countQuery.Where("DATE(check_out_time) <= ?", *req.CheckOutEnd)
	}
	if req.ReserveStart != nil && *req.ReserveStart != "" {
		countQuery = countQuery.Where("reserve_time >= ?", *req.ReserveStart)
	}
	if req.ReserveEnd != nil && *req.ReserveEnd != "" {
		countQuery = countQuery.Where("reserve_time <= ?", *req.ReserveEnd)
	}

	// 使用DISTINCT去重，确保JOIN后的订单计数准确（一个订单可能关联多条扩展记录）
	if err := countQuery.Distinct("hotel_order_main.id").Count(&total).Error; err != nil {
		return nil, err
	}

	// 业务排序：按预定时间倒序排列，最新预定的订单显示在最前面
	query = query.Order("reserve_time DESC")

	// 分页查询：根据偏移量和每页数量限制查询结果集
	query = query.Offset(offset).Limit(req.PageSize)

	// 执行查询，获取符合条件的订单列表（包含所有关联数据）
	var ordersList []hotel_admin.OrderMain
	if err := query.Find(&ordersList).Error; err != nil {
		return nil, err
	}

	// 数据转换：将数据库实体对象转换为业务响应对象，提取关联数据并填充到响应结构中
	orders := make([]OrderInfo, len(ordersList))
	for i, order := range ordersList {
		orderInfo := OrderInfo{
			ID:                order.ID,
			OrderNo:           order.OrderNo,
			BranchID:          order.BranchID,
			GuestID:           order.GuestID,
			RoomID:            order.RoomID,
			RoomTypeID:        order.RoomTypeID,
			GuestSource:       order.GuestSource,
			CheckInTime:       order.CheckInTime,
			CheckOutTime:      order.CheckOutTime,
			ReserveTime:       order.ReserveTime,
			OrderAmount:       order.OrderAmount,
			DepositReceived:   order.DepositReceived,
			OutstandingAmount: order.OutstandingAmount,
			OrderStatus:       order.OrderStatus,
			PayType:           order.PayType,
			PenaltyAmount:     order.PenaltyAmount,
			CreatedAt:         order.CreatedAt,
			UpdatedAt:         order.UpdatedAt,
		}

		// 业务逻辑：从预加载的关联对象中提取业务字段，填充到响应对象中，供前端展示使用
		if order.Branch != nil { // 如果订单关联了分店信息（预加载的数据）
			orderInfo.BranchName = order.Branch.HotelName // 设置分店名称（从关联的分店信息中获取）
		}
		if order.Guest != nil { // 如果订单关联了客人信息（预加载的数据）
			orderInfo.GuestName = order.Guest.Name // 设置客人姓名（从关联的客人信息中获取）
		}
		if order.Room != nil { // 如果订单关联了房源信息（预加载的数据）
			orderInfo.RoomNo = order.Room.RoomNo            // 设置房间号（从关联的房源信息中获取）
			orderInfo.RoomName = order.Room.RoomName        // 设置房间名称（从关联的房源信息中获取）
			orderInfo.RoomNos = []string{order.Room.RoomNo} // 初始化房间号列表，包含当前房源房间号
		}
		if order.RoomType != nil { // 如果订单关联了房型信息（预加载的数据）
			orderInfo.RoomTypeName = order.RoomType.RoomTypeName // 设置房型名称（从关联的房型信息中获取）
		}
		if order.OrderExtension != nil { // 如果订单关联了订单扩展信息（预加载的数据）
			orderInfo.Contact = order.OrderExtension.Contact               // 设置联系人（从关联的订单扩展信息中获取）
			orderInfo.ContactPhone = order.OrderExtension.ContactPhone     // 设置联系电话（从关联的订单扩展信息中获取）
			orderInfo.SpecialRequest = order.OrderExtension.SpecialRequest // 设置特殊要求（从关联的订单扩展信息中获取）
			orderInfo.GuestCount = order.OrderExtension.GuestCount         // 设置客人数量（从关联的订单扩展信息中获取）
			orderInfo.RoomCount = order.OrderExtension.RoomCount           // 设置房间数量（从关联的订单扩展信息中获取）
		}

		orders[i] = orderInfo // 将转换后的订单信息对象添加到列表中
	}

	// 业务逻辑：一个订单可能包含多个房间（同订单号不同房源），查询该订单的所有房间号列表
	for i := range orders { // 遍历订单列表，为每个订单查询同订单号的其他房间号
		var otherRoomNos []string                   // 声明其他房间号列表变量，用于存储查询到的同订单号的其他房间号
		db.MysqlDB.Model(&hotel_admin.OrderMain{}). // 创建订单模型的查询构建器
								Select("room_info.room_no").                                                                                   // 选择房间号字段
								Joins("LEFT JOIN room_info ON hotel_order_main.room_id = room_info.id").                                       // 通过LEFT JOIN关联房源表（获取房间号）
								Where("hotel_order_main.order_no = ? AND hotel_order_main.room_id != ?", orders[i].OrderNo, orders[i].RoomID). // 筛选同订单号但不同房源ID的订单
								Where("hotel_order_main.deleted_at IS NULL").                                                                  // 添加软删除筛选条件（只查询未删除的订单）
								Where("room_info.deleted_at IS NULL").                                                                         // 添加软删除筛选条件（只查询未删除的房源）
								Pluck("room_info.room_no", &otherRoomNos)                                                                      // 提取房间号字段值到列表中（Pluck函数用于提取单列数据）

		orders[i].RoomNos = append(orders[i].RoomNos, otherRoomNos...) // 将查询到的其他房间号追加到当前订单的房间号列表中
	}

	return &ListOrdersResp{ // 返回订单列表响应对象
		List:     orders,        // 设置订单列表（转换后的订单信息列表，包含关联数据和同订单号的多个房间号）
		Total:    uint64(total), // 设置总数（转换为uint64类型）
		Page:     req.Page,      // 设置当前页码（从请求中获取）
		PageSize: req.PageSize,  // 设置每页数量（从请求中获取）
	}, nil // 返回响应对象和无错误
}

// GetOrder 获取订单详情
// 业务功能：根据订单ID查询订单的完整信息，包含所有关联数据（分店、客人、房源、房型、扩展信息等）
// 入参说明：
//   - orderID: 订单ID
//
// 返回值说明：
//   - *OrderInfo: 订单完整信息（包含所有关联对象数据及同订单号的多个房间号列表）
//   - error: 订单不存在或查询失败
func (s *OrderService) GetOrder(orderID uint64) (*OrderInfo, error) {
	var order hotel_admin.OrderMain         // 声明订单实体变量，用于存储查询到的订单信息
	if err := db.MysqlDB.Preload("Branch"). // 预加载所属分店关联数据（JOIN查询分店信息）
						Preload("Guest").                          // 预加载客人信息关联数据（JOIN查询客人信息）
						Preload("Room").                           // 预加载房源信息关联数据（JOIN查询房源信息）
						Preload("RoomType").                       // 预加载房型信息关联数据（JOIN查询房型信息）
						Preload("OrderExtension").                 // 预加载订单扩展信息关联数据（JOIN查询订单扩展信息）
						First(&order, orderID).Error; err != nil { // 通过订单ID查询订单信息（包含所有预加载的关联数据），如果查询失败则说明订单不存在
		return nil, err // 返回nil和数据库查询错误
	}

	orderInfo := &OrderInfo{ // 创建订单信息对象指针
		ID:                order.ID,                // 设置订单ID（从订单实体中获取）
		OrderNo:           order.OrderNo,           // 设置订单号（从订单实体中获取）
		BranchID:          order.BranchID,          // 设置分店ID（从订单实体中获取）
		GuestID:           order.GuestID,           // 设置客人ID（从订单实体中获取）
		RoomID:            order.RoomID,            // 设置房源ID（从订单实体中获取）
		RoomTypeID:        order.RoomTypeID,        // 设置房型ID（从订单实体中获取）
		GuestSource:       order.GuestSource,       // 设置客人来源（从订单实体中获取）
		CheckInTime:       order.CheckInTime,       // 设置入住时间（从订单实体中获取）
		CheckOutTime:      order.CheckOutTime,      // 设置离店时间（从订单实体中获取）
		ReserveTime:       order.ReserveTime,       // 设置预定时间（从订单实体中获取）
		OrderAmount:       order.OrderAmount,       // 设置订单金额（从订单实体中获取）
		DepositReceived:   order.DepositReceived,   // 设置已收押金（从订单实体中获取）
		OutstandingAmount: order.OutstandingAmount, // 设置未付金额（从订单实体中获取）
		OrderStatus:       order.OrderStatus,       // 设置订单状态（从订单实体中获取）
		PayType:           order.PayType,           // 设置支付方式（从订单实体中获取）
		PenaltyAmount:     order.PenaltyAmount,     // 设置违约金（从订单实体中获取）
		CreatedAt:         order.CreatedAt,         // 设置创建时间（从订单实体中获取）
		UpdatedAt:         order.UpdatedAt,         // 设置更新时间（从订单实体中获取）
	}

	if order.Branch != nil { // 如果订单关联了分店信息（预加载的数据）
		orderInfo.BranchName = order.Branch.HotelName // 设置分店名称（从关联的分店信息中获取）
	}
	if order.Guest != nil { // 如果订单关联了客人信息（预加载的数据）
		orderInfo.GuestName = order.Guest.Name // 设置客人姓名（从关联的客人信息中获取）
	}
	if order.Room != nil { // 如果订单关联了房源信息（预加载的数据）
		orderInfo.RoomNo = order.Room.RoomNo            // 设置房间号（从关联的房源信息中获取）
		orderInfo.RoomName = order.Room.RoomName        // 设置房间名称（从关联的房源信息中获取）
		orderInfo.RoomNos = []string{order.Room.RoomNo} // 初始化房间号列表，包含当前房源房间号
	}
	if order.RoomType != nil { // 如果订单关联了房型信息（预加载的数据）
		orderInfo.RoomTypeName = order.RoomType.RoomTypeName // 设置房型名称（从关联的房型信息中获取）
	}
	if order.OrderExtension != nil { // 如果订单关联了订单扩展信息（预加载的数据）
		orderInfo.Contact = order.OrderExtension.Contact               // 设置联系人（从关联的订单扩展信息中获取）
		orderInfo.ContactPhone = order.OrderExtension.ContactPhone     // 设置联系电话（从关联的订单扩展信息中获取）
		orderInfo.SpecialRequest = order.OrderExtension.SpecialRequest // 设置特殊要求（从关联的订单扩展信息中获取）
		orderInfo.GuestCount = order.OrderExtension.GuestCount         // 设置客人数量（从关联的订单扩展信息中获取）
		orderInfo.RoomCount = order.OrderExtension.RoomCount           // 设置房间数量（从关联的订单扩展信息中获取）
	}

	// 业务逻辑：查询同一订单号下的其他房源房间号（同一订单可能包含多个房间），合并到房间号列表中
	var otherRoomNos []string                   // 声明其他房间号列表变量，用于存储查询到的同订单号的其他房间号
	db.MysqlDB.Model(&hotel_admin.OrderMain{}). // 创建订单模型的查询构建器
							Select("room_info.room_no").                                                                  // 选择房间号字段
							Joins("LEFT JOIN room_info ON hotel_order_main.room_id = room_info.id").                      // 通过LEFT JOIN关联房源表（获取房间号）
							Where("hotel_order_main.order_no = ? AND hotel_order_main.id != ?", order.OrderNo, order.ID). // 筛选同订单号但不同订单ID的订单（排除当前订单）
							Where("hotel_order_main.deleted_at IS NULL").                                                 // 添加软删除筛选条件（只查询未删除的订单）
							Where("room_info.deleted_at IS NULL").                                                        // 添加软删除筛选条件（只查询未删除的房源）
							Pluck("room_info.room_no", &otherRoomNos)                                                     // 提取房间号字段值到列表中（Pluck函数用于提取单列数据）

	orderInfo.RoomNos = append(orderInfo.RoomNos, otherRoomNos...) // 将查询到的其他房间号追加到当前订单的房间号列表中

	return orderInfo, nil // 返回订单信息指针和无错误
}
