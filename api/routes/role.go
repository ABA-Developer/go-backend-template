package routes

import (
	"be-dashboard-nba/api/app"
	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/constant"

	authRepo "be-dashboard-nba/internal/modules/auth/repository"
	authUsecase "be-dashboard-nba/internal/modules/auth/usecase"

	roleHandlers "be-dashboard-nba/internal/modules/role/delivery/http"
	roleRepo "be-dashboard-nba/internal/modules/role/repository"
	roleUsecase "be-dashboard-nba/internal/modules/role/usecase"

	"github.com/gofiber/fiber/v2"
)

func RoleRouter(http fiber.Router, application *app.Application) {
	aRepo := authRepo.NewAuthRepository(application.DB)
	aUsecase := authUsecase.NewAuthUsecase(aRepo, application.Log, application.DB)

	rRepo := roleRepo.NewRoleRepository(application.DB)
	rUsecase := roleUsecase.NewRoleUsecase(rRepo, application.Log)
	mdw := middleware.NewEnsureToken(application.DB)

	routes := http.Group("roles")
	routes.Use(mdw.ValidateToken())

	routes.Get("/", middleware.Authorize(aUsecase, constant.MenuSettingsRole, constant.ActionRead), roleHandlers.ReadRoles(rUsecase))

	routes.Post("/", middleware.Authorize(aUsecase, constant.MenuSettingsRole, constant.ActionCreate), roleHandlers.CreateRole(rUsecase, application.Validator))

	routes.Put("/role-access/:role_id", roleHandlers.UpdateRoleAccess(rUsecase, application.Validator))

	routes.Put("/:role_id", middleware.Authorize(aUsecase, constant.MenuSettingsRole, constant.ActionUpdate), roleHandlers.UpdateRole(rUsecase, application.Validator))

	routes.Delete("/:role_id", middleware.Authorize(aUsecase, constant.MenuSettingsRole, constant.ActionDelete), roleHandlers.DeleteRole(rUsecase))

	routes.Get("/role-access/:role_id", middleware.Authorize(aUsecase, constant.MenuSettingsRole, constant.ActionRead), roleHandlers.ReadRoleAccess(rUsecase))

	routes.Get("/:role_id", middleware.Authorize(aUsecase, constant.MenuSettingsRole, constant.ActionRead), roleHandlers.ReadRoleDetail(rUsecase))
}
