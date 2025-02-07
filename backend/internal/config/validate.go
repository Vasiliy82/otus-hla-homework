package config

import (
	"fmt"
)

// Config validation
func (c *Config) Validate() []error {
	var errs []error
	if c.Kafka != nil {
		errs = append(errs, c.Kafka.Validate()...)
	} else {
		errs = append(errs, fmt.Errorf("missing kafka configuration"))
	}

	if c.Posts != nil {
		errs = append(errs, c.Posts.Validate()...)
	} else {
		errs = append(errs, fmt.Errorf("missing posts configuration"))
	}

	if c.JWT != nil {
		errs = append(errs, c.JWT.Validate()...)
	} else {
		errs = append(errs, fmt.Errorf("missing jwt configuration"))
	}

	if c.SQLServer != nil {
		errs = append(errs, c.SQLServer.Validate()...)
	} else {
		errs = append(errs, fmt.Errorf("missing database configuration"))
	}

	if c.API != nil {
		errs = append(errs, c.API.Validate()...)
	} else {
		errs = append(errs, fmt.Errorf("missing api configuration"))
	}

	if c.Metrics != nil {
		errs = append(errs, c.Metrics.Validate()...)
	} else {
		errs = append(errs, fmt.Errorf("missing metrics configuration"))
	}

	if c.Cache != nil {
		errs = append(errs, c.Cache.Validate()...)
	} else {
		errs = append(errs, fmt.Errorf("missing cache configuration"))
	}

	if c.SocialNetwork != nil {
		errs = append(errs, c.SocialNetwork.Validate()...)
	} else {
		errs = append(errs, fmt.Errorf("missing social_network configuration"))
	}

	if c.Log == nil {
		errs = append(errs, fmt.Errorf("missing log configuration"))
	}

	if c.Dialogs != nil {
		errs = append(errs, c.Dialogs.Validate()...)
	} else {
		errs = append(errs, fmt.Errorf("missing dialogs configuration"))
	}

	if c.Tarantool != nil {
		errs = append(errs, c.Tarantool.Validate()...)
	} else {
		errs = append(errs, fmt.Errorf("missing tarantool configuration"))
	}

	return errs
}

// APIConfig validation
func (c *APIConfig) Validate() []error {
	var errs []error
	if c.ServerAddress == "" {
		errs = append(errs, fmt.Errorf("api.server_address must not be empty"))
	}
	if c.ContextTimeout <= 0 {
		errs = append(errs, fmt.Errorf("api.context_timeout must be greater than zero"))
	}
	if c.ShutdownTimeout <= 0 {
		errs = append(errs, fmt.Errorf("api.shutdown_timeout must be greater than zero"))
	}
	if c.FeedDefaultPageSize <= 0 {
		errs = append(errs, fmt.Errorf("api.feed_default_page_size must be greater than zero"))
	}
	if c.FeedMaxPageSize <= 0 || c.FeedMaxPageSize < c.FeedDefaultPageSize {
		errs = append(errs, fmt.Errorf("api.feed_max_page_size must be greater than zero and not less than feed_default_page_size"))
	}
	return errs
}

// CacheConfig validation
func (c *CacheConfig) Validate() []error {
	var errs []error
	if c.Redis == nil {
		errs = append(errs, fmt.Errorf("cache.redis configuration is missing"))
	} else {
		errs = append(errs, c.Redis.Validate()...)
	}
	if c.Expiry <= 0 {
		errs = append(errs, fmt.Errorf("cache.expiry must be greater than zero"))
	}
	return errs
}

// RedisConfig validation
func (c *RedisConfig) Validate() []error {
	var errs []error
	if c.Host == "" {
		errs = append(errs, fmt.Errorf("redis.host must not be empty"))
	}
	if c.Port <= 0 {
		errs = append(errs, fmt.Errorf("redis.port must be greater than zero"))
	}
	return errs
}

// JWTConfig validation
func (c *JWTConfig) Validate() []error {
	var errs []error
	if c.PrivateKey == "" {
		errs = append(errs, fmt.Errorf("jwt.private_key must not be empty"))
	}
	if c.PublicKey == "" {
		errs = append(errs, fmt.Errorf("jwt.public_key must not be empty"))
	}
	if c.TokenExpiry <= 0 {
		errs = append(errs, fmt.Errorf("jwt.token_expiry must be greater than zero"))
	}
	return errs
}

// MetricsConfig validation
func (c *MetricsConfig) Validate() []error {
	var errs []error
	if c.UpdateInterval <= 0 {
		errs = append(errs, fmt.Errorf("metrics.update_interval must be greater than zero"))
	}
	if len(c.BucketsHttpRequestDuration) == 0 {
		errs = append(errs, fmt.Errorf("metrics.buckets_http_request_duration must not be empty"))
	}
	return errs
}

// SocialNetworkConfig validation
func (c *SocialNetworkConfig) Validate() []error {
	var errs []error

	if len(c.RoutingConfig) == 0 {
		errs = append(errs, fmt.Errorf("routing configuration is empty"))
	} else {
		for i, route := range c.RoutingConfig {
			routeErrs := route.Validate()
			for _, err := range routeErrs {
				errs = append(errs, fmt.Errorf("route %d: %v", i+1, err))
			}
		}
	}

	if c.FeedLength <= 0 {
		errs = append(errs, fmt.Errorf("social_network.feed_length must be greater than zero"))
	}
	if c.MaxPostCreatedPerWorker <= 1 {
		errs = append(errs, fmt.Errorf("social_network.max_post_created_per_worker must be greater than 1"))
	}
	if c.PostCreatedPacketSize <= 0 {
		errs = append(errs, fmt.Errorf("social_network.post_created_packet_size must be greater than 1"))
	}
	if c.SvcDialogsURL == "" {
		errs = append(errs, fmt.Errorf("social_network.svc_dialogs_url must not be empty"))
	}
	if c.SvcPostsWsURL == "" {
		errs = append(errs, fmt.Errorf("social_network.svc_posts_ws_url must not be empty"))
	}
	return errs
}

