package api

import (
	"net/http"
	"os"
	"time"

	"github.com/thimc/hotel-backend/db"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const TokenHeader = "X-Api-Token"

/*
JWTAuthentication will:
  - Expect that the HTTP request has a header of name `TokenHeader`
  - Get the user claims (email, expiration date and user id)
  - Verify that the token is valid
  - Verify that the token hasn't expired
  - Pass the user as types.User struct to the fiber context so
    that it may be used in the room handler.
*/
func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()[TokenHeader]
		if !ok {
			return ErrorUnauthorized()
		}

		claims, err := validateToken(token)
		if err != nil {
			return NewError(http.StatusUnauthorized, "Invalid token")
		}

		expires := int64(claims["expires"].(float64))
		if time.Now().Unix() > expires {
			return ErrorTokenExpired()
		}

		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return ErrorUnauthorized()
		}

		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrorUnauthorized()
		}
		secret := os.Getenv(db.ENV_JWT_SECRET)
		return []byte(secret), nil
	})
	if err != nil {
		return nil, ErrorUnauthorized()
	}

	if !token.Valid {
		return nil, ErrorUnauthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrorUnauthorized()
	}

	return claims, nil
}
