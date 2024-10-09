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
		fmt.Fprintln(os.Stderr, "ERROR: ", err)
		os.Exit(1)
	}
}

func actualMain(ctx context.Context, w io.Writer, args []string) error {
	// Handle OS signals gracefully
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Loading the .env contents
	godotenv.Load()

	// Creating a database connection pool
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		return err
	}
	defer dbPool.Close()

	dbQuery := database.New(dbPool)

	return nil //TODO: Placeholder value
}
