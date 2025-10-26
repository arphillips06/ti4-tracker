package services

import (
	"time"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
)

func ApplyShardOfTheThrone(gameID uint, newHolderID uint) error {
	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		return err
	}

	var lastShardScore models.Score
	err := database.DB.
		Where("game_id = ? AND type = ? AND relic_title = ?", gameID, "relic", "Shard of the Throne").
		Order("created_at desc").
		First(&lastShardScore).Error

	if err == nil && lastShardScore.PlayerID != newHolderID {
		if err := helpers.CreateRelicScore(gameID, lastShardScore.PlayerID, -1, "Shard of the Throne"); err != nil {
			return err
		}
	}

	if err := helpers.CreateRelicScore(gameID, newHolderID, 1, "Shard of the Throne"); err != nil {
		return err
	}

	return MaybeFinishGameFromScore(&game, newHolderID)
}

func ApplyCrownOfEmphidia(gameID, playerID uint) error {
	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		return err
	}

	if err := helpers.CreateRelicScore(gameID, playerID, 1, "The Crown of Emphidia"); err != nil {
		return err
	}

	return MaybeFinishGameFromScore(&game, playerID)

}

func ApplyObsidian(gameID, playerID uint) error {
	score := models.Score{
		GameID:     gameID,
		PlayerID:   playerID,
		Points:     0,
		Type:       "relic",
		RelicTitle: "The Obsidian",
		CreatedAt:  time.Now(),
	}

	return database.DB.Create(&score).Error
}

func ApplyBookOfLatvina(gameID, playerID uint) error {
	score := models.Score{
		GameID:     gameID,
		PlayerID:   playerID,
		Points:     1,
		Type:       "relic",
		RelicTitle: "Book Of Latvina",
		CreatedAt:  time.Now(),
	}
	return database.DB.Create(&score).Error
}
