package configReader

import (
	"github.com/BurntSushi/toml"
	"github.com/FackOff25/GoToTeamGradSuggester/pkg/config"
)

type ConfigReader struct {
	configFilePath string
}

func NewConfigReader(configFilePath string) (config.ConfigReaderInterface) {
	return &ConfigReader{configFilePath: configFilePath}
}

func (cr *ConfigReader) ParseConfig() (*config.Config, error) {
	var config config.Config

	if _, err := toml.DecodeFile(cr.configFilePath, &config); err != nil {
		return nil, err
	} 
	return &config, nil
}
