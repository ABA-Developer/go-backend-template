package handlers

import (
	"be-dashboard-nba/internal/presentation/presenter"
	menuPresenter "be-dashboard-nba/internal/presentation/presenter/menu"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// ReadMenuParent godoc
// @Summary      Get Parent Menus
// @Description  Retrieves a list of all root-level menus (menus without a parent), typically for use in a dropdown. Requires Bearer token.
// @Tags         Menu
// @Produce      json
// @Success      200 {object} presenter.ResponsePayloadData{data=[]menuPresenter.MenuParent} "Successfully get parent menu"
// @Failure      401 {object} presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      500 {object} presenter.ResponsePayloadMessage "Failed get parent menu"
// @Security     BearerAuth
// @Router       /menus/parent [get]
func ReadMenuParent(svc contract.MenuUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		data, err := svc.ReadMenuParentUseCase(c.UserContext())
		if err != nil {

			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusInternalServerError,
				Message: "Failed get parent menu",
			})
		}
		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    menuPresenter.ToReadMenuParentResponses(data),
			Message: "Successfully get parent menu",
		})
	}
}
