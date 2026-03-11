package routes

import (
	"be-dashboard-nba/api/app"
	handlers "be-dashboard-nba/api/handlers/menu"
	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/constant"
	authService "be-dashboard-nba/usecase/auth"
	menuService "be-dashboard-nba/usecase/menu"

	"github.com/gofiber/fiber/v2"
)

func MenuRouter(http fiber.Router, application *app.Application) {
	authSvc := authService.NewService(application.DB, application.Log)
	menuSvc := menuService.NewService(application.DB, application.Log)
	mdw := middleware.NewEnsureToken(application.DB)

	routes := http.Group("/menus")
	routes.Use(mdw.ValidateToken())

	routes.Get("/", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionRead), handlers.ReadListMenu(menuSvc))

	routes.Post("/", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionCreate), handlers.CreateMenu(menuSvc, application.Validator))

	routes.Get("/sidebar", handlers.ReadMenuSidebar(menuSvc))

	routes.Get("/parent", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionRead), handlers.ReadMenuParent(menuSvc))

	routes.Put("/reorder", handlers.UpdateMenuOrder(menuSvc, application.Validator))

	routes.Get("/:menu_id", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionRead), handlers.ReadMenuDetail(menuSvc))

	routes.Delete("/:menu_id", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionDelete), handlers.DeleteMenu(menuSvc))

	routes.Put("/:menu_id", handlers.UpdateMenu(menuSvc, application.Validator))
}
