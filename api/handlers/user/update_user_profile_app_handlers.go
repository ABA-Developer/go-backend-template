package handlers

import (
	"be-dashboard-nba/api/presenter"
	userPresenter "be-dashboard-nba/api/presenter/user"
	"be-dashboard-nba/internal/auth"
	"be-dashboard-nba/internal/validator"
	"be-dashboard-nba/pkg/user/service"
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// UpdateUserApp godoc
// @Summary      Update User
// @Description  Updates a user's profile details by their ID. Requires Bearer token.
// @Tags         User Profile
// @Accept       json
// @Produce      json
// @Param        user body     presenter.UpdateUserRequest true "User data to update"
// @Success      200  {object} presenter.ResponsePayloadData     "Successfully update user"
// @Failure      400  {object} presenter.ResponsePayloadMessage       "Bad Request (Invalid ID, Invalid JSON, or Validation Error)"
// @Failure      401  {object} presenter.ResponsePayloadMessage       "Unauthorized"
// @Failure      500  {object} presenter.ResponsePayloadMessage       "Internal Server Error"
// @Security     BearerAuth
// @Router       /users/me [put]
func UpdateProfileApp(svc service.Service, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request userPresenter.UpdateUserRequest

		if err = c.BodyParser(&request); err != nil {
			log.Error(err, "Failed parse request body")
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusBadRequest,
				Message: "Failed parse request",
			})
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

		userIDFromToken := ah.GetClaims().UserID
		ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
		defer cancel()
		err = svc.UpdateUserService(ctx, request, userIDFromToken)
		if err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusInternalServerError,
				Message: "Failed update user",
			})
		}

		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully update user",
		})
	}
}
