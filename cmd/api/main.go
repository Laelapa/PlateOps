package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/DemetriusPapas/PlateOps/internal/database"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v5"
)

// Using main() as a thin wrapper around actualMain()
func main() {
	ctx := context.Background()
	if err := actualMain(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "FATAL: ", err)
		os.Exit(1)
	}
}

func actualMain(osSignalCtx context.Context, w io.Writer, args []string) error {
	// Handle OS signals gracefully
	osSignalCtx, cancel := signal.NotifyContext(osSignalCtx, os.Interrupt)
	defer cancel()

	// Loading the .env contents
	err := godotenv.Load()
	if err != nil {
		return err
	}

	// Creating a database connection pool
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		return err
	}
	defer dbPool.Close()

	dbQuery := database.New(dbPool)

	return nil //TODO: Placeholder value
}
