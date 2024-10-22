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

// Using main() as a thin wrapper around run()
func main() {
	
	if err := run(); err != nil {
		log.Fatalf("FATAL: %v\n", err)
	}
}

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
		return fmt.Errorf("error initializing database pool: %v", err)
	}
	defer dbPool.Close()

	queries := repository.New(dbPool)

	app := app.New(ctx, logger, queries)
	if err = app.LaunchServer(); err != nil {
		return fmt.Errorf("error launching server: %v", err)
	}

	return nil
}
