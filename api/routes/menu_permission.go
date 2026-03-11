package routes

import (
	"be-dashboard-nba/api/app"
	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/constant"
	authRepo "be-dashboard-nba/internal/modules/auth/repository"
	authUsecase "be-dashboard-nba/internal/modules/auth/usecase"
	menuRepo "be-dashboard-nba/internal/modules/menu/repository"
	handlers "be-dashboard-nba/internal/modules/menu_permission/delivery/http"
	menuPermissionRepo "be-dashboard-nba/internal/modules/menu_permission/repository"
	menuPermissionUsecase "be-dashboard-nba/internal/modules/menu_permission/usecase"

	"github.com/gofiber/fiber/v2"
)

func MenuPermissionRouter(http fiber.Router, application *app.Application) {
	aRepo := authRepo.NewAuthRepository(application.DB)
	mRepo := menuRepo.NewMenuRepository(application.DB)
	mprRepo := menuPermissionRepo.NewMenuRepository(application.DB)

	aUsecase := authUsecase.NewAuthUsecase(aRepo, application.Log, application.DB)
	mprUsecase := menuPermissionUsecase.NewMenuPermissionUsecase(mprRepo, mRepo, application.Log)
	mdw := middleware.NewEnsureToken(application.DB)

	routes := http.Group("/menu-permissions")
	routes.Use(mdw.ValidateToken())

	routes.Post("/:menu_id", middleware.Authorize(aUsecase, constant.MenuSettingsMenu, constant.ActionCreateMenuPermission), handlers.CreateMenuPermission(mprUsecase, application.Validator))

	routes.Get("/:menu_id", middleware.Authorize(aUsecase, constant.MenuSettingsMenu, constant.ActionReadMenuPermission), handlers.ReadMenuPermissionService(mprUsecase))

	routes.Put("/:menu_permission_id", middleware.Authorize(aUsecase, constant.MenuSettingsMenu, constant.ActionUpdateMenuPermission), handlers.UpdateMenuPermission(mprUsecase, application.Validator))

	routes.Get("/detail/:menu_permission_id", middleware.Authorize(aUsecase, constant.MenuSettingsMenu, constant.ActionReadMenuPermission), handlers.ReadMenuPermissionDetail(mprUsecase))

	routes.Delete("/:menu_permission_id", middleware.Authorize(aUsecase, constant.MenuSettingsMenu, constant.ActionDeleteMenuPermission), handlers.DeleteMenuPermission(mprUsecase))

}
