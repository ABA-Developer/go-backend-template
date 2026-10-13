package handler

import (
	contract "github.com/ABA-Developer/go-backend-template/internal/domain/contract/usecase"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/logger"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/menu/mapper"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/menu/presenter"
	response "github.com/ABA-Developer/go-backend-template/internal/presentation/response"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// @Summary      Get List of All Menus (Admin)
// @Description  Retrieves a paginated and searchable list of all menus in the system.
// @Tags         Menu
// @Produce      json
// @Param        search query    string false "Search term for name, description, or URL"
// @Success      200  {object}  response.DataPayload{data=presenter.ReadMenuListResponse} "Successfully get menu list"
// @Failure      400  {object}  response.ErrorPayload "Invalid query parameters"
// @Failure      500  {object}  response.ErrorPayload "Failed get list menu"
// @Router       /menus [get]
func ReadListMenu(svc contract.MenuUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request presenter.ReadMenuListRequest
		if err = c.QueryParser(&request); err != nil {
			logger.Errorw("error parse request", err)
			return response.Error(c, response.ErrorParam{
				Code:    http.StatusInternalServerError,
				Message: "Invalid query parameters",
			})
		}

		data, err := svc.ReadListMenuUseCase(c.UserContext(), mapper.ToReadMenuListPayload(request))
		if err != nil {
			return response.Error(c, response.ErrorParam{
				Code:    http.StatusInternalServerError,
				Message: "Failed get list menu",
			})
		}

		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    presenter.ToReadMenuListResponse(data),
			Message: "Successfully get menu list",
		})
	}
}
