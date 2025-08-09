// Package category contains the category model.
package category

// Category of a product.
type Category struct {
	Code string
	Name string
	ID   uint64
}

// Categories contains a list of product categories.
type Categories []Category
