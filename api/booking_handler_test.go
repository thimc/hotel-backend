package api

import (
	"encoding/json"
	"fmt"
	"hotel/db/fixtures"
	"hotel/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := fixtures.AddUser(&tdb.Store, "James", "Solo", true)
	hotel := fixtures.AddHotel(&tdb.Store, "Hotel", "None", 5, nil)
	room := fixtures.AddRoom(&tdb.Store, "medium", 99.99, hotel.ID)
	booking := fixtures.AddBooking(&tdb.Store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	_ = booking

	app := fiber.New()
	admin := app.Group("/", JWTAuthentication(tdb.User), AdminAuth)
	bookingHandler := NewBookingHandler(&tdb.Store)
	admin.Get("/", bookingHandler.HandleGetBookings)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add(TokenHeader, CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	have := bookings[0]

	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking, got %d", len(bookings))
	}
	if have.ID != booking.ID {
		t.Fatalf("expected booking id %s, got %s", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("expected user id %s, got %s", booking.ID, have.ID)
	}
}

func TestAdminGetBookingNormalUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := fixtures.AddUser(&tdb.Store, "James", "Solo", false)
	hotel := fixtures.AddHotel(&tdb.Store, "Hotel", "None", 5, nil)
	room := fixtures.AddRoom(&tdb.Store, "medium", 99.99, hotel.ID)
	_ = fixtures.AddBooking(&tdb.Store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))

	app := fiber.New()
	admin := app.Group("/", JWTAuthentication(tdb.User), AdminAuth)
	bookingHandler := NewBookingHandler(&tdb.Store)
	admin.Get("/", bookingHandler.HandleGetBookings)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add(TokenHeader, CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected http status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestUserGetBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := fixtures.AddUser(&tdb.Store, "James", "Solo", false)
	hotel := fixtures.AddHotel(&tdb.Store, "Hotel", "None", 5, nil)
	room := fixtures.AddRoom(&tdb.Store, "medium", 99.99, hotel.ID)
	booking := fixtures.AddBooking(&tdb.Store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))

	app := fiber.New()
	route := app.Group("/", JWTAuthentication(tdb.User))
	bookingHandler := NewBookingHandler(&tdb.Store)
	route.Get("/:id", bookingHandler.HandleGetBooking)

	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add(TokenHeader, CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var have types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&have); err != nil {
		t.Fatal(err)
	}

	if have.ID != booking.ID {
		t.Fatalf("expected booking id %s, got %s", booking.ID, have.ID)
	}

	if have.UserID != booking.UserID {
		t.Fatalf("expected user id %s, got %s", booking.UserID, have.UserID)
	}
}
