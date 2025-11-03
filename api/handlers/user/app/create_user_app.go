package app

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"be-dashboard-nba/api/handlers/user/payload"
	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/internal/auth"
	"be-dashboard-nba/internal/validator"
	"be-dashboard-nba/pkg/user/service"
)

func createUserApp(svc *service.Service, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request payload.CreateUserRequest
		if err = c.QueryParser(&request); err != nil {
			return presenter.ResponseError(c, err, "Failed parse request")
		}

		if err := validate.Validate(request); err != nil {
			log.WithContext(c.UserContext()).Error("error validation create user request %v", err)
			return presenter.ResponseErrorValidate(c, err)
		}

		ah, err := auth.GetAuth(c)
		if err != nil {
			log.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		err = svc.CreateUserService(c.UserContext(), request, ah.GetClaims().UserID)
		if err != nil {
			return presenter.ResponseError(c, err, "Failed create user")
		}

		return presenter.ResponseData(c, presenter.ResponsePayload{
			Code:    http.StatusCreated,
			Data:    nil,
			Message: "Successfully create user",
		})
	}
}
