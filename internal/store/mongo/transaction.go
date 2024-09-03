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

func (m *Mongo) GetTransactions(ctx context.Context, transactionsFilter *filter.GetTransactionsFilter) ([]*model.Transaction, int64, error) {
	errGroup, gCtx := errgroup.WithContext(ctx)
	collection := m.database.Collection(CollectionTransactions)
	filter := bson.M{
		"$or": []any{
			bson.M{
				"fromAccount": transactionsFilter.FromAccountID,
			},
			bson.M{
				"toAccount": transactionsFilter.ToAccountID,
			},
			bson.M{
				"category": transactionsFilter.CategoryID,
			},
		},
		"deleted": false,
	}
	paginationFindOptions := newPaginationFindOptions(transactionsFilter.Pagination)

	transactions := []*model.Transaction{}
	errGroup.Go(func() error {
		cursor, err := collection.Find(ctx, filter, paginationFindOptions)
		if err != nil {
			log.Error().Err(err).Msgf("mongo: failed to get transactions. %v", err)
			return err
		}

		if err := cursor.All(gCtx, &transactions); err != nil {
			log.Error().Err(err).Msgf("mongo: failed to decode transactions. %v", err)
			return err
		}

		return nil
	})

	var totalCount int64
	errGroup.Go(func() error {
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			log.Error().Err(err).Msgf("mongo: failed to count transactions. %v", err)
			return err
		}

		totalCount = count
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return nil, 0, err
	}

	return transactions, totalCount, nil
}

func (m *Mongo) GetTransaction(ctx context.Context, id string) (*model.Transaction, error) {
	filter, err := newNotDeletedByIDFilter(id)
	if err != nil {
		return nil, store.ErrNotFound
	}

	return m.getTransaction(ctx, filter)
}

func (m *Mongo) CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	currentTime := time.Now()
	transaction.Dates = model.Dates{
		Created:  &currentTime,
		Modified: &currentTime,
	}

	err := m.createTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	result, err := m.database.
		Collection(CollectionTransactions).
		InsertOne(ctx, transaction)
	if err != nil {
		log.Error().Err(err).Msgf("mongo: failed to create transaction. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	newID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		err = fmt.Errorf("id of the inserted document %q is not an object id", result.InsertedID)
		return nil, err
	}

	transaction.ID = newID.Hex()
	return transaction, nil
}

func (m *Mongo) UpdateTransaction(ctx context.Context, id string, updatedTransaction *model.Transaction) (*model.Transaction, error) {
	filter, err := newNotDeletedByIDFilter(id)
	if err != nil {
		return nil, store.ErrNotFound
	}

	transaction, err := m.getTransaction(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = m.revertTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	err = m.createTransaction(ctx, updatedTransaction)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	update := bson.M{
		"$set": bson.M{
			"type":           updatedTransaction.Type,
			"fromAccount":    updatedTransaction.FromAccountID,
			"toAccount":      updatedTransaction.ToAccountID,
			"category":       updatedTransaction.CategoryID,
			"amount":         updatedTransaction.Amount,
			"dates.modified": &currentTime,
			"description":    updatedTransaction.Description,
		},
	}

	result, err := m.database.
		Collection(CollectionTransactions).
		UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, store.ErrNotFound
	}

	return m.getTransaction(ctx, filter)
}

func (m *Mongo) DeleteTransaction(ctx context.Context, id string) error {
	filter, err := newNotDeletedByIDFilter(id)
	if err != nil {
		return store.ErrNotFound
	}

	transaction, err := m.getTransaction(ctx, filter)
	if err != nil {
		return err
	}

	err = m.revertTransaction(ctx, transaction)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"dates.deleted": time.Now(),
			"deleted":       true,
		},
	}
	updateResult, err := m.database.
		Collection(CollectionTransactions).
		UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if updateResult.MatchedCount == 0 {
		return store.ErrNotFound
	}

	return nil
}

func (m *Mongo) getTransaction(ctx context.Context, filter bson.M) (*model.Transaction, error) {
	res := m.database.
		Collection(CollectionTransactions).
		FindOne(ctx, filter)

	transaction := &model.Transaction{}
	err := res.Decode(transaction)

	if err != nil {
		log.Error().Err(err).Msgf("mongo getTransaction: error while decoding the database object to a transaction. %v", err)
		return nil, mongoErrorToDBError(err)
	}

	return transaction, nil
}

func (m *Mongo) createTransaction(ctx context.Context, transaction *model.Transaction) error {
	return m.makeTransaction(ctx, &makeTransactionParams{
		fromAccountID:      transaction.FromAccountID,
		fromAccountAmmount: transaction.Amount * -1,
		toAccountID:        transaction.ToAccountID,
		toAccountAmmount:   transaction.Amount,
		categoryID:         transaction.CategoryID,
		categoryAmmount:    transaction.Amount * -1,
	})
}

func (m *Mongo) revertTransaction(ctx context.Context, transaction *model.Transaction) error {
	return m.makeTransaction(ctx, &makeTransactionParams{
		fromAccountID:      transaction.FromAccountID,
		fromAccountAmmount: transaction.Amount,
		toAccountID:        transaction.ToAccountID,
		toAccountAmmount:   transaction.Amount * -1,
		categoryID:         transaction.CategoryID,
		categoryAmmount:    transaction.Amount,
	})
}

type makeTransactionParams struct {
	fromAccountID      string
	fromAccountAmmount float64
	toAccountID        string
	toAccountAmmount   float64
	categoryID         string
	categoryAmmount    float64
}

func (m *Mongo) makeTransaction(ctx context.Context, params *makeTransactionParams) error {
	var fromAccount, toAccount *model.Account
	var category *model.Category

	if params.fromAccountID != "" {
		accountDB, err := m.GetAccount(ctx, params.fromAccountID)
		if err != nil {
			return err
		}

		fromAccount = accountDB
	}

	if params.toAccountID != "" {
		accountDB, err := m.GetAccount(ctx, params.toAccountID)
		if err != nil {
			return err
		}

		toAccount = accountDB
	}

	if params.categoryID != "" {
		categoryDB, err := m.GetCategory(ctx, params.categoryID)
		if err != nil {
			return err
		}

		category = categoryDB
	}

	if fromAccount != nil {
		fromAccount.Balance += params.fromAccountAmmount

		_, err := m.UpdateAccount(ctx, fromAccount.ID, fromAccount)
		if err != nil {
			return err
		}
	}

	if toAccount != nil {
		toAccount.Balance += params.toAccountAmmount

		_, err := m.UpdateAccount(ctx, toAccount.ID, toAccount)
		if err != nil {
			return err
		}
	}

	if category != nil {
		category.Balance += params.categoryAmmount

		_, err := m.UpdateCategory(ctx, category.ID, category)
		if err != nil {
			return err
		}
	}

	return nil
}
