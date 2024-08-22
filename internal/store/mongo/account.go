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

func (m *Mongo) GetAccounts(ctx context.Context, filter *store.GetAccountsFilter) ([]*model.Account, error) {
	res, err := m.database.
		Collection(CollectionAccounts).
		Find(ctx, bson.M{"owner": filter.Owner}, nil)

	if err != nil {
		log.Error().Err(err).Msgf("mongo: failed to get accounts. %v", err)
		return nil, err
	}

	dtoAccounts := []*dto.Account{}
	for res.Next(ctx) {
		dtoAccount := &dto.Account{}
		err = res.Decode(&dtoAccount)
		if err != nil {
			return nil, err
		}

		dtoAccounts = append(dtoAccounts, dtoAccount)
	}

	return dto.DtoAccountsToModelAccounts(dtoAccounts), nil
}

func (m *Mongo) GetAccount(ctx context.Context, id string) (*model.Account, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, store.ErrNotFound
	}

	filter := getNotDeletedByIDFilter(oid)
	return m.getAccount(ctx, filter)
}

func (m *Mongo) CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	currentTime := time.Now()

	dtoAccount := dto.ModelAccountToDtoAccount(account)
	dtoAccount.Dates = dto.Dates{Created: &currentTime, Modified: &currentTime}

	_, err := m.database.
		Collection(CollectionAccounts).
		InsertOne(ctx, dtoAccount)

	if err != nil {
		log.Error().Err(err).Msgf("mongo: failed to create account. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	createdAccount := dto.DtoAccountToModelAccount(dtoAccount)

	return createdAccount, nil
}

func (m *Mongo) UpdateAccount(ctx context.Context, id string, account *model.Account) (*model.Account, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, store.ErrNotFound
	}

	filter := getNotDeletedByIDFilter(oid)
	currentTime := time.Now()
	dtoAccount := dto.ModelAccountToDtoAccount(account)
	dtoAccount.Dates.Modified = &currentTime
	update := bson.M{
		"$set": bson.M{
			"name":           dtoAccount.Name,
			"dates.modified": dtoAccount.Dates.Modified,
			"balance":        dtoAccount.Balance,
		},
	}

	updateResult, err := m.
		database.
		Collection(CollectionAccounts).
		UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if updateResult.MatchedCount == 0 {
		return nil, store.ErrNotFound
	}

	return m.getAccount(ctx, filter)
}

func (m *Mongo) DeleteAccount(ctx context.Context, id string) error {
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

	account := dto.Account{}
	err := res.Decode(&account)

	if err != nil {
		log.Error().Err(err).Msgf("mongo getAccount: error while decoding the database object to a account. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	return dto.DtoAccountToModelAccount(&account), nil
}
