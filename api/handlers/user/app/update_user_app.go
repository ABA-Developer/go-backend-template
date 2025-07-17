package app

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"be-dashboard-nba/api/handlers/user/payload"
	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/auth"
	"be-dashboard-nba/internal/validator"
	"be-dashboard-nba/pkg/user/service"
)

func updateUserApp(svc *service.Service, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request payload.UpdateUserRequest

		idParam := c.Params("id")
		if idParam == "" {
			return presenter.ResponseError(c, constant.ErrIDNull, "Param ID cannot be null")
		}

		request.ID, err = strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			return
		}

		if err = c.QueryParser(&request); err != nil {
			return presenter.ResponseError(c, err, "Failed parse request")
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

		err = svc.UpdateUserService(c.UserContext(), request, ah.GetClaims().UserID)
		if err != nil {
			return presenter.ResponseError(c, err, "Failed update user")
		}

		return presenter.ResponseData(c, presenter.ResponsePayload{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully update user",
		})
	}
}
