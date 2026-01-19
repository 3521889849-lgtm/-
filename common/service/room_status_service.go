package service

import (
	"errors"
	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"
	"time"
)

// RoomStatusService 房源状态管理服务
// 负责处理房源状态的更新、查询等核心业务逻辑，
// 包括房源状态验证（启用/停用/维修）、批量状态更新、日历化房态查询等。
type RoomStatusService struct{}

// UpdateRoomStatus 更新房源状态（启用、停用、维修）
// 业务功能：修改指定房源的状态，用于房源启用、停用、维修等业务场景
// 入参说明：
//   - roomID: 房源ID
//   - status: 房源状态（ACTIVE-启用、INACTIVE-停用、MAINTENANCE-维修）
//
// 返回值说明：
//   - error: 房源不存在、状态值无效或数据库操作错误
func (s *RoomStatusService) UpdateRoomStatus(roomID uint64, status string) error {
	// 业务规则：状态值必须有效，只允许使用预定义的状态值
	validStatuses := map[string]bool{ // 创建有效状态值映射表，用于验证状态值是否有效
		"ACTIVE":      true, // 启用状态
		"INACTIVE":    true, // 停用状态
		"MAINTENANCE": true, // 维修状态
	}
	if !validStatuses[status] { // 如果状态值不在有效状态值映射表中（无效状态）
		return errors.New("无效的状态值，支持：ACTIVE（启用）、INACTIVE（停用）、MAINTENANCE（维修）") // 返回错误信息，表示无效的状态值（包含支持的状态值列表）
	}

	// 业务规则：房源必须存在，验证房源是否存在
	var roomInfo hotel_admin.RoomInfo                                 // 声明房源实体变量，用于存储查询到的房源信息
	if err := db.MysqlDB.First(&roomInfo, roomID).Error; err != nil { // 通过房源ID查询房源信息，如果查询失败则说明房源不存在
		return errors.New("房源不存在") // 返回错误信息，表示房源不存在
	}

	roomInfo.Status = status                // 更新房源状态（从参数中获取）
	return db.MysqlDB.Save(&roomInfo).Error // 保存房源信息到数据库，返回保存操作的结果（成功为nil，失败为error）
}

// BatchUpdateRoomStatus 批量更新房源状态
func (s *RoomStatusService) BatchUpdateRoomStatus(roomIDs []uint64, status string) error {
	// 验证状态值
	validStatuses := map[string]bool{ // 创建有效状态值映射表，用于验证状态值是否有效
		"ACTIVE":      true, // 启用状态
		"INACTIVE":    true, // 停用状态
		"MAINTENANCE": true, // 维修状态
	}
	if !validStatuses[status] { // 如果状态值不在有效状态值映射表中（无效状态）
		return errors.New("无效的状态值") // 返回错误信息，表示无效的状态值
	}

	if len(roomIDs) == 0 { // 如果房源ID列表为空（长度为0）
		return errors.New("房源ID列表不能为空") // 返回错误信息，表示房源ID列表不能为空
	}

	return db.MysqlDB.Model(&hotel_admin.RoomInfo{}). // 创建房源模型的查询构建器
								Where("id IN ?", roomIDs).     // 添加筛选条件：房源ID在列表中（使用IN子句）
								Update("status", status).Error // 批量更新房源状态字段，返回更新操作的结果（成功为nil，失败为error）
}

// CalendarRoomStatusReq 日历化房态查询请求
type CalendarRoomStatusReq struct {
	BranchID  *uint64   `json:"branch_id,omitempty"` // 分店ID，可选
	StartDate time.Time `json:"start_date"`          // 开始日期
	EndDate   time.Time `json:"end_date"`            // 结束日期
	RoomNo    *string   `json:"room_no,omitempty"`   // 房间号，可选
	Status    *string   `json:"status,omitempty"`    // 房态筛选，可选
}

