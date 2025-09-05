// yabase/logger/init.go
package logger

import (
	"github.com/3086953492/YaBase/configs"
)

// InitWithConfig 使用配置初始化logger并设置为默认logger
func InitWithConfig(config configs.LogConfig) error {
	logger, err := NewBuilder().WithConfig(config).Build()
	if err != nil {
		return err
	}

	SetDefault(logger)
	return nil
}
