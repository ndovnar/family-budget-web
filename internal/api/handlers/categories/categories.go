package categories

import (
	"context"

	"github.com/ndovnar/family-budget-api/internal/authz"
	"github.com/ndovnar/family-budget-api/internal/filter"
	"github.com/ndovnar/family-budget-api/internal/model"
)

type Categories struct {
	authz *authz.Authz
	store Store
}

type Store interface {
	GetCategories(ctx context.Context, filter *filter.GetCategoriesFilter) ([]*model.Category, int64, error)
	GetCategory(ctx context.Context, id string) (*model.Category, error)
	CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error)
	UpdateCategory(ctx context.Context, id string, category *model.Category) (*model.Category, error)
	DeleteCategory(ctx context.Context, id string) error
}

func New(authz *authz.Authz, store Store) *Categories {
	return &Categories{
		authz: authz,
		store: store,
	}
}
