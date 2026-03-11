package routes

import (
	"be-dashboard-nba/api/app"
	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/constant"
	authRepo "be-dashboard-nba/internal/modules/auth/repository"
	authUsecase "be-dashboard-nba/internal/modules/auth/usecase"
	handlers "be-dashboard-nba/internal/modules/menu/delivery/http"
	menuRepo "be-dashboard-nba/internal/modules/menu/repository"
	menuUsecase "be-dashboard-nba/internal/modules/menu/usecase"

	"github.com/gofiber/fiber/v2"
)

func MenuRouter(http fiber.Router, application *app.Application) {
	aRepo := authRepo.NewAuthRepository(application.DB)
	aUsecase := authUsecase.NewAuthUsecase(aRepo, application.Log, application.DB)
	mRepo := menuRepo.NewMenuRepository(application.DB)
	mUsecase := menuUsecase.NewMenuUsecase(mRepo, application.Log)
	mdw := middleware.NewEnsureToken(application.DB)

	routes := http.Group("/menus")
	routes.Use(mdw.ValidateToken())

	routes.Get("/", middleware.Authorize(aUsecase, constant.MenuSettingsMenu, constant.ActionRead), handlers.ReadListMenu(mUsecase))

	routes.Post("/", middleware.Authorize(aUsecase, constant.MenuSettingsMenu, constant.ActionCreate), handlers.CreateMenu(mUsecase, application.Validator))

	routes.Get("/sidebar", handlers.ReadMenuSidebar(mUsecase))

	routes.Get("/parent", middleware.Authorize(aUsecase, constant.MenuSettingsMenu, constant.ActionRead), handlers.ReadMenuParent(mUsecase))

	routes.Put("/reorder", middleware.Authorize(aUsecase, constant.MenuSettingsMenu, constant.ActionUpdate), handlers.UpdateMenuOrder(mUsecase, application.Validator))

	routes.Get("/:menu_id", middleware.Authorize(aUsecase, constant.MenuSettingsMenu, constant.ActionRead), handlers.ReadMenuDetail(mUsecase))

	routes.Delete("/:menu_id", middleware.Authorize(aUsecase, constant.MenuSettingsMenu, constant.ActionDelete), handlers.DeleteMenu(mUsecase))

	routes.Put("/:menu_id", middleware.Authorize(aUsecase, constant.MenuSettingsMenu, constant.ActionUpdate), handlers.UpdateMenu(mUsecase, application.Validator))
}