// DialogServiceConfig validation
func (c *DialogServiceConfig) Validate() []error {
	var errs []error
	if c.DefaultPageSize <= 0 {
		errs = append(errs, fmt.Errorf("dialogs.default_page_size must be greater than zero"))
	}
	if c.MaxPageSize <= 0 || c.MaxPageSize < c.DefaultPageSize {
		errs = append(errs, fmt.Errorf("dialogs.max_page_size must be greater than zero and not less than default_page_size"))
	}
	// No need to validate UseInmem as it's a boolean
	return errs
}

// DatabaseConfig validation
func (c *DatabaseConfig) Validate() []error {
	var errs []error
	if c.Master == nil {
		errs = append(errs, fmt.Errorf("database.master configuration is missing"))
	} else {
		errs = append(errs, c.Master.Validate()...)
	}
	for i, replica := range c.Replicas {
		if replica == nil {
			errs = append(errs, fmt.Errorf("database.replica[%d] configuration is missing", i))
		} else {
			errs = append(errs, replica.Validate()...)
		}
	}
	return errs
}

// DBInstanceConfig validation
func (c *DBInstanceConfig) Validate() []error {
	var errs []error
	if c.Host == "" {
		errs = append(errs, fmt.Errorf("database.host must not be empty"))
	}
	if c.Port == "" {
		errs = append(errs, fmt.Errorf("database.port must not be empty"))
	}
	if c.User == "" {
		errs = append(errs, fmt.Errorf("database.user must not be empty"))
	}
	if c.Name == "" {
		errs = append(errs, fmt.Errorf("database.name must not be empty"))
	}
	if c.MaxConns <= 0 {
		errs = append(errs, fmt.Errorf("database.max_conns must be greater than zero"))
	}
	if c.MaxConnLifetime <= 0 {
		errs = append(errs, fmt.Errorf("database.max_conn_lifetime must be greater than zero"))
	}
	return errs
}

// KafkaConfig validation
func (c *KafkaConfig) Validate() []error {
	var errs []error
	if c.Brokers == "" {
		errs = append(errs, fmt.Errorf("kafka.brokers must not be empty"))
	}
	if c.Acks == "" {
		errs = append(errs, fmt.Errorf("kafka.acks must not be empty"))
	}
	if c.Retries < 0 {
		errs = append(errs, fmt.Errorf("kafka.retries must be zero or greater"))
	}
	if c.LingerMs < 0 {
		errs = append(errs, fmt.Errorf("kafka.linger_ms must be zero or greater"))
	}
	if c.TopicPostModified == "" {
		errs = append(errs, fmt.Errorf("kafka.topic_post_modified must not be empty"))
	}
	if c.TopicFeedChanged == "" {
		errs = append(errs, fmt.Errorf("kafka.topic_feed_changed must not be empty"))
	}
	if c.TopicFollowerNotify == "" {
		errs = append(errs, fmt.Errorf("kafka.topic_follower_notify must not be empty"))
	}
	if c.CGPostModified == "" {
		errs = append(errs, fmt.Errorf("kafka.consumergroup_post_modified must not be empty"))
	}
	if c.CGFeedChanged == "" {
		errs = append(errs, fmt.Errorf("kafka.consumergroup_feed_changed must not be empty"))
	}
	if c.CGFollowerNotify == "" {
		errs = append(errs, fmt.Errorf("kafka.consumergroup_follower_notify must not be empty"))
	}
	if c.NumWorkersPostModified <= 0 {
		errs = append(errs, fmt.Errorf("kafka.num_workers_post_modified must be greater than zero"))
	}
	if c.NumWorkersFeedChanged <= 0 {
		errs = append(errs, fmt.Errorf("kafka.num_workers_feed_changed must be greater than zero"))
	}

	if c.NumWorkersFollowerNotify <= 0 {
		errs = append(errs, fmt.Errorf("kafka.num_workers_follower_notify must be greater than zero"))
	}
	// No need to validate EnableIdempotence as it's a boolean
	return errs
}

// PostsConfig validation
func (c *PostsConfig) Validate() []error {
	var errs []error
	if c.WebsocketPingInterval <= 0 {
		errs = append(errs, fmt.Errorf("posts.websocket_ping_interval must be greater than zero"))
	}
	if c.WebsocketPongWait <= 0 {
		errs = append(errs, fmt.Errorf("posts.websocket_pong_wait must be greater than zero"))
	}
	return errs
}

// Validate проверяет корректность конфигурации Tarantool
func (c *TarantoolConfig) Validate() []error {
	var errs []error
	if c.Host == "" {
		errs = append(errs, fmt.Errorf("tarantool.host must not be empty"))
	}
	if c.Port <= 0 {
		errs = append(errs, fmt.Errorf("tarantool.port must be greater than zero"))
	}
	if c.User == "" {
		errs = append(errs, fmt.Errorf("tarantool.login must not be empty"))
	}
	if c.Password == "" {
		errs = append(errs, fmt.Errorf("tarantool.password must not be empty"))
	}
	return errs
}
