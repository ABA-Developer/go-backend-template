package routes

import (
	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/api/app"
	"be-dashboard-nba/api/middleware"

	userHandlers "be-dashboard-nba/internal/modules/user/delivery/http"
	userRepo "be-dashboard-nba/internal/modules/user/repository"
	userUsecase "be-dashboard-nba/internal/modules/user/usecase"
)

func UserRouter(http fiber.Router, application *app.Application) {
	repo := userRepo.NewUserRepository(application.DB)
	svc := userUsecase.NewUserUsecase(repo, application.Log)
	mdw := middleware.NewEnsureToken(application.DB)

	routes := http.Group("/users")
	routes.Use(mdw.ValidateToken())

	routes.Get("/", userHandlers.ReadUsers(svc))
	routes.Post("/", userHandlers.CreateUser(svc, application.Validator))
	routes.Put("/:user_id", userHandlers.UpdateUser(svc, application.Validator))
	routes.Delete("/:user_id", userHandlers.DeleteUser(svc))
	routes.Put("/me", userHandlers.UpdateProfileApp(svc, application.Validator))
	routes.Get("/me", userHandlers.ReadProfileApp(svc))
	routes.Get("/:user_id", userHandlers.ReadUserDetail(svc))

}
