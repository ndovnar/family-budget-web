package store

import (
	"context"

	"github.com/ndovnar/family-budget-api/internal/model"
)

type Store interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)

	GetSessionByID(ctx context.Context, id string) (*model.Session, error)
	CreateSession(ctx context.Context, session *model.Session) (*model.Session, error)

	CreateAccount(ctx context.Context, session *model.Account) (*model.Account, error)
}
