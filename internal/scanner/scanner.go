package scanner

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var htmlSpamTags = [][]byte{
	[]byte("<html"),
	[]byte("<body"),
	[]byte("<head"),
	[]byte("<link"),
	[]byte("<title"),
}

// Result 保存单个URL的扫描结果
type Result struct {
	URL        string
	StatusCode int
	BodySize   int64
	Error      string
	IsHTML     bool
	SubInfo    *SubscriptionInfo
}

// Scanner URL扫描器
type Scanner struct {
	client    *http.Client
	userAgent string
	threads   int
	checkHTML bool
}

// New 创建新的Scanner实例
func New(timeout time.Duration, userAgent string, threads int, checkHTML bool) *Scanner {
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
		checkHTML: checkHTML,
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

		seen := make(map[string]bool) // 用于去重
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			// 跳过重复URL
			if seen[line] {
				continue
			}
			seen[line] = true

			if isValidURL(line) {
				urls <- line
			} else {
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
	result.SubInfo = ExtractSubscriptionInfo(resp.Header)

	// 读取并计算Body大小，按需流式检测HTML垃圾标签
	bodySize, isHTML, err := drainAndCheckHTML(resp.Body, s.checkHTML)
	if err != nil {
		result.Error = fmt.Sprintf("读取响应失败: %v", err)
		return result
	}
	result.BodySize = bodySize
	result.IsHTML = isHTML

	return result
}

// drainAndCheckHTML 流式排空响应体统计字节数，并在开启检测时通过滑动窗口匹配 HTML 标签
func drainAndCheckHTML(r io.Reader, checkHTML bool) (int64, bool, error) {
	if !checkHTML {
		n, err := io.Copy(io.Discard, r)
		return n, false, err
	}

	const maxOverlap = 5 // 最长匹配词 "<title" 长度为 6，跨块边界需保留末尾 5 字节
	buf := make([]byte, maxOverlap+32*1024)
	overlap := 0
	var total int64

	for {
		n, err := r.Read(buf[overlap:])
		if n > 0 {
			total += int64(n)
			window := buf[:overlap+n]
			if containsHTMLTag(window) {
				if err == io.EOF {
					return total, true, nil
				}
				if err != nil {
					return total, true, err
				}
				rest, copyErr := io.Copy(io.Discard, r)
				return total + rest, true, copyErr
			}
			if len(window) >= maxOverlap {
				copy(buf[:maxOverlap], window[len(window)-maxOverlap:])
				overlap = maxOverlap
			} else {
				overlap = len(window)
			}
		}
		if err == io.EOF {
			return total, false, nil
		}
		if err != nil {
			return total, false, err
		}
	}
}

// containsHTMLTag 原地转换 ASCII 小写并检查是否包含任意 HTML 垃圾标签（零额外内存分配）
func containsHTMLTag(b []byte) bool {
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + ('a' - 'A')
		}
	}
	for _, tag := range htmlSpamTags {
		if bytes.Contains(b, tag) {
			return true
		}
	}
	return false
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

// CountUniqueLines 计算文件中的唯一有效行数
// 逻辑与ScanFile保持一致：跳过空行，去重
func CountUniqueLines(path string) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	seen := make(map[string]bool)
	var count int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if !seen[line] {
			seen[line] = true
			count++
		}
	}
	return count, scanner.Err()
}
