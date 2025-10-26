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

// ListGames godoc
// @Summary      List games
// @Description  Returns all games with players and winner info.
// @Tags         games
// @Produce      json
// @Param        search  query     string  false  "Search query (e.g., 'winner:Alice', 'player:Bob')"
// @Success      200     {array}   models.Game
// @Failure      500     {object}  map[string]string  "error"
// @Router       /games [get]
func ListGames(c *gin.Context) (int, any, error) {
	if s := strings.TrimSpace(c.Query("search")); s != "" {
		var games []models.Game

		if err := listGamesWithSearch(s).
			Preload("GamePlayers.Player").
			Preload("Winner").
			Find(&games).Error; err != nil {
			return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
		}
		return http.StatusOK, games, nil
	}

	var games []models.Game
	if err := database.DB.
		Preload("GamePlayers.Player").
		Preload("Winner").
		Find(&games).Error; err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
	}
	return http.StatusOK, games, nil
}

// GetGameByID godoc
// @Summary      Get game detail
// @Description  Returns detailed game state with objective-based scoring breakdown.
// @Tags         games
// @Param        id   path      string  true  "Game ID"
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]string  "error"
// @Router       /games/{id} [get]
func GetGameByID(c *gin.Context) (int, any, error) {
	id := c.Param("id")
	resp, err := services.BuildGameDetailResponse(id)
	if err != nil {
		return http.StatusNotFound, gin.H{"error": err.Error()}, nil
	}
	return http.StatusOK, resp, nil
}

// GetGameObjectives godoc
// @Summary      Get game public objectives
// @Description  Returns all public objectives for a game, including stage and round.
// @Tags         games
// @Param        id   path      string  true  "Game ID"
// @Produce      json
// @Success      200  {array}   map[string]interface{}
// @Failure      500  {object}  map[string]string  "error"
// @Router       /games/{id}/objectives [get]
func GetGameObjectives(c *gin.Context) (int, any, error) {
	gameID := c.Param("id")
	objectives, err := services.GetAllPublicObjectivesForGame(gameID)
	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
	}
	return http.StatusOK, objectives, nil
}

// GetGameExists godoc
// @Summary      Check game exists
// @Description  Returns {"exists": true|false}.
// @Tags         games
// @Param        id   path      string  true  "Game ID"
// @Produce      json
// @Success      200  {object}  map[string]bool
// @Failure      404  {object}  map[string]bool
// @Router       /games/{id}/exists [get]
func GetGameExists(c *gin.Context) (int, any, error) {
	id := c.Param("id")
	var game models.Game
	if err := database.DB.First(&game, id).Error; err != nil {
		return http.StatusNotFound, gin.H{"exists": false}, nil
	}
	return http.StatusOK, gin.H{"exists": true}, nil
}

// CreateGame godoc
// @Summary      Create a new game
// @Description  Creates a new game with players; can optionally generate objectives.
// @Tags         games
// @Accept       json
// @Produce      json
// @Param        body  body      models.CreateGameInput  true  "New game payload"
// @Success      200  {object}  map[string]interface{}  "game, revealed"
// @Failure      400  {object}  map[string]string       "error"
// @Failure      500  {object}  map[string]string       "error"
// @Router       /games [post]
func CreateGame(c *gin.Context) (int, any, error) {
	input, ok := helpers.BindJSON[models.CreateGameInput](c)
	if !ok {
		return http.StatusBadRequest, gin.H{"error": "invalid payload"}, nil
	}
	game, revealed, err := services.CreateNewGameWithPlayers(*input)
	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
	}
	return http.StatusOK, gin.H{"game": game, "revealed": revealed}, nil
}

