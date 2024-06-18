package model

type Session struct {
	ID        string
	UserID    string
	IsRevoked bool
	Dates     Dates
}
