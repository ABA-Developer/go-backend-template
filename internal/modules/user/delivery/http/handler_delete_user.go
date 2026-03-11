package http

import (
	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/auth"
	"be-dashboard-nba/internal/modules/user/domain"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/pkg/errors"
)

func DeleteUser(usecase domain.UserUsecase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		userIDParams := c.Params("user_id")

		ah, err := auth.GetAuth(c)
		if err != nil {
			log.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		userID := ah.GetClaims().UserID

		err = usecase.DeleteUserUsecase(c.UserContext(), userIDParams, userID)
		if err != nil {
			if errors.Is(err, constant.ErrUserIdNotFound) {
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    http.StatusNotFound,
					Message: constant.ErrUserIdNotFound.Message,
				})
			}
			if errors.Is(err, constant.ErrForbiddenSelfDelete) {
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    constant.ErrForbiddenSelfDelete.Code,
					Message: constant.ErrForbiddenSelfDelete.Message,
				})
			}
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusInternalServerError,
				Message: "Failed update user",
			})
		}

		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully delete user",
		})
	}
}
