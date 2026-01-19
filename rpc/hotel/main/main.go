package main

import (
	"context"
	_ "example_shop/common/init"
	"example_shop/kitex_gen/hotel"
	"example_shop/kitex_gen/hotel/hotelservice"
	hotelHandler "example_shop/rpc/hotel"
	"log"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

// HotelServiceImpl 实现 HotelService 接口
type HotelServiceImpl struct{}

// 注意：这里需要实现所有在 IDL 中定义的方法
// 实际的业务逻辑在 handler 中

// CreateRoomType 创建房型字典
func (s *HotelServiceImpl) CreateRoomType(ctx context.Context, req *hotel.CreateRoomTypeReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.CreateRoomType(ctx, req)
}

// UpdateRoomType 更新房型字典
func (s *HotelServiceImpl) UpdateRoomType(ctx context.Context, req *hotel.UpdateRoomTypeReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.UpdateRoomType(ctx, req)
}

// GetRoomType 获取房型详情
func (s *HotelServiceImpl) GetRoomType(ctx context.Context, id int64) (*hotel.RoomType, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetRoomType(ctx, id)
}

// ListRoomTypes 获取房型列表
func (s *HotelServiceImpl) ListRoomTypes(ctx context.Context, req *hotel.ListRoomTypeReq) (*hotel.ListRoomTypeResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.ListRoomTypes(ctx, req)
}

// DeleteRoomType 删除房型字典
func (s *HotelServiceImpl) DeleteRoomType(ctx context.Context, id int64) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.DeleteRoomType(ctx, id)
}

// CreateRoomInfo 创建房源信息
func (s *HotelServiceImpl) CreateRoomInfo(ctx context.Context, req *hotel.CreateRoomInfoReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.CreateRoomInfo(ctx, req)
}

// UpdateRoomInfo 更新房源信息
func (s *HotelServiceImpl) UpdateRoomInfo(ctx context.Context, req *hotel.UpdateRoomInfoReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.UpdateRoomInfo(ctx, req)
}

// GetRoomInfo 获取房源详情
func (s *HotelServiceImpl) GetRoomInfo(ctx context.Context, id int64) (*hotel.RoomInfo, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetRoomInfo(ctx, id)
}

// ListRoomInfos 获取房源列表
func (s *HotelServiceImpl) ListRoomInfos(ctx context.Context, req *hotel.ListRoomInfoReq) (*hotel.ListRoomInfoResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.ListRoomInfos(ctx, req)
}

// DeleteRoomInfo 删除房源信息
func (s *HotelServiceImpl) DeleteRoomInfo(ctx context.Context, id int64) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.DeleteRoomInfo(ctx, id)
}

// UpdateRoomStatus 更新房源状态
func (s *HotelServiceImpl) UpdateRoomStatus(ctx context.Context, req *hotel.UpdateRoomStatusReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.UpdateRoomStatus(ctx, req)
}

// BatchUpdateRoomStatus 批量更新房源状态
func (s *HotelServiceImpl) BatchUpdateRoomStatus(ctx context.Context, req *hotel.BatchUpdateRoomStatusReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.BatchUpdateRoomStatus(ctx, req)
}

// CreateRoomBinding 创建关联房绑定
func (s *HotelServiceImpl) CreateRoomBinding(ctx context.Context, req *hotel.CreateRoomBindingReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.CreateRoomBinding(ctx, req)
}

// BatchCreateRoomBindings 批量创建关联房绑定
func (s *HotelServiceImpl) BatchCreateRoomBindings(ctx context.Context, req *hotel.BatchCreateRoomBindingsReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.BatchCreateRoomBindings(ctx, req)
}

// GetRoomBindings 获取房源的关联房列表
func (s *HotelServiceImpl) GetRoomBindings(ctx context.Context, roomId int64) (*hotel.ListRoomBindingsResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetRoomBindings(ctx, roomId)
}

// DeleteRoomBinding 删除关联房绑定
func (s *HotelServiceImpl) DeleteRoomBinding(ctx context.Context, bindingId int64) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.DeleteRoomBinding(ctx, bindingId)
}

