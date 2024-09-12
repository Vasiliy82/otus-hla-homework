package config

type APIConfig struct {
	ServerAddress  string `yaml:"server_address"`
	ContextTimeout int    `yaml:"context_timeout"`
}