// CalendarRoomStatusItem 日历化房态项
type CalendarRoomStatusItem struct {
	RoomID               uint64    `json:"room_id"`
	RoomNo               string    `json:"room_no"`
	RoomName             string    `json:"room_name"`
	Date                 time.Time `json:"date"`
	RoomStatus           string    `json:"room_status"`             // 空净房/入住房/维修房/锁定房/空账房/预定房
	RemainingCount       uint8     `json:"remaining_count"`         // 当日剩余数量
	CheckedInCount       uint8     `json:"checked_in_count"`        // 已入住人数
	CheckOutPendingCount uint8     `json:"check_out_pending_count"` // 预退房人数
	ReservedPendingCount uint8     `json:"reserved_pending_count"`  // 预定待入住人数
}

// GetCalendarRoomStatus 获取日历化房态数据
func (s *RoomStatusService) GetCalendarRoomStatus(req CalendarRoomStatusReq) ([]CalendarRoomStatusItem, error) {
	// 验证日期范围
	if req.StartDate.After(req.EndDate) { // 如果开始日期晚于结束日期（日期范围无效）
		return nil, errors.New("开始日期不能晚于结束日期") // 返回nil和错误信息，表示开始日期不能晚于结束日期
	}

	// 限制查询范围（最多90天）
	days := int(req.EndDate.Sub(req.StartDate).Hours() / 24) // 计算日期范围的天数（结束日期减去开始日期，转换为天数）
	if days > 90 {                                           // 如果日期范围超过90天
		return nil, errors.New("查询日期范围不能超过90天") // 返回nil和错误信息，表示查询日期范围不能超过90天
	}

	// 构建查询：使用JOIN关联房源表，获取房间号和房间名称
	query := db.MysqlDB.Table("room_status_detail"). // 从房态详情表开始查询
		Select(`room_status_detail.room_id,
			room_status_detail.date,
			room_status_detail.room_status,
			room_status_detail.remaining_count,
			room_status_detail.checked_in_count,
			room_status_detail.check_out_pending_count,
			room_status_detail.reserved_pending_count,
			room_info.room_no,
			room_info.room_name`).
		Joins("LEFT JOIN room_info ON room_status_detail.room_id = room_info.id AND room_info.deleted_at IS NULL"). // 通过LEFT JOIN关联房源表（获取房间号和房间名称），并在JOIN条件中过滤已删除的房源
		Where("room_status_detail.date >= ? AND room_status_detail.date <= ?",     // 添加日期范围筛选条件：日期大于等于开始日期且小于等于结束日期
			req.StartDate.Format("2006-01-02"), // 格式化开始日期为YYYY-MM-DD格式
			req.EndDate.Format("2006-01-02")).  // 格式化结束日期为YYYY-MM-DD格式
		Where("room_status_detail.deleted_at IS NULL") // 添加软删除筛选条件（只查询未删除的房态详情）

	// 分店筛选
	if req.BranchID != nil { // 如果请求中提供了分店ID（指针非空），则添加分店筛选条件
		query = query.Where("room_info.branch_id = ?", *req.BranchID) // 添加分店ID筛选条件，只查询指定分店的房态数据（解引用指针获取值）
	}

	// 房间号筛选
	if req.RoomNo != nil && *req.RoomNo != "" { // 如果请求中提供了房间号（指针非空且值非空），则添加房间号筛选条件
		query = query.Where("room_info.room_no LIKE ?", "%"+*req.RoomNo+"%") // 添加模糊搜索条件，搜索房间号包含关键词的房态数据（解引用指针获取值）
	}

	// 房态筛选
	if req.Status != nil && *req.Status != "" { // 如果请求中提供了房态（指针非空且值非空），则添加房态筛选条件
		query = query.Where("room_status_detail.room_status = ?", *req.Status) // 添加房态筛选条件，只查询指定房态的数据（解引用指针获取值）
	}

	// 排序：先按房间号，再按日期
	query = query.Order("room_info.room_no ASC, room_status_detail.date ASC") // 添加排序条件，先按房间号正序排列，再按日期正序排列

	// 定义查询结果结构体
	type Result struct {
		RoomID               uint64    `gorm:"column:room_id"`
		Date                 time.Time `gorm:"column:date"`
		RoomStatus           string    `gorm:"column:room_status"`
		RemainingCount       uint8     `gorm:"column:remaining_count"`
		CheckedInCount       uint8     `gorm:"column:checked_in_count"`
		CheckOutPendingCount uint8     `gorm:"column:check_out_pending_count"`
		ReservedPendingCount uint8     `gorm:"column:reserved_pending_count"`
		RoomNo               string    `gorm:"column:room_no"`
		RoomName             string    `gorm:"column:room_name"`
	}

	var results []Result // 定义查询结果列表
	if err := query.Scan(&results).Error; err != nil { // 执行查询并将结果扫描到结果列表中，如果查询失败则返回错误
		// 如果查询失败，返回详细的错误信息
		return nil, errors.New("查询日历化房态失败: " + err.Error()) // 返回nil和详细的错误信息
	}
	
	// 如果没有查询到数据，返回空数组而不是错误
	if results == nil {
		results = []Result{} // 初始化为空数组
	}

	// 转换为返回格式
	items := make([]CalendarRoomStatusItem, 0, len(results)) // 创建日历化房态项列表，初始容量为查询结果数量
	for _, r := range results {                             // 遍历查询结果列表
		items = append(items, CalendarRoomStatusItem{ // 将查询结果转换为日历化房态项并添加到列表中
			RoomID:               r.RoomID,               // 设置房间ID（从查询结果中获取）
			RoomNo:               r.RoomNo,               // 设置房间号（从查询结果中获取）
			RoomName:             r.RoomName,             // 设置房间名称（从查询结果中获取）
			Date:                 r.Date,                 // 设置日期（从查询结果中获取）
			RoomStatus:           r.RoomStatus,           // 设置房态状态（从查询结果中获取）
			RemainingCount:       r.RemainingCount,       // 设置剩余数量（从查询结果中获取）
			CheckedInCount:       r.CheckedInCount,       // 设置已入住人数（从查询结果中获取）
			CheckOutPendingCount: r.CheckOutPendingCount, // 设置预退房人数（从查询结果中获取）
			ReservedPendingCount: r.ReservedPendingCount, // 设置预定待入住人数（从查询结果中获取）
		})
	}

	return items, nil // 返回日历化房态项列表和无错误
}

