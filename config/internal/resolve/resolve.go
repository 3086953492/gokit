// Package resolve 提供配置文件路径解析功能
package resolve

import (
	"os"
	"path/filepath"
)

// ResolveOptions 路径解析选项
type ResolveOptions struct {
	// ConfigFile 显式指定配置文件路径（最高优先级）
	ConfigFile string

	// ConfigDir 配置文件目录
	ConfigDir string

	// EnvConfigKey 从环境变量读取配置文件路径的 key
	EnvConfigKey string

	// Formats 支持的配置文件格式
	Formats []string
}

// ResolvePath 按优先级解析配置文件路径
// 优先级：显式路径 > 环境变量 > 默认目录探测
func ResolvePath(opts ResolveOptions) (string, error) {
	// 1. 显式指定的配置文件路径（最高优先级）
	if opts.ConfigFile != "" {
		if _, err := os.Stat(opts.ConfigFile); err == nil {
			return opts.ConfigFile, nil
		}
		return "", os.ErrNotExist
	}

	// 2. 环境变量指定的路径
	if opts.EnvConfigKey != "" {
		if envPath := os.Getenv(opts.EnvConfigKey); envPath != "" {
			if _, err := os.Stat(envPath); err == nil {
				return envPath, nil
			}
		}
	}

	// 3. 在目录下按格式探测默认配置文件（默认 yaml -> json）
	for _, format := range opts.Formats {
		defaultPath := filepath.Join(opts.ConfigDir, "config."+format)
		if _, err := os.Stat(defaultPath); err == nil {
			return defaultPath, nil
		}
	}

	return "", os.ErrNotExist
}

