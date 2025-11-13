package handlers

import (
	"be-dashboard-nba/api/presenter"
	menuPermissionPresenter "be-dashboard-nba/api/presenter/menu_permission"
	"be-dashboard-nba/internal/auth"
	"be-dashboard-nba/internal/validator"
	"be-dashboard-nba/pkg/menu_permission/service"
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
// @Param        permission body      presenter.CreateMenuPermissionRequest true "Permission data to create"
// @Success      201        {object}  presenter.ResponsePayloadData           "Successfully create menu permission"
// @Failure      400        {object}  presenter.ResponsePayloadMessage        "Bad Request (Invalid Menu ID, Invalid JSON, or Validation Error)"
// @Failure      401        {object}  presenter.ResponsePayloadMessage        "Unauthorized"
// @Failure      500        {object}  presenter.ResponsePayloadMessage        "Internal Server Error"
// @Security     BearerAuth
// @Router       /menu-permissions/{menu_id} [post]
func CreateMenuPermission(svc service.Service, validate *validator.Validator) fiber.Handler {
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

		err = svc.CreateMenuPermissionService(c.UserContext(), payload, userID, menuID)

		if err != nil {
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
