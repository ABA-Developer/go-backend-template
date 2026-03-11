package routes

import (
	"be-dashboard-nba/api/app"
	handlers "be-dashboard-nba/api/handlers/role"
	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/constant"
	authService "be-dashboard-nba/usecase/auth"
	roleService "be-dashboard-nba/usecase/role"

	"github.com/gofiber/fiber/v2"
)

func RoleRouter(http fiber.Router, application *app.Application) {
	authSvc := authService.NewService(application.DB, application.Log)
	roleSvc := roleService.NewService(application.DB, application.Log)
	mdw := middleware.NewEnsureToken(application.DB)

	routes := http.Group("roles")
	routes.Use(mdw.ValidateToken())

	routes.Get("/", middleware.Authorize(authSvc, constant.MenuSettingsRole, constant.ActionRead), handlers.ReadRoles(roleSvc))

	routes.Post("/", middleware.Authorize(authSvc, constant.MenuSettingsRole, constant.ActionCreate), handlers.CreateRole(roleSvc, application.Validator))

	routes.Put("/role-access/:role_id", handlers.UpdateRoleMenuPermission(roleSvc, application.Validator))

	routes.Put("/:role_id", middleware.Authorize(authSvc, constant.MenuSettingsRole, constant.ActionUpdate), handlers.UpdateRole(roleSvc, application.Validator))

	routes.Delete("/:role_id", middleware.Authorize(authSvc, constant.MenuSettingsRole, constant.ActionDelete), handlers.DeleteRole(roleSvc))

	routes.Get("/role-access/:role_id", middleware.Authorize(authSvc, constant.MenuSettingsRole, constant.ActionRead), handlers.ReadRoleAccess(roleSvc))

	routes.Get("/:role_id", middleware.Authorize(authSvc, constant.MenuSettingsRole, constant.ActionRead), handlers.ReadRoleDetail(roleSvc))
}
