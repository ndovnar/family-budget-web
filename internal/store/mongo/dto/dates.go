package dto

import (
	"time"

	"github.com/ndovnar/family-budget-api/internal/model"
)

type Dates struct {
	Created  *time.Time `bson:"created"`
	Modified *time.Time `bson:"modified"`
	Deleted  *time.Time `bson:"deleted,omitempty"`
}

func ModelDatesToDtoDates(dates model.Dates) Dates {
	return Dates{
		Created:  dates.Created,
		Modified: dates.Modified,
		Deleted:  dates.Deleted,
	}
}

func DtoDatesToModelDates(dates Dates) model.Dates {
	return model.Dates{
		Created:  dates.Created,
		Modified: dates.Modified,
		Deleted:  dates.Deleted,
	}
}
