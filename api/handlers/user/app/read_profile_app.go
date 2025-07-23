package app

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"be-dashboard-nba/api/handlers/user/payload"
	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/internal/auth"
	"be-dashboard-nba/pkg/user/service"
)

func readProfileApp(svc *service.Service) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		ah, err := auth.GetAuth(c)
		if err != nil {
			log.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		data, err := svc.ReadDetailUserService(c.UserContext(), ah.GetClaims().UserID)
		if err != nil {
			return presenter.ResponseError(c, err, "Failed read profile user")
		}

		return presenter.ResponseData(c, presenter.ResponsePayload{
			Code:    http.StatusOK,
			Data:    payload.ToReadUserResponse(data),
			Message: "Successfully read profile user",
		})
	}
}
