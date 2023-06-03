package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
		t.Fatalf("expected the user to be the inserted user")
	}

}

func TestAuthenticationWithWrongPasswordFailure(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	fixtures.AddUser(&tdb.Store, "James", "Solo", false)

	// start the fiber app
	app := fiber.New()
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

	var genericResp GenericResp
	if err := json.NewDecoder(resp.Body).Decode(&genericResp); err != nil {
		t.Fatal(err)
	}

	if genericResp.Success {
		t.Fatalf("expected success false but got %v", genericResp.Success)
	}

	if genericResp.Msg != "invalid credentials" {
		t.Fatalf("expected response  %s but got %s", "invalid credentials", genericResp.Msg)
	}
}
