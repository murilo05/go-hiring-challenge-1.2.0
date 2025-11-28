package catalog

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/models"

	"github.com/mytheresa/go-hiring-challenge/app/repository"

	"go.uber.org/zap"
)

type Response struct {
	Products []Product `json:"products"`
}

type Product struct {
	Code  string  `json:"code"`
	Price float64 `json:"price"`
}

type CatalogHandler struct {
	repo   *repository.ProductsRepository
	logger *zap.SugaredLogger
}

func NewCatalogHandler(r *repository.ProductsRepository, logger *zap.SugaredLogger) *CatalogHandler {
	return &CatalogHandler{
		repo:   r,
		logger: logger,
	}
}

func (h *CatalogHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filters := h.parseFilters(r)

	res, err := h.repo.GetAllProducts(ctx, filters)
	if err != nil {
		h.logger.Error("failed to get products: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	products := make([]Product, len(res))
	for i, p := range res {
		products[i] = Product{
			Code:  p.Code,
			Price: p.Price.InexactFloat64(),
		}
	}

	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Products: products,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode response: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *CatalogHandler) parseFilters(r *http.Request) models.Filters {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			if l < 1 {
				limit = 1
			} else if l > 100 {
				limit = 100
			} else {
				limit = l
			}
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			if o >= 0 {
				offset = o
			}
		}
	}

	return models.Filters{
		Limit:  limit,
		Offset: offset,
	}
}
