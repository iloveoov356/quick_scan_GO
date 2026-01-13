package scanner

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// Result 保存单个URL的扫描结果
type Result struct {
	URL        string
	StatusCode int
	BodySize   int64
	Error      string
}

// Scanner URL扫描器
type Scanner struct {
	client    *http.Client
	userAgent string
	threads   int
}

// New 创建新的Scanner实例
func New(timeout time.Duration, userAgent string, threads int) *Scanner {
	client := &http.Client{
		Timeout: timeout,
		// 不跟随重定向,记录原始状态码
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &Scanner{
		client:    client,
		userAgent: userAgent,
		threads:   threads,
	}
}

// ScanFile 扫描文件中的所有URL
func (s *Scanner) ScanFile(filePath string) (<-chan Result, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %w", err)
	}

	urls := make(chan string, s.threads*2)
	results := make(chan Result, s.threads*2)

	// 启动worker池
	var wg sync.WaitGroup
	for i := 0; i < s.threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for u := range urls {
				results <- s.scanURL(u)
			}
		}()
	}

	// 读取URL并发送到channel
	go func() {
		defer close(urls)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if isValidURL(line) {
				urls <- line
			} else if line != "" {
				// 非法URL也记录到结果
				results <- Result{
					URL:        line,
					StatusCode: -1,
					BodySize:   0,
					Error:      "无效的URL格式",
				}
			}
		}
	}()

	// 等待所有worker完成后关闭results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	return results, nil
}

// scanURL 扫描单个URL
func (s *Scanner) scanURL(targetURL string) Result {
	result := Result{
		URL:        targetURL,
		StatusCode: -1, // 默认 -1 表示请求未成功
		BodySize:   0,
	}

	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		result.Error = fmt.Sprintf("创建请求失败: %v", err)
		return result
	}

	if s.userAgent != "" {
		req.Header.Set("User-Agent", s.userAgent)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		result.Error = classifyError(err)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	// 读取并计算Body大小
	bodySize, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		result.Error = fmt.Sprintf("读取响应失败: %v", err)
		return result
	}
	result.BodySize = bodySize

	return result
}

// classifyError 分类错误类型
func classifyError(err error) string {
	errStr := err.Error()

	switch {
	case strings.Contains(errStr, "no such host"):
		return "DNS解析失败"
	case strings.Contains(errStr, "connection refused"):
		return "连接被拒绝"
	case strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded"):
		return "请求超时"
	case strings.Contains(errStr, "certificate") || strings.Contains(errStr, "tls"):
		return "TLS/证书错误"
	case strings.Contains(errStr, "connection reset"):
		return "连接被重置"
	default:
		return fmt.Sprintf("请求失败: %v", err)
	}
}

// isValidURL 验证URL格式是否合法
func isValidURL(str string) bool {
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}
