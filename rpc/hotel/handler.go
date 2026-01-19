package hotel

import (
	"context"
	"example_shop/common/model/hotel_admin"
	"example_shop/common/service"
	"example_shop/kitex_gen/hotel"
	"time"
)

type HotelService struct {
	RoomTypeService           *service.RoomTypeService
	RoomInfoService           *service.RoomInfoService
	RoomStatusService         *service.RoomStatusService
	RoomBindingService        *service.RoomBindingService
	RoomImageService          *service.RoomImageService
	FacilityService           *service.FacilityService
	RoomFacilityService       *service.RoomFacilityService
	CancellationPolicyService *service.CancellationPolicyService
	BranchService             *service.BranchService
	ChannelSyncService        *service.ChannelSyncService
	OrderService              *service.OrderService
	GuestService              *service.GuestService
	FinancialService          *service.FinancialService
	UserAccountService        *service.UserAccountService
	RoleService               *service.RoleService
	PermissionService         *service.PermissionService
	ChannelConfigService      *service.ChannelConfigService
	SystemConfigService       *service.SystemConfigService
	BlacklistService          *service.BlacklistService
	MemberService             *service.MemberService
	MemberRightsService       *service.MemberRightsService
	MemberPointsService       *service.MemberPointsService
	OperationLogService       *service.OperationLogService
}

// NewHotelService 创建酒店服务实例
func NewHotelService() *HotelService {
	return &HotelService{
		RoomTypeService:           &service.RoomTypeService{},
		RoomInfoService:           &service.RoomInfoService{},
		RoomStatusService:         &service.RoomStatusService{},
		RoomBindingService:        &service.RoomBindingService{},
		RoomImageService:          &service.RoomImageService{},
		FacilityService:           &service.FacilityService{},
		RoomFacilityService:       &service.RoomFacilityService{},
		CancellationPolicyService: &service.CancellationPolicyService{},
		BranchService:             &service.BranchService{},
		ChannelSyncService:        &service.ChannelSyncService{},
		OrderService:              &service.OrderService{},
		GuestService:              &service.GuestService{},
		FinancialService:          &service.FinancialService{},
		UserAccountService:        &service.UserAccountService{},
		RoleService:               &service.RoleService{},
		PermissionService:         &service.PermissionService{},
		ChannelConfigService:      &service.ChannelConfigService{},
		SystemConfigService:       &service.SystemConfigService{},
		BlacklistService:          &service.BlacklistService{},
		MemberService:             &service.MemberService{},
		MemberRightsService:       &service.MemberRightsService{},
		MemberPointsService:       &service.MemberPointsService{},
		OperationLogService:        &service.OperationLogService{},
	}
}

// CreateRoomType 创建房型字典
func (s *HotelService) CreateRoomType(ctx context.Context, req *hotel.CreateRoomTypeReq) (resp *hotel.BaseResp, err error) {
	serviceReq := &service.CreateRoomTypeReq{
		RoomTypeName:  req.RoomTypeName,
		BedSpec:       req.BedSpec,
		Area:          req.Area,
		HasBreakfast:  req.HasBreakfast,
		HasToiletries: req.HasToiletries,
		DefaultPrice:  req.DefaultPrice,
	}

	_, err = s.RoomTypeService.CreateRoomType(serviceReq)
	if err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "创建成功",
	}, nil
}

// UpdateRoomType 更新房型字典
func (s *HotelService) UpdateRoomType(ctx context.Context, req *hotel.UpdateRoomTypeReq) (resp *hotel.BaseResp, err error) {
	var roomTypeName *string
	var bedSpec string
	var hasBreakfast *bool
	var hasToiletries *bool
	var defaultPrice float64
	var status string

	if req.RoomTypeName != nil {
		roomTypeName = req.RoomTypeName
	}
	if req.BedSpec != nil {
		bedSpec = *req.BedSpec
	}
	if req.HasBreakfast != nil {
		hasBreakfast = req.HasBreakfast
	}
	if req.HasToiletries != nil {
		hasToiletries = req.HasToiletries
	}
	if req.DefaultPrice != nil {
		defaultPrice = *req.DefaultPrice
	}
	if req.Status != nil {
		status = *req.Status
	}

	serviceReq := &service.UpdateRoomTypeReq{
		RoomTypeName:  roomTypeName,
		BedSpec:       bedSpec,
		Area:          req.Area,
		HasBreakfast:  hasBreakfast,
		HasToiletries: hasToiletries,
		DefaultPrice:  defaultPrice,
		Status:        status,
	}

	err = s.RoomTypeService.UpdateRoomType(uint64(req.Id), serviceReq)
	if err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "更新成功",
	}, nil
}

// GetRoomType 获取房型详情
func (s *HotelService) GetRoomType(ctx context.Context, id int64) (resp *hotel.RoomType, err error) {
	roomType, err := s.RoomTypeService.GetRoomType(uint64(id))
	if err != nil {
		return nil, err
	}

	return convertRoomTypeToThrift(roomType), nil
}

// ListRoomTypes 获取房型列表
func (s *HotelService) ListRoomTypes(ctx context.Context, req *hotel.ListRoomTypeReq) (resp *hotel.ListRoomTypeResp, err error) {
	var status, keyword string
	if req.Status != nil {
		status = *req.Status
	}
	if req.Keyword != nil {
		keyword = *req.Keyword
	}

	serviceReq := &service.ListRoomTypeReq{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
		Status:   status,
		Keyword:  keyword,
	}

	roomTypes, total, err := s.RoomTypeService.ListRoomTypes(serviceReq)
	if err != nil {
		return nil, err
	}

	thriftRoomTypes := make([]*hotel.RoomType, 0, len(roomTypes))
	for _, rt := range roomTypes {
		thriftRoomTypes = append(thriftRoomTypes, convertRoomTypeToThrift(&rt))
	}

	return &hotel.ListRoomTypeResp{
		List:     thriftRoomTypes,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// DeleteRoomType 删除房型字典
func (s *HotelService) DeleteRoomType(ctx context.Context, id int64) (resp *hotel.BaseResp, err error) {
	err = s.RoomTypeService.DeleteRoomType(uint64(id))
	if err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "删除成功",
	}, nil
}

// CreateRoomInfo 创建房源信息
func (s *HotelService) CreateRoomInfo(ctx context.Context, req *hotel.CreateRoomInfoReq) (resp *hotel.BaseResp, err error) {
	var cancellationPolicyID *uint64
	if req.CancellationPolicyId != nil {
		id := uint64(*req.CancellationPolicyId)
		cancellationPolicyID = &id
	}

	serviceReq := &service.CreateRoomInfoReq{
		BranchID:             uint64(req.BranchId),
		RoomTypeID:           uint64(req.RoomTypeId),
		RoomNo:               req.RoomNo,
		RoomName:             req.RoomName,
		MarketPrice:          req.MarketPrice,
		CalendarPrice:        req.CalendarPrice,
		RoomCount:            uint8(req.RoomCount),
		Area:                 req.Area,
		BedSpec:              req.BedSpec,
		HasBreakfast:         req.HasBreakfast,
		HasToiletries:        req.HasToiletries,
		CancellationPolicyID: cancellationPolicyID,
		CreatedBy:            uint64(req.CreatedBy),
	}

	_, err = s.RoomInfoService.CreateRoomInfo(serviceReq)
	if err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "创建成功",
	}, nil
}

// UpdateRoomInfo 更新房源信息
func (s *HotelService) UpdateRoomInfo(ctx context.Context, req *hotel.UpdateRoomInfoReq) (resp *hotel.BaseResp, err error) {
	var roomNo, roomName, bedSpec, status string
	var roomCount uint8
	var hasBreakfast *bool
	var hasToiletries *bool
	var marketPrice, calendarPrice float64
	var cancellationPolicyID *uint64

	if req.RoomNo != nil {
		roomNo = *req.RoomNo
	}
	if req.RoomName != nil {
		roomName = *req.RoomName
	}
	if req.BedSpec != nil {
		bedSpec = *req.BedSpec
	}
	if req.Status != nil {
		status = *req.Status
	}
	if req.RoomCount != nil {
		// Thrift 中的 i8 是 int8，需要转换为 uint8
		if *req.RoomCount > 0 {
			roomCount = uint8(*req.RoomCount)
		}
	}
	if req.HasBreakfast != nil {
		hasBreakfast = req.HasBreakfast
	}
	if req.HasToiletries != nil {
		hasToiletries = req.HasToiletries
	}
	if req.MarketPrice != nil {
		marketPrice = *req.MarketPrice
	}
	if req.CalendarPrice != nil {
		calendarPrice = *req.CalendarPrice
	}
	if req.CancellationPolicyId != nil {
		id := uint64(*req.CancellationPolicyId)
		cancellationPolicyID = &id
	}

	serviceReq := &service.UpdateRoomInfoReq{
		RoomNo:               roomNo,
		RoomName:             roomName,
		MarketPrice:          marketPrice,
		CalendarPrice:        calendarPrice,
		RoomCount:            roomCount,
		Area:                 req.Area,
		BedSpec:              bedSpec,
		HasBreakfast:         hasBreakfast,
		HasToiletries:        hasToiletries,
		CancellationPolicyID: cancellationPolicyID,
		Status:               status,
	}

	err = s.RoomInfoService.UpdateRoomInfo(uint64(req.Id), serviceReq)
	if err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "更新成功",
	}, nil
}

// GetRoomInfo 获取房源详情
func (s *HotelService) GetRoomInfo(ctx context.Context, id int64) (resp *hotel.RoomInfo, err error) {
	roomInfo, err := s.RoomInfoService.GetRoomInfo(uint64(id))
	if err != nil {
		return nil, err
	}

	return convertRoomInfoToThrift(roomInfo), nil
}

