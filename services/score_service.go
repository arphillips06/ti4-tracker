package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
)

func ScoreMecatolPoint(gameID, playerID uint) error {
	roundID, err := helpers.GetCurrentRoundID(gameID)
	if err != nil {
		return err
	}

	var existing models.Score
	err = database.DB.
		Where("game_id = ? AND type = ?", gameID, "mecatol").
		First(&existing).Error
	if err == nil {
		return fmt.Errorf("Mecatol Rex point already awarded")
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		return err
	}

	score := models.Score{
		GameID:   gameID,
		RoundID:  roundID,
		PlayerID: playerID,
		Points:   1,
		Type:     "mecatol",
	}
	if err := database.DB.Create(&score).Error; err != nil {
		return err
	}

	return MaybeFinishGameFromScore(&game, playerID)
}

func CheckIfScoreExists(gameID, playerID, objectiveID uint) (bool, error) {
	var existing models.Score
	err := database.DB.
		Where("game_id = ? AND player_id = ? AND objective_id = ?", gameID, playerID, objectiveID).
		First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check for existing score: %w", err)
	}

	return true, nil
}

func WinnerByScore(game *models.Game) error {
	var topScore struct {
		PlayerID uint
		Points   int
	}

	err := database.DB.
		Table("scores").
		Select("player_id, SUM(points) as points").
		Where("game_id = ?", game.ID).
		Group("player_id").
		Order("points DESC").
		Limit(1).
		Scan(&topScore).Error

	if err != nil {
		return err
	}

	if topScore.PlayerID != 0 {
		game.WinnerID = &topScore.PlayerID
		return nil
	}

	return nil // No winner yet
}

func MaybeFinishGameFromExhaustion(game *models.Game) error {
	now := time.Now()
	game.FinishedAt = &now

	if err := WinnerByScore(game); err != nil {
		return err
	}

	return database.DB.Save(game).Error
}

func GetGameAndRounds(gameID uint) (*models.Game, error) {
	var game models.Game
	if err := database.DB.Preload("Rounds").First(&game, gameID).Error; err != nil {
		return nil, err
	}
	return &game, nil
}

func MaybeFinishGameFromScore(game *models.Game, scoringPlayerID uint) error {
	var totalPoints int
	err := database.DB.Model(&models.Score{}).
		Where("game_id = ? AND player_id = ?", game.ID, scoringPlayerID).
		Select("SUM(points)").Scan(&totalPoints).Error
	if err != nil {
		return err
	}

	if totalPoints >= game.WinningPoints {
		now := time.Now()
		game.FinishedAt = &now
		game.WinnerID = &scoringPlayerID
		return database.DB.Save(game).Error
	}

	return nil
}

func RemoveScore(gameID, playerID, objectiveID int) error {
	return database.DB.
		Table("scores").
		Where("game_id = ? AND player_id = ? AND objective_id = ?", gameID, playerID, objectiveID).
		Delete(nil).Error
}
