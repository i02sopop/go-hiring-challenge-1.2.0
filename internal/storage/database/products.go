package database

import (
	"fmt"

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

func (db *Database) GetProducts(limit, offset int) ([]product.Product, error) {
	products := make(Products, 0)
	res := db.session.Preload("Variants").Preload("Category")
	res = res.Limit(limit).Offset(offset).Find(&products)
	if res.Error != nil {
		return nil, fmt.Errorf("unable to fetch the products: %w", res.Error)
	}

	return products.toModel(), nil
}
