package categories

import (
	"context"

	"github.com/ndovnar/family-budget-api/internal/model"
)

type Categories struct {
	store Store
}

type Store interface {
	GetCategories(ctx context.Context) ([]*model.Category, error)
	GetCategory(ctx context.Context, id string) (*model.Category, error)
	CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error)
	UpdateCategory(ctx context.Context, id string, category *model.Category) (*model.Category, error)
	DeleteCategory(ctx context.Context, id string) error
}

func New(store Store) *Categories {
	return &Categories{
		store: store,
	}
}
