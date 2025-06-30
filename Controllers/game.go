package controllers

import (
	"net/http"
	"strconv"

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

func AdvanceRound(c *gin.Context) {
	gameIDstr := c.Param("game_id")
	gameID, err := strconv.ParseUint(gameIDstr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game_id"})
		return
	}

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}
	game.CurrentRound += 1

	round := models.Round{
		GameID: game.ID,
		Number: game.CurrentRound,
	}

	if err := database.DB.Create(&round).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "round_advanded", "current_round": game.CurrentRound})
}
