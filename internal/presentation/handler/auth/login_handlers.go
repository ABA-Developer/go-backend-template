package handlers

import (
	"be-dashboard-nba/internal/presentation/presenter"
	authPresenter "be-dashboard-nba/internal/presentation/presenter/auth"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/infrastructure/validator"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// @Summary      User Login
// @Description  Authenticate user with email and password, returning tokens and user data. This endpoint uses a centralized error handler.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param 		 auth body    authPresenter.LoginRequest true "Login credentials"
// @Success      200 {object} presenter.ResponsePayloadData{data=authPresenter.SessionResponse} "Successful login"
// @Failure      400 {object} presenter.ResponsePayloadMessage "Bad Request (e.g., Validation Error, Invalid JSON)"
// @Failure      401 {object} presenter.ResponsePayloadMessage "Unauthorized (e.g., Wrong email or password)"
// @Failure      500 {object} presenter.ResponsePayloadMessage "Internal Server Error (e.g., Database connection issue)"
// @Router       /auth/login [post]
func Login(svc contract.AuthUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request authPresenter.LoginRequest

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
			iPAddress = c.IP()
		)

		data, user, err := svc.LoginUseCase(c.UserContext(), request, userAgent, iPAddress)
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
			Data:    authPresenter.ToSessionResponse(data, user),
			Message: "Successfully login",
		})
	}
}