// ListRoomInfos 获取房源列表
func (s *HotelService) ListRoomInfos(ctx context.Context, req *hotel.ListRoomInfoReq) (resp *hotel.ListRoomInfoResp, err error) {
	var branchID, roomTypeID uint64
	var status, keyword string

	if req.BranchId != nil {
		branchID = uint64(*req.BranchId)
	}
	if req.RoomTypeId != nil {
		roomTypeID = uint64(*req.RoomTypeId)
	}
	if req.Status != nil {
		status = *req.Status
	}
	if req.Keyword != nil {
		keyword = *req.Keyword
	}

	serviceReq := &service.ListRoomInfoReq{
		Page:       int(req.Page),
		PageSize:   int(req.PageSize),
		BranchID:   branchID,
		RoomTypeID: roomTypeID,
		Status:     status,
		Keyword:    keyword,
	}

	roomInfos, total, err := s.RoomInfoService.ListRoomInfos(serviceReq)
	if err != nil {
		return nil, err
	}

	thriftRoomInfos := make([]*hotel.RoomInfo, 0, len(roomInfos))
	for _, ri := range roomInfos {
		thriftRoomInfos = append(thriftRoomInfos, convertRoomInfoToThrift(&ri))
	}

	return &hotel.ListRoomInfoResp{
		List:     thriftRoomInfos,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// DeleteRoomInfo 删除房源信息
func (s *HotelService) DeleteRoomInfo(ctx context.Context, id int64) (resp *hotel.BaseResp, err error) {
	err = s.RoomInfoService.DeleteRoomInfo(uint64(id))
	if err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "删除成功",
	}, nil
}

// convertRoomTypeToThrift 转换房型模型到 Thrift 结构
func convertRoomTypeToThrift(rt *hotel_admin.RoomTypeDict) *hotel.RoomType {
	return &hotel.RoomType{
		Id:            int64(rt.ID),
		RoomTypeName:  rt.RoomTypeName,
		BedSpec:       rt.BedSpec,
		Area:          rt.Area,
		HasBreakfast:  rt.HasBreakfast,
		HasToiletries: rt.HasToiletries,
		DefaultPrice:  rt.DefaultPrice,
		Status:        rt.Status,
		CreatedAt:     rt.CreatedAt.Format(time.DateTime),
		UpdatedAt:     rt.UpdatedAt.Format(time.DateTime),
	}
}

// convertRoomInfoToThrift 转换房源模型到 Thrift 结构
func convertRoomInfoToThrift(ri *hotel_admin.RoomInfo) *hotel.RoomInfo {
	roomInfo := &hotel.RoomInfo{
		Id:            int64(ri.ID),
		BranchId:      int64(ri.BranchID),
		RoomTypeId:    int64(ri.RoomTypeID),
		RoomNo:        ri.RoomNo,
		RoomName:      ri.RoomName,
		MarketPrice:   ri.MarketPrice,
		CalendarPrice: ri.CalendarPrice,
		RoomCount:     int8(ri.RoomCount),
		Area:          ri.Area,
		BedSpec:       ri.BedSpec,
		HasBreakfast:  ri.HasBreakfast,
		HasToiletries: ri.HasToiletries,
		Status:        ri.Status,
		CreatedAt:     ri.CreatedAt.Format(time.DateTime),
		UpdatedAt:     ri.UpdatedAt.Format(time.DateTime),
	}
	if ri.CancellationPolicyID != nil {
		id := int64(*ri.CancellationPolicyID)
		roomInfo.CancellationPolicyId = &id
	}
	return roomInfo
}

// UpdateRoomStatus 更新房源状态
func (s *HotelService) UpdateRoomStatus(ctx context.Context, req *hotel.UpdateRoomStatusReq) (resp *hotel.BaseResp, err error) {
	if err := s.RoomStatusService.UpdateRoomStatus(uint64(req.RoomId), req.Status); err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "状态更新成功",
	}, nil
}

// BatchUpdateRoomStatus 批量更新房源状态
func (s *HotelService) BatchUpdateRoomStatus(ctx context.Context, req *hotel.BatchUpdateRoomStatusReq) (resp *hotel.BaseResp, err error) {
	roomIDs := make([]uint64, len(req.RoomIds))
	for i, id := range req.RoomIds {
		roomIDs[i] = uint64(id)
	}

	if err := s.RoomStatusService.BatchUpdateRoomStatus(roomIDs, req.Status); err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "批量更新成功",
	}, nil
}

// CreateRoomBinding 创建关联房绑定
func (s *HotelService) CreateRoomBinding(ctx context.Context, req *hotel.CreateRoomBindingReq) (resp *hotel.BaseResp, err error) {
	var bindingDesc *string
	if req.BindingDesc != nil {
		bindingDesc = req.BindingDesc
	}

	if err := s.RoomBindingService.CreateRoomBinding(uint64(req.MainRoomId), uint64(req.RelatedRoomId), bindingDesc); err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "关联房绑定成功",
	}, nil
}

// BatchCreateRoomBindings 批量创建关联房绑定
func (s *HotelService) BatchCreateRoomBindings(ctx context.Context, req *hotel.BatchCreateRoomBindingsReq) (resp *hotel.BaseResp, err error) {
	relatedRoomIDs := make([]uint64, len(req.RelatedRoomIds))
	for i, id := range req.RelatedRoomIds {
		relatedRoomIDs[i] = uint64(id)
	}

	var bindingDesc *string
	if req.BindingDesc != nil {
		bindingDesc = req.BindingDesc
	}

	if err := s.RoomBindingService.BatchCreateRoomBindings(uint64(req.MainRoomId), relatedRoomIDs, bindingDesc); err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "批量关联房绑定成功",
	}, nil
}

// GetRoomBindings 获取房源的关联房列表
func (s *HotelService) GetRoomBindings(ctx context.Context, roomId int64) (resp *hotel.ListRoomBindingsResp, err error) {
	bindings, err := s.RoomBindingService.GetRoomBindings(uint64(roomId))
	if err != nil {
		return nil, err
	}

	thriftBindings := make([]*hotel.RelatedRoomBinding, len(bindings))
	for i, b := range bindings {
		var bindingDesc *string
		if b.BindingDesc != nil {
			bindingDesc = b.BindingDesc
		}
		thriftBindings[i] = &hotel.RelatedRoomBinding{
			Id:            int64(b.ID),
			MainRoomId:    int64(b.MainRoomID),
			RelatedRoomId: int64(b.RelatedRoomID),
			BindingDesc:   bindingDesc,
			CreatedAt:     b.CreatedAt.Format(time.DateTime),
		}
	}

	return &hotel.ListRoomBindingsResp{
		Bindings: thriftBindings,
	}, nil
}

// DeleteRoomBinding 删除关联房绑定
func (s *HotelService) DeleteRoomBinding(ctx context.Context, bindingId int64) (resp *hotel.BaseResp, err error) {
	if err := s.RoomBindingService.DeleteRoomBinding(uint64(bindingId)); err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "删除成功",
	}, nil
}

// GetRoomImages 获取房源图片列表
func (s *HotelService) GetRoomImages(ctx context.Context, roomId int64) (resp *hotel.ListRoomImagesResp, err error) {
	images, err := s.RoomImageService.GetRoomImages(uint64(roomId))
	if err != nil {
		return nil, err
	}

	thriftImages := make([]*hotel.RoomImage, len(images))
	for i, img := range images {
		thriftImages[i] = &hotel.RoomImage{
			Id:          int64(img.ID),
			RoomId:      int64(img.RoomID),
			ImageUrl:    img.ImageURL,
			ImageSize:   img.ImageSize,
			ImageFormat: img.ImageFormat,
			SortOrder:   int8(img.SortOrder),
			UploadTime:  img.UploadTime.Format(time.DateTime),
		}
	}

	return &hotel.ListRoomImagesResp{
		Images: thriftImages,
	}, nil
}

// DeleteRoomImage 删除房源图片
func (s *HotelService) DeleteRoomImage(ctx context.Context, imageId int64) (resp *hotel.BaseResp, err error) {
	if err := s.RoomImageService.DeleteRoomImage(uint64(imageId)); err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "删除成功",
	}, nil
}

// UpdateImageSortOrder 更新图片排序
func (s *HotelService) UpdateImageSortOrder(ctx context.Context, req *hotel.UpdateImageSortOrderReq) (resp *hotel.BaseResp, err error) {
	if err := s.RoomImageService.UpdateImageSortOrder(uint64(req.ImageId), uint8(req.SortOrder)); err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "排序更新成功",
	}, nil
}

// BatchUpdateImageSortOrder 批量更新图片排序
func (s *HotelService) BatchUpdateImageSortOrder(ctx context.Context, req *hotel.BatchUpdateImageSortOrderReq) (resp *hotel.BaseResp, err error) {
	sortOrders := make([]service.ImageSortOrder, len(req.SortOrders))
	for i, so := range req.SortOrders {
		sortOrders[i] = service.ImageSortOrder{
			ImageID:   uint64(so.ImageId),
			SortOrder: uint8(so.SortOrder),
		}
	}

	if err := s.RoomImageService.BatchUpdateImageSortOrder(uint64(req.RoomId), sortOrders); err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "批量排序更新成功",
	}, nil
}

// CreateFacility 创建设施字典
func (s *HotelService) CreateFacility(ctx context.Context, req *hotel.CreateFacilityReq) (resp *hotel.BaseResp, err error) {
	serviceReq := &service.CreateFacilityReq{
		FacilityName: req.FacilityName,
		Description:  req.Description,
	}

	_, err = s.FacilityService.CreateFacility(serviceReq)
	if err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "创建成功",
	}, nil
}

