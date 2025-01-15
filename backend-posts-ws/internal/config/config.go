package config

import (
	"os"
	"time"

	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	API       *APIConfig       `yaml:"api"`
	WebSocket *WebSocketConfig `yaml:"websocket"`
	Broker    *BrokerConfig    `yaml:"broker"`
	Log       *LogConfig       `yaml:"log"`
}

type APIConfig struct {
	ServerAddress   string        `yaml:"server_address"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type WebSocketConfig struct {
	PingInterval time.Duration `yaml:"ping_interval"` // как часто сервер пингует клиента
	PongWait     time.Duration `yaml:"pong_wait"`     // время ожидания ответа (pong) от клиента

}

type BrokerConfig struct {
	Brokers string `yaml:"brokers"`
	Group   string `yaml:"group"`
	Topic   string `yaml:"topic"`
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

// LogConfig представляет настройки логирования
type LogConfig struct {
	Level log.Lvl `yaml:"level"`
}
