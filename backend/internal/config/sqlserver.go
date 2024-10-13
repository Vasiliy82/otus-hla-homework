package config

import "time"

type DBInstanceConfig struct {
	Host            string        `yaml:"host"`
	Port            string        `yaml:"port"`
	User            string        `yaml:"user"`
	Pass            string        `yaml:"password"`
	Name            string        `yaml:"name"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time"`
}

type DatabaseConfig struct {
	Master   *DBInstanceConfig
	Replicas []*DBInstanceConfig
}
