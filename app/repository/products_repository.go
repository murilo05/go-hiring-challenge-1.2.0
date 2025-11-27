package repository

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/models"
)

type ProductsRepository struct {
	db ProductRepository
}

func NewProductsRepository(db ProductRepository) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(ctx context.Context, filters models.Filters) ([]models.Product, error) {
	products, err := r.db.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	return products, nil
}
