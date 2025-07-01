package services

import (
	"fmt"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

func ScoreImperialPoint(gameID, roundID, playerID uint) error {
	score := models.Score{
		GameID:   gameID,
		RoundID:  roundID,
		PlayerID: playerID,
		Points:   1,
		Type:     "imperial",
	}
	return database.DB.Create(&score).Error
}

func ScoreMecatolPoint(gameID, roundID, playerID uint) error {
	var existing models.Score
	err := database.DB.
		Where("game_id = ? AND type = ?", gameID, "mecatol").
		First(&existing).Error
	if err == nil {
		// Mecatol point already awarded
		return fmt.Errorf("Mecatol Rex point already awarded")
	}

	score := models.Score{
		GameID:   gameID,
		RoundID:  roundID,
		PlayerID: playerID,
		Points:   1,
		Type:     "mecatol",
	}
	return database.DB.Create(&score).Error
}
