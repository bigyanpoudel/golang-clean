package response

import (
	"math"

	"github.com/gin-gonic/gin"
)

func JSONWithPagination(c *gin.Context, statusCode int, response map[string]interface{}) {
	limit, _ := c.MustGet("limit").(int64)
	size, _ := c.MustGet("page").(int64)
	pagination := gin.H{"page": size, "has_next": (response["count"].(int64) - limit*size) > 0, "count": response["count"], "last_page": math.Ceil(float64(response["count"].(int64)) / float64(size))}
	c.JSON(
		statusCode,
		gin.H{
			"data":       response["data"],
			"pagination": pagination,
		})
}
