package routes

import (
	"github.com/ABA-Developer/go-backend-template/constant"
	authUseCase "github.com/ABA-Developer/go-backend-template/internal/application/auth/usecase"
	menuPermissionUseCase "github.com/ABA-Developer/go-backend-template/internal/application/menu/menu-permission/usecase"
	menuUseCase "github.com/ABA-Developer/go-backend-template/internal/application/menu/usecase"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/runtime/container"
	handlers "github.com/ABA-Developer/go-backend-template/internal/presentation/menu/handler"
	menuPermissionHandlers "github.com/ABA-Developer/go-backend-template/internal/presentation/menu/menu-permission/handler"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/middleware"

	"github.com/gofiber/fiber/v2"
)

func AddMenuRoutes(http fiber.Router, c *container.Container) {
	authSvc := authUseCase.NewUseCase(c.GetDB())
	menuSvc := menuUseCase.NewUseCase(c.GetDB())
	mpSvc := menuPermissionUseCase.NewUseCase(c.GetDB())
	mdw := middleware.NewEnsureToken(c.GetDB())

	routes := http.Group("/menus")
	routes.Use(mdw.ValidateToken())

	routes.Get("/", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionRead), handlers.ReadListMenu(menuSvc))

	routes.Post("/", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionCreate), handlers.CreateMenu(menuSvc, c.GetValidator()))

	routes.Get("/sidebar", handlers.ReadMenuSidebar(menuSvc))

	routes.Get("/parent", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionRead), handlers.ReadMenuParent(menuSvc))

	routes.Put("/reorder", handlers.UpdateMenuOrder(menuSvc, c.GetValidator()))

	routes.Get("/:menu_id", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionRead), handlers.ReadMenuDetail(menuSvc))

	routes.Delete("/:menu_id", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionDelete), handlers.DeleteMenu(menuSvc))

	routes.Put("/:menu_id", handlers.UpdateMenu(menuSvc, c.GetValidator()))

	permissionRoutes := routes.Group("/:menu_id/permissions")
	permissionRoutes.Get("/", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionReadMenuPermission), menuPermissionHandlers.ReadMenuPermissionUseCase(mpSvc))
	permissionRoutes.Post("/", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionCreateMenuPermission), menuPermissionHandlers.CreateMenuPermission(mpSvc, c.GetValidator()))
	permissionRoutes.Get("/:menu_permission_id", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionReadMenuPermission), menuPermissionHandlers.ReadMenuPermissionDetailUseCase(mpSvc))
	permissionRoutes.Put("/:menu_permission_id", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionUpdateMenuPermission), menuPermissionHandlers.UpdateMenuPermission(mpSvc, c.GetValidator()))
	permissionRoutes.Delete("/:menu_permission_id", middleware.Authorize(authSvc, constant.MenuSettingsMenu, constant.ActionDeleteMenuPermission), menuPermissionHandlers.DeleteMenuPermission(mpSvc))
}