// UpdateCalendarRoomStatus 更新或创建日历化房态
func (s *RoomStatusService) UpdateCalendarRoomStatus(roomID uint64, date time.Time, status string) error {
	// 验证房态值
	validStatuses := map[string]bool{ // 创建有效房态值映射表，用于验证房态值是否有效
		"空净房": true, // 空净房状态
		"入住房": true, // 入住房状态
		"维修房": true, // 维修房状态
		"锁定房": true, // 锁定房状态
		"空账房": true, // 空账房状态
		"预定房": true, // 预定房状态
	}
	if !validStatuses[status] { // 如果房态值不在有效房态值映射表中（无效房态）
		return errors.New("无效的房态值，支持：空净房、入住房、维修房、锁定房、空账房、预定房") // 返回错误信息，表示无效的房态值（包含支持的房态值列表）
	}

	// 检查房间是否存在
	var roomInfo hotel_admin.RoomInfo                                 // 声明房源实体变量，用于存储查询到的房源信息
	if err := db.MysqlDB.First(&roomInfo, roomID).Error; err != nil { // 通过房源ID查询房源信息，如果查询失败则说明房源不存在
		return errors.New("房源不存在") // 返回错误信息，表示房源不存在
	}

	// 使用日期部分（去掉时间）
	dateStr := date.Format("2006-01-02")             // 格式化日期为YYYY-MM-DD格式（只保留日期部分，去掉时间）
	dateOnly, _ := time.Parse("2006-01-02", dateStr) // 解析日期字符串为时间对象（只包含日期部分，时间为00:00:00）

	var statusDetail hotel_admin.RoomStatusDetail                                                    // 声明房态详情实体变量，用于存储查询到的房态详情信息
	err := db.MysqlDB.Where("room_id = ? AND date = ?", roomID, dateOnly).First(&statusDetail).Error // 查询该房源和日期的房态详情，如果查询失败则说明房态详情不存在

	if err != nil { // 如果查询失败（房态详情不存在），则需要创建新的房态详情
		// 不存在则创建
		statusDetail = hotel_admin.RoomStatusDetail{ // 创建房态详情实体对象
			RoomID:     roomID,   // 设置房间ID（从参数中获取）
			Date:       dateOnly, // 设置日期（只包含日期部分）
			RoomStatus: status,   // 设置房态状态（从参数中获取）
		}
		return db.MysqlDB.Create(&statusDetail).Error // 将房态详情保存到数据库，返回保存操作的结果（成功为nil，失败为error）
	} else {
		// 存在则更新
		statusDetail.RoomStatus = status            // 更新房态状态（从参数中获取）
		return db.MysqlDB.Save(&statusDetail).Error // 保存房态详情信息到数据库，返回保存操作的结果（成功为nil，失败为error）
	}
}

