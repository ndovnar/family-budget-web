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

	err := m.proceedTransaction(ctx, transaction)
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

	transactionRevert := &model.Transaction{
		Type:        transaction.Type,
		Category:    transaction.Category,
		FromAccount: transaction.ToAccount,
		ToAccount:   transaction.FromAccount,
		Amount:      transaction.Amount,
	}

	err = m.proceedTransaction(ctx, transactionRevert)
	if err != nil {
		return nil, err
	}

	err = m.proceedTransaction(ctx, updatedTransaction)
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

	transactionRevert := &model.Transaction{
		Type:        transaction.Type,
		Category:    transaction.Category,
		FromAccount: transaction.ToAccount,
		ToAccount:   transaction.FromAccount,
		Amount:      transaction.Amount,
	}

	err = m.proceedTransaction(ctx, transactionRevert)
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
		fmt.Println("HERE 4")
		return err
	}

	if updateResult.MatchedCount == 0 {
		fmt.Println("HERE 5")
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

func (m *Mongo) proceedTransaction(ctx context.Context, transaction *model.Transaction) error {
	switch transaction.Type {
	case model.TransactionTypeTransfer:
		fromAccount, err := m.GetAccount(ctx, transaction.FromAccount)
		if err != nil {
			return err
		}

		toAccount, err := m.GetAccount(ctx, transaction.ToAccount)
		if err != nil {
			return err
		}

		fromAccount.Balance -= transaction.Amount
		toAccount.Balance += transaction.Amount

		_, err = m.UpdateAccount(ctx, fromAccount.ID, fromAccount)
		if err != nil {
			return err
		}

		_, err = m.UpdateAccount(ctx, toAccount.ID, toAccount)
		if err != nil {
			return err
		}
	}

	return nil
}
