package handler

import (
	contract "github.com/ABA-Developer/go-backend-template/internal/domain/contract/usecase"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/logger"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/validator"
	authInternal "github.com/ABA-Developer/go-backend-template/internal/presentation/auth"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/menu/mapper"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/menu/presenter"
	response "github.com/ABA-Developer/go-backend-template/internal/presentation/response"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// UpdateMenuOrder
// @Summary      Update Menu Order (Reorder)
// @Description  Updates the sort order of menus using drag-and-drop. Requires Bearer token.
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Param        order body      presenter.UpdateMenuOrderRequest true "New menu order"
// @Success      200   {object}  response.DataPayload    "Successfully update menu order"
// @Failure      400   {object}  response.ErrorPayload     "Bad Request (Invalid JSON or Validation Error)"
// @Failure      401   {object}  response.ErrorPayload     "Unauthorized"
// @Failure      500   {object}  response.ErrorPayload     "Internal Server Error"
// @Security     BearerAuth
// @Router       /menus/reorder [put]
func UpdateMenuOrder(svc contract.MenuUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request presenter.UpdateMenuOrderRequest

		if err = c.BodyParser(&request); err != nil {
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: "Failed parse request",
				})
		}

		if err := validate.Validate(request); err != nil {
			return response.ErrorValidate(c, err)
		}

		ah, err := authInternal.GetAuth(c)
		if err != nil {
			logger.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}
		userID := ah.GetClaims().UserID

		err = svc.UpdateMenuOrderUseCase(c.UserContext(), mapper.ToUpdateMenuOrderPayload(request), userID)
		if err != nil {
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusInternalServerError,
					Message: "Failed update menu order",
				})
		}

		return response.Data(c,
			response.DataPayload{
				Code:    http.StatusOK,
				Data:    nil,
				Message: "Successfully update menu order",
			})
	}
}
