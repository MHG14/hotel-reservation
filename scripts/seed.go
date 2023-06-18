package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mhg14/hotel-reservation/api"
	"github.com/mhg14/hotel-reservation/db"
	"github.com/mhg14/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBUri))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)

	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   hotelStore,
	}

	user := fixtures.AddUser(store, "james", "foo", false)
	fmt.Println("james ->", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("admin ->", api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(store, "Al Khalifa", "Dubai", 5.0, nil)
	fmt.Println("hotelID ->", hotel.ID)
	room := fixtures.AddRoom(store, "large", true, 98.44, hotel.ID)
	fmt.Println("roomID ->", room.ID)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 4))
	fmt.Println("bookingID ->", booking.ID)
}
