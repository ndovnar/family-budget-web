package mongo

import (
	"github.com/ndovnar/family-budget-api/internal/filter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newNotDeletedByIDFilter(id string) (bson.M, error) {
	primitiveID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return bson.M{
		"_id":     primitiveID,
		"deleted": false,
	}, nil
}

func newByEmailFilter(email string) bson.M {
	return bson.M{
		"email":   email,
		"deleted": false,
	}
}

func newByIDFilter(id string) (bson.M, error) {
	primitiveID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return bson.M{
		"_id": primitiveID,
	}, nil
}

func newPaginationFindOptions(pagination *filter.Pagination) *options.FindOptions {
	return &options.FindOptions{Limit: &pagination.Limit, Skip: &pagination.Offset}
}
