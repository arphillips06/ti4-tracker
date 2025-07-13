package helpers

import (
	"errors"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

//to use when game end needs checking

func IsGameFinished(gameID uint) (bool, error) {
	var game models.Game
	if err := database.DB.Select("finished_at").First(&game, gameID).Error; err != nil {
		return false, err
	}
	return game.FinishedAt != nil, nil
}

func GetCurrentRoundID(gameID uint) (uint, error) {
	var game models.Game
	if err := database.DB.Select("id, current_round").First(&game, gameID).Error; err != nil {
		return 0, err
	}

	var round models.Round
	if err := database.DB.
		Where("game_id = ? AND number = ?", game.ID, game.CurrentRound).
		First(&round).Error; err != nil {
		return 0, errors.New("current round not found")
	}

	return round.ID, nil
}
