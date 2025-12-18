// Package writer 提供日志写入器实现
package writer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/3086953492/gokit/logger/internal/cleanup"
	"gopkg.in/natefinch/lumberjack.v2"
)

// DailyRotateWriter 按日期分割的日志写入器
// 在日期切换或首次写入时自动切换文件，无后台 goroutine
type DailyRotateWriter struct {
	mu          sync.Mutex
	filename    string
	maxSize     int
	maxBackups  int
	maxAge      int
	compress    bool
	currentDate string
	writer      io.WriteCloser
}

// NewDailyRotateWriter 创建按日期分割的日志写入器
func NewDailyRotateWriter(filename string, maxSize, maxBackups, maxAge int, compress bool) *DailyRotateWriter {
	return &DailyRotateWriter{
		filename:   filename,
		maxSize:    maxSize,
		maxBackups: maxBackups,
		maxAge:     maxAge,
		compress:   compress,
	}
}

// Write 实现 io.Writer 接口
func (d *DailyRotateWriter) Write(p []byte) (n int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	currentDate := time.Now().Format("2006-01-02")

	// 日期变化或未初始化时切换文件
	if d.currentDate != currentDate || d.writer == nil {
		if err := d.rotate(currentDate); err != nil {
			return 0, err
		}
	}

	return d.writer.Write(p)
}

// rotate 切换到新日期的文件
func (d *DailyRotateWriter) rotate(date string) error {
	// 关闭旧 writer
	if d.writer != nil {
		d.writer.Close()
	}

	// 生成带日期的文件名：basename-YYYY-MM-DD.ext
	dir := filepath.Dir(d.filename)
	ext := filepath.Ext(d.filename)
	base := d.filename[:len(d.filename)-len(ext)]
	dailyFilename := fmt.Sprintf("%s-%s%s", base, date, ext)

	// 确保目录存在
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 创建新的 lumberjack writer（处理单日内按大小轮转）
	d.writer = &lumberjack.Logger{
		Filename:   dailyFilename,
		MaxSize:    d.maxSize,
		MaxBackups: d.maxBackups,
		MaxAge:     d.maxAge,
		Compress:   d.compress,
	}
	d.currentDate = date

	// 日期切换时执行一次清理（异步，不阻塞写入）
	go cleanup.CleanupOldLogs(d.filename, d.maxAge, d.maxBackups)

	return nil
}

// Close 关闭写入器
func (d *DailyRotateWriter) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.writer != nil {
		err := d.writer.Close()
		d.writer = nil
		return err
	}
	return nil
}

// Sync 同步数据（如果底层支持）
func (d *DailyRotateWriter) Sync() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if syncer, ok := d.writer.(interface{ Sync() error }); ok {
		return syncer.Sync()
	}
	return nil
}

