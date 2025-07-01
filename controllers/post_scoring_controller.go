package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/gin-gonic/gin"
)

// POST /games/:game_id/score
// Description: Submits a score for a player by marking them as having scored a specific objective in a game.
func AddScore(c *gin.Context) {
	var input struct {
		GameID        uint   `json:"game_id"`
		PlayerID      uint   `json:"player_id"`
		RoundID       uint   `json:"round_id"`
		ObjectiveName string `json:"objective_name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//load game
	var game models.Game
	if err := database.DB.Preload("Rounds").First(&game, input.GameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game Not Found"})
		return
	}
	//check if the game is finished
	if game.FinishedAt != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Game is already Finished"})
		return
	}

	//load objective by name
	var objective models.Objective
	if err := database.DB.Where("LOWER(name) = ?", strings.ToLower(input.ObjectiveName)).First(&objective).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Objective Not Found"})
		return
	}
	var round models.Round
	if objective.Type == "Secret" {
		var count int64
		err := database.DB.
			Model(&models.Score{}).
			Joins("JOIN objectives ON objectives.id = scores.objective_id").
			Where("scores.player_id = ? AND scores.round_id = ? AND objectives.type = ? AND objectives.phase = ?",
				input.PlayerID, round.ID, "Secret", objective.Phase).
			Count(&count).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate secret scoring rules"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprint("Player has already scored a secret objective during the %s phase this round", objective.Phase),
			})
			return
		}
	}
	//load current round for the game

	if err := database.DB.Where("game_id = ? AND number = ?", game.ID, game.CurrentRound).First(&round).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Current Round Not Found"})
		return
	}

	//add the score
	score := models.Score{
		GameID:      input.GameID,
		PlayerID:    input.PlayerID,
		ObjectiveID: objective.ID,
		Points:      objective.Points,
		RoundID:     round.ID,
	}

	var existing models.Score
	err := database.DB.Where("game_id = ? AND player_id = ? AND objective_id = ?", input.GameID, input.PlayerID, objective.ID).
		First(&existing).Error

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Objective already scored by this player"})
		return
	}

	if err := database.DB.Create(&score).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//sum all points for the player in a single game instance
	var totalPoints int
	database.DB.Model(&models.Score{}).
		Where("game_id = ? AND player_id = ?", input.GameID, input.PlayerID).
		Select("SUM(points)").Scan(&totalPoints)

	//check if points >= winning points
	if totalPoints >= game.WinningPoints {
		now := time.Now()
		game.FinishedAt = &now
		if input.PlayerID != 0 {
			game.WinnerID = *&input.PlayerID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "playerID is nil"})
			return
		}
		if err := database.DB.Save(&game).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Game finished", "winner": input.PlayerID})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "score added",
		"total_points": totalPoints,
		"round":        game.CurrentRound,
		"objective":    objective.Name,
		"points":       objective.Points,
	})
}
