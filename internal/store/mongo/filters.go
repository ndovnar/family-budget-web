package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getNotDeletedByIDFilter(id string) (bson.M, error) {
	primitiveID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return bson.M{
		"_id":     primitiveID,
		"deleted": false,
	}, nil
}

func getByEmailFilter(email string) bson.M {
	return bson.M{
		"email":   email,
		"deleted": false,
	}
}

func getByIDFilter(id string) (bson.M, error) {
	primitiveID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return bson.M{
		"_id": primitiveID,
	}, nil
}
