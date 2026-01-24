# URL 快速扫描工具

一个使用 Go 开发的高性能命令行工具，批量检测 URL 的 HTTP 状态码和响应体大小。

## 安装

### 下载预编译版本

从 [Releases](../../releases) 页面下载对应平台的版本

### 源码编译

```bash
go build -o scan ./cmd/scan
```

要求：Go 1.21+

## 使用方法

```bash
./scan -f <file_path> [-n <threads>] [-t <timeout>] [-ua <user-agent>] [-o <output>]
```

### 参数

| 参数  | 说明                 | 默认值                           |
| ----- | -------------------- | -------------------------------- |
| `-f`  | URL 文件路径（必填） | -                                |
| `-n`  | 并发数               | 4                                |
| `-t`  | 超时秒数             | 30                               |
| `-ua` | User-Agent           | Mozilla/5.0 (Windows NT 10.0...) |
| `-o`  | 输出 CSV 路径        | result\_年-月-日\_\_时-分-秒.csv |

### 示例

```bash
# 基础用法
./scan -f urls.txt

# 自定义参数
./scan -f urls.txt -n 10 -t 15 -o output.csv
```

## URL 文件格式

每行一个 URL：

```
https://www.google.com
http://example.com
```

## 输出格式

CSV 包含 4 列：

| 列名             | 说明                     |
| ---------------- | ------------------------ |
| 状态码           | HTTP 状态码，错误时为 -1 |
| Body 大小(bytes) | 响应体大小               |
| URL              | 请求的 URL               |
| 错误信息         | DNS 失败/超时/TLS 错误等 |

示例输出：

```csv
状态码,Body大小(bytes),URL,错误信息
200,15234,https://www.google.com,
-1,0,https://nonexistent.com,DNS解析失败
```

## 特性

- ✅ Worker 池并发扫描
- ✅ 自动 URL 去重
- ✅ 自定义超时和 User-Agent
- ✅ 完善的错误分类（DNS/超时/TLS）
- ✅ UTF-8 BOM（Excel 兼容）
- ✅ 实时进度显示
- ✅ 零外部依赖

## 发布新版本

```bash
git tag v1.0.0
git push origin v1.0.0
```

GitHub Actions 自动编译多平台版本。
