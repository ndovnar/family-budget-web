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

func (m *Mongo) GetAccounts(ctx context.Context, accountsFilter *filter.GetAccountsFilter) ([]*model.Account, int64, error) {
	errGroup, gCtx := errgroup.WithContext(ctx)
	collection := m.database.Collection(CollectionAccounts)
	filter := bson.M{
		"owner":   accountsFilter.OwnerID,
		"deleted": accountsFilter.Deleted,
	}
	paginationFindOptions := newPaginationFindOptions(accountsFilter.Pagination)

	accounts := []*model.Account{}
	errGroup.Go(func() error {
		cursor, err := collection.Find(ctx, filter, paginationFindOptions)
		if err != nil {
			log.Error().Err(err).Msgf("mongo: failed to get accounts. %v", err)
			return err
		}

		if err := cursor.All(gCtx, &accounts); err != nil {
			log.Error().Err(err).Msgf("mongo: failed to decode accounts. %v", err)
			return err
		}

		return nil
	})

	var totalCount int64
	errGroup.Go(func() error {
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			log.Error().Err(err).Msgf("mongo: failed to count accounts. %v", err)
			return err
		}

		totalCount = count
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return nil, 0, err
	}

	return accounts, totalCount, nil
}

func (m *Mongo) GetAccount(ctx context.Context, id string) (*model.Account, error) {
	filter, err := newNotDeletedByIDFilter(id)
	if err != nil {
		return nil, store.ErrNotFound
	}

	return m.getAccount(ctx, filter)
}

func (m *Mongo) CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	currentTime := time.Now()
	account.Dates = model.Dates{
		Created:  &currentTime,
		Modified: &currentTime,
	}

	result, err := m.database.
		Collection(CollectionAccounts).
		InsertOne(ctx, account)

	if err != nil {
		log.Error().Err(err).Msgf("mongo: failed to create account. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	newID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		err = fmt.Errorf("id of the inserted document %q is not an object id", result.InsertedID)
		return nil, err
	}

	account.ID = newID.Hex()
	return account, nil
}

func (m *Mongo) UpdateAccount(ctx context.Context, id string, account *model.Account) (*model.Account, error) {
	filter, err := newNotDeletedByIDFilter(id)
	if err != nil {
		return nil, store.ErrNotFound
	}

	currentTime := time.Now()
	update := bson.M{
		"$set": bson.M{
			"name":           account.Name,
			"balance":        account.Balance,
			"dates.modified": &currentTime,
		},
	}

	result, err := m.
		database.
		Collection(CollectionAccounts).
		UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, store.ErrNotFound
	}

	return m.getAccount(ctx, filter)
}

func (m *Mongo) DeleteAccount(ctx context.Context, id string) error {
	filter, err := newNotDeletedByIDFilter(id)
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
		Collection(CollectionAccounts).
		UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if updateResult.MatchedCount == 0 {
		return store.ErrNotFound
	}

	return nil
}

func (m *Mongo) getAccount(ctx context.Context, filter bson.M) (*model.Account, error) {
	res := m.database.
		Collection(CollectionAccounts).
		FindOne(ctx, filter)

	account := &model.Account{}
	err := res.Decode(account)

	if err != nil {
		log.Error().Err(err).Msgf("mongo getAccount: error while decoding the database object to a account. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	return account, nil
}
