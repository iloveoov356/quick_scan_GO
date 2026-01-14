package config

import (
	"flag"
	"fmt"
	"time"
)

const (
	DefaultThreads   = 0 // 0 表示自动根据CPU核心数动态调整
	DefaultTimeout   = 30
	DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	Version          = "0.0.6"
)

// Config 保存所有命令行配置参数
type Config struct {
	FilePath   string        // 输入URL文件路径
	Threads    int           // 并发数
	Timeout    time.Duration // 超时时间
	UserAgent  string        // User-Agent
	OutputPath string        // 输出CSV路径
}

// Parse 解析命令行参数并返回配置
func Parse() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.FilePath, "f", "", "URL文件路径 (必填)")
	flag.IntVar(&cfg.Threads, "n", DefaultThreads, "并发扫描数量 (0=自动，基于CPU核心数和任务量动态调整)")
	timeoutSec := flag.Int("t", DefaultTimeout, "超时秒数")
	flag.StringVar(&cfg.UserAgent, "ua", DefaultUserAgent, "自定义User-Agent")
	flag.StringVar(&cfg.OutputPath, "o", "", "输出CSV路径 (默认: result_<timestamp>.csv)")

	flag.Parse()

	// 验证必填参数
	if cfg.FilePath == "" {
		return nil, fmt.Errorf("必须指定URL文件路径 (-f)")
	}

	// 转换超时时间为 Duration
	cfg.Timeout = time.Duration(*timeoutSec) * time.Second

	// 设置默认输出路径
	if cfg.OutputPath == "" {
		cfg.OutputPath = fmt.Sprintf("result_%s.csv", time.Now().Format("2006-01-02__15-04-05"))
	}

	return cfg, nil
}

// PrintUsage 打印使用帮助
func PrintUsage() {
	fmt.Println("URL扫描工具 - 批量检测URL状态")
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  ./scan -f <file_path> [-n <threads>] [-t <timeout>] [-ua <user-agent>] [-o <output_path>]")
	fmt.Println()
	fmt.Println("参数说明:")
	flag.PrintDefaults()
}
