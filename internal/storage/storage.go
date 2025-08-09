// Package storage defines the storage interface.
package storage

import (
	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/category"
	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/product"
)

// Storage type for the server.
type Storage interface {
	// Connect to the storage.
	Connect() error
	// GetAllProducts gets a list of all the products stored in the storage.
	GetAllProducts() ([]product.Product, error)
	// GetProducts obtains a list of products with a limit and an offset.
	GetProducts(limit, offset int) ([]product.Product, error)
	// GetAllCategories gets a list of all the categories stored in the storage.
	GetAllCategories() (category.Categories, error)
	// Disconnect from the storage.
	Disconnect() error
}
