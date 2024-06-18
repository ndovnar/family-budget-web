package mongo

import (
	"context"
	"time"

	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/ndovnar/family-budget-api/internal/store"
	"github.com/ndovnar/family-budget-api/internal/store/mongo/dto"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *Mongo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, store.ErrNotFound
	}

	filter := getByIDFilter(oid)
	return m.getUser(ctx, filter)
}

func (m *Mongo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	filter := getByEmailFilter(email)
	return m.getUser(ctx, filter)
}

func (m *Mongo) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	currentTime := time.Now()

	dtoUser := dto.ModelUserToDtoUser(user)
	dtoUser.Dates = dto.Dates{Created: &currentTime, Modified: &currentTime}
	dtoUser.IsDeleted = false

	_, err := m.database.
		Collection(CollectionUsers).
		InsertOne(ctx, dtoUser)

	if err != nil {
		log.Info().Msgf("mongo: failed to create user. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	return dto.DtoUserToModelUser(dtoUser), nil
}

func (m *Mongo) getUser(ctx context.Context, filter bson.M) (*model.User, error) {
	res := m.database.
		Collection(CollectionUsers).
		FindOne(ctx, filter)

	user := dto.User{}
	err := res.Decode(&user)

	if err != nil {
		log.Error().Err(err).Msgf("mongo getUser: error while decoding the database object to a user. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	return dto.DtoUserToModelUser(&user), nil
}
