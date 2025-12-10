package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/models"
	"go.uber.org/zap"
)

type MockCategoryRepository struct {
	ListCategoryFunc   func(ctx context.Context) ([]models.Category, error)
	CreateCategoryFunc func(ctx context.Context, category *models.Category) (*models.Category, error)
}

func (m *MockCategoryRepository) ListCategory(ctx context.Context) ([]models.Category, error) {
	if m.ListCategoryFunc != nil {
		return m.ListCategoryFunc(ctx)
	}
	return nil, nil
}

func (m *MockCategoryRepository) CreateCategory(ctx context.Context, category *models.Category) (*models.Category, error) {
	if m.CreateCategoryFunc != nil {
		return m.CreateCategoryFunc(ctx, category)
	}
	return nil, nil
}

func TestGetAllCategories_Success(t *testing.T) {
	expectedCategories := []models.Category{
		{ID: 1, Code: "BOOTS", Name: "boots"},
		{ID: 2, Code: "SANDALS", Name: "sandals"},
		{ID: 3, Code: "SNEAKERS", Name: "sneakers"},
	}

	mockDB := &MockCategoryRepository{
		ListCategoryFunc: func(ctx context.Context) ([]models.Category, error) {
			return expectedCategories, nil
		},
	}

	repo := NewCategoriesRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	categories, err := repo.GetAllCategories(ctx)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(categories) != 3 {
		t.Errorf("expected 3 categories, got %d", len(categories))
	}

	if categories[0].Code != "BOOTS" {
		t.Errorf("expected first category BOOTS, got %s", categories[0].Code)
	}

	if categories[1].Name != "sandals" {
		t.Errorf("expected second category name sandals, got %s", categories[1].Name)
	}
}

func TestGetAllCategories_EmptyResult(t *testing.T) {
	mockDB := &MockCategoryRepository{
		ListCategoryFunc: func(ctx context.Context) ([]models.Category, error) {
			return []models.Category{}, nil
		},
	}

	repo := NewCategoriesRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	categories, err := repo.GetAllCategories(ctx)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(categories) != 0 {
		t.Errorf("expected 0 categories, got %d", len(categories))
	}
}

func TestGetAllCategories_DatabaseError(t *testing.T) {
	expectedError := errors.New("database connection failed")

	mockDB := &MockCategoryRepository{
		ListCategoryFunc: func(ctx context.Context) ([]models.Category, error) {
			return nil, expectedError
		},
	}

	repo := NewCategoriesRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	categories, err := repo.GetAllCategories(ctx)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("expected error '%v', got '%v'", expectedError, err)
	}

	if categories != nil {
		t.Errorf("expected nil categories, got %v", categories)
	}
}

func TestGetAllCategories_ContextCanceled(t *testing.T) {
	mockDB := &MockCategoryRepository{
		ListCategoryFunc: func(ctx context.Context) ([]models.Category, error) {
			return nil, context.Canceled
		},
	}

	repo := NewCategoriesRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	categories, err := repo.GetAllCategories(ctx)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}

	if categories != nil {
		t.Errorf("expected nil categories, got %v", categories)
	}
}

func TestCreateCategory_Success(t *testing.T) {
	inputCategory := &models.Category{
		Code: "BOOTS",
		Name: "boots",
	}

	mockDB := &MockCategoryRepository{
		CreateCategoryFunc: func(ctx context.Context, category *models.Category) (*models.Category, error) {
			category.ID = 1
			return category, nil
		},
	}

	repo := NewCategoriesRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	createdCategory, err := repo.CreateCategory(ctx, inputCategory)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if createdCategory == nil {
		t.Fatal("expected category, got nil")
	}

	if createdCategory.ID != 1 {
		t.Errorf("expected ID 1, got %d", createdCategory.ID)
	}

	if createdCategory.Code != "BOOTS" {
		t.Errorf("expected code BOOTS, got %s", createdCategory.Code)
	}

	if createdCategory.Name != "boots" {
		t.Errorf("expected name boots, got %s", createdCategory.Name)
	}
}

func TestCreateCategory_AllTypes(t *testing.T) {
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
			inputCategory := &models.Category{
				Code: cat.code,
				Name: cat.name,
			}

			mockDB := &MockCategoryRepository{
				CreateCategoryFunc: func(ctx context.Context, category *models.Category) (*models.Category, error) {
					category.ID = uint(i + 1)
					return category, nil
				},
			}

			repo := NewCategoriesRepository(mockDB, zap.NewNop().Sugar())
			ctx := context.Background()

			createdCategory, err := repo.CreateCategory(ctx, inputCategory)

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if createdCategory.Code != cat.code {
				t.Errorf("expected code %s, got %s", cat.code, createdCategory.Code)
			}

			if createdCategory.Name != cat.name {
				t.Errorf("expected name %s, got %s", cat.name, createdCategory.Name)
			}
		})
	}
}

