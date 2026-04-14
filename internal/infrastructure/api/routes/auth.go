package routes

import (
	app "be-dashboard-nba/internal/infrastructure/api"
	handlers "be-dashboard-nba/internal/presentation/handler/auth"
	"be-dashboard-nba/internal/presentation/middleware"
	"be-dashboard-nba/internal/application/usecase/auth"

	"github.com/gofiber/fiber/v2"
)

func AuthRouter(http fiber.Router, application *app.Application) {
	svc := auth.NewUseCase(application.DB, application.Log)
	mdw := middleware.NewEnsureToken(application.DB)

	routes := http.Group("/auth")

	routes.Post("/login", handlers.Login(svc, application.Validator))
	routes.Post("/logout", mdw.ValidateToken(), handlers.Logout(svc))
	routes.Get("/me", mdw.ValidateToken(), handlers.AuthMe(svc))
}
