package category

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/models"
	"go.uber.org/zap"
)

type MockCategoriesRepository struct {
	GetAllCategoriesFunc func(ctx context.Context) ([]models.Category, error)
	CreateCategoryFunc   func(ctx context.Context, category *models.Category) (*models.Category, error)
}

func (m *MockCategoriesRepository) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	if m.GetAllCategoriesFunc != nil {
		return m.GetAllCategoriesFunc(ctx)
	}
	return nil, nil
}

func (m *MockCategoriesRepository) CreateCategory(ctx context.Context, category *models.Category) (*models.Category, error) {
	if m.CreateCategoryFunc != nil {
		return m.CreateCategoryFunc(ctx, category)
	}
	return nil, nil
}

func TestListCategories_Success(t *testing.T) {
	mockCategories := []models.Category{
		{ID: 1, Code: "BOOTS", Name: "boots"},
		{ID: 2, Code: "SANDALS", Name: "sandals"},
		{ID: 3, Code: "SNEAKERS", Name: "sneakers"},
	}

	mockRepo := &MockCategoriesRepository{
		GetAllCategoriesFunc: func(ctx context.Context) ([]models.Category, error) {
			return mockCategories, nil
		},
	}

	handler := NewCategoryHandler(mockRepo, zap.NewNop().Sugar())
	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	handler.ListCategories(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var response struct {
		Data []Category `json:"data"`
	}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response.Data) != 3 {
		t.Errorf("expected 3 categories, got %d", len(response.Data))
	}

	if response.Data[0].Code != "BOOTS" {
		t.Errorf("expected first category BOOTS, got %s", response.Data[0].Code)
	}

	if response.Data[1].Name != "sandals" {
		t.Errorf("expected second category name sandals, got %s", response.Data[1].Name)
	}
}

func TestListCategories_EmptyList(t *testing.T) {
	mockRepo := &MockCategoriesRepository{
		GetAllCategoriesFunc: func(ctx context.Context) ([]models.Category, error) {
			return []models.Category{}, nil
		},
	}

	handler := NewCategoryHandler(mockRepo, zap.NewNop().Sugar())
	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	handler.ListCategories(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response struct {
		Data []Category `json:"data"`
	}
	json.NewDecoder(w.Body).Decode(&response)

	if len(response.Data) != 0 {
		t.Errorf("expected 0 categories, got %d", len(response.Data))
	}
}

func TestListCategories_RepositoryError(t *testing.T) {
	mockRepo := &MockCategoriesRepository{
		GetAllCategoriesFunc: func(ctx context.Context) ([]models.Category, error) {
			return nil, errors.New("database connection failed")
		},
	}

	handler := NewCategoryHandler(mockRepo, zap.NewNop().Sugar())
	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	handler.ListCategories(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("expected error response body, got empty")
	}
}

func TestCreateCategory_Success(t *testing.T) {
	mockRepo := &MockCategoriesRepository{
		CreateCategoryFunc: func(ctx context.Context, category *models.Category) (*models.Category, error) {
			category.ID = 1
			return category, nil
		},
	}

	handler := NewCategoryHandler(mockRepo, zap.NewNop().Sugar())

	reqBody := CreateCategoryRequest{
		Code: "BOOTS",
		Name: "boots",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateCategory(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var response struct {
		Data Category `json:"data"`
	}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Data.ID != 1 {
		t.Errorf("expected ID 1, got %d", response.Data.ID)
	}

	if response.Data.Code != "BOOTS" {
		t.Errorf("expected code BOOTS, got %s", response.Data.Code)
	}

	if response.Data.Name != "boots" {
		t.Errorf("expected name boots, got %s", response.Data.Name)
	}
}

func TestCreateCategory_AllCategories(t *testing.T) {
	categories := []struct {
		code string
		name string
	}{
		{"BOOTS", "boots"},
		{"SANDALS", "sandals"},
		{"SNEAKERS", "sneakers"},
	}

	for i, cat := range categories {
		t.Run(cat.code, func(t *testing.T) {
			mockRepo := &MockCategoriesRepository{
				CreateCategoryFunc: func(ctx context.Context, category *models.Category) (*models.Category, error) {
					category.ID = uint(i + 1)
					return category, nil
				},
			}

			handler := NewCategoryHandler(mockRepo, zap.NewNop().Sugar())

			reqBody := CreateCategoryRequest{
				Code: cat.code,
				Name: cat.name,
			}
			bodyBytes, _ := json.Marshal(reqBody)

			req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateCategory(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", w.Code)
			}

			var response struct {
				Data Category `json:"data"`
			}
			json.NewDecoder(w.Body).Decode(&response)

			if response.Data.Code != cat.code {
				t.Errorf("expected code %s, got %s", cat.code, response.Data.Code)
			}
		})
	}
}

func TestCreateCategory_InvalidJSON(t *testing.T) {
	mockRepo := &MockCategoriesRepository{}
	handler := NewCategoryHandler(mockRepo, zap.NewNop().Sugar())

	invalidJSON := []byte(`{"code": "BOOTS", "name":`)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateCategory(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCreateCategory_MissingCode(t *testing.T) {
	mockRepo := &MockCategoriesRepository{}
	handler := NewCategoryHandler(mockRepo, zap.NewNop().Sugar())

	reqBody := CreateCategoryRequest{
		Code: "",
		Name: "boots",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateCategory(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCreateCategory_MissingName(t *testing.T) {
	mockRepo := &MockCategoriesRepository{}
	handler := NewCategoryHandler(mockRepo, zap.NewNop().Sugar())

	reqBody := CreateCategoryRequest{
		Code: "BOOTS",
		Name: "",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateCategory(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCreateCategory_MissingBothFields(t *testing.T) {
	mockRepo := &MockCategoriesRepository{}
	handler := NewCategoryHandler(mockRepo, zap.NewNop().Sugar())

	reqBody := CreateCategoryRequest{
		Code: "",
		Name: "",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateCategory(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCreateCategory_EmptyBody(t *testing.T) {
	mockRepo := &MockCategoriesRepository{}
	handler := NewCategoryHandler(mockRepo, zap.NewNop().Sugar())

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateCategory(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCreateCategory_RepositoryError(t *testing.T) {
	mockRepo := &MockCategoriesRepository{
		CreateCategoryFunc: func(ctx context.Context, category *models.Category) (*models.Category, error) {
			return nil, errors.New("database error")
		},
	}

	handler := NewCategoryHandler(mockRepo, zap.NewNop().Sugar())

	reqBody := CreateCategoryRequest{
		Code: "BOOTS",
		Name: "boots",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateCategory(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("expected error response body, got empty")
	}
}

func TestCreateCategory_VerifyRepositoryCalled(t *testing.T) {
	var capturedCategory *models.Category
	mockRepo := &MockCategoriesRepository{
		CreateCategoryFunc: func(ctx context.Context, category *models.Category) (*models.Category, error) {
			capturedCategory = category
			category.ID = 5
			return category, nil
		},
	}

	handler := NewCategoryHandler(mockRepo, zap.NewNop().Sugar())

	reqBody := CreateCategoryRequest{
		Code: "SNEAKERS",
		Name: "sneakers",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateCategory(w, req)

	if capturedCategory == nil {
		t.Fatal("expected repository to be called")
	}

	if capturedCategory.Code != "SNEAKERS" {
		t.Errorf("expected code SNEAKERS, got %s", capturedCategory.Code)
	}

	if capturedCategory.Name != "sneakers" {
		t.Errorf("expected name sneakers, got %s", capturedCategory.Name)
	}
}
