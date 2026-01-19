package encrypt

import (
	"strings"
)

// desensitize 通用脱敏函数（保留前后指定字符数，中间用*代替）
func desensitize(s string, prefixLen, suffixLen int) string {
	if s == "" {
		return ""
	}

	runes := []rune(s)
	total := len(runes)

	if total <= prefixLen+suffixLen {
		return strings.Repeat("*", total)
	}

	prefix := string(runes[:prefixLen])
	suffix := string(runes[total-suffixLen:])
	middle := strings.Repeat("*", total-prefixLen-suffixLen)

	return prefix + middle + suffix
}

// DesensitizeIDNumber 脱敏身份证号（显示前3位和后4位）
func DesensitizeIDNumber(idNumber string) string {
	return desensitize(idNumber, 3, 4)
}

// DesensitizePhone 脱敏手机号（显示前3位和后4位）
func DesensitizePhone(phone string) string {
	return desensitize(phone, 3, 4)
}

// DesensitizeName 脱敏姓名（显示第一个字符）
func DesensitizeName(name string) string {
	if name == "" {
		return ""
	}

	runes := []rune(name)
	if len(runes) == 1 {
		return "*"
	}

	return string(runes[0]) + strings.Repeat("*", len(runes)-1)
}

// DesensitizeAddress 脱敏地址（显示前6个字符）
func DesensitizeAddress(address string) string {
	if address == "" {
		return ""
	}

	runes := []rune(address)
	if len(runes) <= 6 {
		return strings.Repeat("*", len(runes))
	}

	return string(runes[:6]) + strings.Repeat("*", len(runes)-6)
}
