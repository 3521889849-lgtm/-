package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// ChannelSyncService 渠道同步服务
// 负责处理酒店房态数据向外部渠道（如OTA平台）的同步等核心业务逻辑，
// 包括渠道配置验证、房态数据查询、HTTP请求发送、同步日志记录等。
type ChannelSyncService struct{}

// SyncRoomStatusReq 同步房态请求
type SyncRoomStatusReq struct {
	BranchID  uint64   `json:"branch_id"`          // 分店ID
	ChannelID uint64   `json:"channel_id"`         // 渠道ID
	StartDate string   `json:"start_date"`         // 开始日期 YYYY-MM-DD
	EndDate   string   `json:"end_date"`           // 结束日期 YYYY-MM-DD
	RoomIDs   []uint64 `json:"room_ids,omitempty"` // 房间ID列表，可选，为空则同步所有房间
}

// SyncRoomStatusResp 同步房态响应
type SyncRoomStatusResp struct {
	SuccessCount int      `json:"success_count"` // 成功数量
	FailCount    int      `json:"fail_count"`    // 失败数量
	SyncLogs     []uint64 `json:"sync_logs"`     // 同步日志ID列表
	Message      string   `json:"message"`       // 提示信息
}

// SyncRoomStatusToChannel 同步房态数据到渠道
// 业务功能：将指定分店、日期范围的房态数据同步到外部渠道（如OTA平台），用于房态数据的实时同步
// 入参说明：
//   - req: 同步房态请求，包含分店ID、渠道ID、开始日期、结束日期、房间ID列表（可选，为空则同步所有房间）
//
// 返回值说明：
//   - *SyncRoomStatusResp: 同步结果（成功数量、失败数量、同步日志ID列表、提示信息）
//   - error: 分店不存在或同步失败错误
func (s *ChannelSyncService) SyncRoomStatusToChannel(req SyncRoomStatusReq) (*SyncRoomStatusResp, error) {
	// 业务规则：分店必须存在，验证分店是否存在
	var branch hotel_admin.HotelBranch                                    // 声明分店实体变量，用于存储查询到的分店信息
	if err := db.MysqlDB.First(&branch, req.BranchID).Error; err != nil { // 通过分店ID查询分店信息，如果查询失败则说明分店不存在
		return nil, errors.New("分店不存在") // 返回nil和错误信息，表示分店不存在
	}

	// 业务规则：渠道必须存在且已启用，验证渠道配置是否存在且状态为启用
	var channel hotel_admin.ChannelConfig                  // 声明渠道配置实体变量，用于存储查询到的渠道配置信息
	err := db.MysqlDB.First(&channel, req.ChannelID).Error // 通过渠道ID查询渠道配置信息，如果查询失败则说明渠道不存在
	if err != nil {                                        // 如果查询失败，说明渠道不存在
		// 业务逻辑：渠道不存在时，返回友好提示信息，不报错（允许未配置渠道的情况）
		log.Printf("渠道ID %d 不存在，返回模拟响应", req.ChannelID) // 记录日志，表示渠道不存在，返回模拟响应
		return &SyncRoomStatusResp{                     // 返回同步响应对象（模拟成功响应）
			SuccessCount: 0,                               // 设置成功数量为0（未同步）
			FailCount:    0,                               // 设置失败数量为0（未同步）
			Message:      "渠道未配置，请先在系统设置中添加渠道配置。本次同步已跳过。", // 设置提示信息，告知用户渠道未配置
		}, nil // 返回响应对象和无错误（不报错，允许未配置渠道的情况）
	}

	// 业务规则：渠道必须处于启用状态，才能进行房态同步
	if channel.Status != "ACTIVE" { // 如果渠道状态不是启用（ACTIVE），则不允许同步
		return &SyncRoomStatusResp{ // 返回同步响应对象（跳过同步）
			SuccessCount: 0,                       // 设置成功数量为0（未同步）
			FailCount:    0,                       // 设置失败数量为0（未同步）
			Message:      "渠道未启用，请先启用渠道。本次同步已跳过。", // 设置提示信息，告知用户渠道未启用
		}, nil // 返回响应对象和无错误（不报错，允许渠道未启用的情况）
	}

	// 业务规则：渠道API URL必须配置，否则无法发送同步请求
	if channel.ApiURL == "" { // 如果渠道API URL为空（未配置），则无法发送同步请求
		log.Printf("渠道 %s (ID=%d) 的API URL为空", channel.ChannelName, channel.ID) // 记录日志，表示渠道API URL为空
		return &SyncRoomStatusResp{                                             // 返回同步响应对象（跳过同步）
			SuccessCount: 0,                                                                        // 设置成功数量为0（未同步）
			FailCount:    0,                                                                        // 设置失败数量为0（未同步）
			Message:      fmt.Sprintf("渠道 %s 的API URL未配置，请先配置API地址。本次同步已跳过。", channel.ChannelName), // 设置提示信息，告知用户渠道API URL未配置（包含渠道名称）
		}, nil // 返回响应对象和无错误（不报错，允许API URL未配置的情况）
	}

	// 业务逻辑：查询指定分店、日期范围内的房态数据，准备同步到渠道
	log.Printf("开始查询房态数据：分店ID=%d, 日期范围=%s ~ %s", req.BranchID, req.StartDate, req.EndDate)

	query := db.MysqlDB.Model(&hotel_admin.RoomStatusDetail{}).
		Select(`
			room_status_detail.room_id,
			room_info.room_no,
			room_info.room_name,
			room_status_detail.date,
			room_status_detail.room_status,
			room_status_detail.remaining_count,
			room_status_detail.checked_in_count,
			room_status_detail.check_out_pending_count,
			room_status_detail.reserved_pending_count
		`).
		Joins("LEFT JOIN room_info ON room_status_detail.room_id = room_info.id").
		Where("room_info.branch_id = ?", req.BranchID).
		Where("room_status_detail.date >= ? AND room_status_detail.date <= ?", req.StartDate, req.EndDate).
		Where("room_status_detail.deleted_at IS NULL").
		Where("room_info.deleted_at IS NULL")

	// 房间筛选
	if len(req.RoomIDs) > 0 { // 如果请求中提供了房间ID列表（长度大于0），则添加房间筛选条件
		query = query.Where("room_status_detail.room_id IN ?", req.RoomIDs) // 添加房间ID筛选条件，只查询指定房间的房态数据（使用IN子句）
	}

	var roomStatuses []struct { // 定义查询结果结构体列表，用于存储从数据库查询到的原始房态数据
		RoomID               uint64    // 房间ID
		RoomNo               string    // 房间号
		RoomName             string    // 房间名称
		Date                 time.Time // 日期
		RoomStatus           string    // 房态状态
		RemainingCount       uint8     // 剩余数量
		CheckedInCount       uint8     // 已入住数量
		CheckOutPendingCount uint8     // 待离店数量
		ReservedPendingCount uint8     // 待入住数量
	}

	if err := query.Scan(&roomStatuses).Error; err != nil { // 执行查询并将结果扫描到结果结构体列表中，如果查询失败则返回错误
		log.Printf("查询房态数据失败: %v", err)             // 记录日志，表示查询房态数据失败
		return nil, fmt.Errorf("查询房态数据失败: %v", err) // 返回nil和错误信息，表示查询房态数据失败
	}

	log.Printf("查询到 %d 条房态数据", len(roomStatuses)) // 记录日志，表示查询到的房态数据条数

	// 如果没有房态数据，直接返回
	if len(roomStatuses) == 0 { // 如果房态数据列表为空（没有需要同步的数据）
		return &SyncRoomStatusResp{ // 返回同步响应对象（跳过同步）
			SuccessCount: 0,                           // 设置成功数量为0（未同步）
			FailCount:    0,                           // 设置失败数量为0（未同步）
			Message:      "当前没有房态数据需要同步。请先创建房源并设置房态。", // 设置提示信息，告知用户当前没有房态数据需要同步
		}, nil // 返回响应对象和无错误（不报错，允许没有房态数据的情况）
	}

	// 调用渠道API同步数据
	client := &http.Client{Timeout: 10 * time.Second} // 创建HTTP客户端，设置超时时间为10秒（防止长时间等待）
	successCount := 0                                 // 声明成功数量变量，初始化为0，用于统计同步成功的房态数量
	failCount := 0                                    // 声明失败数量变量，初始化为0，用于统计同步失败的房态数量
	var syncLogIDs []uint64                           // 声明同步日志ID列表变量，用于存储创建的同步日志ID

	// 构建批量同步数据
	syncData := map[string]interface{}{ // 创建同步数据映射表，用于存储需要同步到渠道的数据
		"branch_code":   branch.BranchCode, // 设置分店编码（从分店信息中获取）
		"branch_name":   branch.HotelName,  // 设置分店名称（从分店信息中获取）
		"start_date":    req.StartDate,     // 设置开始日期（从请求中获取）
		"end_date":      req.EndDate,       // 设置结束日期（从请求中获取）
		"room_statuses": roomStatuses,      // 设置房态数据列表（从查询结果中获取）
	}

	bodyData, err := json.Marshal(syncData) // 将同步数据序列化为JSON格式（用于HTTP请求体）
	if err != nil {                         // 如果序列化失败，则返回错误
		log.Printf("序列化同步数据失败: %v", err)             // 记录日志，表示序列化同步数据失败
		return nil, fmt.Errorf("序列化同步数据失败: %v", err) // 返回nil和错误信息，表示序列化同步数据失败
	}

	log.Printf("准备调用渠道API: %s", channel.ApiURL)                                        // 记录日志，表示准备调用渠道API（包含API URL）
	httpReq, err := http.NewRequest("POST", channel.ApiURL, bytes.NewReader(bodyData)) // 创建HTTP POST请求（使用渠道API URL和序列化后的数据作为请求体）
	if err != nil {                                                                    // 如果创建请求失败，则返回错误
		log.Printf("创建HTTP请求失败: %v", err)         // 记录日志，表示创建HTTP请求失败
		return nil, fmt.Errorf("创建请求失败: %v", err) // 返回nil和错误信息，表示创建请求失败
	}

	httpReq.Header.Set("Content-Type", "application/json") // 设置HTTP请求头，指定内容类型为JSON格式
	resp, err := client.Do(httpReq)                        // 执行HTTP请求（发送到渠道API），如果请求失败则返回错误
	if err != nil {                                        // 如果请求失败，则返回错误
		log.Printf("调用渠道API失败: %v", err)        // 记录日志，表示调用渠道API失败
		return nil, fmt.Errorf("请求失败: %v", err) // 返回nil和错误信息，表示请求失败
	}
	defer resp.Body.Close() // 延迟关闭HTTP响应体（确保资源释放）

	// 读取响应
	body, _ := io.ReadAll(resp.Body) // 读取HTTP响应体的全部内容（忽略读取错误，因为后续会根据状态码判断）

	if resp.StatusCode == http.StatusOK { // 如果HTTP响应状态码为200（成功）
		// 解析响应，获取成功和失败的数量
		var result map[string]interface{}                     // 声明结果映射表变量，用于存储解析后的响应数据
		if err := json.Unmarshal(body, &result); err == nil { // 将响应体JSON数据反序列化为映射表，如果解析成功则进入if块
			if sc, ok := result["success_count"].(float64); ok { // 如果响应中包含success_count字段且类型为float64（JSON数字类型）
				successCount = int(sc) // 将成功数量转换为int类型并赋值给successCount变量
			}
			if fc, ok := result["fail_count"].(float64); ok { // 如果响应中包含fail_count字段且类型为float64（JSON数字类型）
				failCount = int(fc) // 将失败数量转换为int类型并赋值给failCount变量
			}
		} else {
			// 如果无法解析，假设全部成功
			successCount = len(roomStatuses) // 如果解析失败（无法解析响应），则假设所有房态数据都同步成功（设置成功数量为房态数据总数）
		}

		// 为每个房间创建同步日志
		for _, status := range roomStatuses { // 遍历房态数据列表，为每个房间创建同步成功日志
			logID := s.createSyncLog(req.ChannelID, "房态同步", status.RoomID, "成功", "") // 调用创建同步日志函数，创建同步成功日志（不包含失败原因）
			syncLogIDs = append(syncLogIDs, logID)                                   // 将创建的同步日志ID追加到同步日志ID列表中
		}
	} else {
		// 请求失败，记录所有房间为失败
		failCount = len(roomStatuses)         // 如果HTTP响应状态码不是200（失败），则设置失败数量为房态数据总数
		for _, status := range roomStatuses { // 遍历房态数据列表，为每个房间创建同步失败日志
			s.createSyncLog(req.ChannelID, "房态同步", status.RoomID, "失败", fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body))) // 调用创建同步日志函数，创建同步失败日志（包含HTTP状态码和响应内容作为失败原因）
		}
	}

	return &SyncRoomStatusResp{ // 返回同步响应对象
		SuccessCount: successCount, // 设置成功数量（从响应中解析或假设的值）
		FailCount:    failCount,    // 设置失败数量（从响应中解析或计算的值）
		SyncLogs:     syncLogIDs,   // 设置同步日志ID列表（创建的所有同步日志ID）
	}, nil // 返回响应对象和无错误
}

// createSyncLog 创建同步日志
func (s *ChannelSyncService) createSyncLog(channelID uint64, syncType string, syncDataID uint64, syncStatus string, failReason string) uint64 {
	log := hotel_admin.ChannelSyncLog{ // 创建渠道同步日志实体对象
		ChannelID:  channelID,  // 设置渠道ID（从参数中获取）
		SyncType:   syncType,   // 设置同步类型（从参数中获取，如"房态同步"）
		SyncDataID: syncDataID, // 设置同步数据ID（从参数中获取，如房间ID）
		SyncStatus: syncStatus, // 设置同步状态（从参数中获取，如"成功"或"失败"）
		SyncTime:   time.Now(), // 设置同步时间为当前时间（自动生成）
	}

	if failReason != "" { // 如果失败原因不为空（非空字符串）
		log.FailReason = &failReason // 设置失败原因（使用指针类型，解引用后存储失败原因）
	}

	db.MysqlDB.Create(&log) // 将同步日志保存到数据库，如果保存失败则忽略错误（不中断主流程）
	return log.ID           // 返回同步日志ID（创建成功后自动生成的ID）
}
