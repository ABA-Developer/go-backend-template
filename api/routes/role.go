package routes

import (
	"be-dashboard-nba/api/app"
	handlers "be-dashboard-nba/api/handlers/role"
	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/constant"
	authService "be-dashboard-nba/pkg/auth/service"
	roleService "be-dashboard-nba/pkg/role/service"

	"github.com/gofiber/fiber/v2"
)

func RoleRouter(http fiber.Router, application *app.Application) {
	authService := authService.NewService(application.DB, application.Log)
	roleService := roleService.NewService(application.DB, application.Log)
	mdw := middleware.NewEnsureToken(application.DB)

	routes := http.Group("roles")
	routes.Use(mdw.ValidateToken())

	routes.Get("/", middleware.Authorize(authService, constant.MenuSettingsRole, constant.ActionRead), handlers.ReadRoles(roleService))

	routes.Post("/", middleware.Authorize(authService, constant.MenuSettingsRole, constant.ActionCreate), handlers.CreateRole(roleService, application.Validator))

	routes.Put("/role-access/:role_id", handlers.UpdateRoleMenuPermission(roleService, application.Validator))

	routes.Put("/:role_id", middleware.Authorize(authService, constant.MenuSettingsRole, constant.ActionUpdate), handlers.UpdateRole(roleService, application.Validator))

	routes.Delete("/:role_id", middleware.Authorize(authService, constant.MenuSettingsRole, constant.ActionDelete), handlers.DeleteRole(roleService))

	routes.Get("/role-access/:role_id", middleware.Authorize(authService, constant.MenuSettingsRole, constant.ActionRead), handlers.ReadRoleAccess(roleService))

	routes.Get("/:role_id", middleware.Authorize(authService, constant.MenuSettingsRole, constant.ActionRead), handlers.ReadRoleDetail(roleService))
}
