package api

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/thimc/hotel-backend/db"
	"github.com/thimc/hotel-backend/types"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

/* AuthParams are used when the user is trying to authenticate, these fields are mandatory. */
type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

/*
HandleAuthenticate will do the following:
- Parse the body and expect it to be a JSON formated AuthParams
- Get the user information based on the passed email
- Verify that their password is correct
- Return a authentication response with a token
*/
func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return ErrorBadRequest()
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrorUnauthorized()
		}
		return ErrorNotFound("User")
	}

	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return ErrorUnauthorized()
	}

	resp := AuthResponse{
		User:  user,
		Token: CreateTokenFromUser(user),
	}
	return c.JSON(resp)
}

/*
CreateTokenFromUser does the following:
- Create a token that is valid for 4 hours
- Sign it with the environment variable JWT_SECRET
*/
func CreateTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()

	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret:", err)
	}
	return tokenStr
}
