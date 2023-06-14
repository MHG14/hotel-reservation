package main

import (
	"context"
	"log"

	"github.com/mhg14/hotel-reservation/db"
	"github.com/mhg14/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	hotelStore db.HotelStore
	roomStore  db.RoomStore
	userStore  db.UserStore
	client     *mongo.Client
	ctx        = context.Background()
)

func seedHotel(name string, location string, rating float64) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:  "small",
			Price: 99.9,
		},
		{
			Size:  "normal",
			Price: 122.9,
		},
		{
			Size:  "kingsize",
			Price: 222.5,
		},
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func seedUser(fname, lname, email string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  "supersecurepassword",
	})

	if err != nil {
		log.Fatal(err)
	}
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	seedHotel("Bellucia", "France", 3.0)
	seedHotel("Cozy Hotel", "Netherlands", 4.0)
	seedHotel("Al Khalifa", "Dubai", 1.0)
	seedUser("Jamess", "Foo", "james@foo.com")
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBUri))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)

}
