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

// ReadRoleDetail godoc
// @Summary      Get Role Detail
// @Description  Fetches the details of a single role by its ID. Requires Bearer token.
// @Tags         Role
// @Produce      json
// @Param        role_id  path      int      true   "Role ID"
// @Success      200 {object}  presenter.ResponsePayloadData{data=presenter.ReadRoleResponse} "Successfully get role detail"
// @Failure      400 {object}  presenter.ResponsePayloadMessage "Bad Request (Invalid Role ID)"
// @Failure      401 {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      404 {object}  presenter.ResponsePayloadMessage "Not Found (Role ID not found)"
// @Failure      500 {object}  presenter.ResponsePayloadMessage "Failed get role detail"
// @Security     BearerAuth
// @Router       /roles/{role_id} [get]
func ReadRoleDetail(svc service.Service) fiber.Handler {
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

		data, err := svc.ReadRoleDetail(c.UserContext(), roleID)
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
						Message: "Failed get role detail",
					})
			}
		}
		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    rolePresenter.ToReadRoleResponse(data),
			Message: "Successfully get role detail",
		})
	}
}
