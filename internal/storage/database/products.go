package database

import (
	"fmt"

	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/filter"
	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/product"
)

// GetAllProducts returns the list of all products stored in the database.
func (db *Database) GetAllProducts() ([]product.Product, error) {
	products := make(Products, 0)
	res := db.session.Preload("Variants").Preload("Category").Find(&products)
	if res.Error != nil {
		return nil, fmt.Errorf("unable to fetch the products: %w", res.Error)
	}

	return products.toModel(), nil
}

func (db *Database) GetProducts(limit, offset int, filters ...filter.Filter) ([]product.Product, error) {
	products := make(Products, 0)
	res := db.session.Preload("Variants").Preload("Category").Limit(limit).Offset(offset)

	clause := ""
	conditions := make([]any, 0)
	for i := range filters {
		filter := filters[i]
		if i > 0 {
			clause += " AND "
		}

		clause += fmt.Sprintf("%s %s ?", filter.Key, filter.Operation)
		conditions = append(conditions, filter.Value)
	}

	conditions = append([]any{clause}, conditions...)
	res = res.Find(&products, conditions...)
	if res.Error != nil {
		return nil, fmt.Errorf("unable to fetch the products: %w", res.Error)
	}

	return products.toModel(), nil
}
