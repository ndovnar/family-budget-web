package mongo

import (
	"github.com/ndovnar/family-budget-api/internal/store"
	"go.mongodb.org/mongo-driver/mongo"
)

func mongoErrorToDBError(err error) error {
	if err == mongo.ErrNoDocuments {
		return store.ErrNotFound
	}

	if mongo.IsDuplicateKeyError(err) {
		return store.ErrDuplicateKey
	}

	return store.ErrUnknown
}
