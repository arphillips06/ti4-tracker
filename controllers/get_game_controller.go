package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

// GET /games
// Returns all games with their associated players
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

	game, scores, err := services.GetGameAndScores(id)
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

	response := models.GameDetailResponse{
		ID:                game.ID,
		WinningPoints:     game.WinningPoints,
		CurrentRound:      game.CurrentRound,
		FinishedAt:        game.FinishedAt,
		UseObjectiveDecks: game.UseObjectiveDecks,
		Players:           game.GamePlayers,
		Rounds:            game.Rounds,
		Objectives:        game.GameObjectives,
		Scores:            summaryList,
		AllScores:         scores,
	}

	c.JSON(http.StatusOK, response)
}

// GET /games/:id/objectives
// Returns all public objectives tied to this game, including stage and round info
func GetGameObjectives(c *gin.Context) {
	gameID := c.Param("id")

	var gameObjectives []models.GameObjective
	err := database.DB.
		Preload("Objective").
		Preload("Round").
		Where("game_id = ? AND round_id > 0", gameID).
		Find(&gameObjectives).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load objectives for game"})
		return
	}

	c.JSON(http.StatusOK, gameObjectives)
}
