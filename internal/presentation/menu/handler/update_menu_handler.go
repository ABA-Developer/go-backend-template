package handler

import (
	"github.com/ABA-Developer/go-backend-template/constant"
	contract "github.com/ABA-Developer/go-backend-template/internal/domain/contract/usecase"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/logger"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/validator"
	authInternal "github.com/ABA-Developer/go-backend-template/internal/presentation/auth"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/menu/mapper"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/menu/presenter"
	response "github.com/ABA-Developer/go-backend-template/internal/presentation/response"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

// @Summary      Update Menu
// @Description  Updates an existing menu item by its ID. Requires Bearer token.
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Param        menu_id  path      int                             true "Menu ID"
// @Param        menu     body      presenter.UpdateMenuRequest true "Menu data to update"
// @Success      200      {object}  response.DataPayload       "Successfully update menu"
// @Failure      400      {object}  response.ErrorPayload  "Bad Request (Invalid Menu ID, Invalid JSON, or Validation Error)"
// @Failure      401      {object}  response.ErrorPayload  "Unauthorized"
// @Failure      500      {object}  response.ErrorPayload  "Internal Server Error"
// @Security     BearerAuth
// @Router       /menus/{menu_id} [put]
func UpdateMenu(svc contract.MenuUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request presenter.UpdateMenuRequest

		menuIDParams := c.Params("menu_id")
		menuID, err := strconv.Atoi(menuIDParams)
		if err != nil {
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: "menu id params must be in number",
				})
		}
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

		err = svc.UpdateMenuUseCase(c.UserContext(), mapper.ToUpdateMenuPayload(request), userID, menuID)
		if err != nil {
			if errors.Is(err, constant.ErrMenuIdNotFound) {
				return response.Error(c,
					response.ErrorParam{
						Code:    http.StatusNotFound,
						Message: constant.ErrDataNotFound.Message,
					})
			}
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusInternalServerError,
					Message: "Failed update menu",
				})
		}

		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully update menu",
		})
	}
}
