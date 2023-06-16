package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mhg14/hotel-reservation/api"
	"github.com/mhg14/hotel-reservation/db"
	"github.com/mhg14/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	hotelStore   db.HotelStore
	roomStore    db.RoomStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	client       *mongo.Client
	ctx          = context.Background()
)

func seedHotel(name string, location string, rating float64) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	_ = insertedHotel
	return insertedHotel
}

func seedUser(isAdmin bool, fname, lname, email, password string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  password,
	})

	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin
	insertedUser, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("email: %s, token: %s\n", user.Email, api.CreateTokenFromUser(user))
	return insertedUser
}

func seedRoom(size string, seaside bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Price:   price,
		Seaside: seaside,
		HotelID: hotelID,
	}

	insertedRoom, err := roomStore.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func seedBooking(roomID, userID primitive.ObjectID, from, till time.Time) {
	booking := &types.Booking{
		UserID: userID,
		RoomID: roomID,
		From:   from,
		Till:   till,
	}

	resp, err := bookingStore.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("booking", resp.ID)
}

func main() {
	james := seedUser(false, "Jamess", "Foo", "james@foo.com", "supersecurepassword")
	seedUser(true, "MohamadHasan", "Ghamari", "mhg14@foo.com", "adminpassword")
	seedHotel("Bellucia", "France", 3.0)
	seedHotel("Cozy Hotel", "Netherlands", 4.0)
	hotel := seedHotel("Al Khalifa", "Dubai", 1.0)
	seedRoom("small", true, 98.99, hotel.ID)
	seedRoom("normal", false, 122.99, hotel.ID)
	room := seedRoom("kingsize", false, 200.99, hotel.ID)
	seedBooking(room.ID, james.ID, time.Now(), time.Now().AddDate(0, 0, 2))
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
	bookingStore = db.NewMongoBookingStore(client)
}