// UpdateFacility 更新设施字典
func (s *HotelService) UpdateFacility(ctx context.Context, req *hotel.UpdateFacilityReq) (resp *hotel.BaseResp, err error) {
	serviceReq := &service.UpdateFacilityReq{}
	if req.FacilityName != nil {
		serviceReq.FacilityName = *req.FacilityName
	}
	serviceReq.Description = req.Description
	if req.Status != nil {
		serviceReq.Status = *req.Status
	}

	if err := s.FacilityService.UpdateFacility(uint64(req.Id), serviceReq); err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "更新成功",
	}, nil
}

// GetFacility 获取设施详情
func (s *HotelService) GetFacility(ctx context.Context, id int64) (resp *hotel.Facility, err error) {
	facility, err := s.FacilityService.GetFacility(uint64(id))
	if err != nil {
		return nil, err
	}

	return &hotel.Facility{
		Id:           int64(facility.ID),
		FacilityName: facility.FacilityName,
		Description:  facility.Description,
		Status:       facility.Status,
		CreatedAt:    facility.CreatedAt.Format(time.DateTime),
		UpdatedAt:    facility.UpdatedAt.Format(time.DateTime),
	}, nil
}

// ListFacilities 获取设施列表
func (s *HotelService) ListFacilities(ctx context.Context, req *hotel.ListFacilityReq) (resp *hotel.ListFacilitiesResp, err error) {
	serviceReq := &service.ListFacilityReq{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}
	if req.Status != nil {
		serviceReq.Status = *req.Status
	}
	if req.Keyword != nil {
		serviceReq.Keyword = *req.Keyword
	}

	facilities, total, err := s.FacilityService.ListFacilities(serviceReq)
	if err != nil {
		return nil, err
	}

	thriftFacilities := make([]*hotel.Facility, len(facilities))
	for i, f := range facilities {
		thriftFacilities[i] = &hotel.Facility{
			Id:           int64(f.ID),
			FacilityName: f.FacilityName,
			Description:  f.Description,
			Status:       f.Status,
			CreatedAt:    f.CreatedAt.Format(time.DateTime),
			UpdatedAt:    f.UpdatedAt.Format(time.DateTime),
		}
	}

	return &hotel.ListFacilitiesResp{
		List:     thriftFacilities,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// DeleteFacility 删除设施字典
func (s *HotelService) DeleteFacility(ctx context.Context, id int64) (resp *hotel.BaseResp, err error) {
	if err := s.FacilityService.DeleteFacility(uint64(id)); err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "删除成功",
	}, nil
}

// SetRoomFacilities 设置房源的设施
func (s *HotelService) SetRoomFacilities(ctx context.Context, req *hotel.SetRoomFacilitiesReq) (resp *hotel.BaseResp, err error) {
	facilityIDs := make([]uint64, len(req.FacilityIds))
	for i, id := range req.FacilityIds {
		facilityIDs[i] = uint64(id)
	}

	if err := s.RoomFacilityService.SetRoomFacilities(uint64(req.RoomId), facilityIDs); err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "设置成功",
	}, nil
}

// GetRoomFacilities 获取房源的设施列表
func (s *HotelService) GetRoomFacilities(ctx context.Context, roomId int64) (resp *hotel.ListFacilitiesResp, err error) {
	facilities, err := s.RoomFacilityService.GetRoomFacilities(uint64(roomId))
	if err != nil {
		return nil, err
	}

	thriftFacilities := make([]*hotel.Facility, len(facilities))
	for i, f := range facilities {
		thriftFacilities[i] = &hotel.Facility{
			Id:           int64(f.ID),
			FacilityName: f.FacilityName,
			Description:  f.Description,
			Status:       f.Status,
			CreatedAt:    f.CreatedAt.Format(time.DateTime),
			UpdatedAt:    f.UpdatedAt.Format(time.DateTime),
		}
	}

	return &hotel.ListFacilitiesResp{
		List: thriftFacilities,
	}, nil
}

// AddRoomFacility 为房源添加单个设施
func (s *HotelService) AddRoomFacility(ctx context.Context, req *hotel.AddRoomFacilityReq) (resp *hotel.BaseResp, err error) {
	if err := s.RoomFacilityService.AddRoomFacility(uint64(req.RoomId), uint64(req.FacilityId)); err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "添加成功",
	}, nil
}

// RemoveRoomFacility 移除房源的单个设施
func (s *HotelService) RemoveRoomFacility(ctx context.Context, req *hotel.RemoveRoomFacilityReq) (resp *hotel.BaseResp, err error) {
	if err := s.RoomFacilityService.RemoveRoomFacility(uint64(req.RoomId), uint64(req.FacilityId)); err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "移除成功",
	}, nil
}

// CreateCancellationPolicy 创建退订政策
func (s *HotelService) CreateCancellationPolicy(ctx context.Context, req *hotel.CreateCancellationPolicyReq) (resp *hotel.BaseResp, err error) {
	var roomTypeID *uint64
	if req.RoomTypeId != nil {
		id := uint64(*req.RoomTypeId)
		roomTypeID = &id
	}

	serviceReq := &service.CreateCancellationPolicyReq{
		PolicyName:      req.PolicyName,
		RuleDescription: req.RuleDescription,
		PenaltyRatio:    req.PenaltyRatio,
		RoomTypeID:      roomTypeID,
	}

	_, err = s.CancellationPolicyService.CreateCancellationPolicy(serviceReq)
	if err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "创建成功",
	}, nil
}

// UpdateCancellationPolicy 更新退订政策
func (s *HotelService) UpdateCancellationPolicy(ctx context.Context, req *hotel.UpdateCancellationPolicyReq) (resp *hotel.BaseResp, err error) {
	var roomTypeID *uint64
	if req.RoomTypeId != nil {
		id := uint64(*req.RoomTypeId)
		roomTypeID = &id
	}

	serviceReq := &service.UpdateCancellationPolicyReq{
		RoomTypeID: roomTypeID,
	}
	if req.PolicyName != nil {
		serviceReq.PolicyName = *req.PolicyName
	}
	if req.RuleDescription != nil {
		serviceReq.RuleDescription = *req.RuleDescription
	}
	if req.PenaltyRatio != nil {
		serviceReq.PenaltyRatio = *req.PenaltyRatio
	}
	if req.Status != nil {
		serviceReq.Status = *req.Status
	}

	if err := s.CancellationPolicyService.UpdateCancellationPolicy(uint64(req.Id), serviceReq); err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "更新成功",
	}, nil
}

// GetCancellationPolicy 获取退订政策详情
func (s *HotelService) GetCancellationPolicy(ctx context.Context, id int64) (resp *hotel.CancellationPolicy, err error) {
	policy, err := s.CancellationPolicyService.GetCancellationPolicy(uint64(id))
	if err != nil {
		return nil, err
	}

	thriftPolicy := &hotel.CancellationPolicy{
		Id:              int64(policy.ID),
		PolicyName:      policy.PolicyName,
		RuleDescription: policy.RuleDescription,
		PenaltyRatio:    policy.PenaltyRatio,
		Status:          policy.Status,
		CreatedAt:       policy.CreatedAt.Format(time.DateTime),
		UpdatedAt:       policy.UpdatedAt.Format(time.DateTime),
	}

	if policy.RoomTypeID != nil {
		id := int64(*policy.RoomTypeID)
		thriftPolicy.RoomTypeId = &id
	}

	return thriftPolicy, nil
}

// ListCancellationPolicies 获取退订政策列表
func (s *HotelService) ListCancellationPolicies(ctx context.Context, req *hotel.ListCancellationPolicyReq) (resp *hotel.ListCancellationPoliciesResp, err error) {
	var roomTypeID uint64
	if req.RoomTypeId != nil {
		roomTypeID = uint64(*req.RoomTypeId)
	}

	serviceReq := &service.ListCancellationPolicyReq{
		Page:       int(req.Page),
		PageSize:   int(req.PageSize),
		RoomTypeID: roomTypeID,
	}
	if req.Status != nil {
		serviceReq.Status = *req.Status
	}
	if req.Keyword != nil {
		serviceReq.Keyword = *req.Keyword
	}

	policies, total, err := s.CancellationPolicyService.ListCancellationPolicies(serviceReq)
	if err != nil {
		return nil, err
	}

	thriftPolicies := make([]*hotel.CancellationPolicy, len(policies))
	for i, p := range policies {
		thriftPolicies[i] = &hotel.CancellationPolicy{
			Id:              int64(p.ID),
			PolicyName:      p.PolicyName,
			RuleDescription: p.RuleDescription,
			PenaltyRatio:    p.PenaltyRatio,
			Status:          p.Status,
			CreatedAt:       p.CreatedAt.Format(time.DateTime),
			UpdatedAt:       p.UpdatedAt.Format(time.DateTime),
		}
		if p.RoomTypeID != nil {
			id := int64(*p.RoomTypeID)
			thriftPolicies[i].RoomTypeId = &id
		}
	}

	return &hotel.ListCancellationPoliciesResp{
		List:     thriftPolicies,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// DeleteCancellationPolicy 删除退订政策
func (s *HotelService) DeleteCancellationPolicy(ctx context.Context, id int64) (resp *hotel.BaseResp, err error) {
	if err := s.CancellationPolicyService.DeleteCancellationPolicy(uint64(id)); err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "删除成功",
	}, nil
}

// GetCalendarRoomStatus 获取日历化房态
func (s *HotelService) GetCalendarRoomStatus(ctx context.Context, req *hotel.CalendarRoomStatusReq) (resp *hotel.CalendarRoomStatusResp, err error) {
	// 解析日期
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, err
	}

	// 构建服务层请求
	serviceReq := service.CalendarRoomStatusReq{
		StartDate: startDate,
		EndDate:   endDate,
	}

	if req.BranchId != nil {
		id := uint64(*req.BranchId)
		serviceReq.BranchID = &id
	}

	if req.RoomNo != nil {
		serviceReq.RoomNo = req.RoomNo
	}

	if req.Status != nil {
		serviceReq.Status = req.Status
	}

	// 调用服务层
	items, err := s.RoomStatusService.GetCalendarRoomStatus(serviceReq)
	if err != nil {
		return nil, err
	}

	// 转换为 Thrift 格式
	thriftItems := make([]*hotel.CalendarRoomStatusItem, len(items))
	for i, item := range items {
		thriftItems[i] = &hotel.CalendarRoomStatusItem{
			RoomId:               int64(item.RoomID),
			RoomNo:               item.RoomNo,
			RoomName:             item.RoomName,
			Date:                 item.Date.Format("2006-01-02"),
			RoomStatus:           item.RoomStatus,
			RemainingCount:       int8(item.RemainingCount),
			CheckedInCount:       int8(item.CheckedInCount),
			CheckOutPendingCount: int8(item.CheckOutPendingCount),
			ReservedPendingCount: int8(item.ReservedPendingCount),
		}
	}

	return &hotel.CalendarRoomStatusResp{
		Items: thriftItems,
	}, nil
}

