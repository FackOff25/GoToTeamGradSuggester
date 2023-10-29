package config

import (
	"path/filepath"

	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string
	ServerPort    string
	PlacesApiHost string
	LogAppName    string
	LogLevel      int
	LogFormat     string
	LogOutput     string
}

func GetConfig(configPath string) (*Config, error) {
	viper.SetConfigName(filepath.Base(configPath))
	viper.AddConfigPath(filepath.Dir(configPath))
	viper.SetConfigType(filepath.Ext(configPath)[1:])

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("%s", err)
	}

	config := Config{
		ServerAddress: viper.GetString("server.address"),
		ServerPort:    viper.GetString("server.port"),
		PlacesApiHost: viper.GetString("placesApi.host"),
		LogAppName:    "Suggest",
		LogLevel:      viper.GetInt("logger.level"),
		LogFormat:     viper.GetString("logger.format"),
		LogOutput:     viper.GetString("logger.output"),
	}

	return &config, nil
}
