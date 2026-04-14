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

// CreateMenuPermission godoc
// @Summary      Create Menu Permission
// @Description  Adds a new permission (e.g., 'read', 'create') to a specific menu by its ID. Requires Bearer token.
// @Tags         Menu Permission
// @Accept       json
// @Produce      json
// @Param        menu_id    path      int                                     true "Menu ID"
// @Param        permission body      menuPermissionPresenter.CreateMenuPermissionRequest true "Permission data to create"
// @Success      201        {object}  presenter.ResponsePayloadData           "Successfully create menu permission"
// @Failure      400        {object}  presenter.ResponsePayloadMessage        "Bad Request (Invalid Menu ID, Invalid JSON, or Validation Error)"
// @Failure      401        {object}  presenter.ResponsePayloadMessage        "Unauthorized"
// @Failure      500        {object}  presenter.ResponsePayloadMessage        "Internal Server Error"
// @Security     BearerAuth
// @Router       /menus/{menu_id}/permissions [post]
func CreateMenuPermission(svc contract.MenuPermissionUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		menuIDParams := c.Params("menu_id")
		menuID, err := strconv.Atoi(menuIDParams)
		if err != nil {
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusBadRequest,
					Message: "Failed parse params",
				})
		}

		var payload menuPermissionPresenter.CreateMenuPermissionRequest
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

		err = svc.CreateMenuPermissionUseCase(c.UserContext(), payload, userID, menuID)

		if err != nil {
			if errors.Is(err, constant.ErrMenuIdNotFound) {
				return presenter.ResponseMessage(c,
					presenter.ResponsePayloadMessage{
						Code:    http.StatusNotFound,
						Message: constant.ErrDataNotFound.Message,
					})
			}
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusInternalServerError,
					Message: "Failed create menu permission",
				})
		}

		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusCreated,
			Data:    nil,
			Message: "Successfully create menu permission",
		})
	}
}
