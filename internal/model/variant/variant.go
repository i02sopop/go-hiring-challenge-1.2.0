// Package variant defines the product variant model.
// Variants can be used to represent different configurations or options for a product.
package variant

import (
	"github.com/shopspring/decimal"
)

// Variant represents a product variant.
// Variants can be used to represent different configurations or options for a product.
type Variant struct {
	Name      string
	SKU       string
	Price     decimal.Decimal
	ID        uint
	ProductID uint
}
