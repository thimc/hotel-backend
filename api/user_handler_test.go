package api

import (
	"bytes"
	"encoding/json"
	"hotel/types"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestPostUser(t *testing.T) {
	// spin up a temporary database
	tdb := setup(t)
	defer tdb.teardown(t)

	// start the fiber app
	app := fiber.New()
	userHandler := NewUserHandler(tdb.User)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "some@foo.com",
		FirstName: "James",
		LastName:  "Foo",
		Password:  "12345678910",
	}
	b, err := json.Marshal(&params)
	if err != nil {
		t.Error(err)
	}
	// make request
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	// decode response
	var user types.User
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Error(err)
	}

	// validate
	if len(user.ID) == 0 {
		t.Errorf("expected a user id but got %d\n", len(user.ID))
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expected the encrypted password not to be included but got %d bytes\n",
			len(user.EncryptedPassword))
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected first name %s but got %s\n", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected last name %s but got %s\n", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s\n", params.Email, user.Email)
	}
}
