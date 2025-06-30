package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/database/factions"
	"github.com/arphillips06/TI4-stats/models"

	"github.com/gin-gonic/gin"
)

type PlayerInput struct {
	ID      string
	Name    string
	Faction string
}

type CreateGameInput struct {
	WinningPoints int
	Players       []PlayerInput
}

type selectedPlayersWithFaction struct {
	Player  models.Player
	Faction string
}

// creates a new game and assigns players to it using names or IDs
// also assigns factions
func CreateGame(c *gin.Context) {
	var input CreateGameInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.WinningPoints != 10 && input.WinningPoints != 14 {
		input.WinningPoints = 10
	}
	var allPlayers []models.Player
	_ = database.DB.Find(&allPlayers)

	playerMap := make(map[string]models.Player)
	for _, p := range allPlayers {
		playerMap[strconv.Itoa(int(p.ID))] = p
		playerMap[strings.ToLower(p.Name)] = p
	}

	var selected []selectedPlayersWithFaction

	for _, p := range input.Players {
		lookup := strings.ToLower(p.Name)
		if p.ID != "" {
			lookup = p.ID
		}

		player, exists := playerMap[lookup]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("player not found")})
			return
		}

		if !factions.IsValidFaction(p.Faction) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid faction: %s", p.Faction)})
			return
		}

		selected = append(selected, selectedPlayersWithFaction{
			Player:  player,
			Faction: p.Faction,
		})
	}

	game := models.Game{WinningPoints: input.WinningPoints}
	if err := database.DB.Create(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, entry := range selected {
		gamePlayer := models.GamePlayer{
			GameID:   game.ID,
			PlayerID: entry.Player.ID,
			Faction:  entry.Faction,
		}

		if err := database.DB.Create(&gamePlayer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, game)

}

// advances the round counter and creates a record
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

	if err := database.DB.Save(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "round_advanded", "current_round": game.CurrentRound})
}

// list all games
func ListGames(c *gin.Context) {
	var games []models.Game
	if err := database.DB.Preload("GamePlayers").Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, games)
}
