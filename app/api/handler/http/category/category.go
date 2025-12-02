package category

import (
	"net/http"

	httpResponse "github.com/mytheresa/go-hiring-challenge/app/api/handler/http"
	"go.uber.org/zap"
)

type Category struct {
	ID   uint   `json:"id"`
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
