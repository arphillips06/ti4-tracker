package services

import (
	"time"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

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

		err := database.DB.Model(&models.GamePlayer{}).
			Where("game_id = ? AND player_id = ?", game.ID, scoringPlayerID).
			Update("won", true).Error
		if err != nil {
			return err
		}

		return database.DB.Save(game).Error
	}

	return nil
}

func MaybeFinishGameFromExhaustion(game *models.Game) error {
	now := time.Now()
	game.FinishedAt = &now

	if err := WinnerByScore(game); err != nil {
		return err
	}

	return database.DB.Save(game).Error
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

		err := database.DB.Model(&models.GamePlayer{}).
			Where("game_id = ? AND player_id = ?", game.ID, topScore.PlayerID).
			Update("won", true).Error
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}
