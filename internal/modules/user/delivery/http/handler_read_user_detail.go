package http

import (
	"be-dashboard-nba/api/presenter"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/user/domain"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func ReadUserDetail(usecase domain.UserUsecase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		userIDParams := c.Params("user_id")

		data, err := usecase.ReadDetailUserUsecase(c.UserContext(), userIDParams)

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
			Data:    ToReadUserDetailResponse(data),
			Message: "Successfully read profile user",
		})
	}
}
