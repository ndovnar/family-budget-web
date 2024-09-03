package model

type Category struct {
	ID        string  `json:"id,omitempty" bson:"id,omitempty"`
	BudgetID  string  `json:"budget" bson:"budget"`
	Name      string  `json:"name" bson:"name"`
	Currency  string  `json:"currency" bson:"currency"`
	Balance   float64 `json:"balance" bson:"balance"`
	IsDeleted bool    `json:"deleted,omitempty" bson:"deleted"`
	Dates     Dates   `json:"dates" bson:"dates"`
}
