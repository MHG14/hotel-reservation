package api

import (
	"context"
	"log"
	"testing"

	"github.com/mhg14/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testDbUri  = "mongodb://localhost:27017/hotel-reservation"
	testDbName = "hotel-reservation-test"
)

type testDb struct {
	client *mongo.Client
	*db.Store
}

func (tDb *testDb) tearDown(t *testing.T) {
	if err := tDb.client.Database(testDbName).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testDb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testDbUri))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)

	return &testDb{
		client: client,
		Store: &db.Store{
			User:    db.NewMongoUserStore(client),
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Hotel:   hotelStore,
			Booking: db.NewMongoBookingStore(client),
		},
	}
}
