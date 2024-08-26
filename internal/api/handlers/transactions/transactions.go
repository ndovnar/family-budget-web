package transactions

import (
	"context"

	"github.com/ndovnar/family-budget-api/internal/auth"
	"github.com/ndovnar/family-budget-api/internal/model"
)

type Transactions struct {
	auth  *auth.Auth
	store Store
}

type Store interface {
	GetTransactions(ctx context.Context, filter *model.GetTransactionsFilter) ([]*model.Transaction, error)
	GetTransaction(ctx context.Context, id string) (*model.Transaction, error)
	CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error)
	UpdateTransaction(ctx context.Context, id string, transaction *model.Transaction) (*model.Transaction, error)
	DeleteTransaction(ctx context.Context, id string) error
}

func New(auth *auth.Auth, store Store) *Transactions {
	return &Transactions{
		auth:  auth,
		store: store,
	}
}
