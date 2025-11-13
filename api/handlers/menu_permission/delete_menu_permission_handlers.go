package handlers

import (
	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/menu_permission/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// DeleteMenuPermission godoc
// @Summary      Delete Menu Permission
// @Description  Deletes a specific permission from a menu by its unique permission ID. Requires Bearer token.
// @Tags         Menu Permission
// @Produce      json
// @Param        menu_permission_id path      int      true "Menu Permission ID"
// @Success      200 {object}  presenter.ResponsePayloadData "Successfully delete menu permission"
// @Failure      400 {object}  presenter.ResponsePayloadMessage "Bad Request (Invalid ID)"
// @Failure      401 {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      404 {object}  presenter.ResponsePayloadMessage "Not Found (Permission ID not found)"
// @Failure      500 {object}  presenter.ResponsePayloadMessage "Internal Server Error"
// @Security     BearerAuth
// @Router       /menu-permissions/{menu_permission_id} [delete]
func DeleteMenuPermission(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		menuPermissionIDParams := c.Params("menu_permission_id")
		menuPermisionID, err := strconv.Atoi(menuPermissionIDParams)
		if err != nil {
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusBadRequest,
					Message: "Failed parse params",
				})

		}

		err = svc.DeleteMenuPermissionService(c.UserContext(), menuPermisionID)
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
						Message: "Failed update menu permission",
					})
			}
		}
		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully delete menu permission",
		})
	}
}
