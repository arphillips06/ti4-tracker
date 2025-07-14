package controllers

import (
	"net/http"
	"sort"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

// ListGames returns all games, including associated players and winner info.
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
	var custodiansPlayerID *uint
	for _, s := range scores {
		if s.Type == "mecatol" {
			custodiansPlayerID = &s.PlayerID
			break
		}
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
		ID:                 game.ID,
		WinningPoints:      game.WinningPoints,
		CurrentRound:       game.CurrentRound,
		FinishedAt:         game.FinishedAt,
		UseObjectiveDecks:  game.UseObjectiveDecks,
		Players:            game.GamePlayers,
		Rounds:             game.Rounds,
		Objectives:         game.GameObjectives,
		Scores:             summaryList,
		AllScores:          scores,
		Winner:             &game.Winner,
		CustodiansPlayerID: custodiansPlayerID,
	}

	c.JSON(http.StatusOK, response)
}

// GET /games/:id/objectives
// Returns all public objectives tied to this game, including stage and round info
func GetGameObjectives(c *gin.Context) {
	gameID := c.Param("id")
	const (
		ScoreTypeAgenda = "agenda"
		AgendaCDL       = "Classified Document Leaks"
	)

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

	// Step 2: Add CDL-converted secret objectives as pseudo-public objectives
	var scores []models.Score
	err = database.DB.
		Where("game_id = ? AND type = ? AND agenda_title = ?", gameID, ScoreTypeAgenda, AgendaCDL).
		Find(&scores).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load CDL agenda scores"})
		return
	}

	for _, s := range scores {
		// Avoid duplication if objective already present
		if helpers.ContainsCDLObjective(gameObjectives, s.ObjectiveID) {
			continue
		}

		// Load the full Objective from DB
		var fullObj models.Objective
		if err := database.DB.First(&fullObj, s.ObjectiveID).Error; err != nil {
			continue // silently skip invalid/unknown objectives
		}

		// Inject as pseudo-public CDL objective
		cdlObj := models.GameObjective{
			ID:          0,
			GameID:      s.GameID,
			ObjectiveID: s.ObjectiveID,
			Stage:       fullObj.Stage,
			Objective:   fullObj,
			IsCDL:       true,
		}
		gameObjectives = append(gameObjectives, cdlObj)
	}

	// Sort objectives by ID for consistent frontend ordering
	sort.Slice(gameObjectives, func(i, j int) bool {
		return gameObjectives[i].ObjectiveID < gameObjectives[j].ObjectiveID
	})

	// Step 3: Return combined result
	c.JSON(http.StatusOK, gameObjectives)
}
