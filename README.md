# URL 快速扫描工具

一个使用 Go 语言开发的命令行 URL 扫描工具，可批量检测 URL 的 HTTP 状态码和响应 Body 大小。

## 项目结构

```
quick_scan_GO/
├── .github/
│   └── workflows/
│       └── release.yml      # GitHub Actions 自动编译发布
├── cmd/
│   └── scan/
│       └── main.go          # 程序入口
├── internal/
│   ├── config/
│   │   └── config.go        # 配置参数处理
│   ├── scanner/
│   │   └── scanner.go       # 核心扫描逻辑
│   └── writer/
│       └── csv_writer.go    # CSV 输出处理
├── .gitignore
├── go.mod
└── README.md
```

## 安装

### 方式一：下载预编译版本

从 [Releases](../../releases) 页面下载对应平台的版本：

- Windows 64-bit: `scan-vX.X.X-windows-amd64.zip`
- Linux 64-bit: `scan-vX.X.X-linux-amd64.tar.gz`
- macOS Intel: `scan-vX.X.X-darwin-amd64.tar.gz`
- macOS Apple Silicon (M 芯片): `scan-vX.X.X-darwin-arm64.tar.gz`

### 方式二：从源码编译

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
| `-o`  | 输出 CSV 路径           | result\_年-月-日\_\_时-分-秒.csv                                |

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

| 列名             | 说明                                 |
| ---------------- | ------------------------------------ |
| 状态码           | HTTP 状态码，错误时为 -1             |
| Body 大小(bytes) | 响应体大小，错误时为 0               |
| URL              | 请求的 URL（原样保留）               |
| 错误信息         | DNS 失败/超时/TLS 错误等详细错误信息 |

### 输出示例

```csv
状态码,Body大小(bytes),URL,错误信息
200,15234,https://www.google.com,
301,0,http://github.com,
-1,0,invalid-url,无效的URL格式
-1,0,https://nonexistent.domain.com,DNS解析失败
-1,0,https://slow-server.com,请求超时
```

## 特性

- ✅ 基于 worker 池的并发扫描
- ✅ 自动 URL 去重
- ✅ 支持自定义超时时间
- ✅ 自定义 User-Agent
- ✅ 自动过滤无效 URL 并记录
- ✅ 完善的错误分类（DNS/超时/TLS 等）
- ✅ UTF-8 BOM 支持（Excel 兼容）
- ✅ 实时进度显示
- ✅ 不跟随重定向（记录原始状态码）
- ✅ 100% Go 标准库，零外部依赖

## 发布新版本

```bash
git tag v1.0.0
git push origin v1.0.0
```

GitHub Actions 会自动编译并发布 Windows/Linux/macOS 多平台版本。

## 扩展

项目结构清晰，便于未来扩展更多参数，如：

- 自定义请求头
- 代理支持
- 重试机制
- 更多输出格式（JSON 等）
- 响应体大小限制
