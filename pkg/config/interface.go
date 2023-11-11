package config

type ConfigReaderInterface interface {
	GetConfig() (*Config, error)
}
