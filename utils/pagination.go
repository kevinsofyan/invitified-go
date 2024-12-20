package utils

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Offset     int   `json:"offset"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

func GetPagination(c echo.Context) Pagination {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	return Pagination{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

func SetPagination(p *Pagination, total int64) {
	p.Total = total
	p.TotalPages = int((total + int64(p.Limit) - 1) / int64(p.Limit))
}
