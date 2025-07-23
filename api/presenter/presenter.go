package presenter

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ResponsePayload struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func ResponseOKForErrNoRows(c *fiber.Ctx, err error, msg string) error {
	if !errors.Is(err, sql.ErrNoRows) {
		return ResponseError(c, err, msg)
	}

	return c.Status(http.StatusOK).JSON(responseDataPayload{
		Message: msg,
	})
}
