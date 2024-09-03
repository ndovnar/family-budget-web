package accounts

import (
	"context"

	"github.com/ndovnar/family-budget-api/internal/auth"
	"github.com/ndovnar/family-budget-api/internal/authz"
	"github.com/ndovnar/family-budget-api/internal/filter"
	"github.com/ndovnar/family-budget-api/internal/model"
)

type Accounts struct {
	auth  *auth.Auth
	authz *authz.Authz
	store Store
}

type Store interface {
	GetAccounts(ctx context.Context, filter *filter.GetAccountsFilter) ([]*model.Account, int64, error)
	GetAccount(ctx context.Context, id string) (*model.Account, error)
	CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error)
	UpdateAccount(ctx context.Context, id string, account *model.Account) (*model.Account, error)
	DeleteAccount(ctx context.Context, id string) error
}

func New(auth *auth.Auth, authz *authz.Authz, store Store) *Accounts {
	return &Accounts{
		auth:  auth,
		authz: authz,
		store: store,
	}
}
