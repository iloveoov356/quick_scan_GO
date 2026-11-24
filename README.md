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
./scan -f <file_path> [-n <threads>] [-t <timeout>] [-ua <user-agent>] [-o <output>] [-ok] [-no-html] [-no-empty=false]
```

### 参数

| 参数        | 说明                                                                | 默认值                           |
| ----------- | ------------------------------------------------------------------- | -------------------------------- |
| `-f`        | URL 文件路径（必填）                                                | -                                |
| `-n`        | 并发数                                                              | 4                                |
| `-t`        | 超时秒数                                                            | 30                               |
| `-ua`       | User-Agent                                                          | v2rayN/6.60.0                    |
| `-o`        | 输出 CSV 路径                                                       | result\_年-月-日\_\_时-分-秒.csv |
| `-ok`       | 仅保留 HTTP 状态码为 200 的结果                                     | false                            |
| `-no-html`  | 过滤响应体包含 HTML 标签 (`<html`/`<body`/`<head`/`<link`/`<title`) | false                            |
| `-no-empty` | 过滤 Body 大小为 0 的记录（传 `-no-empty=false` 关闭）              | true                             |

### 示例

```bash
# 基础用法（默认过滤 Body 大小为 0 的记录）
./scan -f urls.txt

# 仅保留状态码为 200 且不含 HTML 垃圾页的结果
./scan -f urls.txt -ok -no-html

# 保留所有记录（包括请求失败及空 Body 记录）
./scan -f urls.txt -no-empty=false
```

## URL 文件格式

每行一个 URL：

```
https://www.google.com
http://example.com
```

## 输出格式

CSV 包含以下列（当 HTTP 响应头包含 `Subscription-Userinfo` 时自动解析订阅详情，缺失项留空）：

| 列名             | 说明                                       |
| ---------------- | ------------------------------------------ |
| 状态码           | HTTP 状态码，错误时为 -1                   |
| Body 大小(bytes) | 响应体大小                                 |
| URL              | 请求的 URL                                 |
| 已用流量         | `upload + download`，自动转换单位 (KB/GB)  |
| 总流量           | `total`，自动转换单位 (KB/MB/GB/TB)        |
| 剩余流量         | `total - 已用流量`，自动转换单位           |
| 到期时间         | `expire` 转换为 `YYYY-MM-DD` 或 `长期有效` |
| 重置日           | `reset_day` 流量重置日                     |
| 套餐名称         | `plan_name` 套餐名称（支持 URL 解码）      |
| 官网地址         | `app_url` 对应站点地址                     |
| 错误信息         | DNS 失败/超时/TLS 错误等                   |

示例输出：

```csv
状态码,Body大小(bytes),URL,已用流量,总流量,剩余流量,到期时间,重置日,套餐名称,官网地址,错误信息
200,15234,https://example.com/sub,11.00 KB,100.00 KB,89.00 KB,2100-06-06,14,VIP1,http://a.com,
-1,0,https://nonexistent.com,,,,,,,,,DNS解析失败
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
git tag v0.0.1
git push origin v0.0.1
```

GitHub Actions 自动编译多平台版本。
