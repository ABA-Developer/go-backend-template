package app

import (
	"context"
	"time"
)

func GracefulShutdown(app *Application) {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Graceful shutdown timeout watcher
	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			app.Log.Error().Msg("Graceful shutdown timed out, forcing exit")
		}
	}()

	// Shutdown Fiber server
	if err := app.Server.Shutdown(); err != nil {
		app.Log.Error().Err(err).Msg("Error shutting down Fiber")
	} else {
		app.Log.Info().Msg("Fiber server stopped")
	}

	// Close DB connection
	if app.DB != nil {
		if err := app.DB.Close(); err != nil {
			app.Log.Err(err).Msg("Error closing database")
		} else {
			app.Log.Info().Msg("Database connection closed")
		}
	}

	app.Log.Info().Msg("Server exited gracefully")
}
