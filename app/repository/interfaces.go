package repository

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/models"
)

type ProductRepository interface {
	List(ctx context.Context, queryParams *models.QueryParams) ([]models.Product, error)
	GetByCode(ctx context.Context, code string) (*models.Product, error)
}

type CategoryRepository interface {
	ListCategory(ctx context.Context) ([]models.Category, error)
}