// BatchUpdateCalendarRoomStatus 批量更新日历化房态
func (s *RoomStatusService) BatchUpdateCalendarRoomStatus(updates []struct {
	RoomID uint64
	Date   time.Time
	Status string
}) error {
	if len(updates) == 0 { // 如果更新列表为空（长度为0）
		return errors.New("更新列表不能为空") // 返回错误信息，表示更新列表不能为空
	}

	for _, update := range updates { // 遍历更新列表，为每个更新项调用单个更新方法
		if err := s.UpdateCalendarRoomStatus(update.RoomID, update.Date, update.Status); err != nil { // 调用单个更新方法，更新指定房源和日期的房态，如果更新失败则返回错误
			return err // 返回更新错误（中断批量更新流程）
		}
	}

	return nil // 返回nil表示批量更新成功（所有房态都更新成功）
}

// RealTimeStatisticsReq 实时数据统计查询请求
type RealTimeStatisticsReq struct {
	BranchID   *uint64 `json:"branch_id,omitempty"`    // 分店ID，可选
	Date       *string `json:"date,omitempty"`         // 日期，可选，默认为今日 YYYY-MM-DD
	RoomNo     *string `json:"room_no,omitempty"`      // 房间号筛选，可选
	RoomTypeID *uint64 `json:"room_type_id,omitempty"` // 房型ID筛选，可选
}

// RealTimeStatisticsResp 实时数据统计响应
type RealTimeStatisticsResp struct {
	Date                 string `json:"date"`                    // 统计日期
	TotalRooms           uint64 `json:"total_rooms"`             // 总房间数
	RemainingRooms       uint64 `json:"remaining_rooms"`         // 剩余房间数
	CheckedInCount       uint64 `json:"checked_in_count"`        // 已入住人数
	CheckOutPendingCount uint64 `json:"check_out_pending_count"` // 预退房人数
	ReservedPendingCount uint64 `json:"reserved_pending_count"`  // 预定待入住人数
	OccupiedRooms        uint64 `json:"occupied_rooms"`          // 已入住房间数
	MaintenanceRooms     uint64 `json:"maintenance_rooms"`       // 维修房间数
	LockedRooms          uint64 `json:"locked_rooms"`            // 锁定房间数
	EmptyRooms           uint64 `json:"empty_rooms"`             // 空净房间数
	ReservedRooms        uint64 `json:"reserved_rooms"`          // 预定房间数
	// 按房态分组统计
	StatusBreakdown []struct {
		Status string `json:"status"` // 房态
		Count  uint64 `json:"count"`  // 数量
	} `json:"status_breakdown"`
	// 按房间明细
	RoomDetails []struct {
		RoomID               uint64 `json:"room_id"`
		RoomNo               string `json:"room_no"`
		RoomName             string `json:"room_name"`
		RoomStatus           string `json:"room_status"`
		RemainingCount       uint8  `json:"remaining_count"`
		CheckedInCount       uint8  `json:"checked_in_count"`
		CheckOutPendingCount uint8  `json:"check_out_pending_count"`
		ReservedPendingCount uint8  `json:"reserved_pending_count"`
	} `json:"room_details,omitempty"`
}

