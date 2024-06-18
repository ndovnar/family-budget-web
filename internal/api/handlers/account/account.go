package account

import (
	"context"

	"github.com/ndovnar/family-budget-api/internal/model"
)

type Account struct {
	store Store
}

type Store interface {
	CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error)
}

func NewAccount(store Store) *Account {
	return &Account{
		store: store,
	}
}
