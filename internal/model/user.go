package model

type User struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	Password  string
	IsDeleted bool
	Dates     Dates
}
