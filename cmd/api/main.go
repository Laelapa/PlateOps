package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
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
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)	// Handle OS signals gracefully
	defer cancel()


	return nil //TODO: Placeholder value
}
