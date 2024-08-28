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

func (m *Mongo) GetTransactions(ctx context.Context, getTransactionsFilter *model.GetTransactionsFilter) ([]*model.Transaction, error) {
	filter := bson.M{
		"$or": []any{
			bson.M{
				"fromAccount": getTransactionsFilter.Account,
			},
			bson.M{
				"toAccount": getTransactionsFilter.Account,
			},
		},
		"deleted": getTransactionsFilter.Deleted,
	}

	res, err := m.database.
		Collection(CollectionTransactions).
		Find(ctx, filter, nil)

	if err != nil {
		log.Error().Err(err).Msgf("mongo: failed to get transactions. %v", err)
		return nil, err
	}

	transactions := []*model.Transaction{}
	for res.Next(ctx) {
		transaction := &model.Transaction{}
		err = res.Decode(transaction)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (m *Mongo) GetTransaction(ctx context.Context, id string) (*model.Transaction, error) {
	filter, err := getNotDeletedByIDFilter(id)
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
	filter, err := getNotDeletedByIDFilter(id)
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
			"fromAccount":    updatedTransaction.FromAccount,
			"toAccount":      updatedTransaction.ToAccount,
			"category":       updatedTransaction.Category,
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
	filter, err := getNotDeletedByIDFilter(id)
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
		fromAccount:        transaction.FromAccount,
		fromAccountAmmount: transaction.Amount * -1,
		toAccount:          transaction.ToAccount,
		toAccountAmmount:   transaction.Amount,
		category:           transaction.Category,
		categoryAmmount:    transaction.Amount * -1,
	})
}

func (m *Mongo) revertTransaction(ctx context.Context, transaction *model.Transaction) error {
	return m.makeTransaction(ctx, &makeTransactionParams{
		fromAccount:        transaction.FromAccount,
		fromAccountAmmount: transaction.Amount,
		toAccount:          transaction.ToAccount,
		toAccountAmmount:   transaction.Amount * -1,
		category:           transaction.Category,
		categoryAmmount:    transaction.Amount,
	})
}

type makeTransactionParams struct {
	fromAccount        string
	fromAccountAmmount float64
	toAccount          string
	toAccountAmmount   float64
	category           string
	categoryAmmount    float64
}

func (m *Mongo) makeTransaction(ctx context.Context, params *makeTransactionParams) error {
	var fromAccount, toAccount *model.Account
	var category *model.Category

	if params.fromAccount != "" {
		accountDB, err := m.GetAccount(ctx, params.fromAccount)
		if err != nil {
			return err
		}

		fromAccount = accountDB
	}

	if params.toAccount != "" {
		accountDB, err := m.GetAccount(ctx, params.toAccount)
		if err != nil {
			return err
		}

		toAccount = accountDB
	}

	if params.category != "" {
		categoryDB, err := m.GetCategory(ctx, params.category)
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
