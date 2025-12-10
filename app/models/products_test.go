package models

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestProduct_TableName(t *testing.T) {
	product := Product{}
	assert.Equal(t, "products", product.TableName())
}

func TestProduct_Fields(t *testing.T) {
	price := decimal.NewFromFloat(10.99)
	product := Product{
		ID:         1,
		Code:       "PROD001",
		Price:      price,
		CategoryID: 1,
	}

	assert.Equal(t, uint(1), product.ID)
	assert.Equal(t, "PROD001", product.Code)
	assert.Equal(t, price, product.Price)
	assert.Equal(t, uint(1), product.CategoryID)
}
