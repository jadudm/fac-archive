package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Logger, _ = zap.NewProduction()

func Init() {

	viper.SetConfigName("config")         // name of config file (without extension)
	viper.SetConfigType("yaml")           // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory
	err := viper.ReadInConfig()           // Find and read the config file
	if err != nil {
		Logger.Error("could not read config. exiting.")
	}
}
