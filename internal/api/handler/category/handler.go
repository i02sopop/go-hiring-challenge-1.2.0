// Package category implements the category handler for the API.
package category

import (
	"encoding/json"
	"net/http"

	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/category"
)

// Category model for the API response.
type Category struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type categoryRepository interface {
	// GetAllCategories gets a list of all the categories stored in the storage.
	GetAllCategories() (category.Categories, error)
}

type Handler struct {
	repo categoryRepository
}

// NewHandler returns a new api handler.
func NewHandler(r categoryRepository) *Handler {
	return &Handler{
		repo: r,
	}
}

// CategoriesResponse defines the API response for the list of categories.
type CategoriesResponse struct {
	Categories    []Category `json:"products"`
	NumCategories int        `json:"total"`
}

// HandleGetCategories handle the get of a list of categories.
func (h *Handler) HandleGetCategories(w http.ResponseWriter, req *http.Request) {
	res, err := h.repo.GetAllCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// Map the response.
	categories := make([]Category, 0, len(res))
	for i := range res {
		p := res[i]
		categories = append(categories, Category{
			Code: p.Code,
			Name: p.Name,
		})
	}

	// Return the products as a JSON response.
	w.Header().Set("Content-Type", "application/json")
	response := CategoriesResponse{
		NumCategories: len(categories),
		Categories:    categories,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
