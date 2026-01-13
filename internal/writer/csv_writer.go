package writer

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"quick_scan/internal/scanner"
)

// CSVWriter 用于写入扫描结果到CSV文件
type CSVWriter struct {
	file   *os.File
	writer *csv.Writer
}

// NewCSVWriter 创建新的CSV写入器
func NewCSVWriter(filePath string) (*CSVWriter, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("创建CSV文件失败: %w", err)
	}

	// 写入UTF-8 BOM，确保Excel正确识别中文
	file.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(file)

	// 写入表头
	if err := writer.Write([]string{"状态码", "Body大小(bytes)", "URL", "错误信息"}); err != nil {
		file.Close()
		return nil, fmt.Errorf("写入表头失败: %w", err)
	}

	return &CSVWriter{
		file:   file,
		writer: writer,
	}, nil
}

// WriteResult 写入单条扫描结果
func (w *CSVWriter) WriteResult(result scanner.Result) error {
	statusCode := ""
	if result.StatusCode > 0 {
		statusCode = strconv.Itoa(result.StatusCode)
	}

	bodySize := ""
	if result.BodySize > 0 {
		bodySize = strconv.FormatInt(result.BodySize, 10)
	}

	record := []string{
		statusCode,
		bodySize,
		result.URL,
		result.Error,
	}

	return w.writer.Write(record)
}

// Close 关闭CSV写入器
func (w *CSVWriter) Close() error {
	w.writer.Flush()
	if err := w.writer.Error(); err != nil {
		return fmt.Errorf("刷新CSV缓冲区失败: %w", err)
	}
	return w.file.Close()
}
