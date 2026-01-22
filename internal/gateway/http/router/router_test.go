package router_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"example_shop/internal/gateway/http/handler"
	"example_shop/internal/gateway/http/router"
	"example_shop/kitex_gen/common"
	"example_shop/kitex_gen/orderapi/orderservice"
	"example_shop/kitex_gen/ticketapi/ticketservice"
	kitexuser "example_shop/kitex_gen/userapi"
	"example_shop/kitex_gen/userapi/userservice"
	"example_shop/pkg/jwt"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/kitex/client/callopt"
)

type mockUserClient struct {
	registerFn       func(ctx context.Context, req *kitexuser.RegisterReq) (*kitexuser.RegisterResp, error)
	loginFn          func(ctx context.Context, req *kitexuser.LoginReq) (*kitexuser.LoginResp, error)
	verifyRealNameFn func(ctx context.Context, req *kitexuser.VerifyRealNameReq) (*common.BaseResp, error)
	getUserInfoFn    func(ctx context.Context, req *kitexuser.GetUserInfoReq) (*kitexuser.GetUserInfoResp, error)
}

func (m *mockUserClient) Register(ctx context.Context, req *kitexuser.RegisterReq, _ ...callopt.Option) (*kitexuser.RegisterResp, error) {
	return m.registerFn(ctx, req)
}

func (m *mockUserClient) Login(ctx context.Context, req *kitexuser.LoginReq, _ ...callopt.Option) (*kitexuser.LoginResp, error) {
	return m.loginFn(ctx, req)
}

func (m *mockUserClient) VerifyRealName(ctx context.Context, req *kitexuser.VerifyRealNameReq, _ ...callopt.Option) (*common.BaseResp, error) {
	return m.verifyRealNameFn(ctx, req)
}

func (m *mockUserClient) GetUserInfo(ctx context.Context, req *kitexuser.GetUserInfoReq, _ ...callopt.Option) (*kitexuser.GetUserInfoResp, error) {
	return m.getUserInfoFn(ctx, req)
}

func newTestServer(t *testing.T, c userservice.Client) *server.Hertz {
	t.Helper()
	var ticketClient ticketservice.Client
	var orderClient orderservice.Client
	h := server.Default()
	router.RegisterRoutes(h, handler.NewApp(c, ticketClient, orderClient))
	return h
}

func TestRegister_Success(t *testing.T) {
	mock := &mockUserClient{
		registerFn: func(ctx context.Context, req *kitexuser.RegisterReq) (*kitexuser.RegisterResp, error) {
			if req.UserName != "u1" || req.Password != "p1" || req.Phone != "13800138000" {
				t.Fatalf("unexpected req: %#v", req)
			}
			return &kitexuser.RegisterResp{BaseResp: &common.BaseResp{Code: 200, Msg: "ok"}, UserId: "uid-1"}, nil
		},
		loginFn:          func(ctx context.Context, req *kitexuser.LoginReq) (*kitexuser.LoginResp, error) { return nil, nil },
		verifyRealNameFn: func(ctx context.Context, req *kitexuser.VerifyRealNameReq) (*common.BaseResp, error) { return nil, nil },
		getUserInfoFn:    func(ctx context.Context, req *kitexuser.GetUserInfoReq) (*kitexuser.GetUserInfoResp, error) { return nil, nil },
	}
	h := newTestServer(t, mock)

	body := []byte(`{"user_name":"u1","password":"p1","phone":"13800138000"}`)
	w := ut.PerformRequest(h.Engine, "POST", "/api/v1/user/register", &ut.Body{Body: bytes.NewBuffer(body), Len: len(body)}, ut.Header{Key: "Content-Type", Value: "application/json"})
	resp := w.Result()
	if resp.StatusCode() != 200 {
		t.Fatalf("unexpected status: %d, body=%s", resp.StatusCode(), string(resp.Body()))
	}

	var got map[string]any
	_ = json.Unmarshal(resp.Body(), &got)
	if got["user_id"] != "uid-1" {
		t.Fatalf("unexpected body: %s", string(resp.Body()))
	}
}

func TestRegister_BadRequest(t *testing.T) {
	mock := &mockUserClient{
		registerFn:       func(ctx context.Context, req *kitexuser.RegisterReq) (*kitexuser.RegisterResp, error) { return nil, nil },
		loginFn:          func(ctx context.Context, req *kitexuser.LoginReq) (*kitexuser.LoginResp, error) { return nil, nil },
		verifyRealNameFn: func(ctx context.Context, req *kitexuser.VerifyRealNameReq) (*common.BaseResp, error) { return nil, nil },
		getUserInfoFn:    func(ctx context.Context, req *kitexuser.GetUserInfoReq) (*kitexuser.GetUserInfoResp, error) { return nil, nil },
	}
	h := newTestServer(t, mock)

	body := []byte(`{"user_name":"u1","password":"p1"}`)
	w := ut.PerformRequest(h.Engine, "POST", "/api/v1/user/register", &ut.Body{Body: bytes.NewBuffer(body), Len: len(body)}, ut.Header{Key: "Content-Type", Value: "application/json"})
	resp := w.Result()
	if resp.StatusCode() != 400 {
		t.Fatalf("unexpected status: %d, body=%s", resp.StatusCode(), string(resp.Body()))
	}
}

