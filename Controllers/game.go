package controllers

import (
	"net/http"
	"time"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

// CreateGame handles the POST /games endpoint
// It creates a new game with players and optionally generates objectives
func CreateGame(c *gin.Context) {

	var input models.CreateGameInput

	//bind incoming JSON to the CreateGameInput struct
	//the struct is defined in 'models' so it can be reused and keeps this code cleaner
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//sets to true unless otherwise set
	//true means that the code will generate the objectives
	useDecks := true
	if input.UseObjectiveDecks != nil {
		useDecks = *input.UseObjectiveDecks
	}
	//set 10 points as default unless 14 is given
	if input.WinningPoints != 10 && input.WinningPoints != 14 {
		input.WinningPoints = 10
	}
	//basic validation on player names
	//calls a function from 'services'
	selected, err := services.ParseAndValidatePlayers(input.Players)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//creates game and first round and adds it to the database
	game, round1, err := services.CreateGameAndRound(input.WinningPoints, useDecks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//save player & faction combo to the struct 'GamePlayer'
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
	//Reload to make sure UseObjectiveDecks persists correctly from DB
	if err := database.DB.First(&game, game.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload game"})
		return
	}
	// If we're using the internal decks, assign public objectives now
	if game.UseObjectiveDecks {
		if err := services.AssignObjectivesToGame(game, round1); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Return the new game and any initially revealed objectives
	//TO-DO: currently the game will only let you score a 'revealed' objective
	//however when querying to see what obj were assingned to a game it will
	//always show you all of them. Needs fixing
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
// Returns full game data including player names, factions, and total points scored.
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

// POST /games/:game_id/advance-round
// Advances the round and reveals a public objective unless none remain (in which case, ends the game)
func AdvanceRound(c *gin.Context) {
	gameID := c.Param("game_id")

	game, err := services.GetGameByID(gameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}
	if game.FinishedAt != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game already finished"})
		return
	}
	newRound, err := services.CreateNewRound(game)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new round"})
		return
	}

	stage := services.DetermineStageToReveal(game.ID)

	err = services.RevealNextObjective(game.ID, newRound.ID, stage)
	if err != nil {
		now := time.Now()
		game.FinishedAt = &now
		database.DB.Save(&game)
		c.JSON(http.StatusOK, gin.H{
			"message":       "Game Ended",
			"round":         game.CurrentRound,
			"totalRevealed": services.CountRevealedObjectives(game.ID),
		})
		return
	}
	totalRevealed := services.CountRevealedObjectives(game.ID)
	c.JSON(http.StatusOK, gin.H{
		"message":       "round_advanced",
		"current_round": game.CurrentRound,
		"revealed":      stage,
		"totalRevealed": totalRevealed,
	})
}
