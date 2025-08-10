package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

// ListPlayersInGame godoc
// @Summary      List players in a game
// @Tags         players
// @Param        id   path      string  true  "Game ID"
// @Produce      json
// @Success      200  {array}   map[string]interface{}
// @Failure      500  {object}  map[string]string  "error"
// @Router       /games/{id}/players [get]
func ListPlayersInGame(c *gin.Context) (int, any, error) {
	gameID := c.Param("id")
	players, err := services.GetPlayersInGame(gameID)
	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
	}
	return http.StatusOK, players, nil
}

// GetPlayerGames godoc
// @Summary      Get a player's games
// @Tags         players
// @Param        id   path      string  true  "Player ID"
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "player, games"
// @Failure      404  {object}  map[string]string       "error"
// @Router       /players/{id}/games [get]
func GetPlayerGames(c *gin.Context) (int, any, error) {
	playerID := c.Param("id")
	player, err := services.GetGamesForPlayer(playerID)
	if err != nil {
		return http.StatusNotFound, gin.H{"error": "Player not found"}, nil
	}
	return http.StatusOK, gin.H{"player": player.Name, "games": player.Games}, nil
}

// ListPlayers godoc
// @Summary      List players
// @Tags         players
// @Produce      json
// @Success      200  {array}   map[string]interface{}
// @Failure      500  {object}  map[string]string  "error"
// @Router       /players [get]
func ListPlayers(c *gin.Context) (int, any, error) {
	players, err := services.ListAllPlayers()
	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
	}
	return http.StatusOK, players, nil
}
