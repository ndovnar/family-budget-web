package categories

import "github.com/ndovnar/family-budget-api/internal/model"

type categoriesResponse struct {
	Values []*model.Category
	Meta   *meta
}

type meta struct {
	Count int64 `json:"count"`
}

func newCategoriesResponse(categories []*model.Category, count int64) *categoriesResponse {
	return &categoriesResponse{
		Values: categories,
		Meta: &meta{
			Count: count,
		},
	}
}
