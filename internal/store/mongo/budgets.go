package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/ndovnar/family-budget-api/internal/filter"
	"github.com/ndovnar/family-budget-api/internal/model"
	"github.com/ndovnar/family-budget-api/internal/store"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/sync/errgroup"
)

func (m *Mongo) GetBudgets(ctx context.Context, budgetsFilter *filter.GetBudgetsFilter) ([]*model.Budget, int64, error) {
	errGroup, gCtx := errgroup.WithContext(ctx)
	collection := m.database.Collection(CollectionBudgets)
	filter := bson.M{
		"owner":   budgetsFilter.OwnerID,
		"deleted": budgetsFilter.Deleted,
	}
	paginationFindOptions := newPaginationFindOptions(budgetsFilter.Pagination)

	budgets := []*model.Budget{}
	errGroup.Go(func() error {
		cursor, err := collection.Find(ctx, filter, paginationFindOptions)
		if err != nil {
			log.Error().Err(err).Msgf("mongo: failed to get budgets. %v", err)
			return err
		}

		if err := cursor.All(gCtx, &budgets); err != nil {
			log.Error().Err(err).Msgf("mongo: failed to decode budgets. %v", err)
			return err
		}

		return nil
	})

	var totalCount int64
	errGroup.Go(func() error {
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			log.Error().Err(err).Msgf("mongo: failed to count budgets. %v", err)
			return err
		}

		totalCount = count
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return nil, 0, err
	}

	return budgets, totalCount, nil
}

func (m *Mongo) GetBudget(ctx context.Context, id string) (*model.Budget, error) {
	filter, err := newNotDeletedByIDFilter(id)
	if err != nil {
		return nil, store.ErrNotFound
	}

	return m.getBudget(ctx, filter)
}

func (m *Mongo) CreateBudget(ctx context.Context, budget *model.Budget) (*model.Budget, error) {
	currentTime := time.Now()
	budget.Dates = model.Dates{
		Created:  &currentTime,
		Modified: &currentTime,
	}

	result, err := m.database.
		Collection(CollectionBudgets).
		InsertOne(ctx, budget)
	if err != nil {
		log.Error().Err(err).Msgf("mongo: failed to create budget. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	newID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		err = fmt.Errorf("id of the inserted document %q is not an object id", result.InsertedID)
		return nil, err
	}

	budget.ID = newID.Hex()
	return budget, nil
}

func (m *Mongo) UpdateBudget(ctx context.Context, id string, budget *model.Budget) (*model.Budget, error) {
	filter, err := newNotDeletedByIDFilter(id)
	if err != nil {
		return nil, store.ErrNotFound
	}

	currentTime := time.Now()
	update := bson.M{
		"$set": bson.M{
			"name":           budget.Name,
			"dates.modified": &currentTime,
		},
	}

	result, err := m.
		database.
		Collection(CollectionBudgets).
		UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, store.ErrNotFound
	}

	return m.getBudget(ctx, filter)
}

func (m *Mongo) DeleteBudget(ctx context.Context, id string) error {
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

	budget := &model.Budget{}
	err := res.Decode(budget)

	if err != nil {
		log.Error().Err(err).Msgf("mongo getBudget: error while decoding the database object to a budget. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	return budget, nil
}
