// Package builder 提供 zap logger 构建功能
package builder

import (
	"io"
	"os"
	"path/filepath"

	"github.com/3086953492/gokit/logger/internal/writer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// RotateStrategy 轮转策略常量
const (
	RotateBySize = 0
	RotateByDate = 1
)

// FileConfig 文件配置
type FileConfig struct {
	Filename       string
	MaxSize        int
	MaxBackups     int
	MaxAge         int
	Compress       bool
	RotateStrategy int
}

// Config 构建配置
type Config struct {
	Level      int8
	Console    bool
	File       *FileConfig
	AddCaller  bool
	CallerSkip int
}

// Build 构建 zap.Logger 并返回需要关闭的 io.Closer 列表
func Build(cfg Config) (*zap.Logger, []io.Closer, error) {
	level := zapcore.Level(cfg.Level)
	encoderConfig := createEncoderConfig()

	var cores []zapcore.Core
	var closers []io.Closer

	// 控制台输出
	if cfg.Console {
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 文件输出
	if cfg.File != nil {
		fileWriter, err := createFileWriter(cfg.File)
		if err != nil {
			return nil, nil, err
		}
		closers = append(closers, fileWriter)

		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(fileWriter),
			level,
		)
		cores = append(cores, fileCore)
	}

	// 合并核心
	core := zapcore.NewTee(cores...)

	// 构建选项
	var zapOpts []zap.Option
	if cfg.AddCaller {
		zapOpts = append(zapOpts, zap.AddCaller())
		if cfg.CallerSkip > 0 {
			zapOpts = append(zapOpts, zap.AddCallerSkip(cfg.CallerSkip))
		}
	}
	zapOpts = append(zapOpts, zap.AddStacktrace(zapcore.ErrorLevel))

	logger := zap.New(core, zapOpts...)
	return logger, closers, nil
}

// createEncoderConfig 创建编码器配置
func createEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// createFileWriter 创建文件写入器
func createFileWriter(cfg *FileConfig) (io.WriteCloser, error) {
	// 确保目录存在
	dir := filepath.Dir(cfg.Filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	if cfg.RotateStrategy == RotateByDate {
		return writer.NewDailyRotateWriter(
			cfg.Filename,
			cfg.MaxSize,
			cfg.MaxBackups,
			cfg.MaxAge,
			cfg.Compress,
		), nil
	}

	// 默认按大小轮转
	return &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}, nil
}

