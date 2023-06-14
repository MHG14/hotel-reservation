package db

import (
	"context"
	"fmt"

	"github.com/mhg14/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "user"

type Dropper interface {
	Drop(context.Context) error
}

type UserStore interface {
	Dropper

	GetUserByEmail(context.Context, string) (*types.User, error)
	GetUserById(context.Context, primitive.ObjectID) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	UpdateUser(context.Context, bson.M, types.UpdateUserParams) error
	DeleteUser(context.Context, primitive.ObjectID) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(DBName).Collection(userColl),
	}
}

func (m *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := m.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*types.User

	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (m *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	if err := m.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *MongoUserStore) GetUserById(ctx context.Context, id primitive.ObjectID) (*types.User, error) {
	var user types.User
	if err := m.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := m.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil

}

func (m *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error {
	update := bson.M{
		"$set": params.ToBSON(),
	}
	_, err := m.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoUserStore) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	_, err := m.coll.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		return err
	}
	return nil
}

func (m *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("--- DROPPING USER COLLECTION")
	return m.coll.Drop(ctx)
}
