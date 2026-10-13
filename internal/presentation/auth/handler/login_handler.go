package handler

import (
	"github.com/ABA-Developer/go-backend-template/constant"
	contract "github.com/ABA-Developer/go-backend-template/internal/domain/contract/usecase"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/validator"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/auth/mapper"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/auth/presenter"
	response "github.com/ABA-Developer/go-backend-template/internal/presentation/response"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// @Summary      User Login
// @Description  Authenticate user with email and password, returning tokens and user data. This endpoint uses a centralized error handler.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param 		 auth body    presenter.LoginRequest true "Login credentials"
// @Success      200 {object} response.DataPayload{data=presenter.SessionResponse} "Successful login"
// @Failure      400 {object} response.ErrorPayload "Bad Request (e.g., Validation Error, Invalid JSON)"
// @Failure      401 {object} response.ErrorPayload "Unauthorized (e.g., Wrong email or password)"
// @Failure      500 {object} response.ErrorPayload "Internal Server Error (e.g., Database connection issue)"
// @Router       /auth/login [post]
func Login(svc contract.AuthUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request presenter.LoginRequest

		if err = c.BodyParser(&request); err != nil {
			return response.Error(c, response.ErrorParam{
				Code:    http.StatusBadRequest,
				Message: "Failed parse request",
			})
		}

		if err := validate.Validate(request); err != nil {
			return response.ErrorValidate(c, err)
		}

		var (
			userAgent = string(c.Request().Header.UserAgent())
			iPAddress = c.IP()
		)

		data, user, err := svc.LoginUseCase(c.UserContext(), mapper.ToLoginPayload(request), userAgent, iPAddress)
		if err != nil {
			if errors.Is(err, constant.ErrWrongEmailOrPassword) {

				return response.Error(c, response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: constant.ErrWrongEmailOrPassword.Error(),
				})
			} else {
				return response.Error(c, response.ErrorParam{
					Code:    http.StatusInternalServerError,
					Message: "Failed login",
				})
			}
		}

		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    presenter.ToSessionResponse(data, user),
			Message: "Successfully login",
		})
	}
}
