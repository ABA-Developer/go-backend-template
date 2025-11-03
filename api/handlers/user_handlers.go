package handlers

import (
	"be-dashboard-nba/api/presenter"
	userPresenter "be-dashboard-nba/api/presenter/user"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/auth"
	"be-dashboard-nba/internal/validator"
	"be-dashboard-nba/pkg/user/service"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func CreateUserApp(svc service.Service, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request userPresenter.CreateUserRequest
		if err = c.QueryParser(&request); err != nil {
			return presenter.ResponseError(c, err, "Failed parse request")
		}

		if err := validate.Validate(request); err != nil {
			log.WithContext(c.UserContext()).Error("error validation create user request %v", err)
			return presenter.ResponseErrorValidate(c, err)
		}

		ah, err := auth.GetAuth(c)
		if err != nil {
			log.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		err = svc.CreateUserService(c.UserContext(), request, ah.GetClaims().UserID)
		if err != nil {
			return presenter.ResponseError(c, err, "Failed create user")
		}

		return presenter.ResponseData(c, presenter.ResponsePayload{
			Code:    http.StatusCreated,
			Data:    nil,
			Message: "Successfully create user",
		})
	}
}

func DeleteUserApp(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		idParam := c.Params("id")
		if idParam == "" {
			return presenter.ResponseError(c, constant.ErrIDNull, "Param ID cannot be null")
		}

		err = svc.DeleteUserService(c.UserContext(), idParam)
		if err != nil {
			return presenter.ResponseError(c, err, "Failed delete user")
		}

		return presenter.ResponseData(c, presenter.ResponsePayload{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully delete user",
		})
	}
}

// ReadProfileApp godoc
// @Summary      Get Current User Profile
// @Description  Fetches the profile details of the currently authenticated user. Requires Bearer token.
// @Tags         User Profile
// @Produce      json
// @Success      200 {object} presenter.ResponsePayload{data=presenter.ReadUserResponse} "Successfully retrieved user profile"
// @Failure      401 {object} presenter.responseErrorPayload "Unauthorized (Invalid or missing token)"
// @Failure      404 {object} presenter.responseErrorPayload "User profile not found"
// @Failure      500 {object} presenter.responseErrorPayload "Internal Server Error"
// @Security     BearerAuth
// @Router       /users/me [get]
func ReadProfileApp(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		ah, err := auth.GetAuth(c)
		if err != nil {
			log.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}
		ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
		defer cancel()

		data, err := svc.ReadDetailUserService(ctx, ah.GetClaims().UserID)

		if err != nil {
			return err
		}

		fmt.Println(data)

		return presenter.ResponseData(c, presenter.ResponsePayload{
			Code:    http.StatusOK,
			Data:    userPresenter.ToReadUserResponse(data),
			Message: "Successfully read profile user",
		})
	}
}

// UpdateUserApp godoc
// @Summary      Update User
// @Description  Updates a user's profile details by their ID. Requires Bearer token.
// @Tags         User Profile
// @Accept       json
// @Produce      json
// @Param        id   path     string                        true "User ID"
// @Param        user body     presenter.UpdateUserRequest true "User data to update"
// @Success      200  {object} presenter.ResponsePayload     "Successfully update user"
// @Failure      400  {object} presenter.responseErrorPayload       "Bad Request (Invalid ID, Invalid JSON, or Validation Error)"
// @Failure      401  {object} presenter.responseErrorPayload       "Unauthorized"
// @Failure      500  {object} presenter.responseErrorPayload       "Internal Server Error"
// @Security     BearerAuth
// @Router       /users/{id} [put]
func UpdateUserApp(svc service.Service, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request userPresenter.UpdateUserRequest

		idParam := c.Params("id")
		if idParam == "" {
			return presenter.ResponseError(c, constant.ErrIDNull, "Param ID cannot be null")
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

// UpdateUserApp godoc
// @Summary      Update User
// @Description  Updates a user's profile details by their ID. Requires Bearer token.
// @Tags         User Profile
// @Accept       json
// @Produce      json
// @Param        user body     presenter.UpdateUserRequest true "User data to update"
// @Success      200  {object} presenter.ResponsePayload     "Successfully update user"
// @Failure      400  {object} presenter.responseErrorPayload       "Bad Request (Invalid ID, Invalid JSON, or Validation Error)"
// @Failure      401  {object} presenter.responseErrorPayload       "Unauthorized"
// @Failure      500  {object} presenter.responseErrorPayload       "Internal Server Error"
// @Security     BearerAuth
// @Router       /users/me [put]
func UpdateProfileApp(svc service.Service, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request userPresenter.UpdateUserRequest

		if err = c.BodyParser(&request); err != nil {
			log.Error(err, "Failed parse request body")
			return presenter.ResponseError(c, err, "Failed parse request")
		}
		fmt.Println(request)

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

		err = svc.UpdateUserService(c.UserContext(), request, userIDFromToken)
		if err != nil {
			return err
		}

		return presenter.ResponseData(c, presenter.ResponsePayload{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully update user",
		})
	}
}