// UpdateCalendarRoomStatus 更新日历化房态
func (s *HotelService) UpdateCalendarRoomStatus(ctx context.Context, req *hotel.UpdateCalendarRoomStatusReq) (resp *hotel.BaseResp, err error) {
	// 解析日期
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return &hotel.BaseResp{
			Code: 400,
			Msg:  "日期格式错误，应为 YYYY-MM-DD",
		}, nil
	}

	// 调用服务层
	if err := s.RoomStatusService.UpdateCalendarRoomStatus(uint64(req.RoomId), date, req.Status); err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "更新成功",
	}, nil
}

// BatchUpdateCalendarRoomStatus 批量更新日历化房态
func (s *HotelService) BatchUpdateCalendarRoomStatus(ctx context.Context, req *hotel.BatchUpdateCalendarRoomStatusReq) (resp *hotel.BaseResp, err error) {
	updates := make([]struct {
		RoomID uint64
		Date   time.Time
		Status string
	}, len(req.Updates))

	for i, update := range req.Updates {
		date, err := time.Parse("2006-01-02", update.Date)
		if err != nil {
			return &hotel.BaseResp{
				Code: 400,
				Msg:  "日期格式错误，应为 YYYY-MM-DD",
			}, nil
		}

		updates[i] = struct {
			RoomID uint64
			Date   time.Time
			Status string
		}{
			RoomID: uint64(update.RoomId),
			Date:   date,
			Status: update.Status,
		}
	}

	// 调用服务层
	if err := s.RoomStatusService.BatchUpdateCalendarRoomStatus(updates); err != nil {
		return &hotel.BaseResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &hotel.BaseResp{
		Code: 200,
		Msg:  "批量更新成功",
	}, nil
}

// GetRealTimeStatistics 获取实时数据统计
func (s *HotelService) GetRealTimeStatistics(ctx context.Context, req *hotel.RealTimeStatisticsReq) (resp *hotel.RealTimeStatisticsResp, err error) {
	// 构建服务层请求
	serviceReq := service.RealTimeStatisticsReq{}

	if req.BranchId != nil {
		id := uint64(*req.BranchId)
		serviceReq.BranchID = &id
	}

	if req.Date != nil {
		serviceReq.Date = req.Date
	}

	if req.RoomNo != nil {
		serviceReq.RoomNo = req.RoomNo
	}

	if req.RoomTypeId != nil {
		id := uint64(*req.RoomTypeId)
		serviceReq.RoomTypeID = &id
	}

	// 调用服务层
	result, err := s.RoomStatusService.GetRealTimeStatistics(serviceReq)
	if err != nil {
		return nil, err
	}

	// 转换为 Thrift 格式
	thriftResp := &hotel.RealTimeStatisticsResp{
		Date:                 result.Date,
		TotalRooms:           int64(result.TotalRooms),
		RemainingRooms:       int64(result.RemainingRooms),
		CheckedInCount:       int64(result.CheckedInCount),
		CheckOutPendingCount: int64(result.CheckOutPendingCount),
		ReservedPendingCount: int64(result.ReservedPendingCount),
		OccupiedRooms:        int64(result.OccupiedRooms),
		MaintenanceRooms:     int64(result.MaintenanceRooms),
		LockedRooms:          int64(result.LockedRooms),
		EmptyRooms:           int64(result.EmptyRooms),
		ReservedRooms:        int64(result.ReservedRooms),
	}

	// 转换房态分组统计
	thriftResp.StatusBreakdown = make([]*hotel.StatusBreakdown, len(result.StatusBreakdown))
	for i, breakdown := range result.StatusBreakdown {
		thriftResp.StatusBreakdown[i] = &hotel.StatusBreakdown{
			Status: breakdown.Status,
			Count:  int64(breakdown.Count),
		}
	}

	// 转换房间明细（如果存在）
	if len(result.RoomDetails) > 0 {
		thriftResp.RoomDetails = make([]*hotel.RoomDetailStat, len(result.RoomDetails))
		for i, detail := range result.RoomDetails {
			thriftResp.RoomDetails[i] = &hotel.RoomDetailStat{
				RoomId:               int64(detail.RoomID),
				RoomNo:               detail.RoomNo,
				RoomName:             detail.RoomName,
				RoomStatus:           detail.RoomStatus,
				RemainingCount:       int8(detail.RemainingCount),
				CheckedInCount:       int8(detail.CheckedInCount),
				CheckOutPendingCount: int8(detail.CheckOutPendingCount),
				ReservedPendingCount: int8(detail.ReservedPendingCount),
			}
		}
	}

	return thriftResp, nil
}

// ListBranches 获取分店列表
func (s *HotelService) ListBranches(ctx context.Context, req *hotel.ListBranchesReq) (resp *hotel.ListBranchesResp, err error) {
	serviceReq := service.ListBranchesReq{}

	if req.Status != nil {
		serviceReq.Status = req.Status
	}

	branches, err := s.BranchService.ListBranches(serviceReq)
	if err != nil {
		return nil, err
	}

	thriftBranches := make([]*hotel.Branch, len(branches))
	for i, branch := range branches {
		thriftBranches[i] = &hotel.Branch{
			Id:           int64(branch.ID),
			HotelName:    branch.HotelName,
			BranchCode:   branch.BranchCode,
			Address:      branch.Address,
			Contact:      branch.Contact,
			ContactPhone: branch.ContactPhone,
			Status:       branch.Status,
			CreatedAt:    branch.CreatedAt,
			UpdatedAt:    branch.UpdatedAt,
		}
	}

	return &hotel.ListBranchesResp{
		Branches: thriftBranches,
	}, nil
}

// GetBranch 获取分店详情
func (s *HotelService) GetBranch(ctx context.Context, branchID int64) (resp *hotel.Branch, err error) {
	branch, err := s.BranchService.GetBranch(uint64(branchID))
	if err != nil {
		return nil, err
	}

	return &hotel.Branch{
		Id:           int64(branch.ID),
		HotelName:    branch.HotelName,
		BranchCode:   branch.BranchCode,
		Address:      branch.Address,
		Contact:      branch.Contact,
		ContactPhone: branch.ContactPhone,
		Status:       branch.Status,
		CreatedAt:    branch.CreatedAt,
		UpdatedAt:    branch.UpdatedAt,
	}, nil
}

// SyncRoomStatusToChannel 同步房态数据到渠道
func (s *HotelService) SyncRoomStatusToChannel(ctx context.Context, req *hotel.SyncRoomStatusToChannelReq) (resp *hotel.SyncRoomStatusToChannelResp, err error) {
	serviceReq := service.SyncRoomStatusReq{
		BranchID:  uint64(req.BranchId),
		ChannelID: uint64(req.ChannelId),
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	if len(req.RoomIds) > 0 {
		roomIDs := make([]uint64, len(req.RoomIds))
		for i, id := range req.RoomIds {
			roomIDs[i] = uint64(id)
		}
		serviceReq.RoomIDs = roomIDs
	}

	result, err := s.ChannelSyncService.SyncRoomStatusToChannel(serviceReq)
	if err != nil {
		return nil, err
	}

	syncLogIDs := make([]int64, len(result.SyncLogs))
	for i, id := range result.SyncLogs {
		syncLogIDs[i] = int64(id)
	}

	return &hotel.SyncRoomStatusToChannelResp{
		SuccessCount: int32(result.SuccessCount),
		FailCount:    int32(result.FailCount),
		SyncLogs:     syncLogIDs,
	}, nil
}

// ListOrders 获取订单列表
func (s *HotelService) ListOrders(ctx context.Context, req *hotel.ListOrdersReq) (resp *hotel.ListOrdersResp, err error) {
	serviceReq := service.ListOrdersReq{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}

	if req.BranchId != nil {
		id := uint64(*req.BranchId)
		serviceReq.BranchID = &id
	}
	if req.GuestSource != nil {
		serviceReq.GuestSource = req.GuestSource
	}
	if req.OrderNo != nil {
		serviceReq.OrderNo = req.OrderNo
	}
	if req.Phone != nil {
		serviceReq.Phone = req.Phone
	}
	if req.Keyword != nil {
		serviceReq.Keyword = req.Keyword
	}
	if req.OrderStatus != nil {
		serviceReq.OrderStatus = req.OrderStatus
	}
	if req.CheckInStart != nil {
		serviceReq.CheckInStart = req.CheckInStart
	}
	if req.CheckInEnd != nil {
		serviceReq.CheckInEnd = req.CheckInEnd
	}
	if req.CheckOutStart != nil {
		serviceReq.CheckOutStart = req.CheckOutStart
	}
	if req.CheckOutEnd != nil {
		serviceReq.CheckOutEnd = req.CheckOutEnd
	}
	if req.ReserveStart != nil {
		serviceReq.ReserveStart = req.ReserveStart
	}
	if req.ReserveEnd != nil {
		serviceReq.ReserveEnd = req.ReserveEnd
	}

	result, err := s.OrderService.ListOrders(serviceReq)
	if err != nil {
		return nil, err
	}

	thriftOrders := make([]*hotel.Order, len(result.List))
	for i, order := range result.List {
		thriftOrders[i] = &hotel.Order{
			Id:                int64(order.ID),
			OrderNo:           order.OrderNo,
			BranchId:          int64(order.BranchID),
			GuestId:           int64(order.GuestID),
			RoomId:            int64(order.RoomID),
			RoomTypeId:        int64(order.RoomTypeID),
			GuestSource:       order.GuestSource,
			CheckInTime:       order.CheckInTime.Format("2006-01-02 15:04:05"),
			CheckOutTime:      order.CheckOutTime.Format("2006-01-02 15:04:05"),
			ReserveTime:       order.ReserveTime.Format("2006-01-02 15:04:05"),
			OrderAmount:       order.OrderAmount,
			DepositReceived:   order.DepositReceived,
			OutstandingAmount: order.OutstandingAmount,
			OrderStatus:       order.OrderStatus,
			PayType:           order.PayType,
			PenaltyAmount:     order.PenaltyAmount,
			CreatedAt:         order.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:         order.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if order.BranchName != "" {
			thriftOrders[i].BranchName = &order.BranchName
		}
		if order.GuestName != "" {
			thriftOrders[i].GuestName = &order.GuestName
		}
		if order.RoomNo != "" {
			thriftOrders[i].RoomNo = &order.RoomNo
		}
		if order.RoomName != "" {
			thriftOrders[i].RoomName = &order.RoomName
		}
		if order.RoomTypeName != "" {
			thriftOrders[i].RoomTypeName = &order.RoomTypeName
		}
		if order.Contact != "" {
			thriftOrders[i].Contact = &order.Contact
		}
		if order.ContactPhone != "" {
			thriftOrders[i].ContactPhone = &order.ContactPhone
		}
		if order.SpecialRequest != nil {
			thriftOrders[i].SpecialRequest = order.SpecialRequest
		}
		if order.GuestCount > 0 {
			count := int8(order.GuestCount)
			thriftOrders[i].GuestCount = &count
		}
		if order.RoomCount > 0 {
			count := int8(order.RoomCount)
			thriftOrders[i].RoomCount = &count
		}
		if len(order.RoomNos) > 0 {
			thriftOrders[i].RoomNos = order.RoomNos
		}
	}

	return &hotel.ListOrdersResp{
		List:     thriftOrders,
		Total:    int64(result.Total),
		Page:     int32(result.Page),
		PageSize: int32(result.PageSize),
	}, nil
}

// GetOrder 获取订单详情
func (s *HotelService) GetOrder(ctx context.Context, orderID int64) (resp *hotel.Order, err error) {
	order, err := s.OrderService.GetOrder(uint64(orderID))
	if err != nil {
		return nil, err
	}

	thriftOrder := &hotel.Order{
		Id:                int64(order.ID),
		OrderNo:           order.OrderNo,
		BranchId:          int64(order.BranchID),
		GuestId:           int64(order.GuestID),
		RoomId:            int64(order.RoomID),
		RoomTypeId:        int64(order.RoomTypeID),
		GuestSource:       order.GuestSource,
		CheckInTime:       order.CheckInTime.Format("2006-01-02 15:04:05"),
		CheckOutTime:      order.CheckOutTime.Format("2006-01-02 15:04:05"),
		ReserveTime:       order.ReserveTime.Format("2006-01-02 15:04:05"),
		OrderAmount:       order.OrderAmount,
		DepositReceived:   order.DepositReceived,
		OutstandingAmount: order.OutstandingAmount,
		OrderStatus:       order.OrderStatus,
		PayType:           order.PayType,
		PenaltyAmount:     order.PenaltyAmount,
		CreatedAt:         order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if order.BranchName != "" {
		thriftOrder.BranchName = &order.BranchName
	}
	if order.GuestName != "" {
		thriftOrder.GuestName = &order.GuestName
	}
	if order.RoomNo != "" {
		thriftOrder.RoomNo = &order.RoomNo
	}
	if order.RoomName != "" {
		thriftOrder.RoomName = &order.RoomName
	}
	if order.RoomTypeName != "" {
		thriftOrder.RoomTypeName = &order.RoomTypeName
	}
	if order.Contact != "" {
		thriftOrder.Contact = &order.Contact
	}
	if order.ContactPhone != "" {
		thriftOrder.ContactPhone = &order.ContactPhone
	}
	if order.SpecialRequest != nil {
		thriftOrder.SpecialRequest = order.SpecialRequest
	}
	if order.GuestCount > 0 {
		count := int8(order.GuestCount)
		thriftOrder.GuestCount = &count
	}
	if order.RoomCount > 0 {
		count := int8(order.RoomCount)
		thriftOrder.RoomCount = &count
	}
	if len(order.RoomNos) > 0 {
		thriftOrder.RoomNos = order.RoomNos
	}

	return thriftOrder, nil
}

// ListInHouseGuests 获取在住客人列表
func (s *HotelService) ListInHouseGuests(ctx context.Context, req *hotel.ListInHouseGuestsReq) (resp *hotel.ListInHouseGuestsResp, err error) {
	serviceReq := service.ListInHouseGuestsReq{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}

	if req.BranchId != nil {
		id := uint64(*req.BranchId)
		serviceReq.BranchID = &id
	}
	if req.Province != nil {
		serviceReq.Province = req.Province
	}
	if req.City != nil {
		serviceReq.City = req.City
	}
	if req.District != nil {
		serviceReq.District = req.District
	}
	if req.Name != nil {
		serviceReq.Name = req.Name
	}
	if req.Phone != nil {
		serviceReq.Phone = req.Phone
	}
	if req.IdNumber != nil {
		serviceReq.IDNumber = req.IdNumber
	}
	if req.RoomNo != nil {
		serviceReq.RoomNo = req.RoomNo
	}

	result, err := s.GuestService.ListInHouseGuests(serviceReq)
	if err != nil {
		return nil, err
	}

	thriftGuests := make([]*hotel.InHouseGuest, len(result.List))
	for i, guest := range result.List {
		thriftGuests[i] = &hotel.InHouseGuest{
			Id:                int64(guest.ID),
			GuestId:           int64(guest.GuestID),
			Name:              guest.Name,
			IdType:            guest.IDType,
			IdNumber:          guest.IDNumber,
			Phone:             guest.Phone,
			CheckInTime:       guest.CheckInTime.Format("2006-01-02"),
			CheckOutTime:      guest.CheckOutTime.Format("2006-01-02"),
			OrderId:           int64(guest.OrderID),
			OrderNo:           guest.OrderNo,
			GuestSource:       guest.GuestSource,
			RoomId:            int64(guest.RoomID),
			RoomNo:            guest.RoomNo,
			RoomTypeId:        int64(guest.RoomTypeID),
			RoomTypeName:      guest.RoomTypeName,
			OrderAmount:       guest.OrderAmount,
			DepositReceived:   guest.DepositReceived,
			OutstandingAmount: guest.OutstandingAmount,
		}

		if guest.Province != nil {
			thriftGuests[i].Province = guest.Province
		}
		if guest.Address != nil {
			thriftGuests[i].Address = guest.Address
		}
		if guest.Ethnicity != nil {
			thriftGuests[i].Ethnicity = guest.Ethnicity
		}
	}

	return &hotel.ListInHouseGuestsResp{
		List:     thriftGuests,
		Total:    int64(result.Total),
		Page:     int32(result.Page),
		PageSize: int32(result.PageSize),
	}, nil
}

// ListFinancialFlows 获取收支流水列表
func (s *HotelService) ListFinancialFlows(ctx context.Context, req *hotel.ListFinancialFlowsReq) (resp *hotel.ListFinancialFlowsResp, err error) {
	serviceReq := service.ListFinancialFlowsReq{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}

	if req.BranchId != nil {
		id := uint64(*req.BranchId)
		serviceReq.BranchID = &id
	}
	if req.FlowType != nil {
		serviceReq.FlowType = req.FlowType
	}
	if req.FlowItem != nil {
		serviceReq.FlowItem = req.FlowItem
	}
	if req.PayType != nil {
		serviceReq.PayType = req.PayType
	}
	if req.OperatorId != nil {
		id := uint64(*req.OperatorId)
		serviceReq.OperatorID = &id
	}
	if req.OccurStart != nil {
		serviceReq.OccurStart = req.OccurStart
	}
	if req.OccurEnd != nil {
		serviceReq.OccurEnd = req.OccurEnd
	}

	result, err := s.FinancialService.ListFinancialFlows(serviceReq)
	if err != nil {
		return nil, err
	}

	thriftFlows := make([]*hotel.FinancialFlow, len(result.List))
	for i, flow := range result.List {
		thriftFlows[i] = &hotel.FinancialFlow{
			Id:         int64(flow.ID),
			BranchId:   int64(flow.BranchID),
			FlowType:   flow.FlowType,
			FlowItem:   flow.FlowItem,
			PayType:    flow.PayType,
			Amount:     flow.Amount,
			OccurTime:  flow.OccurTime.Format("2006-01-02 15:04:05"),
			OperatorId: int64(flow.OperatorID),
			CreatedAt:  flow.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if flow.OrderID != nil {
			orderID := int64(*flow.OrderID)
			thriftFlows[i].OrderId = &orderID
		}
		if flow.RoomID != nil {
			roomID := int64(*flow.RoomID)
			thriftFlows[i].RoomId = &roomID
		}
		if flow.GuestID != nil {
			guestID := int64(*flow.GuestID)
			thriftFlows[i].GuestId = &guestID
		}
		if flow.Remark != nil {
			thriftFlows[i].Remark = flow.Remark
		}
		if flow.RoomNo != "" {
			thriftFlows[i].RoomNo = &flow.RoomNo
		}
		if flow.GuestName != "" {
			thriftFlows[i].GuestName = &flow.GuestName
		}
		if flow.ContactPhone != "" {
			thriftFlows[i].ContactPhone = &flow.ContactPhone
		}
		if flow.OperatorName != "" {
			thriftFlows[i].OperatorName = &flow.OperatorName
		}
		if flow.OrderNo != "" {
			thriftFlows[i].OrderNo = &flow.OrderNo
		}
	}

	// 转换汇总统计
	summary := &hotel.FinancialSummaryResp{
		Income: &hotel.FinancialSummary{
			Total:           result.Summary.Income.Total,
			Cash:            result.Summary.Income.Cash,
			Alipay:          result.Summary.Income.Alipay,
			Wechat:          result.Summary.Income.WeChat,
			Unionpay:        result.Summary.Income.UnionPay,
			CardSwipe:       result.Summary.Income.CardSwipe,
			TuyouCollection: result.Summary.Income.TuyouCollection,
			CtripCollection: result.Summary.Income.CtripCollection,
			QunarCollection: result.Summary.Income.QunarCollection,
		},
		Expense: &hotel.FinancialSummary{
			Total:           result.Summary.Expense.Total,
			Cash:            result.Summary.Expense.Cash,
			Alipay:          result.Summary.Expense.Alipay,
			Wechat:          result.Summary.Expense.WeChat,
			Unionpay:        result.Summary.Expense.UnionPay,
			CardSwipe:       result.Summary.Expense.CardSwipe,
			TuyouCollection: result.Summary.Expense.TuyouCollection,
			CtripCollection: result.Summary.Expense.CtripCollection,
			QunarCollection: result.Summary.Expense.QunarCollection,
		},
		Balance: &hotel.FinancialSummary{
			Total:           result.Summary.Balance.Total,
			Cash:            result.Summary.Balance.Cash,
			Alipay:          result.Summary.Balance.Alipay,
			Wechat:          result.Summary.Balance.WeChat,
			Unionpay:        result.Summary.Balance.UnionPay,
			CardSwipe:       result.Summary.Balance.CardSwipe,
			TuyouCollection: result.Summary.Balance.TuyouCollection,
			CtripCollection: result.Summary.Balance.CtripCollection,
			QunarCollection: result.Summary.Balance.QunarCollection,
		},
	}

	return &hotel.ListFinancialFlowsResp{
		List:     thriftFlows,
		Total:    int64(result.Total),
		Page:     int32(result.Page),
		PageSize: int32(result.PageSize),
		Summary:  summary,
	}, nil
}

// ========== 账号管理 ==========

// CreateUserAccount 创建账号
func (s *HotelService) CreateUserAccount(ctx context.Context, req *hotel.CreateUserAccountReq) (resp *hotel.BaseResp, err error) {
	serviceReq := service.CreateUserAccountReq{
		Username:     req.Username,
		Password:     req.Password,
		RealName:     req.RealName,
		ContactPhone: req.ContactPhone,
		RoleID:       uint64(req.RoleId),
	}

	if req.BranchId != nil {
		id := uint64(*req.BranchId)
		serviceReq.BranchID = &id
	}
	if req.Status != nil {
		serviceReq.Status = *req.Status
	}

	if err := s.UserAccountService.CreateUserAccount(serviceReq); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "创建成功"}, nil
}

// UpdateUserAccount 更新账号
func (s *HotelService) UpdateUserAccount(ctx context.Context, req *hotel.UpdateUserAccountReq) (resp *hotel.BaseResp, err error) {
	serviceReq := service.UpdateUserAccountReq{
		ID: uint64(req.Id),
	}

	if req.Username != nil {
		serviceReq.Username = req.Username
	}
	if req.Password != nil {
		serviceReq.Password = req.Password
	}
	if req.RealName != nil {
		serviceReq.RealName = req.RealName
	}
	if req.ContactPhone != nil {
		serviceReq.ContactPhone = req.ContactPhone
	}
	if req.RoleId != nil {
		id := uint64(*req.RoleId)
		serviceReq.RoleID = &id
	}
	if req.BranchId != nil {
		id := uint64(*req.BranchId)
		serviceReq.BranchID = &id
	}
	if req.Status != nil {
		serviceReq.Status = req.Status
	}

	if err := s.UserAccountService.UpdateUserAccount(serviceReq); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "更新成功"}, nil
}

// GetUserAccount 获取账号详情
func (s *HotelService) GetUserAccount(ctx context.Context, id int64) (resp *hotel.UserAccount, err error) {
	account, err := s.UserAccountService.GetUserAccount(uint64(id))
	if err != nil {
		return nil, err
	}

	thriftAccount := &hotel.UserAccount{
		Id:           int64(account.ID),
		Username:     account.Username,
		RealName:     account.RealName,
		ContactPhone: account.ContactPhone,
		RoleId:       int64(account.RoleID),
		Status:       account.Status,
		CreatedAt:    account.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if account.RoleName != "" {
		thriftAccount.RoleName = &account.RoleName
	}
	if account.BranchID != nil {
		branchID := int64(*account.BranchID)
		thriftAccount.BranchId = &branchID
	}
	if account.BranchName != "" {
		thriftAccount.BranchName = &account.BranchName
	}
	if account.LastLoginAt != nil {
		loginAt := account.LastLoginAt.Format("2006-01-02 15:04:05")
		thriftAccount.LastLoginAt = &loginAt
	}

	return thriftAccount, nil
}

// ListUserAccounts 获取账号列表
func (s *HotelService) ListUserAccounts(ctx context.Context, req *hotel.ListUserAccountsReq) (resp *hotel.ListUserAccountsResp, err error) {
	serviceReq := service.ListUserAccountsReq{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}

	if req.RoleId != nil {
		id := uint64(*req.RoleId)
		serviceReq.RoleID = &id
	}
	if req.BranchId != nil {
		id := uint64(*req.BranchId)
		serviceReq.BranchID = &id
	}
	if req.Status != nil {
		serviceReq.Status = req.Status
	}
	if req.Keyword != nil {
		serviceReq.Keyword = req.Keyword
	}

	result, err := s.UserAccountService.ListUserAccounts(serviceReq)
	if err != nil {
		return nil, err
	}

	thriftAccounts := make([]*hotel.UserAccount, len(result.List))
	for i, account := range result.List {
		thriftAccounts[i] = &hotel.UserAccount{
			Id:           int64(account.ID),
			Username:     account.Username,
			RealName:     account.RealName,
			ContactPhone: account.ContactPhone,
			RoleId:       int64(account.RoleID),
			Status:       account.Status,
			CreatedAt:    account.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if account.RoleName != "" {
			thriftAccounts[i].RoleName = &account.RoleName
		}
		if account.BranchID != nil {
			branchID := int64(*account.BranchID)
			thriftAccounts[i].BranchId = &branchID
		}
		if account.BranchName != "" {
			thriftAccounts[i].BranchName = &account.BranchName
		}
		if account.LastLoginAt != nil {
			loginAt := account.LastLoginAt.Format("2006-01-02 15:04:05")
			thriftAccounts[i].LastLoginAt = &loginAt
		}
	}

	return &hotel.ListUserAccountsResp{
		List:     thriftAccounts,
		Total:    int64(result.Total),
		Page:     int32(result.Page),
		PageSize: int32(result.PageSize),
	}, nil
}

// DeleteUserAccount 删除账号
func (s *HotelService) DeleteUserAccount(ctx context.Context, id int64) (resp *hotel.BaseResp, err error) {
	if err := s.UserAccountService.DeleteUserAccount(uint64(id)); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "删除成功"}, nil
}

// ========== 角色管理 ==========

// CreateRole 创建角色
func (s *HotelService) CreateRole(ctx context.Context, req *hotel.CreateRoleReq) (resp *hotel.BaseResp, err error) {
	serviceReq := service.CreateRoleReq{
		RoleName: req.RoleName,
	}

	if req.Description != nil {
		serviceReq.Description = req.Description
	}
	if req.Status != nil {
		serviceReq.Status = *req.Status
	}
	if req.PermissionIds != nil {
		permissionIDs := make([]uint64, len(req.PermissionIds))
		for i, id := range req.PermissionIds {
			permissionIDs[i] = uint64(id)
		}
		serviceReq.PermissionIDs = permissionIDs
	}

	if err := s.RoleService.CreateRole(serviceReq); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "创建成功"}, nil
}

// UpdateRole 更新角色
func (s *HotelService) UpdateRole(ctx context.Context, req *hotel.UpdateRoleReq) (resp *hotel.BaseResp, err error) {
	serviceReq := service.UpdateRoleReq{
		ID: uint64(req.Id),
	}

	if req.RoleName != nil {
		serviceReq.RoleName = req.RoleName
	}
	if req.Description != nil {
		serviceReq.Description = req.Description
	}
	if req.Status != nil {
		serviceReq.Status = req.Status
	}
	if req.PermissionIds != nil {
		permissionIDs := make([]uint64, len(req.PermissionIds))
		for i, id := range req.PermissionIds {
			permissionIDs[i] = uint64(id)
		}
		serviceReq.PermissionIDs = permissionIDs
	}

	if err := s.RoleService.UpdateRole(serviceReq); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "更新成功"}, nil
}

// GetRole 获取角色详情
func (s *HotelService) GetRole(ctx context.Context, id int64) (resp *hotel.Role, err error) {
	role, err := s.RoleService.GetRole(uint64(id))
	if err != nil {
		return nil, err
	}

	thriftRole := &hotel.Role{
		Id:        int64(role.ID),
		RoleName:  role.RoleName,
		Status:    role.Status,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}

	if role.Description != nil {
		thriftRole.Description = role.Description
	}

	// 转换权限ID列表
	if len(role.Permissions) > 0 {
		permissionIDs := make([]int64, len(role.Permissions))
		for i, perm := range role.Permissions {
			permissionIDs[i] = int64(perm.ID)
		}
		thriftRole.PermissionIds = permissionIDs
	} else {
		thriftRole.PermissionIds = []int64{}
	}

	return thriftRole, nil
}

// ListRoles 获取角色列表
func (s *HotelService) ListRoles(ctx context.Context, req *hotel.ListRolesReq) (resp *hotel.ListRolesResp, err error) {
	serviceReq := service.ListRolesReq{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}

	if req.Status != nil {
		serviceReq.Status = req.Status
	}
	if req.Keyword != nil {
		serviceReq.Keyword = req.Keyword
	}

	result, err := s.RoleService.ListRoles(serviceReq)
	if err != nil {
		return nil, err
	}

	thriftRoles := make([]*hotel.Role, len(result.List))
	for i, role := range result.List {
		thriftRoles[i] = &hotel.Role{
			Id:        int64(role.ID),
			RoleName:  role.RoleName,
			Status:    role.Status,
			CreatedAt: role.CreatedAt,
			UpdatedAt: role.UpdatedAt,
		}

		if role.Description != nil {
			thriftRoles[i].Description = role.Description
		}
	}

	return &hotel.ListRolesResp{
		List:     thriftRoles,
		Total:    int64(result.Total),
		Page:     int32(result.Page),
		PageSize: int32(result.PageSize),
	}, nil
}

// DeleteRole 删除角色
func (s *HotelService) DeleteRole(ctx context.Context, id int64) (resp *hotel.BaseResp, err error) {
	if err := s.RoleService.DeleteRole(uint64(id)); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "删除成功"}, nil
}

// ========== 权限管理 ==========

// ListPermissions 获取权限列表
func (s *HotelService) ListPermissions(ctx context.Context, req *hotel.ListPermissionsReq) (resp *hotel.ListPermissionsResp, err error) {
	serviceReq := service.ListPermissionsReq{}

	if req.PermissionType != nil {
		serviceReq.PermissionType = req.PermissionType
	}
	if req.ParentId != nil {
		id := uint64(*req.ParentId)
		serviceReq.ParentID = &id
	}
	if req.Status != nil {
		serviceReq.Status = req.Status
	}

	result, err := s.PermissionService.ListPermissions(serviceReq)
	if err != nil {
		return nil, err
	}

	thriftPermissions := make([]*hotel.Permission, len(result.List))
	for i, perm := range result.List {
		thriftPermissions[i] = convertPermissionToThrift(perm)
	}

	return &hotel.ListPermissionsResp{
		List: thriftPermissions,
	}, nil
}

// convertPermissionToThrift 转换权限信息
func convertPermissionToThrift(perm service.PermissionInfo) *hotel.Permission {
	thriftPerm := &hotel.Permission{
		Id:             int64(perm.ID),
		PermissionName: perm.PermissionName,
		PermissionUrl:  perm.PermissionURL,
		PermissionType: perm.PermissionType,
		Status:         perm.Status,
	}

	if perm.ParentID != nil {
		parentID := int64(*perm.ParentID)
		thriftPerm.ParentId = &parentID
	}

	if len(perm.Children) > 0 {
		children := make([]*hotel.Permission, len(perm.Children))
		for i, child := range perm.Children {
			children[i] = convertPermissionToThrift(child)
		}
		thriftPerm.Children = children
	}

	return thriftPerm
}

// ========== 渠道配置管理 ==========

// CreateChannelConfig 创建渠道配置
func (s *HotelService) CreateChannelConfig(ctx context.Context, req *hotel.CreateChannelConfigReq) (resp *hotel.BaseResp, err error) {
	serviceReq := service.CreateChannelConfigReq{
		ChannelName: req.ChannelName,
		ChannelCode: req.ChannelCode,
		ApiURL:      req.ApiUrl,
	}

	if req.SyncRule != nil {
		serviceReq.SyncRule = *req.SyncRule
	}
	if req.Status != nil {
		serviceReq.Status = *req.Status
	}

	if err := s.ChannelConfigService.CreateChannelConfig(serviceReq); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "创建成功"}, nil
}

