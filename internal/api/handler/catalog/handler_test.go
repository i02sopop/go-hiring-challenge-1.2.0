package catalog

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gitlab.com/flimzy/testy"

	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/storage/database"
)

const (
	dbName     = "users"
	dbUser     = "user"
	dbPassword = "password"
)

func TestHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		product            string
		expectedStatusCode int
	}{
		{
			name:               "empty product",
			product:            "",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "non existing product",
			product:            "unknown",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "existing product",
			product:            "PROD005",
			expectedStatusCode: http.StatusOK,
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Storage initialization.
			ctx := context.TODO()
			postgresContainer := startDatabase(ctx, t)
			defer stopDatabase(t, postgresContainer)

			dbPort, err := postgresContainer.MappedPort(ctx, "5432")
			if err != nil {
				t.Errorf("unable to inspect the postgres container: %s", err)
				t.FailNow()
			}

			db := database.New(dbUser, dbPassword, dbName, dbPort.Port())
			db.Connect()
			defer db.Disconnect()

			cat := NewHandler(db)
			recorder := httptest.NewRecorder()
			mux := http.NewServeMux()
			mux.HandleFunc("GET /catalog/{code}", cat.HandleGetProduct)

			req, err := http.NewRequest("GET", fmt.Sprintf("/catalog/%s", tc.product), nil)
			if err != nil {
				t.Fatal(err)
				t.SkipNow()
			}

			mux.ServeHTTP(recorder, req)
			if tc.expectedStatusCode != recorder.Code {
				t.Errorf("Unexpected status code: expected %d, got %d", tc.expectedStatusCode,
					recorder.Code)
				t.Errorf("body: %s", recorder.Body)
				t.FailNow()
			}

			if tc.expectedStatusCode == http.StatusOK {
				d := testy.DiffAsJSON(testy.Snapshot(t), recorder.Body)
				if d != nil {
					t.Error(d)
				}
			}
		})
	}
}

func startDatabase(ctx context.Context, t *testing.T) *postgres.PostgresContainer {
	t.Helper()
	container, err := postgres.Run(ctx, "postgres:17-alpine",
		postgres.WithInitScripts(filepath.Join("testdata", "init-db.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	return container
}

func stopDatabase(t *testing.T, container *postgres.PostgresContainer) {
	t.Helper()
	if err := testcontainers.TerminateContainer(container); err != nil {
		t.Errorf("failed to terminate container: %s", err)
		t.FailNow()
	}
}
