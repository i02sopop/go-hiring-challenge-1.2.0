// Package category contains the category model.
package category

import "errors"

var ErrInvalidCategory = errors.New("category not valid")

// Category of a product.
type Category struct {
	Code string
	Name string
	ID   uint64
}

// Categories contains a list of product categories.
type Categories []Category
