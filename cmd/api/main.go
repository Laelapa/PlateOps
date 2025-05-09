// Package main implements the entry point for the PlateOps API server.
// It handles the initialization of core components including logging,
// database connection, and the HTTP server.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Laelapa/PlateOps/internal/app"
	"github.com/Laelapa/PlateOps/internal/repository"
	"go.uber.org/zap"

	"github.com/Laelapa/GoHome/logging"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/jackc/pgx/v5"
)

const (
	DefaultShutdownTimeout = 5 * time.Second // Time until forceful shutdown
)

// main serves as the entry point for the application and acts as a thin wrapper
// around the run function. It will terminate the application with a fatal log
// if run encounters an error.
func main() {
	if err := run(); err != nil {
		log.Fatalf("FATAL: %v\n", err)
	}
}

// run initializes and orchestrates all components of the application:
//   - Sets up signal handling for graceful shutdown
//   - Loads environment variables
//   - Initializes the logger
//   - Establishes database connection
//   - Creates repository queries
//   - Launches the HTTP server
//
// Returns an error if any initialization step fails.
func run() error {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// FIXME: add error handling, automate loading of .env file or pulling from secrets manager
	//
	// Disabled for production, use your service's secrets instead
	// Refer to the dotenv example file for the required environment variables
	// Uncomment the following lines to load environment variables from a .env file in a local development environment
	//
	// err := godotenv.Load()
	// if err != nil {
	// 	return fmt.Errorf("error loading .env file: %w", err)
	// }

	// TODO: More configurable init based on .env options / Integrate viper

	logger, err := logging.NewLogger(os.Getenv("ENVIRONMENT"))
	if err != nil {
		return fmt.Errorf("error creating logger: %w", err) // FIXME: add error handling
	}

	defer func() {
		if syncErr := logger.Sync(); syncErr != nil { // FIXME: handle case of writing to unbuffered output that doesnt support sync
			// Print the error without crashing the program
			fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", syncErr)
		}
	}()

	dbPool, err := pgxpool.New(ctx, os.Getenv("DB_URL"))
	if err != nil {
		return fmt.Errorf("database initialization failed: %w", err)
	}
	defer dbPool.Close()

	// Verify database connection
	if err := dbPool.Ping(ctx); err != nil {
		return fmt.Errorf("database connection check failed: %w", err)
	}

	queries := repository.New(dbPool)

	// Parse the server shutdown timeout from the environment
	shutdownTimeout, err := time.ParseDuration(os.Getenv("SERVER_SHUTDOWN_TIMEOUT") + "s")
	if err != nil {
		shutdownTimeout = DefaultShutdownTimeout // fallback default
		logger.LogAppWarn(
			"Failed to parse SERVER_SHUTDOWN_TIMEOUT, falling back to default",
			zap.Duration("shutdown timeout", shutdownTimeout),
		)
	}

		app := app.New(
		ctx,
		logger,
		os.Getenv("SERVER_PORT"),
		os.Getenv("STATIC_DIR"), // FIXME: check if this is a valid path
		queries,
		shutdownTimeout,
	)
	if err = app.LaunchServer(); err != nil {
		return fmt.Errorf("error launching server: %w", err) // FIXME: add error handling
	}
	return nil
}
