package routes

import app "be-dashboard-nba/internal/infrastructure/api"

func Routes(app *app.Application) {
	api := app.Server.Group("/api").Group("/v1")
	AuthRouter(api, app)
	UserRouter(api, app)
	MenuRouter(api, app)
	RoleRouter(api, app)
}