// AdvanceRound godoc
// @Summary      Advance game round
// @Description  Advances the round; reveals a public objective unless none remain (then ends the game).
// @Tags         games
// @Param        game_id  path      string  true  "Game ID"
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string  "error"
// @Failure      404  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /games/{game_id}/advance-round [post]
func AdvanceRound(c *gin.Context) (int, any, error) {
	gameIDStr := c.Param("game_id")
	gameIDUint, err := strconv.ParseUint(gameIDStr, 10, 64)
	if err != nil {
		return http.StatusBadRequest, gin.H{"error": "invalid game ID"}, nil
	}
	response, err := services.AdvanceGameRound(uint(gameIDUint))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "game not found" {
			status = http.StatusNotFound
		} else if err.Error() == "game already finished" {
			status = http.StatusBadRequest
		}
		return status, gin.H{"error": err.Error()}, nil
	}
	return http.StatusOK, response, nil
}

// AssignObjective godoc
// @Summary      Manually assign a public objective to a round
// @Description  Admin action to attach an objective to a specific game round.
// @Tags         games
// @Accept       json
// @Produce      json
// @Param        body  body      models.AssignObjectiveRequest  true  "Assignment payload"
// @Success      200  {object}  map[string]string  "message"
// @Failure      400  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /games/assign-objective [post]
func AssignObjective(c *gin.Context) (int, any, error) {
	req, ok := helpers.BindJSON[models.AssignObjectiveRequest](c)
	if !ok {
		return http.StatusBadRequest, gin.H{"error": "invalid payload"}, nil
	}
	if err := services.ManuallyAssignObjective(req.GameID, uint(req.RoundID), req.ObjectiveID); err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
	}
	return http.StatusOK, gin.H{"message": "objective assigned"}, nil
}

// RandomiseSpeaker godoc
// @Summary      Randomise speaker
// @Description  Randomly selects a speaker for the game.
// @Tags         games
// @Param        id   path      string  true  "Game ID"
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "speaker_id, speaker_name"
// @Failure      400  {object}  map[string]string       "error"
// @Failure      500  {object}  map[string]string       "error"
// @Router       /games/{id}/speaker/randomise [post]
func RandomiseSpeaker(c *gin.Context) (int, any, error) {
	gameIDParam := c.Param("id")
	gameID, err := strconv.ParseUint(gameIDParam, 10, 64)
	if err != nil {
		return http.StatusBadRequest, gin.H{"error": "Invalid game ID"}, nil
	}
	speaker, err := services.RandomiseSpeaker(uint(gameID))
	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
	}
	return http.StatusOK, gin.H{"speaker_id": speaker.ID, "speaker_name": speaker.Name}, nil
}

// PostAssignSpeaker godoc
// @Summary      Assign speaker
// @Description  Assigns a speaker (initial or current) for a specific round.
// @Tags         games
// @Accept       json
// @Produce      json
// @Param        game_id  path      string  true  "Game ID"
// @Param        body     body      object  true  "player_id, round_id, is_initial"
// @Success      200  {object}  map[string]string  "message"
// @Failure      400  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /games/{game_id}/speaker [post]
func PostAssignSpeaker(c *gin.Context) (int, any, error) {
	gameID, _ := strconv.Atoi(c.Param("game_id"))
	type AssignSpeakerRequest struct {
		PlayerID  uint `json:"player_id"`
		RoundID   uint `json:"round_id"`
		IsInitial bool `json:"is_initial"`
	}
	req, ok := helpers.BindJSON[AssignSpeakerRequest](c)
	if !ok {
		return http.StatusBadRequest, gin.H{"error": "invalid payload"}, nil
	}
	if err := services.AssignSpeaker(uint(gameID), req.RoundID, req.PlayerID); err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
	}
	return http.StatusOK, gin.H{"message": "Speaker assigned"}, nil
}

// DeleteGameHandler godoc
// @Summary      Delete a game
// @Description  Permanently deletes a game by its ID.
// @Tags         games
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Game ID"
// @Success      200  {object}  map[string]interface{} "status and deleted game ID"
// @Failure      400  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /games/{id} [delete]
func DeleteGameHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
		return
	}

	if err := helpers.DeleteGame(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted", "game_id": id})
}
