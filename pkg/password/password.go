/*
 * 密码工具包
 *
 * 功能说明：
 * - 封装密码加密和验证功能
 * - 使用bcrypt算法进行密码加密
 * - 提供统一的密码处理接口
 *
 * 使用场景：
 * - 用户注册时加密密码
 * - 用户登录时验证密码
 */
package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 加密密码
//
// 参数：
//   - password: 明文密码
//
// 返回值：
//   - string: 加密后的密码（可存储到数据库）
//   - error: 如果加密失败，返回错误信息
//
// 算法说明：
// - 使用bcrypt算法，默认成本因子为bcrypt.DefaultCost(10)
// - 每次加密结果不同，但验证时能正确匹配
// - 安全性高，但性能略慢（单次加密约100ms）
func HashPassword(password string) (string, error) {
	// 参数校验
	if password == "" {
		return "", errors.New("密码不能为空")
	}

	// 使用bcrypt加密密码
	// bcrypt.DefaultCost: 成本因子为10，平衡安全性和性能
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// 返回加密后的密码字符串
	return string(hashedBytes), nil
}

// VerifyPassword 验证密码
//
// 参数：
//   - hashedPassword: 数据库中存储的加密密码
//   - password: 用户输入的明文密码
//
// 返回值：
//   - bool: true表示密码正确，false表示密码错误
//   - error: 如果验证过程出错，返回错误信息
//
// 使用场景：
// - 用户登录时验证密码
func VerifyPassword(hashedPassword, password string) (bool, error) {
	// 参数校验
	if hashedPassword == "" || password == "" {
		return false, errors.New("密码不能为空")
	}

	// 使用bcrypt验证密码
	// CompareHashAndPassword: 将明文密码与加密密码对比
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// 密码不匹配
		return false, nil
	}

	// 密码匹配
	return true, nil
}

// MustHashPassword 加密密码（失败时panic）
//
// 注意：仅在确定密码有效的情况下使用，避免程序崩溃
func MustHashPassword(password string) string {
	hashed, err := HashPassword(password)
	if err != nil {
		panic("密码加密失败: " + err.Error())
	}
	return hashed
}
