package http

import (
	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/internal/modules/menu/domain"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// ReadListMenu godoc
// @Summary      Get List of All Menus (Admin)
// @Description  Retrieves a paginated and searchable list of all menus in the system.
// @Tags         Menu
// @Produce      json
// @Param        search query    string false "Search term for name, description, or URL"
// @Success      200  {object}  presenter.ResponsePayloadData{data=presenter.ReadMenuListResponse} "Successfully get menu list"
// @Failure      400  {object}  presenter.ResponsePayloadMessage "Invalid query parameters"
// @Failure      500  {object}  presenter.ResponsePayloadMessage "Failed get list menu"
// @Router       /menus [get]
func ReadListMenu(usecase domain.MenuUsecase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request ReadMenuRequest
		if err = c.QueryParser(&request); err != nil {
			log.Errorw("error parse request", err)
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusInternalServerError,
				Message: "Invalid query parameters",
			})
		}
		filter := request.ToDomainFilter()
		data, err := usecase.ReadListMenuUsecase(c.UserContext(), filter)
		if err != nil {
			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusInternalServerError,
				Message: "Failed get list menu",
			})
		}

		if data == nil {
			data = []domain.Menu{}
		}

		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    ToReadMenuListResponse(data),
			Message: "Successfully get menu list",
		})
	}
}
