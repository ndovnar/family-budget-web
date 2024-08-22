package categories

import "github.com/ndovnar/family-budget-api/internal/model"

type categoryResponse struct {
	ID       string  `json:"id"`
	BudgetID string  `json:"budgetId"`
	Name     string  `json:"name"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

func newCategoryResponse(category *model.Category) *categoryResponse {
	return &categoryResponse{
		ID:       category.ID,
		BudgetID: category.BudgetID,
		Name:     category.Name,
		Balance:  category.Balance,
		Currency: category.Currency,
	}
}

func newCategoriesResponse(categories []*model.Category) []*categoryResponse {
	categoriesResponse := make([]*categoryResponse, 0, len(categories))
	for _, category := range categories {
		categoriesResponse = append(categoriesResponse, newCategoryResponse(category))
	}

	return categoriesResponse
}
