package model

type Account struct {
	ID        string  `json:"id,omitempty," bson:"_id,omitempty"`
	OwnerID   string  `json:"owner" bson:"owner"`
	Name      string  `json:"name" bson:"name"`
	Balance   float64 `json:"balance" bson:"balance"`
	IsDeleted bool    `json:"deleted,omitempty" bson:"deleted"`
	Dates     Dates   `json:"dates" bson:"dates"`
}
