package handlers

import (
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/presentation/presenter"
	userPresenter "be-dashboard-nba/internal/presentation/presenter/user"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// Summary      Get List of Users
// @Description  Retrieves a paginated and searchable list of users in the system. Requires Bearer token.
// @Tags         User
// @Produce      json
// @Param        search query    string false "Search term for name or email"
// @Param        page   query    int    false "Page number (default: 1)"
// @Param        limit  query    int    false "Items per page (default: 10)"
// @Param        order  query    string false "Sort order (e.g., name ASC)"
// @Success      200  {object}  presenter.ResponsePayloadPaginate{data=[]userPresenter.ReadUserResponse} "Successfully get users"
// @Failure      400  {object}  presenter.ResponsePayloadMessage "Invalid query parameters"
// @Failure      401  {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      500  {object}  presenter.ResponsePayloadMessage "Failed get users"
// @Security     BearerAuth
// @Router       /users [get]
func ReadUsers(svc contract.UserUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request userPresenter.ReadUserRequest
		if err = c.QueryParser(&request); err != nil {
			log.Errorw("error parse request", err)
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusBadRequest,
					Message: "Invalid query parameters",
				})
		}

		data, err := svc.ReadUsersUseCase(c.UserContext(), request)
		if err != nil {
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusInternalServerError,
					Message: "Failed get users",
				})
		}

		return presenter.ResponsePaginate(c, presenter.ResponsePayloadPaginate{
			Code:       http.StatusOK,
			Message:    "Successfully get users",
			Data:       userPresenter.ToReadUserResponses(data.Data),
			Pagination: data.Pagination,
		})
	}
}
