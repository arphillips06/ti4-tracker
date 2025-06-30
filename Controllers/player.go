package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/gin-gonic/gin"
)

func CreatePlayer(c *gin.Context) {
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

func AssignPlayerToGame(c *gin.Context) {
	var input struct {
		GameID   uint   `json:"game_id"`
		PlayerID uint   `json:"player_id"`
		Faction  string `json:"faction"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gp := models.GamePlayer{
		GameID:   input.GameID,
		PlayerID: input.PlayerID,
		Faction:  input.Faction,
	}

	if err := database.DB.Create(&gp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gp)
}

func ListPlayersInGame(c *gin.Context) {
	gameID := c.Param("id")
	var gamePlayers []models.GamePlayer

	if err := database.DB.Where("game_id = ?", gameID).Preload("Player").Find(&gamePlayers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gamePlayers)
}

func GetPlayerGames(c *gin.Context) {
	playerID := c.Param("id")

	var player models.Player

	err := database.DB.
		Preload("Games.Game").                    // Load Game inside each GamePlayer
		Preload("Games.Game.GamePlayers.Player"). // Load all players of that Game
		First(&player, playerID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"player": player.Name,
		"games":  player.Games,
	})
}
