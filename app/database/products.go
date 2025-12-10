package database

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/mytheresa/go-hiring-challenge/app/repository"
)

var _ repository.ProductRepository = &PG{}

func (pg *PG) ListProducts(ctx context.Context, queryParams *models.QueryParams) ([]models.Product, error) {
	pg.logger.Info("Database: listing and counting products")
	var products []models.Product

	query := pg.Model(&models.Product{})

	if queryParams.Filters.CategoryCode != "" && queryParams.Filters.PriceLessThan != nil {
		query = pg.createFilteredQuery(queryParams)
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

func (pg *PG) GetProductByCode(ctx context.Context, code string) (*models.Product, error) {
	pg.logger.Info("Database: getting product by code %s", code)
	var product models.Product

	err := pg.WithContext(ctx).
		Preload("Variants").
		Preload("Category").
		Where("code = ?", code).
		First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			pg.logger.Warn("product not found with code: ", code)
			return nil, nil
		}
		pg.logger.Error("database failed to get product: ", err)
		return nil, err
	}

	for i := range product.Variants {
		if product.Variants[i].Price.IsZero() || product.Variants[i].Price.Equal(decimal.Zero) {
			product.Variants[i].Price = product.Price
		}
	}

	return &product, nil
}

func (pg *PG) createFilteredQuery(queryParams *models.QueryParams) *gorm.DB {
	query := pg.Model(&models.Product{})

	if queryParams.Filters.CategoryCode != "" {
		query = query.Joins("JOIN categories ON products.category_id = categories.id").
			Where("categories.code = ?", queryParams.Filters.CategoryCode)
	}

	if queryParams.Filters.PriceLessThan != nil {
		query = query.Where("products.price < ?", queryParams.Filters.PriceLessThan)
	}

	return query
}
