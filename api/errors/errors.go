package errors

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

func (e Error) Error() string {
	return e.Message
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func ErrorInvalidID() Error {
	return Error{
		Code:    http.StatusBadRequest,
		Message: "Invalid id",
	}
}

func ErrorUnauthorized() Error {
	return Error{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized",
	}
}

func ErrorBadRequest() Error {
	return Error{
		Code:    http.StatusBadRequest,
		Message: "Bad request",
	}
}

func ErrorNotFound(resource string) Error {
	return Error{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("Resource %s not found", resource),
	}
}
