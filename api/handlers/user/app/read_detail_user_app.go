package app

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/api/handlers/user/payload"
	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/user/service"
)

func readDetailUserApp(svc *service.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		idParam := c.Params("id")
		if idParam == "" {
			return presenter.ResponseError(c, constant.ErrIDNull, "Param ID cannot be null")
		}

		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			return
		}

		data, err := svc.ReadDetailUserService(c.UserContext(), id)
		if err != nil {
			return presenter.ResponseError(c, err, "Failed read detail user")
		}

		return presenter.ResponseData(c, presenter.ResponsePayload{
			Code:    http.StatusOK,
			Data:    payload.ToReadUserResponse(data),
			Message: "Successfully read detail user",
		})
	}
}
