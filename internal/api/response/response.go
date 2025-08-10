// Package response defines the api handler responses.
package response

import (
	"encoding/json"
	"net/http"
)

func OKResponse(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

type errorResponse struct {
	Message string `json:"error"`
}

func ErrorResponse(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(&errorResponse{Message: msg}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
