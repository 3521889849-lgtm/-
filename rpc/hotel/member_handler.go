package hotel

import (
	"context"
	"example_shop/common/service"
	"example_shop/kitex_gen/hotel"
	"time"
)

// ========== 会员管理 ==========

// CreateMember 创建会员
func (s *HotelService) CreateMember(ctx context.Context, req *hotel.CreateMemberReq) (resp *hotel.BaseResp, err error) {
	serviceReq := service.CreateMemberReq{
		GuestID:     uint64(req.GuestId),
		MemberLevel: req.MemberLevel,
	}

	if req.PointsBalance != nil {
		serviceReq.PointsBalance = uint64(*req.PointsBalance)
	}
	if req.Status != nil {
		serviceReq.Status = *req.Status
	}

	if err := s.MemberService.CreateMember(serviceReq); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "创建成功"}, nil
}

// UpdateMember 更新会员
func (s *HotelService) UpdateMember(ctx context.Context, req *hotel.UpdateMemberReq) (resp *hotel.BaseResp, err error) {
	serviceReq := service.UpdateMemberReq{
		ID: uint64(req.Id),
	}

	if req.MemberLevel != nil {
		serviceReq.MemberLevel = req.MemberLevel
	}
	if req.PointsBalance != nil {
		balance := uint64(*req.PointsBalance)
		serviceReq.PointsBalance = &balance
	}
	if req.Status != nil {
		serviceReq.Status = req.Status
	}

	if err := s.MemberService.UpdateMember(serviceReq); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "更新成功"}, nil
}

