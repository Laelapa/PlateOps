package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/DemetriusPapas/PlateOps/internal/repository"
	"go.uber.org/zap"
)

type serverOptions struct {
	shutdownTimeout time.Duration
}

type App struct {
	ctx           context.Context
	logger        *zap.SugaredLogger
	queries       *repository.Queries
	server        *http.Server
	serverOptions *serverOptions
}

func New(
	ctx context.Context,
	logger *zap.SugaredLogger,
	queries *repository.Queries,
) *App {

	return &App{
		ctx:     ctx,
		logger:  logger,
		queries: queries,
		server: &http.Server{
			Addr:    ":8080", // TODO: grab from .env instead
			Handler: newMux(),
		},
		serverOptions: &serverOptions{
			shutdownTimeout: 5 * time.Second,
		},
	}
}

func newMux() http.Handler {

	var mux http.Handler = http.NewServeMux()
	// TODO: setup routes
	mux = attachBasicMiddleware(mux)

	return mux
}

func attachBasicMiddleware(handler http.Handler) http.Handler {

	// TODO: handler = middleware(handler)

	return handler
}

func (app *App) SetServerShutdownTimeout(t time.Duration) {

	app.serverOptions.shutdownTimeout = t
}

func (app *App) LaunchServer() error {

	errChan := make(chan error, 1)
	defer close(errChan)

	go func() {

		app.logger.Infof("Server running on %s", app.server.Addr)
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Errorf("Error by ListenAndServe(): %v\n", err)
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:

		return fmt.Errorf("server failed to start: %v", err)

	case <-app.ctx.Done():

		app.ShutdownServer()
		return nil
	}
}

// Attempts to gracefully shut the server down via http.Server.Shutdown() in an amount
// of time determined by the shutdownTimeout option. If that fails then the server gets
// forcefully terminated by calling http.Server.Close().
func (app *App) ShutdownServer() {

	ctxServerShutdown, cancel := context.WithTimeout(context.Background(), app.serverOptions.shutdownTimeout)
	defer cancel()

	if err := app.server.Shutdown(ctxServerShutdown); err != nil && err != http.ErrServerClosed {
		app.logger.Errorf("Error during server shutdown: %v\n", err)
		app.logger.Infof("Closing server forcefully\n")
		app.server.Close()
	} else {
		app.logger.Infof("Server shut down successfully\n")
	}

}