// UpdateChannelConfig 更新渠道配置
func (s *HotelService) UpdateChannelConfig(ctx context.Context, req *hotel.UpdateChannelConfigReq) (resp *hotel.BaseResp, err error) {
	serviceReq := service.UpdateChannelConfigReq{
		ID: uint64(req.Id),
	}

	if req.ChannelName != nil {
		serviceReq.ChannelName = req.ChannelName
	}
	if req.ChannelCode != nil {
		serviceReq.ChannelCode = req.ChannelCode
	}
	if req.ApiUrl != nil {
		serviceReq.ApiURL = req.ApiUrl
	}
	if req.SyncRule != nil {
		serviceReq.SyncRule = req.SyncRule
	}
	if req.Status != nil {
		serviceReq.Status = req.Status
	}

	if err := s.ChannelConfigService.UpdateChannelConfig(serviceReq); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "更新成功"}, nil
}

// GetChannelConfig 获取渠道配置详情
func (s *HotelService) GetChannelConfig(ctx context.Context, id int64) (resp *hotel.ChannelConfig, err error) {
	config, err := s.ChannelConfigService.GetChannelConfig(uint64(id))
	if err != nil {
		return nil, err
	}

	return &hotel.ChannelConfig{
		Id:          int64(config.ID),
		ChannelName: config.ChannelName,
		ChannelCode: config.ChannelCode,
		ApiUrl:      config.ApiURL,
		SyncRule:    config.SyncRule,
		Status:      config.Status,
		CreatedAt:   config.CreatedAt,
		UpdatedAt:   config.UpdatedAt,
	}, nil
}

