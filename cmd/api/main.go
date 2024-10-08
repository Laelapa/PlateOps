package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5"
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
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt) // Handle OS signals gracefully
	defer cancel()

	dbConnection, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		return err
	}
	defer dbConnection.Close(context.Background())

	return nil //TODO: Placeholder value
}
