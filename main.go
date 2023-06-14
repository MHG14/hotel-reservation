package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mhg14/hotel-reservation/api"
	"github.com/mhg14/hotel-reservation/api/middlewares"
	"github.com/mhg14/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config(fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
})

func main() {
	listenPort := flag.String("listenPort", ":5000", "The listen port of the api server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBUri))
	if err != nil {
		log.Fatal(err)
	}

	var (
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)
		userStore  = db.NewMongoUserStore(client)

		store = &db.Store{
			Room:  roomStore,
			Hotel: hotelStore,
			User:  userStore,
		}
		userHandler  = api.NewUserHandler(userStore)
		authHandler  = api.NewAuthHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		app          = fiber.New(config)
		auth         = app.Group("api")
		apiv1        = app.Group("api/v1", middlewares.JWTAuthentication)
	)

	// auth handlres
	auth.Post("/auth", authHandler.HandleAuthnticate)

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

	app.Listen(*listenPort)
}
