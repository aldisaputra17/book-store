package entities

import "math"

type Pagination struct {
	TotalRecords int `json:"total_records"`
	TotalPages   int `json:"total_pages"`
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
}

func CalculatePagination(totalRecords, currentPage, pageSize int) Pagination {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return Pagination{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  currentPage,
		PageSize:     pageSize,
	}
}

func CalculateOffset(currentPage, pageSize int) int {
	return (currentPage - 1) * pageSize
}
