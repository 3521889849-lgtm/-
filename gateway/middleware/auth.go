// Package middleware 提供 Gateway 网关的 HTTP 中间件
// 该包实现以下功能：
// 1. JWT 身份认证（登录校验、Token 生成与解析）
// 2. 角色权限控制（管理员/客服权限验证）
// 3. 请求链路追踪（TraceID 生成与传递）
// 4. 请求日志记录
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// ============ 上下文键定义 ============

// 上下文键定义，用于在 context 中存储用户信息
type contextKey string

const (
	ContextKeyUserID   contextKey = "user_id"   // 用户ID上下文键
	ContextKeyUserName contextKey = "user_name" // 用户名上下文键
	ContextKeyRoleCode contextKey = "role_code" // 角色编码上下文键
)

// 角色编码常量（与model.RoleAdmin等保持一致）
const (
	RoleAdmin           = "admin"            // 管理员
	RoleCustomerService = "customer_service" // 客服
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// respondJSON 返回JSON响应
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// AuthMiddleware 登录校验中间件
// 验证请求头中的Token，将用户信息存入上下文
// 返回包装后的Handler
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 获取Authorization头
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondJSON(w, http.StatusOK, Response{Code: 401, Msg: "请先登录"})
			return
		}

		// 解析Token
		claims, err := ParseToken(authHeader)
		if err != nil {
			respondJSON(w, http.StatusOK, Response{Code: 401, Msg: "登录已过期，请重新登录"})
			return
		}

		// 将用户信息存入上下文
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyUserName, claims.UserName)
		ctx = context.WithValue(ctx, ContextKeyRoleCode, claims.RoleCode)

		// 继续处理请求
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthMiddlewareFunc 登录校验中间件（函数版本）
// 用于包装 http.HandlerFunc
func AuthMiddlewareFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取Authorization头
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondJSON(w, http.StatusOK, Response{Code: 401, Msg: "请先登录"})
			return
		}

		// 解析Token
		claims, err := ParseToken(authHeader)
		if err != nil {
			respondJSON(w, http.StatusOK, Response{Code: 401, Msg: "登录已过期，请重新登录"})
			return
		}

		// 将用户信息存入上下文
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyUserName, claims.UserName)
		ctx = context.WithValue(ctx, ContextKeyRoleCode, claims.RoleCode)

		// 继续处理请求
		next(w, r.WithContext(ctx))
	}
}

// CheckAdminPermission 管理员权限校验中间件
// 仅允许管理员访问，客服访问返回403
func CheckAdminPermission(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从上下文获取角色编码
		roleCode, ok := r.Context().Value(ContextKeyRoleCode).(string)
		if !ok || roleCode == "" {
			respondJSON(w, http.StatusOK, Response{Code: 401, Msg: "请先登录"})
			return
		}

		// 校验是否为管理员
		if roleCode != RoleAdmin {
			respondJSON(w, http.StatusOK, Response{Code: 403, Msg: "权限不足，仅管理员可操作"})
			return
		}

		// 继续处理请求
		next(w, r)
	}
}

// GetUserIDFromContext 从上下文获取用户ID
func GetUserIDFromContext(ctx context.Context) int64 {
	if userID, ok := ctx.Value(ContextKeyUserID).(int64); ok {
		return userID
	}
	return 0
}

// GetUserNameFromContext 从上下文获取用户名
func GetUserNameFromContext(ctx context.Context) string {
	if userName, ok := ctx.Value(ContextKeyUserName).(string); ok {
		return userName
	}
	return ""
}

// GetRoleCodeFromContext 从上下文获取角色编码
func GetRoleCodeFromContext(ctx context.Context) string {
	if roleCode, ok := ctx.Value(ContextKeyRoleCode).(string); ok {
		return roleCode
	}
	return ""
}

// IsAdmin 判断当前用户是否为管理员
func IsAdmin(ctx context.Context) bool {
	return GetRoleCodeFromContext(ctx) == RoleAdmin
}

// IsCustomerService 判断当前用户是否为客服
func IsCustomerService(ctx context.Context) bool {
	roleCode := GetRoleCodeFromContext(ctx)
	return roleCode == RoleCustomerService || roleCode == RoleAdmin
}

// extractBearerToken 从Authorization头中提取Token
func extractBearerToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	return authHeader
}
