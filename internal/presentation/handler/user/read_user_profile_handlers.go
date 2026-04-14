package handlers

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/auth"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/presentation/presenter"
	userPresenter "be-dashboard-nba/internal/presentation/presenter/user"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// ReadProfileApp godoc
// @Summary      Get Current User Profile
// @Description  Fetches the profile details of the currently authenticated user. Requires Bearer token.
// @Tags         User
// @Produce      json
// @Success      200 {object} presenter.ResponsePayloadData{data=userPresenter.ReadUserResponse} "Successfully retrieved user profile"
// @Failure      401 {object} presenter.ResponsePayloadMessage "Unauthorized (Invalid or missing token)"
// @Failure      404 {object} presenter.ResponsePayloadMessage "User profile not found"
// @Failure      500 {object} presenter.ResponsePayloadMessage "Internal Server Error"
// @Security     BearerAuth
// @Router       /users/me [get]
func ReadProfileApp(svc contract.UserUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		ah, err := auth.GetAuth(c)
		if err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusUnauthorized,
				Message: "Failed to get auth claims",
			})
		}

		data, err := svc.ReadUserProfileUseCase(c.UserContext(), ah.GetClaims().UserID)

		if err != nil {
			if errors.Is(err, constant.ErrUserIdNotFound) {
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    constant.ErrUserIdNotFound.Code,
					Message: constant.ErrUserIdNotFound.Message,
				})
			} else {
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    http.StatusInternalServerError,
					Message: "failed to get user profile",
				})
			}
		}

		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    userPresenter.ToReadUserProfileUseCaseResponse(data),
			Message: "Successfully read user profile",
		})
	}
}
