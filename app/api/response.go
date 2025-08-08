package api

import (
	"net/http"
)

func OKResponse(_ http.ResponseWriter, _ any) {
}

func ErrorResponse(_ http.ResponseWriter, _ int, _ string) {
}
