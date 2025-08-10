package database

import (
	"fmt"

	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/category"
)

func (db *Database) GetAllCategories() (category.Categories, error) {
	categories := make(Categories, 0)
	res := db.session.Find(&categories)
	if res.Error != nil {
		return nil, fmt.Errorf("unable to fetch the categories: %w", res.Error)
	}

	return categories.toModel(), nil
}

func (db *Database) AddCategory(cat *category.Category) error {
	if cat.Name == "" {
		return category.ErrInvalidCategory
	}

	if cat.Code == "" {
		return category.ErrInvalidCategory
	}

	c := Category{
		Name: cat.Name,
		Code: cat.Code,
	}
	res := db.session.Create(&c)
	if res.Error != nil {
		return fmt.Errorf("unable to create the category: %w", res.Error)
	}

	return nil
}
