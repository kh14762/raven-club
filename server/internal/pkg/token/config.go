package token

import "time"

type Config struct {
	SecretKey            []byte
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func NewConfig() *Config {
	return &Config{
		SecretKey:            []byte("iPQpgJWX6ccgPR9aIq5+2AygsmXKXYhXS7JUJWGOdhM="), // TODO get from .env file
		AccessTokenDuration:  time.Hour,
		RefreshTokenDuration: time.Hour * 24 * 30,
	}
}
