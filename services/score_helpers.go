package services

import (
	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

func GetGameAndRounds(gameID uint) (*models.Game, error) {
	var game models.Game
	if err := database.DB.Preload("Rounds").First(&game, gameID).Error; err != nil {
		return nil, err
	}
	return &game, nil
}

func RemoveScore(gameID, playerID, objectiveID int) error {
	return database.DB.
		Table("scores").
		Where("game_id = ? AND player_id = ? AND objective_id = ?", gameID, playerID, objectiveID).
		Delete(nil).Error
}
