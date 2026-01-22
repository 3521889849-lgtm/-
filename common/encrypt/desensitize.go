package encrypt

import "unicode/utf8"

func MaskPhone(phone string) string {
	if len(phone) < 7 {
		return phone
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

func MaskIDCard(idCard string) string {
	if len(idCard) < 8 {
		return idCard
	}
	return idCard[:3] + "***********" + idCard[len(idCard)-4:]
}

func MaskName(name string) string {
	if name == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(name)
	if r == utf8.RuneError && size == 0 {
		return ""
	}
	return string(r) + "**"
}
