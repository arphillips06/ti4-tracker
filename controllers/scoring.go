package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/arphillips06/TI4-stats/services"

	"github.com/gin-gonic/gin"
)

// POST /games/:game_id/score
// Description: Submits a score for a player by marking them as having scored a specific objective in a game.
func AddScore(c *gin.Context) {
	var input struct {
		GameID      uint `json:"game_id"`
		PlayerID    uint `json:"player_id"`
		ObjectiveID uint `json:"objective_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := services.SubmitScore(input.GameID, input.PlayerID, input.ObjectiveID)
	if err != nil {
		switch err.Error() {
		case "game not found", "objective not found", "current round not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "game is already finished", "objective already scored by this player":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

func ScoreImperialPoint(c *gin.Context) {
	var input struct {
		GameID   uint `json:"game_id"`
		PlayerID uint `json:"player_id"`
		RoundID  uint `json:"round_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ScoreImperialPoint(input.GameID, input.PlayerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func ScoreMecatolPoint(c *gin.Context) {
	var input struct {
		GameID   uint `json:"game_id"`
		PlayerID uint `json:"player_id"`
		RoundID  uint `json:"round_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ScoreMecatolPoint(input.GameID, input.PlayerID); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func DeleteScore(c *gin.Context) {
	var req struct {
		GameID      int `json:"game_id"`
		PlayerID    int `json:"player_id"`
		ObjectiveID int `json:"objective_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.RemoveScore(req.GameID, req.PlayerID, req.ObjectiveID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status((http.StatusNoContent))
}

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
	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	err := services.HandleSupportForTheThrone(uint(gameID), uint(playerID), req.Action)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// GetScoreSummary handles GET /games/:id/scores
// Returns total points per player in a game.
func GetScoreSummary(c *gin.Context) {
	id := c.Param("id")

	summary, err := services.GetScoreSummaryByPlayer(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetScoresByRound handles GET /games/:id/scores/rounds
// Returns a round-by-round list of scores with source info.
func GetScoresByRound(c *gin.Context) {
	id := c.Param("id")

	groupedScores, err := services.GetScoresGroupedByRound(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load scores"})
		return
	}

	c.JSON(http.StatusOK, groupedScores)
}

// GetObjectiveScoreSummary handles GET /games/:id/scores/objectives
// Returns scoring summary grouped by public and secret objectives.
func GetObjectiveScoreSummary(c *gin.Context) {
	gameID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	summary, err := services.GetObjectiveScoreSummary(uint(gameID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func ScoreImperialRiderPoint(c *gin.Context) {
	var input struct {
		GameID   uint `json:"game_id"`
		PlayerID uint `json:"player_id"`
		RoundID  uint `json:"round_id"` // Optional
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ScoreImperialRiderPoint(input.GameID, input.RoundID, input.PlayerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
