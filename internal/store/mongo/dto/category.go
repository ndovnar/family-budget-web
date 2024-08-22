package dto

import (
	"github.com/ndovnar/family-budget-api/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID        primitive.ObjectID `bson:"_id"`
	BudgetID  string             `bson:"budgetId"`
	Name      string             `bson:"name"`
	Currency  string             `bson:"currency"`
	Balance   float64            `bson:"balance"`
	IsDeleted bool               `bson:"deleted"`
	Dates     Dates              `bson:"dates,omitempty"`
}

func ModelCategoryToDtoCategory(category *model.Category) *Category {
	return &Category{
		ID:        primitive.NewObjectID(),
		BudgetID:  category.BudgetID,
		Name:      category.Name,
		Currency:  category.Currency,
		Balance:   category.Balance,
		IsDeleted: category.IsDeleted,
		Dates:     ModelDatesToDtoDates(category.Dates),
	}
}

func DtoCategoryToModelCategory(category *Category) *model.Category {
	return &model.Category{
		ID:        category.ID.Hex(),
		BudgetID:  category.BudgetID,
		Name:      category.Name,
		Currency:  category.Currency,
		Balance:   category.Balance,
		IsDeleted: category.IsDeleted,
		Dates:     DtoDatesToModelDates(category.Dates),
	}
}

func DtoCategoriesToModelCategories(src []*Category) []*model.Category {
	dst := make([]*model.Category, 0, len(src))
	for _, item := range src {
		innerItem := item
		dst = append(dst, DtoCategoryToModelCategory(innerItem))
	}
	return dst
}
