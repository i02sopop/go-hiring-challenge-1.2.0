package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/api/handler/catalog"
	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/storage"
)

const (
	readHeaderTimeout = 1 * time.Second
)

// Server implements the api http server.
type Server struct {
	logger  *slog.Logger
	st      storage.Storage
	srv     *http.Server
	address string
}

// NewServer initializes the api server.
func NewServer(addr string, st storage.Storage) *Server {
	return &Server{
		address: addr,
		logger:  slog.Default().With("address", addr),
		st:      st,
		srv: &http.Server{
			Addr:              addr,
			Handler:           router(st),
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}
}

func router(st storage.Storage) http.Handler {
	cat := catalog.NewHandler(st)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", cat.HandleGetProducts)
	mux.HandleFunc("GET /catalog/{code}", cat.HandleGetProduct)

	return mux
}

// Start the server.
func (s *Server) Start(ctx context.Context) error {
	// Initialize database connection
	if err := s.st.Connect(); err != nil {
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

	return s.st.Disconnect()
}
