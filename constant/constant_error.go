package constant

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

)

// error.
var (
	ErrFailedParseRequest = fiber.NewError(http.StatusBadRequest, "failed to parse request")

	ErrHeaderTokenNotFound = fiber.NewError(http.StatusUnauthorized, "header authorization not found")
	ErrHeaderTokenInvalid  = fiber.NewError(http.StatusUnauthorized, "invalid header token")
	ErrTokenInvalid        = fiber.NewError(http.StatusUnauthorized, "invalid token")
	ErrTokenMissing        = fiber.NewError(http.StatusUnauthorized, "missing token")
	ErrTokenExpired        = fiber.NewError(http.StatusUnauthorized, "expired token")
	ErrTokenUnauthorized   = fiber.NewError(http.StatusUnauthorized, "unauthorized token")

	ErrDataNotFound = fiber.NewError(http.StatusNotFound, "data not found")

	ErrUnknownSource = fiber.NewError(http.StatusInternalServerError, "an error occurred, please try again later")
)

// error message.
const (
	ErrMsgValidate      = "There are some errors in your request"
	ErrMsgUnknownSource = "an error occurred, please try again later"
)

// error form field.
var (
	// 400.
	ErrPasswordIncorrect  = errors.New("password incorrect")
	ErrAccountNotFound    = errors.New("account not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)
