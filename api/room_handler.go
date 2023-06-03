package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/thimc/hotel-backend/db"
	"github.com/thimc/hotel-backend/types"

	"github.com/gofiber/fiber/v2"
)

/*
The User ID and Room ID isn't specified here because the User ID is already in the Context
*/
type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	UntilDate  time.Time `json:"untilDate"`
	NumPersons int       `json:"numPersons"`
}

/*
validate will:
- Check that the from and until date are valid
- Check that the user has specified at least one customer
*/
func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.UntilDate) {
		return ErrorBadRequest()
	}
	if p.NumPersons < 1 {
		return fmt.Errorf("")
	}
	return nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

/*
HandleGetRooms will return all the rooms with no filter
*/
func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), map[string]any{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

/*
HandleGetRoom will return th rooms that matches id
*/
func (h *RoomHandler) HandleGetRoom(c *fiber.Ctx) error {
	id := c.Params("id")
	room, err := h.store.Room.GetRoomByID(c.Context(), id)
	if err != nil {
		return ErrorNotFound("Room")
	}
	return c.JSON(room)
}

/*
HandleBookRoom will do the following:
- Expect that the body is a JSON marshaled `BookRoomParams`
- Validate the parameters
- Convert the room ID from the HTTP URL to a primitive MongoDB id
- Get the user data from the fiber context (this is being passed from the middleware)
- Check that the room is available
- Generate and insert a booking
*/
func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return ErrorBadRequest()
	}

	if err := params.validate(); err != nil {
		return ErrorBadRequest()
	}

	roomID := c.Params("id")

	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.JSON(Response{
			Success: false,
			Message: "internal server error",
		})
	}

	ok, err := h.isRoomAvailable(c, roomID, params)
	if err != nil {
		return err
	}
	if !ok {
		return NewError(http.StatusBadRequest, fmt.Sprintf("room %s is already booked", roomID))
	}

	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomID,
		FromDate:   params.FromDate,
		UntilDate:  params.UntilDate,
		NumPersons: params.NumPersons,
	}

	insertedBooking, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}

	return c.JSON(insertedBooking)
}

/* isRoomAvailable will validate that the room is available */
func (h *RoomHandler) isRoomAvailable(c *fiber.Ctx, roomID string, params BookRoomParams) (bool, error) {
	filter := map[string]any{
		"roomID": roomID,
		"fromDate": map[string]any{
			"$gte": params.FromDate,
		},
		"untilDate": map[string]any{
			"$lte": params.UntilDate,
		},
	}
	bookings, err := h.store.Booking.GetBookings(c.Context(), filter)
	if err != nil {
		return false, err
	}
	return len(bookings) == 0, nil
}
