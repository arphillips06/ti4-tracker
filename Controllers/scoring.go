package controllers

import (
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
		Points        int    `json:"points"`
		ObjectiveName string `json:"objective_name"`
		ObjectiveID   uint   `json:"objective_id"`
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

	//load current round for the game
	var round models.Round
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
