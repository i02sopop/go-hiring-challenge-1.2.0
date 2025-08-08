package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOKResponse(t *testing.T) {
	t.Parallel()
	type sampleResponse struct {
		Message string `json:"message"`
	}

	sample := sampleResponse{Message: "Success"}

	t.Run("successful http200 json response", func(t *testing.T) {
		t.Parallel()
		recorder := httptest.NewRecorder()
		OKResponse(recorder, sample)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200 OK")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"),
			"Expected Content-Type to be application/json")

		expected := `{"message":"Success"}`
		assert.JSONEq(t, expected, recorder.Body.String(), "Response body does not match expected")
	})
}

func TestErrorResponse(t *testing.T) {
	t.Parallel()
	t.Run("json response for a given http status code", func(t *testing.T) {
		t.Parallel()
		recorder := httptest.NewRecorder()
		ErrorResponse(recorder, http.StatusInternalServerError, "Some error occurred")

		assert.Equal(t, http.StatusInternalServerError, recorder.Code,
			"Expected status code 500 Internal Server Error")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"),
			"Expected Content-Type to be application/json")

		expected := `{"error":"Some error occurred"}`
		assert.JSONEq(t, expected, recorder.Body.String(), "Response body does not match expected")
	})
}
