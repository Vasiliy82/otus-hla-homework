package config

import "time"

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

type DatabaseConfig struct {
	Master   *DBInstanceConfig
	Replicas []*DBInstanceConfig
}
