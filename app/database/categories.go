package database

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/models"
	"github.com/mytheresa/go-hiring-challenge/app/repository"
)

var _ repository.CategoryRepository = &PG{}

func (pg *PG) ListCategory(ctx context.Context) ([]models.Category, error) {
	pg.logger.Info("Database: listing categories")
	var categories []models.Category

	err := pg.WithContext(ctx).
		Find(&categories).Error
	if err != nil {
		pg.logger.Error("database failed to list categories: ", err)
		return nil, err
	}

	return categories, nil
}
