package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mhg14/hotel-reservation/db"
	"github.com/mhg14/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddHotel(store *db.Store, name, location string, rating float64, rooms []primitive.ObjectID) *types.Hotel {
	var roomIDs = rooms
	if rooms == nil {
		roomIDs = []primitive.ObjectID{}
	}
	hotel := &types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    roomIDs,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.Insert(context.TODO(), hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddUser(store *db.Store, fname, lname string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", fname, lname),
		FirstName: fname,
		LastName:  lname,
		Password:  fmt.Sprintf("%s_%s", fname, lname),
	})

	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = admin
	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func AddRoom(store *db.Store, size string, seaside bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Price:   price,
		Seaside: seaside,
		HotelID: hotelID,
	}

	insertedRoom, err := store.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddBooking(store *db.Store, userID, roomID primitive.ObjectID, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID: userID,
		RoomID: roomID,
		From:   from,
		Till:   till,
	}

	insertedBooking, err := store.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}
