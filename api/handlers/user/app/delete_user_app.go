package app

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/user/service"
)

func deleteUserApp(svc *service.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		idParam := c.Params("id")
		if idParam == "" {
			return presenter.ResponseError(c, constant.ErrIDNull, "Param ID cannot be null")
		}

		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			return
		}

		err = svc.DeleteUserService(c.UserContext(), id)
		if err != nil {
			return presenter.ResponseError(c, err, "Failed delete user")
		}

		return presenter.ResponseData(c, presenter.ResponsePayload{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully delete user",
		})
	}
}
