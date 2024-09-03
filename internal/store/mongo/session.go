package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/ndovnar/family-budget-api/internal/store"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *Mongo) GetSessionByID(ctx context.Context, id string) (*model.Session, error) {
	filter, err := newByIDFilter(id)
	if err != nil {
		return nil, store.ErrNotFound
	}

	return m.getSession(ctx, filter)
}

func (m *Mongo) CreateSession(ctx context.Context, session *model.Session) (*model.Session, error) {
	currentTime := time.Now()
	session.Dates = model.Dates{
		Created:  &currentTime,
		Modified: &currentTime,
	}

	result, err := m.database.
		Collection(CollectionSessions).
		InsertOne(ctx, session)
	if err != nil {
		log.Info().Msgf("mongo: failed to create session. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	newID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		err = fmt.Errorf("id of the inserted document %q is not an object id", result.InsertedID)
		return nil, err
	}

	session.ID = newID.Hex()
	return session, nil
}

func (m *Mongo) DeleteSession(ctx context.Context, id string) error {
	filter, err := newNotDeletedByIDFilter(id)
	if err != nil {
		return store.ErrNotFound
	}

	update := bson.M{
		"$set": bson.M{
			"deleted":       true,
			"dates.deleted": time.Now(),
		},
	}

	updateResult, err := m.database.
		Collection(CollectionSessions).
		UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if updateResult.MatchedCount == 0 {
		return store.ErrNotFound
	}

	return nil
}

func (m *Mongo) getSession(ctx context.Context, filter bson.M) (*model.Session, error) {
	res := m.database.
		Collection(CollectionSessions).
		FindOne(ctx, filter)

	session := &model.Session{}
	err := res.Decode(session)

	if err != nil {
		log.Error().Err(err).Msgf("mongo getSession: error while decoding the database object to a session. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	return session, nil
}
