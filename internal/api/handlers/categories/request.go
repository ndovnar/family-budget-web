package categories

type categoryRequest struct {
	BudgetID string  `json:"budgetId" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Currency string  `json:"currency" binding:"required"`
	Balance  float64 `json:"balance"`
}
