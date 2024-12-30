package model

type Transaction struct {
	ID            string  `json:"id,omitempty" bson:"_id,omitempty"`
	FromAccountID string  `json:"fromAccount" bson:"fromAccount"`
	ToAccountID   string  `json:"toAccount" bson:"toAccount"`
	CategoryID    string  `json:"category" bson:"category"`
	UserID        string  `json:"user" bson:"user"`
	Amount        float64 `json:"amount" bson:"amount"`
	Description   string  `json:"description" bson:"description"`
	IsDeleted     bool    `json:"deleted,omitempty" bson:"deleted"`
	Dates         Dates   `json:"dates" bson:"dates"`
}
