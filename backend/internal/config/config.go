package config

import (
	"fmt"
	"os"
	"time"

	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Kafka         *KafkaConfig         `yaml:"kafka"`
	Posts         *PostsConfig         `yaml:"posts"`
	JWT           *JWTConfig           `yaml:"jwt"`
	SQLServer     *DatabaseConfig      `yaml:"database"`
	API           *APIConfig           `yaml:"api"`
	Metrics       *MetricsConfig       `yaml:"metrics"`
	Cache         *CacheConfig         `yaml:"cache"`
	SocialNetwork *SocialNetworkConfig `yaml:"social_network"`
	Log           *LogConfig           `yaml:"log"`
	Dialogs       *DialogServiceConfig `yaml:"dialogs"`
	Tarantool     *TarantoolConfig     `yaml:"tarantool"`
}

type APIConfig struct {
	ServerAddress       string        `yaml:"server_address"`
	ContextTimeout      time.Duration `yaml:"context_timeout"`
	ShutdownTimeout     time.Duration `yaml:"shutdown_timeout"`
	FeedDefaultPageSize int           `yaml:"feed_default_page_size"`
	FeedMaxPageSize     int           `yaml:"feed_max_page_size"`
}
type CacheConfig struct {
	Redis              *RedisConfig  `yaml:"redis"`
	Expiry             time.Duration `yaml:"expiry"`
	CacheWarmupEnabled bool          `yaml:"cache_warmup_enabled"`
	CacheWarmupPeriod  time.Duration `yaml:"cache_warmup_period"`
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
	RoutingConfig           []*RouteConfig `yaml:"routing"`
	FeedLength              int            `yaml:"feed_length"`
	SvcDialogsURL           string         `yaml:"svc_dialogs_url"`
	SvcPostsWsURL           string         `yaml:"svc_posts_ws_url"`
	MaxPostCreatedPerWorker int            `yaml:"max_post_created_per_worker"`
	PostCreatedPacketSize   int            `yaml:"post_created_packet_size"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

type KafkaConfig struct {
	Brokers                  string `yaml:"brokers"`
	Acks                     string `yaml:"acks"`               // Гарантия доставки
	Retries                  int    `yaml:"retries"`            // Количество повторов в случае ошибки
	LingerMs                 int    `yaml:"linger_ms"`          // Снижение нагрузки за счет небольшого ожидания перед отправкой
	EnableIdempotence        bool   `yaml:"enable_idempotence"` // Идемпотентность продюсера
	TopicPostModified        string `yaml:"topic_post_modified"`
	TopicFeedChanged         string `yaml:"topic_feed_changed"`
	TopicFollowerNotify      string `yaml:"topic_follower_notify"`
	CGPostModified           string `yaml:"consumergroup_post_modified"`
	CGFeedChanged            string `yaml:"consumergroup_feed_changed"`
	CGFollowerNotify         string `yaml:"consumergroup_follower_notify"`
	NumWorkersPostModified   int    `yaml:"num_workers_post_modified"`
	NumWorkersFeedChanged    int    `yaml:"num_workers_feed_changed"`
	NumWorkersFollowerNotify int    `yaml:"num_workers_follower_notify"`
}

type PostsConfig struct {
	WebsocketPingInterval time.Duration `yaml:"websocket_ping_interval"` // как часто сервер пингует клиента
	WebsocketPongWait     time.Duration `yaml:"websocket_pong_wait"`     // время ожидания ответа (pong) от клиента

}

type TarantoolConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// LogConfig представляет настройки логирования
type LogConfig struct {
	Level log.Lvl `yaml:"level"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file %s: %w", configPath, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file %s: %w", configPath, err)
	}
	if config.Log == nil {
		return nil, fmt.Errorf("missing required 'log' section in file %s", configPath)
	}

	return &config, nil
}
