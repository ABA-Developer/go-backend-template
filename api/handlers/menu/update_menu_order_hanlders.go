package handlers

import (
	"be-dashboard-nba/api/presenter"
	menuPresenter "be-dashboard-nba/api/presenter/menu"
	"be-dashboard-nba/internal/auth"
	"be-dashboard-nba/internal/validator"
	"be-dashboard-nba/pkg/menu/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// UpdateMenuOrder
// @Summary      Update Menu Order (Reorder)
// @Description  Updates the sort order of menus using drag-and-drop. Requires Bearer token.
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Param        order body      presenter.UpdateMenuOrderRequest true "New menu order"
// @Success      200   {object}  presenter.ResponsePayloadData    "Successfully update menu order"
// @Failure      400   {object}  presenter.ResponsePayloadMessage     "Bad Request (Invalid JSON or Validation Error)"
// @Failure      401   {object}  presenter.ResponsePayloadMessage     "Unauthorized"
// @Failure      500   {object}  presenter.ResponsePayloadMessage     "Internal Server Error"
// @Security     BearerAuth
// @Router       /menus/reorder [put]
func UpdateMenuOrder(svc service.Service, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request menuPresenter.UpdateMenuOrderRequest

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

		err = svc.UpdateMenuOrderService(c.UserContext(), request, userID)
		if err != nil {
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusInternalServerError,
					Message: "Failed update menu order",
				})
		}

		return presenter.ResponseData(c,
			presenter.ResponsePayloadData{
				Code:    http.StatusOK,
				Data:    nil,
				Message: "Successfully update menu order",
			})
	}
}
