package handlers

import (
	"be-dashboard-nba/api/presenter"
	userPresenter "be-dashboard-nba/api/presenter/user"
	"be-dashboard-nba/pkg/user/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func ReadUsers(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request userPresenter.ReadUserRequest
		if err = c.QueryParser(&request); err != nil {
			log.Errorw("error parse request", err)
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusBadRequest,
					Message: "Invalid query parameters",
				})
		}

		data, err := svc.ReadUsersService(c.UserContext(), request)
		if err != nil {
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusInternalServerError,
					Message: "Failed get users",
				})
		}

		return presenter.ResponsePaginate(c, presenter.ResponsePayloadPaginate{
			Code:       http.StatusOK,
			Message:    "Successfully get users",
			Data:       userPresenter.ToReadUserResponses(data.Data),
			Pagination: data.Pagination,
		})
	}
}
