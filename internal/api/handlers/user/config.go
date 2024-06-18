package user

import (
	"time"

	"github.com/ndovnar/family-budget-api/internal/auth"
)

type Config struct {
	Auth                 *auth.Auth
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	Store                Store
}
