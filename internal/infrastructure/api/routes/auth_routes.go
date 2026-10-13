package routes

import (
	"github.com/gofiber/fiber/v2"

	usecase "github.com/ABA-Developer/go-backend-template/internal/application/auth/usecase"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/runtime/container"
	handlers "github.com/ABA-Developer/go-backend-template/internal/presentation/auth/handler"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/middleware"
)

func AddAuthRoutes(http fiber.Router, c *container.Container) {
	svc := usecase.NewUseCase(c.GetDB())
	mdw := middleware.NewEnsureToken(c.GetDB())

	routes := http.Group("/auth")

	routes.Post("/login", handlers.Login(svc, c.GetValidator()))
	routes.Post("/logout", mdw.ValidateToken(), handlers.Logout(svc))
	routes.Get("/me", mdw.ValidateToken(), handlers.AuthMe(svc))
}
