package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/gin-gonic/gin"
)

func createPlayer(c *gin.Context) {
	var input models.Player
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var player models.Player
	result := database.DB.Where("name = ?", input.Name).First(&player)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			player = input
			if err := database.DB.Create(&player).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, player)
}

func assignObjectiveToPlayer(c *gin.Context) {
	var input struct {
		PlayerName    string `json:"player_name"`
		ObjectiveName string `json:"objective_name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

}
