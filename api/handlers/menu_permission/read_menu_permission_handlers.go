package handlers

import (
	"be-dashboard-nba/api/presenter"
	menuPermissionPresenter "be-dashboard-nba/api/presenter/menu_permission"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/menu_permission/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// ReadMenuPermissionListParams godoc
// @Summary      Get Menu Permission List
// @Description  Retrieves a paginated and searchable list of permissions (e.g., read, create) for a specific menu ID. Requires Bearer token.
// @Tags         Menu Permission
// @Produce      json
// @Param        menu_id  path      int    true   "ID of the Menu"
// @Param        search   query     string false  "Search term for action_name"
// @Param        page     query     int    false  "Page number (default: 1)"
// @Param        limit    query     int    false  "Items per page (default: 10)"
// @Param        order    query     string false  "Sort order (e.g., action_name ASC)"
// @Success      200      {object}  presenter.ResponsePayloadPaginate{data=[]presenter.ReadMenuPermissionResponse} "Successfully get menu permission list"
// @Failure      400      {object}  presenter.ResponsePayloadMessage "Bad Request (Invalid Menu ID)"
// @Failure      401      {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      404      {object}  presenter.ResponsePayloadMessage "Not Found (Menu ID not found)"
// @Failure      500      {object}  presenter.ResponsePayloadMessage "Internal Server Error (Query parsing or service failure)"
// @Security     BearerAuth
// @Router       /menu-permissions/{menu_id} [get]
func ReadMenuPermissionListParams(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		menuIDParams := c.Params("menu_id")
		menuID, err := strconv.Atoi(menuIDParams)
		if err != nil {
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusBadRequest,
					Message: "menu id params must be in number",
				})
		}

		var request menuPermissionPresenter.ReadMenuPermissionListRequest
		if err = c.QueryParser(&request); err != nil {
			log.Errorw("error parse request", err)
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusBadRequest,
					Message: "Invalid query parameters",
				})
		}

		data, err := svc.ReadMenuPermissionListParams(c.UserContext(), request, menuID)
		if err != nil {
			if errors.Is(err, constant.ErrMenuIdNotFound) {
				return presenter.ResponseMessage(c,
					presenter.ResponsePayloadMessage{
						Code:    http.StatusNotFound,
						Message: constant.ErrMenuIdNotFound.Message,
					})
			} else {

				return presenter.ResponseMessage(c,
					presenter.ResponsePayloadMessage{
						Code:    http.StatusInternalServerError,
						Message: "Failed get menu permission",
					})
			}
		}
		return presenter.ResponsePaginate(c, presenter.ResponsePayloadPaginate{
			Code:       http.StatusOK,
			Message:    "Successfully get menu permission list",
			Data:       menuPermissionPresenter.ToReadMenuPermissionListResponses(data.Data),
			Pagination: data.Pagination,
		})
	}
}
