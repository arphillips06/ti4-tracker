package handle

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handle maps an error into an HTTP JSON response.
func Handle(c *gin.Context, err error) {
	if err == nil {
		return
	}

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	log.Printf("ERROR: %v", err)

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func ParseID(c *gin.Context, param string) (uint, error) {
	idStr := c.Param(param)
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt <= 0 {
		return 0, errors.New("invalid ID")
	}
	return uint(idInt), nil
}
