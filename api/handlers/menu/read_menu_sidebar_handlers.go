package handlers

import (
	"be-dashboard-nba/api/presenter"
	menuPresenter "be-dashboard-nba/api/presenter/menu"
	"be-dashboard-nba/internal/auth"
	"be-dashboard-nba/pkg/menu/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// ReadMenuSidebar godoc
// @Summary      Get User's Sidebar Menu
// @Description  Retrieves the menu list specifically for the authenticated user, filtered by their 'read' permissions and formatted in a nested tree structure for the sidebar.
// @Tags         Menu
// @Produce      json
// @Success      200  {object}  presenter.ResponsePayloadData{data=[]presenter.ReadMenuListResponse} "Successfully get sidebar menu"
// @Failure      401  {object}  presenter.ResponsePayloadMessage "Unauthorized"
// @Failure      500  {object}  presenter.ResponsePayloadMessage "Failed get sidebar menu"
// @Security     BearerAuth
// @Router       /menus/sidebar [get]
func ReadMenuSidebar(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		ah, err := auth.GetAuth(c)
		if err != nil {
			log.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		userID := ah.GetClaims().UserID

		data, err := svc.ReadSidebarMenuService(c.UserContext(), userID)
		if err != nil {

			return presenter.ResponseMessage(c, presenter.ResponsePayloadMessage{
				Code:    http.StatusInternalServerError,
				Message: "Failed get sidebar menu",
			})
		}
		return presenter.ResponseData(c, presenter.ResponsePayloadData{
			Code:    http.StatusOK,
			Data:    menuPresenter.ToReadMenuListResponse(data),
			Message: "Successfully get sidebar menu",
		})
	}
}
