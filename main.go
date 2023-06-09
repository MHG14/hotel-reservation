package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/mhg14/hotel-reservation/api"
	"github.com/mhg14/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Configuration
// 1. MongoDB endpoint
// 2. ListenAddress of http server
// 3. JWT secret
// 4. MongoDB name

var config = fiber.Config(fiber.Config{
	ErrorHandler: api.ErrorHandler,
})

func main() {
	mongoURI := os.Getenv("MONGO_DB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)

		store = &db.Store{
			Room:    roomStore,
			Hotel:   hotelStore,
			User:    userStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(userStore)
		authHandler    = api.NewAuthHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app            = fiber.New(config)
		auth           = app.Group("api")
		apiv1          = app.Group("api/v1", api.JWTAuthentication(userStore))
		admin          = apiv1.Group("admin", api.AdminAuth)
	)

	// auth handlres
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// Versioned API routes
	// user handlers
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Post("/users", userHandler.HandlePostUser)
	apiv1.Get("/users/:id", userHandler.HandleGetUser)
	apiv1.Put("/users/:id", userHandler.HandlePutUser)
	apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)

	// hotel handlers
	apiv1.Get("hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("hotels/:id/rooms", hotelHandler.HandleGetRooms)
	apiv1.Get("hotels/:id", hotelHandler.HandleGetHotelById)

	// room handlers
	apiv1.Post("rooms/:id/book", roomHandler.HandleBookRoom)
	apiv1.Get("/rooms", roomHandler.HandleGetRooms)

	// booking handlers
	apiv1.Get("/bookings/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/bookings/:id/cancel", bookingHandler.HandleCancelBooking)

	// admin handlers
	admin.Get("/bookings", bookingHandler.HandleGetBookings)

	listenPort := os.Getenv("HTTP_LISTEN_PORT")
	app.Listen(listenPort)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
