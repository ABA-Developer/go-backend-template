package handlers

import (
	"be-dashboard-nba/constant"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/presentation/presenter"
	userPresenter "be-dashboard-nba/internal/presentation/presenter/user"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

// Summary      Get User Detail
// @Description  Fetches the details of a single user by their ID. Requires Bearer token.
// @Tags         User
// @Produce      json
// @Param        user_id  path      string   true "User ID"
// @Success      200 {object} presenter.ResponsePayloadData{data=userPresenter.ReadUserDetailResponse} "Successfully get user detail"
// @Failure      400 {object} presenter.ResponsePayloadMessage "Bad Request (Invalid User ID)"
// @Failure      401 {object} presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      404 {object} presenter.ResponsePayloadMessage "Not Found (User ID not found)"
// @Failure      500 {object} presenter.ResponsePayloadMessage "Internal Server Error"
// @Security     BearerAuth
// @Router       /users/{user_id} [get]
func ReadUserDetail(svc contract.UserUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		userIDParams := c.Params("user_id")

		data, err := svc.ReadDetailUserUseCase(c.UserContext(), userIDParams)
		if err != nil {
			if errors.Is(err, constant.ErrUserIdNotFound) {
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    constant.ErrUserIdNotFound.Code,
					Message: constant.ErrUserIdNotFound.Message,
				})
			}

			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusInternalServerError,
				Message: "failed to get profile",
			})
		}

		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    userPresenter.ToReadUserDetailResponse(data),
			Message: "Successfully read profile user",
		})
	}
}
