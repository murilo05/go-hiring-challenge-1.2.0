package repository

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/api/handler/http/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/models"

	"go.uber.org/zap"
)

type ProductsRepository struct {
	db     ProductRepository
	logger *zap.SugaredLogger
}

func NewProductsRepository(db ProductRepository, logger *zap.SugaredLogger) *ProductsRepository {
	return &ProductsRepository{
		db:     db,
		logger: logger,
	}
}

var _ catalog.ProductsRepository = &ProductsRepository{}

func (r *ProductsRepository) GetAllProducts(ctx context.Context, queryParams *models.QueryParams) ([]models.Product, error) {
	r.logger.Info("Repository: fetching products")

	products, err := r.db.ListProducts(ctx, queryParams)
	if err != nil {
		r.logger.Error("Repository: failed to fetch products: ", err)
		return nil, err
	}

	return products, nil
}

func (r *ProductsRepository) GetProduct(ctx context.Context, code string) (*models.Product, error) {
	r.logger.Infof("Repository: fetching product with code %s", code)

	product, err := r.db.GetProductByCode(ctx, code)
	if err != nil {
		r.logger.Error("Repository: failed to fetch product: ", err)
		return nil, err
	}

	return product, nil
}
