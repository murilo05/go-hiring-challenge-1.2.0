package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/models"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type MockProductsRepository struct {
	GetAllProductsFunc func(ctx context.Context, params *models.QueryParams) ([]models.Product, error)
	GetProductFunc     func(ctx context.Context, code string) (*models.Product, error)
}

func (m *MockProductsRepository) GetAllProducts(ctx context.Context, params *models.QueryParams) ([]models.Product, error) {
	if m.GetAllProductsFunc != nil {
		return m.GetAllProductsFunc(ctx, params)
	}
	return nil, nil
}

func (m *MockProductsRepository) GetProduct(ctx context.Context, code string) (*models.Product, error) {
	if m.GetProductFunc != nil {
		return m.GetProductFunc(ctx, code)
	}
	return nil, nil
}

func getRealCatalogData() []models.Product {
	return []models.Product{
		{Code: "PROD001", Price: decimal.NewFromFloat(89.99), Category: models.Category{ID: 1, Name: "boots"}},
		{Code: "PROD002", Price: decimal.NewFromFloat(129.99), Category: models.Category{ID: 1, Name: "boots"}},
		{Code: "PROD003", Price: decimal.NewFromFloat(59.99), Category: models.Category{ID: 2, Name: "sandals"}},
		{Code: "PROD004", Price: decimal.NewFromFloat(79.99), Category: models.Category{ID: 2, Name: "sandals"}},
		{Code: "PROD005", Price: decimal.NewFromFloat(99.99), Category: models.Category{ID: 3, Name: "sneakers"}},
		{Code: "PROD006", Price: decimal.NewFromFloat(149.99), Category: models.Category{ID: 3, Name: "sneakers"}},
		{Code: "PROD007", Price: decimal.NewFromFloat(199.99), Category: models.Category{ID: 1, Name: "boots"}},
		{Code: "PROD008", Price: decimal.NewFromFloat(39.99), Category: models.Category{ID: 2, Name: "sandals"}},
	}
}

func TestGetProduct_SuccessWithVariants(t *testing.T) {
	mockRepo := &MockProductsRepository{
		GetProductFunc: func(ctx context.Context, code string) (*models.Product, error) {
			if code == "PROD005" {
				return &models.Product{
					Code:     "PROD005",
					Price:    decimal.NewFromFloat(99.99),
					Category: models.Category{ID: 3, Name: "sneakers"},
					Variants: []models.Variant{},
				}, nil
			}
			return nil, nil
		},
	}

	handler := NewCatalogHandler(mockRepo, zap.NewNop().Sugar())
	req := httptest.NewRequest(http.MethodGet, "/products/PROD005", nil)
	req.SetPathValue("code", "PROD005")
	w := httptest.NewRecorder()

	handler.GetProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var response struct {
		Data GetProductResponse `json:"data"`
	}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Data.Product.Code != "PROD005" {
		t.Errorf("expected code PROD005, got %s", response.Data.Product.Code)
	}

	if response.Data.Product.Price != 99.99 {
		t.Errorf("expected price 99.99, got %.2f", response.Data.Product.Price)
	}

	if response.Data.Product.Category != "sneakers" {
		t.Errorf("expected category sneakers, got %s", response.Data.Product.Category)
	}
}

