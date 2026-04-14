package handlers

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/auth"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/infrastructure/validator"
	"be-dashboard-nba/internal/presentation/presenter"
	menuPermissionPresenter "be-dashboard-nba/internal/presentation/presenter/menu/menu-permission"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// UpdateMenuPermission godoc
// @Summary      Update Menu Permission
// @Description  Updates an existing menu permission (e.g., 'code', 'action_name') by its unique ID. Requires Bearer token.
// @Tags         Menu Permission
// @Accept       json
// @Produce      json
// @Param        menu_id            path      int                                     true "Menu ID"
// @Param        menu_permission_id path      int                                     true "Menu Permission ID"
// @Param        permission         body      menuPermissionPresenter.UpdateMenuPermissionRequest true "Permission data to update"
// @Success      200 {object}  presenter.ResponsePayloadData           "Successfully update menu permission"
// @Failure      400 {object}  presenter.ResponsePayloadMessage        "Bad Request (Invalid ID, Invalid JSON, or Validation Error)"
// @Failure      401 {object}  presenter.ResponsePayloadMessage        "Unauthorized"
// @Failure      404 {object}  presenter.ResponsePayloadMessage        "Not Found (Permission ID not found)"
// @Failure      500 {object}  presenter.ResponsePayloadMessage        "Internal Server Error"
// @Security     BearerAuth
// @Router       /menus/{menu_id}/permissions/{menu_permission_id} [put]
func UpdateMenuPermission(svc contract.MenuPermissionUseCase, validate *validator.Validator) fiber.Handler {
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

		var payload menuPermissionPresenter.UpdateMenuPermissionRequest
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

		err = svc.UpdateMenuPermissionUseCase(c.UserContext(), payload, userID, menuPermissionID)

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
			Message: "Successfully update menu permission",
		})
	}
}
