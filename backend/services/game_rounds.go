package services

import (
	"errors"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
)

// Creates and advances to a new round
func CreateNewRound(game *models.Game) (*models.Round, error) {
	newRound := models.Round{
		GameID: game.ID,
		Number: game.CurrentRound + 1,
	}
	if err := database.DB.Create(&newRound).Error; err != nil {
		return nil, err
	}
	game.CurrentRound = newRound.Number
	if err := database.DB.Save(&game).Error; err != nil {
		return nil, err
	}
	return &newRound, nil
}

// Determines if we should reveal a Stage I or Stage II objective this round
func DetermineStageToReveal(gameID uint) string {
	var count int64
	database.DB.Model(&models.GameObjective{}).
		Where("game_id = ? AND stage = ? AND round_id > 0", gameID, "I").
		Count(&count)
	if count >= 5 {
		return "II"
	}
	return "I"
}

// Marks the next unrevealed objective of the given stage as revealed in the current round
func RevealNextObjective(gameID, roundID uint, stage string) error {
	var obj models.GameObjective
	err := database.DB.
		Where("game_id = ? AND round_id = 0 AND stage = ? AND revealed = false", gameID, stage).
		First(&obj).Error
	if err != nil {
		return err
	}
	obj.RoundID = roundID
	obj.Revealed = true

	return database.DB.Save(&obj).Error
}

// Counts total number of revealed public objectives for a game
func CountRevealedObjectives(gameID uint) int64 {
	var count int64
	database.DB.Model(&models.GameObjective{}).
		Where("game_id = ? AND round_id > 0", gameID).
		Count(&count)
	return count
}

func AdvanceGameRound(gameID uint) (map[string]interface{}, error) {
	game, err := helpers.GetUnfinishedGame(gameID)
	if err != nil {
		return nil, err
	}

	if game.CurrentRound >= 9 {
		if err := MaybeFinishGameFromExhaustion(game); err != nil {
			return nil, errors.New("failed to finish game")
		}
		return map[string]interface{}{
			"message":       "Game Ended",
			"round":         game.CurrentRound,
			"totalRevealed": CountRevealedObjectives(game.ID),
			"winner_id":     game.WinnerID,
		}, nil
	}

	newRound, err := CreateNewRound(game)
	if err != nil {
		return nil, errors.New("failed to create new round")
	}

	stage := DetermineStageToReveal(game.ID)
	_ = RevealNextObjective(game.ID, newRound.ID, stage)

	return map[string]interface{}{
		"message":       "round_advanced",
		"current_round": game.CurrentRound,
		"revealed":      stage,
		"totalRevealed": CountRevealedObjectives(game.ID),
	}, nil
}
