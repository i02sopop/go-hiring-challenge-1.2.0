package database

import (
	"github.com/shopspring/decimal"

	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/category"
	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/product"
	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/variant"
)

// Product represents a product in the catalog.
// It includes a unique code and a price.
type Product struct {
	Code       string          `gorm:"uniqueIndex;not null"`
	Price      decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	Category   Category        `gorm:"foreignKey:CategoryID"`
	Variants   Variants        `gorm:"foreignKey:ProductID"`
	CategoryID uint            `gorm:"column:category"`
	ID         uint            `gorm:"primaryKey"`
}

// TableName returns the table name for the Products.
func (p *Product) TableName() string {
	return "products"
}

func (p *Product) toModel() *product.Product {
	if p == nil {
		return nil
	}

	return &product.Product{
		Code:     p.Code,
		Price:    p.Price,
		Variants: p.Variants.toModel(),
		ID:       p.ID,
		Category: p.Category.toModel(),
	}
}

type Products []Product

func (p Products) toModel() []product.Product {
	res := make([]product.Product, 0, len(p))
	for i := range p {
		elem := p[i].toModel()
		if elem != nil {
			res = append(res, *elem)
		}
	}

	return res
}

// Variant represents a product variant in the catalog.
// It includes a unique name, SKU, and an optional price.
// Variants can be used to represent different configurations or options for a product.
type Variant struct {
	Name      string          `gorm:"not null"`
	SKU       string          `gorm:"uniqueIndex;not null"`
	Price     decimal.Decimal `gorm:"type:decimal(10,2);null"`
	ID        uint            `gorm:"primaryKey"`
	ProductID uint            `gorm:"not null"`
}

// TableName returns the database table name for the Variants.
func (v *Variant) TableName() string {
	return "product_variants"
}

func (v *Variant) toModel() *variant.Variant {
	if v == nil {
		return nil
	}

	return &variant.Variant{
		Name:      v.Name,
		SKU:       v.SKU,
		Price:     v.Price,
		ID:        v.ID,
		ProductID: v.ProductID,
	}
}

type Variants []Variant

func (v Variants) toModel() []variant.Variant {
	res := make([]variant.Variant, 0, len(v))
	for i := range v {
		elem := v[i].toModel()
		if elem != nil {
			res = append(res, *elem)
		}
	}

	return res
}

// Category of a product.
type Category struct {
	Code string `gorm:"uniqueIndex;not null"`
	Name string `gorm:"not null"`
	ID   uint64 `gorm:"primaryKey"`
}

func (c *Category) TableName() string {
	return "category"
}

func (c *Category) toModel() *category.Category {
	if c == nil {
		return nil
	}

	return &category.Category{
		ID:   c.ID,
		Code: c.Code,
		Name: c.Name,
	}
}

type Categories []Category

func (c Categories) toModel() category.Categories {
	res := make(category.Categories, 0, len(c))
	for i := range c {
		elem := c[i].toModel()
		if elem != nil {
			res = append(res, *elem)
		}
	}

	return res
}
