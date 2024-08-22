package model

type Budget struct {
	ID        string
	Owner     string
	Name      string
	IsDeleted bool
	Dates     Dates
}
