package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

// GET /games/:id/players
// Returns a list of players and their factions in a specific game.
func ListPlayersInGame(c *gin.Context) {
	gameID := c.Param("id")
	players, err := services.GetPlayersInGame(gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, players)
}

// GET /players/:id/games
// Returns a list of games that a specific player has participated in,
// including which other players were in those games.
func GetPlayerGames(c *gin.Context) {
	playerID := c.Param("id")
	player, err := services.GetGamesForPlayer(playerID)
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
	players, err := services.ListAllPlayers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, players)
}