// ListChannelConfigs 获取渠道配置列表
func (s *HotelService) ListChannelConfigs(ctx context.Context, req *hotel.ListChannelConfigsReq) (resp *hotel.ListChannelConfigsResp, err error) {
	serviceReq := service.ListChannelConfigsReq{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}

	if req.Status != nil {
		serviceReq.Status = req.Status
	}
	if req.Keyword != nil {
		serviceReq.Keyword = req.Keyword
	}

	result, err := s.ChannelConfigService.ListChannelConfigs(serviceReq)
	if err != nil {
		return nil, err
	}

	thriftConfigs := make([]*hotel.ChannelConfig, len(result.List))
	for i, config := range result.List {
		thriftConfigs[i] = &hotel.ChannelConfig{
			Id:          int64(config.ID),
			ChannelName: config.ChannelName,
			ChannelCode: config.ChannelCode,
			ApiUrl:      config.ApiURL,
			SyncRule:    config.SyncRule,
			Status:      config.Status,
			CreatedAt:   config.CreatedAt,
			UpdatedAt:   config.UpdatedAt,
		}
	}

	return &hotel.ListChannelConfigsResp{
		List:     thriftConfigs,
		Total:    int64(result.Total),
		Page:     int32(result.Page),
		PageSize: int32(result.PageSize),
	}, nil
}

