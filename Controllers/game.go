package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/gin-gonic/gin"
)

func CreateGame(c *gin.Context) {
	var input struct {
		WinningPoints int `json:"winning_points"`
		//to do - accept inital players
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.WinningPoints != 10 && input.WinningPoints != 14 {
		input.WinningPoints = 10
	}

	game := models.Game{
		WinningPoints: input.WinningPoints,
	}
	if err := database.DB.Create(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, game)
}
