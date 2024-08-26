package tokens

import (
	"context"

	"github.com/ndovnar/family-budget-api/internal/auth"
	"github.com/ndovnar/family-budget-api/internal/model"
)

type Tokens struct {
	store Store
	auth  *auth.Auth
}

type Store interface {
	GetSessionByID(ctx context.Context, id string) (*model.Session, error)
	CreateSession(ctx context.Context, params *model.Session) (*model.Session, error)
	DeleteSession(ctx context.Context, id string) error
}

func New(auth *auth.Auth, store Store) *Tokens {
	return &Tokens{
		auth:  auth,
		store: store,
	}
}
