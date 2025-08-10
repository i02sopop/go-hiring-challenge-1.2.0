// Package category implements the category handler for the API.
package category

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/api/response"
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
	// AddCategory adds a new category to the storage.
	AddCategory(cat *category.Category) error
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
func (h *Handler) HandleGetCategories(w http.ResponseWriter, _ *http.Request) {
	res, err := h.repo.GetAllCategories()
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())

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
	resp := CategoriesResponse{
		NumCategories: len(categories),
		Categories:    categories,
	}

	response.OKResponse(w, resp)
}

// HandlePostCategories handle the creation of a new category.
func (h *Handler) HandlePostCategories(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var cat Category
	err := decoder.Decode(&cat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	newCategory := &category.Category{
		Name: cat.Name,
		Code: cat.Code,
	}
	err = h.repo.AddCategory(newCategory)
	if err != nil {
		if errors.Is(err, category.ErrInvalidCategory) {
			response.ErrorResponse(w, http.StatusBadRequest, err.Error())
		} else {
			response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
	}

	response.OKResponse(w, newCategory)
}
