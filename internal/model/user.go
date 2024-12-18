package model

type User struct {
	ID        string `json:"id,omitempty" bson:"_id,omitempty"`
	Email     string `json:"email" bson:"email"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Password  string `json:"-" bson:"password"`
	IsDeleted bool   `json:"deleted,omitempty" bson:"deleted"`
	Dates     Dates  `json:"dates" bson:"dates"`
}
