package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/gin-gonic/gin"
)

// GET /games/:id/players
// Returns a list of players and their factions in a specific game.
func ListPlayersInGame(c *gin.Context) {
	gameID := c.Param("id")
	var gamePlayers []models.GamePlayer

	if err := database.DB.Where("game_id = ?", gameID).Preload("Player").Find(&gamePlayers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gamePlayers)
}

// GET /players/:id/games
// Returns a list of games that a specific player has participated in,
// including which other players were in those games.
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

// GET /players
// Returns a list of all players in the system.
func ListPlayers(c *gin.Context) {
	var players []models.Player
	if err := database.DB.Find(&players).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, players)
}
