package ticket

import (
	"regexp"
	"strings"
	"time"
)

var safeTextRe = regexp.MustCompile(`^[\p{Han}A-Za-z0-9_\-]+$`)

func validateTextField(val, fieldName string) (string, string) {
	v := strings.TrimSpace(val)
	if v == "" {
		return "", fieldName + "不能为空"
	}
	if len(v) > 64 {
		return "", fieldName + "过长"
	}
	if !safeTextRe.MatchString(strings.ReplaceAll(v, " ", "")) {
		return "", fieldName + "包含非法字符"
	}
	return v, ""
}

func validateSeatType(seatType string) (string, string) {
	v := strings.TrimSpace(seatType)
	if v == "" {
		return "", ""
	}
	allow := map[string]struct{}{
		"硬座":  {},
		"二等座": {},
		"一等座": {},
		"商务座": {},
		"硬卧":  {},
		"软卧":  {},
	}
	if _, ok := allow[v]; !ok {
		return "", "座位类型无效"
	}
	return v, ""
}

// ValidateSeatType 校验座位类型并返回规范化后的值（白名单）。
func ValidateSeatType(seatType string) (string, string) {
	return validateSeatType(seatType)
}

// parseTravelDate 解析 YYYY-MM-DD 为 [dayStart, dayEnd)。
func parseTravelDate(dateStr string) (time.Time, time.Time, error) {
	d, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	start := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
	end := start.Add(24 * time.Hour)
	return start, end, nil
}

// ParseTravelDate 解析 YYYY-MM-DD 为 [dayStart, dayEnd)。
func ParseTravelDate(dateStr string) (time.Time, time.Time, error) {
	return parseTravelDate(dateStr)
}

// parseTimeOfDay 解析 HH:mm 为一天内的偏移量。
func parseTimeOfDay(s string) (time.Duration, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	t, err := time.Parse("15:04", s)
	if err != nil {
		return 0, false
	}
	return time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute, true
}

func applyDepartTimeWindow(dayStart, dayEnd time.Time, startStr, endStr string) (time.Time, time.Time, string) {
	start := dayStart
	end := dayEnd
	if d, ok := parseTimeOfDay(startStr); ok {
		start = dayStart.Add(d)
	}
	if d, ok := parseTimeOfDay(endStr); ok {
		end = dayStart.Add(d)
		if end.Before(start) {
			return time.Time{}, time.Time{}, "出发时段不合法"
		}
	}
	return start, end, ""
}

