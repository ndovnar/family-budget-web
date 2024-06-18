package dto

import (
	"github.com/ndovnar/family-budget-api/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	FirstName string             `bson:"firstName"`
	LastName  string             `bson:"lastName"`
	IsDeleted bool               `bson:"deleted"`
	Dates     Dates              `bson:"dates,omitempty"`
}

func ModelUserToDtoUser(user *model.User) *User {
	return &User{
		ID:        primitive.NewObjectID(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
		IsDeleted: user.IsDeleted,
		Dates:     ModelDatesToDtoDates(user.Dates),
	}
}

func DtoUserToModelUser(user *User) *model.User {
	return &model.User{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
		IsDeleted: user.IsDeleted,
		Dates:     DtoDatesToModelDates(user.Dates),
	}
}
