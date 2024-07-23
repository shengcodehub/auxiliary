package config

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	FilePath = "./etc/"
)

var configViper = viper.New()
var paasViper = viper.New()

func Setup() {
	configViper.SetConfigName("home")
	configViper.SetConfigType("yaml")
	configViper.AddConfigPath(FilePath)
	if err := configViper.ReadInConfig(); err != nil {
		logx.Errorf("fatal error config file: %v", err)
	}
	configViper.OnConfigChange(
		func(e fsnotify.Event) {
			fmt.Println("------------------------------------------------------------")
			fmt.Println("Config file changed:", e.Name)
			fmt.Println("config has been reloaded at ",
				time.Now().Format("2006-01-02 15:04:05.000"))
			RestoreViper()
		},
	)
	configViper.WatchConfig()

	tmpViper := viper.New()
	if err := tmpViper.MergeConfigMap(configViper.AllSettings()); err != nil {
		logx.Errorf("merge config config file: %v", err)
	}
	if err := tmpViper.MergeConfigMap(paasViper.AllSettings()); err != nil {
		logx.Errorf("merge config passwd file: %v", err)
	}
	viperDefault := viper.GetViper()
	*viperDefault = *tmpViper
}

func RestoreViper() {
	tmpViper := viper.New()
	if err := tmpViper.MergeConfigMap(configViper.AllSettings()); err != nil {
		logx.Errorf("merge config config file: %v", err)
	}
	if err := tmpViper.MergeConfigMap(paasViper.AllSettings()); err != nil {
		logx.Errorf("merge config passwd file: %v", err)
	}
	viperDefault := viper.GetViper()
	*viperDefault = *tmpViper
}
