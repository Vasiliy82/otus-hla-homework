package config

import (
	"os"
	"time"

	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	WebSocket     *WebSocketConfig     `yaml:"websocket"`
	JWT           *JWTConfig           `yaml:"jwt"`
	SQLServer     *DatabaseConfig      `yaml:"database"`
	API           *APIConfig           `yaml:"api"`
	Metrics       *MetricsConfig       `yaml:"metrics"`
	Cache         *CacheConfig         `yaml:"cache"`
	SocialNetwork *SocialNetworkConfig `yaml:"social_network"`
	Log           *LogConfig           `yaml:"log"`
	Dialogs       *DialogServiceConfig `yaml:"dialogs"`
}

type APIConfig struct {
	ServerAddress       string        `yaml:"server_address"`
	ContextTimeout      time.Duration `yaml:"context_timeout"`
	ShutdownTimeout     time.Duration `yaml:"shutdown_timeout"`
	FeedDefaultPageSize int           `yaml:"feed_default_page_size"`
	FeedMaxPageSize     int           `yaml:"feed_max_page_size"`
}
type CacheConfig struct {
	Redis                   *RedisConfig  `yaml:"redis"`
	Kafka                   *KafkaConfig  `yaml:"kafka"`
	Expiry                  time.Duration `yaml:"expiry"`
	InvalidateNumWorkers    int           `yaml:"invalidate_num_workers"`
	InvalidateTopic         string        `yaml:"invalidate_topic"`
	InvalidateConsumerGroup string        `yaml:"invalidate_consumer_group"`
	CacheWarmupEnabled      bool          `yaml:"cache_warmup_enabled"`
	CacheWarmupPeriod       time.Duration `yaml:"cache_warmup_period"`
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
type SocialNetworkConfig struct {
	FeedLength int `yaml:"feed_length"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

type KafkaConfig struct {
	Brokers           string `yaml:"brokers"`
	Acks              string `yaml:"acks"`               // Гарантия доставки
	Retries           int    `yaml:"retries"`            // Количество повторов в случае ошибки
	LingerMs          int    `yaml:"linger_ms"`          // Снижение нагрузки за счет небольшого ожидания перед отправкой
	EnableIdempotence bool   `yaml:"enable_idempotence"` // Идемпотентность продюсера
}

type WebSocketConfig struct {
	PingInterval time.Duration `yaml:"ping_interval"` // как часто сервер пингует клиента
	PongWait     time.Duration `yaml:"pong_wait"`     // время ожидания ответа (pong) от клиента

}

// LogConfig представляет настройки логирования
type LogConfig struct {
	Level log.Lvl `yaml:"level"`
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
