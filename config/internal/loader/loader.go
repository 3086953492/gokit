// Package loader 提供配置文件加载功能
package loader

import (
	"fmt"

	"github.com/spf13/viper"
)

// Load 从指定路径加载配置文件到目标结构体
// 使用独立的 viper 实例，避免全局状态污染
func Load(configPath string, target any) error {
	v := viper.New()
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	if err := v.Unmarshal(target); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	return nil
}

