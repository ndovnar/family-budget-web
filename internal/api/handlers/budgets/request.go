package budgets

type createBudgetRequest struct {
	Name string `json:"name" binding:"required"`
}

type updateBudgetRequest struct {
	ID   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}
