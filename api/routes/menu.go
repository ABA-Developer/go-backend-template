package routes

import (
	"be-dashboard-nba/api/app"
	handlers "be-dashboard-nba/api/handlers/menu"
	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/constant"
	authService "be-dashboard-nba/pkg/auth/service"
	menuService "be-dashboard-nba/pkg/menu/service"

	"github.com/gofiber/fiber/v2"
)

func MenuRouter(http fiber.Router, application *app.Application) {
	authService := authService.NewService(application.DB, application.Log)
	menuService := menuService.NewService(application.DB, application.Log)
	mdw := middleware.NewEnsureToken(application.DB)

	routes := http.Group("/menus")
	routes.Use(mdw.ValidateToken())

	routes.Get("/", middleware.Authorize(authService, constant.MenuSettingsMenu, constant.ActionRead), handlers.ReadListMenu(menuService))

	routes.Post("/", middleware.Authorize(authService, constant.MenuSettingsMenu, constant.ActionCreate), handlers.CreateMenu(menuService, application.Validator))

	routes.Get("/sidebar", handlers.ReadMenuSidebar(menuService))

	routes.Get("/parent", middleware.Authorize(authService, constant.MenuSettingsMenu, constant.ActionRead), handlers.ReadMenuParent(menuService))

	routes.Put("/reorder", middleware.Authorize(authService, constant.MenuSettingsMenu, constant.ActionUpdate), handlers.UpdateMenuOrder(menuService, application.Validator))

	routes.Get("/:menu_id", middleware.Authorize(authService, constant.MenuSettingsMenu, constant.ActionRead), handlers.ReadMenuDetail(menuService))

	routes.Delete("/:menu_id", middleware.Authorize(authService, constant.MenuSettingsMenu, constant.ActionDelete), handlers.DeleteMenu(menuService))

	routes.Put("/:menu_id", middleware.Authorize(authService, constant.MenuSettingsMenu, constant.ActionUpdate), handlers.UpdateMenu(menuService, application.Validator))
}