func TestGetProduct_NotFound(t *testing.T) {
	mockRepo := &MockProductsRepository{
		GetProductFunc: func(ctx context.Context, code string) (*models.Product, error) {
			return nil, nil
		},
	}

	handler := NewCatalogHandler(mockRepo, zap.NewNop().Sugar())
	req := httptest.NewRequest(http.MethodGet, "/products/PROD999", nil)
	req.SetPathValue("code", "PROD999")
	w := httptest.NewRecorder()

	handler.GetProduct(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestGetProduct_EmptyCode(t *testing.T) {
	handler := NewCatalogHandler(&MockProductsRepository{}, zap.NewNop().Sugar())
	req := httptest.NewRequest(http.MethodGet, "/products/", nil)
	w := httptest.NewRecorder()

	handler.GetProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestGetProduct_RepositoryError(t *testing.T) {
	mockRepo := &MockProductsRepository{
		GetProductFunc: func(ctx context.Context, code string) (*models.Product, error) {
			return nil, errors.New("database error")
		},
	}

	handler := NewCatalogHandler(mockRepo, zap.NewNop().Sugar())
	req := httptest.NewRequest(http.MethodGet, "/products/PROD001", nil)
	req.SetPathValue("code", "PROD001")
	w := httptest.NewRecorder()

	handler.GetProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestListProducts_AllProducts(t *testing.T) {
	catalog := getRealCatalogData()

	mockRepo := &MockProductsRepository{
		GetAllProductsFunc: func(ctx context.Context, params *models.QueryParams) ([]models.Product, error) {
			params.Pagination.Total = int64(len(catalog))
			return catalog, nil
		},
	}

	handler := NewCatalogHandler(mockRepo, zap.NewNop().Sugar())
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler.ListProducts(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var response struct {
		Data ListProductsResponse `json:"data"`
	}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response.Data.Products) != 8 {
		t.Errorf("expected 8 products, got %d", len(response.Data.Products))
	}

	if response.Data.Pagination.Total != 8 {
		t.Errorf("expected total 8, got %d", response.Data.Pagination.Total)
	}

	if response.Data.Pagination.Limit != 10 {
		t.Errorf("expected default limit 10, got %d", response.Data.Pagination.Limit)
	}

	if response.Data.Products[0].Code != "PROD001" {
		t.Errorf("expected first product PROD001, got %s", response.Data.Products[0].Code)
	}
}

func TestListProducts_WithPagination(t *testing.T) {
	catalog := getRealCatalogData()

	mockRepo := &MockProductsRepository{
		GetAllProductsFunc: func(ctx context.Context, params *models.QueryParams) ([]models.Product, error) {
			params.Pagination.Total = int64(len(catalog))
			start := params.Pagination.Offset
			end := start + params.Pagination.Limit
			if end > len(catalog) {
				end = len(catalog)
			}
			return catalog[start:end], nil
		},
	}

	handler := NewCatalogHandler(mockRepo, zap.NewNop().Sugar())
	req := httptest.NewRequest(http.MethodGet, "/products?limit=3&offset=2", nil)
	w := httptest.NewRecorder()

	handler.ListProducts(w, req)

	var response struct {
		Data ListProductsResponse `json:"data"`
	}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response.Data.Products) != 3 {
		t.Errorf("expected 3 products, got %d", len(response.Data.Products))
	}

	if response.Data.Products[0].Code != "PROD003" {
		t.Errorf("expected PROD003 at offset 2, got %s", response.Data.Products[0].Code)
	}

	if response.Data.Pagination.Limit != 3 {
		t.Errorf("expected limit 3, got %d", response.Data.Pagination.Limit)
	}

	if response.Data.Pagination.Offset != 2 {
		t.Errorf("expected offset 2, got %d", response.Data.Pagination.Offset)
	}
}

func TestListProducts_FilterByCategory(t *testing.T) {
	catalog := getRealCatalogData()

	mockRepo := &MockProductsRepository{
		GetAllProductsFunc: func(ctx context.Context, params *models.QueryParams) ([]models.Product, error) {
			filtered := []models.Product{}
			for _, p := range catalog {
				if p.Category.Name == params.Filters.CategoryCode {
					filtered = append(filtered, p)
				}
			}
			params.Pagination.Total = int64(len(filtered))
			return filtered, nil
		},
	}

	handler := NewCatalogHandler(mockRepo, zap.NewNop().Sugar())
	req := httptest.NewRequest(http.MethodGet, "/products?category=boots", nil)
	w := httptest.NewRecorder()

	handler.ListProducts(w, req)

	var response struct {
		Data ListProductsResponse `json:"data"`
	}
	json.NewDecoder(w.Body).Decode(&response)

	if len(response.Data.Products) != 3 {
		t.Errorf("expected 3 boots products, got %d", len(response.Data.Products))
	}

	for _, p := range response.Data.Products {
		if p.Category != "boots" {
			t.Errorf("expected category boots, got %s for product %s", p.Category, p.Code)
		}
	}
}

func TestListProducts_FilterByPrice(t *testing.T) {
	catalog := getRealCatalogData()

	mockRepo := &MockProductsRepository{
		GetAllProductsFunc: func(ctx context.Context, params *models.QueryParams) ([]models.Product, error) {
			filtered := []models.Product{}
			for _, p := range catalog {
				if params.Filters.PriceLessThan == nil || p.Price.LessThan(*params.Filters.PriceLessThan) {
					filtered = append(filtered, p)
				}
			}
			params.Pagination.Total = int64(len(filtered))
			return filtered, nil
		},
	}

	handler := NewCatalogHandler(mockRepo, zap.NewNop().Sugar())
	req := httptest.NewRequest(http.MethodGet, "/products?price_less_than=100", nil)
	w := httptest.NewRecorder()

	handler.ListProducts(w, req)

	var response struct {
		Data ListProductsResponse `json:"data"`
	}
	json.NewDecoder(w.Body).Decode(&response)

	for _, p := range response.Data.Products {
		if p.Price >= 100 {
			t.Errorf("expected price < 100, got %.2f for %s", p.Price, p.Code)
		}
	}
}

func TestListProducts_FilterByCategoryAndPrice(t *testing.T) {
	catalog := getRealCatalogData()

	mockRepo := &MockProductsRepository{
		GetAllProductsFunc: func(ctx context.Context, params *models.QueryParams) ([]models.Product, error) {
			filtered := []models.Product{}
			for _, p := range catalog {
				matchCategory := params.Filters.CategoryCode == "" || p.Category.Name == params.Filters.CategoryCode
				matchPrice := params.Filters.PriceLessThan == nil || p.Price.LessThan(*params.Filters.PriceLessThan)
				if matchCategory && matchPrice {
					filtered = append(filtered, p)
				}
			}
			params.Pagination.Total = int64(len(filtered))
			return filtered, nil
		},
	}

	handler := NewCatalogHandler(mockRepo, zap.NewNop().Sugar())
	req := httptest.NewRequest(http.MethodGet, "/products?category=sandals&price_less_than=70", nil)
	w := httptest.NewRecorder()

	handler.ListProducts(w, req)

	var response struct {
		Data ListProductsResponse `json:"data"`
	}
	json.NewDecoder(w.Body).Decode(&response)

	for _, p := range response.Data.Products {
		if p.Category != "sandals" {
			t.Errorf("expected category sandals, got %s", p.Category)
		}
		if p.Price >= 70 {
			t.Errorf("expected price < 70, got %.2f", p.Price)
		}
	}
}

func TestListProducts_RepositoryError(t *testing.T) {
	mockRepo := &MockProductsRepository{
		GetAllProductsFunc: func(ctx context.Context, params *models.QueryParams) ([]models.Product, error) {
			return nil, errors.New("database connection failed")
		},
	}

	handler := NewCatalogHandler(mockRepo, zap.NewNop().Sugar())
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler.ListProducts(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("expected error response body, got empty")
	}

	var errorResponse map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
		t.Logf("Response body: %s", body)
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errorMsg, ok := errorResponse["error"].(string); ok {
		if errorMsg == "" {
			t.Error("expected non-empty error message")
		}
	} else if errorMsg, ok := errorResponse["message"].(string); ok {
		if errorMsg == "" {
			t.Error("expected non-empty error message")
		}
	} else {
		t.Logf("Error response: %+v", errorResponse)
	}
}

func TestListProducts_DatabaseTimeout(t *testing.T) {
	mockRepo := &MockProductsRepository{
		GetAllProductsFunc: func(ctx context.Context, params *models.QueryParams) ([]models.Product, error) {
			return nil, errors.New("context deadline exceeded")
		},
	}

	handler := NewCatalogHandler(mockRepo, zap.NewNop().Sugar())
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler.ListProducts(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("expected error response body, got empty")
	}
}

func TestListProducts_EmptyResult(t *testing.T) {
	mockRepo := &MockProductsRepository{
		GetAllProductsFunc: func(ctx context.Context, params *models.QueryParams) ([]models.Product, error) {
			params.Pagination.Total = 0
			return []models.Product{}, nil
		},
	}

	handler := NewCatalogHandler(mockRepo, zap.NewNop().Sugar())
	req := httptest.NewRequest(http.MethodGet, "/products?category=nonexistent", nil)
	w := httptest.NewRecorder()

	handler.ListProducts(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response struct {
		Data ListProductsResponse `json:"data"`
	}
	json.NewDecoder(w.Body).Decode(&response)

	if len(response.Data.Products) != 0 {
		t.Errorf("expected 0 products, got %d", len(response.Data.Products))
	}

	if response.Data.Pagination.Total != 0 {
		t.Errorf("expected total 0, got %d", response.Data.Pagination.Total)
	}
}

func TestParsePaginationParams_DefaultValues(t *testing.T) {
	handler := &CatalogHandler{}
	req := httptest.NewRequest(http.MethodGet, "/products", nil)

	pagination := handler.parsePaginationParams(req)

	if pagination.Limit != 10 {
		t.Errorf("expected default limit 10, got %d", pagination.Limit)
	}

	if pagination.Offset != 0 {
		t.Errorf("expected default offset 0, got %d", pagination.Offset)
	}
}

func TestParsePaginationParams_OffsetBoundaries(t *testing.T) {
	handler := &CatalogHandler{}

	tests := []struct {
		name     string
		offset   string
		expected int
	}{
		{"zero offset", "0", 0},
		{"positive offset", "10", 10},
		{"negative offset", "-5", 0},
		{"invalid string", "abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/products?offset="+tt.offset, nil)
			pagination := handler.parsePaginationParams(req)

			if pagination.Offset != tt.expected {
				t.Errorf("for offset=%s, expected %d, got %d", tt.offset, tt.expected, pagination.Offset)
			}
		})
	}
}

func TestParseFilterParams_ValidCategory(t *testing.T) {
	handler := &CatalogHandler{}

	categories := []string{"boots", "sandals", "sneakers"}
	for _, category := range categories {
		t.Run(category, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/products?category="+category, nil)
			filters := handler.parseFilterParams(req)

			if filters.CategoryCode != category {
				t.Errorf("expected category %s, got %s", category, filters.CategoryCode)
			}
		})
	}
}

func TestParseFilterParams_ValidPrice(t *testing.T) {
	handler := &CatalogHandler{}
	req := httptest.NewRequest(http.MethodGet, "/products?price_less_than=99.99", nil)

	filters := handler.parseFilterParams(req)

	if filters.PriceLessThan == nil {
		t.Fatal("expected price filter to be set")
	}

	expected := decimal.NewFromFloat(99.99)
	if !filters.PriceLessThan.Equal(expected) {
		t.Errorf("expected price 99.99, got %s", filters.PriceLessThan.String())
	}
}

func TestParseFilterParams_InvalidPrice(t *testing.T) {
	handler := &CatalogHandler{}

	tests := []struct {
		name  string
		price string
	}{
		{"invalid string", "invalid"},
		{"negative", "-10"},
		{"zero", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/products?price_less_than="+tt.price, nil)
			filters := handler.parseFilterParams(req)

			if filters.PriceLessThan != nil {
				t.Errorf("for price=%s, expected nil filter, got %v", tt.price, filters.PriceLessThan)
			}
		})
	}
}

func TestParseFilterParams_NoFilters(t *testing.T) {
	handler := &CatalogHandler{}
	req := httptest.NewRequest(http.MethodGet, "/products", nil)

	filters := handler.parseFilterParams(req)

	if filters.CategoryCode != "" {
		t.Errorf("expected empty category, got %s", filters.CategoryCode)
	}

	if filters.PriceLessThan != nil {
		t.Errorf("expected nil price filter, got %v", filters.PriceLessThan)
	}
}
