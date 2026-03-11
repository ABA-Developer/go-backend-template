package routes

import (
	"be-dashboard-nba/api/app"
	handlers "be-dashboard-nba/api/handlers/menu_permission"
	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/constant"
	authService "be-dashboard-nba/usecase/auth"
	menuPermissionService "be-dashboard-nba/usecase/menu_permission"

	"github.com/gofiber/fiber/v2"
)

func MenuPermissionRouter(http fiber.Router, application *app.Application) {
	mpSvc := menuPermissionService.NewService(application.DB, application.Log)
	authSvc := authService.NewService(application.DB, application.Log)
	mdw := middleware.NewEnsureToken(application.DB)

	routes := http.Group("/menu-permissions")
	routes.Use(mdw.ValidateToken())

	routes.Get("/:menu_id", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionReadMenuPermission), handlers.ReadMenuPermissionService(mpSvc))

	routes.Get("/detail/:menu_permission_id", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionReadMenuPermission), handlers.ReadMenuPermissionDetail(mpSvc))

	routes.Post("/:menu_id", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionCreateMenuPermission), handlers.CreateMenuPermission(mpSvc, application.Validator))

	routes.Put("/:menu_permission_id", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionUpdateMenuPermission), handlers.UpdateMenuPermission(mpSvc, application.Validator))

	routes.Delete("/:menu_permission_id", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionDeleteMenuPermission), handlers.DeleteMenuPermission(mpSvc))
}
