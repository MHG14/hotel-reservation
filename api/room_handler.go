package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mhg14/hotel-reservation/db"
	"github.com/mhg14/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	store *db.Store
}

type BookRoomParams struct {
	From     time.Time `json:"from"`
	Till     time.Time `json:"till"`
	Capacity int       `json:"capacity"`
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (b *BookRoomParams) validate() error {
	now := time.Now()
	if now.After(b.From) || now.After(b.Till) {
		return fmt.Errorf("you can not book in this range of time")
	}
	return nil
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.validate(); err != nil {
		return err
	}
	roomID := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  "internal server error",
		})
	}

	ok, err = h.isRoomValidForBooking(c.Context(), oid, params)
	if err != nil {
		return err
	}

	if !ok {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  fmt.Sprintf("room %s is already booked", roomID),
		})
	}

	booking := types.Booking{
		RoomID:   oid,
		UserID:   user.ID,
		From:     params.From,
		Till:     params.Till,
		Capacity: params.Capacity,
	}

	insertedBooking, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}

	return c.JSON(insertedBooking)
}

func (h *RoomHandler) isRoomValidForBooking(c context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"roomID": roomID,
		"from": bson.M{
			"$gte": params.From,
		},
		"till": bson.M{
			"$lte": params.Till,
		},
	}
	bookings, err := h.store.Booking.GetBookings(c, where)
	if err != nil {
		return false, err
	}
	fmt.Println(len(bookings))

	ok := len(bookings) == 0
	return ok, nil
}
