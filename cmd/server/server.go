package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/i02sopop/go-hiring-challenge-1.2.0/app/catalog"
	"github.com/i02sopop/go-hiring-challenge-1.2.0/app/database"
	"github.com/i02sopop/go-hiring-challenge-1.2.0/models"
)

const (
	readHeaderTimeout = 1 * time.Second
)

type Server struct {
	logger  *slog.Logger
	db      *database.Database
	srv     *http.Server
	address string
}

func NewServer(addr string) *Server {
	return &Server{
		address: addr,
		logger:  slog.Default().With("address", addr),
		db: database.New(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PORT")),
		srv: &http.Server{
			Addr:              addr,
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	// Initialize database connection
	if err := s.db.Connect(); err != nil {
		return fmt.Errorf("unable to connect to the database: %w", err)
	}

	dbSession, err := s.db.Session()
	if err != nil {
		return fmt.Errorf("unable to obtain the database session: %w", err)
	}

	// Initialize handlers
	prodRepo := models.NewProductsRepository(dbSession)
	cat := catalog.NewHandler(prodRepo)

	// Set up routing
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", cat.HandleGet)
	s.srv.Handler = mux

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

func (s *Server) Stop(ctx context.Context) error {
	if s.srv == nil {
		return nil
	}

	slog.InfoContext(ctx, "Shutting down the server", "server", s.srv)
	if err := s.db.Disconnect(); err != nil {
		return err
	}

	return s.srv.Shutdown(ctx)
}