// GetRoomImages 获取房源图片列表
func (s *HotelServiceImpl) GetRoomImages(ctx context.Context, roomId int64) (*hotel.ListRoomImagesResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetRoomImages(ctx, roomId)
}

// DeleteRoomImage 删除房源图片
func (s *HotelServiceImpl) DeleteRoomImage(ctx context.Context, imageId int64) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.DeleteRoomImage(ctx, imageId)
}

// UpdateImageSortOrder 更新图片排序
func (s *HotelServiceImpl) UpdateImageSortOrder(ctx context.Context, req *hotel.UpdateImageSortOrderReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.UpdateImageSortOrder(ctx, req)
}

// BatchUpdateImageSortOrder 批量更新图片排序
func (s *HotelServiceImpl) BatchUpdateImageSortOrder(ctx context.Context, req *hotel.BatchUpdateImageSortOrderReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.BatchUpdateImageSortOrder(ctx, req)
}

// CreateFacility 创建设施字典
func (s *HotelServiceImpl) CreateFacility(ctx context.Context, req *hotel.CreateFacilityReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.CreateFacility(ctx, req)
}

// UpdateFacility 更新设施字典
func (s *HotelServiceImpl) UpdateFacility(ctx context.Context, req *hotel.UpdateFacilityReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.UpdateFacility(ctx, req)
}

// GetFacility 获取设施详情
func (s *HotelServiceImpl) GetFacility(ctx context.Context, id int64) (*hotel.Facility, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetFacility(ctx, id)
}

// ListFacilities 获取设施列表
func (s *HotelServiceImpl) ListFacilities(ctx context.Context, req *hotel.ListFacilityReq) (*hotel.ListFacilitiesResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.ListFacilities(ctx, req)
}

// DeleteFacility 删除设施字典
func (s *HotelServiceImpl) DeleteFacility(ctx context.Context, id int64) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.DeleteFacility(ctx, id)
}

// SetRoomFacilities 设置房源的设施
func (s *HotelServiceImpl) SetRoomFacilities(ctx context.Context, req *hotel.SetRoomFacilitiesReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.SetRoomFacilities(ctx, req)
}

// GetRoomFacilities 获取房源的设施列表
func (s *HotelServiceImpl) GetRoomFacilities(ctx context.Context, roomId int64) (*hotel.ListFacilitiesResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetRoomFacilities(ctx, roomId)
}

// AddRoomFacility 为房源添加单个设施
func (s *HotelServiceImpl) AddRoomFacility(ctx context.Context, req *hotel.AddRoomFacilityReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.AddRoomFacility(ctx, req)
}

// RemoveRoomFacility 移除房源的单个设施
func (s *HotelServiceImpl) RemoveRoomFacility(ctx context.Context, req *hotel.RemoveRoomFacilityReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.RemoveRoomFacility(ctx, req)
}

// CreateCancellationPolicy 创建退订政策
func (s *HotelServiceImpl) CreateCancellationPolicy(ctx context.Context, req *hotel.CreateCancellationPolicyReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.CreateCancellationPolicy(ctx, req)
}

// UpdateCancellationPolicy 更新退订政策
func (s *HotelServiceImpl) UpdateCancellationPolicy(ctx context.Context, req *hotel.UpdateCancellationPolicyReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.UpdateCancellationPolicy(ctx, req)
}

// GetCancellationPolicy 获取退订政策详情
func (s *HotelServiceImpl) GetCancellationPolicy(ctx context.Context, id int64) (*hotel.CancellationPolicy, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetCancellationPolicy(ctx, id)
}

// ListCancellationPolicies 获取退订政策列表
func (s *HotelServiceImpl) ListCancellationPolicies(ctx context.Context, req *hotel.ListCancellationPolicyReq) (*hotel.ListCancellationPoliciesResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.ListCancellationPolicies(ctx, req)
}

// DeleteCancellationPolicy 删除退订政策
func (s *HotelServiceImpl) DeleteCancellationPolicy(ctx context.Context, id int64) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.DeleteCancellationPolicy(ctx, id)
}

// GetCalendarRoomStatus 获取日历化房态
func (s *HotelServiceImpl) GetCalendarRoomStatus(ctx context.Context, req *hotel.CalendarRoomStatusReq) (*hotel.CalendarRoomStatusResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetCalendarRoomStatus(ctx, req)
}

