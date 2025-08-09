package database

import (
	"errors"
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

// GetProducts returns the list of products stored in the database and filtered by the
// limit, offset and other filters.
func (db *Database) GetProducts(limit, offset int,
	filters ...filter.Filter,
) ([]product.Product, error) {
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

	if len(conditions) > 0 {
		conditions = append([]any{clause}, conditions...)
	} else if len(clause) > 0 {
		conditions = []any{clause}
	}

	products := make(Products, 0)
	res := db.session.Preload("Variants").Preload("Category").Limit(limit).Offset(offset)
	res = res.Find(&products, conditions...)
	if res.Error != nil {
		return nil, fmt.Errorf("unable to fetch the products: %w", res.Error)
	}

	return products.toModel(), nil
}

func (db *Database) GetProduct(productCode string) (*product.Product, error) {
	if db == nil {
		return nil, errors.New("no connection to the storage available")
	}

	if db.session == nil {
		if err := db.Connect(); err != nil {
			return nil, err
		}
	}

	var prod Product
	res := db.session.Preload("Variants").Preload("Category").Find(&prod,
		map[string]any{"code": productCode})
	if res.Error != nil {
		return nil, fmt.Errorf("unable to fetch the product %s: %w", productCode, res.Error)
	}

	if res.RowsAffected == 0 {
		return nil, product.ErrNotFound
	}

	return prod.toModel(), nil
}
