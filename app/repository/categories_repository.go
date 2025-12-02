package repository

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/api/handler/http/category"
	"github.com/mytheresa/go-hiring-challenge/app/models"
	"go.uber.org/zap"
)

type CategoriesRepository struct {
	db     CategoryRepository
	logger *zap.SugaredLogger
}

func NewCategoriesRepository(db CategoryRepository, logger *zap.SugaredLogger) *CategoriesRepository {
	return &CategoriesRepository{
		db:     db,
		logger: logger,
	}
}

var _ category.CategoriesRepository = &CategoriesRepository{}

func (r *CategoriesRepository) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	r.logger.Info("Repository: fetching categories")

	categories, err := r.db.ListCategory(ctx)
	if err != nil {
		r.logger.Error("Repository: failed to fetch categories: ", err)
		return nil, err
	}

	return categories, nil
}
