package catalog

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/models"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type Response struct {
	Products   []Product  `json:"products"`
	Pagination Pagination `json:"pagination"`
}

type Product struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type Pagination struct {
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Total  int64 `json:"total"`
}

type CatalogHandler struct {
	repo   ProductsRepository
	logger *zap.SugaredLogger
}

func NewCatalogHandler(r ProductsRepository, logger *zap.SugaredLogger) *CatalogHandler {
	return &CatalogHandler{
		repo:   r,
		logger: logger,
	}
}

func (h *CatalogHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	paginationParams := h.parsePaginationParams(r)
	filtersParams := h.parseFilterParams(r)

	queryParams := models.QueryParams{
		Pagination: paginationParams,
		Filters:    filtersParams,
	}

	res, err := h.repo.GetAllProducts(ctx, &queryParams)
	if err != nil {
		h.logger.Error("failed to get products: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	products := make([]Product, len(res))
	for i, p := range res {
		products[i] = Product{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: p.Category.Name,
		}
	}

	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Products:   products,
		Pagination: Pagination(queryParams.Pagination),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode response: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *CatalogHandler) parsePaginationParams(r *http.Request) models.Pagination {
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

	return models.Pagination{
		Limit:  limit,
		Offset: offset,
	}
}

func (h *CatalogHandler) parseFilterParams(r *http.Request) models.Filters {
	filters := models.Filters{}

	if categoryCode := r.URL.Query().Get("category"); categoryCode != "" {
		filters.CategoryCode = categoryCode
	}

	if priceStr := r.URL.Query().Get("price_less_than"); priceStr != "" {
		if price, err := decimal.NewFromString(priceStr); err == nil && price.IsPositive() {
			filters.PriceLessThan = &price
		}
	}

	return filters
}
