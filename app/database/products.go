package database

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/models"
	"github.com/mytheresa/go-hiring-challenge/app/repository"
)

var _ repository.ProductRepository = &PG{}

func (pg *PG) List(ctx context.Context, filters models.Filters) ([]models.Product, error) {
	var products []models.Product

	if err := pg.Preload("Variants").
		Limit(filters.Limit).
		Offset(filters.Offset).
		Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}
