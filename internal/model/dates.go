package model

import "time"

type Dates struct {
	Created  *time.Time `json:"created,omitempty" bson:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty" bson:"modified,omitempty"`
	Deleted  *time.Time `json:"deleted,omitempty" bson:"deleted,omitempty"`
}
