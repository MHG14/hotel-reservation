package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mhg14/hotel-reservation/api"
	"github.com/mhg14/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName = "hotel-reservation"
	dbUri = "mongodb://localhost:27017/hotel-reservation"
)

var config = fiber.Config(fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
})

func main() {
	listenPort := flag.String("listenPort", ":5000", "The listen port of the api server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}

	// handlers initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client, dbName))

	app := fiber.New(config)
	apiv1 := app.Group("api/v1")

	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Post("/users", userHandler.HandlePostUser)
	apiv1.Get("/users/:id", userHandler.HandleGetUser)
	apiv1.Put("/users/:id", userHandler.HandlePutUser)
	apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)



	app.Listen(*listenPort)
}
