package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"be-dashboard-nba/api/presenter"
	authPresenter "be-dashboard-nba/api/presenter/auth"
	"be-dashboard-nba/internal/auth"
	"be-dashboard-nba/internal/validator"
	"be-dashboard-nba/pkg/auth/service"
)

// @Summary      User Login
// @Description  Authenticate user with email and password, returning tokens and user data. This endpoint uses a centralized error handler.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param 		 auth body    presenter.LoginRequest true "Login credentials"
// @Success      200 {object} presenter.ResponsePayload{data=presenter.SessionResponse} "Successful login"
// @Failure      400 {object} presenter.responseErrorPayload "Bad Request (e.g., Validation Error, Invalid JSON)"
// @Failure      401 {object} presenter.responseErrorPayload "Unauthorized (e.g., Wrong email or password)"
// @Failure      500 {object} presenter.responseErrorPayload "Internal Server Error (e.g., Database connection issue)"
// @Router       /auth/login [post]
func Login(svc service.Service, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request authPresenter.LoginRequest

		if err = c.BodyParser(&request); err != nil {
			return presenter.ResponseError(c, err, "Failed parse request")
		}
		fmt.Println(request)
		if err := validate.Validate(request); err != nil {
			return presenter.ResponseErrorValidate(c, err)
		}

		var (
			userAgent = string(c.Request().Header.UserAgent())
			iPAddress = c.IP()
		)

		data, user, err := svc.LoginService(c.UserContext(), request, userAgent, iPAddress)
		if err != nil {
			return err
		}

		return presenter.ResponseData(c, presenter.ResponsePayload{
			Code:    http.StatusOK,
			Data:    authPresenter.ToSessionResponse(data, user),
			Message: "Successfully login",
		})
	}
}

// Logout godoc
// @Summary      User Logout
// @Description  Logs out the current user by invalidating their session/token. Requires Bearer token.
// @Tags         Authentication
// @Produce      json
// @Success      200 {object} presenter.ResponsePayload "Successfully logged out"
// @Failure      400 {object} presenter.responseErrorPayload "Bad Request (e.g., Validation Error, Invalid JSON)"
// @Failure      500 {object} presenter.responseErrorPayload "Internal Server Error (e.g., Database connection issue)"
// @Security     BearerAuth
// @Router       /auth/logout [post]
func Logout(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		ah, err := auth.GetAuth(c)
		if err != nil {
			return
		}
		err = svc.LogoutService(c.UserContext(), ah.GetClaims(), c.IP())
		if err != nil {
			return presenter.ResponseError(c, err, "Failed logout")
		}

		return presenter.ResponseData(c, presenter.ResponsePayload{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully logout",
		})
	}
}

// ReadProfileApp godoc
// @Summary      Get Current User Profile
// @Description  Fetches the profile details of the currently authenticated user. Requires Bearer token.
// @Tags         Authentication
// @Produce      json
// @Success      200 {object} presenter.ResponsePayload{data=presenter.UserResponse} "Successfully retrieved user profile"
// @Failure      401 {object} presenter.responseErrorPayload "Unauthorized (Invalid or missing token)"
// @Failure      404 {object} presenter.responseErrorPayload "User profile not found"
// @Failure      500 {object} presenter.responseErrorPayload "Internal Server Error"
// @Security     BearerAuth
// @Router       /auth/me [get]
func AuthMe(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		ah, err := auth.GetAuth(c)
		if err != nil {
			log.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		data, err := svc.AuthMeService(c.UserContext(), ah.GetClaims().UserID)

		if err != nil {
			return err
		}

		return presenter.ResponseData(c, presenter.ResponsePayload{
			Code:    http.StatusOK,
			Data:    authPresenter.ToReadAuthMeResponse(data),
			Message: "Successfully read profile user",
		})
	}
}
