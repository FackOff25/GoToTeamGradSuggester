package config

type Config struct {
	ServerAddress string `toml:"serverAddres"`
	ServerPort    string `toml:"serverPort"`
}

type ConfigReaderInterface interface {
	ParseConfig() (*Config, error)
}