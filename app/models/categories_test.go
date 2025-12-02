package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategory_TableName(t *testing.T) {
	category := Category{}
	assert.Equal(t, "categories", category.TableName())
}

func TestCategory_Fields(t *testing.T) {
	category := Category{
		ID:   1,
		Code: "CLOTHING",
		Name: "Clothing",
	}

	assert.Equal(t, uint(1), category.ID)
	assert.Equal(t, "CLOTHING", category.Code)
	assert.Equal(t, "Clothing", category.Name)
}
