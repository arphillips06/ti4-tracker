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
	id := c.Param("id")

	_, scores, err := services.GetGameAndScores(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	scoreSummaryMap := make(map[uint]models.PlayerScoreSummary)
	for _, s := range scores {
		summary := scoreSummaryMap[s.PlayerID]
		summary.PlayerID = s.PlayerID
		summary.PlayerName = s.Player.Name
		summary.Points += s.Points
		scoreSummaryMap[s.PlayerID] = summary
	}

	var summaryList []models.PlayerScoreSummary
	for _, s := range scoreSummaryMap {
		summaryList = append(summaryList, s)
	}

	c.JSON(http.StatusOK, summaryList)
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
