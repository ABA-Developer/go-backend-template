package handlers

import (
	"be-dashboard-nba/api/presenter"
	authPresenter "be-dashboard-nba/api/presenter/auth"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/auth"
	"be-dashboard-nba/pkg/auth/service"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// ReadProfileApp godoc
// @Summary      Get Current User Profile
// @Description  Fetches the profile details of the currently authenticated user. Requires Bearer token.
// @Tags         Authentication
// @Produce      json
// @Success      200 {object} presenter.ResponsePayloadData{data=presenter.UserResponse} "Successfully retrieved user profile"
// @Failure      401 {object} presenter.ResponsePayloadMessage "Unauthorized (Invalid or missing token)"
// @Failure      404 {object} presenter.ResponsePayloadMessage "User profile not found"
// @Failure      500 {object} presenter.ResponsePayloadMessage "Internal Server Error"
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
			if errors.Is(err, constant.ErrUserIdNotFound) {
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    http.StatusBadRequest,
					Message: constant.ErrUserIdNotFound.Error(),
				})
			}
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusInternalServerError,
				Message: "Failed logout",
			})
		}

		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    authPresenter.ToReadAuthMeResponse(data),
			Message: "Successfully read profile user",
		})
	}
}