func TestLogin_Success(t *testing.T) {
	mock := &mockUserClient{
		registerFn: func(ctx context.Context, req *kitexuser.RegisterReq) (*kitexuser.RegisterResp, error) {
			return nil, nil
		},
		loginFn: func(ctx context.Context, req *kitexuser.LoginReq) (*kitexuser.LoginResp, error) {
			if req.Phone != "13800138000" || req.Password != "p1" {
				t.Fatalf("unexpected req: %#v", req)
			}
			return &kitexuser.LoginResp{
				BaseResp: &common.BaseResp{Code: 200, Msg: "ok"},
				UserId:   "uid-1",
				Token:    "t-1",
				UserInfo: &kitexuser.UserInfo{UserId: "uid-1", Phone: "138****8000"},
			}, nil
		},
		verifyRealNameFn: func(ctx context.Context, req *kitexuser.VerifyRealNameReq) (*common.BaseResp, error) { return nil, nil },
		getUserInfoFn:    func(ctx context.Context, req *kitexuser.GetUserInfoReq) (*kitexuser.GetUserInfoResp, error) { return nil, nil },
	}
	h := newTestServer(t, mock)

	body := []byte(`{"phone":"13800138000","password":"p1"}`)
	w := ut.PerformRequest(h.Engine, "POST", "/api/v1/user/login", &ut.Body{Body: bytes.NewBuffer(body), Len: len(body)}, ut.Header{Key: "Content-Type", Value: "application/json"})
	resp := w.Result()
	if resp.StatusCode() != 200 {
		t.Fatalf("unexpected status: %d, body=%s", resp.StatusCode(), string(resp.Body()))
	}
	var got map[string]any
	_ = json.Unmarshal(resp.Body(), &got)
	if got["token"] != "t-1" {
		t.Fatalf("unexpected body: %s", string(resp.Body()))
	}
}

func TestVerifyRealName_Success(t *testing.T) {
	mock := &mockUserClient{
		registerFn: func(ctx context.Context, req *kitexuser.RegisterReq) (*kitexuser.RegisterResp, error) { return nil, nil },
		loginFn:    func(ctx context.Context, req *kitexuser.LoginReq) (*kitexuser.LoginResp, error) { return nil, nil },
		verifyRealNameFn: func(ctx context.Context, req *kitexuser.VerifyRealNameReq) (*common.BaseResp, error) {
			if req.UserId != "uid-1" || req.RealName != "张三" || req.IdCard != "110101199001011234" || req.Phone != "13800138000" {
				t.Fatalf("unexpected req: %#v", req)
			}
			return &common.BaseResp{Code: 200, Msg: "ok"}, nil
		},
		getUserInfoFn: func(ctx context.Context, req *kitexuser.GetUserInfoReq) (*kitexuser.GetUserInfoResp, error) { return nil, nil },
	}
	h := newTestServer(t, mock)

	token, err := jwt.GenerateToken("uid-1", "13800138000")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	body := []byte(`{"real_name":"张三","id_card":"110101199001011234","phone":"13800138000"}`)
	w := ut.PerformRequest(h.Engine, "POST", "/api/v1/user/verify_realname", &ut.Body{Body: bytes.NewBuffer(body), Len: len(body)},
		ut.Header{Key: "Content-Type", Value: "application/json"},
		ut.Header{Key: "Authorization", Value: "Bearer " + token},
	)
	resp := w.Result()
	if resp.StatusCode() != 200 {
		t.Fatalf("unexpected status: %d, body=%s", resp.StatusCode(), string(resp.Body()))
	}
}

func TestGetUserInfo_Success(t *testing.T) {
	mock := &mockUserClient{
		registerFn:       func(ctx context.Context, req *kitexuser.RegisterReq) (*kitexuser.RegisterResp, error) { return nil, nil },
		loginFn:          func(ctx context.Context, req *kitexuser.LoginReq) (*kitexuser.LoginResp, error) { return nil, nil },
		verifyRealNameFn: func(ctx context.Context, req *kitexuser.VerifyRealNameReq) (*common.BaseResp, error) { return nil, nil },
		getUserInfoFn: func(ctx context.Context, req *kitexuser.GetUserInfoReq) (*kitexuser.GetUserInfoResp, error) {
			if req.UserId != "uid-1" {
				t.Fatalf("unexpected req: %#v", req)
			}
			return &kitexuser.GetUserInfoResp{BaseResp: &common.BaseResp{Code: 200, Msg: "ok"}, UserInfo: &kitexuser.UserInfo{UserId: "uid-1"}}, nil
		},
	}
	h := newTestServer(t, mock)

	token, err := jwt.GenerateToken("uid-1", "13800138000")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	w := ut.PerformRequest(h.Engine, "GET", "/api/v1/user/info", nil, ut.Header{Key: "Authorization", Value: "Bearer " + token})
	resp := w.Result()
	if resp.StatusCode() != 200 {
		t.Fatalf("unexpected status: %d, body=%s", resp.StatusCode(), string(resp.Body()))
	}
}
