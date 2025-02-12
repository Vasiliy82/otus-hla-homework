package config_counters

import (
	"fmt"
)

// Validate проверяет конфигурацию и возвращает список ошибок
func (c *ConfigCounters) Validate() []error {
	var errs []error

	if c.Kafka != nil {
		errs = append(errs, c.Kafka.Validate()...)
	} else {
		errs = append(errs, fmt.Errorf("missing kafka configuration"))
	}

	if c.API != nil {
		errs = append(errs, c.API.Validate()...)
	} else {
		errs = append(errs, fmt.Errorf("missing api configuration"))
	}

	if c.Redis != nil {
		errs = append(errs, c.Redis.Validate()...)
	} else {
		errs = append(errs, fmt.Errorf("missing redis configuration"))
	}

	if c.Postgres != nil {
		errs = append(errs, c.Postgres.Validate()...)
	} else {
		errs = append(errs, fmt.Errorf("missing postgres configuration"))
	}

	return errs
}

// Validate проверяет конфигурацию Kafka
func (k *KafkaConfig) Validate() []error {
	var errs []error
	if k.Brokers == "" {
		errs = append(errs, fmt.Errorf("kafka: missing brokers"))
	}
	if k.TopicSagaBus == "" {
		errs = append(errs, fmt.Errorf("kafka: missing topic_saga_bus"))
	}
	if k.CGSagaBus == "" {
		errs = append(errs, fmt.Errorf("kafka: missing consumergroup_saga_bus"))
	}
	if k.NumWorkersSagaBus <= 0 {
		errs = append(errs, fmt.Errorf("kafka: num_workers_saga_bus must be greater than 0"))
	}
	return errs
}

// Validate проверяет конфигурацию API
func (a *APIConfig) Validate() []error {
	var errs []error
	if a.ServerAddress == "" {
		errs = append(errs, fmt.Errorf("api: missing server_address"))
	}
	if a.ContextTimeout <= 0 {
		errs = append(errs, fmt.Errorf("api: context_timeout must be greater than 0"))
	}
	if a.ShutdownTimeout <= 0 {
		errs = append(errs, fmt.Errorf("api: shutdown_timeout must be greater than 0"))
	}
	return errs
}

// Validate проверяет конфигурацию Redis
func (r *RedisConfig) Validate() []error {
	var errs []error
	if r.Host == "" {
		errs = append(errs, fmt.Errorf("redis: missing host"))
	}
	if r.Port <= 0 {
		errs = append(errs, fmt.Errorf("redis: invalid port"))
	}
	return errs
}

// Validate проверяет конфигурацию Postgres
func (d *DBInstanceConfig) Validate() []error {
	var errs []error
	if d.Host == "" {
		errs = append(errs, fmt.Errorf("postgres: missing host"))
	}
	if d.Port == "" {
		errs = append(errs, fmt.Errorf("postgres: missing port"))
	}
	if d.User == "" {
		errs = append(errs, fmt.Errorf("postgres: missing user"))
	}
	if d.Pass == "" {
		errs = append(errs, fmt.Errorf("postgres: missing password"))
	}
	if d.Name == "" {
		errs = append(errs, fmt.Errorf("postgres: missing database name"))
	}
	if d.MaxConns <= 0 {
		errs = append(errs, fmt.Errorf("postgres: max_conns must be greater than 0"))
	}
	if d.MinConns < 0 {
		errs = append(errs, fmt.Errorf("postgres: min_conns must be 0 or greater"))
	}
	if d.MaxConnIdleTime < 0 {
		errs = append(errs, fmt.Errorf("postgres: max_conn_idle_time must be 0 or greater"))
	}
	if d.MaxConnLifetime < 0 {
		errs = append(errs, fmt.Errorf("postgres: max_conn_lifetime must be 0 or greater"))
	}
	return errs
}
