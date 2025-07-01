package controllers

import (
	"net/http"
	"strconv"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

// GET /games/:game_id/scores
// Description: Returns a grouped list of round-by-round scores for a game, including player names and objectives scored.
func GetScoreSummary(c *gin.Context) {
	gameID := c.Param("id")

	var summaries []models.PlayerScoreSummary

	rows, err := database.DB.
		Table("scores").
		Select("players.id as player_id, players.name as player_name, SUM(scores.points) as points").
		Joins("JOIN players ON scores.player_id = players.id").
		Where("scores.game_id = ?", gameID).
		Group("players.id").
		Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not calculate scores"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var s models.PlayerScoreSummary
		database.DB.ScanRows(rows, &s)
		summaries = append(summaries, s)
	}

	c.JSON(http.StatusOK, summaries)
}

func GetScoresByRound(c *gin.Context) {
	gameID := c.Param("id")
	type rawScore struct {
		RoundNumber int
		PlayerName  string
		Objective   string
		Points      int
	}

	var results []rawScore

	err := database.DB.
		Table("scores").
		Select("rounds.number as round_number, players.name as player_name, objectives.name as objective, scores.points").
		Joins("JOIN players ON scores.player_id = players.id").
		Joins("JOIN objectives ON scores.objective_id = objectives.id").
		Joins("JOIN rounds ON scores.round_id = rounds.id").
		Where("scores.game_id = ?", gameID).
		Order("rounds.number, players.name").
		Scan(&results).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load scores"})
		return
	}

	grouped := make(map[int][]models.RoundScore)

	for _, r := range results {
		grouped[r.RoundNumber] = append(grouped[r.RoundNumber], models.RoundScore{
			Player:    r.PlayerName,
			Objective: r.Objective,
			Points:    r.Points,
		})
	}

	var response []models.RoundScoresGroup
	for round, scores := range grouped {
		response = append(response, models.RoundScoresGroup{
			Round:  round,
			Scores: scores,
		})
	}
	c.JSON(http.StatusOK, response)
}

func GetObjectiveScoreSummary(c *gin.Context) {
	gameIDStr := c.Param("id")
	gameID, err := strconv.Atoi(gameIDStr)
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
