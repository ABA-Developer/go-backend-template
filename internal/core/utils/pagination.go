package utils

import (
	"be-dashboard-nba/api/presenter"
	"math"
)

func GeneratePagination(totalItems, currentPage, limit int) presenter.Pagination {
	if limit <= 0 {
		limit = 10 // Safe fallback
	}
	if currentPage <= 0 {
		currentPage = 1
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	hasNext := currentPage < totalPages
	hasPrev := currentPage > 1

	return presenter.Pagination{
		Page:       currentPage,
		PageSize:   limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}
}
