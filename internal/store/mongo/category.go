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

func (m *Mongo) GetCategories(ctx context.Context) ([]*model.Category, error) {
	res, err := m.database.
		Collection(CollectionCategories).
		Find(ctx, bson.M{}, nil)

	if err != nil {
		log.Error().Err(err).Msgf("mongo: failed to get categories. %v", err)
		return nil, err
	}

	categories := []*model.Category{}
	for res.Next(ctx) {
		category := &model.Category{}
		err = res.Decode(category)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (m *Mongo) GetCategory(ctx context.Context, id string) (*model.Category, error) {
	filter, err := getNotDeletedByIDFilter(id)
	if err != nil {
		return nil, store.ErrNotFound
	}

	return m.getCategory(ctx, filter)
}

func (m *Mongo) CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error) {
	currentTime := time.Now()
	category.Dates = model.Dates{
		Created:  &currentTime,
		Modified: &currentTime,
	}

	result, err := m.database.
		Collection(CollectionCategories).
		InsertOne(ctx, category)
	if err != nil {
		log.Error().Err(err).Msgf("mongo: failed to create category. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	newID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		err = fmt.Errorf("id of the inserted document %q is not an object id", result.InsertedID)
		return nil, err
	}

	category.ID = newID.Hex()
	return category, nil
}

func (m *Mongo) UpdateCategory(ctx context.Context, id string, category *model.Category) (*model.Category, error) {
	filter, err := getNotDeletedByIDFilter(id)
	if err != nil {
		return nil, store.ErrNotFound
	}

	currentTime := time.Now()
	update := bson.M{
		"$set": bson.M{
			"name":           category.Name,
			"budget":         category.Budget,
			"dates.modified": &currentTime,
			"balance":        category.Balance,
			"currency":       category.Currency,
		},
	}

	result, err := m.
		database.
		Collection(CollectionCategories).
		UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, store.ErrNotFound
	}

	return m.getCategory(ctx, filter)
}

func (m *Mongo) DeleteCategory(ctx context.Context, id string) error {
	filter, err := getNotDeletedByIDFilter(id)
	if err != nil {
		return store.ErrNotFound
	}

	update := bson.M{
		"$set": bson.M{
			"dates.deleted": time.Now(),
			"deleted":       true,
		},
	}
	updateResult, err := m.database.
		Collection(CollectionCategories).
		UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if updateResult.MatchedCount == 0 {
		return store.ErrNotFound
	}

	return nil
}

func (m *Mongo) getCategory(ctx context.Context, filter bson.M) (*model.Category, error) {
	res := m.database.
		Collection(CollectionCategories).
		FindOne(ctx, filter)

	category := &model.Category{}
	err := res.Decode(category)

	if err != nil {
		log.Error().Err(err).Msgf("mongo getCategory: error while decoding the database object to a category. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	return category, nil
}
