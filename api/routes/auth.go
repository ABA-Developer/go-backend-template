package routes

import (
	"be-dashboard-nba/api/app"
	"be-dashboard-nba/api/middleware"

	authHandlers "be-dashboard-nba/internal/modules/auth/delivery/http"
	authRepo "be-dashboard-nba/internal/modules/auth/repository"
	authUsecase "be-dashboard-nba/internal/modules/auth/usecase"

	"github.com/gofiber/fiber/v2"
)

func AuthRouter(http fiber.Router, application *app.Application) {
	repo := authRepo.NewAuthRepository(application.DB)
	svc := authUsecase.NewAuthUsecase(repo, application.Log, application.DB)

	mdw := middleware.NewEnsureToken(application.DB)

	routes := http.Group("/auth")

	routes.Post("/login", authHandlers.Login(svc, application.Validator))
	routes.Post("/logout", mdw.ValidateToken(), authHandlers.Logout(svc))
	routes.Get("/me", mdw.ValidateToken(), authHandlers.AuthMe(svc))
}
