package controllers

import (
	"net/http"
	"strconv"

	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

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
