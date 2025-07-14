package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

// POST /players
// Creates a new player.
// Ensures the name is provided and not already in use.
func CreatePlayer(c *gin.Context) {
	input, ok := helpers.BindJSON[models.Player](c)
	if !ok || strings.TrimSpace(input.Name) == "" {
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
	input, ok := helpers.BindJSON[models.AssignPlayerInput](c)
	if !ok {
		return
	}

	gp, err := services.AssignPlayerToGame(input.GameID, input.PlayerID, input.Faction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gp)
}

func SFTT(c *gin.Context) {
	gameID, _ := strconv.ParseUint(c.Param("game_id"), 10, 64)
	playerID, _ := strconv.ParseUint(c.Param("player_id"), 10, 64)

	type SFTTRequest struct {
		RoundID uint   `json:"round_id"`
		Action  string `json:"action"`
	}

	req, ok := helpers.BindJSON[SFTTRequest](c)
	if !ok {
		return
	}

	err := services.HandleSupportForTheThrone(uint(gameID), uint(playerID), req.Action)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
