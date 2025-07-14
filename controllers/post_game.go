package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

// CreateGame handles the POST /games endpoint
// It creates a new game with players and optionally generates objectives
func CreateGame(c *gin.Context) {

	input, ok := helpers.BindJSON[models.CreateGameInput](c)
	if !ok {
		return
	}
	const (
		DefaultWinningPoints   = 10
		AlternateWinningPoints = 14
	)
	//sets to true unless otherwise set
	//true means that the code will generate the objectives
	useDecks := true
	if input.UseObjectiveDecks != nil {
		useDecks = *input.UseObjectiveDecks
	}
	//set 10 points as default unless 14 is given
	if input.WinningPoints != DefaultWinningPoints && input.WinningPoints != AlternateWinningPoints {
		input.WinningPoints = DefaultWinningPoints
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

	type CreateGameResponse struct {
		Game     models.Game            `json:"game"`
		Revealed []models.GameObjective `json:"revealed"`
	}

	c.JSON(http.StatusOK, CreateGameResponse{
		Game:     game,
		Revealed: revealed,
	})
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

	// After round 9, automatically attempt to finish the game due to round limit
	if game.CurrentRound >= 9 {
		if err := services.MaybeFinishGameFromExhaustion(game); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finish game"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":       "Game Ended",
			"round":         game.CurrentRound,
			"totalRevealed": services.CountRevealedObjectives(game.ID),
			"winner_id":     game.WinnerID,
		})
		return
	}
	newRound, err := services.CreateNewRound(game)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new round"})
		return
	}

	stage := services.DetermineStageToReveal(game.ID)

	_ = services.RevealNextObjective(game.ID, newRound.ID, stage)

	totalRevealed := services.CountRevealedObjectives(game.ID)
	c.JSON(http.StatusOK, gin.H{
		"message":       "round_advanced",
		"current_round": game.CurrentRound,
		"revealed":      stage,
		"totalRevealed": totalRevealed,
	})
}

func AssignObjective(c *gin.Context) {
	req, ok := helpers.BindJSON[models.AssignObjectiveRequest](c)
	if !ok {
		return
	}

	err := services.ManuallyAssignObjective(req.GameID, uint(req.RoundID), req.ObjectiveID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "objective assigned"})
}
