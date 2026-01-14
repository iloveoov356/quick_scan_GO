package utils

import (
	"fmt"
	"strings"
)

// PrintProgressBar 打印进度条
func PrintProgressBar(current, total int64) {
	const barWidth = 40
	var percent float64
	if total > 0 {
		percent = float64(current) / float64(total) * 100
	}

	// 限制百分比最大为 100
	if percent > 100 {
		percent = 100
	}

	filled := int(float64(barWidth) * percent / 100)
	if filled > barWidth {
		filled = barWidth
	}

	bar := strings.Repeat("=", filled) + strings.Repeat("-", barWidth-filled)

	// \r 回车不换行，实现原地刷新
	fmt.Printf("\r进度: [%s] %.1f%% (%d/%d)", bar, percent, current, total)
}