// UpdateCalendarRoomStatus 更新日历化房态
func (s *HotelServiceImpl) UpdateCalendarRoomStatus(ctx context.Context, req *hotel.UpdateCalendarRoomStatusReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.UpdateCalendarRoomStatus(ctx, req)
}

// BatchUpdateCalendarRoomStatus 批量更新日历化房态
func (s *HotelServiceImpl) BatchUpdateCalendarRoomStatus(ctx context.Context, req *hotel.BatchUpdateCalendarRoomStatusReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.BatchUpdateCalendarRoomStatus(ctx, req)
}

// GetRealTimeStatistics 获取实时数据统计
func (s *HotelServiceImpl) GetRealTimeStatistics(ctx context.Context, req *hotel.RealTimeStatisticsReq) (*hotel.RealTimeStatisticsResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetRealTimeStatistics(ctx, req)
}

// ListBranches 获取分店列表
func (s *HotelServiceImpl) ListBranches(ctx context.Context, req *hotel.ListBranchesReq) (*hotel.ListBranchesResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.ListBranches(ctx, req)
}

// GetBranch 获取分店详情
func (s *HotelServiceImpl) GetBranch(ctx context.Context, branchId int64) (*hotel.Branch, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetBranch(ctx, branchId)
}

// SyncRoomStatusToChannel 同步房态数据到渠道
func (s *HotelServiceImpl) SyncRoomStatusToChannel(ctx context.Context, req *hotel.SyncRoomStatusToChannelReq) (*hotel.SyncRoomStatusToChannelResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.SyncRoomStatusToChannel(ctx, req)
}

// ListOrders 获取订单列表
func (s *HotelServiceImpl) ListOrders(ctx context.Context, req *hotel.ListOrdersReq) (*hotel.ListOrdersResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.ListOrders(ctx, req)
}

// GetOrder 获取订单详情
func (s *HotelServiceImpl) GetOrder(ctx context.Context, orderId int64) (*hotel.Order, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetOrder(ctx, orderId)
}

// ListInHouseGuests 获取在住客人列表
func (s *HotelServiceImpl) ListInHouseGuests(ctx context.Context, req *hotel.ListInHouseGuestsReq) (*hotel.ListInHouseGuestsResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.ListInHouseGuests(ctx, req)
}

// ListFinancialFlows 获取收支流水列表
func (s *HotelServiceImpl) ListFinancialFlows(ctx context.Context, req *hotel.ListFinancialFlowsReq) (*hotel.ListFinancialFlowsResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.ListFinancialFlows(ctx, req)
}

// CreateUserAccount 创建账号
func (s *HotelServiceImpl) CreateUserAccount(ctx context.Context, req *hotel.CreateUserAccountReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.CreateUserAccount(ctx, req)
}

// UpdateUserAccount 更新账号
func (s *HotelServiceImpl) UpdateUserAccount(ctx context.Context, req *hotel.UpdateUserAccountReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.UpdateUserAccount(ctx, req)
}

// GetUserAccount 获取账号详情
func (s *HotelServiceImpl) GetUserAccount(ctx context.Context, id int64) (*hotel.UserAccount, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetUserAccount(ctx, id)
}

// ListUserAccounts 获取账号列表
func (s *HotelServiceImpl) ListUserAccounts(ctx context.Context, req *hotel.ListUserAccountsReq) (*hotel.ListUserAccountsResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.ListUserAccounts(ctx, req)
}

// DeleteUserAccount 删除账号
func (s *HotelServiceImpl) DeleteUserAccount(ctx context.Context, id int64) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.DeleteUserAccount(ctx, id)
}

// CreateRole 创建角色
func (s *HotelServiceImpl) CreateRole(ctx context.Context, req *hotel.CreateRoleReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.CreateRole(ctx, req)
}

// UpdateRole 更新角色
func (s *HotelServiceImpl) UpdateRole(ctx context.Context, req *hotel.UpdateRoleReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.UpdateRole(ctx, req)
}

// GetRole 获取角色详情
func (s *HotelServiceImpl) GetRole(ctx context.Context, id int64) (*hotel.Role, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetRole(ctx, id)
}

