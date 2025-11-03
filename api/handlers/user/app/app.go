package app

import (
	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/internal/validator"
	"be-dashboard-nba/pkg/user/service"
)

func UserRouter(app fiber.Router, db db.DB, validate *validator.Validator) {
	svc := service.NewService(db)
	mdw := middleware.NewEnsureToken(db)

	app.Use(mdw.ValidateToken())
	app.Get("", readListUserApp(svc, validate))
	app.Get("/me", readProfileApp(svc))
	app.Get("/:id", readDetailUserApp(svc))
	app.Post("", createUserApp(svc, validate))
	app.Put("/:id", updateUserApp(svc, validate))
	app.Delete("/:id", deleteUserApp(svc))
}
