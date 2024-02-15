package core

import (
	"flag"
	"fmt"
	"market/core/internal"
	"market/global"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Viper(path ...string) *viper.Viper {
	var config string

	// search config path: command -> env -> default
	if len(path) == 0 {
		flag.StringVar(&config, "c", "", "choose a config file")
		flag.Parse()
		if config == "" {
			if configEnv := os.Getenv(internal.ConfigEnv); configEnv == "" {
				switch gin.Mode() {
				case gin.DebugMode:
					config = internal.ConfigDebugFile
				case gin.ReleaseMode:
					config = internal.ConfigReleaseFile
				case gin.TestMode:
					config = internal.ConfigTestFile
				}
			} else {
				config = configEnv
			}
		}
	} else {
		config = path[0]
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		panic(fmt.Errorf("error for config file: %s", err.Error()))
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		global.MARKET_LOG.Info(fmt.Sprintf("config file changed: %s", e.Name))
		if err = v.Unmarshal(&global.MARKET_CONFIG); err != nil {
			global.MARKET_LOG.Error(err.Error())
		}
	})
	if err = v.Unmarshal(&global.MARKET_CONFIG); err != nil {
		global.MARKET_LOG.Error(err.Error())
	}

	return v
}
