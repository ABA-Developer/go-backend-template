package handlers

import (
	"be-dashboard-nba/internal/presentation/presenter"
	menuPresenter "be-dashboard-nba/internal/presentation/presenter/menu"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/auth"
	"be-dashboard-nba/internal/infrastructure/validator"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// UpdateMenu godoc
// @Summary      Update Menu
// @Description  Updates an existing menu item by its ID. Requires Bearer token.
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Param        menu_id  path      int                             true "Menu ID"
// @Param        menu     body      menuPresenter.UpdateMenuRequest true "Menu data to update"
// @Success      200      {object}  presenter.ResponsePayloadData       "Successfully update menu"
// @Failure      400      {object}  presenter.ResponsePayloadMessage  "Bad Request (Invalid Menu ID, Invalid JSON, or Validation Error)"
// @Failure      401      {object}  presenter.ResponsePayloadMessage  "Unauthorized"
// @Failure      500      {object}  presenter.ResponsePayloadMessage  "Internal Server Error"
// @Security     BearerAuth
// @Router       /menus/{menu_id} [put]
func UpdateMenu(svc contract.MenuUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request menuPresenter.UpdateMenuRequest

		menuIDParams := c.Params("menu_id")
		menuID, err := strconv.Atoi(menuIDParams)
		if err != nil {
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusBadRequest,
					Message: "menu id params must be in number",
				})
		}
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

		ah, err := auth.GetAuth(c)
		if err != nil {
			log.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}
		userID := ah.GetClaims().UserID

		err = svc.UpdateMenuUseCase(c.UserContext(), request, userID, menuID)
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
					Message: "Failed update menu",
				})
		}

		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully update menu",
		})
	}
}
