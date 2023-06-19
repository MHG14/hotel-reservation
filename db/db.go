package db

const (
	DBUri      = "mongodb://localhost:27017/hotel-reservation"
	DBName     = "hotel-reservation"
	TestDBName = "hotel-reservation-test"
)

type Pagination struct {
	Limit int64
	Page  int64
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
