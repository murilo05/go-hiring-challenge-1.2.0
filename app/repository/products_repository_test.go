package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/models"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type MockProductRepository struct {
	ListProductsFunc     func(ctx context.Context, queryParams *models.QueryParams) ([]models.Product, error)
	GetProductByCodeFunc func(ctx context.Context, code string) (*models.Product, error)
}

func (m *MockProductRepository) ListProducts(ctx context.Context, queryParams *models.QueryParams) ([]models.Product, error) {
	if m.ListProductsFunc != nil {
		return m.ListProductsFunc(ctx, queryParams)
	}
	return nil, nil
}

func (m *MockProductRepository) GetProductByCode(ctx context.Context, code string) (*models.Product, error) {
	if m.GetProductByCodeFunc != nil {
		return m.GetProductByCodeFunc(ctx, code)
	}
	return nil, nil
}

func getTestProducts() []models.Product {
	return []models.Product{
		{
			Code:     "PROD001",
			Price:    decimal.NewFromFloat(89.99),
			Category: models.Category{ID: 1, Name: "boots"},
		},
		{
			Code:     "PROD002",
			Price:    decimal.NewFromFloat(129.99),
			Category: models.Category{ID: 1, Name: "boots"},
		},
		{
			Code:     "PROD003",
			Price:    decimal.NewFromFloat(59.99),
			Category: models.Category{ID: 2, Name: "sandals"},
		},
	}
}

func TestGetAllProducts_Success(t *testing.T) {
	expectedProducts := getTestProducts()

	mockDB := &MockProductRepository{
		ListProductsFunc: func(ctx context.Context, queryParams *models.QueryParams) ([]models.Product, error) {
			queryParams.Pagination.Total = int64(len(expectedProducts))
			return expectedProducts, nil
		},
	}

	repo := NewProductsRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	queryParams := &models.QueryParams{
		Pagination: models.Pagination{Limit: 10, Offset: 0},
	}

	products, err := repo.GetAllProducts(ctx, queryParams)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(products) != 3 {
		t.Errorf("expected 3 products, got %d", len(products))
	}

	if products[0].Code != "PROD001" {
		t.Errorf("expected first product PROD001, got %s", products[0].Code)
	}

	if !products[0].Price.Equal(decimal.NewFromFloat(89.99)) {
		t.Errorf("expected price 89.99, got %s", products[0].Price.String())
	}
}

func TestGetAllProducts_WithPagination(t *testing.T) {
	allProducts := getTestProducts()

	mockDB := &MockProductRepository{
		ListProductsFunc: func(ctx context.Context, queryParams *models.QueryParams) ([]models.Product, error) {
			start := queryParams.Pagination.Offset
			end := start + queryParams.Pagination.Limit
			if end > len(allProducts) {
				end = len(allProducts)
			}
			queryParams.Pagination.Total = int64(len(allProducts))
			return allProducts[start:end], nil
		},
	}

	repo := NewProductsRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	queryParams := &models.QueryParams{
		Pagination: models.Pagination{Limit: 2, Offset: 1},
	}

	products, err := repo.GetAllProducts(ctx, queryParams)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}

	if products[0].Code != "PROD002" {
		t.Errorf("expected PROD002 at offset 1, got %s", products[0].Code)
	}
}

func TestGetAllProducts_WithFilters(t *testing.T) {
	allProducts := getTestProducts()

	mockDB := &MockProductRepository{
		ListProductsFunc: func(ctx context.Context, queryParams *models.QueryParams) ([]models.Product, error) {
			filtered := []models.Product{}
			for _, p := range allProducts {
				if queryParams.Filters.CategoryCode == "" || p.Category.Name == queryParams.Filters.CategoryCode {
					filtered = append(filtered, p)
				}
			}
			return filtered, nil
		},
	}

	repo := NewProductsRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	queryParams := &models.QueryParams{
		Pagination: models.Pagination{Limit: 10, Offset: 0},
		Filters:    models.Filters{CategoryCode: "boots"},
	}

	products, err := repo.GetAllProducts(ctx, queryParams)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(products) != 2 {
		t.Errorf("expected 2 boots products, got %d", len(products))
	}

	for _, p := range products {
		if p.Category.Name != "boots" {
			t.Errorf("expected category boots, got %s", p.Category.Name)
		}
	}
}

