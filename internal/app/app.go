// Package app provides the core application structure and server management
// functionality for the PlateOps service. It handles HTTP server initialization,
// routing, middleware attachment, and graceful shutdown procedures.
package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Laelapa/PlateOps/internal/env"
	"github.com/Laelapa/PlateOps/internal/middleware"
	"github.com/Laelapa/PlateOps/internal/repository"
	"github.com/Laelapa/PlateOps/internal/routes"

	"github.com/Laelapa/GoHome/logging"
	"github.com/Laelapa/guarddoggo"
	"go.uber.org/zap"
)

type serverOptions struct {
	shutdownTimeout time.Duration
}

type App struct {
	ctx           context.Context
	logger        *logging.Logger
	queries       *repository.Queries
	azor          guarddoggo.GuardDoggo
	server        *http.Server
	serverOptions *serverOptions
}

const (
	defaultReadHeaderTimeout = 10 * time.Second
	defaultReadTimeout       = 30 * time.Second
	defaultWriteTimeout      = 30 * time.Second
	defaultIdleTimeout       = 120 * time.Second
)

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
	logger *logging.Logger,
	queries *repository.Queries,
	azor guarddoggo.GuardDoggo,
	port string,
	staticDir string,
	shutdownTimeout time.Duration,
) *App {

	if staticDir == "" {
		logger.LogAppWarn("Static directory not specified, using default directory 'static'")
		staticDir = "static"
	}

	return &App{
		ctx:     ctx,
		logger:  logger,
		queries: queries,
		azor:    azor,
		server: &http.Server{
			Addr:              fmt.Sprintf(":%s", env.ValidatePort(port, logger)),
			Handler:           newMux(staticDir, logger, queries, azor),
			ReadHeaderTimeout: defaultReadHeaderTimeout, // Prevents slow header attacks
			ReadTimeout:       defaultReadTimeout,       // Prevents slow request attacks
			WriteTimeout:      defaultWriteTimeout,      // Prevents clients from keeping connections open
			IdleTimeout:       defaultIdleTimeout,       // Closes idle connections
		},
		serverOptions: &serverOptions{
			shutdownTimeout: 5 * time.Second,
		},
	}
}

// newMux creates and configures the HTTP request multiplexer with all routes
// and middleware attached.
func newMux(staticDir string, logger *logging.Logger, queries *repository.Queries, azor guarddoggo.GuardDoggo) http.Handler {

	mux := routes.Setup(staticDir, logger, queries, azor)

	return attachBasicMiddleware(mux, logger)
}

// attachBasicMiddleware wraps the provided handler with common middleware
// functions used across all routes.
func attachBasicMiddleware(handler http.Handler, logger *logging.Logger) http.Handler {

	handler = middleware.SecurityResponseHeaders(handler)
	handler = middleware.CacheControlHeader(handler)
	handler = middleware.RequestLogger(handler, logger)

	return handler
}

// SetServerShutdownTimeout configures the duration the server will wait
// during shutdown before forcefully terminating connections.
//
// Parameters:
//   - t: The duration to wait during shutdown in nanoseconds
func (app *App) SetServerShutdownTimeout(t time.Duration) {

	app.serverOptions.shutdownTimeout = t
	app.logger.LogAppInfo(
		"Server shutdown timeout set",
		zap.Duration(logging.FieldDuration, t),
	)
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

		app.logger.LogAppInfo(
			"Server running",
			zap.String(logging.FieldServerAddr, app.server.Addr),
		)
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.LogAppError("Error thrown by ListenAndServe", err)
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:

		return fmt.Errorf("server failed to start: %v", err)

	case <-app.ctx.Done():

		app.logger.LogAppInfo("Shutting down server")
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

	if err := app.server.Shutdown(ctxServerShutdown); err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.logger.LogAppError("Error during server shutdown", err)
		app.logger.LogAppWarn("Closing server forcefully")
		if closeErr := app.server.Close(); closeErr != nil {
			app.logger.LogAppError("Error during forced server close", closeErr)
		}
	} else {
		app.logger.LogAppInfo("Server shut down successfully")
	}
}
