package helpers

import (
	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

func AgendaAlreadyResolved(gameID uint, agendaTitle string) (bool, error) {
	var count int64
	err := database.DB.
		Model(&models.Score{}).
		Where("game_id = ? AND agenda_title = ?", gameID, agendaTitle).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
