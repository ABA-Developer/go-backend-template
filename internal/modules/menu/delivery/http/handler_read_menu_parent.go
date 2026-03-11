package http

import (
	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/internal/modules/menu/domain"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// ReadMenuParent godoc
// @Summary      Get Parent Menus
// @Description  Retrieves a list of all root-level menus (menus without a parent), typically for use in a dropdown. Requires Bearer token.
// @Tags         Menu
// @Produce      json
// @Success      200 {object} presenter.ResponsePayloadData{data=[]presenter.MenuParent} "Successfully get parent menu"
// @Failure      401 {object} presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      500 {object} presenter.ResponsePayloadMessage "Failed get parent menu"
// @Security     BearerAuth
// @Router       /menus/parent [get]
func ReadMenuParent(usecase domain.MenuUsecase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		data, err := usecase.ReadMenuParentUsecase(c.UserContext())
		if err != nil {

			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusInternalServerError,
				Message: "Failed get parent menu",
			})
		}
		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    data,
			Message: "Successfully get parent menu",
		})
	}
}
