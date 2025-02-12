package config_counters

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type ConfigCounters struct {
	Kafka    *KafkaConfig      `yaml:"kafka"`
	API      *APIConfig        `yaml:"api"`
	Redis    *RedisConfig      `yaml:"redis"`
	Postgres *DBInstanceConfig `yaml:"postgres"`
}

type KafkaConfig struct {
	Brokers           string `yaml:"brokers"`
	Acks              string `yaml:"acks"`               // Гарантия доставки
	Retries           int    `yaml:"retries"`            // Количество повторов в случае ошибки
	LingerMs          int    `yaml:"linger_ms"`          // Снижение нагрузки за счет небольшого ожидания перед отправкой
	EnableIdempotence bool   `yaml:"enable_idempotence"` // Идемпотентность продюсера
	TopicSagaBus      string `yaml:"topic_saga_bus"`
	CGSagaBus         string `yaml:"consumergroup_saga_bus"`
	NumWorkersSagaBus int    `yaml:"num_workers_saga_bus"`
}

type APIConfig struct {
	ServerAddress   string        `yaml:"server_address"`
	ContextTimeout  time.Duration `yaml:"context_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

type DBInstanceConfig struct {
	Host            string        `yaml:"host"`
	Port            string        `yaml:"port"`
	User            string        `yaml:"user"`
	Pass            string        `yaml:"password"`
	Name            string        `yaml:"name"`
	MaxConns        int           `yaml:"max_conns"`
	MinConns        int           `yaml:"min_conns"`
	MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime"`
}

func LoadConfig(configPath string) (*ConfigCounters, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file %s: %w", configPath, err)
	}

	var config ConfigCounters
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file %s: %w", configPath, err)
	}

	return &config, nil
}
