package api

import (
	"errors"
	"fmt"
	"github.com/thimc/hotel-backend/db"
	"github.com/thimc/hotel-backend/types"
	"net/http"
	"os"
	"time"

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

/* GenericResp is used when any type of error occurs in all of the handlers and middleware */
type GenericResp struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

/* Wrapper function for when the credentials are invalid */
func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusUnauthorized).JSON(GenericResp{
		Success: false,
		Msg:     "invalid credentials",
	})
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
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}
		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return invalidCredentials(c)
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
