package handlers

import (
	"be-dashboard-nba/api/presenter"
	rolePresenter "be-dashboard-nba/api/presenter/role"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/validator"
	"be-dashboard-nba/pkg/role/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

// UpdateRoleMenuPermission godoc
// @Summary      Update Role Menu Permission
// @Description  Updates menu permission assigned to a specific role. Requires Bearer token.
// @Tags         Role Menu Permission
// @Accept       json
// @Produce      json
// @Param        role_id  path      int                                     true  "Role ID"
// @Param        payload  body      presenter.UpdateRoleAccessRequest   true  "Role menu permission payload"
// @Success      200      {object}  presenter.ResponsePayloadData           "Successfully update role menu permission"
// @Failure      400      {object}  presenter.ResponseErrorPayload          "Bad Request (Invalid ID, Invalid JSON, or Validation Error)"
// @Failure      401      {object}  presenter.ResponsePayloadMessage        "Unauthorized"
// @Failure      404      {object}  presenter.ResponsePayloadMessage        "Role ID or Menu Permission not found"
// @Failure      500      {object}  presenter.ResponsePayloadMessage        "Internal server error"
// @Security     BearerAuth
// @Router       /roles/role-access/{role_id} [put]
func UpdateRoleMenuPermission(svc service.Service, validate *validator.Validator) fiber.Handler {
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
		var request rolePresenter.UpdateRoleAccessRequest
		fmt.Println(request)
		if err = c.BodyParser(&request); err != nil {
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusBadRequest,
					Message: "Failed parse request",
				})
		}
		if err := validate.Validate(request); err != nil {
			return presenter.ResponseErrorValidate(c, err)
		}

		err = svc.UpdateRoleAccessService(c.UserContext(), roleID, request)

		if err != nil {
			switch {
			case errors.Is(err, constant.ErrRoleIdNotFound):
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    http.StatusNotFound,
					Message: constant.ErrRoleIdNotFound.Message,
				})
			case errors.Is(err, constant.ErrMenuPermissionIdNotFound):
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    http.StatusNotFound,
					Message: constant.ErrMenuPermissionIdNotFound.Message,
				})
			case err.Error() == constant.ErrMsgValidate:
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    http.StatusBadRequest,
					Message: constant.ErrMsgValidate,
				})
			default:
				return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
					Code:    http.StatusInternalServerError,
					Message: constant.ErrMsgUnknownSource,
				})
			}
		}

		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully update role menu permission",
		})
	}
}
