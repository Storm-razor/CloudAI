package config

import (
	"log"

	"github.com/spf13/viper"
)

var AppConfigInstance *AppConfig

func InitConfig() {
	AppConfigInstance = &AppConfig{}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./etc")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// 热重载(未使用)
	// viper.WatchConfig()
	// viper.OnConfigChange(func(e fsnotify.Event) {
	// 	if err := viper.Unmarshal(AppConfigInstance); err != nil {
	// 		log.Printf("loadConfig failed, unmarshal config err: %v", err)
	// 	}
	// })

	if err := viper.Unmarshal(AppConfigInstance); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
}

// GetConfig 获取配置
func GetConfig() *AppConfig {
	return AppConfigInstance
}
