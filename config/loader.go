package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	ConfigEnvVar = "CONFIG"
)

var (
	// 支持的配置文件格式
	supportedFormats = []string{"yaml", "json"}

	// Gin模式对应的配置文件名
	modeConfigMap = map[string]string{
		gin.DebugMode:   "config.yaml",
		gin.TestMode:    "config.test.yaml",
		gin.ReleaseMode: "config.release.yaml",
	}
)

// LoadConfig 通用配置加载函数
// cfg: 配置结构指针，会被自动填充
// configDir: 配置文件目录，如 "./configs"
// onReload: 可选的重载回调函数
func LoadConfig(cfg *Config, configDir string, onReload ...func(*Config)) error {
	configPath := determineConfigPath(configDir)

	if err := loadConfigFile(configPath, configDir); err != nil {
		return fmt.Errorf("配置文件读取失败: %v", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("配置解析失败: %v", err)
	}

	setupConfigWatch(cfg, onReload...)
	fmt.Printf("成功读取配置文件: %s\n", viper.ConfigFileUsed())
	return nil
}

// determineConfigPath 确定配置文件路径
func determineConfigPath(configDir string) string {
	// 处理命令行参数
	var config string
	flag.StringVar(&config, "c", "", "指定配置文件路径")
	flag.Parse()

	// 1. 命令行参数优先级最高
	if config != "" {
		fmt.Printf("使用命令行参数指定的配置文件: %s\n", config)
		return config
	}

	// 2. 环境变量次之
	if configEnv := os.Getenv(ConfigEnvVar); configEnv != "" {
		fmt.Printf("使用环境变量指定的配置文件: %s\n", configEnv)
		return configEnv
	}

	// 3. 根据Gin模式自动选择
	if gin.Mode() != "" {
		if configFile, exists := modeConfigMap[gin.Mode()]; exists {
			configPath := filepath.Join(configDir, configFile)
			fmt.Printf("使用Gin模式(%s)对应的配置文件: %s\n", gin.Mode(), configPath)
			return configPath
		}
	}

	// 4. 返回空字符串，使用默认配置目录
	return ""
}

// loadConfigFile 加载配置文件
func loadConfigFile(configPath, configDir string) error {
	if configPath != "" {
		// 使用指定的配置文件
		viper.SetConfigFile(configPath)
		return viper.ReadInConfig()
	}

	// 使用默认配置目录，尝试多种格式
	return loadDefaultConfig(configDir)
}

// loadDefaultConfig 加载默认配置
func loadDefaultConfig(configDir string) error {
	viper.AddConfigPath(configDir)

	for _, format := range supportedFormats {
		viper.SetConfigName("config")
		viper.SetConfigType(format)

		if err := viper.ReadInConfig(); err == nil {
			fmt.Printf("成功读取%s格式的配置文件\n", format)
			return nil
		}
	}

	return fmt.Errorf("无法在目录 %s 中读取配置文件，支持的格式: %v", configDir, supportedFormats)
}

// setupConfigWatch 设置配置文件监听
func setupConfigWatch(cfg *Config, onReload ...func(*Config)) {
	if len(onReload) > 0 && onReload[0] != nil {
		WatchConfig(cfg, func(newCfg any) {
			if config, ok := newCfg.(*Config); ok {
				onReload[0](config)
			}
		})
	}
}
