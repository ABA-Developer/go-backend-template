package http

import (
	"be-dashboard-nba/api/presenter"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/auth"
	"be-dashboard-nba/internal/modules/user/domain"
	"be-dashboard-nba/internal/core/validator"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func UpdateUser(usecase domain.UserUsecase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		userIDParams := c.Params("user_id")

		var request UpdateUserRequest

		if err = c.BodyParser(&request); err != nil {
			log.Error(err, "Failed parse request body")
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusBadRequest,
				Message: "Failed parse request",
			})
		}

		if err := validate.Validate(request); err != nil {
			log.WithContext(c.UserContext()).Error("error validation update user request %v", err)
			return presenter.ResponseErrorValidate(c, err)
		}

		ah, err := auth.GetAuth(c)
		if err != nil {
			log.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		userIDFromToken := ah.GetClaims().UserID

		payload := request.ToDomainPayload(userIDParams, userIDFromToken)
		err = usecase.UpdateUserUsecase(c.UserContext(), payload)
		if err != nil {
			if errors.Is(err, constant.ErrRoleIdNotFound) {
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    http.StatusNotFound,
					Message: constant.ErrRoleIdNotFound.Message,
				})
			}
			if errors.Is(err, constant.ErrUserIdNotFound) {
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    http.StatusNotFound,
					Message: constant.ErrUserIdNotFound.Message,
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
			Message: "Successfully update user",
		})
	}
}
