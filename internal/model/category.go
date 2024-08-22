package model

type Category struct {
	ID        string
	BudgetID  string
	Name      string
	Currency  string
	Balance   float64
	IsDeleted bool
	Dates     Dates
}
