package responses

import (
	"fmt"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"error"`
}

func (r Response) Error() string {
	return r.Message
}

func New(success bool, message string) Response {
	return Response{
		Success: success,
		Message: message,
	}
}

type Error struct {
	Response
	Code int `json:"code"`
}

func (e Error) Error() string {
	return e.Response.Message
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
			Message: fmt.Sprintf("Resource %s not found", resource),
		},
		Code: http.StatusNotFound,
	}
}

func ErrorInternalServer() Error {
	return Error{
		Response: Response{
			Success: false,
			Message: "Internal service error",
		},
		Code: http.StatusInternalServerError,
	}
}
