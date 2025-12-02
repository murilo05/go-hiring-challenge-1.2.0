package catalog

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/models"
)

type ProductsRepository interface {
	GetAllProducts(ctx context.Context, filters *models.Pagination) ([]models.Product, error)
}
