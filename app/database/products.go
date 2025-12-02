package database

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/models"

	"github.com/mytheresa/go-hiring-challenge/app/repository"
)

var _ repository.ProductRepository = &PG{}

func (pg *PG) List(ctx context.Context, queryParams *models.QueryParams) ([]models.Product, error) {
	pg.logger.Info("Database: listing and counting products")
	var products []models.Product

	query := pg.Model(&models.Product{}).Joins("JOIN categories ON products.category_id = categories.id")

	if queryParams.Filters.CategoryCode != "" {
		query = query.Where("categories.code = ?", queryParams.Filters.CategoryCode)
	}

	if queryParams.Filters.PriceLessThan != nil {
		query = query.Where("products.price < ?", queryParams.Filters.PriceLessThan)
	}

	err := query.Count(&queryParams.Pagination.Total).Error
	if err != nil {
		pg.logger.Error("database failed to count products: ", err)
		return nil, err
	}

	if err := query.WithContext(ctx).
		Preload("Variants").
		Preload("Category").
		Limit(queryParams.Pagination.Limit).
		Offset(queryParams.Pagination.Offset).
		Find(&products).Error; err != nil {
		pg.logger.Info("database failed to list products: ", err)
		return nil, err
	}

	return products, nil
}
