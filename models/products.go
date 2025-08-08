package models

import (
	"github.com/shopspring/decimal"
)

// Product represents a product in the catalog.
// It includes a unique code and a price.
type Product struct {
	Code     string          `gorm:"uniqueIndex;not null"`
	Price    decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	Variants []Variant       `gorm:"foreignKey:ProductID"`
	ID       uint            `gorm:"primaryKey"`
}

func (p *Product) TableName() string {
	return "products"
}
