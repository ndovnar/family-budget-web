package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getNotDeletedByIDFilter(id primitive.ObjectID) bson.M {
	return bson.M{
		"_id":     id,
		"deleted": false,
	}
}

func getByEmailFilter(email string) bson.M {
	return bson.M{
		"email": email,
	}
}

func getByIDFilter(id primitive.ObjectID) bson.M {
	return bson.M{
		"_id": id,
	}
}
