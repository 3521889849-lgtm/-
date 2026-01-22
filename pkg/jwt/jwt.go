/*
 * JWT工具包
 * 
 * 功能说明：
 * - 封装JWT Token的生成和验证功能
 * - 使用HS256算法签名
 * - 支持自定义Claims（用户信息、过期时间等）
 * 
 * JWT结构：
 * Header.Payload.Signature
 * - Header: 算法和类型
 * - Payload: 用户信息和过期时间
 * - Signature: 签名（防止篡改）
 * 
 * 使用场景：
 * - 用户登录后生成Token
 * - 后续请求验证Token获取用户信息
 */
package jwt

import (
	"errors"
	"example_shop/common/config"
	"example_shop/pkg/password"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT Claims结构体
// 包含用户信息和JWT标准字段
type Claims struct {
	UserID string `json:"user_id"` // 用户ID
	Phone  string `json:"phone"`   // 手机号（登录账号）
	jwt.RegisteredClaims           // JWT标准Claims（包含过期时间、签发时间等）
}

// secretKey 获取JWT密钥
// 从配置文件中读取，如果未配置则使用默认值
func getSecretKey() []byte {
	secret := config.Cfg.JWT.Secret
	if secret == "" {
		// 默认密钥（生产环境必须修改）
		secret = "default-secret-key-change-in-production"
	}
	return []byte(secret)
}

// GenerateToken 生成JWT Token
// 
// 参数：
//   - userID: 用户ID
//   - phone: 用户手机号
// 
// 返回值：
//   - string: JWT Token字符串
//   - error: 如果生成失败，返回错误信息
// 
// Token有效期：从配置文件读取，默认2小时
func GenerateToken(userID, phone string) (string, error) {
	// 参数校验
	if userID == "" || phone == "" {
		return "", errors.New("用户ID和手机号不能为空")
	}

	// 获取Token过期时间（从配置文件读取，单位：秒）
	expireSeconds := config.Cfg.JWT.Expire
	if expireSeconds <= 0 {
		expireSeconds = 7200 // 默认2小时（7200秒）
	}

	// 创建Claims
	claims := &Claims{
		UserID: userID,
		Phone:  phone,
		RegisteredClaims: jwt.RegisteredClaims{
			// 签发时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// 过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireSeconds) * time.Second)),
			// 签发者
			Issuer: "ticket_system",
		},
	}

	// 创建Token
	// HS256: HMAC-SHA256算法，使用密钥签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获取完整的Token字符串
	tokenString, err := token.SignedString(getSecretKey())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析并验证JWT Token
// 
// 参数：
//   - tokenString: JWT Token字符串
// 
// 返回值：
//   - *Claims: 解析出的用户信息
//   - error: 如果解析或验证失败，返回错误信息
// 
// 验证内容：
// - Token格式是否正确
// - 签名是否有效
// - Token是否过期
func ParseToken(tokenString string) (*Claims, error) {
	// 参数校验
	if tokenString == "" {
		return nil, errors.New("Token不能为空")
	}

	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名算法")
		}
		// 返回密钥用于验证签名
		return getSecretKey(), nil
	})

	// 检查解析错误
	if err != nil {
		return nil, err
	}

	// 提取Claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Token有效，返回Claims
		return claims, nil
	}

	// Token无效
	return nil, errors.New("无效的Token")
}

// RefreshToken 刷新Token
// 
// 功能说明：
// - 如果Token即将过期（剩余时间少于一半），可以刷新Token
// - 生成新的Token，有效期重新计算
// 
// 参数：
//   - tokenString: 旧的Token字符串
// 
// 返回值：
//   - string: 新的Token字符串
//   - error: 如果刷新失败，返回错误信息
func RefreshToken(tokenString string) (string, error) {
	// 解析旧Token
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 检查Token是否即将过期（剩余时间少于一半）
	expireTime := claims.ExpiresAt.Time
	now := time.Now()
	if expireTime.Sub(now) < time.Duration(config.Cfg.JWT.Expire/2)*time.Second {
		// 生成新Token
		return GenerateToken(claims.UserID, claims.Phone)
	}

	// Token还未到刷新时间，返回原Token
	return tokenString, nil
}

// ValidateToken 验证Token是否有效
// 
// 参数：
//   - tokenString: Token字符串
// 
// 返回值：
//   - bool: true表示Token有效，false表示Token无效
//   - string: 如果有效，返回用户ID；如果无效，返回错误信息
func ValidateToken(tokenString string) (bool, string) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return false, err.Error()
	}
	return true, claims.UserID
}

func HashPassword(plain string) (string, error) {
	return password.HashPassword(plain)
}

func VerifyPassword(hashedPassword, plain string) (bool, error) {
	return password.VerifyPassword(hashedPassword, plain)
}
