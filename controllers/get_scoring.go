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
		Round  int    `json:"round"`
		Player string `json:"player"`
		Source string `json:"source"`
		Points int    `json:"points"`
	}

	var results []rawScore

	err := database.DB.
		Table("scores").
		Select(`
			COALESCE(rounds.number, 0) AS round,
			players.name AS player,
			COALESCE(objectives.name, scores.agenda_title, scores.relic_title,
				CASE
					WHEN scores.type = 'imperial' THEN 'Imperial Point'
					WHEN scores.type = 'mecatol' THEN 'Custodians'
					WHEN scores.type = 'Support' THEN 'Support for the Throne'
					ELSE 'Unknown'
				END
			) AS source,
			scores.points
		`).
		Joins("JOIN players ON scores.player_id = players.id").
		Joins("LEFT JOIN rounds ON scores.round_id = rounds.id").
		Joins("LEFT JOIN objectives ON scores.objective_id = objectives.id").
		Where("scores.game_id = ?", gameID).
		Order("round, players.name").
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load scores"})
		return
	}

	grouped := make(map[int][]rawScore)
	for _, r := range results {
		grouped[r.Round] = append(grouped[r.Round], r)
	}

	type RoundScoresGroup struct {
		Round  int        `json:"round"`
		Scores []rawScore `json:"scores"`
	}

	var response []RoundScoresGroup
	for round, scores := range grouped {
		response = append(response, RoundScoresGroup{
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
