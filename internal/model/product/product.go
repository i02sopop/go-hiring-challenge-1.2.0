// Package product defines the product data model.
package product

import (
	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/category"
	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/variant"
	"github.com/shopspring/decimal"
)

// Product represents a product in the catalog.
type Product struct {
	Category *category.Category
	Code     string
	Price    decimal.Decimal
	Variants []variant.Variant
	ID       uint
}
