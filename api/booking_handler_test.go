package api

import (
	"testing"
	"time"

	"github.com/mhg14/hotel-reservation/db/fixtures"
)

func TestGetBookings(t *testing.T) {
	db := setup(t)
	defer db.tearDown(t)
	user := fixtures.AddUser(db.Store, "james", "foo", false)
	hotel := fixtures.AddHotel(db.Store, "some hotel", "tehran", 4.0, nil)
	room := fixtures.AddRoom(db.Store, "small", true, 44, hotel.ID)

	from := time.Now()
	till := time.Now().AddDate(0, 0, 4)

	booking := fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
	_ = booking

	// fmt.Println(booking)
}
