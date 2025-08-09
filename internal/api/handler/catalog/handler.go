// Package catalog implements the catalog handler for the API.
// It defines just one endpoint to get the whole catalog of products from the products repository.
package catalog

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/filter"
	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/product"
)

const (
	limitParamName = "limit"
	defaultLimit   = 10

	offsetParamName = "offset"
	defaultOffset   = 0
)

type Response struct {
	Products    []Product `json:"products"`
	NumProducts int       `json:"total"`
	Offset      int       `json:"offset"`
}

type Product struct {
	Code     string  `json:"code"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
}

type productsRepository interface {
	// GetProducts obtains a list of products from the repository with a limit and an offset.
	GetProducts(limit, offset int, filters ...filter.Filter) ([]product.Product, error)
}

type Handler struct {
	repo productsRepository
}

func NewHandler(r productsRepository) *Handler {
	return &Handler{
		repo: r,
	}
}

func (h *Handler) HandleGet(w http.ResponseWriter, req *http.Request) {
	limit, err := h.getIntQueryParam(req, limitParamName, defaultLimit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	offset, err := h.getIntQueryParam(req, offsetParamName, defaultOffset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	filters := make([]filter.Filter, 0)
	cat := h.getQueryParam(req, "category")
	if cat != "" {
		filters = append(filters, filter.Filter{
			Key:       "category",
			Value:     cat,
			Operation: filter.Equal,
		})
	}

	price := h.getQueryParam(req, "price")
	if price != "" {
		filters = append(filters, filter.Filter{
			Key:       "price",
			Value:     price,
			Operation: filter.LessThan,
		})
	}

	res, err := h.repo.GetProducts(limit, offset, filters...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// Map the response.
	products := make([]Product, 0, len(res))
	for i := range res {
		p := res[i]
		products = append(products, Product{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: p.Category.Name,
		})
	}

	// Return the products as a JSON response.
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		NumProducts: len(products),
		Offset:      offset,
		Products:    products,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) getIntQueryParam(req *http.Request, name string, defaultValue int) (int, error) {
	param := h.getQueryParam(req, name)
	if param == "" {
		return defaultValue, nil
	}

	return strconv.Atoi(param)
}

func (h *Handler) getQueryParam(req *http.Request, name string) string {
	params := req.URL.Query()
	queryParam, ok := params[name]
	if !ok {
		return ""
	}

	return queryParam[0]
}
