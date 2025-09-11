package config

import "sync"

// 全局配置管理器
var (
	globalConfig  *Config
	globalMutex   sync.RWMutex
	isInitialized bool
)

// GetGlobalConfig 获取全局配置
func GetGlobalConfig() *Config {
	globalMutex.RLock()
	defer globalMutex.RUnlock()
	return globalConfig
}

// IsConfigInitialized 检查配置是否已初始化
func IsConfigInitialized() bool {
	globalMutex.RLock()
	defer globalMutex.RUnlock()
	return isInitialized
}

// InitConfig 初始化全局配置
// 只需要在main函数调用一次，后续配置变更会自动更新全局配置
func InitConfig() error {
	cfg := &Config{}
	err := LoadConfig(cfg, func(newCfg *Config) {
		// 配置变更时自动更新全局配置
		globalMutex.Lock()
		globalConfig = newCfg
		globalMutex.Unlock()
	})

	if err == nil {
		// 初始化时设置全局配置
		globalMutex.Lock()
		globalConfig = cfg
		isInitialized = true
		globalMutex.Unlock()
	}

	return err
}
