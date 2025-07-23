package app

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/api/handlers/auth/payload"
	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/internal/validator"
	"be-dashboard-nba/pkg/auth/service"
)

func loginApp(svc *service.Service, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request payload.LoginRequest
		if err = c.BodyParser(&request); err != nil {
			return presenter.ResponseError(c, err, "Failed parse request")
		}

		if err := validate.Validate(request); err != nil {
			log.Printf("error validation login request %v", err)
			return presenter.ResponseErrorValidate(c, err)
		}

		var (
			userAgent = string(c.Request().Header.UserAgent())
			iPAddress = c.IP()
		)

		data, user, err := svc.LoginService(c.UserContext(), request, userAgent, iPAddress)
		if err != nil {
			return presenter.ResponseError(c, err, "Failed login")
		}

		return presenter.ResponseData(c, presenter.ResponsePayload{
			Code:    http.StatusOK,
			Data:    payload.ToSessionResponse(data, user),
			Message: "Successfully login",
		})
	}
}
