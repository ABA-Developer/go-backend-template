package app

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"be-dashboard-nba/api/handlers/user/payload"
	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/internal/validator"
	"be-dashboard-nba/pkg/user/service"
)

func readListUserApp(svc *service.Service, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request payload.ReadUserListRequest
		if err = c.QueryParser(&request); err != nil {
			return presenter.ResponseError(c, err, "Failed parse request")
		}

		if err := validate.Validate(request); err != nil {
			log.WithContext(c.UserContext()).Error("error validation read list user request %v", err)
			return presenter.ResponseErrorValidate(c, err)
		}

		data, totalData, err := svc.ReadListUserService(c.UserContext(), request)
		if err != nil {
			return presenter.ResponseError(c, err, "Failed read list user")
		}

		return presenter.ResponsePaginate(c, request.PaginationPayload, totalData, presenter.ResponsePayload{
			Code:    http.StatusOK,
			Data:    payload.ToReadUserResponses(data),
			Message: "Successfully read list user",
		})
	}
}
