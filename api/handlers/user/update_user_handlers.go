package handlers

import (
	"be-dashboard-nba/api/presenter"
	userPresenter "be-dashboard-nba/api/presenter/user"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/auth"
	"be-dashboard-nba/internal/validator"
	"be-dashboard-nba/pkg/user/service"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func UpdateUser(svc service.Service, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		userIDParams := c.Params("user_id")

		var request userPresenter.UpdateUserRequest

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

		err = svc.UpdateUserService(c.UserContext(), request, userIDFromToken, userIDParams)
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
