# URL 快速扫描工具

一个使用 Go 语言开发的命令行 URL 扫描工具，可批量检测 URL 的 HTTP 状态码和响应 Body 大小。

## 项目结构

```
quick_scan_GO/
├── cmd/
│   └── scan/
│       └── main.go          # 程序入口
├── internal/
│   ├── config/
│   │   └── config.go        # 配置参数处理
│   ├── scanner/
│   │   └── scanner.go       # 核心扫描逻辑
│   └── writer/
│       └── csv_writer.go    # CSV输出处理
├── go.mod
└── README.md
```

## 安装

确保已安装 Go 1.21+，然后运行：

```bash
go build -o scan ./cmd/scan
```

## 使用方法

```bash
./scan -f <file_path> [-n <threads>] [-t <timeout>] [-ua <user-agent>] [-o <output_path>]
```

### 参数说明

| 参数  | 说明                    | 默认值                                                          |
| ----- | ----------------------- | --------------------------------------------------------------- |
| `-f`  | URL 文件路径 **(必填)** | -                                                               |
| `-n`  | 并发扫描数量            | 4                                                               |
| `-t`  | 超时秒数                | 30                                                              |
| `-ua` | 自定义 User-Agent       | Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36... |
| `-o`  | 输出 CSV 路径           | result\_<时间戳>.csv                                            |

### 使用示例

```bash
# 基础用法
./scan -f urls.txt

# 自定义并发和超时
./scan -f urls.txt -n 10 -t 15

# 完整参数
./scan -f urls.txt -n 10 -t 15 -ua "MyBot/1.0" -o output.csv
```

## URL 文件格式

每行一个 URL，支持 HTTP 和 HTTPS：

```
https://www.google.com
https://www.github.com
http://example.com
invalid-url-will-be-logged
```

## 输出格式

CSV 文件包含以下列：

| 列名             | 说明                           |
| ---------------- | ------------------------------ |
| 状态码           | HTTP 状态码 (如 200, 404, 301) |
| Body 大小(bytes) | 响应体大小                     |
| URL              | 请求的 URL                     |
| 错误信息         | 如有错误则显示                 |

## 特性

- ✅ 基于 worker 池的并发扫描
- ✅ 支持自定义超时时间
- ✅ 自定义 User-Agent
- ✅ 自动过滤无效 URL 并记录
- ✅ UTF-8 BOM 支持（Excel 兼容）
- ✅ 实时进度显示
- ✅ 不跟随重定向（记录原始状态码）

## 扩展

项目结构清晰，便于未来扩展更多参数，如：

- 自定义请求头
- 代理支持
- 重试机制
- 更多输出格式（JSON 等）
