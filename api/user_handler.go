package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/mhg14/hotel-reservation/db"
	"github.com/mhg14/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	user, err := h.userStore.GetUserById(c.Context(), oid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "Not Found"})
		}
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		// values bson.M
		params types.UpdateUserParams
		userId = c.Params("id")
	)
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	filter := bson.M{
		"_id": oid,
	}
	if err := h.userStore.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}
	return c.JSON(map[string]string{"updated": userId})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	if err := h.userStore.DeleteUser(c.Context(), oid); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": userId})
}
