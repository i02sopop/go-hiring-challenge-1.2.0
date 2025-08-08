package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/api/handler/catalog"
	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/storage/database"
)

const (
	readHeaderTimeout = 1 * time.Second
)

// Server implements the api http server.
type Server struct {
	logger  *slog.Logger
	db      *database.Database
	srv     *http.Server
	address string
}

// NewServer initializes the api server.
func NewServer(addr string) *Server {
	db := database.New(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PORT"))

	return &Server{
		address: addr,
		logger:  slog.Default().With("address", addr),
		db:      db,
		srv: &http.Server{
			Addr:              addr,
			Handler:           router(db),
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}
}

func router(db *database.Database) http.Handler {
	cat := catalog.NewHandler(db)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", cat.HandleGet)

	return mux
}

// Start the server.
func (s *Server) Start(ctx context.Context) error {
	// Initialize database connection
	if err := s.db.Connect(); err != nil {
		return fmt.Errorf("unable to connect to the database: %w", err)
	}

	// Start the server
	go func() {
		slog.InfoContext(ctx, "Starting the server...", "server", s.srv)
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(ctx, "Server failed", "server", s.srv, "error", err)

			return
		}

		slog.InfoContext(ctx, "Server stopped gracefully...", "server", s.srv)
	}()

	return nil
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) error {
	if s.srv == nil {
		return nil
	}

	slog.InfoContext(ctx, "Shutting down the server", "server", s.srv)
	// We shutdown the API before the database to make sure the database is not
	// used anymore.
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}

	return s.db.Disconnect()
}
