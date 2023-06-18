package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mhg14/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return ErrResourceNotFound("bookings")
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	booking, err := h.store.Booking.GetBookingByID(c.Context(), oid)
	if err != nil {
		return ErrResourceNotFound("booking")
	}

	user, err := GetAuthenticatedUser(c)
	if err != nil {
		return ErrUnauthorized()
	}
	if user.ID != booking.UserID {
		return ErrUnauthorized()
	}

	return c.JSON(booking)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	oid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	booking, err := h.store.Booking.GetBookingByID(c.Context(), oid)
	if err != nil {
		return err
	}

	user, err := GetAuthenticatedUser(c)
	if err != nil {
		return ErrUnauthorized()
	}
	if user.ID != booking.UserID {
		return ErrUnauthorized()
	}
	if err := h.store.Booking.UpdateBooking(c.Context(), booking.ID, bson.M{"canceled": true}); err != nil {
		return err
	}
	return c.JSON(genericResp{
		Type: "msg",
		Msg:  "updated",
	})

}
