package mongo

import (
	"context"
	"time"

	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/ndovnar/family-budget-api/internal/store"
	"github.com/ndovnar/family-budget-api/internal/store/mongo/dto"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *Mongo) GetSessionByID(ctx context.Context, id string) (*model.Session, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, store.ErrNotFound
	}

	filter := getByIDFilter(oid)
	res := m.database.
		Collection(CollectionSessions).
		FindOne(ctx, filter)

	session := dto.Session{}
	err = res.Decode(&session)

	if err != nil {
		log.Error().Err(err).Msgf("mongo getSessionByID: error while decoding the database object to a session. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	return dto.DtoSessionToModelSession(&session), nil
}

func (m *Mongo) CreateSession(ctx context.Context, session *model.Session) (*model.Session, error) {
	currentTime := time.Now()

	dtoSession := dto.ModelSessionToDtoSession(session)
	dtoSession.Dates = dto.Dates{Created: &currentTime, Modified: &currentTime}
	dtoSession.IsRevoked = false

	_, err := m.database.
		Collection(CollectionSessions).
		InsertOne(ctx, dtoSession)

	if err != nil {
		log.Info().Msgf("mongo: failed to create session. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	return dto.DtoSessionToModelSession(dtoSession), nil
}
