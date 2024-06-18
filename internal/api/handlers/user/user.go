package user

import (
	"context"

	"github.com/ndovnar/family-budget-api/internal/auth"
	"github.com/ndovnar/family-budget-api/internal/model"
)

type User struct {
	store Store
	auth  *auth.Auth
}

type Store interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateSession(ctx context.Context, params *model.Session) (*model.Session, error)
}

func NewUser(auth *auth.Auth, store Store) *User {
	return &User{
		auth:  auth,
		store: store,
	}
}
