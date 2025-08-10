package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

type ShardRequest struct {
	GameID      uint `json:"game_id"`
	NewHolderID uint `json:"new_holder_id"`
}

type RelicRequest struct {
	GameID   uint `json:"game_id"`
	PlayerID uint `json:"player_id"`
}

// HandleShardRelic godoc
// @Summary      Update "Shard of the Throne" holder
// @Description  Transfers Shard; grants a point to new holder and removes from previous holder if applicable.
// @Tags         relics
// @Accept       json
// @Produce      json
// @Param        body  body      controllers.ShardRequest  true  "Game ID and new holder ID"
// @Success      200  {object}  map[string]string  "message"
// @Failure      400  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /relics/shard [post]
func HandleShardRelic(c *gin.Context) (int, any, error) {
	var req ShardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return http.StatusBadRequest, gin.H{"error": "Invalid request"}, nil
	}
	if err := services.ApplyShardOfTheThrone(req.GameID, req.NewHolderID); err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
	}
	return http.StatusOK, gin.H{"message": "Shard of the Throne updated"}, nil
}

// HandleCrownRelic godoc
// @Summary      Apply "Crown of Emphidia"
// @Description  Grants 1 point to the specified player (one-time effect).
// @Tags         relics
// @Accept       json
// @Produce      json
// @Param        body  body      controllers.RelicRequest  true  "Game ID and player ID"
// @Success      200  {object}  map[string]string  "message"
// @Failure      400  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /relics/crown [post]
func HandleCrownRelic(c *gin.Context) (int, any, error) {
	var req RelicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return http.StatusBadRequest, gin.H{"error": "Invalid request"}, nil
	}
	if err := services.ApplyCrownOfEmphidia(req.GameID, req.PlayerID); err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
	}
	return http.StatusOK, gin.H{"message": "Crown of Emphidia point assigned"}, nil
}

// HandleObsidianRelic godoc
// @Summary      Apply "The Obsidian"
// @Description  Allows a player to score one additional secret objective this game.
// @Tags         relics
// @Accept       json
// @Produce      json
// @Param        body  body      controllers.RelicRequest  true  "Game ID and player ID"
// @Success      200  {object}  map[string]string  "message"
// @Failure      400  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /relics/obsidian [post]
func HandleObsidianRelic(c *gin.Context) (int, any, error) {
	var req RelicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return http.StatusBadRequest, gin.H{"error": "Invalid request"}, nil
	}
	if err := services.ApplyObsidian(req.GameID, req.PlayerID); err != nil {
		return http.StatusInternalServerError, gin.H{"error": "Failed to record Obsidian relic use"}, nil
	}
	return http.StatusOK, gin.H{"message": "The Obsidian has been granted"}, nil
}
