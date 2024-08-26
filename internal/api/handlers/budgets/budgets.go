package budgets

import (
	"context"

	"github.com/ndovnar/family-budget-api/internal/auth"
	"github.com/ndovnar/family-budget-api/internal/model"
)

type Budgets struct {
	store Store
	auth  *auth.Auth
}

type Store interface {
	GetBudgets(ctx context.Context, filter *model.GetBudgetsFilter) ([]*model.Budget, error)
	GetBudget(ctx context.Context, id string) (*model.Budget, error)
	CreateBudget(ctx context.Context, account *model.Budget) (*model.Budget, error)
	UpdateBudget(ctx context.Context, id string, account *model.Budget) (*model.Budget, error)
	DeleteBudget(ctx context.Context, id string) error
}

func New(auth *auth.Auth, store Store) *Budgets {
	return &Budgets{
		store: store,
		auth:  auth,
	}
}
