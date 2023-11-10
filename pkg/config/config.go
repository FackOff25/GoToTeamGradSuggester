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

	DBurl      string
	DBport     string
	DBuser     string
	DBpassword string
	DBname     string

	LogAppName string
	LogLevel   int
	LogFormat  string
	LogOutput  string
	Categories []string
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

		DBurl:      viper.GetString("db.url"),
		DBport:     viper.GetString("db.port"),
		DBuser:     viper.GetString("db.user"),
		DBpassword: viper.GetString("db.password"),
		DBname:     viper.GetString("db.name"),

		LogAppName: "Suggest",
		LogLevel:   viper.GetInt("logger.level"),
		LogFormat:  viper.GetString("logger.format"),
		LogOutput:  viper.GetString("logger.output"),
		Categories: viper.GetStringSlice("categories"),
	}

	return &config, nil
}
