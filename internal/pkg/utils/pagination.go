package utils

import (
	"github.com/gin-gonic/gin"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
)

func ParsePaginationQuery(c *gin.Context) (model.PaginationRequest, error) {
	var req model.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return req, err
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 1 {
		req.Limit = 10
	}

	return req, nil
}

func BuildPagination(page, limit int, total int64) Pagination {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages < 0 {
		totalPages = 0
	}

	return Pagination{
		Page:       page,
		Limit:      limit,
		TotalData:  total,
		TotalPages: totalPages,
	}
}
