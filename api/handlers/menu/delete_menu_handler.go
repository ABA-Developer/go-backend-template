package handlers

import (
	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/menu/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// DeleteMenu godoc
// @Summary      Update Menu
// @Description  Updates an existing menu item by its ID. Requires Bearer token.
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Param        menu_id  path      int                             true "Menu ID"
// @Success      200      {object}  presenter.ResponsePayloadData       "Successfully delete menu"
// @Failure      400      {object}  presenter.ResponsePayloadMessage  "Bad Request (Invalid Menu ID, Invalid JSON, or Validation Error)"
// @Failure      401      {object}  presenter.ResponsePayloadMessage  "Unauthorized"
// @Failure      500      {object}  presenter.ResponsePayloadMessage  "Internal Server Error"
// @Security     BearerAuth
// @Router       /menus/{menu_id} [delete]
func DeleteMenu(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		menuIDParams := c.Params("menu_id")
		menuID, err := strconv.Atoi(menuIDParams)
		if err != nil {
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusBadRequest,
					Message: "menu id params must be in number",
				})
		}

		err = svc.DeleteMenuService(c.UserContext(), menuID)
		if err != nil {
			if errors.Is(err, constant.ErrMenuIdNotFound) {
				return presenter.ResponseMessage(c,
					presenter.ResponsePayloadMessage{
						Code:    http.StatusNotFound,
						Message: constant.ErrDataNotFound.Message,
					})
			}
		}

		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully delete menu",
		})

	}
}
