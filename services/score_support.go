package services

import (
	"fmt"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
)

func ScoreSupportPoint(gameID, playerID uint) error {
	roundID, err := helpers.GetCurrentRoundID(gameID)
	if err != nil {
		return err
	}

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		return err
	}
	if game.FinishedAt != nil {
		return fmt.Errorf("game is already finished")
	}

	var playerCount int64
	err = database.DB.
		Model(&models.GamePlayer{}).
		Where("game_id = ?", gameID).
		Count(&playerCount).Error
	if err != nil {
		return err
	}

	var supportCount int64
	err = database.DB.
		Model(&models.Score{}).
		Where("game_id = ? AND type = ?", gameID, "Support").
		Select("COALESCE(SUM(points), 0)").
		Row().
		Scan(&supportCount)

	if err != nil {
		return err
	}

	if supportCount >= playerCount-1 {
		return fmt.Errorf("Support for the Throne can only be scored by %d players in a %d-player game", playerCount-1, playerCount)
	}

	score := models.Score{
		GameID:   gameID,
		RoundID:  roundID,
		PlayerID: playerID,
		Points:   1,
		Type:     "Support",
	}
	if err := database.DB.Create(&score).Error; err != nil {
		return err
	}

	return MaybeFinishGameFromScore(&game, playerID)
}

func LoseOneSupportPoint(gameID, playerID uint) error {
	roundID, err := helpers.GetCurrentRoundID(gameID)
	if err != nil {
		return err
	}

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		return err
	}
	if game.FinishedAt != nil {
		return fmt.Errorf("game is already finished")
	}

	// Create a negative support score record
	score := models.Score{
		GameID:   gameID,
		RoundID:  roundID,
		PlayerID: playerID,
		Points:   -1,
		Type:     "Support",
	}

	if err := database.DB.Create(&score).Error; err != nil {
		return err
	}

	return nil
}
