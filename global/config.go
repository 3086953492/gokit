package global

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/3086953492/YaBase/configs"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 全局配置管理器
var (
	globalConfig  *configs.Config
	globalMutex   sync.RWMutex
	isInitialized bool
)

// GetGlobalConfig 获取全局配置
func GetGlobalConfig() *configs.Config {
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
	cfg := &configs.Config{}
	err := LoadConfig(cfg, func(newCfg *configs.Config) {
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

// LoadConfig 通用配置加载函数
// cfg: 配置结构指针，会被自动填充
// onReload: 可选的重载回调函数
func LoadConfig(cfg *configs.Config, onReload ...func(*configs.Config)) error {
	// 处理命令行参数
	var config string
	flag.StringVar(&config, "c", "", "指定配置文件路径")
	flag.Parse()

	// 如果用户通过命令行或环境变量指定了配置文件，使用指定的配置文件
	if config == "" {
		if configEnv := os.Getenv("CONFIG"); configEnv != "" {
			config = configEnv
			fmt.Printf("您正在使用环境变量, 配置文件的路径为%s\n", configEnv)
		} else if gin.Mode() != "" {
			// 根据gin模式自动选择配置文件
			switch gin.Mode() {
			case gin.DebugMode:
				config = "./configs/config.yaml"
			case gin.TestMode:
				config = "./configs/config.test.yaml"
			case gin.ReleaseMode:
				config = "./configs/config.release.yaml"
			}
			if config != "" {
				fmt.Printf("您正在使用gin模式的%s环境名称, 配置文件的路径为%s\n", gin.Mode(), config)
			}
		}
	} else {
		fmt.Printf("您正在使用命令行的-c参数传递的值, 配置文件的路径为%s\n", config)
	}

	// 如果指定了配置文件，直接使用指定的配置文件
	if config != "" {
		viper.SetConfigFile(config)
		err := viper.ReadInConfig()
		if err != nil {
			return fmt.Errorf("配置文件读取失败: %v", err)
		}

		// 监听配置文件变化
		if len(onReload) > 0 && onReload[0] != nil {
			viper.WatchConfig()
			viper.OnConfigChange(func(in fsnotify.Event) {
				fmt.Println("配置文件发生变更: ", in.Name)
				err := viper.Unmarshal(cfg)
				if err != nil {
					fmt.Printf("配置文件重新加载失败: %v\n", err)
				} else {
					fmt.Println("配置文件已重新加载")
					onReload[0](cfg)
				}
			})
		}

		err = viper.Unmarshal(cfg)
		if err != nil {
			return fmt.Errorf("配置文件解析失败: %v", err)
		}
		fmt.Printf("成功读取配置文件: %s\n", config)
		return nil
	}

	// 读取默认位置的配置文件
	viper.AddConfigPath("./configs")

	// 优先读取YAML配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		// 如果YAML读取失败，尝试读取JSON配置
		fmt.Printf("YAML配置读取失败，尝试读取JSON配置: %v\n", err)

		viper.SetConfigName("config")
		viper.SetConfigType("json")

		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("配置文件读取失败，YAML和JSON配置都无法读取: %v", err)
		}

		fmt.Println("成功读取JSON配置文件")
	} else {
		fmt.Println("成功读取YAML配置文件")
	}

	// 为默认配置文件也启用监听功能
	if len(onReload) > 0 && onReload[0] != nil {
		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			fmt.Println("配置文件发生变更: ", in.Name)
			err := viper.Unmarshal(cfg)
			if err != nil {
				fmt.Printf("配置文件重新加载失败: %v\n", err)
			} else {
				fmt.Println("配置文件已重新加载")
				onReload[0](cfg)
			}
		})
	}

	return viper.Unmarshal(cfg)
}
