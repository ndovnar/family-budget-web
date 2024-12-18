package model

type Session struct {
	ID        string `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string `json:"userId" bson:"userId"`
	IsDeleted bool   `json:"deleted,omitempty" bson:"deleted"`
	Dates     Dates  `json:"dates" bson:"dates"`
}
