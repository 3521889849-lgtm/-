package shared

import "math"

// Money2 将金额保留两位小数（四舍五入）。
func Money2(v float64) float64 {
	return math.Round(v*100) / 100
}

