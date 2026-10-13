package routes

import (
	"github.com/gofiber/fiber/v2"

	usecase "github.com/ABA-Developer/go-backend-template/internal/application/user/usecase"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/runtime/container"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/middleware"
	handlers "github.com/ABA-Developer/go-backend-template/internal/presentation/user/handler"
)

func AddUserRoutes(http fiber.Router, c *container.Container) {
	svc := usecase.NewUseCase(c.GetDB())
	mdw := middleware.NewEnsureToken(c.GetDB())

	routes := http.Group("/users")
	routes.Use(mdw.ValidateToken())

	routes.Get("/", handlers.ReadUsers(svc))
	routes.Post("/", handlers.CreateUser(svc, c.GetValidator()))
	routes.Put("/:user_id", handlers.UpdateUser(svc, c.GetValidator()))
	routes.Delete("/:user_id", handlers.DeleteUser(svc))
	routes.Put("/me", handlers.UpdateProfileApp(svc, c.GetValidator()))
	routes.Get("/me", handlers.ReadProfileApp(svc))
	routes.Get("/:user_id", handlers.ReadUserDetail(svc))

}
