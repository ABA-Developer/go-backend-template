package handler

import (
	contract "github.com/ABA-Developer/go-backend-template/internal/domain/contract/usecase"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/logger"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/validator"
	authInternal "github.com/ABA-Developer/go-backend-template/internal/presentation/auth"
	response "github.com/ABA-Developer/go-backend-template/internal/presentation/response"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/user/mapper"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/user/presenter"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// @Summary      Update User
// @Description  Updates a user's profile details by their ID. Requires Bearer token.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user body     presenter.UpdateUserRequest true "User data to update"
// @Success      200  {object} response.DataPayload     "Successfully update user"
// @Failure      400  {object} response.ErrorPayload       "Bad Request (Invalid ID, Invalid JSON, or Validation Error)"
// @Failure      401  {object} response.ErrorPayload       "Unauthorized"
// @Failure      500  {object} response.ErrorPayload       "Internal Server Error"
// @Security     BearerAuth
// @Router       /users/me [put]
func UpdateProfileApp(svc contract.UserUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request presenter.UpdateUserRequest

		if err = c.BodyParser(&request); err != nil {
			logger.Error(err, "Failed parse request body")
			return response.Error(c, response.ErrorParam{
				Code:    http.StatusBadRequest,
				Message: "Failed parse request",
			})
		}

		if err := validate.Validate(request); err != nil {
			logger.WithContext(c.UserContext()).Error(err, "error validation update user request")
			return response.ErrorValidate(c, err)
		}

		ah, err := authInternal.GetAuth(c)
		if err != nil {
			logger.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		userIDFromToken := ah.GetClaims().UserID

		err = svc.UpdateUserUseCase(c.UserContext(), mapper.ToUpdateUserPayload(request), userIDFromToken, "")
		if err != nil {
			return response.Error(c, response.ErrorParam{
				Code:    http.StatusInternalServerError,
				Message: "Failed update user",
			})
		}

		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully update user",
		})
	}
}
