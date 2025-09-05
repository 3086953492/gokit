package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// DailyRotateWriter 按日期分割的日志写入器
type DailyRotateWriter struct {
	mu          sync.Mutex
	filename    string
	maxSize     int
	maxBackups  int
	maxAge      int
	compress    bool
	currentDate string
	writer      io.WriteCloser
	cleanupDone chan struct{}
	stopCleanup chan struct{}
}

// NewDailyRotateWriter 创建按日期分割的日志写入器
func NewDailyRotateWriter(filename string, maxSize, maxBackups, maxAge int, compress bool) *DailyRotateWriter {
	d := &DailyRotateWriter{
		filename:    filename,
		maxSize:     maxSize,
		maxBackups:  maxBackups,
		maxAge:      maxAge,
		compress:    compress,
		cleanupDone: make(chan struct{}),
		stopCleanup: make(chan struct{}),
	}

	// 启动定时清理goroutine
	go d.startCleanupRoutine()

	return d
}

// startCleanupRoutine 启动定时清理routine
func (d *DailyRotateWriter) startCleanupRoutine() {
	// 每天凌晨1点执行清理
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			// 只在凌晨1点执行清理
			if now.Hour() == 1 {
				CleanupOldLogs(d.filename, d.maxAge, d.maxBackups)
			}
		case <-d.stopCleanup:
			close(d.cleanupDone)
			return
		}
	}
}

// Write 实现 io.Writer 接口
func (d *DailyRotateWriter) Write(p []byte) (n int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 获取当前日期
	currentDate := time.Now().Format("2006-01-02")

	// 如果日期变了或者还没有初始化writer，就创建新的writer
	if d.currentDate != currentDate || d.writer == nil {
		if d.writer != nil {
			d.writer.Close()
		}

		// 生成带日期的文件名
		dir := filepath.Dir(d.filename)
		ext := filepath.Ext(d.filename)
		base := d.filename[:len(d.filename)-len(ext)]
		dailyFilename := fmt.Sprintf("%s-%s%s", base, currentDate, ext)

		// 确保目录存在
		if err := os.MkdirAll(dir, 0755); err != nil {
			return 0, err
		}

		// 创建新的lumberjack writer
		d.writer = &lumberjack.Logger{
			Filename:   dailyFilename,
			MaxSize:    d.maxSize,
			MaxBackups: d.maxBackups,
			MaxAge:     d.maxAge,
			Compress:   d.compress,
		}

		d.currentDate = currentDate

		// 执行一次清理（日期变更时）
		go CleanupOldLogs(d.filename, d.maxAge, d.maxBackups)
	}

	return d.writer.Write(p)
}

// Close 关闭writer
func (d *DailyRotateWriter) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 停止清理routine
	close(d.stopCleanup)
	<-d.cleanupDone

	if d.writer != nil {
		return d.writer.Close()
	}
	return nil
}

// Sync 同步数据
func (d *DailyRotateWriter) Sync() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if syncer, ok := d.writer.(interface{ Sync() error }); ok {
		return syncer.Sync()
	}
	return nil
}
