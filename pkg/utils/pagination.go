package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// QueryParams menyimpan parameter query dari request
type QueryParams struct {
	Page        int
	Limit       int
	SortBy      string
	SortOrder   string
	MinPrice    int64
	MaxPrice    int64
	ListingType string
	Bedrooms    int
	Bathrooms   int
	Certificate string
	Location    string
	Title       string
}

// PaginationMetadata menyimpan informasi pagination
type PaginationMetadata struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponse adalah wrapper untuk response yang di-paginate
type PaginatedResponse struct {
	Data       interface{}            `json:"data"`
	Pagination PaginationMetadata     `json:"pagination"`
	Filters    map[string]interface{} `json:"filters,omitempty"`
}

// ParseQueryParams parsing query parameters dari Gin context
func ParseQueryParams(c *gin.Context) QueryParams {
	params := QueryParams{
		Page:      1,
		Limit:     10,
		SortBy:    "created_at",
		SortOrder: "desc",
	}

	// Parse pagination
	if page := c.DefaultQuery("page", "1"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			params.Page = p
		}
	}

	if limit := c.DefaultQuery("limit", "10"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			params.Limit = l
		}
	}

	// Parse sorting
	if sortBy := c.Query("sort_by"); sortBy != "" {
		// Validate sort_by to prevent SQL injection
		validSortFields := map[string]bool{
			"id": true, "created_at": true, "price": true,
			"title": true, "bedrooms": true, "bathrooms": true,
		}
		if validSortFields[sortBy] {
			params.SortBy = sortBy
		}
	}

	if sortOrder := c.Query("sort_order"); sortOrder == "asc" || sortOrder == "desc" {
		params.SortOrder = sortOrder
	}

	// Parse filtering
	if minPrice := c.Query("min_price"); minPrice != "" {
		if p, err := strconv.ParseInt(minPrice, 10, 64); err == nil && p >= 0 {
			params.MinPrice = p
		}
	}

	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if p, err := strconv.ParseInt(maxPrice, 10, 64); err == nil && p > 0 {
			params.MaxPrice = p
		}
	}

	if listing := c.Query("listing_type"); listing == "WTS" || listing == "WTR" {
		params.ListingType = listing
	}

	if bedrooms := c.Query("bedrooms"); bedrooms != "" {
		if b, err := strconv.Atoi(bedrooms); err == nil && b >= 0 {
			params.Bedrooms = b
		}
	}

	if bathrooms := c.Query("bathrooms"); bathrooms != "" {
		if b, err := strconv.Atoi(bathrooms); err == nil && b >= 0 {
			params.Bathrooms = b
		}
	}

	if cert := c.Query("certificate"); cert != "" {
		params.Certificate = cert
	}

	if location := c.Query("location"); location != "" {
		params.Location = location
	}

	if title := c.Query("title"); title != "" {
		params.Title = title
	}

	return params
}

// CalculateOffset menghitung offset untuk LIMIT/OFFSET query
func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}

// CalculateTotalPages menghitung total halaman
func CalculateTotalPages(total int64, limit int) int {
	if limit <= 0 {
		return 0
	}
	return int((total + int64(limit) - 1) / int64(limit))
}

// BuildFiltersMap membuat map dari filters yang diaplikasikan (untuk response)
func BuildFiltersMap(params QueryParams) map[string]interface{} {
	filters := make(map[string]interface{})

	if params.MinPrice > 0 {
		filters["min_price"] = params.MinPrice
	}
	if params.MaxPrice > 0 {
		filters["max_price"] = params.MaxPrice
	}
	if params.ListingType != "" {
		filters["listing_type"] = params.ListingType
	}
	if params.Bedrooms > 0 {
		filters["bedrooms"] = params.Bedrooms
	}
	if params.Bathrooms > 0 {
		filters["bathrooms"] = params.Bathrooms
	}
	if params.Certificate != "" {
		filters["certificate"] = params.Certificate
	}
	if params.Location != "" {
		filters["location"] = params.Location
	}
	if params.Title != "" {
		filters["title"] = params.Title
	}

	return filters
}
