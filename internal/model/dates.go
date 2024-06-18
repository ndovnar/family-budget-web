package model

import "time"

type Dates struct {
	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`
	Expires  *time.Time `json:"expires,omitempty"`
	Deleted  *time.Time `json:"deleted,omitempty"`
}
