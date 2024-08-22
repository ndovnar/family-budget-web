package budgets

import (
	"github.com/ndovnar/family-budget-api/internal/model"
)

type budgetResponse struct {
	ID    string      `json:"id"`
	Owner string      `json:"owner"`
	Name  string      `json:"name"`
	Dates model.Dates `json:"dates"`
}

func newBudgetResponse(budget *model.Budget) *budgetResponse {
	return &budgetResponse{
		ID:    budget.ID,
		Owner: budget.Owner,
		Name:  budget.Name,
		Dates: budget.Dates,
	}
}

func newBudgetsResponse(budgets []*model.Budget) []*budgetResponse {
	budgetsResponse := make([]*budgetResponse, 0, len(budgets))
	for _, budget := range budgets {
		budgetsResponse = append(budgetsResponse, newBudgetResponse(budget))
	}

	return budgetsResponse
}
