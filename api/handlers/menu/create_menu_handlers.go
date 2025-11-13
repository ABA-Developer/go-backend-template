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

// CreateMenu godoc
// @Summary      Create Menu
// @Description  Adds a new menu item to the system. Requires Bearer token.
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Param        menu body     presenter.CreateMenuRequest true "Menu data to create"
// @Success      201  {object} presenter.ResponsePayloadData         "Successfully create menu"
// @Failure      400  {object} presenter.ResponsePayloadMessage    "Bad Request (Invalid JSON or Validation Error)"
// @Failure      401  {object} presenter.ResponsePayloadMessage    "Unauthorized"
// @Failure      500  {object} presenter.ResponsePayloadMessage    "Internal Server Error"
// @Security     BearerAuth
// @Router       /menus [post]
func CreateMenu(svc service.Service, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request menuPresenter.CreateMenuRequest
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

		err = svc.CreateMenuService(c.UserContext(), request, userID)
		if err != nil {
			return presenter.ResponseMessage(c,
				presenter.ResponsePayloadMessage{
					Code:    http.StatusInternalServerError,
					Message: "Failed create menu",
				})
		}
		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusCreated,
			Data:    nil,
			Message: "Successfully create menu",
		})
	}
}