// GetMember 获取会员详情
func (s *HotelService) GetMember(ctx context.Context, id int64) (resp *hotel.Member, err error) {
	member, err := s.MemberService.GetMember(uint64(id))
	if err != nil {
		return nil, err
	}

	thriftMember := &hotel.Member{
		Id:           int64(member.ID),
		GuestId:      int64(member.GuestID),
		MemberLevel:  member.MemberLevel,
		PointsBalance: int64(member.PointsBalance),
		RegisterTime: member.RegisterTime.Format("2006-01-02 15:04:05"),
		Status:       member.Status,
		CreatedAt:    member.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if member.GuestName != "" {
		thriftMember.GuestName = &member.GuestName
	}
	if member.GuestPhone != "" {
		thriftMember.GuestPhone = &member.GuestPhone
	}
	if member.LastCheckInTime != nil {
		checkInTime := member.LastCheckInTime.Format("2006-01-02 15:04:05")
		thriftMember.LastCheckInTime = &checkInTime
	}

	return thriftMember, nil
}

// GetMemberByGuestID 根据客人ID获取会员信息
func (s *HotelService) GetMemberByGuestID(ctx context.Context, guestId int64) (resp *hotel.Member, err error) {
	member, err := s.MemberService.GetMemberByGuestID(uint64(guestId))
	if err != nil {
		return nil, err
	}

	thriftMember := &hotel.Member{
		Id:           int64(member.ID),
		GuestId:      int64(member.GuestID),
		MemberLevel:  member.MemberLevel,
		PointsBalance: int64(member.PointsBalance),
		RegisterTime: member.RegisterTime.Format("2006-01-02 15:04:05"),
		Status:       member.Status,
		CreatedAt:    member.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if member.GuestName != "" {
		thriftMember.GuestName = &member.GuestName
	}
	if member.GuestPhone != "" {
		thriftMember.GuestPhone = &member.GuestPhone
	}
	if member.LastCheckInTime != nil {
		checkInTime := member.LastCheckInTime.Format("2006-01-02 15:04:05")
		thriftMember.LastCheckInTime = &checkInTime
	}

	return thriftMember, nil
}

// ListMembers 获取会员列表
func (s *HotelService) ListMembers(ctx context.Context, req *hotel.ListMembersReq) (resp *hotel.ListMembersResp, err error) {
	serviceReq := service.ListMembersReq{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}

	if req.MemberLevel != nil {
		serviceReq.MemberLevel = req.MemberLevel
	}
	if req.Status != nil {
		serviceReq.Status = req.Status
	}
	if req.Keyword != nil {
		serviceReq.Keyword = req.Keyword
	}

	result, err := s.MemberService.ListMembers(serviceReq)
	if err != nil {
		return nil, err
	}

	thriftMembers := make([]*hotel.Member, len(result.List))
	for i, member := range result.List {
		thriftMembers[i] = &hotel.Member{
			Id:           int64(member.ID),
			GuestId:      int64(member.GuestID),
			MemberLevel:  member.MemberLevel,
			PointsBalance: int64(member.PointsBalance),
			RegisterTime: member.RegisterTime.Format("2006-01-02 15:04:05"),
			Status:       member.Status,
			CreatedAt:    member.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if member.GuestName != "" {
			thriftMembers[i].GuestName = &member.GuestName
		}
		if member.GuestPhone != "" {
			thriftMembers[i].GuestPhone = &member.GuestPhone
		}
		if member.LastCheckInTime != nil {
			checkInTime := member.LastCheckInTime.Format("2006-01-02 15:04:05")
			thriftMembers[i].LastCheckInTime = &checkInTime
		}
	}

	return &hotel.ListMembersResp{
		List:     thriftMembers,
		Total:    int64(result.Total),
		Page:     int32(result.Page),
		PageSize: int32(result.PageSize),
	}, nil
}

// DeleteMember 删除会员
func (s *HotelService) DeleteMember(ctx context.Context, id int64) (resp *hotel.BaseResp, err error) {
	if err := s.MemberService.DeleteMember(uint64(id)); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "删除成功"}, nil
}

// ========== 会员权益管理 ==========

// CreateMemberRights 创建会员权益
func (s *HotelService) CreateMemberRights(ctx context.Context, req *hotel.CreateMemberRightsReq) (resp *hotel.BaseResp, err error) {
	effectiveTime, err := time.Parse("2006-01-02 15:04:05", req.EffectiveTime)
	if err != nil {
		return &hotel.BaseResp{Code: 500, Msg: "生效时间格式错误"}, nil
	}

	serviceReq := service.CreateMemberRightsReq{
		MemberLevel:   req.MemberLevel,
		RightsName:    req.RightsName,
		EffectiveTime: effectiveTime,
	}

	if req.Description != nil {
		serviceReq.Description = req.Description
	}
	if req.DiscountRatio != nil {
		ratio := float64(*req.DiscountRatio)
		serviceReq.DiscountRatio = &ratio
	}
	if req.ExpireTime != nil {
		expireTime, err := time.Parse("2006-01-02 15:04:05", *req.ExpireTime)
		if err == nil {
			serviceReq.ExpireTime = &expireTime
		}
	}
	if req.Status != nil {
		serviceReq.Status = *req.Status
	}

	if err := s.MemberRightsService.CreateMemberRights(serviceReq); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "创建成功"}, nil
}

// UpdateMemberRights 更新会员权益
func (s *HotelService) UpdateMemberRights(ctx context.Context, req *hotel.UpdateMemberRightsReq) (resp *hotel.BaseResp, err error) {
	serviceReq := service.UpdateMemberRightsReq{
		ID: uint64(req.Id),
	}

	if req.MemberLevel != nil {
		serviceReq.MemberLevel = req.MemberLevel
	}
	if req.RightsName != nil {
		serviceReq.RightsName = req.RightsName
	}
	if req.Description != nil {
		serviceReq.Description = req.Description
	}
	if req.DiscountRatio != nil {
		ratio := float64(*req.DiscountRatio)
		serviceReq.DiscountRatio = &ratio
	}
	if req.EffectiveTime != nil {
		effectiveTime, err := time.Parse("2006-01-02 15:04:05", *req.EffectiveTime)
		if err == nil {
			serviceReq.EffectiveTime = &effectiveTime
		}
	}
	if req.ExpireTime != nil {
		expireTime, err := time.Parse("2006-01-02 15:04:05", *req.ExpireTime)
		if err == nil {
			serviceReq.ExpireTime = &expireTime
		}
	}
	if req.Status != nil {
		serviceReq.Status = req.Status
	}

	if err := s.MemberRightsService.UpdateMemberRights(serviceReq); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "更新成功"}, nil
}

