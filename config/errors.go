package config

import "errors"

// 配置包错误定义
var (
	// ErrConfigNotFound 配置文件未找到
	ErrConfigNotFound = errors.New("config file not found")

	// ErrInvalidConfig 配置无效（解析失败或验证失败）
	ErrInvalidConfig = errors.New("invalid config")

	// ErrAlreadyWatching 已经在监听配置变化
	ErrAlreadyWatching = errors.New("already watching config")
)
