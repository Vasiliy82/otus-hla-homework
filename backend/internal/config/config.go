package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	JWT       *JWTConfig      `yaml:"jwt"`
	SQLServer *DatabaseConfig `yaml:"database"`
	API       *APIConfig      `yaml:"api"`
	Metrics   *MetricsConfig  `yaml:"metrics"`
}

type APIConfig struct {
	ServerAddress   string        `yaml:"server_address"`
	ContextTimeout  time.Duration `yaml:"context_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type JWTConfig struct {
	PrivateKey  string        `yaml:"private_key"`
	PublicKey   string        `yaml:"public_key"`
	TokenExpiry time.Duration `yaml:"token_expiry"` // время жизни токена
}

type MetricsConfig struct {
	UpdateInterval             time.Duration `yaml:"update_interval"`
	BucketsHttpRequestDuration []float64     `yaml:"buckets_http_request_duration"`
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
