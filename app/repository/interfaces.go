package repository

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/models"
)

type ProductRepository interface {
	List(ctx context.Context, filters *models.Pagination) ([]models.Product, error)
}
