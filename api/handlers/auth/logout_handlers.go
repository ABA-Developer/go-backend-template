package handlers

import (
	"be-dashboard-nba/api/presenter"
	authInternal "be-dashboard-nba/internal/auth"
	"be-dashboard-nba/usecase/auth"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Logout godoc
// @Summary      User Logout
// @Description  Logs out the current user by invalidating their session/token. Requires Bearer token.
// @Tags         Authentication
// @Produce      json
// @Success      200 {object} presenter.ResponsePayloadData "Successfully logged out"
// @Failure      400 {object} presenter.ResponsePayloadMessage "Bad Request (e.g., Validation Error, Invalid JSON)"
// @Failure      500 {object} presenter.ResponsePayloadMessage "Internal Server Error (e.g., Database connection issue)"
// @Security     BearerAuth
// @Router       /auth/logout [post]
func Logout(svc auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		ah, err := authInternal.GetAuth(c)
		if err != nil {
			return
		}
		err = svc.LogoutService(c.UserContext(), ah.GetClaims(), c.IP())
		if err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusInternalServerError,
				Message: "Failed logout",
			})
		}

		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully logout",
		})
	}
}
