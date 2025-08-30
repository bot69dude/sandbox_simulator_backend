package config

import "time"

type JWTConfig struct {
	SecretKey           string
	AccessTokenDuration time.Duration
}

func NewJWTConfig(secretKey string) *JWTConfig {
	return &JWTConfig{
		SecretKey:           secretKey,
		AccessTokenDuration: 24 * time.Hour, // tokens valid for 24 hours
	}
}
