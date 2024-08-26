package budgets

type budgetRequest struct {
	Name string `json:"name" binding:"required"`
}
