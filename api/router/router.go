package router

import (
	"crypto/rand"
	"encoding/hex"
	"example_shop/api/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			var b [16]byte
			_, _ = rand.Read(b[:])
			requestID = hex.EncodeToString(b[:])
		}
		c.Header("X-Request-Id", requestID)
		c.Set("request_id", requestID)
		c.Next()
	})

	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-Id")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// 初始化处理器
	roomHandler := handler.NewRoomHandler()
	branchHandler := handler.NewBranchHandler()
	userAccountHandler := handler.NewUserAccountHandler()
	roleHandler := handler.NewRoleHandler()
	channelConfigHandler := handler.NewChannelConfigHandler()
	systemConfigHandler := handler.NewSystemConfigHandler()
	blacklistHandler := handler.NewBlacklistHandler()
	operationLogHandler := handler.NewOperationLogHandler()

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 房型管理路由
		roomTypes := api.Group("/room-types")
		{
			roomTypes.POST("", roomHandler.CreateRoomType)       // 创建房型
			roomTypes.GET("", roomHandler.ListRoomTypes)         // 获取房型列表
			roomTypes.GET("/:id", roomHandler.GetRoomType)       // 获取房型详情
			roomTypes.PUT("/:id", roomHandler.UpdateRoomType)    // 更新房型
			roomTypes.DELETE("/:id", roomHandler.DeleteRoomType) // 删除房型
		}

		// 房源管理路由
		roomInfos := api.Group("/room-infos")
		{
			roomInfos.POST("", roomHandler.CreateRoomInfo)                                   // 创建房源
			roomInfos.GET("", roomHandler.ListRoomInfos)                                     // 获取房源列表
			roomInfos.GET("/:id", roomHandler.GetRoomInfo)                                   // 获取房源详情
			roomInfos.PUT("/:id", roomHandler.UpdateRoomInfo)                                // 更新房源
			roomInfos.DELETE("/:id", roomHandler.DeleteRoomInfo)                             // 删除房源
			roomInfos.PUT("/:id/status", roomHandler.UpdateRoomStatus)                       // 更新房源状态
			roomInfos.PUT("/batch-status", roomHandler.BatchUpdateRoomStatus)                // 批量更新房源状态
			roomInfos.POST("/bindings", roomHandler.CreateRoomBinding)                       // 创建关联房绑定
			roomInfos.POST("/batch-bindings", roomHandler.BatchCreateRoomBindings)           // 批量创建关联房绑定
			roomInfos.GET("/:id/bindings", roomHandler.GetRoomBindings)                      // 获取关联房列表
			roomInfos.DELETE("/bindings/:id", roomHandler.DeleteRoomBinding)                 // 删除关联房绑定
			roomInfos.POST("/:id/images", roomHandler.UploadRoomImages)                      // 批量上传房源图片
			roomInfos.GET("/:id/images", roomHandler.GetRoomImages)                          // 获取房源图片列表
			roomInfos.PUT("/:id/images/sort", roomHandler.BatchUpdateImageSortOrder)         // 批量更新图片排序
			roomInfos.DELETE("/images/:id", roomHandler.DeleteRoomImage)                     // 删除房源图片
			roomInfos.PUT("/images/:id/sort", roomHandler.UpdateImageSortOrder)              // 更新图片排序
			roomInfos.PUT("/:id/facilities", roomHandler.SetRoomFacilities)                  // 设置房源设施
			roomInfos.GET("/:id/facilities", roomHandler.GetRoomFacilities)                  // 获取房源设施列表
			roomInfos.POST("/:id/facilities", roomHandler.AddRoomFacility)                   // 为房源添加单个设施
			roomInfos.DELETE("/:id/facilities/:facility_id", roomHandler.RemoveRoomFacility) // 移除房源的单个设施
		}

		// 设施管理路由
		facilities := api.Group("/facilities")
		{
			facilities.POST("", roomHandler.CreateFacility)       // 创建设施
			facilities.GET("", roomHandler.ListFacilities)        // 获取设施列表
			facilities.GET("/:id", roomHandler.GetFacility)       // 获取设施详情
			facilities.PUT("/:id", roomHandler.UpdateFacility)    // 更新设施
			facilities.DELETE("/:id", roomHandler.DeleteFacility) // 删除设施
		}

		// 退订政策管理路由
		cancellationPolicies := api.Group("/cancellation-policies")
		{
			cancellationPolicies.POST("", roomHandler.CreateCancellationPolicy)       // 创建退订政策
			cancellationPolicies.GET("", roomHandler.ListCancellationPolicies)        // 获取退订政策列表
			cancellationPolicies.GET("/:id", roomHandler.GetCancellationPolicy)       // 获取退订政策详情
			cancellationPolicies.PUT("/:id", roomHandler.UpdateCancellationPolicy)    // 更新退订政策
			cancellationPolicies.DELETE("/:id", roomHandler.DeleteCancellationPolicy) // 删除退订政策
		}

		// 日历化房态管理路由
		calendarStatus := api.Group("/calendar-room-status")
		{
			calendarStatus.GET("", roomHandler.GetCalendarRoomStatus)               // 获取日历化房态
			calendarStatus.PUT("", roomHandler.UpdateCalendarRoomStatus)            // 更新日历化房态
			calendarStatus.PUT("/batch", roomHandler.BatchUpdateCalendarRoomStatus) // 批量更新日历化房态
		}

		// 实时数据统计路由
		api.GET("/real-time-statistics", roomHandler.GetRealTimeStatistics) // 获取实时数据统计

		// 分店管理路由
		branches := api.Group("/branches")
		{
			branches.GET("", branchHandler.ListBranches)       // 获取分店列表
			branches.GET("/:id", branchHandler.GetBranch)      // 获取分店详情
			branches.POST("", branchHandler.CreateBranch)      // 创建分店
			branches.PUT("/:id", branchHandler.UpdateBranch)   // 更新分店
			branches.DELETE("/:id", branchHandler.DeleteBranch) // 删除分店
		}

		// 渠道同步路由
		api.POST("/sync-room-status", roomHandler.SyncRoomStatusToChannel) // 同步房态到渠道

		// 订单管理路由
		orders := api.Group("/orders")
		{
			orders.GET("", roomHandler.ListOrders)   // 获取订单列表
			orders.GET("/:id", roomHandler.GetOrder) // 获取订单详情
		}

		// 在住客人管理路由
		api.GET("/in-house-guests", roomHandler.ListInHouseGuests) // 获取在住客人列表

		// 财务管理路由
		api.GET("/financial-flows", roomHandler.ListFinancialFlows) // 获取收支流水列表

		// 账号管理
		userAccounts := api.Group("/user-accounts")
		{
			userAccounts.POST("", userAccountHandler.CreateUserAccount)
			userAccounts.GET("", userAccountHandler.ListUserAccounts)
			userAccounts.GET("/:id", userAccountHandler.GetUserAccount)
			userAccounts.PUT("/:id", userAccountHandler.UpdateUserAccount)
			userAccounts.DELETE("/:id", userAccountHandler.DeleteUserAccount)
		}

		// 角色管理
		roles := api.Group("/roles")
		{
			roles.POST("", roleHandler.CreateRole)
			roles.GET("", roleHandler.ListRoles)
			roles.GET("/:id", roleHandler.GetRole)
			roles.PUT("/:id", roleHandler.UpdateRole)
			roles.DELETE("/:id", roleHandler.DeleteRole)
		}

		// 权限管理
		permissions := api.Group("/permissions")
		{
			permissions.GET("", roleHandler.ListPermissions)
		}

		// 渠道配置管理
		channelConfigs := api.Group("/channel-configs")
		{
			channelConfigs.POST("", channelConfigHandler.CreateChannelConfig)
			channelConfigs.GET("", channelConfigHandler.ListChannelConfigs)
			channelConfigs.GET("/:id", channelConfigHandler.GetChannelConfig)
			channelConfigs.PUT("/:id", channelConfigHandler.UpdateChannelConfig)
			channelConfigs.DELETE("/:id", channelConfigHandler.DeleteChannelConfig)
		}

		// 系统配置管理
		systemConfigs := api.Group("/system-configs")
		{
			systemConfigs.POST("", systemConfigHandler.CreateSystemConfig)
			systemConfigs.GET("", systemConfigHandler.ListSystemConfigs)
			systemConfigs.GET("/:id", systemConfigHandler.GetSystemConfig)
			systemConfigs.PUT("/:id", systemConfigHandler.UpdateSystemConfig)
			systemConfigs.DELETE("/:id", systemConfigHandler.DeleteSystemConfig)
			systemConfigs.GET("/category/:category", systemConfigHandler.GetSystemConfigsByCategory)
		}

		// 黑名单管理
		blacklists := api.Group("/blacklists")
		{
			blacklists.POST("", blacklistHandler.CreateBlacklist)
			blacklists.GET("", blacklistHandler.ListBlacklists)
			blacklists.GET("/:id", blacklistHandler.GetBlacklist)
			blacklists.PUT("/:id", blacklistHandler.UpdateBlacklist)
			blacklists.DELETE("/:id", blacklistHandler.DeleteBlacklist)
		}

		// 会员管理相关 Handler（需要先定义，因为路由中会使用）
		memberHandler := handler.NewMemberHandler()
		memberRightsHandler := handler.NewMemberRightsHandler()
		memberPointsHandler := handler.NewMemberPointsHandler()

		// 会员管理
		members := api.Group("/members")
		{
			members.POST("", memberHandler.CreateMember)
			members.GET("", memberHandler.ListMembers)
			// 注意：更具体的路由要放在通用路由之前
			members.GET("/guest/:guest_id", memberHandler.GetMemberByGuestID)
			members.GET("/:id/points-balance", memberPointsHandler.GetMemberPointsBalance)
			members.GET("/:id", memberHandler.GetMember)
			members.PUT("/:id", memberHandler.UpdateMember)
			members.DELETE("/:id", memberHandler.DeleteMember)
		}

		// 会员权益管理
		memberRights := api.Group("/member-rights")
		{
			memberRights.POST("", memberRightsHandler.CreateMemberRights)
			memberRights.GET("", memberRightsHandler.ListMemberRights)
			memberRights.GET("/:id", memberRightsHandler.GetMemberRights)
			memberRights.GET("/level/:member_level", memberRightsHandler.GetRightsByMemberLevel)
			memberRights.PUT("/:id", memberRightsHandler.UpdateMemberRights)
			memberRights.DELETE("/:id", memberRightsHandler.DeleteMemberRights)
		}

		// 会员积分管理
		pointsRecords := api.Group("/points-records")
		{
			pointsRecords.POST("", memberPointsHandler.CreatePointsRecord)
			pointsRecords.GET("", memberPointsHandler.ListPointsRecords)
		}

		// 操作日志管理
		operationLogs := api.Group("/operation-logs")
		{
			operationLogs.POST("", operationLogHandler.CreateOperationLog)
			operationLogs.GET("", operationLogHandler.ListOperationLogs)
		}
	}

	// 健康检查（在 API 路由组内）
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 根路径健康检查（兼容）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
