package handlers

import (
	"be-dashboard-nba/api/presenter"
	rolePresenter "be-dashboard-nba/api/presenter/role"
	"be-dashboard-nba/pkg/role/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// ReadRoles godoc
// @Summary      Get Roles List
// @Description  Retrieves a paginated and searchable list of all roles. Requires Bearer token.
// @Tags         Role
// @Produce      json
// @Param        search query    string false "Search term for role name"
// @Param        page   query    int    false "Page number (default: 1)"
// @Param        limit  query    int    false "Items per page (default: 10)"
// @Param        order  query    string false "Sort order (e.g., name ASC)"
// @Success      200  {object}  presenter.ResponsePayloadPaginate{data=[]presenter.ReadRoleResponse} "Successfully get roles"
// @Failure      400  {object}  presenter.ResponsePayloadMessage "Invalid query parameters"
// @Failure      401  {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      500  {object}  presenter.ResponsePayloadMessage "Failed get role"
// @Security     BearerAuth
// @Router       /roles [get]
func ReadRoles(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request rolePresenter.ReadRolesRequest
		if err = c.QueryParser(&request); err != nil {
			log.Errorw("error parse request", err)
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusBadRequest,
					Message: "Invalid query parameters",
				})
		}

		data, err := svc.ReadRolesService(c.UserContext(), request)
		if err != nil {
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusInternalServerError,
					Message: "Failed get role",
				})
		}

		return presenter.ResponsePaginate(c, presenter.ResponsePayloadPaginate{
			Code:       http.StatusOK,
			Message:    "Successfully get roles",
			Data:       rolePresenter.ToReadRoleListResponses(data.Data),
			Pagination: data.Pagination,
		})
	}
}
