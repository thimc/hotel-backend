package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}
	apiError := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiError.Code).JSON(apiError)
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Error struct {
	Response
	Code int `json:"error"`
}

func (e Error) Error() string {
	return e.Message
}

func NewError(code int, message string) *Error {
	return &Error{
		Response: Response{
			Success: false,
			Message: message,
		},
		Code: code,
	}
}

func ErrorInvalidID() Error {
	return Error{
		Response: Response{
			Success: false,
			Message: "Invalid id",
		},
		Code: http.StatusBadRequest,
	}
}

func ErrorUnauthorized() Error {
	return Error{
		Response: Response{
			Success: false,
			Message: "Unauthorized",
		},
		Code: http.StatusUnauthorized,
	}
}

func ErrorTokenExpired() Error {
	return Error{
		Response: Response{
			Success: false,
			Message: "Token expired",
		},
		Code: http.StatusUnauthorized,
	}
}

func ErrorBadRequest() Error {
	return Error{
		Response: Response{
			Success: false,
			Message: "Bad request",
		},
		Code: http.StatusBadRequest,
	}
}

func ErrorNotFound(resource string) Error {
	return Error{
		Response: Response{
			Success: false,
			Message: fmt.Sprintf("%s not found", resource),
		},
		Code: http.StatusNotFound,
	}
}
