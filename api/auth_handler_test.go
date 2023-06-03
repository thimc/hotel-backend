package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thimc/hotel-backend/api/responses"
	"github.com/thimc/hotel-backend/db/fixtures"

	"github.com/gofiber/fiber/v2"
)

func TestAuthenticationSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	insertedUser := fixtures.AddUser(&tdb.Store, "James", "Solo", false)

	// start the fiber app
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "James@Solo.com",
		Password: "James_Solo",
	}

	b, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}

	// make HTTP request
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var authResponse AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		t.Fatal(err)
	}

	if authResponse.Token == "" {
		t.Fatalf("expected the JWT token to be present in the auth response")
	}

	// set the EncryptedPassword to an empty string because we don't return
	// it any JSON response.
	insertedUser.EncryptedPassword = ""
	if insertedUser.ID != authResponse.User.ID {
		t.Fatalf("expected the user ID %s to be the inserted user ID %s", authResponse.User.ID, insertedUser.ID)
	}
	if insertedUser.FirstName != authResponse.User.FirstName {
		t.Fatalf("expected the user first name %s to be the inserted user first name %s", authResponse.User.FirstName, insertedUser.FirstName)
	}
	if insertedUser.LastName != authResponse.User.LastName {
		t.Fatalf("expected the user last name %s to be the inserted user last name %s", authResponse.User.LastName, insertedUser.LastName)
	}
}

func TestAuthenticationWithWrongPasswordFailure(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	fixtures.AddUser(&tdb.Store, "James", "Solo", false)

	// start the fiber app
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "James@Solo.com",
		Password: "unknown",
	}

	b, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}

	// make HTTP request
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected http status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Success {
		t.Fatalf("expected success false but got %v", response.Success)
	}

	if response.Message != responses.ErrorUnauthorized().Message {
		t.Fatalf("expected response  %s but got %s", responses.ErrorUnauthorized().Message, response.Message)
	}
}
