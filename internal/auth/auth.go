package auth

import "time"

type Auth struct {
	secretKey            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func New(cfg Config) *Auth {
	return &Auth{
		secretKey:            cfg.SecretKey,
		accessTokenDuration:  cfg.AccessTokenDuration,
		refreshTokenDuration: cfg.RefreshTokenDuration,
	}
}
