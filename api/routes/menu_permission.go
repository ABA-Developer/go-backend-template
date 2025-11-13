package routes

import (
	"be-dashboard-nba/api/app"
	handlers "be-dashboard-nba/api/handlers/menu_permission"
	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/constant"
	authService "be-dashboard-nba/pkg/auth/service"
	menuPermissionService "be-dashboard-nba/pkg/menu_permission/service"

	"github.com/gofiber/fiber/v2"
)

func MenuPermissionRouter(http fiber.Router, application *app.Application) {
	menuPermissionService := menuPermissionService.NewService(application.DB, application.Log)
	authService := authService.NewService(application.DB, application.Log)
	mdw := middleware.NewEnsureToken(application.DB)

	routes := http.Group("/menu-permissions")
	routes.Use(mdw.ValidateToken())

	routes.Get("/:menu_id", middleware.Authorize(authService, constant.MenuSettingsMenu, constant.ActionReadMenuPermission), handlers.ReadMenuPermissionListParams(menuPermissionService))

	routes.Get("/detail/:menu_permission_id", middleware.Authorize(authService, constant.MenuSettingsMenu, constant.ActionReadMenuPermission), handlers.ReadMenuPermissionDetail(menuPermissionService))

	routes.Post("/:menu_id", middleware.Authorize(authService, constant.MenuSettingsMenu, constant.ActionCreateMenuPermission), handlers.CreateMenuPermission(menuPermissionService, application.Validator))

	routes.Put("/:menu_permission_id", middleware.Authorize(authService, constant.MenuSettingsMenu, constant.ActionUpdateMenuPermission), handlers.UpdateMenuPermission(menuPermissionService, application.Validator))

	routes.Delete("/:menu_permission_id", middleware.Authorize(authService, constant.MenuSettingsMenu, constant.ActionDeleteMenuPermission), handlers.DeleteMenuPermission(menuPermissionService))
}
