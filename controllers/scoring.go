package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/arphillips06/TI4-stats/services"

	"github.com/gin-gonic/gin"
)

// AddScore godoc
// @Summary      Score an objective
// @Description  Marks a player as having scored a specific objective in a game.
// @Tags         scoring
// @Accept       json
// @Produce      json
// @Param        game_id     path      string  true  "Game ID"
// @Param        body        body      object  true  "game_id, player_id, objective_id"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string  "error"
// @Failure      403  {object}  map[string]string  "error"
// @Failure      404  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /games/{game_id}/score [post]
func AddScore(c *gin.Context) (int, any, error) {
	var input struct {
		GameID      uint `json:"game_id"`
		PlayerID    uint `json:"player_id"`
		ObjectiveID uint `json:"objective_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		return http.StatusBadRequest, gin.H{"error": err.Error()}, nil
	}

	resp, err := services.SubmitScore(input.GameID, input.PlayerID, input.ObjectiveID)
	if err != nil {
		switch err.Error() {
		case "game not found", "objective not found", "current round not found":
			return http.StatusNotFound, gin.H{"error": err.Error()}, nil
		case "game is already finished", "objective already scored by this player":
			return http.StatusForbidden, gin.H{"error": err.Error()}, nil
		default:
			return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
		}
	}

	return http.StatusOK, resp, nil
}

// ScoreImperialPoint godoc
// @Summary      Score Imperial point
// @Tags         scoring
// @Accept       json
// @Produce      json
// @Param        body  body      object  true  "game_id, player_id, round_id"
// @Success      204
// @Failure      400  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /score/imperial [post]
func ScoreImperialPoint(c *gin.Context) (int, any, error) {
	var input struct {
		GameID   uint `json:"game_id"`
		PlayerID uint `json:"player_id"`
		RoundID  uint `json:"round_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		return http.StatusBadRequest, gin.H{"error": err.Error()}, nil
	}
	if err := services.ScoreImperialPoint(input.GameID, input.PlayerID); err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
	}
	return http.StatusNoContent, nil, nil
}

// ScoreMecatolPoint godoc
// @Summary      Score Custodians (Mecatol) point
// @Tags         scoring
// @Accept       json
// @Produce      json
// @Param        body  body      object  true  "game_id, player_id, round_id"
// @Success      204
// @Failure      400  {object}  map[string]string  "error"
// @Failure      409  {object}  map[string]string  "error"
// @Router       /score/mecatol [post]
func ScoreMecatolPoint(c *gin.Context) (int, any, error) {
	var input struct {
		GameID   uint `json:"game_id"`
		PlayerID uint `json:"player_id"`
		RoundID  uint `json:"round_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		return http.StatusBadRequest, gin.H{"error": err.Error()}, nil
	}
	if err := services.ScoreMecatolPoint(input.GameID, input.PlayerID); err != nil {
		return http.StatusConflict, gin.H{"error": err.Error()}, nil
	}
	return http.StatusNoContent, nil, nil
}

// DeleteScore godoc
// @Summary      Delete a scored objective
// @Tags         scoring
// @Accept       json
// @Produce      json
// @Param        body  body      object  true  "game_id, player_id, objective_id"
// @Success      204
// @Failure      400  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /score [delete]
func DeleteScore(c *gin.Context) (int, any, error) {
	var req struct {
		GameID      int `json:"game_id"`
		PlayerID    int `json:"player_id"`
		ObjectiveID int `json:"objective_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		return http.StatusBadRequest, gin.H{"error": err.Error()}, nil
	}
	if err := services.RemoveScore(req.GameID, req.PlayerID, req.ObjectiveID); err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
	}
	return http.StatusNoContent, nil, nil
}

// CreatePlayer godoc
// @Summary      Create player
// @Tags         players
// @Accept       json
// @Produce      json
// @Param        body  body      models.Player  true  "Player (name required)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /players [post]
func CreatePlayer(c *gin.Context) (int, any, error) {
	input, ok := helpers.BindJSON[models.Player](c)
	if !ok || strings.TrimSpace(input.Name) == "" {
		return http.StatusBadRequest, gin.H{"error": "name is required"}, nil
	}
	player, err := services.CreatePlayer(input.Name)
	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
	}
	return http.StatusOK, player, nil
}

