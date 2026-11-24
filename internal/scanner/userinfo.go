package scanner

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// SubscriptionInfo 解析并格式化后的订阅信息
type SubscriptionInfo struct {
	UsedTraffic   string // 已用流量 (upload + download)
	TotalTraffic  string // 总流量
	RemainTraffic string // 剩余流量
	ExpireDate    string // 到期时间 (YYYY-MM-DD)
	ResetDay      string // 流量重置日
	PlanName      string // 套餐名称
	AppURL        string // 官网地址
}

// ExtractSubscriptionInfo 从 HTTP 响应头中提取并解析 Subscription-Userinfo
func ExtractSubscriptionInfo(header http.Header) *SubscriptionInfo {
	if header == nil {
		return nil
	}
	raw := strings.TrimSpace(header.Get("Subscription-Userinfo"))
	if raw == "" {
		return nil
	}
	return ParseSubscriptionInfo(raw)
}

// ParseSubscriptionInfo 解析键值对字符串为人类可读的 SubscriptionInfo
func ParseSubscriptionInfo(raw string) *SubscriptionInfo {
	var (
		upload, download, total int64
		hasUp, hasDown, hasTot  bool
		info                    SubscriptionInfo
		matchedAny              bool
	)

	parts := strings.Split(raw, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(kv[0]))
		val := strings.TrimSpace(kv[1])
		if val == "" {
			continue
		}

		switch key {
		case "upload":
			if v, err := strconv.ParseInt(val, 10, 64); err == nil && v >= 0 {
				upload = v
				hasUp = true
				matchedAny = true
			}
		case "download":
			if v, err := strconv.ParseInt(val, 10, 64); err == nil && v >= 0 {
				download = v
				hasDown = true
				matchedAny = true
			}
		case "total":
			if v, err := strconv.ParseInt(val, 10, 64); err == nil && v >= 0 {
				total = v
				hasTot = true
				matchedAny = true
			}
		case "expire":
			if ts, err := strconv.ParseInt(val, 10, 64); err == nil {
				if ts <= 0 {
					info.ExpireDate = "长期有效"
				} else {
					// 兼容毫秒级时间戳
					if ts > 1e12 {
						ts /= 1000
					}
					info.ExpireDate = time.Unix(ts, 0).Format("2006-01-02")
				}
				matchedAny = true
			}
		case "reset_day":
			info.ResetDay = val
			matchedAny = true
		case "plan_name":
			if decoded, err := url.QueryUnescape(val); err == nil {
				info.PlanName = decoded
			} else {
				info.PlanName = val
			}
			matchedAny = true
		case "app_url":
			if decoded, err := url.QueryUnescape(val); err == nil {
				info.AppURL = decoded
			} else {
				info.AppURL = val
			}
			matchedAny = true
		}
	}

	if !matchedAny {
		return nil
	}

	if hasUp || hasDown {
		used := upload + download
		info.UsedTraffic = FormatBytes(used)
		if hasTot {
			remain := total - used
			if remain < 0 {
				remain = 0
			}
			info.RemainTraffic = FormatBytes(remain)
		}
	}
	if hasTot {
		info.TotalTraffic = FormatBytes(total)
	}

	return &info
}

// FormatBytes 将字节数格式化为人类可读的流量单位 (B/KB/MB/GB/TB)
func FormatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit && exp < 3; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.2f %s", float64(b)/float64(div), units[exp])
}
