package users

import (
	"context"

	"github.com/ndovnar/family-budget-api/internal/auth"
	"github.com/ndovnar/family-budget-api/internal/model"
)

type Users struct {
	store Store
	auth  *auth.Auth
}

type Store interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateSession(ctx context.Context, params *model.Session) (*model.Session, error)
	DeleteSession(ctx context.Context, id string) error
}

func New(auth *auth.Auth, store Store) *Users {
	return &Users{
		auth:  auth,
		store: store,
	}
}
