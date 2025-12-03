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
)

// ReadMenuPermissionDetail godoc
// @Summary      Get Menu Permission Detail
// @Description  Fetches the details of a single menu permission by its unique ID. Requires Bearer token.
// @Tags         Menu Permission
// @Produce      json
// @Param        menu_permission_id path      int      true "Menu Permission ID"
// @Success      200 {object}  presenter.ResponsePayloadData{data=presenter.ReadMenuPermissionResponse} "Successfully get menu permission detail"
// @Failure      400 {object}  presenter.ResponsePayloadMessage "Bad Request (Invalid ID)"
// @Failure      401 {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      404 {object}  presenter.ResponsePayloadMessage "Not Found (Permission ID not found)"
// @Failure      500 {object}  presenter.ResponsePayloadMessage "Internal Server Error"
// @Security     BearerAuth
// @Router       /menu-permissions/detail/{menu_permission_id} [get]
func ReadMenuPermissionDetail(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		menuPermissionIDParams := c.Params("menu_permission_id")
		menuPermissionID, err := strconv.Atoi(menuPermissionIDParams)
		if err != nil {
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusBadRequest,
					Message: "Failed parse params",
				})

		}

		data, err := svc.ReadMenuPermissionDetail(c.UserContext(), menuPermissionID)
		if err != nil {
			if errors.Is(err, constant.ErrMenuPermissionIdNotFound) {
				return presenter.ResponseMessage(c,
					presenter.ResponsePayloadMessage{
						Code:    http.StatusNotFound,
						Message: constant.ErrMenuPermissionIdNotFound.Message,
					})
			} else {

				return presenter.ResponseMessage(c,
					presenter.ResponsePayloadMessage{
						Code:    http.StatusInternalServerError,
						Message: "Failed get menu permission detail",
					})
			}
		}
		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    menuPermissionPresenter.ToReadMenuPermissionListResponse(data),
			Message: "Successfully get menu permission detail",
		})
	}
}
