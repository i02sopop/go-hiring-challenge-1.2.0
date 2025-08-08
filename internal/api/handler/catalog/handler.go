// Package catalog implements the catalog handler for the API.
// It defines just one endpoint to get the whole catalog of products from the products repository.
package catalog

import (
	"encoding/json"
	"net/http"

	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/product"
)

type Response struct {
	Products []Product `json:"products"`
}

type Product struct {
	Code  string  `json:"code"`
	Price float64 `json:"price"`
}

type productsRepository interface {
	GetAllProducts() ([]product.Product, error)
}

type Handler struct {
	repo productsRepository
}

func NewHandler(r productsRepository) *Handler {
	return &Handler{
		repo: r,
	}
}

func (h *Handler) HandleGet(w http.ResponseWriter, _ *http.Request) {
	res, err := h.repo.GetAllProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// Map the response.
	products := make([]Product, 0, len(res))
	for i := range res {
		p := res[i]
		products = append(products, Product{
			Code:  p.Code,
			Price: p.Price.InexactFloat64(),
		})
	}

	// Return the products as a JSON response.
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Products: products,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
