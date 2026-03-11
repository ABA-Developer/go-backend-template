package http

import (
	"be-dashboard-nba/api/presenter"

	"be-dashboard-nba/internal/modules/user/domain"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func ReadUsers(usecase domain.UserUsecase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request ReadUserRequest
		if err = c.QueryParser(&request); err != nil {
			log.Errorw("error parse request", err)
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusBadRequest,
					Message: "Invalid query parameters",
				})
		}
		filter := request.ToDomainFilter()
		data, err := usecase.ReadUsersUsecase(c.UserContext(), filter)
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
			Data:       ToReadUserResponses(data.Data),
			Pagination: data.Pagination,
		})
	}
}
