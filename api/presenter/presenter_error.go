package presenter

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/internal/validator"
)

type responseErrorPayload struct {
	Code    int32 `json:"code"`
	Error   interface{} `json:"error,omitempty"`
	Message string      `json:"message"`
}

func ResponseErrorValidate(c *fiber.Ctx, err error) error {
	return c.Status(http.StatusBadRequest).JSON(responseErrorPayload{
		Code: http.StatusBadRequest,
		Message: constant.ErrMsgValidate,
		Error:   validator.ValidationErrors(err),
	})
}

func ResponseError(c *fiber.Ctx, err error, msg string) error {
	e := formatError(err)
	if e != nil {
		return c.Status(http.StatusBadRequest).JSON(responseErrorPayload{
			Error:   e,
			Message: msg,
		})
	}

	return c.Status(http.StatusInternalServerError).JSON(responseErrorPayload{
		Message: msg,
	})
}

func formatError(err error) (e map[string]string) {
	switch {
	case errors.Is(err, constant.ErrPasswordIncorrect):
		e = map[string]string{"password": utils.CapitalFirstLetter(err.Error())}
	case errors.Is(err, constant.ErrAccountNotFound) ||
		errors.Is(err, constant.ErrEmailAlreadyExists):
		e = map[string]string{"email": utils.CapitalFirstLetter(err.Error())}
	}

	return
}
