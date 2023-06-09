package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mhg14/hotel-reservation/types"
)

func TestPostUser(t *testing.T) {
	tDb := setup(t)
	defer tDb.tearDown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tDb.User)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "Mohamad",
		LastName:  "Ghamari",
		Password:  "12345678",
		Email:     "test@gmail.com",
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