// DeleteChannelConfig 删除渠道配置
func (s *HotelService) DeleteChannelConfig(ctx context.Context, id int64) (resp *hotel.BaseResp, err error) {
	if err := s.ChannelConfigService.DeleteChannelConfig(uint64(id)); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "删除成功"}, nil
}

// ========== 系统配置管理 ==========

// CreateSystemConfig 创建系统配置
func (s *HotelService) CreateSystemConfig(ctx context.Context, req *hotel.CreateSystemConfigReq) (resp *hotel.BaseResp, err error) {
	serviceReq := service.CreateSystemConfigReq{
		ConfigCategory: req.ConfigCategory,
		ConfigKey:      req.ConfigKey,
		ConfigValue:    req.ConfigValue,
		UpdatedBy:      uint64(req.UpdatedBy),
	}

	if req.Description != nil {
		serviceReq.Description = req.Description
	}
	if req.Status != nil {
		serviceReq.Status = *req.Status
	}

	if err := s.SystemConfigService.CreateSystemConfig(serviceReq); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "创建成功"}, nil
}

// UpdateSystemConfig 更新系统配置
func (s *HotelService) UpdateSystemConfig(ctx context.Context, req *hotel.UpdateSystemConfigReq) (resp *hotel.BaseResp, err error) {
	serviceReq := service.UpdateSystemConfigReq{
		ID:        uint64(req.Id),
		UpdatedBy: uint64(req.UpdatedBy),
	}

	if req.ConfigCategory != nil {
		serviceReq.ConfigCategory = req.ConfigCategory
	}
	if req.ConfigKey != nil {
		serviceReq.ConfigKey = req.ConfigKey
	}
	if req.ConfigValue != nil {
		serviceReq.ConfigValue = req.ConfigValue
	}
	if req.Description != nil {
		serviceReq.Description = req.Description
	}
	if req.Status != nil {
		serviceReq.Status = req.Status
	}

	if err := s.SystemConfigService.UpdateSystemConfig(serviceReq); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "更新成功"}, nil
}

// GetSystemConfig 获取系统配置详情
func (s *HotelService) GetSystemConfig(ctx context.Context, id int64) (resp *hotel.SystemConfig, err error) {
	config, err := s.SystemConfigService.GetSystemConfig(uint64(id))
	if err != nil {
		return nil, err
	}

	thriftConfig := &hotel.SystemConfig{
		Id:             int64(config.ID),
		ConfigCategory: config.ConfigCategory,
		ConfigKey:      config.ConfigKey,
		ConfigValue:    config.ConfigValue,
		Status:         config.Status,
		UpdatedAt:      config.UpdatedAt,
		UpdatedBy:      int64(config.UpdatedBy),
	}

	if config.Description != nil {
		thriftConfig.Description = config.Description
	}

	return thriftConfig, nil
}

