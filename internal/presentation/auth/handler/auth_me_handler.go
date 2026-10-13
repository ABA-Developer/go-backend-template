package handler

import (
	"github.com/ABA-Developer/go-backend-template/constant"
	contract "github.com/ABA-Developer/go-backend-template/internal/domain/contract/usecase"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/logger"
	authInternal "github.com/ABA-Developer/go-backend-template/internal/presentation/auth"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/auth/presenter"
	response "github.com/ABA-Developer/go-backend-template/internal/presentation/response"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// @Summary      Get Current User Profile
// @Description  Fetches the profile details of the currently authenticated user. Requires Bearer token.
// @Tags         Authentication
// @Produce      json
// @Success      200 {object} response.DataPayload{data=presenter.UserResponse} "Successfully retrieved user profile"
// @Failure      401 {object} response.ErrorPayload "Unauthorized (Invalid or missing token)"
// @Failure      404 {object} response.ErrorPayload "User profile not found"
// @Failure      500 {object} response.ErrorPayload "Internal Server Error"
// @Security     BearerAuth
// @Router       /auth/me [get]
func AuthMe(svc contract.AuthUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		ah, err := authInternal.GetAuth(c)
		if err != nil {
			logger.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		data, err := svc.AuthMeUseCase(c.UserContext(), ah.GetClaims().UserID)

		if err != nil {
			if errors.Is(err, constant.ErrUserIdNotFound) {
				return response.Error(c, response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: constant.ErrUserIdNotFound.Error(),
				})
			}
			return response.Error(c, response.ErrorParam{
				Code:    http.StatusInternalServerError,
				Message: "Failed logout",
			})
		}

		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    presenter.ToReadAuthMeResponse(data),
			Message: "Successfully read profile user",
		})
	}
}
