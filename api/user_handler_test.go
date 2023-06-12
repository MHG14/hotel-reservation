package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mhg14/hotel-reservation/db"
	"github.com/mhg14/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const testDbUri = "mongodb://localhost:27017/hotel-reservation"
const testDbName = "hotel-reservation-test"

type testDb struct {
	db.UserStore
}

func (tDb *testDb) tearDown(t *testing.T) {
	if err := tDb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testDb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testDbUri))
	if err != nil {
		log.Fatal(err)
	}

	return &testDb{
		UserStore: db.NewMongoUserStore(client, testDbName),
	}
}

func TestPostUser(t *testing.T) {
	tDb := setup(t)
	defer tDb.tearDown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tDb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "MOhamad",
		LastName:  "sedghi",
		Password:  "gafgrqcecr",
		Email:     "dinamohamadi@gmail.com",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))

	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Log(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

	if len(user.ID) == 0 {
		t.Errorf("Expecting a user ID to be set")
	}

	if len(user.EncryptedPassword) > 0 {
		t.Errorf("Expecting the user encrypted password not to be included json response")
	}

	if user.FirstName != params.FirstName {
		t.Errorf("Expected first name %s but got %s", params.FirstName, user.FirstName)
	}

	if user.LastName != params.LastName {
		t.Errorf("Expected last name %s but got %s", params.LastName, user.LastName)
	}

	if user.Email != params.Email {
		t.Errorf("Expected email name %s but got %s", params.Email, user.Email)
	}
}
