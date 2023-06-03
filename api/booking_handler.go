package api

import (
	"net/http"

	"github.com/thimc/hotel-backend/db"
	"github.com/thimc/hotel-backend/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

/*
HandleGetBookings will return all bookings with no filter
*/
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return err
	}
	if booking.UserID != user.ID {
		return statusUnauthorized(c)
	}
	return c.JSON(booking)
}

/*
HandleCancelBooking will:
- get the id from the HTTP URL
- get the booking based on the id
- get the user from the context (sent via JWT middelware)
- check that the user owns it or the user is an admin
- update the canceled field
*/
func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return statusUnauthorized(c)
	}

	if booking.UserID != user.ID && !user.IsAdmin {
		return statusUnauthorized(c)
	}

	update := bson.M{"canceled": true}
	if err := h.store.Booking.UpdateBooking(c.Context(), id, update); err != nil {
		return err
	}
	return c.JSON(GenericResp{
		Success: true,
		Msg:     "ok",
	})
}

func statusUnauthorized(c *fiber.Ctx) error {
	return c.Status(http.StatusUnauthorized).JSON(GenericResp{
		Success: false,
		Msg:     "unauthorized",
	})
}
