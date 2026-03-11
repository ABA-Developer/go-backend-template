package routes

import (
	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/api/app"
	handlers "be-dashboard-nba/api/handlers/user"
	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/usecase/user"
)

func UserRouter(http fiber.Router, application *app.Application) {

	svc := user.NewService(application.DB, application.Log)
	mdw := middleware.NewEnsureToken(application.DB)

	routes := http.Group("/users")
	routes.Use(mdw.ValidateToken())

	routes.Get("/", handlers.ReadUsers(svc))
	routes.Post("/", handlers.CreateUser(svc, application.Validator))
	routes.Put("/:user_id", handlers.UpdateUser(svc, application.Validator))
	routes.Delete("/:user_id", handlers.DeleteUser(svc))
	routes.Put("/me", handlers.UpdateProfileApp(svc, application.Validator))
	routes.Get("/me", handlers.ReadProfileApp(svc))
	routes.Get("/:user_id", handlers.ReadUserDetail(svc))

}
