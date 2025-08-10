package category

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
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
func TestGetHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		expectedStatusCode int
	}{
		{
			name:               "list of categories",
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

			cats := NewHandler(db)
			recorder := httptest.NewRecorder()
			mux := http.NewServeMux()
			mux.HandleFunc("GET /categories", cats.HandleGetCategories)

			req, err := http.NewRequestWithContext(ctx, "GET", "/categories", nil)
			if err != nil {
				t.Fatal(err)
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

// nolint: funlen, gocognit, gocyclo, cyclop
func TestPostHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		category           io.Reader
		name               string
		expectedStatusCode int
	}{
		{
			name:               "empty",
			category:           strings.NewReader(``),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "without data",
			category:           strings.NewReader(`{}`),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "with emty name and code",
			category:           strings.NewReader(`{"code": "", "name": ""}`),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "with emty name",
			category:           strings.NewReader(`{"code": "", "name": "test"}`),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "with full data",
			category:           strings.NewReader(`{"code": "test", "name": "test"}`),
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

			req, err := http.NewRequestWithContext(ctx, "GET", "/categories", nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder := doRequest(t, db, req)
			if recorder.Code != http.StatusOK {
				t.Errorf("Error doging the request: %d - %s", recorder.Code, recorder.Body)
				t.FailNow()
			}

			d := testy.DiffAsJSON(testy.Snapshot(t, "presave"), recorder.Body)
			if d != nil {
				t.Error(d)
			}

			req, err = http.NewRequestWithContext(ctx, "POST", "/categories", tc.category)
			if err != nil {
				t.Fatal(err)
			}

			recorder = doRequest(t, db, req)
			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Unexpected status code: expected %d, got %d", tc.expectedStatusCode,
					recorder.Code)
				t.FailNow()
			}

			req, err = http.NewRequestWithContext(ctx, "GET", "/categories", nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder = doRequest(t, db, req)
			if recorder.Code != http.StatusOK {
				t.Errorf("Error doging the request: %d - %s", recorder.Code, recorder.Body)
				t.FailNow()
			}

			d = testy.DiffAsJSON(testy.Snapshot(t, "postsave"), recorder.Body)
			if d != nil {
				t.Error(d)
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
	cats := NewHandler(db)
	recorder := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /categories", cats.HandleGetCategories)
	mux.HandleFunc("POST /categories", cats.HandlePostCategories)
	mux.ServeHTTP(recorder, req)

	return recorder
}
