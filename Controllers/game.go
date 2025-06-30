package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/arphillips06/TI4-stats/services"

	"github.com/gin-gonic/gin"
)

type CreateGameInput struct {
	WinningPoints int
	Players       []models.PlayerInput
}

type selectedPlayersWithFaction struct {
	Player  models.Player
	Faction string
}

// creates a new game and assigns players to it using names or IDs
// also assigns factions
func CreateGame(c *gin.Context) {
	var input CreateGameInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.WinningPoints != 10 && input.WinningPoints != 14 {
		input.WinningPoints = 10
	}

	selected, err := services.ParseAndValidatePlayers(input.Players)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game, round1, err := services.CreateGameAndRound(input.WinningPoints)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, entry := range selected {
		gamePlayer := models.GamePlayer{
			GameID:   game.ID,
			PlayerID: entry.Player.ID,
			Faction:  entry.Faction,
		}
		if err := database.DB.Create(&gamePlayer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := services.AssignObjectivesToGame(game, round1); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var revealed []models.GameObjective
	_ = database.DB.
		Preload("Objective").
		Joins("JOIN rounds ON rounds.id = game_objectives.round_id").
		Where("game_objectives.game_id = ?", game.ID).
		Find(&revealed)

	response := gin.H{
		"game":     game,
		"revealed": revealed,
	}
	c.JSON(http.StatusOK, response)
}

// advances the round counter and creates a record
func AdvanceRound(c *gin.Context) {
	gameID := c.Param("game_id")

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}
	if game.FinishedAt != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game already finished"})
		return
	}

	newRound := models.Round{
		GameID: game.ID,
		Number: game.CurrentRound + 1,
	}

	if err := database.DB.Create(&newRound).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new round"})
		return
	}

	game.CurrentRound = newRound.Number
	if err := database.DB.Save(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update game"})
		return
	}

	//Count home many stage 1 objectives have been revealed
	var revealedStage1Count int64
	database.DB.Model(&models.GameObjective{}).
		Where("game_id = ? AND stage = ? AND round_id IS NOT NULL", game.ID, "I").
		Count(&revealedStage1Count)

	//decide what stage to reveal
	stageToReveal := "I"
	if revealedStage1Count >= 5 {
		stageToReveal = "II"
	}

	var unrevealed models.GameObjective
	err := database.DB.
		Where("game_id = ? AND round_id IS NULL AND stage = ?", game.ID, stageToReveal).
		First(&unrevealed).Error

	if err == nil {
		unrevealed.RoundID = newRound.ID
		database.DB.Save(&unrevealed)
	}

	// if err := database.DB.Save(&game).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message":       "round_advanded",
		"current_round": game.CurrentRound,
		"revealed":      stageToReveal,
	})
}

// list all games
func ListGames(c *gin.Context) {
	var games []models.Game
	if err := database.DB.Preload("GamePlayers").Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, games)
}

// list Scores in game
func GetGameByID(c *gin.Context) {
	id := c.Param("id")

	var game models.Game
	if err := database.DB.
		Preload("GamePlayers.Player").
		First(&game, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}
	var scores []PlayerScoreSummary
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
			var s PlayerScoreSummary
			database.DB.ScanRows(rows, &s)
			scores = append(scores, s)
		}
	}
	response := GameDetailResponse{
		ID:            game.ID,
		WinningPoints: game.WinningPoints,
		CurrentRound:  game.CurrentRound,
		FinishedAt:    game.FinishedAt,
		Players:       game.GamePlayers,
		Scores:        scores,
	}

	c.JSON(http.StatusOK, response)
}

func GetGameObjectives(c *gin.Context) {
	gameID := c.Param("id")

	var gameObjectives []models.GameObjective
	err := database.DB.
		Preload("Objective").
		Preload("Round").
		Where("game_id = ?", gameID).
		Find(&gameObjectives).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load objectives for game"})
		return
	}

	c.JSON(http.StatusOK, gameObjectives)
}
