package http

import (
	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/core/validator"
	"be-dashboard-nba/internal/modules/auth/domain"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// @Summary      User Login
// @Description  Authenticate user with email and password, returning tokens and user data. This endpoint uses a centralized error handler.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param 		 auth body    LoginRequest true "Login credentials"
// @Success      200 {object} presenter.ResponsePayloadData{data=SessionResponse} "Successful login"
// @Failure      400 {object} presenter.ResponsePayloadMessage "Bad Request (e.g., Validation Error, Invalid JSON)"
// @Failure      401 {object} presenter.ResponsePayloadMessage "Unauthorized (e.g., Wrong email or password)"
// @Failure      500 {object} presenter.ResponsePayloadMessage "Internal Server Error (e.g., Database connection issue)"
// @Router       /auth/login [post]

type LoginResponse struct {
	Session  domain.Session `json:"session"`
	UserData domain.User    `json:"user_data"`
}

func Login(authUsecase domain.AuthUsecase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request LoginRequest

		if err = c.BodyParser(&request); err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusBadRequest,
				Message: "Failed parse request",
			})
		}

		if err := validate.Validate(request); err != nil {
			return presenter.ResponseErrorValidate(c, err)
		}

		var (
			userAgent = string(c.Request().Header.UserAgent())
			ipAddress = c.IP()
		)

		session, user, err := authUsecase.LoginUsecase(c.Context(), request.Email, request.Password, userAgent, ipAddress)
		if err != nil {
			if errors.Is(err, constant.ErrWrongEmailOrPassword) {

				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    http.StatusBadRequest,
					Message: constant.ErrWrongEmailOrPassword.Error(),
				})
			} else {
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    http.StatusInternalServerError,
					Message: "Failed login",
				})
			}
		}

		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    ToSessionResponse(session, user),
			Message: "Successfully login",
		})
	}
}
