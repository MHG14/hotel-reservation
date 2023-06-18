package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mhg14/hotel-reservation/db/fixtures"
	"github.com/mhg14/hotel-reservation/types"
)

func TestUserGetBooking(t *testing.T) {
	db := setup(t)
	defer db.tearDown(t)

	var (
		notAuthenticatedUser = fixtures.AddUser(db.Store, "mohammad", "ghamari", false)
		user                 = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel                = fixtures.AddHotel(db.Store, "some hotel", "tehran", 4.0, nil)
		room                 = fixtures.AddRoom(db.Store, "small", true, 44, hotel.ID)
		from                 = time.Now()
		till                 = time.Now().AddDate(0, 0, 4)
		booking              = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app                  = fiber.New()
		route                = app.Group("/", JWTAuthentication(db.User))
		bookingHandler       = NewBookingHandler(db.Store)
	)

	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected staus code 200 but got %d", resp.StatusCode)
	}
	var bookingResponse *types.Booking

	if err := json.NewDecoder(resp.Body).Decode(&bookingResponse); err != nil {
		t.Fatal(err)
	}

	if bookingResponse.ID != booking.ID {
		t.Fatalf("expected %s got %s\n", booking.ID, bookingResponse.ID)
	}

	if bookingResponse.UserID != booking.UserID {
		t.Fatalf("expected %s got %s\n", booking.UserID, bookingResponse.UserID)
	}

	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-token", CreateTokenFromUser(notAuthenticatedUser))

	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 status code but got %d\n", resp.StatusCode)
	}

}

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.tearDown(t)

	var (
		adminUser      = fixtures.AddUser(db.Store, "admin", "admin", true)
		user           = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel          = fixtures.AddHotel(db.Store, "some hotel", "tehran", 4.0, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 44, hotel.ID)
		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 4)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		admin          = app.Group("/", JWTAuthentication(db.User), AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)

	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-token", CreateTokenFromUser(adminUser))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200 but got %d\n", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking got %d\n", len(bookings))
	}
	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected %s got %s\n", booking.ID, have.ID)
	}

	if have.UserID != booking.UserID {
		t.Fatalf("expected %s got %s\n", booking.UserID, have.UserID)
	}
	// if !reflect.DeepEqual(booking, bookings[0]) {
	// 	t.Fatal("expected bookings to be euqal")
	// }
	fmt.Println(bookings)

	// test if non-admin can not access bookings
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-token", CreateTokenFromUser(user))

	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized but got %d\n", resp.StatusCode)
	}
	// fmt.Println(booking)
}
