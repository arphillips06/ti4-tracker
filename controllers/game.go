package controllers

import (
	"net/http"
	"strconv"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

// ListGames returns all games, including associated players and winner info.
func ListGames(c *gin.Context) {
	var games []models.Game
	if err := database.DB.
		Preload("GamePlayers.Player").
		Preload("Winner").
		Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, games)
}

// GET /games/:id
// Returns game detail with objective-based scoring breakdown
func GetGameByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := services.BuildGameDetailResponse(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GET /games/:id/objectives
// Returns all public objectives tied to this game, including stage and round info
func GetGameObjectives(c *gin.Context) {
	gameID := c.Param("id")
	objectives, err := services.GetAllPublicObjectivesForGame(gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, objectives)
}

func GetGameExists(c *gin.Context) {
	id := c.Param("id")
	var game models.Game

	if err := database.DB.First(&game, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"exists": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": true})
}

// CreateGame handles the POST /games endpoint
// It creates a new game with players and optionally generates objectives
func CreateGame(c *gin.Context) {
	input, ok := helpers.BindJSON[models.CreateGameInput](c)
	if !ok {
		return
	}

	game, revealed, err := services.CreateNewGameWithPlayers(*input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game":     game,
		"revealed": revealed,
	})
}

// POST /games/:game_id/advance-round
// Advances the round and reveals a public objective unless none remain (in which case, ends the game)
func AdvanceRound(c *gin.Context) {
	gameIDStr := c.Param("game_id")
	gameIDUint, err := strconv.ParseUint(gameIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	response, err := services.AdvanceGameRound(uint(gameIDUint))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "game not found" {
			status = http.StatusNotFound
		} else if err.Error() == "game already finished" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func AssignObjective(c *gin.Context) {
	req, ok := helpers.BindJSON[models.AssignObjectiveRequest](c)
	if !ok {
		return
	}

	err := services.ManuallyAssignObjective(req.GameID, uint(req.RoundID), req.ObjectiveID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "objective assigned"})
}

func RandomiseSpeaker(c *gin.Context) {
	gameIDParam := c.Param("id")
	gameID, err := strconv.ParseUint(gameIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	speaker, err := services.RandomiseSpeaker(uint(gameID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"speaker_id":   speaker.ID,
		"speaker_name": speaker.Name,
	})
}
