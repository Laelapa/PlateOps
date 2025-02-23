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

	"github.com/DemetriusPapas/PlateOps/internal/app"
	"github.com/DemetriusPapas/PlateOps/internal/logging"
	"github.com/DemetriusPapas/PlateOps/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v5"
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

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("could not load env: %v", err) // FIXME: do something more reasonable
	}

	// TODO: More configurable init based on .env options / Integrate viper

	logger, err := logging.NewLogger(os.Getenv("ENVIRONMENT"))
	if err != nil {
		return fmt.Errorf("could not initialize zap logger: %v", err) // FIXME: do something more reasonable
	}

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

	app := app.New(ctx, logger, queries)
	if err = app.LaunchServer(); err != nil {
		return fmt.Errorf("launching server error: %v", err)
	}

	return nil
}
