package handlers

import (
	"be-dashboard-nba/api/presenter"
	rolePresenter "be-dashboard-nba/api/presenter/role"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/role/service"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

// ReadRoleAccess godoc
// @Summary      Get Role Menu Permission
// @Description  Retrieves all menu permissions associated with a specific role. Requires Bearer token.
// @Tags         Role Menu Permission
// @Produce      json
// @Param        role_id  path      int   true  "Role ID"
// @Success      200      {object}  presenter.ResponsePayloadData{data=presenter.RoleAccessResponse} "Successfully get role menu permission"
// @Failure      400      {object}  presenter.ResponsePayloadMessage "Invalid Role ID"
// @Failure      401      {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      404      {object}  presenter.ResponsePayloadMessage "Role ID not found"
// @Failure      500      {object}  presenter.ResponsePayloadMessage "Failed get role menu permission"
// @Security     BearerAuth
// @Router       /roles/role-access/{role_id} [get]
func ReadRoleAccess(svc service.Service) fiber.Handler {
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

		data, err := svc.ReadRoleAccessService(c.UserContext(), roleID)
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
						Message: "Failed get role menu permission",
					})
			}
		}
		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    rolePresenter.ToReadRoleAccessResponse(data),
			Message: "Successfully get role detail",
		})
	}
}
