package utils

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		limit, _ := c.MustGet("limit").(int64)
		page, _ := c.MustGet("page").(int64)

		offset := (page - 1) * limit

		return db.Offset(int(offset)).Limit(int(limit))
	}
}
