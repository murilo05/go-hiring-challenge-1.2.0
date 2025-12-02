package database

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/models"

	"github.com/mytheresa/go-hiring-challenge/app/repository"
)

var _ repository.ProductRepository = &PG{}

func (pg *PG) List(ctx context.Context, filters *models.Pagination) ([]models.Product, error) {
	pg.logger.Info("Database: listing and counting products")
	var products []models.Product

	if err := pg.WithContext(ctx).
		Preload("Variants").
		Preload("Category").
		Limit(filters.Limit).
		Offset(filters.Offset).
		Find(&products).Error; err != nil {
		pg.logger.Info("database failed to list products: ", err)
		return nil, err
	}

	err := pg.Model(&models.Product{}).Count(&filters.Total).Error
	if err != nil {
		pg.logger.Info("database failed to count products: ", err)
		return nil, err
	}

	return products, nil
}
