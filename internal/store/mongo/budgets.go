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

func (m *Mongo) GetBudgets(ctx context.Context, filter *store.GetBudgetsFilter) ([]*model.Budget, error) {
	res, err := m.database.
		Collection(CollectionBudgets).
		Find(ctx, bson.M{"owner": filter.Owner}, nil)

	if err != nil {
		log.Error().Err(err).Msgf("mongo: failed to get budgets. %v", err)
		return nil, err
	}

	dtoBudgets := []*dto.Budget{}
	for res.Next(ctx) {
		dtoBudget := &dto.Budget{}
		err = res.Decode(&dtoBudget)
		if err != nil {
			return nil, err
		}

		dtoBudgets = append(dtoBudgets, dtoBudget)
	}

	return dto.DtoBudgetsToModelBudgets(dtoBudgets), nil
}

func (m *Mongo) GetBudget(ctx context.Context, id string) (*model.Budget, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, store.ErrNotFound
	}

	filter := getNotDeletedByIDFilter(oid)
	return m.getBudget(ctx, filter)
}

func (m *Mongo) CreateBudget(ctx context.Context, budget *model.Budget) (*model.Budget, error) {
	currentTime := time.Now()

	dtoBudget := dto.ModelBudgetToDtoBudget(budget)
	dtoBudget.Dates = dto.Dates{Created: &currentTime, Modified: &currentTime}

	_, err := m.database.
		Collection(CollectionBudgets).
		InsertOne(ctx, dtoBudget)

	if err != nil {
		log.Error().Err(err).Msgf("mongo: failed to create budget. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	createdBudget := dto.DtoBudgetToModelBudget(dtoBudget)

	return createdBudget, nil
}

func (m *Mongo) UpdateBudget(ctx context.Context, id string, budget *model.Budget) (*model.Budget, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, store.ErrNotFound
	}

	filter := getNotDeletedByIDFilter(oid)
	currentTime := time.Now()
	dtoBudget := dto.ModelBudgetToDtoBudget(budget)
	dtoBudget.Dates.Modified = &currentTime
	update := bson.M{
		"$set": bson.M{
			"name":           dtoBudget.Name,
			"dates.modified": dtoBudget.Dates.Modified,
		},
	}

	updateResult, err := m.
		database.
		Collection(CollectionBudgets).
		UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if updateResult.MatchedCount == 0 {
		return nil, store.ErrNotFound
	}

	return m.getBudget(ctx, filter)
}

func (m *Mongo) DeleteBudget(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return store.ErrNotFound
	}

	filter := getNotDeletedByIDFilter(oid)
	update := bson.M{
		"$set": bson.M{
			"dates.modified": time.Now(),
			"deleted":        true,
		},
	}
	updateResult, err := m.database.
		Collection(CollectionBudgets).
		UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if updateResult.MatchedCount == 0 {
		return store.ErrNotFound
	}

	return nil
}

func (m *Mongo) getBudget(ctx context.Context, filter bson.M) (*model.Budget, error) {
	res := m.database.
		Collection(CollectionBudgets).
		FindOne(ctx, filter)

	budget := dto.Budget{}
	err := res.Decode(&budget)

	if err != nil {
		log.Error().Err(err).Msgf("mongo getBudget: error while decoding the database object to a budget. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	return dto.DtoBudgetToModelBudget(&budget), nil
}
