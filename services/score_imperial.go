package services

import (
	"fmt"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
)

func ScoreImperialPoint(gameID, roundID, playerID uint) error {
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

	score := models.Score{
		GameID:   gameID,
		RoundID:  roundID,
		PlayerID: playerID,
		Points:   1,
		Type:     "imperial",
	}
	if err := database.DB.Create(&score).Error; err != nil {
		return err
	}
	return MaybeFinishGameFromScore(&game, playerID)
}
