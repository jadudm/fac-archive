package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() {
	atom := zap.NewAtomicLevel()

	// To keep the example deterministic, disable timestamps in the output.
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""

	var level zapcore.Level
	if viper.GetString("debug_level") == "DEBUG" {
		level = zapcore.DebugLevel
	} else if viper.GetString("debug_level") == "INFO" {
		level = zapcore.InfoLevel
	} else if viper.GetString("debug_level") == "WARN" {
		level = zapcore.WarnLevel
	} else if viper.GetString("debug_level") == "ERROR" {
		level = zapcore.ErrorLevel
	} else {
		level = zapcore.InfoLevel
	}

	atom.SetLevel(level)

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
	// defer undo()

}

func Init() {

	viper.SetConfigName("config")         // name of config file (without extension)
	viper.SetConfigType("yaml")           // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")              // look for config in the working directory first
	viper.AddConfigPath("$HOME/.factool") // then look in home

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Println("could not read config. exiting.")
	}

	InitLogger()

}
