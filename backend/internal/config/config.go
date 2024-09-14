package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	JWT       *JWTConfig      `yaml:"jwt"`
	SQLServer *DatabaseConfig `yaml:"database"`
	API       *APIConfig      `yaml:"api"`
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
