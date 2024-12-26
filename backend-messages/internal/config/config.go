package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SQLServer *DatabaseConfig      `yaml:"database"`
	API       *APIConfig           `yaml:"api"`
	Dialogs   *DialogServiceConfig `yaml:"dialogs"`
}

type APIConfig struct {
	ServerAddress   string        `yaml:"server_address"`
	ContextTimeout  time.Duration `yaml:"context_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}
type DialogServiceConfig struct {
	DefaultPageSize int `yaml:"default_page_size"`
	MaxPageSize     int `yaml:"max_page_size"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
