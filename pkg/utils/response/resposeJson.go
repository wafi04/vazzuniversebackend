package response

import (
	"github.com/gin-gonic/gin"
)

// Success returns a standard success response
func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
	})
}

func Error(c *gin.Context, respErr *ResponseError) {
	c.JSON(respErr.Status, gin.H{
		"success": false,
		"error":   respErr.Code,
		"message": respErr.Message,
		"details": respErr.Details,
	})
}

// Pagination returns a paginated response
func Pagination(c *gin.Context, statusCode int, data interface{}, total int, page int, limit int) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
		"meta": gin.H{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}
