package handlers

import (
	"be-dashboard-nba/api/presenter"
	userPresenter "be-dashboard-nba/api/presenter/user"
	"be-dashboard-nba/internal/auth"
	"be-dashboard-nba/pkg/user/service"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ReadProfileApp godoc
// @Summary      Get Current User Profile
// @Description  Fetches the profile details of the currently authenticated user. Requires Bearer token.
// @Tags         User Profile
// @Produce      json
// @Success      200 {object} presenter.ResponsePayloadData{data=presenter.ReadUserResponse} "Successfully retrieved user profile"
// @Failure      401 {object} presenter.ResponsePayloadMessage "Unauthorized (Invalid or missing token)"
// @Failure      404 {object} presenter.ResponsePayloadMessage "User profile not found"
// @Failure      500 {object} presenter.ResponsePayloadMessage "Internal Server Error"
// @Security     BearerAuth
// @Router       /users/me [get]
func ReadProfileApp(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		ah, err := auth.GetAuth(c)
		if err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusUnauthorized,
				Message: "Failed to get auth claims",
			})
		}
		ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
		defer cancel()

		data, err := svc.ReadDetailUserService(ctx, ah.GetClaims().UserID)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    http.StatusNotFound,
					Message: "Not found",
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
			Data:    userPresenter.ToReadUserResponse(data),
			Message: "Successfully read profile user",
		})
	}
}
