package middleware

import (
	"context"
	"strings"

	"example_shop/internal/gateway/http/dto"
	"example_shop/pkg/jwt"

	"github.com/cloudwego/hertz/pkg/app"
)

const ContextKeyUserID = "user_id"

func GetUserID(c *app.RequestContext) string {
	v, ok := c.Get(ContextKeyUserID)
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

func AuthRequired() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		raw := strings.TrimSpace(string(c.GetHeader("Authorization")))
		if raw == "" {
			raw = strings.TrimSpace(string(c.GetHeader("X-Token")))
		}
		if raw == "" {
			c.JSON(401, dto.BaseHTTPResp{Code: 401, Msg: "缺少Token"})
			c.Abort()
			return
		}

		tokenString := raw
		if strings.HasPrefix(strings.ToLower(raw), "bearer ") {
			tokenString = strings.TrimSpace(raw[7:])
		}

		claims, err := jwt.ParseToken(tokenString)
		if err != nil || claims == nil || strings.TrimSpace(claims.UserID) == "" {
			c.JSON(401, dto.BaseHTTPResp{Code: 401, Msg: "Token无效"})
			c.Abort()
			return
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Request.Header.Set("X-User-ID", claims.UserID)
		c.Next(ctx)
	}
}

