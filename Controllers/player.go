package controllers

import (
	"net/http"
	"strings"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/gin-gonic/gin"
)

// create a new player, checks for blank entries and errors for existing
func CreatePlayer(c *gin.Context) {
	var input models.Player
	if err := c.ShouldBindJSON(&input); err != nil || strings.TrimSpace(input.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	var existing models.Player
	if err := database.DB.
		Where("LOWER(name) = ?", strings.ToLower(input.Name)).
		First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Player name taken"})
		return
	}

	player := models.Player{Name: input.Name}
	if err := database.DB.Create(&player).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, player)
}

// links a player to a game manually
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

// returns a list of players in a specific game
func ListPlayersInGame(c *gin.Context) {
	gameID := c.Param("id")
	var gamePlayers []models.GamePlayer

	if err := database.DB.Where("game_id = ?", gameID).Preload("Player").Find(&gamePlayers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gamePlayers)
}

// list all games that a specific player has been in
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

// list all players
func ListPlayers(c *gin.Context) {
	var players []models.Player
	if err := database.DB.Find(&players).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, players)
}
