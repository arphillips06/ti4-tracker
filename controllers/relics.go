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

func HandleShardRelic(c *gin.Context) {
	var req ShardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := services.ApplyShardOfTheThrone(req.GameID, req.NewHolderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Shard of the Throne updated"})
}

func HandleCrownRelic(c *gin.Context) {
	var req RelicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := services.ApplyCrownOfEmphidia(req.GameID, req.PlayerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Crown of Emphidia point assigned"})
}

func HandleObsidianRelic(c *gin.Context) {
	var req RelicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := services.ApplyObsidian(req.GameID, req.PlayerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record Obsidian relic use"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "The Obsidian has been granted"})
}
