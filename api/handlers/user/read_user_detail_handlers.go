package handlers

import (
	"be-dashboard-nba/api/presenter"
	userPresenter "be-dashboard-nba/api/presenter/user"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/user/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func ReadUserDetail(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		userIDParams := c.Params("user_id")

		data, err := svc.ReadDetailUserService(c.UserContext(), userIDParams)

		if err != nil {
			if errors.Is(err, constant.ErrUserIdNotFound) {
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    constant.ErrUserIdNotFound.Code,
					Message: constant.ErrUserIdNotFound.Message,
				})
			} else {
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    http.StatusInternalServerError,
					Message: "failed to get profile",
				})
			}
		}

		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    userPresenter.ToReadUserDetailResponse(data),
			Message: "Successfully read profile user",
		})
	}
}
