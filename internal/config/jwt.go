package config

import "time"

type JWTConfig struct {
	PrivateKey       string        `yaml:"private_key"`
	PublicKey        string        `yaml:"public_key"`
	TokenExpiry      time.Duration `yaml:"token_expiry"`      // время жизни токена
	PermissionsClaim string        `yaml:"permissions_claim"` // имя для claim прав доступа
}
