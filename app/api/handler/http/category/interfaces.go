package category

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/models"
)

type CategoriesRepository interface {
	GetAllCategories(ctx context.Context) ([]models.Category, error)
	CreateCategory(ctx context.Context, category *models.Category) (*models.Category, error)
}
