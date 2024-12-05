package token

import "time"

type Config struct {
	SecretKey            []byte
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func NewConfig() *Config {
	return &Config{
		SecretKey:            []byte("your-secret-key"),
		AccessTokenDuration:  time.Hour,
		RefreshTokenDuration: time.Hour * 24 * 30,
	}
}
