package accounts

import (
	"context"

	"github.com/ndovnar/family-budget-api/internal/auth"
	"github.com/ndovnar/family-budget-api/internal/model"
)

type Accounts struct {
	store Store
	auth  *auth.Auth
}

type Store interface {
	GetAccounts(ctx context.Context, filter *model.GetAccountsFilter) ([]*model.Account, error)
	GetAccount(ctx context.Context, id string) (*model.Account, error)
	CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error)
	UpdateAccount(ctx context.Context, id string, account *model.Account) (*model.Account, error)
	DeleteAccount(ctx context.Context, id string) error
}

func New(auth *auth.Auth, store Store) *Accounts {
	return &Accounts{
		store: store,
		auth:  auth,
	}
}
