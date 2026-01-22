package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ============ JWT 配置变量 ============

// JWT配置，包括密钥和过期时间
var (
	jwtSecret     = []byte("your-secret-key-change-in-production") // JWT密钥，生产环境应从配置读取
	jwtExpireTime = 24 * time.Hour                                 // Token过期时间，默认24小时
)

// ============ JWT 声明结构体 ============

// Claims JWT声明
// 包含用户基本信息和标准JWT声明
type Claims struct {
	UserID   int64  `json:"user_id"`   // 用户ID（客服记录的主键ID）
	UserName string `json:"user_name"` // 用户名（登录账号）
	RoleCode string `json:"role_code"` // 角色编码（admin/customer_service）
	jwt.RegisteredClaims              // 内嵌标准JWT声明（包含过期时间、签发时间等）
}

// ============ JWT 配置函数 ============

// SetJWTSecret 设置JWT密钥
// 用于在服务启动时从配置文件加载密钥
// 参数:
//   - secret: 密钥字符串，空字符串不更新
func SetJWTSecret(secret string) {
	if secret != "" {
		jwtSecret = []byte(secret)
	}
}

// SetJWTExpireTime 设置Token过期时间
func SetJWTExpireTime(hours int) {
	if hours > 0 {
		jwtExpireTime = time.Duration(hours) * time.Hour
	}
}

// GenerateToken 生成JWT Token
// 参数: 用户ID、用户名、角色编码
// 返回: Token字符串和错误信息
func GenerateToken(userID int64, userName, roleCode string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		UserName: userName,
		RoleCode: roleCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtExpireTime)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(now),                    // 签发时间
			NotBefore: jwt.NewNumericDate(now),                    // 生效时间
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析JWT Token
// 参数: Token字符串
// 返回: Claims声明和错误信息
func ParseToken(tokenString string) (*Claims, error) {
	// 去除Bearer前缀
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetUserIDFromToken 从Token中获取用户ID
func GetUserIDFromToken(tokenString string) (int64, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// GetRoleCodeFromToken 从Token中获取角色编码
func GetRoleCodeFromToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.RoleCode, nil
}
