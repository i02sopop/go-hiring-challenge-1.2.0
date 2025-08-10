// Package catalog implements the catalog handler for the API.
// It defines just one endpoint to get the whole catalog of products from the products repository.
package catalog

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/api/response"
	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/filter"
	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/product"
	"github.com/shopspring/decimal"
)

const (
	limitParamName = "limit"
	defaultLimit   = 10

	offsetParamName = "offset"
	defaultOffset   = 0
)

// Category response from the API.
type Category struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type Variant struct {
	Name  string `json:"name"`
	SKU   string
	Price decimal.Decimal
}

type Product struct {
	Category Category  `json:"category"`
	Code     string    `json:"code"`
	Variants []Variant `json:"variants,omitempty"`
	Price    float64   `json:"price"`
}

type productsRepository interface {
	// GetProducts obtains a list of products from the repository with a limit and an offset.
	GetProducts(limit, offset int, filters ...filter.Filter) ([]product.Product, error)
	// GetProduct obtains a product from the storage by its code.
	GetProduct(productCode string) (*product.Product, error)
}

type Handler struct {
	repo productsRepository
}

// NewHandler returns a new api handler.
func NewHandler(r productsRepository) *Handler {
	return &Handler{
		repo: r,
	}
}

// ProductResponse defines the API response for a single product.
type ProductResponse struct {
	Product Product `json:"product"`
}

// HandleGetProduct handles the get of a product by its code.
func (h *Handler) HandleGetProduct(w http.ResponseWriter, req *http.Request) {
	productCode := req.PathValue("code")
	if productCode == "" {
		http.Error(w, "product code can't be empty", http.StatusBadRequest)

		return
	}

	res, err := h.repo.GetProduct(productCode)
	if err != nil {
		if errors.Is(err, product.ErrNotFound) {
			response.ErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	resp := ProductResponse{
		Product: Product{
			Code: res.Code,
			Category: Category{
				Name: res.Category.Name,
				Code: res.Category.Code,
			},
			Price:    res.Price.InexactFloat64(),
			Variants: make([]Variant, 0),
		},
	}
	for i := range res.Variants {
		variant := res.Variants[i]
		price := variant.Price
		if price.IsZero() {
			price = res.Price
		}

		resp.Product.Variants = append(resp.Product.Variants, Variant{
			Name:  variant.Name,
			SKU:   variant.SKU,
			Price: price,
		})
	}

	response.OKResponse(w, resp)
}

// ProductsResponse defines the API response for the list of product.
type ProductsResponse struct {
	Products    []Product `json:"products"`
	NumProducts int       `json:"total"`
	Offset      int       `json:"offset"`
}

// HandleGetProducts handle the get of a list of products.
// It accepts a page limit and an offset, and it returns the list of products, the
// offset and the number of products returned.
func (h *Handler) HandleGetProducts(w http.ResponseWriter, req *http.Request) {
	limit, err := h.getIntQueryParam(req, limitParamName, defaultLimit)
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err.Error())

		return
	}

	offset, err := h.getIntQueryParam(req, offsetParamName, defaultOffset)
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err.Error())

		return
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
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	// Map the response.
	products := make([]Product, 0, len(res))
	for i := range res {
		p := res[i]
		products = append(products, Product{
			Code:  p.Code,
			Price: p.Price.InexactFloat64(),
			Category: Category{
				Name: p.Category.Name,
				Code: p.Category.Code,
			},
		})
	}

	// Return the products.
	resp := ProductsResponse{
		NumProducts: len(products),
		Offset:      offset,
		Products:    products,
	}

	response.OKResponse(w, resp)
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
