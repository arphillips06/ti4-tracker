package helpers

import (
	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

func DeleteGame(gameID uint) error {
	tx := database.DB.Begin()

	// Delete children in correct order
	if err := tx.Where("game_id = ?", gameID).Delete(&models.Score{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("game_id = ?", gameID).Delete(&models.GameObjective{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("game_id = ?", gameID).Delete(&models.GamePlayer{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("game_id = ?", gameID).Delete(&models.Round{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("game_id = ?", gameID).Delete(&models.SpeakerAssignment{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("game_id = ?", gameID).Delete(&models.PlayerAchievement{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Finally delete the game
	if err := tx.Delete(&models.Game{}, gameID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