// ListRoles 获取角色列表
func (s *HotelServiceImpl) ListRoles(ctx context.Context, req *hotel.ListRolesReq) (*hotel.ListRolesResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.ListRoles(ctx, req)
}

// DeleteRole 删除角色
func (s *HotelServiceImpl) DeleteRole(ctx context.Context, id int64) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.DeleteRole(ctx, id)
}

// ListPermissions 获取权限列表
func (s *HotelServiceImpl) ListPermissions(ctx context.Context, req *hotel.ListPermissionsReq) (*hotel.ListPermissionsResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.ListPermissions(ctx, req)
}

// CreateChannelConfig 创建渠道配置
func (s *HotelServiceImpl) CreateChannelConfig(ctx context.Context, req *hotel.CreateChannelConfigReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.CreateChannelConfig(ctx, req)
}

// UpdateChannelConfig 更新渠道配置
func (s *HotelServiceImpl) UpdateChannelConfig(ctx context.Context, req *hotel.UpdateChannelConfigReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.UpdateChannelConfig(ctx, req)
}

// GetChannelConfig 获取渠道配置详情
func (s *HotelServiceImpl) GetChannelConfig(ctx context.Context, id int64) (*hotel.ChannelConfig, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetChannelConfig(ctx, id)
}

// ListChannelConfigs 获取渠道配置列表
func (s *HotelServiceImpl) ListChannelConfigs(ctx context.Context, req *hotel.ListChannelConfigsReq) (*hotel.ListChannelConfigsResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.ListChannelConfigs(ctx, req)
}

// DeleteChannelConfig 删除渠道配置
func (s *HotelServiceImpl) DeleteChannelConfig(ctx context.Context, id int64) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.DeleteChannelConfig(ctx, id)
}

// CreateSystemConfig 创建系统配置
func (s *HotelServiceImpl) CreateSystemConfig(ctx context.Context, req *hotel.CreateSystemConfigReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.CreateSystemConfig(ctx, req)
}

// UpdateSystemConfig 更新系统配置
func (s *HotelServiceImpl) UpdateSystemConfig(ctx context.Context, req *hotel.UpdateSystemConfigReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.UpdateSystemConfig(ctx, req)
}

// GetSystemConfig 获取系统配置详情
func (s *HotelServiceImpl) GetSystemConfig(ctx context.Context, id int64) (*hotel.SystemConfig, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetSystemConfig(ctx, id)
}

// ListSystemConfigs 获取系统配置列表
func (s *HotelServiceImpl) ListSystemConfigs(ctx context.Context, req *hotel.ListSystemConfigsReq) (*hotel.ListSystemConfigsResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.ListSystemConfigs(ctx, req)
}

// DeleteSystemConfig 删除系统配置
func (s *HotelServiceImpl) DeleteSystemConfig(ctx context.Context, id int64) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.DeleteSystemConfig(ctx, id)
}

// GetSystemConfigsByCategory 按分类获取系统配置
func (s *HotelServiceImpl) GetSystemConfigsByCategory(ctx context.Context, category string) (*hotel.ListSystemConfigsResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetSystemConfigsByCategory(ctx, category)
}

// CreateBlacklist 创建黑名单
func (s *HotelServiceImpl) CreateBlacklist(ctx context.Context, req *hotel.CreateBlacklistReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.CreateBlacklist(ctx, req)
}

// UpdateBlacklist 更新黑名单
func (s *HotelServiceImpl) UpdateBlacklist(ctx context.Context, req *hotel.UpdateBlacklistReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.UpdateBlacklist(ctx, req)
}

// GetBlacklist 获取黑名单详情
func (s *HotelServiceImpl) GetBlacklist(ctx context.Context, id int64) (*hotel.Blacklist, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetBlacklist(ctx, id)
}

// ListBlacklists 获取黑名单列表
func (s *HotelServiceImpl) ListBlacklists(ctx context.Context, req *hotel.ListBlacklistsReq) (*hotel.ListBlacklistsResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.ListBlacklists(ctx, req)
}

// DeleteBlacklist 删除黑名单
func (s *HotelServiceImpl) DeleteBlacklist(ctx context.Context, id int64) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.DeleteBlacklist(ctx, id)
}

