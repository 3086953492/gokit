package logger

import (
	"go.uber.org/zap"
	"sync"
)

var (
	defaultLogger *Logger
	once          sync.Once
)

// Logger 封装的logger
type Logger struct {
	zap *zap.Logger
	mu  sync.RWMutex
}

// Default 获取默认logger实例
func Default() *Logger {
	once.Do(func() {
		zapLogger, _ := zap.NewProduction()
		defaultLogger = &Logger{zap: zapLogger}
	})
	return defaultLogger
}

// SetDefault 设置默认logger
func SetDefault(logger *zap.Logger) {
	Default().SetLogger(logger)
}

// SetLogger 设置zap logger
func (l *Logger) SetLogger(logger *zap.Logger) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if logger != nil {
		l.zap = logger
	}
}

// GetLogger 获取zap logger
func (l *Logger) GetLogger() *zap.Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.zap
}

// 实例方法
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.GetLogger().Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.GetLogger().Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.GetLogger().Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.GetLogger().Error(msg, fields...)
}

// 包级别的便捷函数
func Debug(msg string, fields ...zap.Field) {
	Default().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Default().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Default().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Default().Error(msg, fields...)
}

func GetLogger() *zap.Logger {
	return Default().GetLogger()
}
