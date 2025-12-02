package repository

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/models"
)

type ProductRepository interface {
	ListProducts(ctx context.Context, queryParams *models.QueryParams) ([]models.Product, error)
	GetProductByCode(ctx context.Context, code string) (*models.Product, error)
}

type CategoryRepository interface {
	ListCategory(ctx context.Context) ([]models.Category, error)
	CreateCategory(ctx context.Context, category *models.Category) (*models.Category, error)
}
