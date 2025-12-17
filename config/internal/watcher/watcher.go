// Package watcher 提供配置文件监听功能
package watcher

import (
	"context"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Logger 日志接口
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

// Watcher 配置文件监听器
type Watcher struct {
	configPath string
	watcher    *fsnotify.Watcher
	logger     Logger
}

// New 创建配置文件监听器
func New(configPath string, logger Logger) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// 监听配置文件所在目录，以捕获文件的重写操作
	dir := filepath.Dir(configPath)
	if err := w.Add(dir); err != nil {
		w.Close()
		return nil, err
	}

	return &Watcher{
		configPath: configPath,
		watcher:    w,
		logger:     logger,
	}, nil
}

// Watch 开始监听配置文件变化，阻塞直到 ctx 取消
// onChange 会在配置文件发生变化时被调用
func (w *Watcher) Watch(ctx context.Context, onChange func()) {
	configName := filepath.Base(w.configPath)

	for {
		select {
		case <-ctx.Done():
			w.watcher.Close()
			return
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			// 只关注目标配置文件的写入或创建事件
			if filepath.Base(event.Name) == configName {
				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
					w.logger.Info("config file changed", "file", event.Name)
					onChange()
				}
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			w.logger.Error("watcher error", "error", err)
		}
	}
}

// Close 关闭监听器
func (w *Watcher) Close() error {
	return w.watcher.Close()
}