// GetRealTimeStatistics 获取实时数据统计
func (s *RoomStatusService) GetRealTimeStatistics(req RealTimeStatisticsReq) (*RealTimeStatisticsResp, error) {
	// 确定查询日期，默认为今日
	var dateStr string                      // 声明日期字符串变量，用于存储查询日期
	if req.Date != nil && *req.Date != "" { // 如果请求中提供了日期（指针非空且值非空）
		dateStr = *req.Date // 使用请求中提供的日期（解引用指针获取值）
	} else {
		dateStr = time.Now().Format("2006-01-02") // 如果请求中未提供日期，则使用当前日期（格式化为YYYY-MM-DD格式）
	}

	// 验证日期格式
	_, err := time.Parse("2006-01-02", dateStr) // 解析日期字符串，验证日期格式是否正确
	if err != nil {                             // 如果日期格式解析失败，则返回错误
		return nil, errors.New("日期格式错误，应为 YYYY-MM-DD") // 返回nil和错误信息，表示日期格式错误
	}

	// 构建查询
	query := db.MysqlDB.Table("room_status_detail"). // 从房态详情表开始查询
		Select(`room_status_detail.room_id,
			room_info.room_no,
			room_info.room_name,
			room_status_detail.room_status,
			room_status_detail.remaining_count,
			room_status_detail.checked_in_count,
			room_status_detail.check_out_pending_count,
			room_status_detail.reserved_pending_count`).
		Joins("LEFT JOIN room_info ON room_status_detail.room_id = room_info.id AND room_info.deleted_at IS NULL"). // 通过LEFT JOIN关联房源表（获取房间号和房间名称），并在JOIN条件中过滤已删除的房源
		Where("room_status_detail.date = ?", dateStr).                             // 添加日期筛选条件：日期等于查询日期
		Where("room_status_detail.deleted_at IS NULL")                            // 添加软删除筛选条件（只查询未删除的房态详情）

	// 分店筛选
	if req.BranchID != nil { // 如果请求中提供了分店ID（指针非空），则添加分店筛选条件
		query = query.Where("room_info.branch_id = ?", *req.BranchID) // 添加分店ID筛选条件，只查询指定分店的房态数据（解引用指针获取值）
	}

	// 房间号筛选
	if req.RoomNo != nil && *req.RoomNo != "" { // 如果请求中提供了房间号（指针非空且值非空），则添加房间号筛选条件
		query = query.Where("room_info.room_no LIKE ?", "%"+*req.RoomNo+"%") // 添加模糊搜索条件，搜索房间号包含关键词的房态数据（解引用指针获取值）
	}

	// 房型筛选
	if req.RoomTypeID != nil { // 如果请求中提供了房型ID（指针非空），则添加房型筛选条件
		query = query.Where("room_info.room_type_id = ?", *req.RoomTypeID) // 添加房型ID筛选条件，只查询指定房型的房态数据（解引用指针获取值）
	}

	// 查询明细数据
	var details []struct { // 定义查询结果结构体列表，用于存储从数据库查询到的原始房态数据
		RoomID               uint64 // 房间ID
		RoomNo               string // 房间号
		RoomName             string // 房间名称
		RoomStatus           string // 房态状态
		RemainingCount       uint8  // 剩余数量
		CheckedInCount       uint8  // 已入住人数
		CheckOutPendingCount uint8  // 预退房人数
		ReservedPendingCount uint8  // 预定待入住人数
	}

	if err := query.Scan(&details).Error; err != nil { // 执行查询并将结果扫描到结果结构体列表中，如果查询失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	// 统计汇总
	resp := &RealTimeStatisticsResp{ // 创建实时数据统计响应对象指针
		Date: dateStr, // 设置统计日期（从查询日期中获取）
	}

	// 统计各房态数量
	statusCount := make(map[string]uint64) // 创建房态数量映射表，用于按房态分组统计数量
	var totalRemaining uint64              // 声明总剩余数量变量，用于累计所有房间的剩余数量
	var totalCheckedIn uint64              // 声明总已入住人数变量，用于累计所有房间的已入住人数
	var totalCheckOutPending uint64        // 声明总预退房人数变量，用于累计所有房间的预退房人数
	var totalReservedPending uint64        // 声明总预定待入住人数变量，用于累计所有房间的预定待入住人数

	for _, detail := range details { // 遍历查询到的房态明细数据，进行统计汇总
		resp.TotalRooms++                // 增加总房间数（每遍历一个房间，总房间数加1）
		statusCount[detail.RoomStatus]++ // 增加对应房态的数量（按房态分组统计，每个房态的数量加1）

		totalRemaining += uint64(detail.RemainingCount)             // 累计总剩余数量（将当前房间的剩余数量加到总剩余数量中）
		totalCheckedIn += uint64(detail.CheckedInCount)             // 累计总已入住人数（将当前房间的已入住人数加到总已入住人数中）
		totalCheckOutPending += uint64(detail.CheckOutPendingCount) // 累计总预退房人数（将当前房间的预退房人数加到总预退房人数中）
		totalReservedPending += uint64(detail.ReservedPendingCount) // 累计总预定待入住人数（将当前房间的预定待入住人数加到总预定待入住人数中）

		// 添加房间明细
		resp.RoomDetails = append(resp.RoomDetails, struct { // 将房间明细添加到响应对象的房间明细列表中
			RoomID               uint64 `json:"room_id"`
			RoomNo               string `json:"room_no"`
			RoomName             string `json:"room_name"`
			RoomStatus           string `json:"room_status"`
			RemainingCount       uint8  `json:"remaining_count"`
			CheckedInCount       uint8  `json:"checked_in_count"`
			CheckOutPendingCount uint8  `json:"check_out_pending_count"`
			ReservedPendingCount uint8  `json:"reserved_pending_count"`
		}{
			RoomID:               detail.RoomID,               // 设置房间ID（从查询结果中获取）
			RoomNo:               detail.RoomNo,               // 设置房间号（从查询结果中获取）
			RoomName:             detail.RoomName,             // 设置房间名称（从查询结果中获取）
			RoomStatus:           detail.RoomStatus,           // 设置房态状态（从查询结果中获取）
			RemainingCount:       detail.RemainingCount,       // 设置剩余数量（从查询结果中获取）
			CheckedInCount:       detail.CheckedInCount,       // 设置已入住人数（从查询结果中获取）
			CheckOutPendingCount: detail.CheckOutPendingCount, // 设置预退房人数（从查询结果中获取）
			ReservedPendingCount: detail.ReservedPendingCount, // 设置预定待入住人数（从查询结果中获取）
		})
	}

	// 设置汇总数据
	resp.RemainingRooms = totalRemaining             // 设置剩余房间数（使用累计的总剩余数量）
	resp.CheckedInCount = totalCheckedIn             // 设置已入住人数（使用累计的总已入住人数）
	resp.CheckOutPendingCount = totalCheckOutPending // 设置预退房人数（使用累计的总预退房人数）
	resp.ReservedPendingCount = totalReservedPending // 设置预定待入住人数（使用累计的总预定待入住人数）

	// 设置各房态数量
	resp.OccupiedRooms = statusCount["入住房"]    // 设置已入住房间数（从房态数量映射表中获取"入住房"的数量）
	resp.MaintenanceRooms = statusCount["维修房"] // 设置维修房间数（从房态数量映射表中获取"维修房"的数量）
	resp.LockedRooms = statusCount["锁定房"]      // 设置锁定房间数（从房态数量映射表中获取"锁定房"的数量）
	resp.EmptyRooms = statusCount["空净房"]       // 设置空净房间数（从房态数量映射表中获取"空净房"的数量）
	resp.ReservedRooms = statusCount["预定房"]    // 设置预定房间数（从房态数量映射表中获取"预定房"的数量）

	// 设置房态分组统计
	for status, count := range statusCount { // 遍历房态数量映射表，构建房态分组统计列表
		resp.StatusBreakdown = append(resp.StatusBreakdown, struct { // 将房态分组统计项添加到响应对象的房态分组统计列表中
			Status string `json:"status"`
			Count  uint64 `json:"count"`
		}{
			Status: status, // 设置房态（从映射表的键中获取）
			Count:  count,  // 设置数量（从映射表的值中获取）
		})
	}

	return resp, nil // 返回实时数据统计响应对象和无错误
}
