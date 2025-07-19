package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BindJSON is a generic wrapper around ShouldBindJSON that handles errors
func BindJSON[T any](c *gin.Context) (*T, bool) {
	var obj T
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, false
	}
	return &obj, true
}
