package auth

import "time"

type Config struct {
	SecretKey            string        `required:"true"`
	AccessTokenDuration  time.Duration `required:"true"`
	RefreshTokenDuration time.Duration `required:"true"`
}
