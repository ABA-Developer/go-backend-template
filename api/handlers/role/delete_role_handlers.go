package handlers

import (
	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/role/service"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

// DeleteRole godoc
// @Summary      Delete Role
// @Description  Deletes a role by its ID. Requires Bearer token.
// @Tags         Role
// @Produce      json
// @Param        role_id  path      int      true   "Role ID"
// @Success      200 {object}  presenter.ResponsePayloadData    "Successfully delete role"
// @Failure      400 {object}  presenter.ResponsePayloadMessage "Bad Request (Invalid Role ID)"
// @Failure      401 {object}  presenter.ResponsePayloadMessage "Unauthorized (Missing Auth)"
// @Failure      404 {object}  presenter.ResponsePayloadMessage "Role ID not found"
// @Failure      500 {object}  presenter.ResponsePayloadMessage "Failed delete role"
// @Security     BearerAuth
// @Router       /roles/{role_id} [delete]
func DeleteRole(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		roleIDParams := c.Params("role_id")
		roleID, err := strconv.Atoi(roleIDParams)
		if err != nil {
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusBadRequest,
					Message: "Failed parse params",
				})

		}

		err = svc.DeleteRoleService(c.UserContext(), roleID)

		if err != nil {
			if errors.Is(err, constant.ErrRoleIdNotFound) {
				return presenter.ResponseMessage(c,
					presenter.ResponsePayloadMessage{
						Code:    http.StatusNotFound,
						Message: constant.ErrRoleIdNotFound.Message,
					})
			} else {

				return presenter.ResponseMessage(c,
					presenter.ResponsePayloadMessage{
						Code:    http.StatusInternalServerError,
						Message: "Failed update role",
					})
			}
		}
		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully delete role",
		})
	}
}
