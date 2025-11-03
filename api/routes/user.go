package routes

import (
	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/api/app"
	"be-dashboard-nba/api/handlers"
	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/pkg/user/service"
)

func UserRouter(http fiber.Router, application *app.Application) {

	svc := service.NewService(application.DB, application.Log)
	mdw := middleware.NewEnsureToken(application.DB)

	routes := http.Group("/users")
	routes.Put("/me", mdw.ValidateToken(), handlers.UpdateProfileApp(svc, application.Validator))
	routes.Get("/me", mdw.ValidateToken(), handlers.ReadProfileApp(svc))

	routes.Put("/:id", mdw.ValidateToken(), handlers.UpdateUserApp(svc, application.Validator))

}
