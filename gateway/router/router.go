// Package router 提供 Gateway 网关的路由配置
// 负责将HTTP请求映射到对应的Handler处理函数
// 路由分组：
// - WebSocket: 实时通信接口
// - 公共接口: 无需登录即可访问
// - 用户接口: 需要登录才能访问
// - 业务接口: 客服、排班、会话等核心功能
package router

import (
	"encoding/json"
	"net/http"

	"example_shop/gateway/handler"    // HTTP请求处理器
	"example_shop/gateway/middleware" // 中间件（认证、权限等）
	"example_shop/gateway/ws"         // WebSocket模块
)

// SetupRoutes 配置所有HTTP路由
// 创建并返回一个配置好的HTTP路由器
// 参数:
//   - customerHandler: 客服业务处理器
//   - hub: WebSocket Hub实例
//
// 返回: 配置完成的HTTP多路复用器
func SetupRoutes(customerHandler *handler.CustomerHandler, hub *ws.Hub) *http.ServeMux {
	mux := http.NewServeMux()

	// ============ WebSocket 实时通信接口 ============
	// WebSocket连接，支持客服与用户的实时消息通信
	// Token可从 URL query 或 Sec-WebSocket-Protocol 头传入
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// 优先从 URL 获取 token

		token := r.URL.Query().Get("token")
		if token == "" {
			// 备选：从 WebSocket 子协议头获取
			token = r.Header.Get("Sec-WebSocket-Protocol")
		}
		// 验证 Token
		claims, err := middleware.ParseToken(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// 如果从 Sec-WebSocket-Protocol 获取的 token，需要设置对应的响应头
		if r.Header.Get("Sec-WebSocket-Protocol") != "" {
			w.Header().Set("Sec-WebSocket-Protocol", token)
		}
		// 升级为 WebSocket 连接
		ws.ServeWs(hub, w, r, claims.UserID, claims.RoleCode)
	})

	// 在线状态统计接口（获取当前在线客服和连接数）
	mux.HandleFunc("/api/stats/online", func(w http.ResponseWriter, r *http.Request) {
		stats := hub.GetOnlineStats()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 0,
			"data": stats,
			"msg":  "success",
		})
	})

	// ============ 公共接口（无需登录） ============
	mux.HandleFunc("/api/v1/user/login", customerHandler.Login)       // 用户登录
	mux.HandleFunc("/api/v1/user/register", customerHandler.Register) // 用户注册（仅客服账号）

	// ============ 需要登录的接口 ============
	mux.HandleFunc("/api/v1/user/current", customerHandler.GetCurrentUser) // 获取当前用户信息

	// ============ 客服管理接口 ============
	mux.HandleFunc("/api/customer/get", customerHandler.GetCustomerService)   // 获取单个客服信息
	mux.HandleFunc("/api/customer/list", customerHandler.ListCustomerService) // 查询客服列表

	// ============ 班次配置接口 ============
	mux.HandleFunc("/api/shift/create", customerHandler.CreateShiftConfig) // 创建班次配置
	mux.HandleFunc("/api/shift/list", customerHandler.ListShiftConfig)     // 查询班次列表
	mux.HandleFunc("/api/shift/update", customerHandler.UpdateShiftConfig) // 更新班次配置
	mux.HandleFunc("/api/shift/delete", customerHandler.DeleteShiftConfig) // 删除班次配置

	// ============ 排班管理接口 ============
	mux.HandleFunc("/api/schedule/assign", customerHandler.AssignSchedule)          // 手动分配排班
	mux.HandleFunc("/api/schedule/auto", customerHandler.AutoSchedule)              // 自动排班
	mux.HandleFunc("/api/schedule/grid", customerHandler.ListScheduleGrid)          // 查询排班表格数据
	mux.HandleFunc("/api/schedule/cell/upsert", customerHandler.UpsertScheduleCell) // 更新排班单元格
	mux.HandleFunc("/api/schedule/export", customerHandler.ExportScheduleExcel)     // 导出排班Excel

	// ============ 请假/调班管理接口 ============
	mux.HandleFunc("/api/leave/apply", customerHandler.ApplyLeaveTransfer)     // 提交请假/调班申请
	mux.HandleFunc("/api/leave/approve", customerHandler.ApproveLeaveTransfer) // 审批申请
	mux.HandleFunc("/api/leave/get", customerHandler.GetLeaveTransfer)         // 获取申请详情
	mux.HandleFunc("/api/leave/list", customerHandler.ListLeaveTransfer)       // 查询申请列表

	// ============ 会话管理接口 ============
	mux.HandleFunc("/api/conversation/create", customerHandler.CreateConversation)            // 创建会话
	mux.HandleFunc("/api/conversation/end", customerHandler.EndConversation)                  // 结束会话
	mux.HandleFunc("/api/conversation/transfer", customerHandler.TransferConversation)        // 转接会话
	mux.HandleFunc("/api/conversation/assign", customerHandler.AssignCustomer)                // 自动分配客服（用户发起咨询时调用）
	mux.HandleFunc("/api/conversation/list", customerHandler.ListConversation)                // 查询会话列表
	mux.HandleFunc("/api/conversation/history/list", customerHandler.ListConversationHistory) // 查询历史会话
	mux.HandleFunc("/api/conversation/message/list", customerHandler.ListConversationMessage) // 查询会话消息
	mux.HandleFunc("/api/conversation/message/send", customerHandler.SendConversationMessage) // 发送消息

	// ============ 快捷回复接口 ============
	mux.HandleFunc("/api/quick_reply/list", customerHandler.ListQuickReply)     // 查询快捷回复列表
	mux.HandleFunc("/api/quick_reply/create", customerHandler.CreateQuickReply) // 创建快捷回复
	mux.HandleFunc("/api/quick_reply/update", customerHandler.UpdateQuickReply) // 更新快捷回复
	mux.HandleFunc("/api/quick_reply/delete", customerHandler.DeleteQuickReply) // 删除快捷回复

	// ============ 会话分类管理接口 ============
	mux.HandleFunc("/api/conversation/category/create", customerHandler.CreateConvCategory)         // 创建会话分类
	mux.HandleFunc("/api/conversation/category/list", customerHandler.ListConvCategory)             // 查询分类列表
	mux.HandleFunc("/api/conversation/classify/update", customerHandler.UpdateConversationClassify) // 更新会话分类标签

	// ============ 会话标签管理接口 ============
	mux.HandleFunc("/api/conversation/tag/create", customerHandler.CreateConvTag) // 创建标签
	mux.HandleFunc("/api/conversation/tag/list", customerHandler.ListConvTag)     // 查询标签列表
	mux.HandleFunc("/api/conversation/tag/update", customerHandler.UpdateConvTag) // 更新标签
	mux.HandleFunc("/api/conversation/tag/delete", customerHandler.DeleteConvTag) // 删除标签

	// ============ 统计看板接口 ============
	mux.HandleFunc("/api/conversation/stats", customerHandler.GetConversationStats) // 获取会话统计数据

	// ============ 会话监控与导出接口 ============
	mux.HandleFunc("/api/conversation/monitor", customerHandler.GetConversationMonitor) // 获取会话监控数据（实时）
	mux.HandleFunc("/api/conversation/export", customerHandler.ExportConversations)     // 导出会话记录

	// ============ 消息分类管理接口 ============
	mux.HandleFunc("/api/msg/classify/auto", customerHandler.MsgAutoClassify)     // 消息自动分类
	mux.HandleFunc("/api/msg/classify/adjust", customerHandler.AdjustMsgClassify) // 人工调整分类
	mux.HandleFunc("/api/msg/classify/stats", customerHandler.GetClassifyStats)   // 分类统计数据

	// ============ 消息分类维度管理接口 ============
	mux.HandleFunc("/api/msg/category/create", customerHandler.CreateMsgCategory) // 创建消息分类维度
	mux.HandleFunc("/api/msg/category/list", customerHandler.ListMsgCategory)     // 查询消息分类维度列表
	mux.HandleFunc("/api/msg/category/update", customerHandler.UpdateMsgCategory) // 更新消息分类维度
	mux.HandleFunc("/api/msg/category/delete", customerHandler.DeleteMsgCategory) // 删除消息分类维度

	// ============ 消息加密与脱敏接口 ============
	mux.HandleFunc("/api/msg/encrypt", customerHandler.EncryptMessage)         // 加密消息内容
	mux.HandleFunc("/api/msg/decrypt", customerHandler.DecryptMessage)         // 解密消息内容
	mux.HandleFunc("/api/msg/desensitize", customerHandler.DesensitizeMessage) // 消息脱敏处理

	// ============ 数据归档管理接口 ============
	mux.HandleFunc("/api/archive/conversations", customerHandler.ArchiveConversations) // 归档历史会话
	mux.HandleFunc("/api/archive/task", customerHandler.GetArchiveTask)                // 获取归档任务状态
	mux.HandleFunc("/api/archive/query", customerHandler.QueryArchivedConversation)    // 查询归档会话

	// ============ 系统接口 ============
	// 健康检查：用于探活与部署环境健康监测
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return mux
}
