package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/gin-gonic/gin"
)

func GetAllSecretObjectives(c *gin.Context) {
	var secrets []models.Objective
	if err := database.DB.Where("type = ?", "Secret").Find(&secrets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load secret objectives"})
		return
	}
	c.JSON(http.StatusOK, secrets)
}
