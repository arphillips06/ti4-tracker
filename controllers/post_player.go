package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

// POST /players
// Creates a new player.
// Ensures the name is provided and not already in use.
func CreatePlayer(c *gin.Context) {
	var input models.Player
	if err := c.ShouldBindJSON(&input); err != nil || strings.TrimSpace(input.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	player, err := services.CreatePlayer(input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, player)
}

// POST /games/assign-player
// Manually assigns a player to a game with a faction.
// Useful when not using the automated game setup flow.
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

func SFTT(c *gin.Context) {
	gameID, _ := strconv.ParseUint(c.Param("game_id"), 10, 64)
	playerID, _ := strconv.ParseUint(c.Param("player_id"), 10, 64)

	var req struct {
		RoundID uint   `json:"round_id"`
		Action  string `json:"action"` // "score or "unscore"
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request bosy"})
		return
	}

	switch req.Action {
	case "score":
		err := services.ScoreSupportPoint(uint(gameID), uint(playerID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	case "unscore":
		err := services.LoseOneSupportPoint(uint(gameID), uint(playerID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid actions"})
	}
}
