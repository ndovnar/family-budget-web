package budgets

import (
	"context"

	"github.com/ndovnar/family-budget-api/internal/auth"
	"github.com/ndovnar/family-budget-api/internal/authz"
	"github.com/ndovnar/family-budget-api/internal/filter"
	"github.com/ndovnar/family-budget-api/internal/model"
)

type Budgets struct {
	auth  *auth.Auth
	authz *authz.Authz
	store Store
}

type Store interface {
	GetBudgets(ctx context.Context, filter *filter.GetBudgetsFilter) ([]*model.Budget, int64, error)
	GetBudget(ctx context.Context, id string) (*model.Budget, error)
	CreateBudget(ctx context.Context, account *model.Budget) (*model.Budget, error)
	UpdateBudget(ctx context.Context, id string, account *model.Budget) (*model.Budget, error)
	DeleteBudget(ctx context.Context, id string) error
}

func New(auth *auth.Auth, authz *authz.Authz, store Store) *Budgets {
	return &Budgets{
		auth:  auth,
		authz: authz,
		store: store,
	}
}