func TestGetAllProducts_EmptyResult(t *testing.T) {
	mockDB := &MockProductRepository{
		ListProductsFunc: func(ctx context.Context, queryParams *models.QueryParams) ([]models.Product, error) {
			queryParams.Pagination.Total = 0
			return []models.Product{}, nil
		},
	}

	repo := NewProductsRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	queryParams := &models.QueryParams{
		Pagination: models.Pagination{Limit: 10, Offset: 0},
	}

	products, err := repo.GetAllProducts(ctx, queryParams)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestGetAllProducts_DatabaseError(t *testing.T) {
	expectedError := errors.New("database connection failed")

	mockDB := &MockProductRepository{
		ListProductsFunc: func(ctx context.Context, queryParams *models.QueryParams) ([]models.Product, error) {
			return nil, expectedError
		},
	}

	repo := NewProductsRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	queryParams := &models.QueryParams{
		Pagination: models.Pagination{Limit: 10, Offset: 0},
	}

	products, err := repo.GetAllProducts(ctx, queryParams)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("expected error '%v', got '%v'", expectedError, err)
	}

	if products != nil {
		t.Errorf("expected nil products, got %v", products)
	}
}

func TestGetAllProducts_ContextCanceled(t *testing.T) {
	mockDB := &MockProductRepository{
		ListProductsFunc: func(ctx context.Context, queryParams *models.QueryParams) ([]models.Product, error) {
			return nil, context.Canceled
		},
	}

	repo := NewProductsRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	queryParams := &models.QueryParams{
		Pagination: models.Pagination{Limit: 10, Offset: 0},
	}

	products, err := repo.GetAllProducts(ctx, queryParams)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}

	if products != nil {
		t.Errorf("expected nil products, got %v", products)
	}
}

func TestGetProduct_Success(t *testing.T) {
	expectedProduct := &models.Product{
		Code:     "PROD005",
		Price:    decimal.NewFromFloat(99.99),
		Category: models.Category{ID: 3, Name: "sneakers"},
		Variants: []models.Variant{},
	}

	mockDB := &MockProductRepository{
		GetProductByCodeFunc: func(ctx context.Context, code string) (*models.Product, error) {
			if code == "PROD005" {
				return expectedProduct, nil
			}
			return nil, nil
		},
	}

	repo := NewProductsRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	product, err := repo.GetProduct(ctx, "PROD005")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product == nil {
		t.Fatal("expected product, got nil")
	}

	if product.Code != "PROD005" {
		t.Errorf("expected code PROD005, got %s", product.Code)
	}

	if !product.Price.Equal(decimal.NewFromFloat(99.99)) {
		t.Errorf("expected price 99.99, got %s", product.Price.String())
	}

	if product.Category.Name != "sneakers" {
		t.Errorf("expected category sneakers, got %s", product.Category.Name)
	}
}

func TestGetProduct_AllProducts(t *testing.T) {
	products := map[string]*models.Product{
		"PROD001": {Code: "PROD001", Price: decimal.NewFromFloat(89.99), Category: models.Category{Name: "boots"}},
		"PROD002": {Code: "PROD002", Price: decimal.NewFromFloat(129.99), Category: models.Category{Name: "boots"}},
		"PROD003": {Code: "PROD003", Price: decimal.NewFromFloat(59.99), Category: models.Category{Name: "sandals"}},
	}

	for code, expected := range products {
		t.Run(code, func(t *testing.T) {
			mockDB := &MockProductRepository{
				GetProductByCodeFunc: func(ctx context.Context, code string) (*models.Product, error) {
					return products[code], nil
				},
			}

			repo := NewProductsRepository(mockDB, zap.NewNop().Sugar())
			ctx := context.Background()

			product, err := repo.GetProduct(ctx, code)

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if product.Code != expected.Code {
				t.Errorf("expected code %s, got %s", expected.Code, product.Code)
			}
		})
	}
}

