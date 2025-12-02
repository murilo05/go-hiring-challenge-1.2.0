package repository

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/models"
)

type ProductRepository interface {
	List(ctx context.Context, queryParams *models.QueryParams) ([]models.Product, error)
}