// CreateMember 创建会员
func (s *HotelServiceImpl) CreateMember(ctx context.Context, req *hotel.CreateMemberReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.CreateMember(ctx, req)
}

// UpdateMember 更新会员
func (s *HotelServiceImpl) UpdateMember(ctx context.Context, req *hotel.UpdateMemberReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.UpdateMember(ctx, req)
}

// GetMember 获取会员详情
func (s *HotelServiceImpl) GetMember(ctx context.Context, id int64) (*hotel.Member, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetMember(ctx, id)
}

// GetMemberByGuestID 根据客人ID获取会员信息
func (s *HotelServiceImpl) GetMemberByGuestID(ctx context.Context, guestId int64) (*hotel.Member, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetMemberByGuestID(ctx, guestId)
}

// ListMembers 获取会员列表
func (s *HotelServiceImpl) ListMembers(ctx context.Context, req *hotel.ListMembersReq) (*hotel.ListMembersResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.ListMembers(ctx, req)
}

// DeleteMember 删除会员
func (s *HotelServiceImpl) DeleteMember(ctx context.Context, id int64) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.DeleteMember(ctx, id)
}

// CreateMemberRights 创建会员权益
func (s *HotelServiceImpl) CreateMemberRights(ctx context.Context, req *hotel.CreateMemberRightsReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.CreateMemberRights(ctx, req)
}

// UpdateMemberRights 更新会员权益
func (s *HotelServiceImpl) UpdateMemberRights(ctx context.Context, req *hotel.UpdateMemberRightsReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.UpdateMemberRights(ctx, req)
}

// GetMemberRights 获取会员权益详情
func (s *HotelServiceImpl) GetMemberRights(ctx context.Context, id int64) (*hotel.MemberRights, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetMemberRights(ctx, id)
}

// ListMemberRights 获取会员权益列表
func (s *HotelServiceImpl) ListMemberRights(ctx context.Context, req *hotel.ListMemberRightsReq) (*hotel.ListMemberRightsResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.ListMemberRights(ctx, req)
}

// GetRightsByMemberLevel 根据会员等级获取权益列表
func (s *HotelServiceImpl) GetRightsByMemberLevel(ctx context.Context, memberLevel string) (*hotel.ListMemberRightsResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetRightsByMemberLevel(ctx, memberLevel)
}

// DeleteMemberRights 删除会员权益
func (s *HotelServiceImpl) DeleteMemberRights(ctx context.Context, id int64) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.DeleteMemberRights(ctx, id)
}

// CreatePointsRecord 创建积分记录
func (s *HotelServiceImpl) CreatePointsRecord(ctx context.Context, req *hotel.CreatePointsRecordReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.CreatePointsRecord(ctx, req)
}

// ListPointsRecords 获取积分记录列表
func (s *HotelServiceImpl) ListPointsRecords(ctx context.Context, req *hotel.ListPointsRecordsReq) (*hotel.ListPointsRecordsResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.ListPointsRecords(ctx, req)
}

// GetMemberPointsBalance 获取会员积分余额
func (s *HotelServiceImpl) GetMemberPointsBalance(ctx context.Context, memberId int64) (int64, error) {
	handler := hotelHandler.NewHotelService()
	return handler.GetMemberPointsBalance(ctx, memberId)
}

// CreateOperationLog 创建操作日志
func (s *HotelServiceImpl) CreateOperationLog(ctx context.Context, req *hotel.CreateOperationLogReq) (*hotel.BaseResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.CreateOperationLog(ctx, req)
}

// ListOperationLogs 查询操作日志列表
func (s *HotelServiceImpl) ListOperationLogs(ctx context.Context, req *hotel.ListOperationLogsReq) (*hotel.ListOperationLogsResp, error) {
	handler := hotelHandler.NewHotelService()
	return handler.ListOperationLogs(ctx, req)
}

func main() {
	// 启动 Kitex 服务
	svr := hotelservice.NewServer(
		new(HotelServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "hotel_service",
		}),
	)

	log.Println("✅ 酒店管理服务启动成功！")
	log.Println("服务名称: hotel_service")

	if err := svr.Run(); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
