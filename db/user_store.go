package db

import (
	"context"

	"github.com/mhg14/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "user"

type UserStore interface {
	GetUserById(context.Context, string) (*types.User, error)
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(dbName).Collection(userColl),
	}
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func (m *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	var user types.User
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if err := m.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
