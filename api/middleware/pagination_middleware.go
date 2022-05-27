package middleware

import (
	"go-clean-api/infrastructure"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	logger infrastructure.Logger
}

func NewPagination(logger infrastructure.Logger) Pagination {
	return Pagination{
		logger: logger,
	}
}

func (p Pagination) SetUp() {

}

func (p Pagination) IncludePagination() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page, err := strconv.ParseInt(ctx.Query("page"), 10, 0)
		if err != nil {
			page = 1
		}
		limit, err := strconv.ParseInt(ctx.Query("limit"), 10, 0)
		if err != nil {
			limit = 10
		}
		sort := ctx.Query("sort")
		ctx.Set("page", page)
		ctx.Set("limit", limit)
		ctx.Set("sort", sort)
		ctx.Next()
	}
}
