package v1

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPageLimit = 50
	maxPageLimit     = 500
)

type Pagination struct {
	Limit  int
	Offset int
}

type PaginationMeta struct {
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Total  int64 `json:"total"`
}

type PaginatedResponse[T any] struct {
	Data       []T            `json:"data"`
	Pagination PaginationMeta  `json:"pagination"`
}

func ParsePagination(c *gin.Context) (Pagination, error) {
	pagination := Pagination{
		Limit:  defaultPageLimit,
		Offset: 0,
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return Pagination{}, fmt.Errorf("invalid limit")
		}
		if limit <= 0 {
			limit = defaultPageLimit
		}
		if limit > maxPageLimit {
			return Pagination{}, fmt.Errorf("limit must be %d or less", maxPageLimit)
		}
		pagination.Limit = limit
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return Pagination{}, fmt.Errorf("invalid offset")
		}
		if offset < 0 {
			return Pagination{}, fmt.Errorf("offset must be 0 or greater")
		}
		pagination.Offset = offset
	}

	return pagination, nil
}