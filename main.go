package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/mhg14/hotel-reservation/api"
)

func main() {
	listenPort := flag.String("listenPort", ":5000", "The listen port of the api server")
	flag.Parse()
	app := fiber.New()
	apiv1 := app.Group("api/v1")
	apiv1.Get("/users", api.HandleGetUsers)
	apiv1.Get("/users/:id", api.HandleGetUser)
	app.Listen(*listenPort)
}
