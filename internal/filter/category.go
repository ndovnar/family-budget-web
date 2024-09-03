package filter

type GetCategoriesFilter struct {
	BudgetID   string `form:"budget" binding:"required"`
	Pagination *Pagination
}