func TestCreateCategory_DatabaseError(t *testing.T) {
	expectedError := errors.New("database error")

	inputCategory := &models.Category{
		Code: "BOOTS",
		Name: "boots",
	}

	mockDB := &MockCategoryRepository{
		CreateCategoryFunc: func(ctx context.Context, category *models.Category) (*models.Category, error) {
			return nil, expectedError
		},
	}

	repo := NewCategoriesRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	createdCategory, err := repo.CreateCategory(ctx, inputCategory)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("expected error '%v', got '%v'", expectedError, err)
	}

	if createdCategory != nil {
		t.Errorf("expected nil category, got %v", createdCategory)
	}
}

func TestCreateCategory_DuplicateCode(t *testing.T) {
	expectedError := errors.New("duplicate key value violates unique constraint")

	inputCategory := &models.Category{
		Code: "BOOTS",
		Name: "boots",
	}

	mockDB := &MockCategoryRepository{
		CreateCategoryFunc: func(ctx context.Context, category *models.Category) (*models.Category, error) {
			return nil, expectedError
		},
	}

	repo := NewCategoriesRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	createdCategory, err := repo.CreateCategory(ctx, inputCategory)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if createdCategory != nil {
		t.Errorf("expected nil category, got %v", createdCategory)
	}
}

func TestCreateCategory_ContextCanceled(t *testing.T) {
	inputCategory := &models.Category{
		Code: "BOOTS",
		Name: "boots",
	}

	mockDB := &MockCategoryRepository{
		CreateCategoryFunc: func(ctx context.Context, category *models.Category) (*models.Category, error) {
			return nil, context.Canceled
		},
	}

	repo := NewCategoriesRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	createdCategory, err := repo.CreateCategory(ctx, inputCategory)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}

	if createdCategory != nil {
		t.Errorf("expected nil category, got %v", createdCategory)
	}
}

func TestCreateCategory_VerifyInputPreserved(t *testing.T) {
	var capturedCategory *models.Category

	inputCategory := &models.Category{
		Code: "SNEAKERS",
		Name: "sneakers",
	}

	mockDB := &MockCategoryRepository{
		CreateCategoryFunc: func(ctx context.Context, category *models.Category) (*models.Category, error) {
			capturedCategory = category
			category.ID = 5
			return category, nil
		},
	}

	repo := NewCategoriesRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.Background()

	repo.CreateCategory(ctx, inputCategory)

	if capturedCategory == nil {
		t.Fatal("expected category to be captured")
	}

	if capturedCategory.Code != "SNEAKERS" {
		t.Errorf("expected code SNEAKERS, got %s", capturedCategory.Code)
	}

	if capturedCategory.Name != "sneakers" {
		t.Errorf("expected name sneakers, got %s", capturedCategory.Name)
	}
}

func TestNewCategoriesRepository(t *testing.T) {
	mockDB := &MockCategoryRepository{}
	logger := zap.NewNop().Sugar()

	repo := NewCategoriesRepository(mockDB, logger)

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

func TestGetAllCategories_VerifyContextPassed(t *testing.T) {
	var capturedCtx context.Context

	mockDB := &MockCategoryRepository{
		ListCategoryFunc: func(ctx context.Context) ([]models.Category, error) {
			capturedCtx = ctx
			return []models.Category{}, nil
		},
	}

	repo := NewCategoriesRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.WithValue(context.Background(), "test", "value")

	repo.GetAllCategories(ctx)

	if capturedCtx == nil {
		t.Fatal("expected context to be captured")
	}

	if capturedCtx.Value("test") != "value" {
		t.Error("expected context to be passed correctly")
	}
}

func TestCreateCategory_VerifyContextPassed(t *testing.T) {
	var capturedCtx context.Context

	mockDB := &MockCategoryRepository{
		CreateCategoryFunc: func(ctx context.Context, category *models.Category) (*models.Category, error) {
			capturedCtx = ctx
			category.ID = 1
			return category, nil
		},
	}

	repo := NewCategoriesRepository(mockDB, zap.NewNop().Sugar())
	ctx := context.WithValue(context.Background(), "test", "value")

	inputCategory := &models.Category{Code: "BOOTS", Name: "boots"}
	repo.CreateCategory(ctx, inputCategory)

	if capturedCtx == nil {
		t.Fatal("expected context to be captured")
	}

	if capturedCtx.Value("test") != "value" {
		t.Error("expected context to be passed correctly")
	}
}
