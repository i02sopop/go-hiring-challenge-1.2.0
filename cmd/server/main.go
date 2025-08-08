package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		slog.ErrorContext(ctx, "error loading the environment file", "error", err)
		os.Exit(-1)
	}

	// signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	srv := NewServer(fmt.Sprintf("localhost:%s", os.Getenv("HTTP_PORT")))
	if err := srv.Start(ctx); err != nil {
		slog.ErrorContext(ctx, "Unable to start the http server", "error", err)
		stop()

		os.Exit(-1)
	}

	<-ctx.Done()
	if err := srv.Stop(ctx); err != nil {
		slog.ErrorContext(ctx, "Unable to shutdown the server", "error", err)
		stop()

		os.Exit(-1)
	}

	stop()
}
