package routes

import "be-dashboard-nba/api/app"

func Routes(app *app.Application) {
	api := app.Server.Group("/api").Group("/v1")
	AuthRouter(api, app)
	UserRouter(api, app)
}