func TestGetProduct_NotFound(t *testing.T) {
	mockDB := &MockProductRepository{
		GetProductByCodeFunc: func(ctx context.Context, code string) (*models.Product, error) {
			return nil, nil
		},
	}

	repo := NewProductsRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	product, err := repo.GetProduct(ctx, "PROD999")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product != nil {
		t.Errorf("expected nil product, got %v", product)
	}
}

func TestGetProduct_DatabaseError(t *testing.T) {
	expectedError := errors.New("database error")

	mockDB := &MockProductRepository{
		GetProductByCodeFunc: func(ctx context.Context, code string) (*models.Product, error) {
			return nil, expectedError
		},
	}

	repo := NewProductsRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	product, err := repo.GetProduct(ctx, "PROD001")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("expected error '%v', got '%v'", expectedError, err)
	}

	if product != nil {
		t.Errorf("expected nil product, got %v", product)
	}
}

func TestGetProduct_ContextCanceled(t *testing.T) {
	mockDB := &MockProductRepository{
		GetProductByCodeFunc: func(ctx context.Context, code string) (*models.Product, error) {
			return nil, context.Canceled
		},
	}

	repo := NewProductsRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	product, err := repo.GetProduct(ctx, "PROD001")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}

	if product != nil {
		t.Errorf("expected nil product, got %v", product)
	}
}

func TestGetProduct_EmptyCode(t *testing.T) {
	var capturedCode string

	mockDB := &MockProductRepository{
		GetProductByCodeFunc: func(ctx context.Context, code string) (*models.Product, error) {
			capturedCode = code
			return nil, nil
		},
	}

	repo := NewProductsRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	repo.GetProduct(ctx, "")

	if capturedCode != "" {
		t.Errorf("expected empty code, got '%s'", capturedCode)
	}
}

func TestNewProductsRepository(t *testing.T) {
	mockDB := &MockProductRepository{}
	logger := zap.NewNop().Sugar()

	repo := NewProductsRepository(mockDB, logger)

	if repo == nil {
		t.Fatal("expected repository, got nil")
	}

	if repo.db == nil {
		t.Error("expected db to be set")
	}

	if repo.logger == nil {
		t.Error("expected logger to be set")
	}
}

func TestGetAllProducts_VerifyContextPassed(t *testing.T) {
	var capturedCtx context.Context

	mockDB := &MockProductRepository{
		ListProductsFunc: func(ctx context.Context, queryParams *models.QueryParams) ([]models.Product, error) {
			capturedCtx = ctx
			return []models.Product{}, nil
		},
	}

	repo := NewProductsRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.WithValue(context.Background(), "test", "value")

	queryParams := &models.QueryParams{
		Pagination: models.Pagination{Limit: 10, Offset: 0},
	}

	repo.GetAllProducts(ctx, queryParams)

	if capturedCtx == nil {
		t.Fatal("expected context to be captured")
	}

	if capturedCtx.Value("test") != "value" {
		t.Error("expected context to be passed correctly")
	}
}

func TestGetProduct_VerifyContextPassed(t *testing.T) {
	var capturedCtx context.Context

	mockDB := &MockProductRepository{
		GetProductByCodeFunc: func(ctx context.Context, code string) (*models.Product, error) {
			capturedCtx = ctx
			return nil, nil
		},
	}

	repo := NewProductsRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.WithValue(context.Background(), "test", "value")

	repo.GetProduct(ctx, "PROD001")

	if capturedCtx == nil {
		t.Fatal("expected context to be captured")
	}

	if capturedCtx.Value("test") != "value" {
		t.Error("expected context to be passed correctly")
	}
}

func TestGetAllProducts_VerifyQueryParamsMutated(t *testing.T) {
	mockDB := &MockProductRepository{
		ListProductsFunc: func(ctx context.Context, queryParams *models.QueryParams) ([]models.Product, error) {
			queryParams.Pagination.Total = 100
			return []models.Product{}, nil
		},
	}

	repo := NewProductsRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	queryParams := &models.QueryParams{
		Pagination: models.Pagination{Limit: 10, Offset: 0, Total: 0},
	}

	repo.GetAllProducts(ctx, queryParams)

	if queryParams.Pagination.Total != 100 {
		t.Errorf("expected total to be mutated to 100, got %d", queryParams.Pagination.Total)
	}
}
