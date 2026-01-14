package utils

import (
	"runtime"
)

// CalculateConcurrency 根据系统CPU核心数和任务总量动态计算并发数
// 返回计算出的并发数和CPU核心数
func CalculateConcurrency(totalTaskCount int64) (int, int) {
	cpuCores := runtime.NumCPU()
	// 基础并发: 核心数 * 50 (IO密集型任务)
	suggested := cpuCores * 10

	// 根据任务总量调整
	// 如果只有少量任务，不需要开太多线程
	if totalTaskCount > 0 {
		// 确保并发数不超过任务数的2倍(避免过度并发)
		maxNeeded := int(totalTaskCount)
		if suggested > maxNeeded {
			suggested = maxNeeded
		}
	}

	// 限制最大/最小并发
	if suggested < 4 {
		suggested = 4
	}
	if suggested > 150 {
		suggested = 150 // 硬限制防止过载
	}

	return suggested, cpuCores
}
