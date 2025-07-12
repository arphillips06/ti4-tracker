package controllers

import (
	"net/http"
	"strings"

	"github.com/arphillips06/TI4-stats/database"
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

	game, err := services.GetGameAndRounds(input.GameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}
	if game.FinishedAt != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Game is already finished"})
		return
	}

	var objective models.Objective
	if err := database.DB.First(&objective, input.ObjectiveID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Objective not found"})
		return
	}

	var round models.Round
	if err := database.DB.Where("game_id = ? AND number = ?", game.ID, game.CurrentRound).First(&round).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Current round not found"})
		return
	}

	if err := services.ValidateSecretScoringRules(input.GameID, input.PlayerID, round.ID, objective.ID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	exists, err := services.CheckIfScoreExists(game.ID, input.PlayerID, objective.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Objective already scored by this player"})
		return
	}

	score := models.Score{
		GameID:      game.ID,
		PlayerID:    input.PlayerID,
		ObjectiveID: objective.ID,
		Points:      objective.Points,
		RoundID:     round.ID,
		Type:        strings.ToLower(objective.Type),
	}
	if err := database.DB.Create(&score).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add score"})
		return
	}

	var totalPoints int
	database.DB.Model(&models.Score{}).
		Where("game_id = ? AND player_id = ?", game.ID, input.PlayerID).
		Select("SUM(points)").Scan(&totalPoints)

	if err := services.MaybeFinishGameFromScore(game, input.PlayerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if game.FinishedAt != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Game finished", "winner": game.WinnerID})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Score added",
		"objective":    objective.Name,
		"points":       objective.Points,
		"round":        game.CurrentRound,
		"total_points": totalPoints,
	})
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

	if err := services.ScoreImperialPoint(input.GameID, input.RoundID, input.PlayerID); err != nil {
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
