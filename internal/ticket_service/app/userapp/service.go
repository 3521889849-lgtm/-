// Package userapp 承载“用户域”的应用层用例（注册/登录/实名/用户信息）。
//
// 该层的职责：
// - 校验请求参数
// - 组织领域规则与基础设施访问（DB/第三方实名等）
// - 返回对外稳定的 RPC 响应结构
package userapp

import (
	"context"
	"errors"
	"example_shop/common/config"
	"example_shop/common/db"
	"example_shop/common/encrypt"
	"example_shop/internal/ticket_service/model"
	kitexuser "example_shop/kitex_gen/user"
	"example_shop/pkg/jwt"
	"example_shop/pkg/realname"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Service struct{}

// New 创建用户域应用服务。
func New() *Service {
	return &Service{}
}

// Register 用户注册用例。
func (s *Service) Register(ctx context.Context, req *kitexuser.RegisterReq) (*kitexuser.RegisterResp, error) {
	if req.UserName == "" || req.Password == "" || req.Phone == "" {
		return &kitexuser.RegisterResp{BaseResp: &kitexuser.BaseResp{Code: 400, Msg: "参数不完整"}}, nil
	}

	var count int64
	db.MysqlDB.Model(&model.UserInfo{}).Where("phone = ?", req.Phone).Count(&count)
	if count > 0 {
		return &kitexuser.RegisterResp{BaseResp: &kitexuser.BaseResp{Code: 400, Msg: "手机号已注册"}}, nil
	}

	hashedPwd, err := jwt.HashPassword(req.Password)
	if err != nil {
		return &kitexuser.RegisterResp{BaseResp: &kitexuser.BaseResp{Code: 500, Msg: "密码加密失败: " + err.Error()}}, nil
	}

	userID := uuid.New().String()

	newUser := model.UserInfo{
		ID:               userID,
		Password:         hashedPwd,
		RealName:         req.UserName,
		Phone:            req.Phone,
		UserLevel:        "ORDINARY",
		RealNameVerified: "UNVERIFIED",
		Status:           "NORMAL",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := db.MysqlDB.Create(&newUser).Error; err != nil {
		return &kitexuser.RegisterResp{BaseResp: &kitexuser.BaseResp{Code: 500, Msg: "创建用户失败"}}, nil
	}

	return &kitexuser.RegisterResp{
		BaseResp: &kitexuser.BaseResp{Code: 200, Msg: "注册成功"},
		UserId:   userID,
	}, nil
}

// Login 用户登录用例。
func (s *Service) Login(ctx context.Context, req *kitexuser.LoginReq) (*kitexuser.LoginResp, error) {
	if req.Phone == "" || req.Password == "" {
		return &kitexuser.LoginResp{BaseResp: &kitexuser.BaseResp{Code: 400, Msg: "参数不完整"}}, nil
	}

	var u model.UserInfo
	if err := db.MysqlDB.Where("phone = ?", req.Phone).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &kitexuser.LoginResp{BaseResp: &kitexuser.BaseResp{Code: 404, Msg: "用户不存在"}}, nil
		}
		return &kitexuser.LoginResp{BaseResp: &kitexuser.BaseResp{Code: 500, Msg: "系统错误"}}, nil
	}

	ok, err := jwt.VerifyPassword(u.Password, req.Password)
	if err != nil {
		return &kitexuser.LoginResp{BaseResp: &kitexuser.BaseResp{Code: 500, Msg: "密码验证失败: " + err.Error()}}, nil
	}
	if !ok {
		return &kitexuser.LoginResp{BaseResp: &kitexuser.BaseResp{Code: 401, Msg: "密码错误"}}, nil
	}

	token, err := jwt.GenerateToken(u.ID, u.Phone)
	if err != nil {
		return &kitexuser.LoginResp{BaseResp: &kitexuser.BaseResp{Code: 500, Msg: "Token生成失败: " + err.Error()}}, nil
	}

	idCard := ""
	if u.IDCard != nil {
		idCard = *u.IDCard
	}
	return &kitexuser.LoginResp{
		BaseResp: &kitexuser.BaseResp{Code: 200, Msg: "登录成功"},
		UserId:   u.ID,
		Token:    token,
		UserInfo: &kitexuser.UserInfo{
			UserId:           u.ID,
			RealName:         u.RealName,
			IdCard:           encrypt.MaskIDCard(idCard),
			Phone:            encrypt.MaskPhone(u.Phone),
			UserLevel:        u.UserLevel,
			RealNameVerified: u.RealNameVerified,
			Status:           u.Status,
		},
	}, nil
}

