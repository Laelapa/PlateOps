// Package app provides the core application structure and server management
// functionality for the PlateOps service. It handles HTTP server initialization,
// routing, middleware attachment, and graceful shutdown procedures.
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

// New creates and returns a new App instance with the provided dependencies.
// It initializes the HTTP server with default configuration and prepares it
// for handling requests.
//
// Parameters:
//   - ctx: The context for the application lifecycle
//   - logger: A configured zap logger for application logging
//   - queries: Database query interface for data operations
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

// newMux creates and configures the HTTP request multiplexer with all routes
// and middleware attached.
func newMux() http.Handler {

	var mux http.Handler = http.NewServeMux()
	// TODO: setup routes
	mux = attachBasicMiddleware(mux)

	return mux
}

// attachBasicMiddleware wraps the provided handler with common middleware
// functions used across all routes.
func attachBasicMiddleware(handler http.Handler) http.Handler {
	
	// TODO: handler = middleware(handler)

	return handler
}

// SetServerShutdownTimeout configures the duration the server will wait
// during shutdown before forcefully terminating connections.
//
// Parameters:
//   - t: The duration to wait during shutdown in nanoseconds
func (app *App) SetServerShutdownTimeout(t time.Duration) {

	app.serverOptions.shutdownTimeout = t
}

// LaunchServer starts the HTTP server and manages its lifecycle. It will run
// until either a server error occurs or the application context is cancelled.
// When the context is cancelled, it triggers a graceful shutdown.
//
// Returns an error if the server fails to start or encounters an error while running.
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

// ShutdownServer attempts to gracefully shut down the HTTP server within the
// configured shutdown timeout duration. If graceful shutdown fails, it forces
// the server to close. The shutdown status is logged through the application logger.
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
