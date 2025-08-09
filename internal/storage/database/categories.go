package database

import (
	"fmt"

	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/category"
)

func (db *Database) GetAllCategories() (category.Categories, error) {
	categories := make(Categories, 0)
	res := db.session.Find(&categories)
	if res.Error != nil {
		return nil, fmt.Errorf("unable to fetch the products: %w", res.Error)
	}

	return categories.toModel(), nil
}
