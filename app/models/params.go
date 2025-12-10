package models

import "github.com/shopspring/decimal"

type QueryParams struct {
	Pagination Pagination
	Filters    Filters
}

type Pagination struct {
	Limit  int
	Offset int
	Total  int64
}

type Filters struct {
	CategoryCode  string
	PriceLessThan *decimal.Decimal
}
