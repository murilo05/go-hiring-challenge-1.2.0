package models

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestVariant_TableName(t *testing.T) {
	variant := &Variant{}
	assert.Equal(t, "product_variants", variant.TableName())
}

func TestVariant_Struct(t *testing.T) {
	variant := &Variant{
		ID:        1,
		ProductID: 100,
		Name:      "Red Variant",
		SKU:       "TEST-SKU-001",
		Price:     decimal.NewFromFloat(29.99),
	}

	assert.Equal(t, uint(1), variant.ID)
	assert.Equal(t, uint(100), variant.ProductID)
	assert.Equal(t, "Red Variant", variant.Name)
	assert.Equal(t, "TEST-SKU-001", variant.SKU)
	assert.True(t, variant.Price.Equal(decimal.NewFromFloat(29.99)))
}

func TestVariant_ZeroValues(t *testing.T) {
	variant := &Variant{}

	assert.Equal(t, uint(0), variant.ID)
	assert.Equal(t, uint(0), variant.ProductID)
	assert.Equal(t, "", variant.Name)
	assert.Equal(t, "", variant.SKU)
	assert.True(t, variant.Price.IsZero())
}

func TestVariant_PriceOperations(t *testing.T) {
	variant := &Variant{
		Price: decimal.NewFromFloat(19.99),
	}

	assert.True(t, variant.Price.GreaterThan(decimal.NewFromFloat(10.0)))
	assert.True(t, variant.Price.LessThan(decimal.NewFromFloat(25.0)))

	assert.Equal(t, "19.99", variant.Price.String())
}
