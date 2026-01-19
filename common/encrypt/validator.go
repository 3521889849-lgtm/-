package encrypt

import (
	"regexp"
	"strings"
)

var (
	idNumber18Regex = regexp.MustCompile(`^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`)
	idNumber15Regex = regexp.MustCompile(`^[1-9]\d{5}\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}$`)
	phoneRegex      = regexp.MustCompile(`^1[3-9]\d{9}$`)
)

// ValidateIDNumber 验证身份证号格式（支持18位和15位）
func ValidateIDNumber(idNumber string) bool {
	idNumber = strings.TrimSpace(idNumber)
	if idNumber == "" {
		return false
	}

	switch len(idNumber) {
	case 18:
		return idNumber18Regex.MatchString(idNumber)
	case 15:
		return idNumber15Regex.MatchString(idNumber)
	default:
		return false
	}
}

// ValidatePhone 验证手机号格式（11位，1开头）
func ValidatePhone(phone string) bool {
	phone = strings.TrimSpace(phone)
	return phone != "" && phoneRegex.MatchString(phone)
}
