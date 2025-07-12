package controllers

import (
	"fmt"
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
		Winner:            &game.Winner,
	}

	c.JSON(http.StatusOK, response)
}

// GET /games/:id/objectives
// Returns all public objectives tied to this game, including stage and round info
func GetGameObjectives(c *gin.Context) {
	gameID := c.Param("id")

	// Step 1: Load normal game objectives
	var gameObjectives []models.GameObjective
	err := database.DB.
		Preload("Objective").
		Preload("Round").
		Where("game_id = ? AND revealed = true", gameID).
		Find(&gameObjectives).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load objectives for game"})
		return
	}
	fmt.Println("Objectives returned:")
	for _, obj := range gameObjectives {
		fmt.Printf("ID: %d | Stage: %s | RoundID: %d | Revealed: %v\n",
			obj.ObjectiveID, obj.Stage, obj.RoundID, obj.Revealed)
	}

	// Step 2: Inject CDL objectives
	var scores []models.Score
	err = database.DB.
		Preload("Objective").
		Where("game_id = ? AND type = ? AND agenda_title = ?", gameID, "agenda", "Classified Document Leaks").
		Find(&scores).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load CDL agenda scores"})
		return
	}

	for _, s := range scores {
		// Avoid duplication if objective already present
		alreadyIncluded := false
		for _, existing := range gameObjectives {
			if existing.ObjectiveID == s.ObjectiveID && existing.IsCDL {
				alreadyIncluded = true
				break
			}
		}
		if alreadyIncluded {
			continue
		}

		// Inject a pseudo-objective
		cdlObj := models.GameObjective{
			ID:          0, // placeholder
			GameID:      s.GameID,
			ObjectiveID: s.ObjectiveID,
			Stage:       s.Objective.Stage,
			Objective:   s.Objective,
			IsCDL:       true,
		}
		gameObjectives = append(gameObjectives, cdlObj)
	}

	// Step 3: Return combined result
	c.JSON(http.StatusOK, gameObjectives)
}
