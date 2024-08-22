package dto

import (
	"github.com/ndovnar/family-budget-api/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Budget struct {
	ID        primitive.ObjectID `bson:"_id"`
	Owner     string             `bson:"owner"`
	Name      string             `bson:"name"`
	IsDeleted bool               `bson:"deleted"`
	Dates     Dates              `bson:"dates,omitempty"`
}

func ModelBudgetToDtoBudget(budget *model.Budget) *Budget {
	return &Budget{
		ID:        primitive.NewObjectID(),
		Owner:     budget.Owner,
		Name:      budget.Name,
		IsDeleted: budget.IsDeleted,
		Dates:     ModelDatesToDtoDates(budget.Dates),
	}
}

func DtoBudgetToModelBudget(budget *Budget) *model.Budget {
	return &model.Budget{
		ID:        budget.ID.Hex(),
		Owner:     budget.Owner,
		Name:      budget.Name,
		IsDeleted: budget.IsDeleted,
		Dates:     DtoDatesToModelDates(budget.Dates),
	}
}

func DtoBudgetsToModelBudgets(src []*Budget) []*model.Budget {
	dst := make([]*model.Budget, 0, len(src))
	for _, item := range src {
		dst = append(dst, DtoBudgetToModelBudget(item))
	}
	return dst
}
