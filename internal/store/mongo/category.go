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

func (m *Mongo) GetCategories(ctx context.Context) ([]*model.Category, error) {
	res, err := m.database.
		Collection(CollectionCategories).
		Find(ctx, bson.M{}, nil)

	if err != nil {
		log.Error().Err(err).Msgf("mongo: failed to get categories. %v", err)
		return nil, err
	}

	dtoCategories := []*dto.Category{}
	for res.Next(ctx) {
		dtoCategory := &dto.Category{}
		err = res.Decode(&dtoCategory)
		if err != nil {
			return nil, err
		}

		dtoCategories = append(dtoCategories, dtoCategory)
	}

	return dto.DtoCategoriesToModelCategories(dtoCategories), nil
}

func (m *Mongo) GetCategory(ctx context.Context, id string) (*model.Category, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, store.ErrNotFound
	}

	filter := getNotDeletedByIDFilter(oid)
	return m.getCategory(ctx, filter)
}

func (m *Mongo) CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error) {
	currentTime := time.Now()

	dtoCategory := dto.ModelCategoryToDtoCategory(category)
	dtoCategory.Dates = dto.Dates{Created: &currentTime, Modified: &currentTime}

	_, err := m.database.
		Collection(CollectionCategories).
		InsertOne(ctx, dtoCategory)

	if err != nil {
		log.Error().Err(err).Msgf("mongo: failed to create category. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	createdCategory := dto.DtoCategoryToModelCategory(dtoCategory)

	return createdCategory, nil
}

func (m *Mongo) UpdateCategory(ctx context.Context, id string, category *model.Category) (*model.Category, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, store.ErrNotFound
	}

	filter := getNotDeletedByIDFilter(oid)
	currentTime := time.Now()
	dtoCategory := dto.ModelCategoryToDtoCategory(category)
	dtoCategory.Dates.Modified = &currentTime
	update := bson.M{
		"$set": bson.M{
			"name":           dtoCategory.Name,
			"budgetId":       dtoCategory.BudgetID,
			"dates.modified": dtoCategory.Dates.Modified,
			"balance":        dtoCategory.Balance,
			"currency":       dtoCategory.Currency,
		},
	}

	updateResult, err := m.
		database.
		Collection(CollectionCategories).
		UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if updateResult.MatchedCount == 0 {
		return nil, store.ErrNotFound
	}

	return m.getCategory(ctx, filter)
}

func (m *Mongo) DeleteCategory(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return store.ErrNotFound
	}

	filter := getNotDeletedByIDFilter(oid)
	update := bson.M{
		"$set": bson.M{
			"dates.modified": time.Now(),
			"dates.deleted":  time.Now(),
			"deleted":        true,
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

	category := dto.Category{}
	err := res.Decode(&category)

	if err != nil {
		log.Error().Err(err).Msgf("mongo getCategory: error while decoding the database object to a category. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	return dto.DtoCategoryToModelCategory(&category), nil
}
