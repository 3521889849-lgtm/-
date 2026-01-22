package userapi

import (
	"context"
	"example_shop/internal/ticket_service/app/userapp"
	"example_shop/kitex_gen/common"
	legacy "example_shop/kitex_gen/user"
	"example_shop/kitex_gen/userapi"
)

type UserServiceImpl struct {
	userApp *userapp.Service
}

func NewUserServiceImpl() *UserServiceImpl {
	return &UserServiceImpl{userApp: userapp.New()}
}

func (s *UserServiceImpl) Register(ctx context.Context, req *userapi.RegisterReq) (*userapi.RegisterResp, error) {
	resp, err := s.userApp.Register(ctx, &legacy.RegisterReq{UserName: req.UserName, Password: req.Password, Phone: req.Phone})
	if err != nil {
		return nil, err
	}
	return &userapi.RegisterResp{BaseResp: toCommon(resp.BaseResp), UserId: resp.UserId}, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, req *userapi.LoginReq) (*userapi.LoginResp, error) {
	resp, err := s.userApp.Login(ctx, &legacy.LoginReq{Phone: req.Phone, Password: req.Password})
	if err != nil {
		return nil, err
	}
	return &userapi.LoginResp{BaseResp: toCommon(resp.BaseResp), UserId: resp.UserId, Token: resp.Token, UserInfo: toUserInfo(resp.UserInfo)}, nil
}

func (s *UserServiceImpl) VerifyRealName(ctx context.Context, req *userapi.VerifyRealNameReq) (*common.BaseResp, error) {
	resp, err := s.userApp.VerifyRealName(ctx, &legacy.VerifyRealNameReq{UserId: req.UserId, RealName: req.RealName, IdCard: req.IdCard, Phone: req.Phone})
	if err != nil {
		return nil, err
	}
	return toCommon(resp), nil
}

func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *userapi.GetUserInfoReq) (*userapi.GetUserInfoResp, error) {
	resp, err := s.userApp.GetUserInfo(ctx, &legacy.GetUserInfoReq{UserId: req.UserId})
	if err != nil {
		return nil, err
	}
	return &userapi.GetUserInfoResp{BaseResp: toCommon(resp.BaseResp), UserInfo: toUserInfo(resp.UserInfo)}, nil
}

func toCommon(b *legacy.BaseResp) *common.BaseResp {
	if b == nil {
		return nil
	}
	return &common.BaseResp{Code: b.Code, Msg: b.Msg}
}

func toUserInfo(u *legacy.UserInfo) *userapi.UserInfo {
	if u == nil {
		return nil
	}
	return &userapi.UserInfo{
		UserId:           u.UserId,
		RealName:         u.RealName,
		IdCard:           u.IdCard,
		Phone:            u.Phone,
		UserLevel:        u.UserLevel,
		RealNameVerified: u.RealNameVerified,
		Status:           u.Status,
	}
}
