package model

type Account struct {
	ID        string
	Name      string
	Balance   float64
	IsDeleted bool
	Dates     Dates
}
