package app

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/internal/auth"
	"be-dashboard-nba/pkg/auth/service"
)

func logoutApp(svc *service.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		ah, err := auth.GetAuth(c)
		if err != nil {
			return
		}

		err = svc.LogoutService(c.UserContext(), ah.GetClaims())
		if err != nil {
			return presenter.ResponseError(c, err, "Failed logout")
		}

		return presenter.ResponseData(c, presenter.ResponsePayload{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully logout",
		})
	}
}
