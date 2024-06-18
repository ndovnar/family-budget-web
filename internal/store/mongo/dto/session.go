package dto

import (
	"github.com/ndovnar/family-budget-api/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    string             `bson:"userId"`
	IsRevoked bool               `bson:"revoked"`
	Dates     Dates              `bson:"dates,omitempty"`
}

func ModelSessionToDtoSession(session *model.Session) *Session {
	return &Session{
		ID:        primitive.NewObjectID(),
		UserID:    session.UserID,
		IsRevoked: session.IsRevoked,
		Dates:     ModelDatesToDtoDates(session.Dates),
	}
}

func DtoSessionToModelSession(session *Session) *model.Session {
	return &model.Session{
		ID:        session.ID.Hex(),
		UserID:    session.UserID,
		IsRevoked: session.IsRevoked,
		Dates:     DtoDatesToModelDates(session.Dates),
	}
}
