package controllers

import (
	"net/http"
	"time"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
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

	// Step 1: Find previous holder
	var lastShardScore models.Score
	err := database.DB.Where("game_id = ? AND type = ? AND relic_title = ?", req.GameID, "relic", "Shard of the Throne").
		Order("created_at desc").
		First(&lastShardScore).Error

	// Step 2: If a previous holder exists and it's different, subtract 1 point
	if err == nil && lastShardScore.PlayerID != req.NewHolderID {
		// Remove point from previous holder
		prevScore := models.Score{
			GameID:     req.GameID,
			PlayerID:   lastShardScore.PlayerID,
			Points:     -1,
			Type:       "relic",
			RelicTitle: "Shard of the Throne",
			CreatedAt:  time.Now(),
		}
		if err := database.DB.Create(&prevScore).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deduct point from previous holder"})
			return
		}
	}

	// Step 3: Give 1 point to new holder
	newScore := models.Score{
		GameID:     req.GameID,
		PlayerID:   req.NewHolderID,
		Points:     1,
		Type:       "relic",
		RelicTitle: "Shard of the Throne",
		CreatedAt:  time.Now(),
	}
	if err := database.DB.Create(&newScore).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign Shard"})
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

	score := models.Score{
		GameID:     req.GameID,
		PlayerID:   req.PlayerID,
		Points:     1,
		Type:       "relic",
		RelicTitle: "The Crown of Emphidia",
		CreatedAt:  time.Now(),
	}

	if err := database.DB.Create(&score).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record Crown of Emphidia score"})
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

	// This relic grants an extra secret slot, but no points.
	score := models.Score{
		GameID:     req.GameID,
		PlayerID:   req.PlayerID,
		Points:     0,
		Type:       "relic",
		RelicTitle: "The Obsidian",
		CreatedAt:  time.Now(),
	}

	if err := database.DB.Create(&score).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record Obsidian relic use"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "The Obsidian has been granted"})
}
