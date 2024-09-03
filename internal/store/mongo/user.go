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

func (m *Mongo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	filter, err := newNotDeletedByIDFilter(id)
	if err != nil {
		return nil, store.ErrNotFound
	}

	return m.getUser(ctx, filter)
}

func (m *Mongo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	filter := newByEmailFilter(email)
	return m.getUser(ctx, filter)
}

func (m *Mongo) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	currentTime := time.Now()
	user.Dates = model.Dates{
		Created:  &currentTime,
		Modified: &currentTime,
	}

	result, err := m.database.
		Collection(CollectionUsers).
		InsertOne(ctx, user)
	if err != nil {
		log.Info().Msgf("mongo: failed to create user. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	newID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		err = fmt.Errorf("id of the inserted document %q is not an object id", result.InsertedID)
		return nil, err
	}

	user.ID = newID.Hex()
	return user, nil
}

func (m *Mongo) getUser(ctx context.Context, filter bson.M) (*model.User, error) {
	res := m.database.
		Collection(CollectionUsers).
		FindOne(ctx, filter)

	user := &model.User{}
	err := res.Decode(user)

	if err != nil {
		log.Error().Err(err).Msgf("mongo getUser: error while decoding the database object to a user. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	return user, nil
}
