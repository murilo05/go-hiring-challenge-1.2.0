package catalog

import (
	"net/http"
	"strconv"

	httpResponse "github.com/mytheresa/go-hiring-challenge/app/api/handler/http"
	"github.com/mytheresa/go-hiring-challenge/app/models"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type ListProductsResponse struct {
	Products   []Product  `json:"products"`
	Pagination Pagination `json:"pagination"`
}

type GetProductResponse struct {
	Product Product `json:"product"`
}

type Product struct {
	Code     string           `json:"code"`
	Price    float64          `json:"price"`
	Category string           `json:"category"`
	Variants []models.Variant `json:"variants,omitempty"`
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

// GetProduct godoc
// @Summary      Get product by code
// @Description  Get detailed information about a specific product including its variants
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        code  path      string  true  "Product code"  example(PROD001)
// @Success      200   {object}  httpResponse.response{data=GetProductResponse}  "Product details"
// @Failure      400   {object}  httpResponse.errorResponse                      "Bad request - product code is required"
// @Failure      404   {object}  httpResponse.errorResponse                      "Product not found"
// @Failure      500   {object}  httpResponse.errorResponse                      "Internal server error"
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
		httpResponse.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch products")
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

	response := ListProductsResponse{
		Products:   products,
		Pagination: Pagination(queryParams.Pagination),
	}

	httpResponse.OKResponse(w, response)
}

// GetProduct godoc
// @Summary      Get product by code
// @Description  Get detailed information about a specific product including its variants
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        code  path      string  true  "Product code"  example(PROD001)
// @Success      200   {object}  httpResponse.response{data=GetProductResponse}  "Product details"
// @Failure      400   {object}  httpResponse.errorResponse                      "Bad request - product code is required"
// @Failure      404   {object}  httpResponse.errorResponse                      "Product not found"
// @Failure      500   {object}  httpResponse.errorResponse                      "Internal server error"
// @Router       /products/{code} [get]
func (h *CatalogHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := r.PathValue("code")
	if code == "" {
		httpResponse.ErrorResponse(w, http.StatusBadRequest, "Product code is required")
		return
	}

	res, err := h.repo.GetProduct(ctx, code)
	if err != nil {
		h.logger.Error("failed to get product by code: ", err)
		httpResponse.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch product")
		return
	}

	if res == nil {
		httpResponse.ErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	response := GetProductResponse{
		Product: Product{
			Code:     res.Code,
			Price:    res.Price.InexactFloat64(),
			Category: res.Category.Name,
			Variants: res.Variants,
		},
	}

	httpResponse.OKResponse(w, response)
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
