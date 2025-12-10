package category

import (
	"encoding/json"
	"net/http"

	httpResponse "github.com/mytheresa/go-hiring-challenge/app/api/handler/http"
	"github.com/mytheresa/go-hiring-challenge/app/models"
	"go.uber.org/zap"
)

type Category struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type CreateCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CategoryHandler struct {
	repo   CategoriesRepository
	logger *zap.SugaredLogger
}

func NewCategoryHandler(r CategoriesRepository, logger *zap.SugaredLogger) *CategoryHandler {
	return &CategoryHandler{
		repo:   r,
		logger: logger,
	}
}

// ListCategories godoc
// @Summary      List all categories
// @Description  Get a list of all available product categories
// @Tags         categories
// @Accept       json
// @Produce      json
// @Success      200  {object}  httpResponse.response{data=[]Category}  "List of categories"
// @Failure      500  {object}  httpResponse.errorResponse              "Internal server error"
// @Router       /categories [get]
func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := h.repo.GetAllCategories(ctx)
	if err != nil {
		h.logger.Error("failed to get categories: ", err)
		httpResponse.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}

	categories := make([]Category, len(res))
	for i, c := range res {
		categories[i] = Category{
			ID:   c.ID,
			Code: c.Code,
			Name: c.Name,
		}
	}

	httpResponse.OKResponse(w, categories)
}

// CreateCategory godoc
// @Summary      Create a new category
// @Description  Create a new product category with code and name
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        request  body      CreateCategoryRequest  true  "Category data"
// @Success      200      {object}  httpResponse.response{data=Category}  "Category created successfully"
// @Failure      400      {object}  httpResponse.errorResponse            "Bad request - invalid input or missing fields"
// @Failure      500      {object}  httpResponse.errorResponse            "Internal server error"
// @Router       /categories [post]
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body: ", err)
		httpResponse.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Code == "" || req.Name == "" {
		httpResponse.ErrorResponse(w, http.StatusBadRequest, "Code and name are required")
		return
	}

	category := &models.Category{
		Code: req.Code,
		Name: req.Name,
	}

	savedCategory, err := h.repo.CreateCategory(ctx, category)
	if err != nil {
		h.logger.Error("failed to create category: ", err)
		httpResponse.ErrorResponse(w, http.StatusInternalServerError, "Failed to create category")
		return
	}

	response := Category{
		ID:   savedCategory.ID,
		Code: savedCategory.Code,
		Name: savedCategory.Name,
	}

	httpResponse.OKResponse(w, response)

}