// ListSystemConfigs 获取系统配置列表
func (s *HotelService) ListSystemConfigs(ctx context.Context, req *hotel.ListSystemConfigsReq) (resp *hotel.ListSystemConfigsResp, err error) {
	serviceReq := service.ListSystemConfigsReq{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}

	if req.ConfigCategory != nil {
		serviceReq.ConfigCategory = req.ConfigCategory
	}
	if req.Status != nil {
		serviceReq.Status = req.Status
	}
	if req.Keyword != nil {
		serviceReq.Keyword = req.Keyword
	}

	result, err := s.SystemConfigService.ListSystemConfigs(serviceReq)
	if err != nil {
		return nil, err
	}

	thriftConfigs := make([]*hotel.SystemConfig, len(result.List))
	for i, config := range result.List {
		thriftConfigs[i] = &hotel.SystemConfig{
			Id:             int64(config.ID),
			ConfigCategory: config.ConfigCategory,
			ConfigKey:      config.ConfigKey,
			ConfigValue:    config.ConfigValue,
			Status:         config.Status,
			UpdatedAt:      config.UpdatedAt,
			UpdatedBy:      int64(config.UpdatedBy),
		}

		if config.Description != nil {
			thriftConfigs[i].Description = config.Description
		}
	}

	return &hotel.ListSystemConfigsResp{
		List:     thriftConfigs,
		Total:    int64(result.Total),
		Page:     int32(result.Page),
		PageSize: int32(result.PageSize),
	}, nil
}

// DeleteSystemConfig 删除系统配置
func (s *HotelService) DeleteSystemConfig(ctx context.Context, id int64) (resp *hotel.BaseResp, err error) {
	if err := s.SystemConfigService.DeleteSystemConfig(uint64(id)); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "删除成功"}, nil
}

// GetSystemConfigsByCategory 按分类获取系统配置
func (s *HotelService) GetSystemConfigsByCategory(ctx context.Context, category string) (resp *hotel.ListSystemConfigsResp, err error) {
	configs, err := s.SystemConfigService.GetConfigByCategory(category)
	if err != nil {
		return nil, err
	}

	thriftConfigs := make([]*hotel.SystemConfig, len(configs))
	for i, config := range configs {
		thriftConfigs[i] = &hotel.SystemConfig{
			Id:             int64(config.ID),
			ConfigCategory: config.ConfigCategory,
			ConfigKey:      config.ConfigKey,
			ConfigValue:    config.ConfigValue,
			Status:         config.Status,
			UpdatedAt:      config.UpdatedAt,
			UpdatedBy:      int64(config.UpdatedBy),
		}

		if config.Description != nil {
			thriftConfigs[i].Description = config.Description
		}
	}

	return &hotel.ListSystemConfigsResp{
		List:     thriftConfigs,
		Total:    int64(len(configs)),
		Page:     1,
		PageSize: int32(len(configs)),
	}, nil
}

// ========== 黑名单管理 ==========

// CreateBlacklist 创建黑名单
func (s *HotelService) CreateBlacklist(ctx context.Context, req *hotel.CreateBlacklistReq) (resp *hotel.BaseResp, err error) {
	serviceReq := service.CreateBlacklistReq{
		IDNumber:   req.IdNumber,
		Phone:      req.Phone,
		Reason:     req.Reason,
		OperatorID: uint64(req.OperatorId),
	}

	if req.GuestId != nil {
		id := uint64(*req.GuestId)
		serviceReq.GuestID = &id
	}
	if req.Status != nil {
		serviceReq.Status = *req.Status
	}

	if err := s.BlacklistService.CreateBlacklist(serviceReq); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "创建成功"}, nil
}

// UpdateBlacklist 更新黑名单
func (s *HotelService) UpdateBlacklist(ctx context.Context, req *hotel.UpdateBlacklistReq) (resp *hotel.BaseResp, err error) {
	serviceReq := service.UpdateBlacklistReq{
		ID: uint64(req.Id),
	}

	if req.GuestId != nil {
		id := uint64(*req.GuestId)
		serviceReq.GuestID = &id
	}
	if req.IdNumber != nil {
		serviceReq.IDNumber = req.IdNumber
	}
	if req.Phone != nil {
		serviceReq.Phone = req.Phone
	}
	if req.Reason != nil {
		serviceReq.Reason = req.Reason
	}
	if req.Status != nil {
		serviceReq.Status = req.Status
	}

	if err := s.BlacklistService.UpdateBlacklist(serviceReq); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "更新成功"}, nil
}

// GetBlacklist 获取黑名单详情
func (s *HotelService) GetBlacklist(ctx context.Context, id int64) (resp *hotel.Blacklist, err error) {
	blacklist, err := s.BlacklistService.GetBlacklist(uint64(id))
	if err != nil {
		return nil, err
	}

	thriftBlacklist := &hotel.Blacklist{
		Id:         int64(blacklist.ID),
		IdNumber:   blacklist.IDNumber,
		Phone:      blacklist.Phone,
		Reason:     blacklist.Reason,
		BlackTime:  blacklist.BlackTime.Format("2006-01-02 15:04:05"),
		OperatorId: int64(blacklist.OperatorID),
		Status:     blacklist.Status,
		CreatedAt:  blacklist.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if blacklist.GuestID != nil {
		guestID := int64(*blacklist.GuestID)
		thriftBlacklist.GuestId = &guestID
	}
	if blacklist.GuestName != "" {
		thriftBlacklist.GuestName = &blacklist.GuestName
	}

	return thriftBlacklist, nil
}

// ListBlacklists 获取黑名单列表
func (s *HotelService) ListBlacklists(ctx context.Context, req *hotel.ListBlacklistsReq) (resp *hotel.ListBlacklistsResp, err error) {
	serviceReq := service.ListBlacklistsReq{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}

	if req.Status != nil {
		serviceReq.Status = req.Status
	}
	if req.Keyword != nil {
		serviceReq.Keyword = req.Keyword
	}

	result, err := s.BlacklistService.ListBlacklists(serviceReq)
	if err != nil {
		return nil, err
	}

	thriftBlacklists := make([]*hotel.Blacklist, len(result.List))
	for i, blacklist := range result.List {
		thriftBlacklists[i] = &hotel.Blacklist{
			Id:         int64(blacklist.ID),
			IdNumber:   blacklist.IDNumber,
			Phone:      blacklist.Phone,
			Reason:     blacklist.Reason,
			BlackTime:  blacklist.BlackTime.Format("2006-01-02 15:04:05"),
			OperatorId: int64(blacklist.OperatorID),
			Status:     blacklist.Status,
			CreatedAt:  blacklist.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if blacklist.GuestID != nil {
			guestID := int64(*blacklist.GuestID)
			thriftBlacklists[i].GuestId = &guestID
		}
		if blacklist.GuestName != "" {
			thriftBlacklists[i].GuestName = &blacklist.GuestName
		}
	}

	return &hotel.ListBlacklistsResp{
		List:     thriftBlacklists,
		Total:    int64(result.Total),
		Page:     int32(result.Page),
		PageSize: int32(result.PageSize),
	}, nil
}

// DeleteBlacklist 删除黑名单
func (s *HotelService) DeleteBlacklist(ctx context.Context, id int64) (resp *hotel.BaseResp, err error) {
	if err := s.BlacklistService.DeleteBlacklist(uint64(id)); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "删除成功"}, nil
}

// ========== 操作日志管理 ==========

// CreateOperationLog 创建操作日志
func (s *HotelService) CreateOperationLog(ctx context.Context, req *hotel.CreateOperationLogReq) (resp *hotel.BaseResp, err error) {
	serviceReq := service.CreateOperationLogReq{
		OperatorID:    uint64(req.OperatorId),
		Module:        req.Module,
		OperationType: req.OperationType,
		Content:       req.Content,
		OperationIP:   req.OperationIp,
		IsSuccess:     req.IsSuccess,
	}

	if req.RelatedId != nil {
		id := uint64(*req.RelatedId)
		serviceReq.RelatedID = &id
	}

	if err := s.OperationLogService.CreateOperationLog(serviceReq); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "创建成功"}, nil
}

// ListOperationLogs 查询操作日志列表
func (s *HotelService) ListOperationLogs(ctx context.Context, req *hotel.ListOperationLogsReq) (resp *hotel.ListOperationLogsResp, err error) {
	serviceReq := service.ListOperationLogsReq{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}

	if req.OperatorId != nil {
		id := uint64(*req.OperatorId)
		serviceReq.OperatorID = &id
	}
	if req.Module != nil {
		serviceReq.Module = req.Module
	}
	if req.OperationType != nil {
		serviceReq.OperationType = req.OperationType
	}
	if req.StartTime != nil {
		serviceReq.StartTime = req.StartTime
	}
	if req.EndTime != nil {
		serviceReq.EndTime = req.EndTime
	}
	if req.IsSuccess != nil {
		serviceReq.IsSuccess = req.IsSuccess
	}

	result, err := s.OperationLogService.ListOperationLogs(serviceReq)
	if err != nil {
		return nil, err
	}

	thriftLogs := make([]*hotel.OperationLog, len(result.List))
	for i, log := range result.List {
		operatorName := log.OperatorName
		thriftLogs[i] = &hotel.OperationLog{
			Id:            int64(log.ID),
			OperatorId:    int64(log.OperatorID),
			OperatorName:  &operatorName,
			Module:        log.Module,
			OperationType: log.OperationType,
			Content:       log.Content,
			OperationTime: log.OperationTime.Format("2006-01-02 15:04:05"),
			OperationIp:   log.OperationIP,
			IsSuccess:     log.IsSuccess,
			CreatedAt:     log.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if log.RelatedID != nil {
			id := int64(*log.RelatedID)
			thriftLogs[i].RelatedId = &id
		}
	}

	return &hotel.ListOperationLogsResp{
		List:     thriftLogs,
		Total:    result.Total,
		Page:     int32(result.Page),
		PageSize: int32(result.PageSize),
	}, nil
}
