package controllers

import (
	"net/http"

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