// AssignPlayerToGame godoc
// @Summary      Assign player to game
// @Tags         players,games
// @Accept       json
// @Produce      json
// @Param        body  body      models.AssignPlayerInput  true  "Game ID, Player ID, Faction"
// @Success      200  {object}  map[string]interface{}      "game_player"
// @Failure      400  {object}  map[string]string           "error"
// @Failure      500  {object}  map[string]string           "error"
// @Router       /games/assign-player [post]
func AssignPlayerToGame(c *gin.Context) (int, any, error) {
	input, ok := helpers.BindJSON[models.AssignPlayerInput](c)
	if !ok {
		return http.StatusBadRequest, gin.H{"error": "invalid payload"}, nil
	}
	gp, err := services.AssignPlayerToGame(input.GameID, input.PlayerID, input.Faction)
	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
	}
	return http.StatusOK, gp, nil
}

// SFTT godoc
// @Summary      Support for the Throne
// @Description  Give or revoke Support for the Throne for a player in a game.
// @Tags         scoring
// @Accept       json
// @Produce      json
// @Param        game_id   path      string  true  "Game ID"
// @Param        player_id path      string  true  "Player ID"
// @Param        body      body      object  true  "round_id, action (give|revoke)"
// @Success      200
// @Failure      400  {object}  map[string]string  "error"
// @Failure      404  {object}  map[string]string  "error"
// @Router       /games/{game_id}/players/{player_id}/sftt [post]
func SFTT(c *gin.Context) (int, any, error) {
	gameID, _ := strconv.ParseUint(c.Param("game_id"), 10, 64)
	playerID, _ := strconv.ParseUint(c.Param("player_id"), 10, 64)

	type SFTTRequest struct {
		RoundID uint   `json:"round_id"`
		Action  string `json:"action"`
	}
	req, ok := helpers.BindJSON[SFTTRequest](c)
	if !ok {
		return http.StatusBadRequest, gin.H{"error": "invalid payload"}, nil
	}

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		return http.StatusNotFound, gin.H{"error": "Game not found"}, nil
	}

	if err := services.HandleSupportForTheThrone(uint(gameID), uint(playerID), req.Action); err != nil {
		return http.StatusBadRequest, gin.H{"error": err.Error()}, nil
	}

	return http.StatusOK, nil, nil
}

// GetScoreSummary godoc
// @Summary      Player score summary
// @Tags         scoring,players
// @Param        id   path      string  true  "Player ID"
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]string  "error"
// @Router       /players/{id}/scores/summary [get]
func GetScoreSummary(c *gin.Context) (int, any, error) {
	id := c.Param("id")
	summary, err := services.GetScoreSummaryByPlayer(id)
	if err != nil {
		return http.StatusNotFound, gin.H{"error": err.Error()}, nil
	}
	return http.StatusOK, summary, nil
}

// GetScoresByRound godoc
// @Summary      Player scores by round
// @Tags         scoring,players
// @Param        id   path      string  true  "Player ID"
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string  "error"
// @Router       /players/{id}/scores/by-round [get]
func GetScoresByRound(c *gin.Context) (int, any, error) {
	id := c.Param("id")
	groupedScores, err := services.GetScoresGroupedByRound(id)
	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": "Could not load scores"}, nil
	}
	return http.StatusOK, groupedScores, nil
}

// GetObjectiveScoreSummary godoc
// @Summary      Objective score summary for a game
// @Tags         scoring,games
// @Param        id   path      int  true  "Game ID"
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /games/{id}/scores/objectives/summary [get]
func GetObjectiveScoreSummary(c *gin.Context) (int, any, error) {
	gameID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return http.StatusBadRequest, gin.H{"error": "Invalid game ID"}, nil
	}
	summary, err := services.GetObjectiveScoreSummary(uint(gameID))
	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
	}
	return http.StatusOK, summary, nil
}

// ScoreImperialRiderPoint godoc
// @Summary      Score Imperial Rider point
// @Tags         scoring
// @Accept       json
// @Produce      json
// @Param        body  body      object  true  "game_id, round_id, player_id"
// @Success      204
// @Failure      400  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /score/imperial-rider [post]
func ScoreImperialRiderPoint(c *gin.Context) (int, any, error) {
	var input struct {
		GameID   uint `json:"game_id"`
		PlayerID uint `json:"player_id"`
		RoundID  uint `json:"round_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		return http.StatusBadRequest, gin.H{"error": err.Error()}, nil
	}
	if err := services.ScoreImperialRiderPoint(input.GameID, input.RoundID, input.PlayerID); err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
	}
	return http.StatusNoContent, nil, nil
}
