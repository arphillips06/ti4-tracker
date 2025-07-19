package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

func ServeObjectives(c *gin.Context, objType string) {
	objs, err := services.GetObjectivesByType(objType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load " + objType + " objectives"})
		return
	}
	c.JSON(http.StatusOK, objs)
}

func GetAllSecretObjectives(c *gin.Context) {
	ServeObjectives(c, "Secret")
}

func GetAllPublicObjectives(c *gin.Context) {
	ServeObjectives(c, "Public")
}
