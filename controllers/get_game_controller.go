package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/gin-gonic/gin"
)

// GET /games
// Returns all games with their associated players
func ListGames(c *gin.Context) {
	var games []models.Game
	if err := database.DB.Preload("GamePlayers").Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, games)
}

// GET /games/:id
// Returns game detail with objective-based scoring breakdown
func GetGameByID(c *gin.Context) {
	id := c.Param("id")

	var game models.Game
	if err := database.DB.
		Preload("GamePlayers.Player").
		First(&game, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}
	var scores []models.PlayerScoreSummary
	rows, err := database.DB.
		Table("scores").
		Select("players.id as player_id, players.name as player_name, SUM(scores.points) as points").
		Joins("JOIN players ON scores.player_id = players.id").
		Where("scores.game_id = ?", game.ID).
		Group("players.id").
		Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var s models.PlayerScoreSummary
			database.DB.ScanRows(rows, &s)
			scores = append(scores, s)
		}
	}
	response := models.GameDetailResponse{
		ID:            game.ID,
		WinningPoints: game.WinningPoints,
		CurrentRound:  game.CurrentRound,
		FinishedAt:    game.FinishedAt,
		Players:       game.GamePlayers,
		Scores:        scores,
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
