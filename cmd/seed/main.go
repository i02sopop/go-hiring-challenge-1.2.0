package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/storage/database"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		slog.ErrorContext(ctx, "error loading the environment file", "error", err)
		os.Exit(-1)
	}

	dir := os.Getenv("POSTGRES_SQL_DIR")
	files, err := os.ReadDir(dir)
	if err != nil {
		slog.ErrorContext(ctx, "unable to read the directory", "directory", dir, "error", err)
		os.Exit(-1)
	}

	// Initialize database connection
	db := database.New(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PORT"))
	if err := db.Connect(); err != nil {
		slog.ErrorContext(ctx, "unable to connect to the database", "error", err)

		os.Exit(-1)
	}

	// Filter and sort .sql files
	var sqlFiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file)
		}
	}

	sort.Slice(sqlFiles, func(i, j int) bool {
		return sqlFiles[i].Name() < sqlFiles[j].Name()
	})

	for _, file := range sqlFiles {
		path := filepath.Join(dir, file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			slog.WarnContext(ctx, "unable to read the file", "filename",
				file.Name(), "error", err)

			continue
		}

		sql := string(content)
		if err := db.Exec(sql).Error; err != nil {
			// We are not logging the content in case the query has sensitive data.
			slog.WarnContext(ctx, "unable to execute the query", "file",
				file.Name(), "error", err)
			cleanup(ctx, db)

			os.Exit(-1)
		}

		slog.InfoContext(ctx, "query executed properly", "file", file.Name())
	}

	cleanup(ctx, db)
}

func cleanup(ctx context.Context, db *database.Database) {
	err := db.Disconnect()
	if err != nil {
		slog.ErrorContext(ctx, "unable to close the database connection", "error", err)
		os.Exit(-1)
	}
}
