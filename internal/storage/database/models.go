package database

import (
	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/product"
	"github.com/i02sopop/go-hiring-challenge-1.2.0/internal/model/variant"
	"github.com/shopspring/decimal"
)

// Product represents a product in the catalog.
// It includes a unique code and a price.
type Product struct {
	Code     string          `gorm:"uniqueIndex;not null"`
	Price    decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	Variants Variants        `gorm:"foreignKey:ProductID"`
	ID       uint            `gorm:"primaryKey"`
}

// TableName returns the table name for the Products.
func (p *Product) TableName() string {
	return "products"
}

func (p *Product) toModel() *product.Product {
	return &product.Product{
		Code:     p.Code,
		Price:    p.Price,
		Variants: p.Variants.toModel(),
		ID:       p.ID,
	}
}

type Products []Product

func (p Products) toModel() []product.Product {
	res := make([]product.Product, 0, len(p))
	for i := range p {
		res = append(res, *p[i].toModel())
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
		res = append(res, *v[i].toModel())
	}

	return res
}