// GetMemberRights 获取会员权益详情
func (s *HotelService) GetMemberRights(ctx context.Context, id int64) (resp *hotel.MemberRights, err error) {
	rights, err := s.MemberRightsService.GetMemberRights(uint64(id))
	if err != nil {
		return nil, err
	}

	thriftRights := &hotel.MemberRights{
		Id:            int64(rights.ID),
		MemberLevel:   rights.MemberLevel,
		RightsName:    rights.RightsName,
		EffectiveTime: rights.EffectiveTime.Format("2006-01-02 15:04:05"),
		Status:        rights.Status,
		CreatedAt:     rights.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if rights.Description != nil {
		thriftRights.Description = rights.Description
	}
	if rights.DiscountRatio != nil {
		ratio := *rights.DiscountRatio
		thriftRights.DiscountRatio = &ratio
	}
	if rights.ExpireTime != nil {
		expireTime := rights.ExpireTime.Format("2006-01-02 15:04:05")
		thriftRights.ExpireTime = &expireTime
	}

	return thriftRights, nil
}

// ListMemberRights 获取会员权益列表
func (s *HotelService) ListMemberRights(ctx context.Context, req *hotel.ListMemberRightsReq) (resp *hotel.ListMemberRightsResp, err error) {
	serviceReq := service.ListMemberRightsReq{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}

	if req.MemberLevel != nil {
		serviceReq.MemberLevel = req.MemberLevel
	}
	if req.Status != nil {
		serviceReq.Status = req.Status
	}
	if req.Keyword != nil {
		serviceReq.Keyword = req.Keyword
	}

	result, err := s.MemberRightsService.ListMemberRights(serviceReq)
	if err != nil {
		return nil, err
	}

	thriftRightsList := make([]*hotel.MemberRights, len(result.List))
	for i, rights := range result.List {
		thriftRightsList[i] = &hotel.MemberRights{
			Id:            int64(rights.ID),
			MemberLevel:   rights.MemberLevel,
			RightsName:    rights.RightsName,
			EffectiveTime: rights.EffectiveTime.Format("2006-01-02 15:04:05"),
			Status:        rights.Status,
			CreatedAt:     rights.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if rights.Description != nil {
			thriftRightsList[i].Description = rights.Description
		}
		if rights.DiscountRatio != nil {
			ratio := *rights.DiscountRatio
			thriftRightsList[i].DiscountRatio = &ratio
		}
		if rights.ExpireTime != nil {
			expireTime := rights.ExpireTime.Format("2006-01-02 15:04:05")
			thriftRightsList[i].ExpireTime = &expireTime
		}
	}

	return &hotel.ListMemberRightsResp{
		List:     thriftRightsList,
		Total:    int64(result.Total),
		Page:     int32(result.Page),
		PageSize: int32(result.PageSize),
	}, nil
}

// GetRightsByMemberLevel 根据会员等级获取权益列表
func (s *HotelService) GetRightsByMemberLevel(ctx context.Context, memberLevel string) (resp *hotel.ListMemberRightsResp, err error) {
	rightsList, err := s.MemberRightsService.GetRightsByMemberLevel(memberLevel)
	if err != nil {
		return nil, err
	}

	thriftRightsList := make([]*hotel.MemberRights, len(rightsList))
	for i, rights := range rightsList {
		thriftRightsList[i] = &hotel.MemberRights{
			Id:            int64(rights.ID),
			MemberLevel:   rights.MemberLevel,
			RightsName:    rights.RightsName,
			EffectiveTime: rights.EffectiveTime.Format("2006-01-02 15:04:05"),
			Status:        rights.Status,
			CreatedAt:     rights.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if rights.Description != nil {
			thriftRightsList[i].Description = rights.Description
		}
		if rights.DiscountRatio != nil {
			ratio := *rights.DiscountRatio
			thriftRightsList[i].DiscountRatio = &ratio
		}
		if rights.ExpireTime != nil {
			expireTime := rights.ExpireTime.Format("2006-01-02 15:04:05")
			thriftRightsList[i].ExpireTime = &expireTime
		}
	}

	return &hotel.ListMemberRightsResp{
		List:     thriftRightsList,
		Total:    int64(len(rightsList)),
		Page:     1,
		PageSize: int32(len(rightsList)),
	}, nil
}

// DeleteMemberRights 删除会员权益
func (s *HotelService) DeleteMemberRights(ctx context.Context, id int64) (resp *hotel.BaseResp, err error) {
	if err := s.MemberRightsService.DeleteMemberRights(uint64(id)); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "删除成功"}, nil
}

// ========== 会员积分管理 ==========

// CreatePointsRecord 创建积分记录
func (s *HotelService) CreatePointsRecord(ctx context.Context, req *hotel.CreatePointsRecordReq) (resp *hotel.BaseResp, err error) {
	serviceReq := service.CreatePointsRecordReq{
		MemberID:     uint64(req.MemberId),
		ChangeType:   req.ChangeType,
		PointsValue:  req.PointsValue,
		ChangeReason: req.ChangeReason,
		OperatorID:   uint64(req.OperatorId),
	}

	if req.OrderId != nil {
		orderID := uint64(*req.OrderId)
		serviceReq.OrderID = &orderID
	}

	if err := s.MemberPointsService.CreatePointsRecord(serviceReq); err != nil {
		return &hotel.BaseResp{Code: 500, Msg: err.Error()}, nil
	}

	return &hotel.BaseResp{Code: 200, Msg: "创建成功"}, nil
}

// ListPointsRecords 获取积分记录列表
func (s *HotelService) ListPointsRecords(ctx context.Context, req *hotel.ListPointsRecordsReq) (resp *hotel.ListPointsRecordsResp, err error) {
	serviceReq := service.ListPointsRecordsReq{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}

	if req.MemberId != nil {
		memberID := uint64(*req.MemberId)
		serviceReq.MemberID = &memberID
	}
	if req.OrderId != nil {
		orderID := uint64(*req.OrderId)
		serviceReq.OrderID = &orderID
	}
	if req.ChangeType != nil {
		serviceReq.ChangeType = req.ChangeType
	}
	if req.StartTime != nil {
		serviceReq.StartTime = req.StartTime
	}
	if req.EndTime != nil {
		serviceReq.EndTime = req.EndTime
	}

	result, err := s.MemberPointsService.ListPointsRecords(serviceReq)
	if err != nil {
		return nil, err
	}

	thriftRecords := make([]*hotel.PointsRecord, len(result.List))
	for i, record := range result.List {
		thriftRecords[i] = &hotel.PointsRecord{
			Id:           int64(record.ID),
			MemberId:     int64(record.MemberID),
			ChangeType:   record.ChangeType,
			PointsValue:  record.PointsValue,
			ChangeReason: record.ChangeReason,
			ChangeTime:   record.ChangeTime.Format("2006-01-02 15:04:05"),
			OperatorId:   int64(record.OperatorID),
		}

		if record.MemberName != "" {
			thriftRecords[i].MemberName = &record.MemberName
		}
		if record.OrderID != nil {
			orderID := int64(*record.OrderID)
			thriftRecords[i].OrderId = &orderID
		}
	}

	return &hotel.ListPointsRecordsResp{
		List:     thriftRecords,
		Total:    int64(result.Total),
		Page:     int32(result.Page),
		PageSize: int32(result.PageSize),
	}, nil
}

// GetMemberPointsBalance 获取会员积分余额
func (s *HotelService) GetMemberPointsBalance(ctx context.Context, memberId int64) (resp int64, err error) {
	balance, err := s.MemberPointsService.GetMemberPointsBalance(uint64(memberId))
	if err != nil {
		return 0, err
	}
	return int64(balance), nil
}
