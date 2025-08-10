package catalog

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/storage/database"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gitlab.com/flimzy/testy"
)

const (
	dbName     = "users"
	dbUser     = "user"
	dbPassword = "password"
)

// nolint: funlen
func TestProductDetailsHandler(t *testing.T) {
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
			err = db.Connect()
			if err != nil {
				t.Errorf("unable to connect to the database: %s", err)
				t.FailNow()
			}

			defer func() {
				err = db.Disconnect()
				if err != nil {
					t.Errorf("unable to connect to the database: %s", err)
					t.FailNow()
				}
			}()

			req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("/catalog/%s", tc.product), nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder := doRequest(t, db, req)
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

// nolint: funlen
func TestProducsListHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		request            string
		expectedStatusCode int
	}{
		{
			name:               "list products",
			request:            "/catalog",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "with limit",
			request:            "/catalog?limit=2",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "with offset",
			request:            "/catalog?offset=2",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "with limit and offset",
			request:            "/catalog?limit=2&offset=2",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "with category",
			request:            "/catalog?category=1",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "with price",
			request:            "/catalog?price=10",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "with category and price",
			request:            "/catalog?category=1&price=15.1",
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
			err = db.Connect()
			if err != nil {
				t.Errorf("unable to connect to the database: %s", err)
				t.FailNow()
			}

			defer func() {
				err = db.Disconnect()
				if err != nil {
					t.Errorf("unable to connect to the database: %s", err)
					t.FailNow()
				}
			}()

			req, err := http.NewRequestWithContext(ctx, "GET", tc.request, nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder := doRequest(t, db, req)
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

func doRequest(t *testing.T, db *database.Database, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	cat := NewHandler(db)
	recorder := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", cat.HandleGetProducts)
	mux.HandleFunc("GET /catalog/{code}", cat.HandleGetProduct)
	mux.ServeHTTP(recorder, req)

	return recorder
}
