package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

func serveObjectives(objType string) (int, any, error) {
	objs, err := services.GetObjectivesByType(objType)
	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": "Failed to load " + objType + " objectives"}, nil
	}
	return http.StatusOK, objs, nil
}

// GetAllSecretObjectives godoc
// @Summary      List secret objectives
// @Tags         objectives
// @Produce      json
// @Success      200  {array}   map[string]interface{}
// @Failure      500  {object}  map[string]string  "error"
// @Router       /objectives/secret [get]
func GetAllSecretObjectives(c *gin.Context) (int, any, error) {
	return serveObjectives("Secret")
}

// GetAllPublicObjectives godoc
// @Summary      List public objectives
// @Tags         objectives
// @Produce      json
// @Success      200  {array}   map[string]interface{}
// @Failure      500  {object}  map[string]string  "error"
// @Router       /objectives/public [get]

func GetAllPublicObjectives(c *gin.Context) (int, any, error) {
	return serveObjectives("Public")
}
