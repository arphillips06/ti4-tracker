package services

import (
	"time"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

func ApplyShardOfTheThrone(gameID uint, newHolderID uint) error {
	// Step 1: Load game
	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		return err
	}

	// Step 2: Find previous holder
	var lastShardScore models.Score
	err := database.DB.
		Where("game_id = ? AND type = ? AND relic_title = ?", gameID, "relic", "Shard of the Throne").
		Order("created_at desc").
		First(&lastShardScore).Error

	if err == nil && lastShardScore.PlayerID != newHolderID {
		deduct := models.Score{
			GameID:     gameID,
			PlayerID:   lastShardScore.PlayerID,
			Points:     -1,
			Type:       "relic",
			RelicTitle: "Shard of the Throne",
			CreatedAt:  time.Now(),
		}
		if err := database.DB.Create(&deduct).Error; err != nil {
			return err
		}
	}

	// Step 3: Add point to new holder
	add := models.Score{
		GameID:     gameID,
		PlayerID:   newHolderID,
		Points:     1,
		Type:       "relic",
		RelicTitle: "Shard of the Throne",
		CreatedAt:  time.Now(),
	}
	if err := database.DB.Create(&add).Error; err != nil {
		return err
	}

	return MaybeFinishGameFromScore(&game, newHolderID)
}

func ApplyCrownOfEmphidia(gameID, playerID uint) error {
	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		return err
	}

	add := models.Score{
		GameID:     gameID,
		PlayerID:   playerID,
		Points:     1,
		Type:       "relic",
		RelicTitle: "The Crown of Emphidia",
		CreatedAt:  time.Now(),
	}
	if err := database.DB.Create(&add).Error; err != nil {
		return err
	}

	return MaybeFinishGameFromScore(&game, playerID)

}

func ApplyObsidian(gameID, playerID uint) error {
	score := models.Score{
		GameID:     gameID,
		PlayerID:   playerID,
		Points:     0, // no point
		Type:       "relic",
		RelicTitle: "The Obsidian",
		CreatedAt:  time.Now(),
	}

	return database.DB.Create(&score).Error
}
