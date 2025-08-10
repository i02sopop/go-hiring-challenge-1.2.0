// Package response defines the api handler responses.
package response

import (
	"net/http"
)

func OKResponse(_ http.ResponseWriter, _ any) {
}

func ErrorResponse(_ http.ResponseWriter, _ int, _ string) {
}
