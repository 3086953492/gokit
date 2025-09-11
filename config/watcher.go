package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func WatchConfig(cfg any, onReload ...func(any)) {
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
}
