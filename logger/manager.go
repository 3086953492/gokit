package logger

import (
	"io"
	"sync"

	"github.com/3086953492/gokit/logger/internal/builder"
	"go.uber.org/zap"
)

// Manager 日志管理器，提供线程安全的结构化日志功能
type Manager struct {
	sugar   *zap.SugaredLogger
	closers []io.Closer
	once    sync.Once
}

// NewManager 创建日志管理器
func NewManager(opts ...Option) (*Manager, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	// 校验配置
	if !options.Console && options.File == nil {
		return nil, ErrNoOutput
	}
	if options.File != nil && options.File.Filename == "" {
		return nil, ErrEmptyFilename
	}

	// 构建 zap logger
	zapLogger, closers, err := builder.Build(builder.Config{
		Level:          int8(options.Level),
		Console:        options.Console,
		File:           toBuilderFileConfig(options.File),
		AddCaller:      options.AddCaller,
		CallerSkip:     options.CallerSkip,
	})
	if err != nil {
		return nil, err
	}

	return &Manager{
		sugar:   zapLogger.Sugar(),
		closers: closers,
	}, nil
}

// toBuilderFileConfig 转换文件配置到内部 builder 格式
func toBuilderFileConfig(cfg *FileConfig) *builder.FileConfig {
	if cfg == nil {
		return nil
	}
	return &builder.FileConfig{
		Filename:       cfg.Filename,
		MaxSize:        cfg.MaxSize,
		MaxBackups:     cfg.MaxBackups,
		MaxAge:         cfg.MaxAge,
		Compress:       cfg.Compress,
		RotateStrategy: int(cfg.RotateStrategy),
	}
}

// Debug 输出调试级别日志
func (m *Manager) Debug(msg string, kv ...any) {
	m.sugar.Debugw(msg, kv...)
}

// Info 输出信息级别日志
func (m *Manager) Info(msg string, kv ...any) {
	m.sugar.Infow(msg, kv...)
}

// Warn 输出警告级别日志
func (m *Manager) Warn(msg string, kv ...any) {
	m.sugar.Warnw(msg, kv...)
}

// Error 输出错误级别日志
func (m *Manager) Error(msg string, kv ...any) {
	m.sugar.Errorw(msg, kv...)
}

// With 派生带固定字段的 logger（共享底层 zap）
func (m *Manager) With(kv ...any) *Manager {
	return &Manager{
		sugar:   m.sugar.With(kv...),
		closers: nil, // 派生实例不持有 closer
	}
}

// Sync 刷新缓冲区
func (m *Manager) Sync() error {
	return m.sugar.Sync()
}

// Close 关闭日志管理器，释放资源
func (m *Manager) Close() error {
	var firstErr error
	m.once.Do(func() {
		// 先 sync
		if err := m.sugar.Sync(); err != nil && firstErr == nil {
			firstErr = err
		}
		// 关闭所有 writer
		for _, c := range m.closers {
			if err := c.Close(); err != nil && firstErr == nil {
				firstErr = err
			}
		}
	})
	return firstErr
}

