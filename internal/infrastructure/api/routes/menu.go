package routes

import (
	"be-dashboard-nba/constant"
	authUseCase "be-dashboard-nba/internal/application/usecase/auth"
	menuUseCase "be-dashboard-nba/internal/application/usecase/menu"
	menuPermissionUseCase "be-dashboard-nba/internal/application/usecase/menu/menu-permission"
	app "be-dashboard-nba/internal/infrastructure/api"
	handlers "be-dashboard-nba/internal/presentation/handler/menu"
	menuPermissionHandlers "be-dashboard-nba/internal/presentation/handler/menu/menu-permission"
	"be-dashboard-nba/internal/presentation/middleware"

	"github.com/gofiber/fiber/v2"
)

func MenuRouter(http fiber.Router, application *app.Application) {
	authSvc := authUseCase.NewUseCase(application.DB, application.Log)
	menuSvc := menuUseCase.NewUseCase(application.DB, application.Log)
	mpSvc := menuPermissionUseCase.NewUseCase(application.DB, application.Log)
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

	permissionRoutes := routes.Group("/:menu_id/permissions")
	permissionRoutes.Get("/", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionReadMenuPermission), menuPermissionHandlers.ReadMenuPermissionUseCase(mpSvc))
	permissionRoutes.Post("/", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionCreateMenuPermission), menuPermissionHandlers.CreateMenuPermission(mpSvc, application.Validator))
	permissionRoutes.Get("/:menu_permission_id", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionReadMenuPermission), menuPermissionHandlers.ReadMenuPermissionDetailUseCase(mpSvc))
	permissionRoutes.Put("/:menu_permission_id", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionUpdateMenuPermission), menuPermissionHandlers.UpdateMenuPermission(mpSvc, application.Validator))
	permissionRoutes.Delete("/:menu_permission_id", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionDeleteMenuPermission), menuPermissionHandlers.DeleteMenuPermission(mpSvc))
}
