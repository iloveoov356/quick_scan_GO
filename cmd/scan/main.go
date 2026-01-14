package main

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"quick_scan/internal/config"
	"quick_scan/internal/scanner"
	"quick_scan/internal/utils"
	"quick_scan/internal/writer"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	fmt.Printf("开始扫描...\n")
	fmt.Printf("  输入文件: %s\n", cfg.FilePath)
	// 计算总行数以用于进度条和动态并发数
	totalLines, err := scanner.CountUniqueLines(cfg.FilePath)
	if err != nil {
		fmt.Printf("无法计算文件行数: %v\n", err)
	}

	// 动态计算并发数
	if cfg.Threads == 0 {
		var cpuCores int
		cfg.Threads, cpuCores = utils.CalculateConcurrency(totalLines)
		fmt.Printf("  并发数: %d (自动调整: %d核CPU)\n", cfg.Threads, cpuCores)
	} else {
		fmt.Printf("  并发数: %d\n", cfg.Threads)
	}

	fmt.Printf("  超时: %v\n", cfg.Timeout)
	fmt.Printf("  输出文件: %s\n", cfg.OutputPath)
	fmt.Println()

	// 创建扫描器
	s := scanner.New(cfg.Timeout, cfg.UserAgent, cfg.Threads)

	// 开始扫描
	startTime := time.Now()
	results, err := s.ScanFile(cfg.FilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "扫描失败: %v\n", err)
		os.Exit(1)
	}

	// 创建CSV写入器
	csvWriter, err := writer.NewCSVWriter(cfg.OutputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "创建输出文件失败: %v\n", err)
		os.Exit(1)
	}
	defer csvWriter.Close()

	// 处理结果
	var count int64
	for result := range results {
		if err := csvWriter.WriteResult(result); err != nil {
			fmt.Fprintf(os.Stderr, "写入结果失败: %v\n", err)
		}
		current := atomic.AddInt64(&count, 1)
		utils.PrintProgressBar(current, totalLines)
	}
	fmt.Println() // 进度条完成后换行

	elapsed := time.Since(startTime)
	fmt.Printf("\n扫描完成!\n")
	fmt.Printf("  总计: %d 条URL\n", count)
	fmt.Printf("  耗时: %v\n", elapsed)
	fmt.Printf("  结果已保存到: %s\n", cfg.OutputPath)
}
