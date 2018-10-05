package iamhere

import (
	. "github.com/520lly/iamhere/app/modules"
	"github.com/fsnotify/fsnotify"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

var Config Configuration
var logger echo.Logger

func loadConfigFile() error {
	viper.SetConfigName("app")
	viper.AddConfigPath("/etc/iamhere/config")
	if err := viper.ReadInConfig(); err != nil {
		logger.Error("Error reading config file, %s", err)
		return err
	}
	err := viper.Unmarshal(&Config)
	if err != nil {
		logger.Error("unable to decode into struct,", err)
		return err
	}
	logger.Debug("Configuration:", Config)
	return nil
}

func InitConfig(log echo.Logger) error {
	logger = log
	logger.Debug("Initialization of config")

	if err := loadConfigFile(); err != nil {
		logger.Error("Error loadConfigFile config file, %s", err)
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Warn("Config file changed:", e.Name)
		if err := loadConfigFile(); err != nil {
			logger.Error("Error loadConfigFile config file, %s", err)
		}
	})
	return nil
}
