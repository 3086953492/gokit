// yabase/logger/builder.go
package logger

import (
	"fmt"
	"os"

	"github.com/3086953492/YaBase/config/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Builder logger构建器
type Builder struct {
    config types.LogConfig
}

// NewBuilder 创建新的构建器
func NewBuilder() *Builder {
    return &Builder{
        config: types.DefaultConfig(),
    }
}

// WithConfig 设置完整配置
func (b *Builder) WithConfig(config types.LogConfig) *Builder {
    b.config = config
    return b
}

// WithLevel 设置日志级别
func (b *Builder) WithLevel(level string) *Builder {
    b.config.Level = level
    return b
}

// WithFilename 设置日志文件名
func (b *Builder) WithFilename(filename string) *Builder {
    b.config.Filename = filename
    return b
}

// WithRotateDaily 设置是否按日期轮转
func (b *Builder) WithRotateDaily(daily bool) *Builder {
    b.config.RotateDaily = daily
    return b
}

// WithConsole 设置是否输出到控制台
func (b *Builder) WithConsole(console bool) *Builder {
    b.config.Console = console
    return b
}

// WithRotationConfig 设置轮转配置
func (b *Builder) WithRotationConfig(maxSize, maxBackups, maxAge int, compress bool) *Builder {
    b.config.MaxSize = maxSize
    b.config.MaxBackups = maxBackups
    b.config.MaxAge = maxAge
    b.config.Compress = compress
    return b
}

// Build 构建logger
func (b *Builder) Build() (*zap.Logger, error) {
    if err := b.config.Validate(); err != nil {
        return nil, err
    }

    // 确保日志目录存在
    if err := os.MkdirAll(b.config.LogsDir, os.ModePerm); err != nil {
        return nil, fmt.Errorf("create logs directory failed: %v", err)
    }

    // 解析日志级别
    level := b.parseLevel(b.config.Level)

    // 创建编码器配置
    encoderConfig := b.createEncoderConfig()

    // 创建核心
    cores := make([]zapcore.Core, 0, 2)

    // 文件输出核心
    if b.config.Filename != "" {
        fileCore, err := b.createFileCore(encoderConfig, level)
        if err != nil {
            return nil, err
        }
        cores = append(cores, fileCore)
    }

    // 控制台输出核心
    if b.config.Console {
        consoleCore := b.createConsoleCore(encoderConfig, level)
        cores = append(cores, consoleCore)
    }

    if len(cores) == 0 {
        return nil, fmt.Errorf("no output configured")
    }

    // 创建logger
    core := zapcore.NewTee(cores...)
    logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

    return logger, nil
}

// parseLevel 解析日志级别
func (b *Builder) parseLevel(levelStr string) zapcore.Level {
    switch levelStr {
    case "debug":
        return zapcore.DebugLevel
    case "info":
        return zapcore.InfoLevel
    case "warn":
        return zapcore.WarnLevel
    case "error":
        return zapcore.ErrorLevel
    default:
        return zapcore.InfoLevel
    }
}

// createEncoderConfig 创建编码器配置
func (b *Builder) createEncoderConfig() zapcore.EncoderConfig {
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

// createFileCore 创建文件输出核心
func (b *Builder) createFileCore(encoderConfig zapcore.EncoderConfig, level zapcore.Level) (zapcore.Core, error) {
    var fileWriter zapcore.WriteSyncer

    if b.config.RotateDaily {
        // 使用按日期分割的写入器
        dailyWriter := NewDailyRotateWriter(
            b.config.Filename,
            b.config.MaxSize,
            b.config.MaxBackups,
            b.config.MaxAge,
            b.config.Compress,
        )
        fileWriter = zapcore.AddSync(dailyWriter)
    } else {
        // 使用lumberjack按大小分割
        lumberJackLogger := &lumberjack.Logger{
            Filename:   b.config.Filename,
            MaxSize:    b.config.MaxSize,
            MaxBackups: b.config.MaxBackups,
            MaxAge:     b.config.MaxAge,
            Compress:   b.config.Compress,
        }
        fileWriter = zapcore.AddSync(lumberJackLogger)
    }

    return zapcore.NewCore(
        zapcore.NewJSONEncoder(encoderConfig),
        fileWriter,
        level,
    ), nil
}

// createConsoleCore 创建控制台输出核心
func (b *Builder) createConsoleCore(encoderConfig zapcore.EncoderConfig, level zapcore.Level) zapcore.Core {
    return zapcore.NewCore(
        zapcore.NewConsoleEncoder(encoderConfig),
        zapcore.AddSync(os.Stdout),
        level,
    )
}