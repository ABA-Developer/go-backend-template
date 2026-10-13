package routes

import (
	"github.com/ABA-Developer/go-backend-template/constant"
	authUseCase "github.com/ABA-Developer/go-backend-template/internal/application/auth/usecase"
	roleUseCase "github.com/ABA-Developer/go-backend-template/internal/application/role/usecase"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/runtime/container"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/middleware"
	handlers "github.com/ABA-Developer/go-backend-template/internal/presentation/role/handler"

	"github.com/gofiber/fiber/v2"
)

func AddRoleRoutes(http fiber.Router, c *container.Container) {
	authSvc := authUseCase.NewUseCase(c.GetDB())
	roleSvc := roleUseCase.NewUseCase(c.GetDB())
	mdw := middleware.NewEnsureToken(c.GetDB())

	routes := http.Group("roles")
	routes.Use(mdw.ValidateToken())

	routes.Get("/", middleware.Authorize(authSvc, constant.MenuSettingsRole, constant.ActionRead), handlers.ReadRoles(roleSvc))

	routes.Post("/", middleware.Authorize(authSvc, constant.MenuSettingsRole, constant.ActionCreate), handlers.CreateRole(roleSvc, c.GetValidator()))

	routes.Put("/role-access/:role_id", handlers.UpdateRoleMenuPermission(roleSvc, c.GetValidator()))

	routes.Put("/:role_id", middleware.Authorize(authSvc, constant.MenuSettingsRole, constant.ActionUpdate), handlers.UpdateRole(roleSvc, c.GetValidator()))

	routes.Delete("/:role_id", middleware.Authorize(authSvc, constant.MenuSettingsRole, constant.ActionDelete), handlers.DeleteRole(roleSvc))

	routes.Get("/role-access/:role_id", middleware.Authorize(authSvc, constant.MenuSettingsRole, constant.ActionRead), handlers.ReadRoleAccess(roleSvc))

	routes.Get("/:role_id", middleware.Authorize(authSvc, constant.MenuSettingsRole, constant.ActionRead), handlers.ReadRoleDetailUseCase(roleSvc))
}
