package user

import (
	"context"
	"example_shop/internal/ticket_service/app/orderapp"
	"example_shop/internal/ticket_service/app/ticketapp"
	"example_shop/internal/ticket_service/app/userapp"
	kitexuser "example_shop/kitex_gen/user"
)

// UserServiceImpl 是 Kitex 生成的 user.UserService 的实现。
// 当前设计为“薄适配层”：只做 RPC -> 应用层分发，业务逻辑落在各领域 app 包中。
type UserServiceImpl struct {
	userApp   *userapp.Service
	ticketApp *ticketapp.Service
	orderApp  *orderapp.Service
}

// NewUserServiceImpl 构造一个带领域应用层依赖的 RPC 实现。
func NewUserServiceImpl() *UserServiceImpl {
	return &UserServiceImpl{
		userApp:   userapp.New(),
		ticketApp: ticketapp.New(),
		orderApp:  orderapp.New(),
	}
}

func (s *UserServiceImpl) Register(ctx context.Context, req *kitexuser.RegisterReq) (*kitexuser.RegisterResp, error) {
	return s.userApp.Register(ctx, req)
}

func (s *UserServiceImpl) Login(ctx context.Context, req *kitexuser.LoginReq) (*kitexuser.LoginResp, error) {
	return s.userApp.Login(ctx, req)
}

func (s *UserServiceImpl) VerifyRealName(ctx context.Context, req *kitexuser.VerifyRealNameReq) (*kitexuser.BaseResp, error) {
	return s.userApp.VerifyRealName(ctx, req)
}

func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *kitexuser.GetUserInfoReq) (*kitexuser.GetUserInfoResp, error) {
	return s.userApp.GetUserInfo(ctx, req)
}

func (s *UserServiceImpl) GetTrainDetail(ctx context.Context, req *kitexuser.GetTrainDetailReq) (*kitexuser.GetTrainDetailResp, error) {
	return s.ticketApp.GetTrainDetail(ctx, req)
}

func (s *UserServiceImpl) CreateOrder(ctx context.Context, req *kitexuser.CreateOrderReq) (*kitexuser.CreateOrderResp, error) {
	return s.orderApp.CreateOrder(ctx, req)
}

func (s *UserServiceImpl) PayOrder(ctx context.Context, req *kitexuser.PayOrderReq) (*kitexuser.PayOrderResp, error) {
	return s.orderApp.PayOrder(ctx, req)
}

func (s *UserServiceImpl) ConfirmPay(ctx context.Context, req *kitexuser.ConfirmPayReq) (*kitexuser.ConfirmPayResp, error) {
	return s.orderApp.ConfirmPay(ctx, req)
}

func (s *UserServiceImpl) CancelOrder(ctx context.Context, req *kitexuser.CancelOrderReq) (*kitexuser.BaseResp, error) {
	return s.orderApp.CancelOrder(ctx, req)
}

func (s *UserServiceImpl) RefundOrder(ctx context.Context, req *kitexuser.RefundOrderReq) (*kitexuser.RefundOrderResp, error) {
	return s.orderApp.RefundOrder(ctx, req)
}

func (s *UserServiceImpl) ChangeOrder(ctx context.Context, req *kitexuser.ChangeOrderReq) (*kitexuser.ChangeOrderResp, error) {
	return s.orderApp.ChangeOrder(ctx, req)
}

func (s *UserServiceImpl) GetOrder(ctx context.Context, req *kitexuser.GetOrderReq) (*kitexuser.GetOrderResp, error) {
	return s.orderApp.GetOrder(ctx, req)
}

func (s *UserServiceImpl) ListOrders(ctx context.Context, req *kitexuser.ListOrdersReq) (*kitexuser.ListOrdersResp, error) {
	return s.orderApp.ListOrders(ctx, req)
}
