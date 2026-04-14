package handlers

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/auth"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/presentation/presenter"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/pkg/errors"
)

// Summary      Delete User
// @Description  Deletes a user by their ID. Requires Bearer token. Users cannot delete themselves.
// @Tags         User
// @Produce      json
// @Param        user_id path      string true "User ID"
// @Success      200     {object}  presenter.ResponsePayloadData "Successfully delete user"
// @Failure      400     {object}  presenter.ResponsePayloadMessage "Bad Request (Invalid User ID)"
// @Failure      401     {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      403     {object}  presenter.ResponsePayloadMessage "Forbidden (Attempt to delete own account)"
// @Failure      404     {object}  presenter.ResponsePayloadMessage "Not Found (User ID not found)"
// @Failure      500     {object}  presenter.ResponsePayloadMessage "Internal Server Error"
// @Security     BearerAuth
// @Router       /users/{user_id} [delete]
func DeleteUser(svc contract.UserUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		userIDParams := c.Params("user_id")

		ah, err := auth.GetAuth(c)
		if err != nil {
			log.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		userID := ah.GetClaims().UserID

		err = svc.DeleteUserUseCase(c.UserContext(), userIDParams, userID)
		if err != nil {
			if errors.Is(err, constant.ErrUserIdNotFound) {
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    http.StatusNotFound,
					Message: constant.ErrUserIdNotFound.Message,
				})
			}
			if errors.Is(err, constant.ErrForbiddenSelfDelete) {
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    constant.ErrForbiddenSelfDelete.Code,
					Message: constant.ErrForbiddenSelfDelete.Message,
				})
			}
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusInternalServerError,
				Message: "Failed update user",
			})
		}

		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully delete user",
		})
	}
}