// VerifyRealName 实名认证用例：调用第三方实名能力并落库更新用户状态。
func (s *Service) VerifyRealName(ctx context.Context, req *kitexuser.VerifyRealNameReq) (*kitexuser.BaseResp, error) {
	if req.UserId == "" || req.RealName == "" || req.IdCard == "" || req.Phone == "" {
		return &kitexuser.BaseResp{Code: 400, Msg: "参数不完整"}, nil
	}

	secretID := strings.TrimSpace(config.Cfg.RealName.SecretID)
	secretKey := strings.TrimSpace(config.Cfg.RealName.SecretKey)
	if secretID == "" || secretKey == "" {
		return &kitexuser.BaseResp{Code: 500, Msg: "实名认证配置缺失(v2)"}, nil
	}

	raw, err := realname.VerifyRealName(secretID, secretKey, req.RealName, req.IdCard, req.Phone)
	if err != nil {
		return &kitexuser.BaseResp{Code: 502, Msg: "实名认证调用失败(v2): " + err.Error()}, nil
	}
	ok2, msg := realname.ParseVerifyResult(raw)
	if !ok2 {
		return &kitexuser.BaseResp{Code: 400, Msg: "实名认证未通过: " + msg}, nil
	}

	result := db.MysqlDB.Model(&model.UserInfo{}).Where("user_id = ?", req.UserId).Updates(map[string]interface{}{
		"real_name":          req.RealName,
		"id_card":            req.IdCard,
		"phone":              req.Phone,
		"real_name_verified": "VERIFIED",
	})

	if result.Error != nil {
		return &kitexuser.BaseResp{Code: 500, Msg: "数据库更新失败: " + result.Error.Error()}, nil
	}
	if result.RowsAffected == 0 {
		return &kitexuser.BaseResp{Code: 404, Msg: "用户不存在"}, nil
	}

	if config.Cfg.RealName.DebugReturn {
		safe := realname.ExtractSafeInfo(raw)
		if safe != "" {
			return &kitexuser.BaseResp{Code: 200, Msg: "实名认证成功: " + safe}, nil
		}
	}
	return &kitexuser.BaseResp{Code: 200, Msg: "实名认证成功"}, nil
}

// GetUserInfo 查询用户信息用例。
func (s *Service) GetUserInfo(ctx context.Context, req *kitexuser.GetUserInfoReq) (*kitexuser.GetUserInfoResp, error) {
	if req.UserId == "" {
		return &kitexuser.GetUserInfoResp{
			BaseResp: &kitexuser.BaseResp{Code: 400, Msg: "UserId不能为空"},
		}, nil
	}

	var u model.UserInfo
	err := db.MysqlDB.Where("user_id = ?", req.UserId).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &kitexuser.GetUserInfoResp{
				BaseResp: &kitexuser.BaseResp{Code: 404, Msg: "用户不存在"},
			}, nil
		}
		return &kitexuser.GetUserInfoResp{
			BaseResp: &kitexuser.BaseResp{Code: 500, Msg: "查询失败: " + err.Error()},
		}, nil
	}

	idCard := ""
	if u.IDCard != nil {
		idCard = *u.IDCard
	}
	return &kitexuser.GetUserInfoResp{
		BaseResp: &kitexuser.BaseResp{Code: 200, Msg: "success"},
		UserInfo: &kitexuser.UserInfo{
			UserId:           u.ID,
			RealName:         u.RealName,
			IdCard:           encrypt.MaskIDCard(idCard),
			Phone:            encrypt.MaskPhone(u.Phone),
			UserLevel:        u.UserLevel,
			RealNameVerified: u.RealNameVerified,
			Status:           u.Status,
		},
	}, nil
}
