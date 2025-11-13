package handlers

import (
	"be-dashboard-nba/api/presenter"
	rolePresenter "be-dashboard-nba/api/presenter/role"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/auth"
	"be-dashboard-nba/internal/validator"
	"be-dashboard-nba/pkg/role/service"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/pkg/errors"
)

// UpdateRole godoc
// @Summary      Update Role
// @Description  Updates an existing role by its ID. Requires Bearer token.
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        role_id  path      int                       true "Role ID"
// @Param        role     body      presenter.UpdateRoleRequest true "Role data to update"
// @Success      200      {object}  presenter.ResponsePayloadData   "Successfully update role"
// @Failure      400      {object}  presenter.ResponseErrorPayload  "Bad Request (Invalid ID, Invalid JSON, or Validation Error)"
// @Failure      401      {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      404      {object}  presenter.ResponsePayloadMessage "Role ID not found"
// @Failure      500      {object}  presenter.ResponsePayloadMessage "Failed update role"
// @Security     BearerAuth
// @Router       /roles/{role_id} [put]
func UpdateRole(svc service.Service, validate *validator.Validator) fiber.Handler {
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
		var payload rolePresenter.UpdateRoleRequest
		if err = c.BodyParser(&payload); err != nil {
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusBadRequest,
					Message: "Failed parse request",
				})
		}
		if err := validate.Validate(payload); err != nil {
			return presenter.ResponseErrorValidate(c, err)
		}

		ah, err := auth.GetAuth(c)
		if err != nil {
			log.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		userID := ah.GetClaims().UserID

		err = svc.UpdateRoleService(c.UserContext(), payload, userID, roleID)
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
			Message: "Successfully update role",
		})
	}
}
