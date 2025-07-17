package server

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/constant"
)

func errorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			code := fiberErr.Code
			message := strings.ToUpper(fiberErr.Message[:1]) + fiberErr.Message[1:]
			return c.Status(code).JSON(fiber.Map{"code": code, "message": message})
		}

		if errors.Is(err, sql.ErrNoRows) {
			err = constant.ErrDataNotFound
		}

		code := mapErrorCode(err)
		message := mapErrorMessage(err, code)
		if len(message) > 0 {
			message = strings.ToUpper(message[:1]) + message[1:]
		}

		return c.Status(code).JSON(fiber.Map{"code": code, "message": message})
	}
}

func mapErrorCode(err error) int {
	switch {
	case errors.Is(err, constant.ErrFailedParseRequest):
		return http.StatusBadRequest
	case errors.Is(err, constant.ErrHeaderTokenNotFound),
		errors.Is(err, constant.ErrHeaderTokenInvalid),
		errors.Is(err, constant.ErrTokenUnauthorized),
		errors.Is(err, constant.ErrTokenInvalid),
		errors.Is(err, constant.ErrTokenExpired):
		return http.StatusUnauthorized
	case errors.Is(err, constant.ErrForbiddenRole),
		errors.Is(err, constant.ErrForbiddenPermission):
		return http.StatusForbidden
	case errors.Is(err, constant.ErrDataNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func mapErrorMessage(err error, code int) string {
	appDebug, _ := strconv.ParseBool(os.Getenv("APP_DEBUG"))
	message := err.Error()

	if !appDebug && code == http.StatusInternalServerError {
		message = constant.ErrUnknownSource.Error()
	}

	return message
}
