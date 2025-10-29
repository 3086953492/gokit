// gokit/logger/init.go
package logger

import (
	"github.com/3086953492/gokit/config/types"
)

// InitWithConfig 使用配置初始化logger并设置为默认logger
func InitWithConfig(config types.LogConfig) error {
	logger, err := NewBuilder().WithConfig(config).Build()
	if err != nil {
		return err
	}

	SetDefault(logger)
	return nil
}
